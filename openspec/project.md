# Project Context

## Purpose

ChaseApp is a real-time notification and chat platform for live events including:
- Police chases
- Rocket launches
- Weather events and disasters
- Aviation/aircraft tracking

The platform delivers real-time alerts and enables community discussion around breaking events through web, iOS, and Android applications.

**Website**: https://chaseapp.tv

## Tech Stack

### Frontend - Web (`web/`)
- **Framework**: Nuxt.js 2 (Bridge) + Vue.js 2/3
- **Language**: TypeScript 4.6
- **Styling**: TailwindCSS via WindiCSS
- **State Management**: Vuex + vuex-persistedstate
- **Mapping**: Mapbox GL JS (or MapLibre GL)
- **Real-time**: NATS WebSocket client
- **Search**: Typesense (self-hosted)
- **Build**: Nuxi, Vite
- **PWA**: @nuxtjs/pwa with service workers

### Frontend - Mobile (`mobile/`)
- **Framework**: Flutter 2.17+
- **Language**: Dart (SDK >=2.17.1 <3.0.0)
- **State Management**: Riverpod
- **Chat**: Custom WebSocket chat client
- **Mapping**: Mapbox GL
- **Push Notifications**: ntfy/Gotify client + APNs/FCM direct
- **Code Generation**: Freezed, json_serializable
- **Monetization**: RevenueCat (purchases_flutter)

### Backend (`api/`)
- **Language**: Go 1.21+
- **Architecture**: Microservices
- **API Gateway**: Kong
- **Authentication**: Custom OAuth 2.0 service
- **Database**: PostgreSQL
- **Messaging/Real-time**: NATS (JetStream for persistence)
- **Search**: Typesense
- **Object Storage**: MinIO (S3-compatible)
- **Push Notifications**: ntfy or Gotify (self-hosted)
- **Chat**: Custom Go WebSocket service

### Infrastructure
- **Build System**: Bazel monorepo
- **Deployment**: Self-hosted Kubernetes (Kustomize for base/staging/prod)
- **Container**: Docker
- **CI/CD**: GitHub Actions (or self-hosted runner)
- **SSL**: Let's Encrypt via cert-manager
- **Ingress**: Kong or nginx-ingress
- **Observability**: Prometheus + Grafana + Loki

## Project Conventions

### Code Style

#### TypeScript/JavaScript (Web)
- ESLint with vue-eslint-parser
- Extends: `plugin:nuxt/recommended`, `plugin:vue/recommended`, `@typescript-eslint/recommended`
- Trailing commas required on multiline (`comma-dangle: always-multiline`)
- Curly braces only for multi-line statements
- camelCase naming disabled (allows snake_case from APIs)
- No default props required for Vue components

#### Dart (Mobile)
- very_good_analysis lint rules
- dart_code_metrics for quality enforcement
- Freezed for immutable data classes
- json_serializable for JSON parsing

#### Go (Backend)
- Standard Go formatting (gofmt/goimports)
- golangci-lint for static analysis
- Service-per-domain microservice pattern
- Interface-based dependency injection
- Context propagation for tracing/cancellation

### Architecture Patterns

#### Web Application
- **SSR-first**: Server-side rendering enabled (`ssr: true`)
- **Plugin-based**: Client/server plugins in `plugins/` directory
- **Component organization**: Grouped by feature (cards, chat, birds, firehose, events)
- **Composables**: Nuxt composition API patterns

#### Mobile Application
- **Feature-based modules**: Code organized under `lib/src/modules/`
- **Riverpod providers**: State management via providers
- **Repository pattern**: Data access abstraction
- **Custom shaders**: GLSL shaders for visual effects

#### Backend Services
- **Microservices**: Domain-driven Go services communicating via NATS
- **API Gateway**: Kong for routing, rate limiting, and auth enforcement
- **Event-driven**: NATS pub/sub for inter-service communication
- **CQRS-lite**: Separate read/write paths where beneficial
- **Repository pattern**: Database access abstraction with PostgreSQL

### Testing Strategy

#### Web
- Vue Test Utils (@vue/test-utils)
- Lighthouse CI for performance auditing
- Playwright for E2E testing

#### Mobile
- Mock-based unit tests for auth and data layers
- Integration tests via Flutter test framework
- Widget testing for UI components

#### Backend
- Go standard testing package
- testcontainers-go for integration tests (PostgreSQL, NATS, etc.)
- golangci-lint in CI pipeline
- Docker Compose for local development environment

### Git Workflow

- **Main branch**: `main` (protected)
- **Feature branches**: `<type>/<description>` (e.g., `chore/k8s`, `feature/new-chat`)
- **Branch naming**: lowercase with hyphens, prefixed by type
- **Merges**: Pull request workflow via GitHub

## Domain Context

### Event Types
- **Chases**: Live police pursuit tracking with real-time location updates
- **Rockets**: Space launch notifications (SpaceX, NASA, etc.)
- **Weather**: Severe weather alerts, storm tracking
- **Birds/Aircraft**: Aviation tracking, airport monitoring

### Real-time Features
- Push notifications via ntfy/Gotify (self-hosted) + direct APNs/FCM
- WebSocket connections via NATS for live event updates
- Custom chat rooms per event (Go WebSocket service)
- Activity feeds via custom implementation

### User Features
- Social authentication (Google, Apple, Facebook, Twitter) via custom OAuth
- User profiles and preferences
- Subscription/premium features via RevenueCat
- Full-text search via Typesense

## Important Constraints

### Technical
- Go 1.21+ for backend services
- Flutter SDK >=2.17.1 <3.0.0
- PostgreSQL 15+
- NATS 2.9+ with JetStream enabled
- Node.js 18+ for web build tools (Volta pinned)

### Platform
- iOS: Sign in with Apple required for App Store
- Android: Google Play policies for notifications
- Web: HTTPS required (local dev uses mkcert certificates)

### Performance
- Map rendering performance critical for chase tracking
- Real-time latency requirements for notifications (<100ms target)
- PWA offline support for web
- NATS message delivery within 50ms

### Infrastructure
- Self-hosted Kubernetes cluster
- No vendor lock-in to cloud providers
- All services must be containerized
- Secrets managed via Kubernetes secrets or Vault

## External Dependencies

### Self-Hosted Services
- **PostgreSQL**: Primary database
- **NATS**: Messaging and real-time pub/sub (JetStream for persistence)
- **Typesense**: Full-text search engine
- **MinIO**: S3-compatible object storage
- **Kong**: API gateway and auth enforcement
- **ntfy/Gotify**: Push notification server
- **Prometheus/Grafana/Loki**: Observability stack

### Third-Party Services (Minimal)
- **Mapbox**: Maps and geocoding (or self-hosted tile server)
- **RevenueCat**: In-app purchases and subscriptions
- **APNs**: Apple Push Notification service (direct integration)
- **FCM**: Firebase Cloud Messaging for Android push (direct integration, no Firebase SDK)

### APIs
- Weather data providers
- Aviation/flight tracking
- Rocket launch schedules
