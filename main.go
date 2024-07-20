package main

import (
	"fmt"
	"go_blogs/configs"
	"go_blogs/connections"
	"go_blogs/models"
	"go_blogs/routes"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "go_blogs/docs"
)

// @title			Golang Blog CRUD
// @version		1.0
// @description	The simple CRUD project
func main() {
	configs.InitEnv()

	connections.InitDatabaseConnection()

	app := fiber.New(fiber.Config{
		AppName:     "Go Blogs",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			fmt.Println("ErrorHandler:", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Message: "Something went wrong",
			})
		},
	})

	routes.InitRoute(app)

	if configs.Env.AppEnv != "production" {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	err := app.Listen(fmt.Sprintf(":%d", configs.Env.Port))
	if err != nil {
		panic(err)
	}
}
