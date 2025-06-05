package handlers

import (
	"context"
	"prueba_tecnica_go_guarapo/api/models"

	"github.com/stretchr/testify/mock"
)

type mockTaskService struct {
	mock.Mock
}

func (m *mockTaskService) GetTasksByUser(ctx context.Context, username string) ([]*models.Task, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]*models.Task), args.Error(1)
}
func (m *mockTaskService) GetTaskByID(ctx context.Context, id int, username string) (*models.Task, error) {
	args := m.Called(ctx, id, username)
	return args.Get(0).(*models.Task), args.Error(1)
}
func (m *mockTaskService) CreateTask(ctx context.Context, title string, username string) (*models.Task, error) {
	args := m.Called(ctx, title, username)
	return args.Get(0).(*models.Task), args.Error(1)
}
func (m *mockTaskService) UpdateTask(ctx context.Context, id int, title string, completed bool, username string) (*models.Task, error) {
	args := m.Called(ctx, id, title, completed, username)
	return args.Get(0).(*models.Task), args.Error(1)
}
func (m *mockTaskService) DeleteTask(ctx context.Context, id int, username string) error {
	args := m.Called(ctx, id, username)
	return args.Error(0)
}
