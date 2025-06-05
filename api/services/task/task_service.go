package services

import (
	"context"
	"errors"
	"prueba_tecnica_go_guarapo/api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskService interface {
	GetTasksByUser(ctx context.Context, username string) ([]*models.Task, error)
	GetTaskByID(ctx context.Context, id int, username string) (*models.Task, error)
	CreateTask(ctx context.Context, title string, username string) (*models.Task, error)
	UpdateTask(ctx context.Context, id int, title string, completed bool, username string) (*models.Task, error)
	DeleteTask(ctx context.Context, id int, username string) error
}

type taskService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewTaskService(db *gorm.DB, logger *logrus.Logger) TaskService {
	return &taskService{
		db:     db,
		logger: logger,
	}
}

func (s *taskService) GetTasksByUser(ctx context.Context, username string) ([]*models.Task, error) {
	var tasks []*models.Task
	if err := s.db.WithContext(ctx).Where("owner = ?", username).Find(&tasks).Error; err != nil {
		s.logger.Error("[Layer: task_service] [Method: GetTasksByUser] Error: ", err)
		return nil, err
	}
	s.logger.Infof("[Layer: task_service] [Method: GetTasksByUser] Info: User '%s' requested their tasks", username)
	return tasks, nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id int, username string) (*models.Task, error) {
	var task models.Task
	if err := s.db.WithContext(ctx).Where("id = ? AND owner = ?", id, username).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warnf("[Layer: task_service] [Method: GetTaskByID] Warning: Task '%d' not found or not owned by user '%s'", id, username)
			return nil, ErrTaskNotFound
		}
		s.logger.Error("[Layer: task_service] [Method: GetTaskByID] Error: ", err)
		return nil, err
	}
	s.logger.Infof("[Layer: task_service] [Method: GetTaskByID] Info: Task '%d' retrieved for user '%s'", id, username)
	return &task, nil
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

	task := &models.Task{
		Title:     title,
		Completed: false,
		Owner:     username,
	}
	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		s.logger.Error("[Layer: task_service] [Method: CreateTask] Error: ", err)
		return nil, err
	}
	s.logger.Infof("[Layer: task_service] [Method: CreateTask] Info: Task '%d' created for user '%s'", task.ID, username)
	return task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, id int, title string, completed bool, username string) (*models.Task, error) {
	var task models.Task
	if err := s.db.WithContext(ctx).Where("id = ? AND owner = ?", id, username).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warnf("[Layer: task_service] [Method: UpdateTask] Warning: Task '%d' not found or not owned by user '%s'", id, username)
			return nil, ErrTaskNotFound
		}
		s.logger.Error("[Layer: task_service] [Method: UpdateTask] Error: ", err)
		return nil, err
	}

	task.Title = title
	task.Completed = completed

	if err := s.db.WithContext(ctx).Save(&task).Error; err != nil {
		s.logger.Error("[Layer: task_service] [Method: UpdateTask] Error: ", err)
		return nil, err
	}
	s.logger.Infof("[Layer: task_service] [Method: UpdateTask] Info: Task '%d' updated for user '%s'", id, username)
	return &task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id int, username string) error {
	var task models.Task
	if err := s.db.WithContext(ctx).Where("id = ? AND owner = ?", id, username).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warnf("[Layer: task_service] [Method: DeleteTask] Warning: Task '%d' not found or not owned by user '%s'", id, username)
			return ErrTaskNotFound
		}
		s.logger.Error("[Layer: task_service] [Method: DeleteTask] Error: ", err)
		return err
	}

	if err := s.db.WithContext(ctx).Delete(&task).Error; err != nil {
		s.logger.Error("[Layer: task_service] [Method: DeleteTask] Error: ", err)
		return err
	}
	s.logger.Infof("[Layer: task_service] [Method: DeleteTask] Info: Task '%d' deleted for user '%s'", id, username)
	return nil
}
