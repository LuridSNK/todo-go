package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/luridsnk/todo-go/app"
	"github.com/luridsnk/todo-go/common"
	"github.com/luridsnk/todo-go/store"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var secret string

func UseEndpoints(application *app.App, store *store.Store, s string) {
	secret = s
	todoGroup := application.Group("api/v1/account")
	{
		todoGroup.Post("/login", login(store))
		todoGroup.Post("/register", register(store))
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

		row, err := store.QueryRow("select * from users where email = $1 limit 1", l.Email)
		if err != nil {
			return err
		}
		var u User
		err = row.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.IsDeleted)
		if err != nil {
			c.JSON(fiber.Map{
				"error": "user not found",
				"descr": err.Error(),
			})
			return c.SendStatus(404)
		}

		if !common.CheckPasswordHash(l.Password, u.PasswordHash) {
			c.JSON(fiber.Map{"error": "credentials are incorrect"})
			return c.SendStatus(401)
		}

		claims := jwt.MapClaims{
			"sub":   u.Id.String(),
			"email": u.Email,
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

		row, err := store.QueryRow(
			"select exists(select 1 from users where email=$1)",
			l.Email)
		if err != nil {
			return err
		}
		var exists bool
		row.Scan(&exists)
		if exists {
			c.Status(415)
			return c.SendString(fmt.Sprintf("User exists"))
		}

		row, err = store.QueryRow(
			"insert into users (email, passwordHash) values ($1, $2) returning id",
			l.Email,
			pwdHash)
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
