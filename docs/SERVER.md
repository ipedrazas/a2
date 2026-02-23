# A2 Server Mode

A2 can run as a web server with an HTTP API and React UI for on-demand repository analysis.

## Running the Server

```bash
# Run server locally (development)
a2 server

# Run with custom configuration
a2 server --port 3000 --host localhost --workspace-dir ./cache

# Run with Docker Compose
docker-compose up

# Run with Docker
docker run -p 8080:8080 -v a2-cache:/workspace/a2-cache a2-server:latest server
```

## Server Options

| Flag | Description | Default |
|------|-------------|---------|
| `--host` | Host to bind to | `0.0.0.0` |
| `--port` | Port to listen on | `8080` |
| `--workspace-dir` | Directory for cloned repos | `./a2-cache` |
| `--cleanup-after` | Clean up workspace after job | `true` |
| `--max-concurrent` | Maximum concurrent jobs | `5` |
| `--cleanup-interval` | Old job cleanup interval | `1h` |
| `--job-history-max-age` | Maximum job history age | `24h` |

## API Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Submit a check
curl -X POST http://localhost:8080/api/check \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/ipedrazas/a2",
    "profile": "api",
    "target": "production"
  }'

# Get job status
curl http://localhost:8080/api/check/{job_id}
```

## Web UI

The server includes a React UI accessible at `http://localhost:8080` with:
- GitHub URL input form
- Optional profile and target selection
- Real-time job status updates
- Results display with filtering and sorting
- JSON export capability

## Development Tasks

```bash
# Build UI
task server:ui:build

# Run UI dev server (hot reload)
task server:dev:ui

# Run server locally
task server:dev

# Build Docker image
task server:build

# Run with Docker Compose
docker-compose up
```
