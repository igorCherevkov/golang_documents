package utils

import "github.com/gofiber/fiber/v2"

type ApiError struct {
	Code	int	`json:"code"`
	Text	string	`json:"text"`
}

type Envelope struct {
	Error	*ApiError	`json:"error,omitempty"`
	Response	any	`json:"response,omitempty"`
	Data	any	`json:"data,omitempty"`
}

func SendError(c *fiber.Ctx, status int, text string) error {
	return c.Status(status).JSON(Envelope{
		Error: &ApiError{Code: status, Text: text},
	})
}

func SendResponse(c *fiber.Ctx, status int, response any) error {
	return c.Status(status).JSON(Envelope{Response: response})
}

func SendData(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(Envelope{Data: data})
}

func ErrBadRequest(c *fiber.Ctx, text string) error {
	return SendError(c, fiber.StatusBadRequest, text)
}

func ErrUnauthorized(c *fiber.Ctx, text string) error {
	return SendError(c, fiber.StatusUnauthorized, text)
}

func ErrForbidden(c *fiber.Ctx, text string) error {
	return SendError(c, fiber.StatusForbidden, text)
}

func ErrNotImplemented(c *fiber.Ctx, text string) error {
	return SendError(c, fiber.StatusNotImplemented, text)
}

func ErrInternal(c *fiber.Ctx, text string) error {
	return SendError(c, fiber.StatusInternalServerError, text)
}