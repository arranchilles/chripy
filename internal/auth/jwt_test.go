package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "I am So Seceretive", time.Second*30)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	/*if len(token) != 256 {
		t.Errorf("Invalid, expected token string length :%d", len(token))
	}*/

	returnId, err := ValidateJWT(token, "I am So Seceretive")

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if returnId != id {
		t.Errorf("local ID and token ID are different local :%s, token:%s", id.String(), returnId.String())
	}

	fmt.Print(id.String()+"\n", returnId.String())

}

func TestJWTFailDifferenKey(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "I am So Seceretive", time.Second*30)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	/*if len(token) != 256 {
		t.Errorf("Invalid, expected token string length :%d", len(token))
	}*/

	returnId, err := ValidateJWT(token, "I am So retarded")

	if err == nil {
		t.Errorf("Key should be invalid")
	}

	if returnId == id {
		t.Errorf("local ID and token ID are the same local :%s, token:%s", id.String(), returnId.String())
	}

}

func TestGetBearerTokenPass(t *testing.T) {

	id := uuid.New()
	token, err := MakeJWT(id, "I am So Seceretive", time.Second*30)
	if err != nil {
		t.Error(err)
	}
	header := http.Header{}
	header.Add("Authorization", token)
	returnedToken, err := GetBearerToken(header)
	if err != nil {
		t.Error(err)
	}
	if returnedToken != token {
		t.Errorf("Tokens are not the same %s \n %s", token, returnedToken)
	}
}
