# ChaseApp

Real-time notification and community platform for live events.

**Website**: https://chaseapp.tv

## Overview

ChaseApp delivers real-time alerts and enables community discussion around breaking events:

- **Police Chases** - Live pursuit tracking with real-time location updates
- **Rocket Launches** - Space launch notifications (SpaceX, NASA, etc.)
- **Weather Events** - Severe weather alerts and storm tracking
- **Aircraft Tracking** - Aviation monitoring with clustering detection

The platform consists of web, iOS, and Android applications backed by a self-hosted infrastructure designed for low-latency real-time updates.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Clients                                  │
├─────────────┬─────────────┬─────────────────────────────────────┤
│   Web App   │  iOS App    │  Android App                        │
│   (Nuxt.js) │  (Flutter)  │  (Flutter)                          │
└──────┬──────┴──────┬──────┴──────┬──────────────────────────────┘
       │             │             │
       └─────────────┼─────────────┘
                     │ HTTPS
                     ▼
              ┌──────────────┐
              │     Kong     │  API Gateway
              │   (Ingress)  │  Auth, Rate Limiting
              └──────┬───────┘
                     │
       ┌─────────────┼─────────────┐
       │             │             │
       ▼             ▼             ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│   API    │  │   Chat   │  │   Auth   │
│  (Go)    │  │   (Go)   │  │  (OAuth) │
└────┬─────┘  └────┬─────┘  └──────────┘
     │             │
     └──────┬──────┘
            │
     ┌──────┴──────┐
     │    NATS     │  Event Bus
     │ (JetStream) │  Real-time Messaging
     └──────┬──────┘
            │
┌───────────┼───────────┬───────────────┐
│           │           │               │
▼           ▼           ▼               ▼
┌────────┐ ┌─────────┐ ┌──────────┐ ┌───────┐
│PostgreSQL│ │Typesense│ │  MinIO   │ │ ntfy  │
│   (DB)   │ │ (Search)│ │(Storage) │ │(Push) │
└──────────┘ └─────────┘ └──────────┘ └───────┘
```

## Project Structure

```
chaseapp/
├── api/                    # Go monolithic API server
│   ├── cmd/server/         # Application entry point
│   ├── internal/           # Private application code
│   │   ├── config/         # Configuration
│   │   ├── database/       # Database connection & migrations
│   │   ├── handler/        # HTTP handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── model/          # Domain models
│   │   └── repository/     # Data access layer
│   ├── migrations/         # PostgreSQL migrations
│   └── pkg/                # Shared packages
├── web/                    # Nuxt.js web application
│   ├── components/         # Vue components
│   ├── pages/              # Route pages
│   ├── plugins/            # Nuxt plugins
│   └── store/              # Vuex store
├── mobile/                 # Flutter mobile app (iOS & Android)
│   ├── lib/                # Dart source code
│   ├── ios/                # iOS-specific code
│   └── android/            # Android-specific code
├── k8s/                    # Kubernetes manifests
│   ├── base/               # Base manifests
│   ├── staging/            # Staging overlay
│   └── prod/               # Production overlay
├── shared/                 # Shared code/configs
├── openspec/               # API specifications & change proposals
└── tools/                  # Build tools & scripts
```

## Tech Stack

### Backend (`api/`)
| Component | Technology |
|-----------|------------|
| Language | Go 1.21+ |
| HTTP Router | Gorilla Mux |
| Database | PostgreSQL 15+ (pgx driver) |
| Migrations | golang-migrate |
| Messaging | NATS JetStream |
| Search | Typesense |
| Object Storage | MinIO (S3-compatible) |
| Push Notifications | ntfy/Gotify + APNs + FCM |

### Web (`web/`)
| Component | Technology |
|-----------|------------|
| Framework | Nuxt.js 2 (Bridge) + Vue.js |
| Language | TypeScript |
| Styling | TailwindCSS (WindiCSS) |
| State | Vuex |
| Maps | Mapbox GL JS |
| Real-time | NATS WebSocket |

### Mobile (`mobile/`)
| Component | Technology |
|-----------|------------|
| Framework | Flutter 2.17+ |
| Language | Dart |
| State | Riverpod |
| Maps | Mapbox GL |
| Push | APNs/FCM direct |
| Payments | RevenueCat |

### Infrastructure
| Component | Technology |
|-----------|------------|
| Build System | Bazel |
| Container Runtime | Docker |
| Orchestration | Kubernetes |
| Config Management | Kustomize |
| API Gateway | Kong |
| SSL | Let's Encrypt (cert-manager) |
| Observability | Prometheus + Grafana + Loki |

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+ (via Volta)
- Flutter SDK 2.17+
- Docker & Docker Compose
- Bazel (optional, for full builds)

### Quick Start (API Development)

```bash
# Clone the repository
git clone https://github.com/carverauto/chaseapp.git
cd chaseapp

