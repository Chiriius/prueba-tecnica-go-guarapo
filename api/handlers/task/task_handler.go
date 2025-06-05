package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"prueba_tecnica_go_guarapo/api/models"
	services "prueba_tecnica_go_guarapo/api/services/task"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TaskHandler interface {
	GetTasks(c *gin.Context)
	GetTask(c *gin.Context)
	CreateTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

type taskHandler struct {
	taskService services.TaskService
	logger      *logrus.Logger
}

func NewTaskHandler(taskService services.TaskService, logger *logrus.Logger) TaskHandler {
	return &taskHandler{
		taskService: taskService,
		logger:      logger,
	}
}

func (h *taskHandler) GetTasks(c *gin.Context) {
	username, _ := c.Get("username")
	tasks, err := h.taskService.GetTasksByUser(c.Request.Context(), username.(string))
	if err != nil {
		h.logger.Error("[Layer: task_handler] [Method: GetTasks] Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tareas"})
		return
	}
	if tasks == nil {
		tasks = []*models.Task{}
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *taskHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: GetTask] ID inválido: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	username, _ := c.Get("username")
	task, err := h.taskService.GetTaskByID(c.Request.Context(), id, username.(string))
	if err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: GetTask] No encontrada: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *taskHandler) CreateTask(c *gin.Context) {
	var req models.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: CreateTask] Datos inválidos: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		h.logger.Warn("[Layer: task_handler] [Method: CreateTask] Título vacío")
		c.JSON(http.StatusBadRequest, gin.H{"error": "El título no puede estar vacío"})
		return
	}
	username, _ := c.Get("username")
	newTask, err := h.taskService.CreateTask(c.Request.Context(), req.Title, username.(string))
	if err != nil {
		h.logger.Error("[Layer: task_handler] [Method: CreateTask] Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la tarea"})
		return
	}
	c.JSON(http.StatusCreated, newTask)
}

func (h *taskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: UpdateTask] ID inválido: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var req models.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: UpdateTask] Datos inválidos: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		h.logger.Warn("[Layer: task_handler] [Method: UpdateTask] Título vacío")
		c.JSON(http.StatusBadRequest, gin.H{"error": "El título no puede estar vacío"})
		return
	}
	username, _ := c.Get("username")
	updatedTask, err := h.taskService.UpdateTask(c.Request.Context(), id, req.Title, req.Completed, username.(string))
	if err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: UpdateTask] No encontrada: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

func (h *taskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: DeleteTask] ID inválido: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	username, _ := c.Get("username")
	err = h.taskService.DeleteTask(c.Request.Context(), id, username.(string))
	if err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: DeleteTask] No encontrada: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tarea eliminada exitosamente"})
}
