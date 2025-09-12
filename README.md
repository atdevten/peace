# Peace - Mental Wellness Application

A full-stack mental health tracking application with Go backend and Next.js frontend.

## Project Structure

```
/peace
├── backend/                    # Go backend application
│   ├── cmd/                   # Application entrypoints
│   ├── internal/              # Private application code
│   ├── pkg/                   # Public library code
│   ├── configs/               # Configuration files
│   ├── migrations/            # Database migrations
│   ├── Dockerfile             # Backend Docker config
│   └── Makefile               # Backend build commands
├── web/                       # Next.js frontend application
├── docker-compose.local.yml   # Local development
├── docker-compose.production.yml # Production environment
├── DOCKER.md                  # Docker setup guide
└── Makefile                   # Root Makefile (delegates to backend)
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

2. **Start local development:**
   ```bash
   # Start everything in containers (recommended)
   make docker-up
   
   # Access:
   # - Frontend: http://localhost (via Traefik)
   # - Backend API: http://api.localhost (via Traefik)
   # - WebSocket Server: http://ws.localhost (via Traefik)
   # - Traefik Dashboard: http://traefik.localhost:8080
   ```

   **Alternative: Run backend locally**
   ```bash
   # Start only database and redis containers
   docker-compose -f docker-compose.local.yml up -d postgres redis
   
   # Run backend locally (in another terminal)
   make dev
   
   # Run frontend locally (in another terminal)
   make web-dev
   ```

3. **Install frontend dependencies** (if running frontend locally)
   ```bash
   make web-install
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

### Backend Configuration
- `backend/configs/config.yml` - Main configuration file
- Environment variables (see docker-compose files)

### Frontend Configuration
- `web/.env.local` - Local environment variables
- `web/next.config.ts` - Next.js configuration

### Production Deployment
1. **Setup Environment Files:**
   ```bash
   cp env.staging.example .env.staging
   cp env.production.example .env.production
   # Edit the files with your actual values
   ```

2. **Deploy to Staging:**
   ```bash
   ./scripts/deploy.sh staging
   ```

3. **Deploy to Production:**
   ```bash
   ./scripts/deploy.sh production
   ```

### Traefik Configuration
- **Auto SSL**: Let's Encrypt certificates automatically generated
- **Dashboard**: Available at `http://your-domain:8080` (staging) or `https://traefik.your-domain` (production)
- **Routes**: Automatically configured based on Docker labels

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
