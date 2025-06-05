package handlers

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, username string) string {
	args := m.Called(ctx, username)
	return args.String(0)
}
func (m *mockAuthService) ValidateToken(token string) (string, bool) {
	args := m.Called(token)
	return args.String(0), args.Bool(1)
}
