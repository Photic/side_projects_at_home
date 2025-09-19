package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Photic/side_projects_at_home/internal/database"
	"github.com/Photic/side_projects_at_home/internal/handlers"
	"github.com/Photic/side_projects_at_home/internal/models"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db, err := database.New()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize services
	taskService := models.NewTaskService(db.DB)
	deviceService := models.NewDeviceService(db.DB)

	// Initialize handlers
	handler := handlers.NewHandler(taskService, deviceService)

	// Setup routes
	r := mux.NewRouter()

	// Task routes
	r.HandleFunc("/", handler.TasksHome).Methods("GET")
	r.HandleFunc("/tasks/new", handler.NewTaskForm).Methods("GET")
	r.HandleFunc("/tasks", handler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}/toggle", handler.ToggleTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", handler.DeleteTask).Methods("DELETE")

	// Device routes
	r.HandleFunc("/devices", handler.DevicesHome).Methods("GET")
	r.HandleFunc("/devices/new", handler.NewDeviceForm).Methods("GET")
	r.HandleFunc("/devices", handler.CreateDevice).Methods("POST")
	r.HandleFunc("/devices/{id}/status", handler.UpdateDeviceStatus).Methods("POST")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Visit: http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}