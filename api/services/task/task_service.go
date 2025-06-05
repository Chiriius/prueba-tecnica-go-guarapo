package services

import (
	"context"
	"prueba_tecnica_go_guarapo/api/models"
	"sync"

	"github.com/sirupsen/logrus"
)

type TaskService interface {
	GetTasksByUser(ctx context.Context, username string) ([]*models.Task, error)
	GetTaskByID(ctx context.Context, id int, username string) (*models.Task, error)
	CreateTask(ctx context.Context, title string, username string) (*models.Task, error)
	UpdateTask(ctx context.Context, id int, title string, completed bool, username string) (*models.Task, error)
	DeleteTask(ctx context.Context, id int, username string) error
}

type taskService struct {
	tasks  map[int]*models.Task
	nextID int
	mu     sync.RWMutex
	logger *logrus.Logger
}

func NewTaskService(logger *logrus.Logger) TaskService {
	return &taskService{
		tasks:  make(map[int]*models.Task),
		nextID: 1,
		logger: logger,
	}
}

func (s *taskService) GetTasksByUser(ctx context.Context, username string) ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userTasks []*models.Task
	for _, task := range s.tasks {
		if task.Owner == username {
			userTasks = append(userTasks, task)
		}
	}
	s.logger.Infof("[Layer: task_service] [Method: GetTasksByUser] Info: User '%s' requested their tasks", username)
	return userTasks, nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id int, username string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists || task.Owner != username {
		s.logger.Warnf("[Layer: task_service] [Method: GetTaskByID] Warning: Task '%d' not found or not owned by user '%s'", id, username)
		return nil, ErrTaskNotFound
	}
	s.logger.Infof("[Layer: task_service] [Method: GetTaskByID] Info: Task '%d' retrieved for user '%s'", id, username)
	return task, nil
}

func (s *taskService) CreateTask(ctx context.Context, title string, username string) (*models.Task, error) {
	if title == "" {
		s.logger.Errorln("[Layer: task_service] [Method: CreateTask] Error: Title is required")
		return nil, ErrTitleRequired
	}
	if username == "" {
		s.logger.Errorln("[Layer: task_service] [Method: CreateTask] Error: UserName is required")
		return nil, ErrUserRequired
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := &models.Task{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
		Owner:     username,
	}
	s.tasks[s.nextID] = task
	s.logger.Infof("[Layer: task_service] [Method: CreateTask] Info: Task '%d' created for user '%s'", s.nextID, username)
	s.nextID++
	return task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, id int, title string, completed bool, username string) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists || task.Owner != username {
		s.logger.Warnf("[Layer: task_service] [Method: UpdateTask] Warning: Task '%d' not found or not owned by user '%s'", id, username)
		return nil, ErrTaskNotFound
	}

	task.Title = title
	task.Completed = completed
	s.logger.Infof("[Layer: task_service] [Method: UpdateTask] Info: Task '%d' updated for user '%s'", id, username)
	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id int, username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists || task.Owner != username {
		s.logger.Warnf("[Layer: task_service] [Method: DeleteTask] Warning: Task '%d' not found or not owned by user '%s'", id, username)
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	s.logger.Infof("[Layer: task_service] [Method: DeleteTask] Info: Task '%d' deleted for user '%s'", id, username)
	return nil
}
