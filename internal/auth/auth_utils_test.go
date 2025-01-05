package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword123!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct Password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect Password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty Password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPasswordHash(tt.password, tt.hash); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	tokenSecret := "SECRET"
	wrongSecret := "WRONG_SECRET"
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, tokenSecret, time.Hour)
	expiredToken, _ := MakeJWT(userID, tokenSecret, -1*time.Hour)
	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			tokenSecret: tokenSecret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Expired Token",
			tokenString: expiredToken,
			tokenSecret: tokenSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Invalid Token",
			tokenString: "invalid.token.string",
			tokenSecret: validToken,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong Secret",
			tokenString: validToken,
			tokenSecret: wrongSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, wantUserID = %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer Token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization Header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization Header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, wantToken = %v", gotToken, tt.wantToken)
			}
		})
	}
}
