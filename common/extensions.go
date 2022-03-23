package common

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const TokenKey string = "user"

func ReadJson[T any](bytes []byte) (*T, error) {
	var value T
	if err := json.Unmarshal(bytes, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func ToJson(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GetUserFromContext(c *fiber.Ctx) string {
	return c.Locals(TokenKey).(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
}
