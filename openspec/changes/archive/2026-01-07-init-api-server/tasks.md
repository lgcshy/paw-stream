# Implementation Tasks: Initialize API Server Project

## 1. Project Initialization
- [x] 1.1 Initialize Go module in `server/api/`
- [x] 1.2 Create directory structure according to `docs/backend_project_layout.md`
- [x] 1.3 Add `.gitignore` for Go projects (include `*.db` for SQLite files)
- [x] 1.4 Create `go.mod` and install core dependencies (Fiber, modernc.org/sqlite, viper, zerolog)
- [x] 1.5 Add `README.md` with project overview and setup instructions

## 2. Configuration Management
- [x] 2.1 Create `internal/config/config.go` with YAML config struct (include log config)
- [x] 2.2 Create `internal/config/defaults.go` with default values
- [x] 2.3 Add `config.yaml.example` file (with log rotation settings)
- [x] 2.4 Implement config loading with Viper (support env var override)
- [x] 2.5 Add log configuration struct (path, max_size, max_backups, max_age, compress)

## 3. Application Bootstrap
- [x] 3.1 Create `cmd/api/main.go` with application entry point
- [x] 3.2 Initialize logger with lumberjack for log rotation
- [x] 3.3 Create `internal/app/api/app.go` for Fiber app initialization
- [x] 3.4 Create `internal/app/api/routes.go` for route registration
- [x] 3.5 Implement graceful shutdown handling (close log files properly)

## 4. Middleware Layer
- [x] 4.1 Create `internal/transport/http/middleware/request_id.go`
- [x] 4.2 Create `internal/transport/http/middleware/logger.go` (zerolog integration with file output)
- [x] 4.3 Create `internal/transport/http/middleware/cors.go`
- [x] 4.4 Create `internal/transport/http/middleware/auth_user.go` (JWT validation stub)
- [x] 4.5 Create `internal/transport/http/middleware/recovery.go` (panic recovery with stack trace logging)

## 5. Basic HTTP Handlers
- [x] 5.1 Create `internal/transport/http/handlers/health_handler.go` (GET /health)
- [x] 5.2 Create `internal/transport/http/handlers/auth_handler.go` (stub for /api/register, /api/login)
- [x] 5.3 Create `internal/transport/http/handlers/device_handler.go` (stub for /api/devices)
- [x] 5.4 Create `internal/transport/http/handlers/path_handler.go` (stub for /api/paths)
- [x] 5.5 Create `internal/transport/http/handlers/mediamtx_auth_handler.go` (stub for /mediamtx/auth)

## 6. Domain Layer - User
- [x] 6.1 Create `internal/domain/user/model.go` (User struct)
- [x] 6.2 Create `internal/domain/user/repo.go` (Repository interface)
- [x] 6.3 Create `internal/domain/user/service.go` (UserService with Register, Login, GetByID)

## 7. Domain Layer - Device
- [x] 7.1 Create `internal/domain/device/model.go` (Device struct)
- [x] 7.2 Create `internal/domain/device/repo.go` (Repository interface)
- [x] 7.3 Create `internal/domain/device/service.go` (DeviceService with Create, GetByID, ListByOwner)

## 8. Domain Layer - ACL
- [x] 8.1 Create `internal/domain/acl/policy.go` (CanPublish, CanRead functions)
- [x] 8.2 Create `internal/domain/acl/service.go` (ACLService)

## 9. Storage Layer - SQLite Implementation
- [x] 9.1 Create `internal/store/sqlite/db.go` (SQLite connection with modernc.org/sqlite)
- [x] 9.2 Create `internal/store/sqlite/user_repo.go` (SQLite user repository)
- [x] 9.3 Create `internal/store/sqlite/device_repo.go` (SQLite device repository)
- [x] 9.4 Implement connection pool with sql.DB
- [x] 9.5 Add automatic database file creation
- [x] 9.6 Create `migrations/001_init_schema.up.sql` (users and devices tables for SQLite)
- [x] 9.7 Create `migrations/001_init_schema.down.sql` (rollback script)
- [x] 9.8 Implement auto-migration on startup

## 11. Utility Packages
- [x] 11.1 Create `internal/pkg/jwtutil/jwt.go` (JWT sign and verify functions)
- [x] 11.2 Create `internal/pkg/idgen/id.go` (UUID and secret generation)
- [x] 11.3 Create `internal/pkg/errors/errors.go` (error types and wrapping)
- [x] 11.4 Add password hashing utilities (bcrypt)
- [x] 11.5 Create `internal/pkg/logger/logger.go` (logger initialization with lumberjack)

## 12. MediaMTX Integration
- [x] 12.1 Create `internal/integration/mediamtx/types.go` (AuthRequest, AuthResponse structs)
- [x] 12.2 Create `internal/integration/mediamtx/authz.go` (authorization logic for MediaMTX callbacks)
- [x] 12.3 Define action types: publish, read, playback

## 13. Deployment Configuration
- [x] 13.1 Create `deployments/docker-compose.yaml` (API server + MediaMTX, no database service needed)
- [x] 13.2 Create `deployments/Dockerfile` for API server (multi-stage build, CGO_ENABLED=0)
- [x] 13.3 Update `deployments/mediamtx.yml` with auth callback configuration
- [x] 13.4 Add volume mount for SQLite database file in docker-compose

## 14. Scripts
- [x] 14.1 Create `scripts/dev_run.sh` (run with hot reload)
- [x] 14.2 Create `scripts/migrate.sh` (run database migrations)
- [x] 14.3 Create `scripts/build.sh` (compile binary)
- [x] 14.4 Make scripts executable

## 15. Testing & Validation
- [x] 15.1 Verify `go mod tidy` succeeds
- [x] 15.2 Verify `CGO_ENABLED=0 go build ./cmd/api` compiles successfully (static binary)
- [x] 15.3 Verify `go run cmd/api/main.go` starts server and creates SQLite database
- [x] 15.4 Test `GET /health` returns 200 OK
- [x] 15.5 Test graceful shutdown (Ctrl+C) and database connection cleanup
- [x] 15.6 Verify structured logging works (JSON output to file)
- [x] 15.7 Test log file is created at configured path
- [x] 15.8 Test log rotation (create large log and verify rotation occurs)
- [x] 15.9 Run `go fmt ./...` and `go vet ./...`
- [x] 15.10 Validate project structure matches `docs/backend_project_layout.md`
- [x] 15.11 Verify SQLite database file is created at configured path
- [x] 15.12 Test migrations are applied automatically on first run

## 16. Documentation
- [x] 16.1 Write `server/api/README.md` with:
  - Project structure explanation
  - Setup instructions (Go version, dependencies)
  - Development workflow (run, test, build)
  - Configuration guide (SQLite database path, log settings)
  - API endpoints overview
  - SQLite advantages and limitations
  - Log file management and rotation
- [x] 16.2 Add code comments for exported functions and types
- [x] 16.3 Create `api/openapi.yaml` skeleton (optional for Phase 3)
- [x] 16.4 Document MediaMTX callback integration points
- [x] 16.5 Document SQLite database backup and migration procedures
- [x] 16.6 Document log file rotation and cleanup procedures
