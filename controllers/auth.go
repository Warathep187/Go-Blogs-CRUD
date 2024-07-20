package controllers

import (
	"context"
	"errors"
	"go_blogs/configs"
	"go_blogs/connections"
	"go_blogs/libs"
	"go_blogs/models"
	"go_blogs/utils"
	"go_blogs/validators"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type authController interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	GetUserData(c *fiber.Ctx) error
}

type AuthController struct {
	MongoUserColl *mongo.Collection
}

func NewAuthControllers() authController {
	return &AuthController{
		MongoUserColl: connections.NewMongoCollection(
			connections.MongoClient.Database(configs.Env.MongoDatabase),
			"users",
		),
	}
}

// @summary		Login
// @description	User Login
// @tags			auth
// @id				Login
// @accept			json
// @produce		json
// @param			email		body		string	true	"email"
// @param			password	body		string	true	"password"	minlength(6)	maxlength(32)
// @success		200			{object}	models.UserSessionData
// @failure		400			{object}	models.ErrorResponse			"some condition failed"
// @failure		422			{array}		models.ValidationErrorResponse	"validation failed"
// @failure		500			{object}	models.ErrorResponse			"something went wrong"
// @router			/api/auth/login [post]
func (ctr *AuthController) Login(c *fiber.Ctx) error {
	payload := c.Locals("payload").(*validators.LoginPayload)

	var result models.User
	err := ctr.MongoUserColl.FindOne(context.TODO(), bson.M{
		"email": payload.Email,
	}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Email or password is invalid",
		})
	} else if err != nil {
		return utils.NewAppError(err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(payload.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Message: "Email or password is invalid",
		})
	}

	userSessionData := models.UserSessionData{
		ID:    result.ID,
		Email: result.Email,
		Name:  result.Name,
	}
	err = libs.SetUserSessionData(c, userSessionData)
	if err != nil {
		return utils.NewAppError(err)
	}

	return c.JSON(userSessionData)
}

// @summary		Register
// @description	Registration
// @tags			auth
// @id				Register
// @accept			json
// @produce		json
// @param			email		body		string	true	"email"
// @param			password	body		string	true	"password"	minlength(5)	maxlength(32)
// @param			name		body		string	true	"Name"
// @success		200			{object}	models.UserSessionData
// @failure		400			{object}	models.ErrorResponse			"some condition failed"
// @failure		422			{array}		models.ValidationErrorResponse	"validation failed"
// @failure		500			{object}	models.ErrorResponse			"something went wrong"
// @router			/api/auth/register [post]
func (ctr *AuthController) Register(c *fiber.Ctx) error {
	payload := c.Locals("payload").(*validators.RegisterPayload)

	ctx := context.TODO()

	var user *models.User

	filter := bson.M{
		"email": payload.Email,
	}
	err := ctr.MongoUserColl.FindOne(ctx, filter).Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return utils.NewAppError(err)
	}
	if user != nil {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Message: "Email has been used. Please use another email",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewAppError(err)
	}

	document := bson.D{
		{Key: "email", Value: payload.Email},
		{Key: "password", Value: string(hashedPassword)},
		{Key: "name", Value: payload.Name},
	}
	_, err = ctr.MongoUserColl.InsertOne(ctx, document)
	if err != nil {
		return utils.NewAppError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Message: "Registered",
	})
}

// @summary		Logout
// @description	User Logout
// @tags			auth
// @id				Logout
// @accept			json
// @produce		json
// @success		200	{object}	models.UserSessionData
// @failure		401	{object}	models.ErrorResponse	"unauthorized"
// @failure		500	{object}	models.ErrorResponse	"something went wrong"
// @router			/api/auth/logout [post]
func (ctr *AuthController) Logout(c *fiber.Ctx) error {
	if err := libs.DestroyUserSessionData(c); err != nil {
		return utils.NewAppError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (ctr *AuthController) GetUserData(c *fiber.Ctx) error {
	userData, err := libs.GetUserSessionData(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Message: fiber.ErrUnauthorized.Message,
		})
	}
	return c.JSON(userData)
}
