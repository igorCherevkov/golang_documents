package api

import (
	"documents/internal/api/utils"
	"documents/internal/store"
	"errors"

	"github.com/gofiber/fiber/v2"
)

const localsUserID = "userID"
const localsToken = "token"

func extractToken(c *fiber.Ctx) string {
	if auth := c.Get("Authorization"); auth != "" {
		if len(auth) > 7 && auth[:7] == "Bearer " {
			return auth[7:]
		}

		return auth
	}

	return ""
}

func requireAuth(tokens *store.TokenStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractToken(c)
		if token == "" {
			return utils.ErrUnauthorized(c, "authorization token is required")
		}

		userID, err := tokens.UserIDByToken(c.Context(), token)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return utils.ErrUnauthorized(c, "invalid or expired token")
			}

			return utils.ErrInternal(c, "internal error")
		}

		c.Locals(localsUserID, userID)
		c.Locals(localsToken, token)

		return c.Next()
	}
}

func userIDFromCtx(ctx *fiber.Ctx) string {
	userID, _ := ctx.Locals(localsUserID).(string)
	return userID
}

func tokenFromCtx(ctx *fiber.Ctx) string {
	token, _ := ctx.Locals(localsToken).(string)

	return token
}