package models

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
	Owner     string `json:"-"` // el username dueño de la tarea
}
