## Notes App (Go + Gin + GORM + JWT)

A clean and minimal REST API for a Notes application built with Go, Gin, GORM, PostgreSQL, and JWT-based authentication. It supports user signup/login with secure password hashing and cookie-based auth, along with CRUD for notes.

---

### ✨ Features
- **JWT auth (access + refresh)**: Access token in JSON, refresh token in HttpOnly cookie
- **Refresh flow & logout**: `/refreshToken` issues new access tokens; `/logout` clears refresh cookie
- **Protected route**: `/validate` guarded by middleware
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
  │   ├─ autho.go                 # Signup, Login, Logout, Validate
  │   ├─ controllers.go           # User create, Notes CRUD
  │   └─ refreshControllers.go    # Refresh access token
  ├─ intializers/
  │   ├─ datbase.go               # ConnectToDB (Postgres via GORM)
  │   ├─ LoadEnvVar.go            # Load .env via godotenv
  │   └─ SyncDatabase.go          # AutoMigrate(User, Note, RefreshToken)
  ├─ middleware/
  │   └─ Reqautho.go              # RequireAuth middleware (reads Authorization cookie)
  ├─ models/
  │   ├─ notes.go                 # Note model
  │   ├─ refresh.go               # RefreshToken model
  │   └─ user.go                  # User model
  ├─ utils/
  │   └─ jwt.go                   # JWT helpers (access & refresh)
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

// models/refresh.go
 type RefreshToken struct {
   ID        uint `gorm:"primaryKey"`
   Token     string
   UserID    uint
   ExpiresAt time.Time
 }
```
- Relations: One `User` has many `Note`.
- Auto-migrations are executed during startup via `SyncDataBase()`.

---

### 🔐 Authentication Flow
- `POST /signup`: Creates a user with a bcrypt-hashed password.
- `POST /login`: Verifies credentials, then returns a short-lived access token in JSON and sets a long-lived refresh token in an HttpOnly cookie named `refresh_token`.
- `POST /refreshToken`: Reads `refresh_token` cookie and issues a new access token in JSON.
- `GET /logout`: Clears the `refresh_token` cookie and deletes it from storage.
- `GET /validate`: Protected route using `RequireAuth` middleware; it expects an `Authorization` cookie containing a valid access token.

JWT details:
- Algorithm: HS256
- Access token claims: `user_id`, `exp`, `type: "access"` (default lifespan ~15m)
- Refresh token claims: `user_id`, `exp`, `type: "refresh"` (default lifespan ~7d)
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
  - Sets cookie: `refresh_token=<JWT>` (HttpOnly)
  - Response: `200 { "access token": string }`

- `POST /refreshToken`
  - Reads cookie: `refresh_token`
  - Response: `200 { "access_token": string }`

- `GET /logout`
  - Clears cookie: `refresh_token`
  - Response: `200 { "message": string }`

- `GET /validate` (protected)
  - Cookie: `Authorization=<ACCESS_TOKEN>` must be present (see cURL below)
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
It will also create `RefreshToken` if referenced by the migrator.

---

### 🧰 cURL Examples
```bash
# Signup
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"Email":"alice@example.com","Password":"SuperSecret123"}'

# Login: gets access token JSON and sets refresh_token cookie
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"Email":"alice@example.com","Password":"SuperSecret123"}' -i

# Suppose the login response JSON was: {"access token":"ACCESS_TOKEN_VALUE"}
# Set the access token as Authorization cookie for protected routes

# Validate (send Authorization cookie)
curl http://localhost:8080/validate \
  --cookie "Authorization=ACCESS_TOKEN_VALUE"

# Refresh access token (uses refresh_token cookie set by login)
curl -X POST http://localhost:8080/refreshToken

# Logout (clears refresh_token cookie)
curl http://localhost:8080/logout

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

