# Go Task Manager API

A lightweight RESTful Task Manager built with Go, SQLite, and JWT Authentication. This backend is ready to connect with a React (Vite) frontend or any frontend framework.

---

## Features

- JWT-based authentication (login required to access tasks)
- Role-based access ready (store `username` and `role` in JWT)
- Secure password hashing using `bcrypt`
- SQLite support with safe query practices (prevents SQL injection)
- CORS enabled for frontend communication
- Full CRUD operations on tasks
- Clean and modular code structure

---

## Tech Stack

- Go (Golang)
- SQLite (via `modernc.org/sqlite`)
- JWT (`github.com/golang-jwt/jwt/v5`)
- Gorilla Mux (router)
- Bcrypt for password hashing

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/ht21992/Go-language-Task-Manager-API.git
cd Go-language-Task-Manager-API
```

2. Install Dependencies

```
go mod tidy
```

3. Run the Server

```
go run main.go
```

The API runs by default at: http://localhost:8080

⸻

API Endpoints

### Auth

#### POST /register — Add new user

Payload:

```
{ "username": "user", "password": "pass", "role": "admin" }
```

#### POST /login — Authenticate and get a JWT token

Payload:

```
{ "username": "user", "password": "pass" }
```

### Tasks (Requires JWT in Authorization: Bearer <token> header)

    •	GET /api/tasks — Get all tasks
    •	POST /api/tasks — Create a new task
    •	PUT /api/tasks/{id} — Update task (supports partial update)
    •	DELETE /api/tasks/{id} — Delete task

Create Task payload

```
{"title":"new task","description":"new desc","completed":true}
```

⸻

Environment Configuration

- No environment file needed for development.
- You can hardcode or extract the JWT secret from an .env file using os.Getenv("JWT_SECRET") in production.

⸻

Security Features

- Passwords are hashed with bcrypt
- SQL Injection protection via ? prepared statements
- JWT tokens validate both username and role
- CORS enabled for React or other frontend integrations
