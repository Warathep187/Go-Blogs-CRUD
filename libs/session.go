package libs

import (
	"context"
	"encoding/json"
	"fmt"
	"go_blogs/connections"
	"go_blogs/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var sessionStorage = session.New(session.Config{
	Expiration:     time.Hour * 6,
	CookieHTTPOnly: true,
})

func getRequestSession(c *fiber.Ctx) (*session.Session, error) {
	sess, err := sessionStorage.Get(c)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func GetUserSessionData(c *fiber.Ctx) (*models.UserSessionData, error) {
	sess, err := getRequestSession(c)
	if err != nil {
		return &models.UserSessionData{}, err
	}

	key := fmt.Sprintf("sess:%s", sess.ID())

	result, err := connections.RedisClient.Get(context.TODO(), key).Result()
	if err != nil {
		return &models.UserSessionData{}, err
	}

	var unmarshalSessionData *models.UserSessionData

	resultInJSON := []byte(result)
	if err = json.Unmarshal(resultInJSON, &unmarshalSessionData); err != nil {
		return &models.UserSessionData{}, err
	}

	return unmarshalSessionData, nil
}

func SetUserSessionData(c *fiber.Ctx, data models.UserSessionData) error {
	sess, err := getRequestSession(c)
	if err != nil {
		return err
	}

	sessionKey := fmt.Sprintf("sess:%s", sess.ID())

	marshaledSessionData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err = connections.RedisClient.Set(context.TODO(), sessionKey, string(marshaledSessionData), time.Hour*6).Err(); err != nil {
		return err
	}

	if err = sess.Save(); err != nil {
		return err
	}

	return nil
}

func DestroyUserSessionData(c *fiber.Ctx) error {
	sess, err := getRequestSession(c)
	if err != nil {
		return err
	}

	sessionKey := fmt.Sprintf("sess:%s", sess.ID())
	if err = connections.RedisClient.Del(context.TODO(), sessionKey).Err(); err != nil {
		return err
	}

	if err = sess.Destroy(); err != nil {
		return err
	}

	return nil
}
