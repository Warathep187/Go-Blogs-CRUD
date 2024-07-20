package validators

import (
	"go_blogs/constants"
	"go_blogs/utils"

	"github.com/gofiber/fiber/v2"
)

// Query
type GetBlogsQuery struct {
	From int `json:"from" validate:"gte=0"`
}

// Params
type GetBlogByIDParams struct {
	ID string `json:"id" validate:"mongodb"`
}

type UpdateBlogParams struct {
	ID string `json:"id" validate:"mongodb"`
}

type DeleteBlogParams struct {
	ID string `json:"id" validate:"mongodb"`
}

// Body
type CreateBlogPayload struct {
	Title   string `json:"title" validate:"required,min=10"`
	Content string `json:"content" validate:"required"`
}

type UpdateBlogPayload struct {
	Title   string `json:"title" validate:"required,min=10"`
	Content string `json:"content" validate:"required"`
}

func ValidateBlogQuery(routeName string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var query interface{}

		switch routeName {
		case constants.RouteName.GET_BLOGS:
			query = new(GetBlogsQuery)
		}

		if err := c.QueryParser(query); err != nil {
			return utils.NewAppError(err)
		}

		errors := validate.Struct(query)

		formattedErrorResponse := utils.TransformValidationErrorFormat(errors)

		if len(formattedErrorResponse) > 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(formattedErrorResponse)
		}

		c.Locals("query", query)

		return c.Next()
	}
}

func ValidateBlogParams(routeName string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var params interface{}

		switch routeName {
		case constants.RouteName.GET_BLOG_BY_ID:
			params = new(GetBlogByIDParams)
		case constants.RouteName.UPDATE_BLOG:
			params = new(UpdateBlogParams)
		case constants.RouteName.DELETE_BLOG:
			params = new(DeleteBlogParams)
		}

		if err := c.ParamsParser(params); err != nil {
			return utils.NewAppError(err)
		}

		errors := validate.Struct(params)

		formattedErrorResponse := utils.TransformValidationErrorFormat(errors)

		if len(formattedErrorResponse) > 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(formattedErrorResponse)
		}

		c.Locals("params", params)

		return c.Next()
	}
}

func ValidateBlogPayload(routeName string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var body interface{}

		switch routeName {
		case constants.RouteName.CREATE_BLOG:
			body = new(CreateBlogPayload)
		case constants.RouteName.UPDATE_BLOG:
			body = new(UpdateBlogPayload)
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
