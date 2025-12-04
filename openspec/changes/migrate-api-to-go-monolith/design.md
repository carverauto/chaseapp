# Design: Go Monolithic API Architecture

## Context

ChaseApp currently runs 20+ Firebase Cloud Functions across three languages (Node.js, TypeScript, Go). This architecture has served the product but introduces:

- **Cold start latency**: Serverless functions have 500ms-2s cold starts
- **Vendor lock-in**: Deep Firebase/GCP dependency
- **Operational complexity**: Multiple deployment targets, mixed languages
- **Cost inefficiency**: Pay-per-invocation adds up with real-time features
- **Testing difficulty**: Firebase emulator limitations

The team wants to self-host on Kubernetes with full control over infrastructure.

## Goals / Non-Goals

### Goals
- Single deployable Go binary for all API functionality
- PostgreSQL as primary database (replacing Firestore)
- NATS for real-time messaging and event-driven workflows
- Kong API gateway for auth enforcement and rate limiting
- Horizontal scalability via Kubernetes
- Clean separation of concerns with domain-driven packages

### Non-Goals
- Microservices architecture (monolith-first approach)
- GraphQL API (REST with JSON is sufficient)
- Real-time database sync (NATS pub/sub replaces RTDB)
- Preserving Firebase SDK compatibility in clients

## Decisions

### 1. Router: Gorilla Mux

**Decision**: Use `github.com/gorilla/mux` for HTTP routing.

**Rationale**:
- Battle-tested, widely adopted in Go ecosystem
- Supports path variables, query params, method matching
- Middleware chaining for auth, logging, CORS
- No magic, explicit route registration

**Alternatives considered**:
- `chi`: Similar features, slightly newer. Gorilla has more community examples.
- `gin`: Higher performance but more opinionated. Unnecessary for our scale.
- `echo`: Good but smaller ecosystem.

### 2. Project Structure

```
api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/                  # Configuration loading
│   ├── server/                  # HTTP server setup
│   ├── middleware/              # Auth, logging, CORS, rate limiting
│   ├── handler/                 # HTTP handlers (thin layer)
│   │   ├── chase.go
│   │   ├── aircraft.go
│   │   ├── weather.go
│   │   ├── stream.go
│   │   ├── webhook.go
│   │   └── push.go
│   ├── service/                 # Business logic
│   │   ├── chase/
│   │   ├── aircraft/
│   │   ├── weather/
│   │   ├── stream/
│   │   ├── notification/
│   │   └── search/
│   ├── repository/              # Data access (PostgreSQL)
│   │   ├── chase.go
│   │   ├── user.go
│   │   └── aircraft.go
│   ├── worker/                  # Background jobs
│   │   ├── indexer.go           # Typesense indexing
│   │   ├── stats.go             # Statistics aggregation
│   │   └── scraper.go           # Stream URL extraction
│   ├── realtime/                # NATS pub/sub
│   │   ├── publisher.go
│   │   └── subscriber.go
│   └── external/                # Third-party integrations
│       ├── noaa/                # Weather alerts
│       ├── usgs/                # Earthquakes
│       ├── aishub/              # Vessel tracking
│       ├── discord/             # Webhooks
│       └── apns/                # Apple push
├── pkg/
│   ├── dbscan/                  # Clustering algorithm
│   ├── geojson/                 # GeoJSON utilities
│   └── scraper/                 # Web scraping (Colly)
├── migrations/                  # PostgreSQL migrations
├── Dockerfile
├── go.mod
└── go.sum
```

### 3. Database: PostgreSQL with pgx

**Decision**: Use PostgreSQL with `jackc/pgx` driver and raw SQL.

**Rationale**:
- Direct SQL gives full control over queries
- pgx is the fastest PostgreSQL driver for Go
- PostGIS extension for geospatial queries
- JSONB columns for flexible schema where needed

**Schema approach**:
- Normalized tables for core entities (users, chases, aircraft)
- JSONB for semi-structured data (weather alerts, stream metadata)
- Database migrations via `golang-migrate/migrate`

### 4. Event-Driven Patterns via NATS

**Decision**: Use NATS JetStream for async workflows replacing Firestore triggers.

**Patterns**:
```
# Chase lifecycle events
chases.created    → Index in Typesense, notify subscribers
chases.updated    → Update search index, check live status
chases.ended      → Send summary notification, update stats

# User events
users.created     → Send Discord webhook, initialize preferences

# Aircraft events
aircraft.clustered → Update BoF data, notify interested users
```

**Rationale**:
- Decouples event producers from consumers
- JetStream provides persistence and replay
- Native Go client with low latency
- Scales horizontally

### 5. Authentication Flow

**Decision**: Kong validates JWTs, Go API trusts validated requests.

**Flow**:
```
Client → Kong (JWT validation) → Go API (trusts X-User-ID header)
```

