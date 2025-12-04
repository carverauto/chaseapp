# API Specification

## ADDED Requirements

### Requirement: Go Monolithic API Server

The system SHALL provide a single Go HTTP server using Gorilla Mux that handles all API requests for the ChaseApp platform.

#### Scenario: Server starts successfully
- **WHEN** the API server process starts
- **THEN** it SHALL bind to the configured port
- **AND** it SHALL establish connections to PostgreSQL, NATS, and Typesense
- **AND** it SHALL respond to health check requests at `/health`

#### Scenario: Graceful shutdown
- **WHEN** the server receives SIGTERM or SIGINT
- **THEN** it SHALL stop accepting new requests
- **AND** it SHALL complete in-flight requests within 30 seconds
- **AND** it SHALL close database and messaging connections cleanly

---

### Requirement: Chase Management API

The system SHALL provide RESTful endpoints for managing chase resources.

#### Scenario: List chases with pagination
- **WHEN** a client sends `GET /api/v1/chases?page=1&limit=20`
- **THEN** the system SHALL return a paginated list of chases
- **AND** the response SHALL include total count and pagination metadata

#### Scenario: Create a new chase
- **WHEN** an authenticated client sends `POST /api/v1/chases` with valid chase data
- **THEN** the system SHALL create the chase in PostgreSQL
- **AND** the system SHALL publish a `chases.created` event to NATS
- **AND** the system SHALL return the created chase with HTTP 201

#### Scenario: Update chase live status
- **WHEN** a chase is updated with `Live: true`
- **THEN** the system SHALL publish a `chases.live` event to NATS
- **AND** subscribers SHALL receive push notifications

#### Scenario: Chase ends
- **WHEN** a chase `Live` field changes from true to false
- **THEN** the system SHALL set `EndedAt` timestamp automatically
- **AND** the system SHALL publish a `chases.ended` event to NATS

---

### Requirement: Aircraft Tracking API

The system SHALL provide endpoints for aircraft data and clustering.

#### Scenario: List aircraft in region
- **WHEN** a client sends `GET /api/v1/aircraft?bounds=lat1,lon1,lat2,lon2`
- **THEN** the system SHALL return aircraft within the bounding box

#### Scenario: Cluster aircraft using DBSCAN
- **WHEN** a client sends `POST /api/v1/aircraft/cluster` with ADSB data
- **THEN** the system SHALL apply DBSCAN clustering algorithm
- **AND** the system SHALL return clusters with media aircraft highlighted

---

### Requirement: External Data Aggregation API

The system SHALL provide endpoints for aggregating external data sources.

#### Scenario: Fetch earthquake data
- **WHEN** a client sends `GET /api/v1/quakes`
- **THEN** the system SHALL fetch data from USGS Earthquake API
- **AND** the system SHALL return GeoJSON formatted earthquake data

#### Scenario: Fetch vessel data
- **WHEN** a client sends `GET /api/v1/boats`
- **THEN** the system SHALL fetch data from AISHub API
- **AND** the system SHALL return vessel positions and metadata

#### Scenario: Fetch rocket launch data
- **WHEN** a client sends `GET /api/v1/launches`
- **THEN** the system SHALL return upcoming and recent rocket launches

#### Scenario: Fetch weather alerts
- **WHEN** a client sends `GET /api/v1/weather/alerts`
- **THEN** the system SHALL fetch active alerts from NOAA/NWS
- **AND** the system SHALL return alerts with GeoJSON geometries

---

### Requirement: Stream URL Extraction API

The system SHALL provide endpoints for extracting live stream URLs from news networks.

#### Scenario: Extract stream URLs
- **WHEN** a client sends `POST /api/v1/streams/extract` with a chase ID
- **THEN** the system SHALL scrape configured news network pages
- **AND** the system SHALL return extracted m3u8 and video URLs
- **AND** the system SHALL update the chase record with stream URLs

#### Scenario: Supported networks
- **WHEN** extracting streams
- **THEN** the system SHALL support NBC LA, ABC7, and CBS News
- **AND** the system SHALL handle network-specific page structures

---

### Requirement: Geospatial Utilities API

The system SHALL provide geospatial calculation endpoints.

#### Scenario: Calculate minimum bounding rectangle
- **WHEN** a client sends `POST /api/v1/geo/bounding-rect` with GeoJSON features
- **THEN** the system SHALL calculate the minimum-area bounding rectangle
- **AND** the system SHALL return the rectangle as a GeoJSON polygon

---

### Requirement: Authentication Token API

