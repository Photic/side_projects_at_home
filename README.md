# Home Automation Dashboard

A collection of smaller web pages for home automation tasks built with Go, templ, and SQLite.

## Features

- **Task Management**: Create, update, and manage home automation tasks
- **Device Control**: Monitor and control various home automation devices
- **Modern UI**: Clean, responsive interface using Tailwind CSS
- **Real-time Updates**: HTMX-powered dynamic updates without page refreshes
- **SQLite Database**: Lightweight, file-based database for data persistence

## Tech Stack

- **Backend**: Go 1.24+ with Gorilla Mux router
- **Frontend**: templ templating engine with Tailwind CSS
- **Database**: SQLite with pure Go driver
- **JavaScript**: HTMX for dynamic interactions
- **Build Tool**: Make for easy development workflow

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Make (optional, but recommended)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Photic/side_projects_at_home.git
cd side_projects_at_home
```

2. Install templ CLI tool:
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

3. Setup the project:
```bash
make setup
```

### Running the Application

1. Build and run:
```bash
make run
```

2. Or for development with auto-reload:
```bash
make dev
```

3. Visit http://localhost:8080

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   └── db.go              # Database connection and migrations
│   ├── handlers/
│   │   └── handlers.go        # HTTP handlers
│   └── models/
│       └── models.go          # Data models and services
├── templates/
│   ├── templates.templ        # templ templates
│   └── templates_templ.go     # Generated Go code
├── static/
│   └── js/
│       └── app.js            # Frontend JavaScript
├── data/                     # SQLite database files (created automatically)
├── Makefile                  # Build and development tasks
└── README.md
```

## API Endpoints

### Tasks
- `GET /` - Home page with task list
- `GET /tasks/new` - New task form
- `POST /tasks` - Create new task
- `POST /tasks/{id}/toggle` - Toggle task completion
- `DELETE /tasks/{id}` - Delete task

### Devices
- `GET /devices` - Device management page
- `GET /devices/new` - New device form
- `POST /devices` - Create new device
- `POST /devices/{id}/status` - Update device status

## Available Make Commands

- `make build` - Build the application
- `make run` - Build and run the application
- `make dev` - Run in development mode with auto-reload
- `make generate` - Generate templ templates
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make fmt` - Format code
- `make setup` - Setup project dependencies

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure code is formatted
5. Submit a pull request

## License

This project is open source and available under the MIT License.
