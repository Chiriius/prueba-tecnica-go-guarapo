package services

import (
	"context"
	"prueba_tecnica_go_guarapo/api/models"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&models.Task{}))
	return db
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	service := NewTaskService(db, logrus.New())
	ctx := context.Background()

	task, err := service.CreateTask(ctx, "", "user1")
	assert.Nil(t, task)
	assert.ErrorIs(t, err, ErrTitleRequired)

	task, err = service.CreateTask(ctx, "Test", "")
	assert.Nil(t, task)
	assert.ErrorIs(t, err, ErrUserRequired)

	// Caso: éxito
	task, err = service.CreateTask(ctx, "Test Task", "user1")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, false, task.Completed)
	assert.Equal(t, "user1", task.Owner)
}

func TestGetTasksByUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewTaskService(db, logrus.New())
	ctx := context.Background()

	tasks, err := service.GetTasksByUser(ctx, "user1")
	assert.NoError(t, err)
	assert.Len(t, tasks, 0)

	_, _ = service.CreateTask(ctx, "Task 1", "user1")
	_, _ = service.CreateTask(ctx, "Task 2", "user1")
	_, _ = service.CreateTask(ctx, "Task 3", "user2")

	tasks, err = service.GetTasksByUser(ctx, "user1")
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
}

func TestGetTaskByID(t *testing.T) {
	db := setupTestDB(t)
	service := NewTaskService(db, logrus.New())
	ctx := context.Background()

	task, _ := service.CreateTask(ctx, "Task", "user1")

	got, err := service.GetTaskByID(ctx, int(task.ID), "user1")
	assert.NoError(t, err)
	assert.Equal(t, task.ID, got.ID)

	_, err = service.GetTaskByID(ctx, 999, "user1")
	assert.ErrorIs(t, err, ErrTaskNotFound)

	_, err = service.GetTaskByID(ctx, int(task.ID), "otro")
	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	service := NewTaskService(db, logrus.New())
	ctx := context.Background()

	task, _ := service.CreateTask(ctx, "Task", "user1")

	updated, err := service.UpdateTask(ctx, int(task.ID), "Updated", true, "user1")
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.Title)
	assert.True(t, updated.Completed)

	_, err = service.UpdateTask(ctx, 999, "Nope", false, "user1")
	assert.ErrorIs(t, err, ErrTaskNotFound)

	_, err = service.UpdateTask(ctx, int(task.ID), "Nope", false, "otro")
	assert.ErrorIs(t, err, ErrTaskNotFound)
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	service := NewTaskService(db, logrus.New())
	ctx := context.Background()

	task, _ := service.CreateTask(ctx, "Task", "user1")

	err := service.DeleteTask(ctx, int(task.ID), "user1")
	assert.NoError(t, err)

	err = service.DeleteTask(ctx, int(task.ID), "user1")
	assert.ErrorIs(t, err, ErrTaskNotFound)

	err = service.DeleteTask(ctx, 999, "user1")
	assert.ErrorIs(t, err, ErrTaskNotFound)

	task2, _ := service.CreateTask(ctx, "Task2", "user2")
	err = service.DeleteTask(ctx, int(task2.ID), "user1")
	assert.ErrorIs(t, err, ErrTaskNotFound)
}
