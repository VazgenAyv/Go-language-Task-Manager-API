package main

import (
	"log"
	"net/http"

	"github.com/ht21992/go-task-manager/database"
	"github.com/ht21992/go-task-manager/handlers"
	"github.com/ht21992/go-task-manager/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	r := mux.NewRouter()

	// Auth route
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/register", handlers.Register).Methods("POST")

	// Task routes (protected)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware) // Apply JWT middleware
	api.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")

	// CORS middleware for frontend compatibility (e.g., React)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // React Vite origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", handler)
}
