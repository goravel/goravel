# Goravel Framework - Go Web Application

Goravel is a Laravel-style web application framework for Go with complete functions and good scalability. This repository provides a starting scaffolding for Gopher developers to quickly build applications.

**ALWAYS reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Bootstrap, Build, and Test
- Copy environment configuration: `cp .env.example .env`
- Generate application key: `./artisan key:generate`
- Download dependencies: `go mod download` -- takes 30-60 seconds
- Build the application: `go build -o main .` -- takes 60-90 seconds. NEVER CANCEL. Set timeout to 180+ seconds.
- Run tests: `go test ./...` -- takes 10-15 seconds. NEVER CANCEL. Set timeout to 60+ seconds.

### Run the Application
- **ALWAYS run the bootstrapping steps first** (copy .env, generate key, download deps)
- Production mode: `go run .` -- starts HTTP server on http://127.0.0.1:3000
- Development with hot reload: `/home/runner/go/bin/air` -- requires `go install github.com/air-verse/air@latest`
- Direct binary execution: `./main` (after building)

### Artisan CLI Commands
The `./artisan` script provides extensive CLI functionality:
- Generate app key: `./artisan key:generate`
- List all commands: `./artisan`
- Show routes: `./artisan route:list`
- Create controller: `./artisan make:controller ControllerName`
- Create model: `./artisan make:model ModelName`
- Create migration: `./artisan make:migration create_table_name`
- Run migrations: `./artisan migrate`
- Database seeding: `./artisan db:seed`

## Validation

### Manual Testing
- **ALWAYS manually validate any new code** by running the complete application startup sequence
- Test the home page: `curl http://127.0.0.1:3000/` (should return HTML welcome page)
- Test API endpoint: `curl http://127.0.0.1:3000/users/123` (should return `{"Hello":"Goravel"}`)
- **ALWAYS run through at least one complete end-to-end scenario** after making changes

### Code Quality
- Format code: `go fmt ./...` -- instant, shows files that were formatted
- Vet code: `go vet ./...` -- instant, may show warnings about signal handling in main.go (expected, not a bug)
- Install golangci-lint: `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2`
- Lint code: `/home/runner/go/bin/golangci-lint run --timeout 3m` -- takes 2-3 minutes. NEVER CANCEL.
- **Note**: golangci-lint may show import errors due to framework version compatibility issues, these do not prevent functionality.

### Critical Timing Information
- **NEVER CANCEL BUILDS OR LONG-RUNNING COMMANDS**
- Initial build: 60-90 seconds. Subsequent builds: 2-5 seconds (cached). Use timeout of 180+ seconds.
- Test commands: 2-10 seconds. Use timeout of 60+ seconds.
- Linting: 2-3 minutes. Use timeout of 300+ seconds.
- Dependency download: First time 30-60 seconds, subsequent <5 seconds. Use timeout of 120+ seconds.
- Artisan commands: <1 second each. Use timeout of 60+ seconds.

## Project Structure

### Key Directories
- `app/` - Application source code
  - `http/controllers/` - HTTP controllers (e.g., user_controller.go)
  - `http/middleware/` - HTTP middleware
  - `models/` - Data models
  - `providers/` - Service providers
  - `console/` - Console commands
  - `grpc/` - GRPC services
- `config/` - Configuration files (app.go, http.go, session.go, etc.)
- `routes/` - Route definitions (web.go, api.go, grpc.go)
- `database/` - Database migrations and seeders
- `tests/feature/` - Feature tests using testify/suite
- `resources/views/` - Template files (.tmpl format)
- `storage/` - File storage and logs
- `bootstrap/` - Application bootstrapping

### Important Files
- `main.go` - Application entry point, starts HTTP server
- `go.mod` - Go module definition, uses Go 1.23.0+
- `.env.example` - Environment configuration template
- `artisan` - CLI tool bash script wrapper
- `.air.toml` - Hot reload configuration
- `docker-compose.yml` - Docker setup (may have issues in some environments)

## Common Outputs

### Repository Root Contents
```
.air.toml          .env.example       .git               .github
.gitignore         Dockerfile         LICENSE            README.md
README_zh.md       app                artisan            bootstrap
config             database           docker-compose.yml go.mod
go.sum             main.go            public             resources
routes             storage            tests
```

### Available Artisan Commands (Partial List)
```
Available commands:
  about              Display basic information about your application
  build              Build the application
  migrate            Run the database migrations
 cache:
  cache:clear        Flush the application cache
 make:
  make:controller    Create a new controller class
  make:model         Create a new model class
  make:migration     Create a new migration file
  make:middleware    Create a new middleware class
  make:test          Create a new test class
 route:
  route:list         List all registered routes
```

### Default Routes
```
GET|HEAD     / ........................ goravel/routes.Web.func1  
GET|HEAD     users/{id} ................. goravel/app/http/controllers.(*UserController).Show  
```

## Environment and Dependencies

### Go Version
- Required: Go 1.23.0+ (toolchain go1.24.0)
- Current system typically provides compatible version

### Key Dependencies
- github.com/goravel/framework v1.16.1 - Core framework
- github.com/goravel/gin v1.4.0 - HTTP driver
- github.com/gin-gonic/gin v1.10.1 - Web framework
- github.com/stretchr/testify v1.10.0 - Testing framework

### Development Tools
- Air for hot reload: `go install github.com/air-verse/air@latest`
- golangci-lint for comprehensive linting
- Standard Go tools: gofmt, go vet, go test

## Known Issues and Limitations

### Database Connectivity
- Application shows PostgreSQL connection warnings: "failed to connect to user=root database=goravel"
- This is expected when no database is configured and does NOT prevent application functionality
- Application runs successfully without database for basic HTTP operations

### Docker Build
- Docker build may fail in some environments due to certificate issues with Go module proxy
- Alternative: Build locally and copy binary to container

### Code Analysis
- golangci-lint may show import errors due to framework version compatibility
- go vet shows expected warning about unbuffered signal channel in main.go
- These do not prevent successful build and run

## Framework-Specific Guidelines

### Creating New Components
- Always use artisan commands to generate boilerplate: `./artisan make:controller`, `./artisan make:model`, etc.
- Follow Laravel-like structure and naming conventions
- Controllers in `app/http/controllers/`, models in `app/models/`

### Route Registration
- Web routes: Define in `routes/web.go`
- API routes: Define in `routes/api.go`
- GRPC routes: Define in `routes/grpc.go`

### Configuration
- Environment variables: Defined in `.env` file
- Application config: `config/app.go`
- HTTP config: `config/http.go`
- Follow config.Env() pattern for environment variable access

### Testing
- Feature tests: Place in `tests/feature/`
- Use testify/suite pattern as shown in example_test.go
- Tests extend tests.TestCase which handles framework bootstrapping