package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "testing", time.Hour)
	if err != nil {
		t.Errorf("Unable to gen token:\n%v", err)
	}

	testID, err := ValidateJWT(token, "testing")
	if err != nil {
		t.Errorf("Unable to validate token:\n%v", err)
	}

	if testID != userID {
		t.Errorf("IDs do not match:\n%v\n%v", userID, testID)
	}
}

func TestHeaderToken(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "testing", time.Hour)
	if err != nil {
		t.Errorf("Unable to gen token:\n%v", err)
	}

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	findToken, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("Couldn't get token:\n%v", err)
	}

	if findToken != token {
		t.Errorf("tokens do not match: \n%v\n%v", token, findToken)
	}
}
