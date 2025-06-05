package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"prueba_tecnica_go_guarapo/api/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testScenarios := []struct {
		testName       string
		requestBody    interface{}
		mockSetup      func(*mockAuthService)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:    "Login exitoso",
			requestBody: models.LoginRequest{Username: "user1"},
			mockSetup: func(m *mockAuthService) {
				m.On("Login", mock.Anything, "user1").Return("token123")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"token123"}`,
		},
		{
			testName:       "JSON inválido",
			requestBody:    `{bad json}`,
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Username es requerido"}`,
		},
		{
			testName:       "Username vacío",
			requestBody:    models.LoginRequest{Username: ""},
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Username es requerido"}`,
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mockAuthService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			logger := logrus.New()
			handler := NewAuthHandler(mockService, logger)

			router := gin.New()
			router.POST("/login", handler.Login)

			var reqBody []byte
			var err error
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
