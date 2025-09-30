## Notes App (Go + Gin + GORM + JWT)

A minimal REST API for a Notes application built with Go, Gin, GORM, PostgreSQL, and cookie-based JWT auth. It supports signup/login (hashed passwords), refresh tokens in HttpOnly cookies, and authenticated Notes CRUD.

---

### ✨ Features
- **JWT auth (access + refresh)**: Access token returned in JSON; refresh token stored as HttpOnly cookie `refresh_token`.
- **Refresh & logout**: `/refreshToken` issues a new access token; `/logout` clears the refresh cookie.
- **Protected Notes CRUD**: All notes endpoints are mounted under `/api` and guarded by middleware.
- **PostgreSQL via GORM**: Auto-migration for `User`, `Note`, `RefreshToken`.
- **12-factor config**: Environment-driven config.

---

### 🧰 Tech Stack
- **Language**: Go
- **Framework**: `github.com/gin-gonic/gin`
- **ORM**: `gorm.io/gorm` + `gorm.io/driver/postgres`
- **Auth**: `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`
- **Env**: `github.com/joho/godotenv`

---

### 📦 Project Structure
```text
/Notes-App
  ├─ controllers/
  │   ├─ autho.go                 # Signup, Login, Logout
  │   ├─ controllers.go           # User create, Notes CRUD handlers
  │   └─ refreshControllers.go    # Issue new access token from refresh cookie
  ├─ intializers/
  │   ├─ datbase.go               # ConnectToDB (Postgres via GORM)
  │   ├─ LoadEnvVar.go            # Load .env via godotenv
  │   └─ SyncDatabase.go          # AutoMigrate(User, Note, RefreshToken)
  ├─ middleware/
  │   └─ Reqautho.go              # RequireAuth (checks refresh_token cookie)
  ├─ models/
  │   ├─ notes.go                 # Note model
  │   ├─ refresh.go               # RefreshToken model
  │   └─ user.go                  # User model
  ├─ utils/
  │   └─ jwt.go                   # JWT helpers (access & refresh)
  ├─ main.go                      # Router, routes, bootstrapping
  ├─ Dockerfile                   # App container image
  ├─ Docker-compose.yml           # Postgres + app
  ├─ go.mod
  └─ go.sum
```

---

### 🗄️ Database Models (simplified)
```go
// models/user.go
type User struct {
  ID       uint `gorm:"primaryKey"`
  Name     string
  Email    string `gorm:"uniqueIndex"`
  Password string
  Notes    []Note
}

// models/notes.go
type Note struct {
  ID      uint `gorm:"primaryKey"`
  Title   string
  Content string
  UserID  uint
}

// models/refresh.go
type RefreshToken struct {
  ID        uint `gorm:"primaryKey"`
  Token     string
  UserID    uint
  ExpiresAt time.Time
}
```

---

### 🔐 Authentication Flow
- `POST /signup`: Creates a user (bcrypt-hashed password).
- `POST /login`: Returns a short-lived access token in JSON and sets `refresh_token` cookie.
- `POST /refreshToken`: Uses `refresh_token` cookie to return a new access token in JSON.
- `GET /logout`: Clears the `refresh_token` cookie and deletes it from storage.

Important:
- The middleware currently guards routes by requiring the `refresh_token` cookie. You do not need to send the access token to call protected routes.

JWT details:
- HS256, claims include `user_id`, `exp`, and `type` (either `access` or `refresh`).
- Signing secret comes from `SECRET` env var.

---

### 🌐 REST API

Base URL: `http://localhost:<PORT>` (default `8080`)

#### Auth
- `POST /signup`
  - Body (form or JSON): `{ "Email": string, "Password": string }`
  - Response: `201 { user }`

- `POST /login`
  - Body: `{ "Email": string, "Password": string }`
  - Sets cookie: `refresh_token=<JWT>` (HttpOnly)
  - Response: `200 { "access token": string }`

- `POST /refreshToken`
  - Reads cookie: `refresh_token`
  - Response: `200 { "access_token": string }`

- `GET /logout`
  - Clears cookie: `refresh_token`
  - Response: `200 { "message": string }`

#### Users
- `POST /user` (public)
  - Creates a user record (primarily for testing).

#### Notes (protected under `/api`, RequireAuth)
- `POST /api/notes`
  - Body: `{ "title": string, "content": string, "userID": number }`
  - Requires `refresh_token` cookie

- `GET /api/notes`
- `GET /api/notes/:id`
- `PUT /api/notes/:id`
- `DELETE /api/notes/:id`
- `DELETE /api/user/:id`

---

### 🧪 Quick Start (local)

#### 1) Prerequisites
- Go 1.24+
- PostgreSQL running and reachable

#### 2) Environment
Create a `.env` file in the project root:
```env
PORT=8080
DB_url=postgres://USER:PASSWORD@HOST:5432/DBNAME?sslmode=disable
SECRET=your-very-strong-jwt-secret
```

#### 3) Install deps
```bash
go mod tidy
```

#### 4) Run the server
```bash
go run main.go
```
You should see logs like:
- "Initializing database"
- "welcome to notes app"
- and a successful DB connection

GORM auto-migrates `User`, `Note`, and `RefreshToken` on startup.

---

### 🐳 Run with Docker Compose

This repo includes `Dockerfile` and `Docker-compose.yml` to run Postgres and the app together.

1) Create `.env` with both app and database variables:
```env
# App
PORT=3500            # set to 3500 to match Docker port mapping
DB_url=postgres://postgres:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable
SECRET=your-very-strong-jwt-secret

# Postgres (used by docker-compose)
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=notes
```

2) Start services:
```bash
docker compose up --build
```

3) Access the API at `http://localhost:3500`.

Notes:
- Compose maps `3500:3500`. Ensure `PORT=3500` in `.env` so the app listens on the same port inside the container.
- The provided `Dockerfile` builds the app; you may want to adjust it to name the binary and `ENTRYPOINT` explicitly in production.

---

### 🧰 cURL Examples
```bash
# Signup
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"Email":"alice@example.com","Password":"SuperSecret123"}'

# Login: returns access token JSON and sets refresh_token cookie
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"Email":"alice@example.com","Password":"SuperSecret123"}' -i

# Call protected notes (cookie is stored by your client; with curl you may need -c/-b to persist)
curl -X POST http://localhost:8080/api/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"My Note","content":"Hello","userID":1}' \
  -b cookies.txt -c cookies.txt

# List notes (protected)
curl http://localhost:8080/api/notes -b cookies.txt -c cookies.txt

# Refresh access token (uses refresh_token cookie set by login)
curl -X POST http://localhost:8080/refreshToken -b cookies.txt -c cookies.txt

# Logout (clears refresh_token cookie)
curl http://localhost:8080/logout -b cookies.txt -c cookies.txt
```

---

### ⚙️ Configuration
- **PORT**: Server port (default `8080`)
- **DB_url**: Postgres DSN consumed by GORM
- **SECRET**: JWT signing secret

---

### 🔒 Security Notes
- Refresh token is stored in an HttpOnly cookie to mitigate XSS token theft.
- In production, set cookie `Secure` and serve over HTTPS.
- Use a strong, rotated `SECRET`.

---

### 🚀 Roadmap Ideas
- Scope notes by authenticated user (current middleware checks refresh cookie only)
- Add pagination and search to `GET /api/notes`
- Add OpenAPI/Swagger docs

---

### 🙌 Contributing
1. Fork the repo and create a feature branch
2. Make your changes with clear commit messages
3. Open a PR with context and screenshots (if applicable)

---

### 📄 License
This project is open-source. Use it as you like.

