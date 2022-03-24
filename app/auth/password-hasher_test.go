package auth

import (
	"testing"
)

func TestPasswordHashingShouldVerify(t *testing.T) {

	password := "test123"
	hash, _ := hasher.HashPassword("test123")
	verified := hasher.CheckPasswordHash(password, hash)
	if !verified {
		t.Error("not verified")
	}
}

func TestPasswordHashingShouldNotVerify(t *testing.T) {

	password := "test123"
	hash, _ := hasher.HashPassword("test1234")
	verified := hasher.CheckPasswordHash(password, hash)
	if verified {
		t.Error("should not be verified")
	}
}
