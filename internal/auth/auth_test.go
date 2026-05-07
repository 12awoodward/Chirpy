package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "testing", time.Hour)
	if err != nil {
		t.Errorf("Unable to gen token: %v", err)
	}

	testID, err := ValidateJWT(token, "testing")
	if err != nil {
		t.Errorf("Unable to validate token: %v", err)
	}

	if testID != userID {
		t.Errorf("IDs do not match: \n%v\n%v", userID, testID)
	}
}