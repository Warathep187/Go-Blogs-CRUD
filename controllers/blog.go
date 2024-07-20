package controllers

import (
	"context"
	"errors"
	"go_blogs/configs"
	"go_blogs/connections"
	"go_blogs/models"
	"go_blogs/utils"
	"go_blogs/validators"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogController interface {
	GetBlogs(c *fiber.Ctx) error
	GetBlogByID(c *fiber.Ctx) error
	CreateBlog(c *fiber.Ctx) error
	UpdateBlog(c *fiber.Ctx) error
	DeleteBlog(c *fiber.Ctx) error
}

type BlogController struct {
	MongoBlogColl *mongo.Collection
}

func NewBlogControllers() blogController {
	return &BlogController{
		MongoBlogColl: connections.NewMongoCollection(
			connections.MongoClient.Database(configs.Env.MongoDatabase),
			"blogs",
		),
	}
}

// @summary		Get blogs
// @description	Get all blogs (number of blogs per query is 10)
// @id				GetBlogs
// @tags			blogs
// @accept			json
// @produce		json
// @param			from	query		int	true	"blog offset"	default(0)	maximum(0)
// @success		200		{array}		models.Blog
// @failure		422		{array}		models.ValidationErrorResponse	"validation failed"
// @failure		500		{object}	models.ErrorResponse			"something went wrong"
// @router			/api/blogs [get]
func (ctr *BlogController) GetBlogs(c *fiber.Ctx) error {
	query := c.Locals("query").(*validators.GetBlogsQuery)

	ctx := context.TODO()

	queryLimit := 10
	opts := options.Find().SetSkip(int64(query.From)).SetLimit(int64(queryLimit)).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := ctr.MongoBlogColl.Find(ctx, bson.M{}, opts)
	if err != nil {
		return utils.NewAppError(err)
	}

	var blogs []models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return utils.NewAppError(err)
	}

	return c.Status(fiber.StatusOK).JSON(blogs)
}

// @summary		Get blog by ID
// @description	Get blog by ID
// @id				GetByID
// @tags			blogs
// @accept			json
// @produce		json
// @param			id	path		string	true	"blog's ID"
// @success		200	{object}	models.Blog
// @failure		404	{object}	models.ErrorResponse			"blog not found"
// @failure		422	{array}		models.ValidationErrorResponse	"validation failed"
// @failure		500	{object}	models.ErrorResponse			"something went wrong"
// @router			/api/blogs/:id [get]
func (ctr *BlogController) GetBlogByID(c *fiber.Ctx) error {
	params := c.Locals("params").(*validators.GetBlogByIDParams)
	blogObjectID, _ := primitive.ObjectIDFromHex(params.ID)

	var blog *models.Blog

	filter := bson.M{
		"_id": blogObjectID,
	}
	if err := ctr.MongoBlogColl.FindOne(context.TODO(), filter).Decode(&blog); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Message: "Blog not found",
			})
		}
		return utils.NewAppError(err)
	}

	return c.Status(fiber.StatusOK).JSON(blog)
}

// @summary		Create blog
// @description	Create new blog
// @id				CreateBlog
// @tags			blogs
// @accept			json
// @produce		json
// @param			title	body		string	true	"blog's title"
// @param			content	body		string	true	"blog's content"
// @success		201		{object}	string
// @failure		422		{array}		models.ValidationErrorResponse	"validation failed"
// @failure		500		{object}	models.ErrorResponse			"something went wrong"
// @router			/api/blogs [post]
func (ctr *BlogController) CreateBlog(c *fiber.Ctx) error {
	payload := c.Locals("payload").(*validators.CreateBlogPayload)
	user := c.Locals("user").(*models.UserSessionData)

	document := bson.D{
		{Key: "title", Value: payload.Title},
		{Key: "content", Value: payload.Content},
		{Key: "createdBy", Value: user.ID},
		{Key: "createdAt", Value: time.Now()},
	}
	_, err := ctr.MongoBlogColl.InsertOne(context.TODO(), document)
	if err != nil {
		return utils.NewAppError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Message: "Created",
	})
}

// @summary		Update blog
// @description	Update new blog
// @id				UpdateBlog
// @tags			blogs
// @accept			json
// @produce		json
// @param			id		path		string	true	"blog's ID"
// @param			title	body		string	true	"blog's title"
// @param			content	body		string	true	"blog's content"
// @success		200		{object}	string
// @failure		404		{object}	models.ErrorResponse			"blog not found"
// @failure		422		{array}		models.ValidationErrorResponse	"validation failed"
// @failure		500		{object}	models.ErrorResponse			"something went wrong"
// @router			/api/blogs/:id [put]
func (ctr *BlogController) UpdateBlog(c *fiber.Ctx) error {
	params := c.Locals("params").(*validators.UpdateBlogParams)
	body := c.Locals("payload").(*validators.UpdateBlogPayload)
	user := c.Locals("user").(*models.UserSessionData)

	ctx := context.TODO()

	blogObjectID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return utils.NewAppError(err)
	}

	var blog *models.Blog
	filter := bson.M{
		"_id": blogObjectID,
	}
	opts := options.FindOne().SetProjection(bson.M{"createdBy": 1})
	if err = ctr.MongoBlogColl.FindOne(ctx, filter, opts).Decode(&blog); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Message: "Blog not found",
			})
		}
		return utils.NewAppError(err)
	}
	if blog.CreatedBy != user.ID {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Message: "Access Denied",
		})
	}

	document := bson.M{
		"$set": bson.D{
			{Key: "title", Value: body.Title},
			{Key: "content", Value: body.Content},
		},
	}
	_, err = ctr.MongoBlogColl.UpdateByID(ctx, blogObjectID, document)
	if err != nil {
		return utils.NewAppError(err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Message: "Updated",
	})
}

// @summary		Delete blog
// @description	Delete new blog
// @id				DeleteBlog
// @tags			blogs
// @accept			json
// @produce		json
// @param			id	path		string	true	"blog's ID"
// @success		200	{object}	string
// @failure		404	{object}	models.ErrorResponse			"blog not found"
// @failure		422	{array}		models.ValidationErrorResponse	"validation failed"
// @failure		500	{object}	models.ErrorResponse			"something went wrong"
// @router			/api/blogs/:id [delete]
func (ctr *BlogController) DeleteBlog(c *fiber.Ctx) error {
	params := c.Locals("params").(*validators.DeleteBlogParams)
	user := c.Locals("user").(*models.UserSessionData)

	blogObjectID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return utils.NewAppError(err)
	}

	ctx := context.TODO()

	var blog *models.Blog

	filter := bson.M{
		"_id": blogObjectID,
	}
	opts := options.FindOne().SetProjection(bson.D{{Key: "createdBy", Value: 1}})
	if err = ctr.MongoBlogColl.FindOne(ctx, filter, opts).Decode(&blog); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Message: "Blog not found",
			})
		}
	}

	if user.ID != blog.CreatedBy {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Message: "Access Denied",
		})
	}

	_, err = ctr.MongoBlogColl.DeleteOne(ctx, filter)
	if err != nil {
		return utils.NewAppError(err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Message: "Deleted",
	})
}
