package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ht21992/go-task-manager/database"
	"github.com/ht21992/go-task-manager/middleware"
	"github.com/ht21992/go-task-manager/models"
)

// GET /tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, description, completed, maintask FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.MainTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// POST /tasks
func CreateTask(w http.ResponseWriter, r *http.Request) {

	role := r.Context().Value(middleware.RoleContextKey)
	if role != "admin" {
		http.Error(w, "Only admin can create tasks", http.StatusForbidden)
		return
	}

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	statement, err := database.DB.Prepare("INSERT INTO tasks (title, description, completed, maintask) VALUES (?,?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println((task.MainTask))
	result, err := statement.Exec(task.Title, task.Description, task.Completed, task.MainTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	task.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

// DELETE /tasks/{id}
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	role := r.Context().Value(middleware.RoleContextKey)
	if role != "admin" {
		http.Error(w, "Only admin can delete tasks", http.StatusForbidden)
		return
	}

	statement, err := database.DB.Prepare("DELETE FROM tasks WHERE id = ?")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task Deleted"})

}

// PUT /tasks/{id}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	role := r.Context().Value(middleware.RoleContextKey)
	if role != "admin" {
		http.Error(w, "Only admin can update tasks", http.StatusForbidden)
		return
	}

	// Fetch the existing task
	var existing models.Task
	row := database.DB.QueryRow("SELECT id, title, description, completed FROM tasks WHERE id = ?", id)
	err := row.Scan(&existing.ID, &existing.Title, &existing.Description, &existing.Completed)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Decode the request body
	var incoming models.Task
	err = json.NewDecoder(r.Body).Decode(&incoming)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Invalid input for task", http.StatusBadRequest)
		return
	}

	// Use existing values if fields are nil (not included in JSON)
	if incoming.Title != nil {
		existing.Title = incoming.Title
	}
	if incoming.Description != nil {
		existing.Description = incoming.Description
	}
	if incoming.Completed != nil {
		existing.Completed = incoming.Completed
	}

	// Update the task
	statement, err := database.DB.Prepare(`
        UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?
    `)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec(existing.Title, existing.Description, existing.Completed, id)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task Updated"})
}
