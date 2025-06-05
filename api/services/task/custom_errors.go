package services

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found or not owned by user")
	ErrTitleRequired = errors.New("title is required")
	ErrUserRequired  = errors.New("UserName is required")
)
