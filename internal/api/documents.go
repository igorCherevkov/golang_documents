package api

import (
	"documents/internal/api/dto"
	"documents/internal/api/utils"
	"documents/internal/models"
	"documents/internal/store"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) uploadDocument(ctx *fiber.Ctx) error {
	userID := userIDFromCtx(ctx)

	metaRaw := ctx.FormValue("meta")
	if metaRaw == "" {
		return utils.ErrBadRequest(ctx, "meta is required")
	}

	var meta dto.DocMeta
	if err := json.Unmarshal([]byte(metaRaw), &meta); err != nil {
		return utils.ErrBadRequest(ctx, "meta must be valid json")
	}

	if err := s.validate.Struct(meta); err != nil {
		return utils.ErrBadRequest(ctx, utils.ValidationErrorText(err))
	}

	var jsonData []byte
	if raw := ctx.FormValue("json"); raw != "" {
		if !json.Valid([]byte(raw)) {
			return utils.ErrBadRequest(ctx, "json must be valid")
		}
		jsonData = []byte(raw)
	}

	doc := models.Document{
		ID: uuid.New().String(),
		UserID: userID,
		Name: meta.Name,
		Mime: meta.Mime,
		IsFile: meta.File,
		IsPublic: meta.Public,
		JSONData: jsonData,
	}

	var fileName string 
	if meta.File {
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			return utils.ErrBadRequest(ctx, "file is required when meta.file is true")
		}

		if err := os.MkdirAll(s.cfg.StorageDir, 0o755); err != nil {
			return utils.ErrInternal(ctx, "internal error")
		}

		path := filepath.Join(s.cfg.StorageDir, doc.ID + filepath.Ext(meta.Name))
		if err := ctx.SaveFile(fileHeader, path); err != nil {
			return utils.ErrInternal(ctx, "internal error")
		}

		doc.FilePath = path
		fileName = meta.Name
	}

	if _, err := s.documents.Create(ctx.Context(), doc, meta.Grant); err != nil {
		return utils.ErrInternal(ctx, "internal error")
	}

	data := fiber.Map{}
	if jsonData != nil {
		data["json"] = json.RawMessage(jsonData)
	}
	
	if fileName != "" {
		data["file"] = fileName
	}

	return utils.SendData(ctx, fiber.StatusOK, data)
}

func (s *Server) listDocuments(ctx *fiber.Ctx) error {
	userID := userIDFromCtx(ctx)

	var q dto.ListDocsQuery
	if err := ctx.QueryParser(&q); err != nil {
		return utils.ErrBadRequest(ctx, "invalid query")
	}

	if err := s.validate.Struct(q); err != nil {
		return utils.ErrBadRequest(ctx, utils.ValidationErrorText(err))
	}

	params := store.ListParams{
		RequesterUserID: userID,
		UserLogin: q.Login,
		FilterKey: q.Key,
		FilterValue: q.Value,
		Limit: q.Limit,
	}

	docs, err := s.documents.List(ctx.Context(), params)
	if err != nil {
		if errors.Is(err, store.ErrBadFilterKey) {
			return utils.ErrBadRequest(ctx, "invalid filter key")
		}

		return utils.ErrInternal(ctx, "internal error")
	}

	if ctx.Method() == fiber.MethodHead {
		return ctx.SendStatus(fiber.StatusOK)
	}

	response := make([]dto.DocResponse, len(docs))
	for i, doc := range docs {
		response[i] = dto.ToDocResponse(doc)
	}

	return utils.SendData(ctx, fiber.StatusOK, fiber.Map{"docs": response})
}

func (s *Server) getDocument(ctx *fiber.Ctx) error {
	userID := userIDFromCtx(ctx)

	id := ctx.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return utils.ErrBadRequest(ctx, "invalid document id")
	}

	doc, err := s.documents.GetByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return utils.ErrForbidden(ctx, "document not found or access denied")
		}

		return utils.ErrInternal(ctx, "internal error")
	}

	requester, err := s.users.GetUserById(ctx.Context(), userID)
	if err != nil {
		return utils.ErrInternal(ctx, "internal error")
	}

	if !utils.CanAccessDocument(doc, userID, requester.Login) {
		return utils.ErrForbidden(ctx, "document not found or access denied")
	}

	if doc.IsFile {
		if ctx.Method() == fiber.MethodHead {
			ctx.Set(fiber.HeaderContentType, doc.Mime)

			return ctx.SendStatus(fiber.StatusOK)
		}

		ctx.Set(fiber.HeaderContentType, doc.Mime)
		return ctx.SendFile(doc.FilePath, false)
	}

	if ctx.Method() == fiber.MethodHead {
		return ctx.SendStatus(fiber.StatusOK)
	}

	var payload any
	if len(doc.JSONData) > 0 {
		payload = json.RawMessage(doc.JSONData)
	} else {
		payload = fiber.Map{}
	}

	return utils.SendData(ctx, fiber.StatusOK, payload)
}

func (s *Server) deleteDocument(ctx *fiber.Ctx) error {
	userID := userIDFromCtx(ctx)

	id := ctx.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return utils.ErrBadRequest(ctx, "invalid document id")
	}

	doc, err := s.documents.GetByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return utils.ErrForbidden(ctx, "document not found or access denied")
		}

		return utils.ErrInternal(ctx, "internal error")
	}

	if doc.UserID != userID {
		return utils.ErrForbidden(ctx, "only the owner can delete document")
	}

	if err := s.documents.Delete(ctx.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return utils.ErrForbidden(ctx, "document not found or access denied")
		}

		return utils.ErrInternal(ctx, "internal error")
	}

	if doc.IsFile && doc.FilePath != "" {
		_ = os.Remove(doc.FilePath)
	}

	return utils.SendResponse(ctx, fiber.StatusOK, fiber.Map{"success": true})
}