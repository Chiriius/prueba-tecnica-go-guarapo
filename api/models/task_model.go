package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
	Owner     string `json:"-"` // el username dueño de la tarea
}
