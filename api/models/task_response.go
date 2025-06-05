package models

type TaskResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Owner     string `json:"owner"`
}
