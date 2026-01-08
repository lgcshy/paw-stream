# API Server Capability Specification

## ADDED Requirements

### Requirement: Project Structure
The API server SHALL follow Go standard project layout with clear separation of concerns across transport, domain, and storage layers.

#### Scenario: Directory structure compliance
- **WHEN** the project structure is examined
- **THEN** it contains `cmd/`, `internal/`, `deployments/`, `scripts/`, and `migrations/` directories
- **AND** the structure matches `docs/backend_project_layout.md` specification
- **AND** all Go code follows `internal/` package visibility rules

#### Scenario: Module initialization
- **WHEN** running `go mod tidy`
- **THEN** all dependencies resolve correctly
- **AND** no unused dependencies remain
- **AND** `go.mod` specifies minimum Go version 1.21+

---

### Requirement: Application Bootstrap
The API server SHALL provide a main entry point that initializes the Fiber application with proper dependency injection and graceful shutdown.

#### Scenario: Server startup
- **WHEN** running `go run cmd/api/main.go`
- **THEN** the Fiber server starts successfully on the configured port (default 3000)
- **AND** all middleware are registered in correct order
- **AND** all route handlers are mounted
- **AND** database connections are established (if configured)
- **AND** structured logs indicate successful startup

#### Scenario: Graceful shutdown
- **WHEN** the server receives SIGINT or SIGTERM signal
- **THEN** it stops accepting new requests
- **AND** waits for in-flight requests to complete (with timeout)
- **AND** closes database connections cleanly
- **AND** logs shutdown completion before exiting

---

### Requirement: Configuration Management
The API server SHALL load configuration from YAML files with support for environment variable overrides.

#### Scenario: Config file loading
- **WHEN** `config.yaml` exists in the working directory
- **THEN** the server reads all configuration values from the file
- **AND** default values are used for missing keys
- **AND** invalid YAML format causes startup failure with clear error message

#### Scenario: Environment variable override
- **WHEN** environment variables with prefix `PAWSTREAM_` are set
- **THEN** they override corresponding YAML config values
- **AND** nested config keys use underscore separator (e.g., `PAWSTREAM_DB_HOST`)

#### Scenario: Config validation
- **WHEN** critical configuration is missing or invalid
- **THEN** the server fails to start with descriptive error message
- **AND** logs indicate which config value is problematic

---

### Requirement: HTTP Middleware
The API server SHALL provide standard HTTP middleware for logging, request tracking, CORS, and authentication.

#### Scenario: Request ID injection
- **WHEN** any HTTP request is received
- **THEN** a unique request ID is generated
- **AND** the request ID is included in all log entries for that request
- **AND** the request ID is returned in `X-Request-ID` response header

#### Scenario: Structured logging
- **WHEN** the server handles a request
- **THEN** request details are logged in JSON format (method, path, status, duration)
- **AND** logs include request ID, user ID (if authenticated), and timestamp
- **AND** error logs include stack traces for debugging
- **AND** logs are written to both console and file (if file output is configured)

#### Scenario: CORS handling
- **WHEN** a browser makes a cross-origin request
- **THEN** appropriate CORS headers are set based on configuration
- **AND** preflight OPTIONS requests return 204 No Content
- **AND** configured origins are allowed in production

#### Scenario: Panic recovery
- **WHEN** a handler panics during request processing
- **THEN** the panic is caught by recovery middleware
- **AND** a 500 Internal Server Error response is returned
- **AND** the panic stack trace is logged
- **AND** the server continues running

---

### Requirement: Log Management
The API server SHALL support configurable log output with automatic log rotation to prevent disk space exhaustion.

#### Scenario: Log file configuration
- **WHEN** the server is configured with log file output
- **THEN** logs are written to the configured file path (e.g., `logs/api.log`)
- **AND** the log directory is created automatically if it doesn't exist
- **AND** console output can be enabled/disabled independently

#### Scenario: Log rotation by size
- **WHEN** the log file reaches the configured maximum size (e.g., 100MB)
- **THEN** the current log file is rotated (renamed with timestamp)
- **AND** a new log file is created
- **AND** old log files are preserved up to the configured max_backups limit

#### Scenario: Log rotation by age
- **WHEN** log files exceed the configured max_age (e.g., 30 days)
- **THEN** old log files are automatically deleted
- **AND** only files within the age limit are retained

