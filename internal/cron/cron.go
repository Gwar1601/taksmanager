// internal/cron/cron.go
package cron

import (
	"log"
	"time"

	"github.com/gwar1601/taskmanager/internal/services"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	taskService services.TaskService
}

func NewCronService(taskService services.TaskService) *CronService {
	return &CronService{taskService}
}

func (c *CronService) Start() {
	scheduler := cron.New(cron.WithSeconds())
	_, err := scheduler.AddFunc("@daily", func() {
		log.Println("Running daily task to update task statuses")
		c.updateTaskStatuses()
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}
	scheduler.Start()
}

func (c *CronService) updateTaskStatuses() {
	// Pobierz zadania, które przekroczyły termin i nie są ukończone
	today := time.Now().Format("2006-01-02")
	tasks, err := c.taskService.GetTasksByDueDate(0, today) // 0 oznacza wszystkich użytkowników
	if err != nil {
		log.Printf("Error fetching tasks for auto-completion: %v", err)
		return
	}

	for _, task := range tasks {
		if task.Status != "Ukończone" {
			if err := c.taskService.UpdateTaskStatus(task.ID, "Ukończone"); err != nil {
				log.Printf("Error updating task status for task ID %d: %v", task.ID, err)
			} else {
				log.Printf("Task ID %d has been marked as completed", task.ID)
			}
		}
	}
}
