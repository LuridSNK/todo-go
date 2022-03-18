package common

import (
	"testing"
)

// go test -run ”

func TestPasswordHashingShouldVerify(t *testing.T) {

	password := "test123"
	hash, _ := HashPassword("test123")
	verified := CheckPasswordHash(password, hash)
	if !verified {
		t.Error("not verified")
	}
}

func TestPasswordHashingShouldNotVerify(t *testing.T) {

	password := "test123"
	hash, _ := HashPassword("test1234")
	verified := CheckPasswordHash(password, hash)
	if verified {
		t.Error("should not be verified")
	}
}
