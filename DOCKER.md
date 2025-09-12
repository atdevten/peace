# Docker Compose Setup Guide

This project uses multiple Docker Compose files for different environments and use cases.

## 📁 **Docker Compose Files Overview**

| File | Purpose | Use Case |
|------|---------|----------|
| `docker-compose.local.yml` | **Local development** | Complete stack with hot reload |
| `docker-compose.production.yml` | **Production environment** | Full production setup with Traefik |

---

## 🚀 **Quick Start Commands**

### **Local Development**
```bash
# Start everything (backend + frontend + database + redis)
make docker-up

# Or manually:
docker-compose -f docker-compose.local.yml up -d

# View logs
docker-compose -f docker-compose.local.yml logs -f

# Stop everything
make docker-down
```

### **Production Environment**
```bash
# Start production environment
make prod-up

# Or manually:
docker-compose -f docker-compose.production.yml up -d
```

---

## 🏗️ **Development Workflow**

### **Full Container Development (Recommended)**
```bash
# Start everything in containers
make docker-up

# Access:
# - Frontend: http://localhost (via Traefik)
# - Backend API: http://api.localhost (via Traefik)
# - WebSocket Server: http://ws.localhost (via Traefik)
# - Traefik Dashboard: http://traefik.localhost:8080
# - Database: localhost:5432
# - Redis: localhost:6379
```

### **Alternative: Run Backend Locally**
```bash
# Start only database and redis containers
docker-compose -f docker-compose.local.yml up -d postgres redis

# Run backend locally (with hot reload)
cd backend
make dev

# Run frontend locally (with hot reload)
cd web
npm run dev
```

---

## 🌐 **Traefik Routing**

### **Local Development**
- **Frontend**: `http://localhost` → Frontend container
- **Backend API**: `http://api.localhost` → Backend container
- **WebSocket Server**: `http://ws.localhost` → WebSocket container
- **Traefik Dashboard**: `http://traefik.localhost:8080` → Traefik dashboard

### **Production**
- **Frontend**: `https://peace.com` → Frontend container
- **Backend API**: `https://api.peace.com` → Backend container
- **Traefik Dashboard**: `https://traefik.peace.com` → Traefik dashboard

---

## 🔧 **Environment Configuration**

### **Local Development**
- **Database**: `peace_local`
- **Redis**: No password
- **JWT Secret**: `local-jwt-secret-key`
- **Log Level**: `debug`
- **Hot Reload**: Enabled
- **Traefik**: Enabled (HTTP only, no SSL)

### **Production**
- **Database**: `peace_production`
- **Redis**: Password protected
- **JWT Secret**: From environment variables
- **Log Level**: `warn`
- **Traefik**: Enabled with SSL + logging

---

## 🐳 **Docker Images**

### **Local Development**
- **Backend**: Built from `./backend/Dockerfile.dev`
- **Frontend**: Built from `./web/Dockerfile`
- **Database**: `postgres:15-alpine`
- **Redis**: `redis:7-alpine`
- **Traefik**: `traefik:v3.0`

### **Production**
- **Backend**: `ghcr.io/atdevten/peace/backend:latest`
- **Frontend**: `ghcr.io/atdevten/peace/frontend:latest`
- **Database**: `postgres:15-alpine`
- **Redis**: `redis:7-alpine`
- **Traefik**: `traefik:v3.0`

---

## 🔍 **Debugging**

### **Backend Debugging**
```bash
# Start with debugger port exposed
docker-compose -f docker-compose.local.yml up -d

# Connect debugger to localhost:2345
```

### **View Logs**
```bash
# All services
docker-compose -f docker-compose.local.yml logs -f

# Specific service
docker-compose -f docker-compose.local.yml logs -f backend
docker-compose -f docker-compose.local.yml logs -f frontend
docker-compose -f docker-compose.local.yml logs -f postgres
```

### **Access Containers**
```bash
# Backend container
docker-compose -f docker-compose.local.yml exec backend sh

# Database container
docker-compose -f docker-compose.local.yml exec postgres psql -U postgres -d peace_local
```

---

## 🧹 **Cleanup Commands**

```bash
# Stop and remove containers
make docker-down

# Stop and remove containers + volumes
make docker-clean

# Remove all unused Docker resources
docker system prune -a

# Remove specific volumes
docker volume rm peace_postgres_local_data
docker volume rm peace_redis_local_data
```

---

## 🚨 **Troubleshooting**

### **Port Conflicts**
```bash
# Check what's using ports
lsof -i :3000
lsof -i :8080
lsof -i :5432
lsof -i :6379

# Stop conflicting services
sudo lsof -ti:3000 | xargs kill -9
```

### **Database Connection Issues**
```bash
# Check database health
docker-compose -f docker-compose.local.yml exec postgres pg_isready -U postgres

# Reset database
docker-compose -f docker-compose.local.yml down -v
docker-compose -f docker-compose.local.yml up -d
```

### **Build Issues**
```bash
# Rebuild without cache
docker-compose -f docker-compose.local.yml build --no-cache

# Clean build
make docker-clean
make docker-build
```

---

## 📝 **Environment Variables**

### **Required for Production**
```bash
# Copy example file
cp env.production.example .env.production

# Edit with your values
nano .env.production
```

### **Key Variables**
- `POSTGRES_PASSWORD`: Database password
- `JWT_SECRET`: JWT signing secret
- `REDIS_PASSWORD`: Redis password (production only)
- `ACME_EMAIL`: Email for Let's Encrypt SSL certificates

---

## 🎯 **Best Practices**

1. **Use `docker-compose.local.yml`** for local development
2. **Always use `make` commands** for consistency
3. **Check logs** when things don't work
4. **Clean up regularly** with `make docker-clean`
5. **Use environment files** for production
6. **Test locally first** before deploying

---

## 🔗 **Useful Links**

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Traefik Documentation](https://doc.traefik.io/traefik/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Redis Docker Image](https://hub.docker.com/_/redis)
