package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Login(ctx context.Context, username string) string
	ValidateToken(token string) (string, bool)
}

type authService struct {
	tokens map[string]string
	mutex  sync.RWMutex
	logger *logrus.Logger
}

func NewAuthService(logger *logrus.Logger) AuthService {
	return &authService{
		tokens: make(map[string]string),
		logger: logger,
	}
}

func (s *authService) Login(ctx context.Context, username string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	token := s.generateToken()
	s.tokens[token] = username
	s.logger.Infof("[Layer: auth_service] [Method: Login] Info: User '%s' logged in with token '%s'\n", username, token)
	return token
}

func (s *authService) ValidateToken(token string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	username, exists := s.tokens[token]
	if exists {
		s.logger.Infof("[Layer: auth_service] [Method: ValidateToken] Info: Token '%s' is valid for user '%s'\n", token, username)
	} else {
		s.logger.Warnf("[Layer: auth_service] [Method: ValidateToken] Warning: Invalid token '%s'\n", token)
	}
	return username, exists
}

func (s *authService) generateToken() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		s.logger.Errorf("[Layer: auth_service] [Method: generateToken] Error: Failed to generate random bytes: %v\n", err)
		return ""
	}
	token := hex.EncodeToString(bytes)
	s.logger.Infof("[Layer: auth_service] [Method: generateToken] Info: Generated token '%s'\n", token)
	return token
}
