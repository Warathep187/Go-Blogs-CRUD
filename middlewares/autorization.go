package middlewares

import (
	"go_blogs/libs"
	"go_blogs/models"

	"github.com/gofiber/fiber/v2"
)

func AuthorizeUser(c *fiber.Ctx) error {
	userData, err := libs.GetUserSessionData(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Message: fiber.ErrUnauthorized.Message,
		})
	}

	c.Locals("user", userData)

	return c.Next()
}
