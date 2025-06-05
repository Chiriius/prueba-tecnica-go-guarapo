package handlers

import (
	"net/http"
	"strconv"

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

// CreateTask godoc
// @Summary      Crear tarea
// @Description  Crea una nueva tarea para el usuario autenticado
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body models.CreateTaskRequest true "Datos de la tarea"
// @Success      201 {object} models.TaskResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks [post]
func (h *taskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: CreateTask] Datos inválidos: ", err)
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
	resp := models.TaskResponse{
		ID:        newTask.ID,
		Title:     newTask.Title,
		Completed: newTask.Completed,
		Owner:     newTask.Owner,
	}
	c.JSON(http.StatusCreated, resp)
}

// UpdateTask godoc
// @Summary      Actualizar tarea
// @Description  Actualiza una tarea existente del usuario autenticado
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la tarea"
// @Param        request body models.UpdateTaskRequest true "Datos de la tarea"
// @Success      200 {object} models.TaskResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks/{id} [put]
func (h *taskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: UpdateTask] ID inválido: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("[Layer: task_handler] [Method: UpdateTask] Datos inválidos: ", err)
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
	resp := models.TaskResponse{
		ID:        updatedTask.ID,
		Title:     updatedTask.Title,
		Completed: updatedTask.Completed,
		Owner:     updatedTask.Owner,
	}
	c.JSON(http.StatusOK, resp)
}

// GetTasks godoc
// @Summary      Listar tareas
// @Description  Obtiene todas las tareas del usuario autenticado
// @Tags         tasks
// @Produce      json
// @Success      200 {array} models.TaskResponse
// @Failure      500 {object} map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks [get]
func (h *taskHandler) GetTasks(c *gin.Context) {
	username, _ := c.Get("username")
	tasks, err := h.taskService.GetTasksByUser(c.Request.Context(), username.(string))
	if err != nil {
		h.logger.Error("[Layer: task_handler] [Method: GetTasks] Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tareas"})
		return
	}
	var resp []models.TaskResponse
	for _, t := range tasks {
		resp = append(resp, models.TaskResponse{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			Owner:     t.Owner,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// GetTask godoc
// @Summary      Obtener tarea
// @Description  Obtiene una tarea específica del usuario autenticado
// @Tags         tasks
// @Produce      json
// @Param        id path int true "ID de la tarea"
// @Success      200 {object} models.TaskResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks/{id} [get]
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
	resp := models.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		Completed: task.Completed,
		Owner:     task.Owner,
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteTask godoc
// @Summary      Eliminar tarea
// @Description  Elimina una tarea del usuario autenticado
// @Tags         tasks
// @Produce      json
// @Param        id path int true "ID de la tarea"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks/{id} [delete]
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
