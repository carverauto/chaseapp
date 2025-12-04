# ChaseApp API

Go monolithic API server for ChaseApp.

> For overall project documentation, see the [root README](../README.md).

## Architecture

This API is built as a monolithic Go service that consolidates the previous Firebase Cloud Functions into a single deployable unit. It follows a clean architecture pattern with:

- **Gorilla Mux** for HTTP routing
- **PostgreSQL** for persistent storage (via pgx)
- **NATS JetStream** for event-driven messaging
- **Kong** API Gateway for authentication and rate limiting

### Project Structure

```
api/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration loading
│   ├── database/        # Database connection and migrations
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # HTTP middleware (auth, logging, CORS, rate limiting)
│   ├── model/           # Domain models
│   └── repository/      # Data access layer
├── migrations/          # PostgreSQL migrations (golang-migrate)
├── pkg/                 # Shared packages (dbscan, geojson, scraper)
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- NATS 2.9+ (with JetStream)

## Quick Start

### 1. Clone and setup environment

```bash
cd api
cp .env.example .env
# Edit .env with your configuration
```

### 2. Start dependencies with Docker Compose

```bash
docker-compose up -d postgres nats typesense
```

This starts:
- PostgreSQL on port 5432
- NATS on port 4222 (with JetStream enabled)
- Typesense on port 8108
- MinIO on port 9000 (S3-compatible storage)
- ntfy on port 8090 (push notifications)

### 3. Run database migrations

```bash
# Install golang-migrate if needed
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://chaseapp:chaseapp_dev@localhost:5432/chaseapp?sslmode=disable" up
```

Or enable auto-migration by setting `DB_AUTO_MIGRATE=true` in your environment.

### 4. Run the server

```bash
go run cmd/server/main.go
```

The server starts on `http://localhost:8080` by default.

## API Endpoints

### Health & Metrics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check (always returns 200 if running) |
| GET | `/ready` | Readiness check (verifies database connectivity) |
| GET | `/metrics` | Prometheus metrics |

### Chases

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/chases` | List chases with pagination and filtering |
| POST | `/api/v1/chases` | Create a new chase |
| GET | `/api/v1/chases/{id}` | Get a single chase |
| PUT | `/api/v1/chases/{id}` | Update a chase |
| DELETE | `/api/v1/chases/{id}` | Delete a chase (soft delete) |
| GET | `/api/v1/chases/bundle` | Get offline data bundle |

**Query Parameters for List:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)
- `live` - Filter by live status (true/false)
- `type` - Filter by chase type (chase, rocket, weather, aircraft)
- `city` - Filter by city
- `state` - Filter by state

### Aircraft

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/aircraft` | List aircraft with filtering |
| POST | `/api/v1/aircraft/cluster` | DBSCAN clustering (WIP) |

**Query Parameters for List:**
- `page`, `limit` - Pagination
- `category` - Filter by category (media, law_enforcement, military, etc.)
- `cluster_id` - Filter by cluster
- `on_ground` - Filter by ground status
- `min_lat`, `max_lat`, `min_lng`, `max_lng` - Bounding box filter

### Push Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/push/subscribe` | Register device for push notifications |
| POST | `/api/v1/push/unsubscribe` | Unsubscribe from notifications |
| GET | `/api/v1/push/safari-package` | Safari push package (WIP) |

### External Data (WIP)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/quakes` | USGS earthquake data |
| GET | `/api/v1/boats` | AISHub vessel data |
| GET | `/api/v1/launches` | Rocket launch data |
| GET | `/api/v1/weather/alerts` | NOAA/NWS weather alerts |

### Other Endpoints (WIP)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/streams/extract` | Extract stream URLs from pages |
| POST | `/api/v1/geo/bounding-rect` | Calculate minimum bounding rectangle |
| POST | `/api/v1/auth/chat-token` | Generate chat service tokens |
| POST | `/api/v1/webhooks/discord` | Send Discord webhook |

## Configuration

