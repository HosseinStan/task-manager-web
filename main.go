package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux" // Import Gorilla Mux for handling routes
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	DueDate     time.Time `json:"due_date"`
}

type TaskManager struct {
	sync.Mutex
	Tasks []Task
}

// Initialize a task manager
var tm = TaskManager{Tasks: []Task{}}

// List all tasks
func listTasks(w http.ResponseWriter, r *http.Request) {
	tm.Lock()
	defer tm.Unlock()

	json.NewEncoder(w).Encode(tm.Tasks)
}

// Add a new task
func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)

	tm.Lock()
	defer tm.Unlock()

	task.ID = len(tm.Tasks) + 1
	task.CreatedAt = time.Now()

	tm.Tasks = append(tm.Tasks, task)
	json.NewEncoder(w).Encode(task)
}

// Mark a task as completed
func completeTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	tm.Lock()
	defer tm.Unlock()

	for i, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks[i].Completed = true
			json.NewEncoder(w).Encode(tm.Tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// Process tasks with workers
func processTasks(w http.ResponseWriter, r *http.Request) {
	// For simplicity, simulate processing tasks
	tm.Lock()
	defer tm.Unlock()

	for i := range tm.Tasks {
		if !tm.Tasks[i].Completed {
			tm.Tasks[i].Completed = true
		}
	}

	fmt.Fprintf(w, "All tasks processed")
}

func main() {
	// Initialize Router
	r := mux.NewRouter()

	// Define Routes
	r.HandleFunc("/tasks", listTasks).Methods("GET")
	r.HandleFunc("/tasks", addTask).Methods("POST")
	r.HandleFunc("/tasks/{id}/complete", completeTask).Methods("PUT")
	r.HandleFunc("/tasks/process", processTasks).Methods("GET")

	// Start server
	http.ListenAndServe(":8080", r)
}
