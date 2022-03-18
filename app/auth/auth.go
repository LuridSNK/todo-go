package auth

import (
	"strings"
	"time"
	"todo_app/app"
	"todo_app/common"
	"todo_app/store"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const secret string = "secret"

func UseEndpoints(application *app.App, store *store.Store) {
	todoGroup := application.Group("api/v1/account")
	{
		todoGroup.Post("/", login(store))
	}
}

func login(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type login struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		l, err := common.ReadJson[login](c.Body())
		if err != nil {
			return err
		}

		row, err := store.QueryRow("select * from users where email = $1", l.Email)
		if err != nil {
			return err
		}
		var user User
		err = row.Scan(&user)
		if err != nil {
			return err
		}

		if !common.CheckPasswordHash(l.Password, user.PasswordHash) {
			c.JSON(fiber.Map{"error": "credentials are incorrect"})
			return c.SendStatus(401)
		}

		claims := jwt.MapClaims{
			"sub":   user.Id.String(),
			"email": user.Email,
			"exp":   time.Now().UTC().Add(time.Hour * 48).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}

func register(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type login struct {
			Email                string `json:"email"`
			Password             string `json:"password"`
			PasswordConfirmation string `json:"passwordConfirmation"`
		}
		l, err := common.ReadJson[login](c.Body())
		if err != nil {
			return err
		}

		if r := strings.Compare(l.Password, l.PasswordConfirmation); r != 0 {
			c.JSON(fiber.Map{"error": "provided passwords don't match"})
			return c.SendStatus(400)
		}

		pwdHash, err := common.HashPassword(l.Password)
		if err != nil {
			return err
		}

		row, err := store.QueryRow("insert into users (email, passwordHash) values ($1, $2) returning id", l.Email, pwdHash)
		if err != nil {
			return err
		}
		var id uuid.UUID
		err = row.Scan(&id)
		if err != nil {
			return err
		}

		claims := jwt.MapClaims{
			"sub":   id.String(),
			"email": l.Email,
			"exp":   time.Now().UTC().Add(time.Hour * 48).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