Configuration is loaded from environment variables. See `.env.example` for all options.

### Server

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Listen host |
| `SERVER_PORT` | `8080` | Listen port |
| `SERVER_READ_TIMEOUT` | `30s` | HTTP read timeout |
| `SERVER_WRITE_TIMEOUT` | `30s` | HTTP write timeout |
| `SERVER_SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown timeout |

### Database

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `chaseapp` | Database user |
| `DB_PASSWORD` | - | Database password |
| `DB_NAME` | `chaseapp` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `DB_MAX_CONNS` | `25` | Max pool connections |
| `DB_AUTO_MIGRATE` | `false` | Run migrations on startup |
| `DB_MIGRATIONS_PATH` | `migrations` | Path to migration files |

### NATS

| Variable | Default | Description |
|----------|---------|-------------|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `NATS_CLUSTER_ID` | `chaseapp` | Cluster ID |
| `NATS_CLIENT_ID` | `api-server` | Client ID |

### Typesense

| Variable | Default | Description |
|----------|---------|-------------|
| `TYPESENSE_HOST` | `localhost` | Typesense host |
| `TYPESENSE_PORT` | `8108` | Typesense port |
| `TYPESENSE_API_KEY` | - | API key |

### Push Notifications

| Variable | Description |
|----------|-------------|
| `NTFY_URL` | ntfy/Gotify server URL |
| `APNS_KEY_ID` | Apple Push Notification Key ID |
| `APNS_TEAM_ID` | Apple Developer Team ID |
| `APNS_KEY_PATH` | Path to APNs auth key (.p8) |
| `APNS_BUNDLE_ID` | iOS app bundle ID |
| `FCM_PROJECT_ID` | Firebase project ID |
| `FCM_KEY_PATH` | Path to FCM service account JSON |

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o bin/api cmd/server/main.go
```

### Docker Build

```bash
docker build -t chaseapp-api .
```

### Full Docker Compose Stack

```bash
docker-compose up --build
```

## Authentication

The API integrates with Kong API Gateway for authentication. Kong validates tokens upstream and forwards user information via headers:

- `X-User-ID` - User's UUID
- `X-User-Email` - User's email address

The auth middleware extracts these headers and makes them available to handlers via request context.

## Database Migrations

Migrations use [golang-migrate](https://github.com/golang-migrate/migrate).

```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq create_table_name

# Apply migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback last migration
migrate -path migrations -database "$DATABASE_URL" down 1

# Check current version
migrate -path migrations -database "$DATABASE_URL" version
```

## Deployment

### Kubernetes

Kubernetes manifests are located in `k8s/` directory with Kustomize overlays for staging and production.

```bash
# Deploy to staging
kubectl apply -k k8s/overlays/staging

# Deploy to production
kubectl apply -k k8s/overlays/prod
```

### Environment Variables in K8s

Sensitive configuration should be stored in Kubernetes Secrets:

```bash
kubectl create secret generic api-secrets \
  --from-literal=DB_PASSWORD=your_password \
  --from-literal=TYPESENSE_API_KEY=your_key
```

## Migration from Firebase

This API replaces the following Firebase Cloud Functions:

| Old Function | New Endpoint |
|--------------|--------------|
| `API/ListChases` | `GET /api/v1/chases` |
| `API/AddChase` | `POST /api/v1/chases` |
| `API/UpdateChase` | `PUT /api/v1/chases/{id}` |
| `API/DeleteChase` | `DELETE /api/v1/chases/{id}` |
| `createBundle` | `GET /api/v1/chases/bundle` |
| `bof/findBofs` | `POST /api/v1/aircraft/cluster` |
| `manageTokens` | `POST /api/v1/push/subscribe` |
| `pushPackage` | `GET /api/v1/push/safari-package` |
| `rocketAPI/GetLaunches` | `GET /api/v1/launches` |
| `weatherAPI/GetWeatherAlerts` | `GET /api/v1/weather/alerts` |

## License

Proprietary - ChaseApp
