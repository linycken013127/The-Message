# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

"The Message" is an online implementation of a board game with a full-stack architecture:
- **Backend**: Go with Gin web framework, GORM ORM, MySQL database
- **Frontend**: React + TypeScript + Vite with Recoil state management

### Architecture Pattern

The backend uses a clean architecture with the following layers:
- **HTTP Handlers** (`service/delivery/http/v1/`): Route definitions and request/response handling
- **Services** (`service/service/`): Business logic for GameService, PlayerService, CardService, DeckService
- **Repositories** (`service/repository/mysql/`): Data access layer using GORM
- **Models & Enums** (`enums/`): Game entities and constants
- **Database** (`database/`): Migrations and seeders for MySQL

The frontend uses Recoil atoms for global state management, organized into:
- `states/globalCards/`: Shared card state (functionalCardStates, lineUpCardStates)
- `states/personalCards/`: Player-specific state
- `components/`: UI components
- `pages/`: Game pages (entrance, game board, map, playerSection, table)
- `hooks/`: Custom React hooks for card-related logic

## Backend Development

### Prerequisites
- Go 1.22+
- Docker for MySQL
- goblin/migrate for database migrations

### Setup Database

```bash
cd Backend
docker-compose up -d           # Start MySQL container
go run ./cmd/migrate/migrate.go # Run migrations
go run ./cmd/migrate/game_card_seeder.go # Seed game cards
```

### Common Backend Commands

**Database Operations:**
```bash
go run ./cmd/migrate/migrate.go   # Run migrations up
go run ./cmd/migrate/rollback.go  # Roll back one migration
go run ./cmd/migrate/refresh.go   # Clear DB and re-run migrations
go run ./cmd/migrate/game_card_seeder.go # Seed card data
```

**Running the Server:**
```bash
cd Backend
go run ./cmd/app/main.go  # Start server on port 8080
```

**Generate Swagger API Documentation:**
```bash
cd Backend
swag init -g ./cmd/app/main.go -output ./cmd/app/docs
# View at http://127.0.0.1:8080/swagger/index.html
```

**Code Quality:**
```bash
cd Backend
go mod tidy                    # Update go.mod and go.sum
goimports -l -w .            # Format imports
golangci-lint run ./...       # Run linter
```

**Running Tests:**
```bash
cd Backend
go test ./tests/e2e/...       # Run all e2e tests
go test -v ./tests/e2e/game_api_test.go # Run specific test file
```

Note: Tests use testify/suite for integration testing against a test database defined in `config/test_database.go`. Tests are set up via `suite_test.go` with automatic database seeding in `SetupTest()`.

## Frontend Development

### Setup

```bash
cd Frontend/the-message
npm install
```

### Common Frontend Commands

**Development:**
```bash
npm run dev      # Start Vite dev server with hot reload
```

**Build & Preview:**
```bash
npm run build    # Build for production (runs TypeScript check + Vite)
npm run preview  # Preview production build locally
```

**Code Quality:**
```bash
npm run lint     # Run ESLint with auto-fix, fail on warnings
npm run test     # Run Jest tests
```

**Running Single Tests:**
```bash
npm run test -- App.test.tsx  # Run specific test file
```

## API Integration

The frontend communicates with the backend at `localhost:8080/api/v1/`. Key endpoints:

- `POST /api/v1/games` - Create new game (requires player list in JSON body)
- `GET /api/v1/cards` - Fetch card data
- `GET /api/v1/players/{id}` - Player info
- `POST /api/v1/players/{id}/cards` - Player card actions

Swagger documentation automatically available at: `http://127.0.0.1:8080/swagger/index.html`

## Database Migrations

- Migration files location: `Backend/database/migrations/`
- Use golang-migrate for consistency
- Create new migrations with: `migrate create -ext sql -dir Backend/database/migrations -seq <migration_name>`

## Key Dependencies & Tools

**Backend:**
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM
- `github.com/golang-migrate/migrate/v4` - Database migrations
- `github.com/swaggo/gin-swagger` - Swagger documentation
- `github.com/stretchr/testify/suite` - Testing framework

**Frontend:**
- `react` & `react-dom` - UI framework
- `recoil` - State management
- `chakra-ui` - Component library
- `vite` - Build tool
- `jest` & `@testing-library/react` - Testing

## Development Workflow Notes

- Use `rg` for fast text search instead of grep (respects gitignore)
- Database state is isolated per test via transaction rollback in `TearDownTest()`
- Swagger docs are auto-generated from comments in handler files
- TypeScript strict mode is enforced in frontend builds
- ESLint warnings cause build failure (`--max-warnings 0`)
