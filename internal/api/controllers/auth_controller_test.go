package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zhaoyikaiii/docmind/pkg/auth"
	"github.com/Zhaoyikaiii/docmind/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.LoadConfig()
}

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ac := &AuthController{}
	r.POST("/login", ac.Login)

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
	}{
		{
			name: "Valid login",
			requestBody: LoginRequest{
				Username: "testuser",
				Password: "password",
			},
			expectedCode: 200,
		},
		{
			name: "Missing username",
			requestBody: LoginRequest{
				Password: "password",
			},
			expectedCode: 400,
		},
		{
			name: "Missing password",
			requestBody: LoginRequest{
				Username: "testuser",
			},
			expectedCode: 400,
		},
		{
			name:         "Invalid JSON",
			requestBody:  "invalid json",
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			var body bytes.Buffer
			if jsonBody, ok := tt.requestBody.(LoginRequest); ok {
				json.NewEncoder(&body).Encode(jsonBody)
			} else {
				body.WriteString(tt.requestBody.(string))
			}

			req, _ := http.NewRequest("POST", "/login", &body)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == 200 {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
				assert.Contains(t, response, "token_type")
				assert.Equal(t, "Bearer", response["token_type"])
			}
		})
	}
}

func TestAuthController_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ac := &AuthController{}
	r.POST("/refresh", ac.RefreshToken)

	tokenInfo, _ := auth.GenerateToken(1, "testuser", "user")

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
	}{
		{
			name: "Valid refresh token",
			requestBody: RefreshRequest{
				RefreshToken: tokenInfo.RefreshToken,
			},
			expectedCode: 200,
		},
		{
			name: "Invalid refresh token",
			requestBody: RefreshRequest{
				RefreshToken: "invalid.token.here",
			},
			expectedCode: 401,
		},
		{
			name: "Empty refresh token",
			requestBody: RefreshRequest{
				RefreshToken: "",
			},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			req, _ := http.NewRequest("POST", "/refresh", &body)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == 200 {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "access_token")
				assert.Contains(t, response, "refresh_token")
			}
		})
	}
}
