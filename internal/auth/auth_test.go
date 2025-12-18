package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	secret := "AllYourBase"
	userID := uuid.MustParse("59d051d3-63e5-4eae-a854-f7d12a1097b0")

	jwtValid, err1 := MakeJWT(userID, secret, 5*time.Minute)
	jwtExpired, err2 := MakeJWT(userID, secret, 0)
	if err1 != nil || err2 != nil {
		t.Fatalf("TestValidateJWT() jwt setup failed err1 = %v err2 = %v", err1, err2)
	}

	tests := []struct {
		name    string
		jwt     string
		wantErr bool
	}{
		{
			name:    "Validate valid",
			jwt:     jwtValid,
			wantErr: false,
		},
		{
			name:    "Validate expired",
			jwt:     jwtExpired,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validatedUserID, err := ValidateJWT(tt.jwt, secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && validatedUserID != userID {
				t.Errorf("ValidateJWT() expects %v, got %v", userID, validatedUserID)
			}
		})
	}
}
