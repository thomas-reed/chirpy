package auth

import (
	"net/http"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct{
		name string
		header http.Header
		wantErr bool
	}{
		{
			name: "Exact Bearer token",
			header: http.Header{
				"Authorization": []string{"Bearer TOKEN_STRING"},
			},
			wantErr: false,
		},
		{
			name: "Has Bearer token amidst other headers",
			header: http.Header{
				"Authorization": []string{"Bearer TOKEN_STRING", "Other SOME_OTHER_STRING"},
				"Non-Authorization": []string{"Bearer NON_AUTH_TOKEN_STRING", "Another SOMETHING_ELSE"},
			},
			wantErr: false,
		},
		{
			name: "Bearer token is missing amidst other headers",
			header: http.Header{
				"Authorization": []string{"Access-token TOKEN_STRING", "Other SOME_OTHER_STRING"},
				"Non-Authorization": []string{"Bearer NON_AUTH_TOKEN_STRING", "Another SOMETHING_ELSE"},
			},
			wantErr: true,
		},
		{
			name: "Has no Authorization header",
			header: http.Header{
				"Non-Authorization": []string{"Bearer NON_AUTH_TOKEN_STRING", "Another SOMETHING_ELSE"},
			},
			wantErr: true,
		},
		{
			name: "Bearer token has extra whitespace",
			header: http.Header{
				"Authorization": []string{"  Bearer      TOKEN_STRING", "Other SOME_OTHER_STRING"},
				"Non-Authorization": []string{"Bearer NON_AUTH_TOKEN_STRING", "Another SOMETHING_ELSE"},
			},
			wantErr: false,
		},
		{
			name: "Bearer with no token",
			header: http.Header{
				"Authorization": []string{"Bearer"},
			},
			wantErr: true,
		},
		{
			name: "Bearer with only spaces",
			header: http.Header{
				"Authorization": []string{"Bearer      "},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T){
			token, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && token != "TOKEN_STRING" {
				t.Errorf("GetBearerToken() token = %v, Expected 'TOKEN_STRING'", token)
			}
		})
	}
}