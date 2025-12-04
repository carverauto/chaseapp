# Change: Migrate Firebase Functions to Go Monolithic API

## Why

The current API layer consists of 20+ Firebase Cloud Functions spread across Node.js, TypeScript, and Go, creating operational complexity, cold-start latency, and vendor lock-in to Google Cloud. Consolidating into a single Go monolithic API with Gorilla Mux will reduce infrastructure costs, improve performance, simplify deployment, and enable self-hosted Kubernetes deployment.

## What Changes

### **BREAKING** - Complete Backend Architecture Migration

- **Remove all Firebase Cloud Functions** - Eliminate 20+ serverless functions
- **Create single Go monolith** - Unified API service using Gorilla Mux router
- **Replace Firestore triggers** - Convert to PostgreSQL with event-driven patterns via NATS
- **Replace Firebase Auth** - Integrate with custom OAuth service via Kong gateway
- **Replace FCM/Pusher** - Use NATS for real-time + ntfy/Gotify for push notifications
- **Replace Algolia** - Use Typesense for search indexing
- **Replace Stream Chat** - Custom WebSocket chat service

### Functions Being Migrated

| Current Function | New Endpoint/Handler | Notes |
|-----------------|---------------------|-------|
| API (Go) | Core routes | Expand existing Go code |
| getStreams (Go) | `POST /api/v1/streams/extract` | Stream URL scraping |
| rocketAPI (Go) | `/api/v1/quakes`, `/api/v1/boats`, `/api/v1/launches` | External data aggregation |
| weatherAPI (Go) | `GET /api/v1/weather/alerts` | NOAA/NWS integration |
| bof/bofTS | `POST /api/v1/aircraft/cluster` | DBSCAN clustering |
| algoliaIndex | Background worker | Typesense indexing |
| tokens | `POST /api/v1/auth/chat-token` | Chat authentication |
| webhooks | `POST /api/v1/webhooks/discord` | Discord notifications |
| pushPackage | `GET /api/v1/push/safari-package` | Safari push package |
| getRectangle | `POST /api/v1/geo/bounding-rect` | Geospatial calculation |
| stats | Background worker | Statistics aggregation |
| updateChase | Database trigger via NATS | Chase lifecycle |
| updateRTDB | Background sync worker | Aircraft data sync |
| addUser | Event handler via NATS | User registration events |
| fcmMessaging | Push service | Notification delivery |
| manageTokens | Push service | Token management |
| createBundle | `GET /api/v1/chases/bundle` | Offline data bundle |
| chaseStats | Background worker | Statistics |
| fireGetMP4LinkTrigger | Background worker | Media extraction |

### Functions Being Removed (No Migration)

- `firebaseBackups` - Not needed without Firebase
- `firebaseBackupsNuevo` - Not needed without Firebase
- `pusher` (empty) - Abandoned placeholder
- `getStreamsTest` - Development tool only

## Impact

- **Affected specs**: `api` (new capability)
- **Affected code**:
  - `api/` - All subdirectories will be consolidated/archived
  - `web/` - API client calls need updating
  - `mobile/` - API client calls need updating
  - `k8s/` - New deployment manifests needed
- **Infrastructure**:
  - Remove Firebase project dependency for backend
  - Add PostgreSQL, NATS, Typesense, MinIO to K8s cluster
  - Configure Kong API gateway routes
