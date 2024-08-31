package auth_test

import (
	"testing"

	"github.com/online-bnsp/backend/middleware/auth"
)

func TestEncodeDecodeJWT(t *testing.T) {
	user_id := 12345
	userEmail := "user@test.com"
	token, err := auth.EncodeJWT(int32(user_id), userEmail)
	if err != nil {
		t.Error("unable to encode jwt")
	}

	claim, err := auth.DecodeJWT(token)
	if err != nil {
		t.Error("unable to decode jwt")
	}

	if claim.UserID != int32(user_id) {
		t.Errorf("invalid user id, expect %v, got %v", user_id, claim.UserID)
	}
	// if claim.Email != userEmail {
	// 	t.Errorf("invalid user id, expect %v, got %v", userEmail, claim.UserID)
	// }
}
