# Tasks: Migrate Firebase Functions to Go Monolithic API

## 1. Project Setup

- [ ] 1.1 Initialize Go module at `api/` with `go mod init chaseapp.tv/api`
- [ ] 1.2 Create directory structure (`cmd/`, `internal/`, `pkg/`, `migrations/`)
- [ ] 1.3 Add Gorilla Mux dependency (`github.com/gorilla/mux`)
- [ ] 1.4 Add PostgreSQL driver (`github.com/jackc/pgx/v5`)
- [ ] 1.5 Add NATS client (`github.com/nats-io/nats.go`)
- [ ] 1.6 Set up configuration loading with environment variables
- [ ] 1.7 Create Dockerfile for the API service
- [ ] 1.8 Create docker-compose.yml for local development (PostgreSQL, NATS, Typesense)

## 2. Database Layer

- [ ] 2.1 Design PostgreSQL schema for chases table
- [ ] 2.2 Design PostgreSQL schema for users table
- [ ] 2.3 Design PostgreSQL schema for aircraft/ADSB table
- [ ] 2.4 Design PostgreSQL schema for push tokens table
- [ ] 2.5 Design PostgreSQL schema for statistics table
- [ ] 2.6 Create initial migration files with golang-migrate
- [ ] 2.7 Implement chase repository (CRUD operations)
- [ ] 2.8 Implement user repository
- [ ] 2.9 Implement aircraft repository
- [ ] 2.10 Implement token repository

## 3. Core HTTP Server

- [ ] 3.1 Create `cmd/server/main.go` entry point
- [ ] 3.2 Implement graceful shutdown handling
- [ ] 3.3 Set up Gorilla Mux router with `/api/v1` prefix
- [ ] 3.4 Implement CORS middleware
- [ ] 3.5 Implement request logging middleware
- [ ] 3.6 Implement auth middleware (extract user from Kong headers)
- [ ] 3.7 Implement rate limiting middleware
- [ ] 3.8 Implement health check endpoint (`/health`)
- [ ] 3.9 Implement readiness check endpoint (`/ready`)

## 4. Chase Endpoints (from API, updateChase, chaseStats, createBundle)

- [ ] 4.1 Implement `GET /api/v1/chases` - List chases with pagination
- [ ] 4.2 Implement `POST /api/v1/chases` - Create chase
- [ ] 4.3 Implement `GET /api/v1/chases/{id}` - Get single chase
- [ ] 4.4 Implement `PUT /api/v1/chases/{id}` - Update chase
- [ ] 4.5 Implement `DELETE /api/v1/chases/{id}` - Delete chase
- [ ] 4.6 Implement `GET /api/v1/chases/bundle` - Offline data bundle
- [ ] 4.7 Implement chase lifecycle logic (Live status, EndedAt timestamp)
- [ ] 4.8 Publish chase events to NATS on create/update/delete

## 5. Aircraft Endpoints (from API, bof, bofTS, updateRTDB)

- [ ] 5.1 Implement `GET /api/v1/aircraft` - List aircraft
- [ ] 5.2 Implement `POST /api/v1/aircraft/cluster` - DBSCAN clustering
- [ ] 5.3 Port DBSCAN algorithm to `pkg/dbscan/`
- [ ] 5.4 Implement BoF (Birds of a Feather) clustering logic
- [ ] 5.5 Implement aircraft data sync worker

## 6. External Data Endpoints (from rocketAPI, weatherAPI)

- [ ] 6.1 Implement `GET /api/v1/quakes` - USGS earthquake data
- [ ] 6.2 Implement `GET /api/v1/boats` - AISHub vessel data
- [ ] 6.3 Implement `GET /api/v1/launches` - Rocket launch data
- [ ] 6.4 Implement `GET /api/v1/weather/alerts` - NOAA/NWS alerts
- [ ] 6.5 Create external API clients in `internal/external/`
- [ ] 6.6 Implement caching layer for external API responses

## 7. Stream Extraction (from getStreams)

- [ ] 7.1 Implement `POST /api/v1/streams/extract` - Extract stream URLs
- [ ] 7.2 Port Colly web scraper to `pkg/scraper/`
- [ ] 7.3 Implement NBC LA stream extraction
- [ ] 7.4 Implement ABC7 m3u8 extraction
- [ ] 7.5 Implement CBS News video extraction
- [ ] 7.6 Add network-specific parsers as pluggable modules

## 8. Geospatial Utilities (from getRectangle)

- [ ] 8.1 Implement `POST /api/v1/geo/bounding-rect` - Minimum bounding rectangle
- [ ] 8.2 Port rotating calipers algorithm to `pkg/geojson/`
- [ ] 8.3 Add GeoJSON parsing utilities

