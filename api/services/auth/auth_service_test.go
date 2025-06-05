package services

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_LoginAndValidateToken(t *testing.T) {
	logger := logrus.New()
	service := NewAuthService(logger)

	testScenarios := []struct {
		testName   string
		username   string
		wantExists bool
	}{
		{
			testName:   "Login and validate valid token",
			username:   "user1",
			wantExists: true,
		},
		{
			testName:   "Validate invalid token",
			username:   "",
			wantExists: false,
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			ctx := context.Background()
			var token string
			if tt.username != "" {
				token = service.Login(ctx, tt.username)
				assert.NotEmpty(t, token)
			} else {
				token = "invalidtoken"
			}
			username, exists := service.ValidateToken(token)
			if tt.wantExists {
				assert.True(t, exists)
				assert.Equal(t, tt.username, username)
			} else {
				assert.False(t, exists)
				assert.Empty(t, username)
			}
		})
	}
}
