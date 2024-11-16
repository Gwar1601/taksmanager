// internal/handlers/tasks.go
package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"taskmanager/internal/models"
	"taskmanager/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{service}
}

// CreateTaskHandler tworzy nowe zadanie
// @Summary Create a new task
// @Description Create a new task with a due date
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Create Task"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (h *TaskHandler) CreateTaskHandler(c *gin.Context) {
	// Implementacja funkcji
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	task.Date = time.Now()
	task.DueDate = task.Date.AddDate(0, 0, 7) // Termin wykonania na 7 dni

	if err := h.service.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

// GetTasksHandler pobiera wszystkie zadania użytkownika
func (h *TaskHandler) GetTasksHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	tasks, err := h.service.GetTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTaskByIDHandler pobiera zadanie po ID
func (h *TaskHandler) GetTaskByIDHandler(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.service.GetTaskByID(taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// UpdateTaskStatusHandler aktualizuje status zadania
func (h *TaskHandler) UpdateTaskStatusHandler(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.UpdateTaskStatus(taskID, input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task status updated"})
}

// DeleteTaskHandler usuwa zadanie
func (h *TaskHandler) DeleteTaskHandler(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.service.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// GetTasksByDueDateHandler pobiera zadania na określony dzień
func (h *TaskHandler) GetTasksByDueDateHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	dueDate := c.Param("due_date") // Format: YYYY-MM-DD

	tasks, err := h.service.GetTasksByDueDate(userID, dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TaskHandler) GetTasksForToday(c *gin.Context) {
	userID := c.GetInt("user_id")
	today := time.Now().Format("2006-01-02")
	tasks, err := h.service.GetTasksByDueDate(userID, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TaskHandler) GetTasksForTomorrow(c *gin.Context) {
	userID := c.GetInt("user_id")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	tasks, err := h.service.GetTasksByDueDate(userID, tomorrow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TaskHandler) GetTasksByDueDateHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	dueDate := c.Param("due_date") // Format: YYYY-MM-DD

	// Walidacja formatu daty
	_, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	tasks, err := h.service.GetTasksByDueDate(userID, dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
