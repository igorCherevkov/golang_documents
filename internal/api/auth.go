package api

import (
	"documents/internal/api/dto"
	"documents/internal/api/utils"
	"documents/internal/store"
	"errors"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) register(ctx *fiber.Ctx) error {
	var req dto.RegisterDto
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrBadRequest(ctx, "invalid body")
	}

	if req.Token != s.cfg.AdminToken {
		return utils.ErrUnauthorized(ctx, "invalid admin token")
	}

	if err := s.validate.Struct(req); err != nil {
		return utils.ErrBadRequest(ctx, utils.ValidationErrorText(err))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrInternal(ctx, "internal error")
	}

	user, err := s.users.Create(ctx.Context(), req.Login, string(hash))
	if err != nil {
		if errors.Is(err, store.ErrLoginTaken) {
			return utils.ErrBadRequest(ctx, "login is already taken")
		}

		return utils.ErrInternal(ctx, "internal error")
	}

	return utils.SendResponse(ctx, fiber.StatusCreated, fiber.Map{"login": user.Login})
}

func (s *Server) auth(ctx *fiber.Ctx) error {
	var req dto.AuthDto
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrBadRequest(ctx, "invalid body")
	}

	if err := s.validate.Struct(req); err != nil {
		return utils.ErrBadRequest(ctx, utils.ValidationErrorText(err))
	}

	user, err := s.users.GetUserByLogin(ctx.Context(), req.Login)
	if err !=  nil {
		if errors.Is(err, store.ErrNotFound) {
			return utils.ErrUnauthorized(ctx, "invalid login or password")
		}

		return utils.ErrUnauthorized(ctx, "internal error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.ErrUnauthorized(ctx, "invalid password")
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return utils.ErrInternal(ctx, "internal error")
	}

	if err := s.tokens.Create(ctx.Context(), token, user.ID, s.cfg.TokenTTL); err != nil {
		return utils.ErrInternal(ctx, "internal error")
	}

	return utils.SendResponse(ctx, fiber.StatusOK, fiber.Map{"token": token})
}

func (s *Server) logout(ctx *fiber.Ctx) error {
	token := tokenFromCtx(ctx)

	if err := s.tokens.DeleteSession(ctx.Context(), token); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return utils.ErrUnauthorized(ctx, "invalid or expired token")
		}

		return utils.ErrInternal(ctx, "internal error")
	}

	return utils.SendResponse(ctx, fiber.StatusOK, fiber.Map{"success": true})
}