# Start infrastructure services
cd api
cp .env.example .env
docker-compose up -d postgres nats typesense

# Run database migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations -database "postgres://chaseapp:chaseapp_dev@localhost:5432/chaseapp?sslmode=disable" up

# Start the API server
go run cmd/server/main.go
```

API available at `http://localhost:8080`

### Web Development

```bash
cd web
npm install
npm run dev
```

Web app available at `http://localhost:3000`

### Mobile Development

```bash
cd mobile
flutter pub get
flutter run
```

## API Endpoints

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |
| GET | `/metrics` | Prometheus metrics |

### Chases API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/chases` | List chases (paginated) |
| POST | `/api/v1/chases` | Create chase |
| GET | `/api/v1/chases/{id}` | Get chase |
| PUT | `/api/v1/chases/{id}` | Update chase |
| DELETE | `/api/v1/chases/{id}` | Delete chase |
| GET | `/api/v1/chases/bundle` | Offline data bundle |

### Aircraft API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/aircraft` | List aircraft |
| POST | `/api/v1/aircraft/cluster` | DBSCAN clustering |

### Push Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/push/subscribe` | Register device |
| POST | `/api/v1/push/unsubscribe` | Unregister device |

### External Data

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/quakes` | USGS earthquake data |
| GET | `/api/v1/launches` | Rocket launch data |
| GET | `/api/v1/weather/alerts` | NOAA weather alerts |
| GET | `/api/v1/boats` | AIS vessel data |

## Configuration

### Environment Variables

See `api/.env.example` for full configuration options.

**Key variables:**

```bash
# Server
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=chaseapp
DB_PASSWORD=secret
DB_NAME=chaseapp

# NATS
NATS_URL=nats://localhost:4222

# Typesense
TYPESENSE_HOST=localhost
TYPESENSE_API_KEY=your_key
```

## Deployment

### Docker

```bash
# Build API image
docker build -t chaseapp-api ./api

# Run with Docker Compose
docker-compose up
```

### Kubernetes

```bash
# Deploy to staging
kubectl apply -k k8s/staging

# Deploy to production
kubectl apply -k k8s/prod
```

## Development Workflow

### Branch Naming

- `feature/<description>` - New features
- `fix/<description>` - Bug fixes
- `chore/<description>` - Maintenance tasks

### Code Style

- **Go**: `gofmt` + `golangci-lint`
- **TypeScript**: ESLint + Prettier
- **Dart**: `very_good_analysis`

### Testing

```bash
# Go tests
cd api && go test ./...

# Web tests
cd web && npm test

# Mobile tests
cd mobile && flutter test
```

## Documentation

- `api/README.md` - API-specific documentation
- `openspec/` - API specifications and change proposals
- `openspec/project.md` - Detailed project context

## Contributing

1. Create a feature branch from `main`
2. Make changes following code style guidelines
3. Write/update tests as needed
4. Submit a pull request

## License

Proprietary - CarverAuto / ChaseApp
