package routes

import (
	"go_blogs/constants"
	"go_blogs/controllers"
	"go_blogs/middlewares"
	"go_blogs/validators"

	"github.com/gofiber/fiber/v2"
)

func InitRoute(app *fiber.App) {
	authControllers := controllers.NewAuthControllers()
	blogControllers := controllers.NewBlogControllers()

	api := app.Group("/api")

	// /api/auth/
	authApi := api.Group("/auth")
	authApi.Post(
		"/login",
		validators.ValidateAuthPayload(constants.RouteName.LOGIN),
		authControllers.Login,
	)
	authApi.Post(
		"/register",
		validators.ValidateAuthPayload(constants.RouteName.REGISTER),
		authControllers.Register,
	)
	authApi.Post("/logout", middlewares.AuthorizeUser, authControllers.Logout)
	authApi.Get("/user", authControllers.GetUserData)

	// /api/blogs
	blogsApi := api.Group("/blogs")
	blogsApi.Get(
		"/",
		validators.ValidateBlogQuery(constants.RouteName.GET_BLOGS),
		blogControllers.GetBlogs,
	)
	blogsApi.Get(
		"/:id",
		validators.ValidateBlogParams(constants.RouteName.GET_BLOG_BY_ID),
		blogControllers.GetBlogByID,
	)
	blogsApi.Post("/",
		middlewares.AuthorizeUser,
		validators.ValidateBlogPayload(constants.RouteName.CREATE_BLOG),
		blogControllers.CreateBlog,
	)
	blogsApi.Put("/:id",
		middlewares.AuthorizeUser,
		validators.ValidateBlogParams(constants.RouteName.UPDATE_BLOG),
		validators.ValidateBlogPayload(constants.RouteName.UPDATE_BLOG),
		blogControllers.UpdateBlog,
	)
	blogsApi.Delete("/:id",
		middlewares.AuthorizeUser,
		validators.ValidateBlogParams(constants.RouteName.DELETE_BLOG),
		blogControllers.DeleteBlog,
	)
}
