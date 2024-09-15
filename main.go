package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Task represents a task with a description, status, and due date
type Task struct {
	ID          int
	Description string
	Status      string
	DueDate     time.Time
}

// TaskManager manages a list of tasks
type TaskManager struct {
	tasks []Task
	mutex sync.Mutex
}

// AddTask adds a new task to the task manager
func (tm *TaskManager) AddTask(description string, dueDate time.Time) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	task := Task{
		ID:          len(tm.tasks) + 1,
		Description: description,
		Status:      "Pending",
		DueDate:     dueDate,
	}
	tm.tasks = append(tm.tasks, task)
}

// CompleteTask marks a task as completed
func (tm *TaskManager) CompleteTask(taskID int) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	for i, task := range tm.tasks {
		if task.ID == taskID {
			tm.tasks[i].Status = "Completed"
			break
		}
	}
}

// ListTasks returns a list of all tasks
func (tm *TaskManager) ListTasks() []Task {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	return tm.tasks
}

// ProcessTasksWithWorkers processes tasks with a given number of workers
func (tm *TaskManager) ProcessTasksWithWorkers(numWorkers int) {
	tasks := tm.ListTasks()
	var wg sync.WaitGroup

	taskChan := make(chan Task, len(tasks))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range taskChan {
				fmt.Printf("Worker %d starting task: %s\n", workerID, task.Description)
				time.Sleep(2 * time.Second) // Simulate task processing
				tm.CompleteTask(task.ID)
				fmt.Printf("Worker %d completed task: %s\n", workerID, task.Description)
			}
		}(i + 1)
	}

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	wg.Wait()
	fmt.Println("All tasks completed!")
}

var taskManager = TaskManager{}

// Web handler to display tasks in the browser
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Task Manager</h1>")
	tasks := taskManager.ListTasks()
	for _, task := range tasks {
		fmt.Fprintf(w, "Task ID: %d, Description: %s, Status: %s, Due: %s<br>",
			task.ID, task.Description, task.Status, task.DueDate.Format("2006-01-02"))
	}
}

// Main function to run the web server
func main() {
	// Adding some example tasks for testing purposes
	taskManager.AddTask("Learn Go", time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC))
	taskManager.AddTask("Learn Concurrency", time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC))
	taskManager.AddTask("Build Task Manager", time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC))

	// HTTP handler
	http.HandleFunc("/", handler)

	// Get the port from the environment variable, if not set default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT environment variable is not set
	}

	// Start the HTTP server
	fmt.Println("Listening on port:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
