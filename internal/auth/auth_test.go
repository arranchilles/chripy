package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "GiveMeYourPasswordYouFuckingLemming!42"
	hash, err := HashPassword(password)
	if err != nil {
		errorValue := err.Error()
		t.Errorf("%v", errorValue)
	}

	isPassword, err := CheckPasswordHash(password, hash)

	if err != nil {
		errorValue := err.Error()
		t.Errorf("%v", errorValue)
	}

	if isPassword != true {
		t.Errorf("Incorrect denail of password")
	}

}

func TestHashPasswordBadPassword(t *testing.T) {
	password := "GiveMeYourPasswordYouFuckingLemming!42"
	hash, err := HashPassword(password)
	if err != nil {
		errorValue := err.Error()
		t.Errorf("%v", errorValue)
	}

	isPassword, err := CheckPasswordHash("cheeky bugger", hash)

	if err != nil {
		errorValue := err.Error()
		t.Errorf("%v", errorValue)
	}

	if isPassword == true {
		t.Errorf("Incorrect acceptance of false password")
	}

}
