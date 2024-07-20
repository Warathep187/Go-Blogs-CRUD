package validators

import (
	"go_blogs/constants"
	"go_blogs/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Body
type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
	Name     string `json:"name" validate:"required"`
}

var validate *validator.Validate = validator.New()

func ValidateAuthPayload(routeName string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var body interface{}

		switch routeName {
		case constants.RouteName.LOGIN:
			body = new(LoginPayload)
		case constants.RouteName.REGISTER:
			body = new(RegisterPayload)
		}

		if err := c.BodyParser(body); err != nil {
			return utils.NewAppError(err)
		}

		errors := validate.Struct(body)

		formattedErrorResponse := utils.TransformValidationErrorFormat(errors)

		if len(formattedErrorResponse) > 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(formattedErrorResponse)
		}

		c.Locals("payload", body)

		return c.Next()
	}
}
