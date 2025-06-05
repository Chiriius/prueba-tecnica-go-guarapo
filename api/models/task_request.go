package models

type CreateTaskRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateTaskRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}
