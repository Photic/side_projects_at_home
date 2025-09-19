package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Device struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) Create(task *Task) error {
	query := `INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)`
	result, err := s.db.Exec(query, task.Title, task.Description, task.Completed)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	task.ID = int(id)
	
	// Get the created task with timestamps
	query = `SELECT created_at, updated_at FROM tasks WHERE id = ?`
	err = s.db.QueryRow(query, task.ID).Scan(&task.CreatedAt, &task.UpdatedAt)
	return err
}

func (s *TaskService) GetAll() ([]Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskService) Update(task *Task) error {
	query := `
		UPDATE tasks 
		SET title = ?, description = ?, completed = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := s.db.Exec(query, task.Title, task.Description, task.Completed, task.ID)
	return err
}

func (s *TaskService) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

type DeviceService struct {
	db *sql.DB
}

func NewDeviceService(db *sql.DB) *DeviceService {
	return &DeviceService{db: db}
}

func (s *DeviceService) Create(device *Device) error {
	query := `INSERT INTO devices (name, type, status, location) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, device.Name, device.Type, device.Status, device.Location)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	device.ID = int(id)
	
	// Get the created device with timestamps
	query = `SELECT created_at, updated_at FROM devices WHERE id = ?`
	err = s.db.QueryRow(query, device.ID).Scan(&device.CreatedAt, &device.UpdatedAt)
	return err
}

func (s *DeviceService) GetAll() ([]Device, error) {
	query := `SELECT id, name, type, status, location, created_at, updated_at FROM devices ORDER BY name`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		err := rows.Scan(&device.ID, &device.Name, &device.Type, &device.Status, &device.Location, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s *DeviceService) UpdateStatus(id int, status string) error {
	query := `UPDATE devices SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := s.db.Exec(query, status, id)
	return err
}