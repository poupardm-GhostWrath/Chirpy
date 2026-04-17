package auth

import (
	"testing"
	"time"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name 			string
		password 		string
		hash 			string
		wantErr 		bool
		matchPassword 	bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	user := uuid.New()
	tokenSecret1 := "monkey"
	tokenSecret2 := "banana"
	duration1, _ := time.ParseDuration("1h")
	duration2, _ := time.ParseDuration("1ms")
	tokenString1, _ := MakeJWT(user, tokenSecret1, duration1) // Default Token
	tokenString2, _ := MakeJWT(user, tokenSecret1, duration2) // Short Duration Token

	tests := []struct {
		name		string
		tokenString	string
		tokenSecret	string
		wantErr		bool
	}{
		{
			name:			"correct token",
			tokenString:	tokenString1,
			tokenSecret:	tokenSecret1,
			wantErr:		false,
		},
		{
			name:			"incorrect token secret",
			tokenString:	tokenString1,
			tokenSecret:	tokenSecret2,
			wantErr:		true,
		},
		{
			name:			"expired token",
			tokenString:	tokenString2,
			tokenSecret:	tokenSecret1,
			wantErr:		true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && userID != user {
				t.Errorf("ValidateJWT() expects %v, got %v", user, userID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tokenString := "Monkey"
	header1 := make(http.Header)
	header2 := make(http.Header)
	headerValue := strings.Join([]string{"Bearer", tokenString}, " ")
	header1.Add("Authorization", headerValue)

	tests := []struct {
		name	string
		header	http.Header
		wantErr	bool
	}{
		{
			name:		"Valid Header",
			header: 	header1,
			wantErr:	false,
		},
		{
			name:		"No Header",
			header:		header2,
			wantErr:	true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token_string, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token_string != tokenString {
				t.Errorf("GetBearerToken() expects %v, got %v", tokenString, token_string)
			}
		})
	}
}