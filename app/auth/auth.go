package auth

import (
	"strings"
	"time"

	"github.com/luridsnk/todo-go/app"
	"github.com/luridsnk/todo-go/common"
	"github.com/luridsnk/todo-go/config"
	"github.com/luridsnk/todo-go/store"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var secret string
var hasher passwordHasher = passwordHasher{}
var tokenExp int64

func UseEndpoints(application *app.App, store *store.Store, appConfig *config.ApplicationConfig) {
	secret = appConfig.Secret
	tokenExp = time.Now().UTC().Add(time.Hour * appConfig.TokenExpiry).Unix()

	todoGroup := application.Group("api/v1/account")
	{
		srv := &AuthService{store: store, hasher: &hasher}
		todoGroup.Post("/login", login(srv))
		todoGroup.Post("/register", register(srv))
	}
}

func login(authSrv *AuthService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loginDto, err := common.ReadJson[LoginDto](c.Body())
		if err != nil {
			return err
		}

		u, err := authSrv.ProcessLogin(loginDto)
		if err != nil {
			return err
		}
		claims := jwt.MapClaims{
			"sub":   u.Id.String(),
			"email": u.Email,
			"exp":   tokenExp,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}

func register(authSrv *AuthService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		reg, err := common.ReadJson[RegisterDto](c.Body())
		if err != nil {
			return err
		}

		if r := strings.Compare(reg.Password, reg.PasswordConfirmation); r != 0 {
			c.JSON(fiber.Map{"error": "provided passwords don't match"})
			return c.SendStatus(400)
		}

		id, err := authSrv.ProcessRegister(reg)
		if err != nil {
			return err
		}

		claims := jwt.MapClaims{
			"sub":   id,
			"email": reg.Email,
			"exp":   tokenExp,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
