package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	authHandlers "prueba_tecnica_go_guarapo/api/handlers/auth"
	taskHandlers "prueba_tecnica_go_guarapo/api/handlers/task"
	"prueba_tecnica_go_guarapo/api/models"
	authServices "prueba_tecnica_go_guarapo/api/services/auth"
	taskServices "prueba_tecnica_go_guarapo/api/services/task"
	middleware "prueba_tecnica_go_guarapo/api/utils"
)

type Server struct {
	router *gin.Engine
	logger *logrus.Logger
	db     *gorm.DB
}

func NewServer(logger *logrus.Logger) *Server {
	router := gin.Default()
	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		logger.Fatal("No se pudo conectar a la base de datos:", err)
	}
	db.AutoMigrate(&models.Task{})
	return &Server{
		router: router,
		logger: logger,
		db:     db,
	}
}

func (s *Server) Start(addr string) {
	authService := authServices.NewAuthService(s.logger)
	taskService := taskServices.NewTaskService(s.db, s.logger)

	authHandler := authHandlers.NewAuthHandler(authService, s.logger)
	taskHandler := taskHandlers.NewTaskHandler(taskService, s.logger)

	api := s.router.Group("/api")
	{
		api.POST("/login", authHandler.Login)

		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware(authService))
		{
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.POST("", taskHandler.CreateTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}

	s.logger.Infof("[Layer: Server] [Method: Start] Server listened in %s", addr)
	s.router.Run(addr)
}
