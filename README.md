# Personal Blog

A web-based personal blog application built with Go that allows users to create, read, edit, and delete blog posts. The application provides a clean and intuitive HTML interface with secure authentication and JSON-based data persistence.

Project reference: https://roadmap.sh/projects/personal-blog

## Features

- Create, read, update, and delete blog posts (CRUD operations)
- Secure Basic Authentication for admin access
- Web-based interface with HTML templates
- Template caching for improved performance
- JSON file-based data storage
- Environment-based configuration management
- Concurrent-safe file operations with mutex locking
- Responsive HTML templates
- Comprehensive unit testing for core layers

## Supported Operations

| Category | Operations                                              |
| -------- | ------------------------------------------------------- |
| Public   | View all posts, View individual post                    |
| Admin    | Create post, Edit post, Delete post, Manage all posts   |
| Storage  | JSON file persistence, Auto-incrementing post IDs       |
| Auth     | Basic HTTP Authentication with configurable credentials |

## Requirements

- Go 1.22 or later

## Installation

Clone the repository:

```bash
git clone https://github.com/hazubeep/personal-blog.git
cd personal-blog
```

Install dependencies:

```bash
go mod tidy
```

## Configuration

Set environment variables to customize the application:

```bash
# Server port (default: :8080)
export PORT=:8080

# Admin credentials (defaults: admin/admin)
export ADMIN_USER=your_username
export ADMIN_PASS=your_password

# Storage path for posts (default: ./data/posts.json)
export STORAGE_PATH=./data/posts.json
```

## Usage

Run the application with:

```bash
go run ./cmd/web
```

The application will start a web server on `http://localhost:8080`. Open this URL in your browser to view the blog.

To access the admin panel, navigate to `http://localhost:8080/admin` and use the configured credentials.

Alternatively, build a binary first:

```bash
go build -o bin/web ./cmd/web
./bin/web
```

## Testing

Run the test suite to verify core functionality:

```bash
go test ./internal/modules/blog/...
```

For verbose test output with coverage:

```bash
go test -v ./internal/modules/blog/ -cover
```

Current test coverage: 36.2% (18 tests covering repository and service layers)

## Project Structure

```
personal-blog/
├── cmd/
│   └── web/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── middleware/
│   │   └── auth.go              # Basic authentication
│   └── modules/
│       └── blog/
│           ├── handler.go           # HTTP handlers
│           ├── service.go           # Business logic
│           ├── repository.go        # Data persistence
│           ├── model.go             # Data models
│           ├── repository_test.go   # Repository tests
│           └── service_test.go      # Service tests
├── templates/
│   ├── home.html         # Blog homepage
│   ├── article.html      # Post view
│   ├── admin.html        # Admin dashboard
│   ├── form.html         # Edit post form
│   └── add-form.html     # New post form
├── data/
│   └── posts.json        # Posts storage file
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # This file
```

## API Routes

### Public Routes

- `GET /` - View all blog posts
- `GET /posts/{id}` - View specific post

### Protected Routes (Requires Basic Auth)

- `GET /admin` - Admin dashboard
- `GET /new` - New post form
- `POST /new` - Create post
- `GET /edit/{id}` - Edit post form
- `POST /edit/{id}` - Update post
- `GET /delete/{id}` - Delete post

## Development

The project follows a layered architecture pattern:

1. Handler Layer - HTTP request processing
2. Service Layer - Business logic and validation
3. Repository Layer - Data persistence
4. Model Layer - Data structures

### Project Status

Current Version: 0.1.0 (MVP)

Suitable for:

- Learning Go web development
- Personal portfolio demonstration
- Self-hosted blogging platform

## Known Limitations

- Basic HTTP authentication only
- JSON file storage (single file)
- No markdown support
- No search functionality
- No post categories or tags
- No comment system
- Single admin user

## Future Improvements

- Add handler and middleware tests
- Implement input validation
- Add structured logging
- Database backend integration (SQLite, PostgreSQL)
- Markdown support for posts
- Post categories and tags system
- Full-text search functionality
- RSS feed generation
- Multiple user support
- Post scheduling capability

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
