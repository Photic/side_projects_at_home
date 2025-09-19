package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Photic/side_projects_at_home/internal/models"
	"github.com/Photic/side_projects_at_home/templates"
	"github.com/gorilla/mux"
)

type Handler struct {
	taskService   *models.TaskService
	deviceService *models.DeviceService
}

func NewHandler(taskService *models.TaskService, deviceService *models.DeviceService) *Handler {
	return &Handler{
		taskService:   taskService,
		deviceService: deviceService,
	}
}

func (h *Handler) TasksHome(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetAll()
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := templates.TasksPage(tasks)
	component.Render(r.Context(), w)
}

func (h *Handler) NewTaskForm(w http.ResponseWriter, r *http.Request) {
	component := templates.NewTaskForm()
	component.Render(r.Context(), w)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Completed:   false,
	}

	if err := h.taskService.Create(task); err != nil {
		log.Printf("Error creating task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := templates.TaskCard(*task)
	component.Render(r.Context(), w)

	// Close modal
	w.Header().Set("HX-Trigger", "closeModal")
}

func (h *Handler) ToggleTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Get current task to toggle completion
	tasks, err := h.taskService.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var task *models.Task
	for _, t := range tasks {
		if t.ID == id {
			task = &t
			break
		}
	}

	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task.Completed = !task.Completed
	if err := h.taskService.Update(task); err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := templates.TaskCard(*task)
	component.Render(r.Context(), w)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.taskService.Delete(id); err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DevicesHome(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceService.GetAll()
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := templates.DevicesPage(devices)
	component.Render(r.Context(), w)
}

func (h *Handler) NewDeviceForm(w http.ResponseWriter, r *http.Request) {
	component := templates.NewDeviceForm()
	component.Render(r.Context(), w)
}

func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	device := &models.Device{
		Name:     r.FormValue("name"),
		Type:     r.FormValue("type"),
		Status:   "offline",
		Location: r.FormValue("location"),
	}

	if err := h.deviceService.Create(device); err != nil {
		log.Printf("Error creating device: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := templates.DeviceCard(*device)
	component.Render(r.Context(), w)

	// Close modal
	w.Header().Set("HX-Trigger", "closeModal")
}

func (h *Handler) UpdateDeviceStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	if err := h.deviceService.UpdateStatus(id, status); err != nil {
		log.Printf("Error updating device status: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get updated device to return
	devices, err := h.deviceService.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var device *models.Device
	for _, d := range devices {
		if d.ID == id {
			device = &d
			break
		}
	}

	if device == nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	component := templates.DeviceCard(*device)
	component.Render(r.Context(), w)
}