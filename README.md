# Peace - Mental Wellness Application

A full-stack mental health tracking application with Go backend and Next.js frontend.

## Project Structure

```
/peace
├── backend/           # Go backend application
│   ├── cmd/          # Application entrypoints
│   ├── internal/     # Private application code
│   ├── pkg/          # Public library code
│   ├── configs/      # Configuration files
│   ├── migrations/   # Database migrations
│   ├── Dockerfile    # Backend Docker config
│   └── Makefile      # Backend build commands
├── web/              # Next.js frontend application
├── docker-compose.yml # Production Docker setup
├── docker-compose.dev.yml # Development Docker setup
└── Makefile          # Root Makefile (delegates to backend)
```

## Quick Start

### Prerequisites
- Go 1.24+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL (for local development)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd peace
   ```

2. **Start backend development server**
   ```bash
   make dev
   ```

3. **Start frontend development server** (in another terminal)
   ```bash
   cd web
   npm install
   npm run dev
   ```

4. **Or start both together**
   ```bash
   make full-dev
   ```

### Docker Development

```bash
# Start all services (backend + database + redis)
make docker-up

# View logs
make logs

# Stop all services
make docker-down
```

## Available Commands

### Backend Commands
```bash
make dev              # Start development server with hot reload
make build            # Build the application
make test             # Run tests
make test-coverage    # Run tests with coverage
make health           # Check application health
```

### Frontend Commands
```bash
make web-dev          # Start frontend development
make web-build        # Build frontend for production
make web-install      # Install frontend dependencies
```

### Docker Commands
```bash
make docker-up        # Start all services
make docker-down      # Stop all services
make docker-build     # Build Docker images
make docker-clean     # Clean Docker resources
```

### Database Commands
```bash
make migrate-up       # Run database migrations
make migrate-down     # Rollback migrations
make migrate-status   # Show migration status
```

## API Endpoints

- **Health Check**: `GET /health`
- **Authentication**: `POST /api/auth/login`, `POST /api/auth/register`
- **Mental Health Records**: `GET|POST /api/mental-health-records`
- **Streak**: `GET /api/mental-health-records/streak`
- **Quotes**: `GET /api/quotes/random`

## Configuration

Backend configuration is managed through:
- `backend/configs/config.yml` - Main configuration file
- Environment variables (see docker-compose files)

Frontend configuration:
- `web/.env.local` - Local environment variables
- `web/next.config.ts` - Next.js configuration

## Development

### Backend Development
The backend follows Clean Architecture principles:
- `cmd/` - Application entrypoints
- `internal/domain/` - Business logic and entities
- `internal/application/` - Use cases and commands
- `internal/infrastructure/` - External dependencies
- `internal/interfaces/` - HTTP handlers and middleware

### Frontend Development
The frontend is built with:
- Next.js 14 with App Router
- TypeScript
- Tailwind CSS
- React Hook Form with Zod validation

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests: `make test`
4. Submit a pull request

## License

MIT License
