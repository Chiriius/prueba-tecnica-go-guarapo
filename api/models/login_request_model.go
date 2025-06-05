package models

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
}