**Implementation**:
- Kong JWT plugin validates tokens from custom OAuth service
- Kong adds `X-User-ID`, `X-User-Email` headers to upstream
- Go middleware extracts user context from headers
- Internal service-to-service calls use separate auth

### 6. API Versioning

**Decision**: URL path versioning (`/api/v1/...`)

**Routes**:
```go
r := mux.NewRouter()
api := r.PathPrefix("/api/v1").Subrouter()

// Chases
api.HandleFunc("/chases", h.ListChases).Methods("GET")
api.HandleFunc("/chases", h.CreateChase).Methods("POST")
api.HandleFunc("/chases/{id}", h.GetChase).Methods("GET")
api.HandleFunc("/chases/{id}", h.UpdateChase).Methods("PUT")
api.HandleFunc("/chases/{id}", h.DeleteChase).Methods("DELETE")
api.HandleFunc("/chases/bundle", h.GetChasesBundle).Methods("GET")

// Aircraft
api.HandleFunc("/aircraft", h.ListAircraft).Methods("GET")
api.HandleFunc("/aircraft/cluster", h.ClusterAircraft).Methods("POST")

// Weather
api.HandleFunc("/weather/alerts", h.GetWeatherAlerts).Methods("GET")

// External data
api.HandleFunc("/quakes", h.GetQuakes).Methods("GET")
api.HandleFunc("/boats", h.GetBoats).Methods("GET")
api.HandleFunc("/launches", h.GetLaunches).Methods("GET")

// Streams
api.HandleFunc("/streams/extract", h.ExtractStreamURLs).Methods("POST")

// Geo utilities
api.HandleFunc("/geo/bounding-rect", h.GetBoundingRectangle).Methods("POST")

// Auth
api.HandleFunc("/auth/chat-token", h.GetChatToken).Methods("POST")

// Push notifications
api.HandleFunc("/push/safari-package", h.GetSafariPushPackage).Methods("GET")
api.HandleFunc("/push/subscribe", h.SubscribePush).Methods("POST")
api.HandleFunc("/push/unsubscribe", h.UnsubscribePush).Methods("POST")

// Webhooks
api.HandleFunc("/webhooks/discord", h.SendDiscordWebhook).Methods("POST")

// Health
r.HandleFunc("/health", h.HealthCheck).Methods("GET")
r.HandleFunc("/ready", h.ReadinessCheck).Methods("GET")
```

### 7. Background Workers

**Decision**: Run workers as goroutines within the same process, coordinated via NATS.

**Workers**:
- **Indexer**: Listens to chase events, updates Typesense
- **Stats aggregator**: Periodic stats calculation
- **Stream scraper**: On-demand URL extraction triggered by NATS
- **Weather poller**: Periodic NOAA/NWS fetching
- **Aircraft sync**: Periodic ADSB data processing

**Rationale**:
- Simpler deployment (single binary)
- Shared connection pools
- Easy local development
- Can extract to separate services later if needed

### 8. Observability

**Decision**: OpenTelemetry for tracing, Prometheus for metrics.

**Implementation**:
- `go.opentelemetry.io/otel` for distributed tracing
- `/metrics` endpoint for Prometheus scraping
- Structured logging with `log/slog` (Go 1.21+)
- Trace IDs propagated through NATS messages

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Monolith becomes unwieldy | Clean package boundaries, can extract services later |
| Single point of failure | Multiple replicas behind Kong, health checks |
| Database bottleneck | Connection pooling, read replicas if needed |
| Migration data loss | Parallel running period, comprehensive testing |
| Client breaking changes | Version API, deprecation period for old endpoints |

## Migration Plan

### Phase 1: Foundation
1. Set up Go project structure with Gorilla Mux
2. Implement PostgreSQL repository layer
3. Set up NATS connection and basic pub/sub
4. Deploy to staging K8s cluster

### Phase 2: Core APIs
1. Migrate chase CRUD operations
2. Migrate aircraft endpoints
3. Migrate weather/quakes/boats/launches
4. Implement search indexing with Typesense

### Phase 3: Real-time & Push
1. Implement NATS-based event handlers
2. Set up ntfy/Gotify push service
3. Migrate WebSocket chat functionality
4. Implement Safari push package generation

### Phase 4: Cutover
1. Update web client to use new API
2. Update mobile client to use new API
3. Run parallel with Firebase (shadow mode)
4. Decommission Firebase functions

### Rollback Plan
- Keep Firebase functions deployed during migration
- Feature flags to switch between old/new APIs
- Database sync from PostgreSQL back to Firestore if needed

## Open Questions

1. **Chat service**: Build custom WebSocket service or evaluate Matrix/other?
2. **Push notification tokens**: Migration strategy for existing FCM tokens?
3. **Data migration**: ETL from Firestore to PostgreSQL approach?
4. **Safari push**: Continue supporting or deprecate Safari web push?