## 9. Authentication & Tokens (from tokens)

- [ ] 9.1 Implement `POST /api/v1/auth/chat-token` - Generate chat tokens
- [ ] 9.2 Implement custom JWT signing for chat service
- [ ] 9.3 Integrate with Kong for upstream auth validation

## 10. Push Notifications (from fcmMessaging, manageTokens, pushPackage)

- [ ] 10.1 Implement `POST /api/v1/push/subscribe` - Subscribe to notifications
- [ ] 10.2 Implement `POST /api/v1/push/unsubscribe` - Unsubscribe
- [ ] 10.3 Implement `GET /api/v1/push/safari-package` - Safari push package
- [ ] 10.4 Port Safari push package generation (P12 signing, ZIP creation)
- [ ] 10.5 Implement ntfy/Gotify client for push delivery
- [ ] 10.6 Implement direct APNs client for iOS push
- [ ] 10.7 Implement direct FCM HTTP v1 API for Android push

## 11. Webhooks (from webhooks, addUser)

- [ ] 11.1 Implement `POST /api/v1/webhooks/discord` - Send Discord webhook
- [ ] 11.2 Create Discord embed builder utility
- [ ] 11.3 Implement user registration event handler (NATS subscriber)

## 12. Search Integration (from algoliaIndex)

- [ ] 12.1 Set up Typesense client
- [ ] 12.2 Create chase search schema in Typesense
- [ ] 12.3 Implement search indexing worker (NATS subscriber)
- [ ] 12.4 Implement `GET /api/v1/search` - Search endpoint (optional)

## 13. Background Workers

- [ ] 13.1 Implement worker manager for goroutine lifecycle
- [ ] 13.2 Implement Typesense indexer worker
- [ ] 13.3 Implement statistics aggregation worker
- [ ] 13.4 Implement weather polling worker
- [ ] 13.5 Implement MP4 link extraction worker (from fireGetMP4LinkTrigger)

## 14. NATS Event Handlers

- [ ] 14.1 Set up NATS JetStream streams and consumers
- [ ] 14.2 Implement `chases.created` event handler
- [ ] 14.3 Implement `chases.updated` event handler
- [ ] 14.4 Implement `chases.ended` event handler
- [ ] 14.5 Implement `users.created` event handler
- [ ] 14.6 Implement `aircraft.updated` event handler

## 15. Observability

- [ ] 15.1 Add OpenTelemetry tracing
- [ ] 15.2 Add Prometheus metrics endpoint (`/metrics`)
- [ ] 15.3 Implement structured logging with `log/slog`
- [ ] 15.4 Add request ID propagation

## 16. Testing

- [ ] 16.1 Set up testcontainers for integration tests
- [ ] 16.2 Write unit tests for repository layer
- [ ] 16.3 Write unit tests for service layer
- [ ] 16.4 Write integration tests for HTTP handlers
- [ ] 16.5 Write integration tests for NATS event handlers
- [ ] 16.6 Set up CI pipeline with tests

## 17. Kubernetes Deployment

- [ ] 17.1 Create K8s Deployment manifest for API
- [ ] 17.2 Create K8s Service manifest
- [ ] 17.3 Create K8s ConfigMap for configuration
- [ ] 17.4 Create K8s Secret for sensitive config
- [ ] 17.5 Update Kong Ingress routes
- [ ] 17.6 Create HorizontalPodAutoscaler
- [ ] 17.7 Add to Kustomize overlays (staging, prod)

## 18. Data Migration

- [ ] 18.1 Create Firestore to PostgreSQL migration script
- [ ] 18.2 Migrate chase data
- [ ] 18.3 Migrate user data
- [ ] 18.4 Migrate push token data
- [ ] 18.5 Validate data integrity post-migration

## 19. Client Updates

- [ ] 19.1 Update web client API base URL
- [ ] 19.2 Update web client API calls to new endpoints
- [ ] 19.3 Update mobile client API base URL
- [ ] 19.4 Update mobile client API calls to new endpoints
- [ ] 19.5 Remove Firebase SDK from clients (where applicable)

## 20. Cutover & Cleanup

- [ ] 20.1 Deploy to staging and validate
- [ ] 20.2 Run shadow mode (both APIs active)
- [ ] 20.3 Switch production traffic to new API
- [ ] 20.4 Monitor for errors and performance
- [ ] 20.5 Decommission Firebase Cloud Functions
- [ ] 20.6 Archive old `api/` function directories
- [ ] 20.7 Update documentation
