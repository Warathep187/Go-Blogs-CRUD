package utils

import "github.com/gofiber/fiber/v2"

func NewAppError(err error) *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
}
