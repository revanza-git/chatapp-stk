# Docker Setup Guide

This guide will help you run the entire Chat App using Docker containers.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed on your system
- [Docker Compose](https://docs.docker.com/compose/install/) (usually included with Docker Desktop)
- At least 4GB of available RAM for all services

## Quick Start

### For Linux/macOS:
```bash
chmod +x docker-start.sh
./docker-start.sh
```

### For Windows:
```cmd
docker-start.bat
```

## Manual Setup

### 1. Environment Configuration

Copy the example environment file and configure it:
```bash
cp .env.example .env
```

Edit `.env` with your preferred settings:
- Update `JWT_SECRET` for production use
- Modify database credentials if needed
- Adjust ports if there are conflicts

### 2. Start Services

Start all services using Docker Compose:
```bash
docker-compose up --build -d
```

### 3. Initialize Ollama (Optional)

Pull the required language model:
```bash
docker-compose exec ollama ollama pull llama2
```

## Services

The Docker setup includes the following services:

| Service | Port | Description |
|---------|------|-------------|
| **Frontend** | 3000 | Next.js React application |
| **Backend** | 8080 | Go API server |
| **PostgreSQL** | 5433 | Database |
| **Ollama** | 11434 | Language model server |

## Usage

Once all services are running:

1. **Frontend**: Visit http://localhost:3000
2. **Backend API**: Available at http://localhost:8080
3. **Database**: Connect to `localhost:5433` with credentials from `.env`
4. **Ollama**: API available at http://localhost:11434

## Development

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Rebuilding Services

```bash
# Rebuild and restart all services
docker-compose up --build

# Rebuild specific service
docker-compose up --build backend
```

### Database Access

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U chatbot_user -d chatbot_db
```

### Accessing Containers

```bash
# Backend container
docker-compose exec backend sh

# Frontend container  
docker-compose exec frontend sh
```

## Managing Data

### Persistent Volumes

Data is persisted in Docker volumes:
- `postgres-data`: Database files
- `ollama-data`: Downloaded models
- `uploads-data`: Uploaded documents

### Backup Database

```bash
docker-compose exec postgres pg_dump -U chatbot_user chatbot_db > backup.sql
```

### Restore Database

```bash
docker-compose exec -T postgres psql -U chatbot_user -d chatbot_db < backup.sql
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Update ports in `docker-compose.yml`
2. **Out of memory**: Increase Docker memory limit
3. **Services not starting**: Check logs with `docker-compose logs`

### Reset Everything

```bash
# Stop and remove all containers, networks, and volumes
docker-compose down -v

# Remove images (optional)
docker-compose down --rmi all -v

# Start fresh
docker-compose up --build -d
```

### Health Checks

```bash
# Check service status
docker-compose ps

# Test backend health
curl http://localhost:8080/api/health

# Test frontend
curl http://localhost:3000
```

## Production Considerations

1. **Security**:
   - Change default passwords in `.env`
   - Use a strong `JWT_SECRET`
   - Configure firewall rules

2. **Performance**:
   - Use Docker Swarm or Kubernetes for scaling
   - Configure resource limits
   - Use a reverse proxy (nginx/traefik)

3. **Monitoring**:
   - Add logging aggregation
   - Set up health check monitoring
   - Configure backup strategies

## Stopping Services

```bash
# Stop services (keeps containers)
docker-compose stop

# Stop and remove containers
docker-compose down

# Stop and remove everything including volumes
docker-compose down -v
``` 