The system SHALL provide endpoints for generating service-specific tokens.

#### Scenario: Generate chat authentication token
- **WHEN** an authenticated client sends `POST /api/v1/auth/chat-token`
- **THEN** the system SHALL generate a signed JWT for the chat service
- **AND** the token SHALL include user ID and permissions

---

### Requirement: Push Notification API

The system SHALL provide endpoints for managing push notification subscriptions and delivery.

#### Scenario: Subscribe to push notifications
- **WHEN** a client sends `POST /api/v1/push/subscribe` with device token
- **THEN** the system SHALL store the token in PostgreSQL
- **AND** the system SHALL associate the token with the user

#### Scenario: Unsubscribe from push notifications
- **WHEN** a client sends `POST /api/v1/push/unsubscribe` with device token
- **THEN** the system SHALL remove the token from the database

#### Scenario: Generate Safari push package
- **WHEN** a client sends `GET /api/v1/push/safari-package`
- **THEN** the system SHALL generate a signed ZIP package
- **AND** the package SHALL contain website.json, manifest.json, icons, and signature

#### Scenario: Deliver push notification
- **WHEN** a `chases.live` event is received
- **THEN** the system SHALL send push notifications via ntfy/Gotify
- **AND** the system SHALL send APNs notifications to iOS devices
- **AND** the system SHALL send FCM notifications to Android devices

---

### Requirement: Webhook Integration API

The system SHALL provide endpoints for sending webhooks to external services.

#### Scenario: Send Discord webhook
- **WHEN** a client sends `POST /api/v1/webhooks/discord` with payload
- **THEN** the system SHALL format an embed message
- **AND** the system SHALL POST to the configured Discord webhook URL

---

### Requirement: Search Integration

The system SHALL maintain a search index using Typesense.

#### Scenario: Index chase on creation
- **WHEN** a `chases.created` event is received
- **THEN** the system SHALL index the chase document in Typesense

#### Scenario: Update index on chase update
- **WHEN** a `chases.updated` event is received
- **THEN** the system SHALL update the chase document in Typesense

#### Scenario: Remove from index on deletion
- **WHEN** a `chases.deleted` event is received
- **THEN** the system SHALL remove the chase document from Typesense

---

### Requirement: Event-Driven Architecture

The system SHALL use NATS JetStream for asynchronous event processing.

#### Scenario: Publish domain events
- **WHEN** a chase is created, updated, or deleted
- **THEN** the system SHALL publish an event to the appropriate NATS subject
- **AND** the event SHALL include the full resource payload

#### Scenario: Subscribe to events
- **WHEN** the API server starts
- **THEN** it SHALL subscribe to configured NATS subjects
- **AND** it SHALL process events with at-least-once delivery semantics

#### Scenario: Event replay on failure
- **WHEN** an event handler fails
- **THEN** NATS JetStream SHALL redeliver the event
- **AND** the handler SHALL be idempotent to handle redelivery

---

### Requirement: API Authentication via Kong

The system SHALL delegate authentication to Kong API Gateway.

#### Scenario: Authenticated request
- **WHEN** Kong validates a JWT token successfully
- **THEN** Kong SHALL add `X-User-ID` and `X-User-Email` headers
- **AND** the Go API SHALL extract user context from these headers

#### Scenario: Unauthenticated request
- **WHEN** Kong rejects a request due to invalid/missing token
- **THEN** the request SHALL NOT reach the Go API
- **AND** Kong SHALL return HTTP 401

#### Scenario: Public endpoints
- **WHEN** a request is made to a public endpoint (health, certain GETs)
- **THEN** Kong SHALL allow the request without authentication

---

### Requirement: Observability

The system SHALL provide observability endpoints and instrumentation.

#### Scenario: Health check
- **WHEN** a client sends `GET /health`
- **THEN** the system SHALL return HTTP 200 if the server is running

#### Scenario: Readiness check
- **WHEN** a client sends `GET /ready`
- **THEN** the system SHALL return HTTP 200 if all dependencies are connected
- **AND** it SHALL return HTTP 503 if any dependency is unavailable

#### Scenario: Prometheus metrics
- **WHEN** a client sends `GET /metrics`
- **THEN** the system SHALL return Prometheus-formatted metrics
- **AND** metrics SHALL include request counts, latencies, and error rates

#### Scenario: Distributed tracing
- **WHEN** a request is processed
- **THEN** the system SHALL create OpenTelemetry spans
- **AND** trace IDs SHALL propagate through NATS messages
