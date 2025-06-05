package handlers

import (
	"net/http"
	"prueba_tecnica_go_guarapo/api/models"
	services "prueba_tecnica_go_guarapo/api/services/auth"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler interface {
	Login(c *gin.Context)
}

type authHandler struct {
	authService services.AuthService
	logger      *logrus.Logger
}

func NewAuthHandler(authService services.AuthService, logger *logrus.Logger) AuthHandler {
	return &authHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *authHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("[Layer: auth_handler] [Method: Login] Invalid JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username es requerido"})
		return
	}

	if strings.TrimSpace(req.Username) == "" {
		h.logger.Warn("[Layer: auth_handler] [Method: Login] Username vacío")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username no puede estar vacío"})
		return
	}

	token := h.authService.Login(c.Request.Context(), req.Username)
	h.logger.Infof("[Layer: auth_handler] [Method: Login] Usuario '%s' autenticado", req.Username)
	response := models.LoginResponse{Token: token}

	c.JSON(http.StatusOK, response)
}
