## Notes App (Go + Gin + GORM + JWT)

A clean and minimal REST API for a Notes application built with Go, Gin, GORM, PostgreSQL, and JWT-based authentication. It supports user signup/login with secure password hashing and cookie-based auth, along with CRUD for notes.

---

### ✨ Features
- **JWT auth with HttpOnly cookies**: Signup, login, and a protected validation route
- **Notes CRUD**: Create, read, update, and delete notes
- **PostgreSQL with GORM**: Auto-migration for models
- **12-factor ready**: Configuration via environment variables
- **Production-ready routing**: Powered by `gin-gonic`

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
  │   ├─ autho.go           # Signup, Login, Validate
  │   └─ controllers.go     # User create, Notes CRUD
  ├─ intializers/
  │   ├─ datbase.go         # ConnectToDB (Postgres via GORM)
  │   ├─ LoadEnvVar.go      # Load .env via godotenv
  │   └─ SyncDatabase.go    # AutoMigrate(User, Note)
  ├─ middleware/
  │   └─ Reqautho.go        # RequireAuth middleware (JWT cookie)
  ├─ models/
  │   ├─ notes.go           # Note model
  │   └─ user.go            # User model
  ├─ main.go                # Router, routes, bootstrapping
  ├─ go.mod
  └─ go.sum
```

---

### 🗄️ Database Models
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
   UserID  uint  // FK to User
 }
```
- Relations: One `User` has many `Note`.
- Auto-migrations are executed during startup via `SyncDataBase()`.

---

### 🔐 Authentication Flow
- `POST /signup`: Creates a user with a bcrypt-hashed password.
- `POST /login`: Verifies credentials, then returns a signed JWT stored in an `Authorization` cookie (HttpOnly, SameSite=Lax).
- `GET /validate`: Protected route using `RequireAuth` middleware; returns a simple success payload if token is valid and not expired.

JWT details:
- Algorithm: HS256
- Claims: `sub` (user ID), `exp` (30 days)
- Secret: `SECRET` environment variable

---

### 🌐 REST API

Base URL: `http://localhost:<PORT>` (defaults to `8080`)

#### Auth
- `POST /signup`
  - Body (form or JSON): `{ "Email": string, "Password": string }`
  - Response: `201 { user }`

- `POST /login`
  - Body (form or JSON): `{ "Email": string, "Password": string }`
  - Sets cookie: `Authorization=<JWT>`
  - Response: `200 {}`

- `GET /validate` (protected)
  - Header/Cookie: `Authorization` cookie must be present
  - Response: `200 { "message": "i am logged in" }`

#### Users
- `POST /user`
  - Creates a user record (name/email/notes). Primarily for testing basic user creation (separate from auth signup).

#### Notes
- `POST /notes`
  - Body (JSON): `{ "title": string, "content": string, "userID": number }`
  - Response: `200 { note }`

- `GET /notes`
  - Response: `200 { notes: Note[] }`

- `GET /notes/:id`
  - Response: `200 { notes: Note[] }` (fetches by id)

- `PUT /notes/:id`
  - Body (JSON): `{ "title": string, "content": string }`
  - Response: `200 { note }`

- `DELETE /notes/:id`
  - Response: `200 { note }`

Note: Current CRUD endpoints are public in `main.go`. You can wrap them with `RequireAuth` if you want to make them private per user.

---

### 🧪 Quick Start

#### 1) Prerequisites
- Go 1.21+
- PostgreSQL running and reachable

#### 2) Environment
Create a `.env` file in the project root:
```env
PORT=8080
DB_url=postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
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
You should see:
- "Initializing database"
- "welcome to notes app"
- and a successful DB connection log

GORM will auto-migrate `User` and `Note` tables on startup.

---

### 🧰 cURL Examples
```bash
# Signup
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"Email":"alice@example.com","Password":"SuperSecret123"}'

# Login (stores cookie)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"Email":"alice@example.com","Password":"SuperSecret123"}' -i

# Validate (send cookie from login response)
curl http://localhost:8080/validate \
  --cookie "Authorization=REPLACE_WITH_JWT"

# Create a note
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"My Note","content":"Hello","userID":1}'

# List notes
curl http://localhost:8080/notes

# Get one note by id
curl http://localhost:8080/notes/1

# Update note
curl -X PUT http://localhost:8080/notes/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","content":"Updated body"}'

# Delete note
curl -X DELETE http://localhost:8080/notes/1
```

---

### ⚙️ Configuration
- **PORT**: Server port (default `8080`)
- **DB_url**: Postgres DSN consumed by GORM
- **SECRET**: JWT signing secret

---

### 🔒 Security Notes
- Auth token is stored in an HttpOnly cookie to mitigate XSS token theft.
- In production, set the cookie `Secure` flag and serve over HTTPS.
- Use a strong, rotated `SECRET` value.

---

### 🚀 Roadmap Ideas
- Scope notes by authenticated user (protect CRUD with `RequireAuth`)
- Add pagination and search to `GET /notes`
- Add refresh tokens / logout endpoint
- Add OpenAPI/Swagger docs

---

### 🙌 Contributing
1. Fork the repo and create a feature branch
2. Make your changes with clear commit messages
3. Open a PR with context and screenshots (if applicable)

---

### 📄 License
This project is open-source. Use it as you like.