#### Scenario: Log compression
- **WHEN** log rotation occurs and compression is enabled
- **THEN** rotated log files are compressed using gzip
- **AND** compressed files have `.gz` extension
- **AND** disk space is conserved

#### Scenario: Log levels
- **WHEN** different log levels are used (debug, info, warn, error)
- **THEN** logs are filtered based on configured minimum level
- **AND** development mode defaults to debug level
- **AND** production mode defaults to info level

---

### Requirement: Health Check Endpoint
The API server SHALL expose a health check endpoint for monitoring and load balancer probes.

#### Scenario: Basic health check
- **WHEN** `GET /health` is requested
- **THEN** the server returns 200 OK with JSON response
- **AND** the response includes `{"status": "ok", "timestamp": "..."}` 
- **AND** the check does not require authentication

#### Scenario: Health check with dependencies
- **WHEN** `GET /health?detailed=true` is requested
- **THEN** the response includes database connection status
- **AND** the response indicates MediaMTX reachability (if configured)
- **AND** overall status is "ok" only if all dependencies are healthy

---

### Requirement: Domain Layer - User Management
The API server SHALL implement user domain with repository pattern for data access abstraction.

#### Scenario: User model definition
- **WHEN** user domain types are examined
- **THEN** User struct includes id, username, nickname, password_hash, disabled, created_at, updated_at
- **AND** password is never stored in plain text
- **AND** UUIDs are used for user IDs

#### Scenario: User repository interface
- **WHEN** accessing user data
- **THEN** operations go through UserRepository interface
- **AND** interface defines Create, GetByID, GetByUsername, Update, Delete methods
- **AND** SQLite repository implements all interface methods

#### Scenario: User service operations
- **WHEN** UserService.Register is called with username and password
- **THEN** password is hashed using bcrypt
- **AND** a new user is created in the repository
- **AND** duplicate username returns appropriate error

---

### Requirement: Domain Layer - Device Management
The API server SHALL implement device domain with secret management for streaming authentication.

#### Scenario: Device model definition
- **WHEN** device domain types are examined
- **THEN** Device struct includes id, owner_user_id, name, location, publish_path, secret_hash, secret_cipher, secret_version, disabled, created_at, updated_at
- **AND** device secrets are stored encrypted
- **AND** publish_path is unique across all devices

#### Scenario: Device creation
- **WHEN** DeviceService.Create is called for a user
- **THEN** a unique device_id (UUID) is generated
- **AND** a cryptographically secure device_secret is generated
- **AND** publish_path is set to `dogcam/<device_id>` format
- **AND** secret is hashed for authentication and encrypted for retrieval

#### Scenario: Device secret rotation
- **WHEN** DeviceService.RotateSecret is called for a device
- **THEN** a new device_secret is generated
- **AND** secret_version is incremented
- **AND** old secret is invalidated
- **AND** device owner is notified of the change

---

### Requirement: Domain Layer - Access Control
The API server SHALL implement ACL logic to determine who can publish streams and who can view them.

#### Scenario: Publish authorization
- **WHEN** ACLService.CanPublish is called with device_secret and path
- **THEN** it verifies the device secret matches the path
- **AND** it checks the device is not disabled
- **AND** it returns true only if both conditions are met

#### Scenario: Read authorization
- **WHEN** ACLService.CanRead is called with user_token and path
- **THEN** it verifies the JWT token is valid
- **AND** it checks the user owns a device with that publish_path
- **AND** it checks the user is not disabled
- **AND** it returns true only if all conditions are met

---

### Requirement: Storage Layer - SQLite Implementation
The API server SHALL use SQLite as the database with CGO-free driver for simplified deployment and static compilation.

#### Scenario: Database initialization
- **WHEN** the server starts for the first time
- **THEN** SQLite database file is created at the configured path (default: `data/pawstream.db`)
- **AND** the data directory is created automatically if it doesn't exist
- **AND** connection to the database is established using modernc.org/sqlite driver
- **AND** migrations are applied automatically

#### Scenario: Database connection
- **WHEN** SQLite store is configured
- **THEN** a sql.DB connection pool is created
- **AND** connection parameters are set (max open/idle connections, connection lifetime)
- **AND** WAL mode is enabled for better concurrent access
- **AND** foreign keys are enabled

#### Scenario: Database migrations
- **WHEN** the server starts
- **THEN** migration files in `migrations/` are applied automatically if needed
- **AND** schema version is tracked in a migrations table
- **AND** failed migrations prevent server startup with clear error message
- **AND** migration status is logged

