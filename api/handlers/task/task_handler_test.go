package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"prueba_tecnica_go_guarapo/api/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskHandler_CreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testScenarios := []struct {
		testName       string
		requestBody    interface{}
		mockSetup      func(*mockTaskService)
		username       string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:    "Crear tarea exitosa",
			requestBody: models.CreateTaskRequest{Title: "Nueva tarea"},
			username:    "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("CreateTask", mock.Anything, "Nueva tarea", "user1").
					Return(&models.Task{Title: "Nueva tarea", Completed: false, Owner: "user1"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"title":"Nueva tarea"`,
		},
		{
			testName:       "Título vacío",
			requestBody:    models.CreateTaskRequest{Title: ""},
			username:       "user1",
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"El título no puede estar vacío"`,
		},
		{
			testName:       "JSON inválido",
			requestBody:    `{bad json}`,
			username:       "user1",
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"El título no puede estar vacío"`,
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mockTaskService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			logger := logrus.New()
			handler := NewTaskHandler(mockService, logger)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("username", tt.username)
			})
			router.POST("/tasks", handler.CreateTask)

			var reqBody []byte
			var err error
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testScenarios := []struct {
		testName       string
		id             string
		requestBody    interface{}
		mockSetup      func(*mockTaskService)
		username       string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:    "Actualizar tarea exitosa",
			id:          "1",
			requestBody: models.UpdateTaskRequest{Title: "Actualizada", Completed: true},
			username:    "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("UpdateTask", mock.Anything, 1, "Actualizada", true, "user1").
					Return(&models.Task{Title: "Actualizada", Completed: true, Owner: "user1"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Actualizada"`,
		},
		{
			testName:       "ID inválido",
			id:             "abc",
			requestBody:    models.UpdateTaskRequest{Title: "Actualizada", Completed: true},
			username:       "user1",
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"ID inválido"`,
		},
		{
			testName:       "Título vacío",
			id:             "1",
			requestBody:    models.UpdateTaskRequest{Title: "", Completed: true},
			username:       "user1",
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"El título no puede estar vacío"`,
		},
		{
			testName:    "Tarea no encontrada",
			id:          "1",
			requestBody: models.UpdateTaskRequest{Title: "Actualizada", Completed: true},
			username:    "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("UpdateTask", mock.Anything, 1, "Actualizada", true, "user1").
					Return((*models.Task)(nil), errors.New("Tarea no encontrada"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"Tarea no encontrada"`,
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mockTaskService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			logger := logrus.New()
			handler := NewTaskHandler(mockService, logger)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("username", tt.username)
			})
			router.PUT("/tasks/:id", handler.UpdateTask)

			var reqBody []byte
			var err error
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest(http.MethodPut, "/tasks/"+tt.id, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestTaskHandler_GetTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testScenarios := []struct {
		testName       string
		id             string
		mockSetup      func(*mockTaskService)
		username       string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Obtener tarea exitosa",
			id:       "1",
			username: "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("GetTaskByID", mock.Anything, 1, "user1").
					Return(&models.Task{Title: "Tarea", Completed: false, Owner: "user1"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Tarea"`,
		},
		{
			testName:       "ID inválido",
			id:             "abc",
			username:       "user1",
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"ID inválido"`,
		},
		{
			testName: "Tarea no encontrada",
			id:       "1",
			username: "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("GetTaskByID", mock.Anything, 1, "user1").
					Return((*models.Task)(nil), errors.New("Tarea no encontrada"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"Tarea no encontrada"`,
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mockTaskService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			logger := logrus.New()
			handler := NewTaskHandler(mockService, logger)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("username", tt.username)
			})
			router.GET("/tasks/:id", handler.GetTask)

			req, _ := http.NewRequest(http.MethodGet, "/tasks/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestTaskHandler_GetTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testScenarios := []struct {
		testName       string
		mockSetup      func(*mockTaskService)
		username       string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Obtener lista de tareas",
			username: "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("GetTasksByUser", mock.Anything, "user1").
					Return([]*models.Task{
						{Title: "Tarea 1", Completed: false, Owner: "user1"},
						{Title: "Tarea 2", Completed: true, Owner: "user1"},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Tarea 1"`,
		},
		{
			testName: "Error al obtener tareas",
			username: "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("GetTasksByUser", mock.Anything, "user1").
					Return([]*models.Task{}, errors.New("Error al obtener tareas"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"Error al obtener tareas"`,
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mockTaskService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			logger := logrus.New()
			handler := NewTaskHandler(mockService, logger)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("username", tt.username)
			})
			router.GET("/tasks", handler.GetTasks)

			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testScenarios := []struct {
		testName       string
		id             string
		mockSetup      func(*mockTaskService)
		username       string
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "Eliminar tarea exitosa",
			id:       "1",
			username: "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("DeleteTask", mock.Anything, 1, "user1").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"Tarea eliminada exitosamente"`,
		},
		{
			testName:       "ID inválido",
			id:             "abc",
			username:       "user1",
			mockSetup:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"ID inválido"`,
		},
		{
			testName: "Tarea no encontrada",
			id:       "1",
			username: "user1",
			mockSetup: func(m *mockTaskService) {
				m.On("DeleteTask", mock.Anything, 1, "user1").
					Return(errors.New("Tarea no encontrada"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"Tarea no encontrada"`,
		},
	}

	for _, tt := range testScenarios {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(mockTaskService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			logger := logrus.New()
			handler := NewTaskHandler(mockService, logger)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("username", tt.username)
			})
			router.DELETE("/tasks/:id", handler.DeleteTask)

			req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
