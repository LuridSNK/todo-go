package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/luridsnk/todo-go/store"
)

type AuthService struct {
	store  *store.Store
	hasher *passwordHasher
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (auth *AuthService) ProcessLogin(login *LoginDto) (*User, *fiber.Error) {
	row, err := auth.store.QueryRow("select * from users where email = $1 limit 1", login.Email)
	if err != nil {
		return nil, &fiber.Error{Message: err.Error(), Code: 400}
	}
	var u User
	err = row.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.IsDeleted)
	if err != nil {
		return nil, &fiber.Error{Message: err.Error(), Code: 404}
	}

	if !hasher.CheckPasswordHash(login.Password, u.PasswordHash) {
		return nil, &fiber.Error{Message: "incorrect credentials", Code: 401}
	}
	return &u, nil
}

type RegisterDto struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (auth *AuthService) ProcessRegister(register *RegisterDto) (string, *fiber.Error) {
	pwdHash, err := auth.hasher.HashPassword(register.Password)
	if err != nil {
		return "", &fiber.Error{Message: err.Error(), Code: 404}
	}

	row, err := auth.store.QueryRow(
		"select exists(select 1 from users where email=$1)",
		register.Email)
	if err != nil {
		return "", &fiber.Error{Message: err.Error(), Code: 404}
	}
	var exists bool
	row.Scan(&exists)
	if exists {
		return "", &fiber.Error{Message: "user exists", Code: 415}
	}

	row, err = auth.store.QueryRow(
		"insert into users (email, passwordHash) values ($1, $2) returning id",
		register.Email,
		pwdHash)
	if err != nil {
		return "", &fiber.Error{Message: err.Error(), Code: 404}
	}
	var id string
	err = row.Scan(&id)
	if err != nil {
		return "", &fiber.Error{Message: err.Error(), Code: 404}
	}
	return id, nil
}