#### Scenario: CGO-free compilation
- **WHEN** the binary is compiled with `CGO_ENABLED=0`
- **THEN** the build succeeds without CGO dependencies
- **AND** the resulting binary is statically linked
- **AND** SQLite functionality works correctly

---

### Requirement: JWT Authentication
The API server SHALL provide JWT token generation and validation for business user authentication.

#### Scenario: Token generation
- **WHEN** a user successfully logs in
- **THEN** a JWT token is generated with user_id and username claims
- **AND** the token is signed with a secret key from configuration
- **AND** the token has a configurable expiration time (default 24 hours)

#### Scenario: Token validation
- **WHEN** a request includes JWT token in Authorization header
- **THEN** the auth_user middleware validates the token signature
- **AND** expired tokens are rejected with 401 Unauthorized
- **AND** valid tokens populate request context with user information

---

### Requirement: MediaMTX Integration
The API server SHALL implement MediaMTX authentication callback endpoint to authorize streaming operations.

#### Scenario: Publish callback handling
- **WHEN** MediaMTX calls `POST /mediamtx/auth` with action=publish
- **THEN** the server extracts device_secret from query params or basic auth
- **AND** the server validates the secret against the requested path
- **AND** returns 200 OK if authorized, 403 Forbidden if not
- **AND** logs the authorization decision with device_id and path

#### Scenario: Read callback handling
- **WHEN** MediaMTX calls `POST /mediamtx/auth` with action=read
- **THEN** the server extracts user_token from query params or header
- **AND** the server validates the user has permission to view the path
- **AND** returns 200 OK if authorized, 403 Forbidden if not

#### Scenario: Callback request format
- **WHEN** MediaMTX auth callback is received
- **THEN** the request includes JSON body with fields: action, path, protocol, ip, user, password, token
- **AND** the server handles both query and playback actions
- **AND** unknown actions are logged and denied by default

---

### Requirement: HTTP API Handlers (Stubs)
The API server SHALL provide placeholder HTTP handlers for business API endpoints to be implemented in Phase 3.

#### Scenario: Auth endpoints stub
- **WHEN** `POST /api/register` or `POST /api/login` is requested
- **THEN** the handler returns 501 Not Implemented
- **AND** logs indicate the endpoint is planned for Phase 3 implementation

#### Scenario: Device endpoints stub
- **WHEN** `GET /api/devices` or `POST /api/devices` is requested
- **THEN** the handler returns 501 Not Implemented
- **AND** response includes a message indicating future implementation

#### Scenario: Paths endpoint stub
- **WHEN** `GET /api/paths` is requested
- **THEN** the handler returns 501 Not Implemented
- **AND** the endpoint structure is defined for future implementation

---

### Requirement: Error Handling
The API server SHALL implement consistent error handling with structured error responses.

#### Scenario: Error response format
- **WHEN** an error occurs during request processing
- **THEN** the response includes JSON with fields: error, message, request_id
- **AND** appropriate HTTP status code is returned (4xx for client errors, 5xx for server errors)
- **AND** sensitive details are not exposed in production

#### Scenario: Domain error mapping
- **WHEN** domain layer returns a business error (e.g., UserNotFound)
- **THEN** the transport layer maps it to appropriate HTTP status code
- **AND** user-friendly error message is returned
- **AND** internal error details are logged but not exposed to client

---

### Requirement: Development Experience
The API server SHALL provide scripts and tooling for efficient local development.

#### Scenario: Development server
- **WHEN** `scripts/dev_run.sh` is executed
- **THEN** the server starts with hot reload enabled
- **AND** code changes trigger automatic restart
- **AND** logs are formatted for human readability in development mode

#### Scenario: Build script
- **WHEN** `scripts/build.sh` is executed
- **THEN** a production binary is compiled to `bin/api`
- **AND** the binary is statically linked (if possible)
- **AND** version information is embedded via ldflags

---

### Requirement: Documentation
The API server SHALL include comprehensive documentation for setup, development, and deployment.

#### Scenario: README completeness
- **WHEN** a new developer reads `server/api/README.md`
- **THEN** they understand the project architecture and structure
- **AND** they know how to install dependencies and run the server
- **AND** they understand the configuration options
- **AND** they know where to find OpenSpec project conventions
- **AND** they can follow API endpoint documentation
