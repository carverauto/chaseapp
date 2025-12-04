// Package server provides the HTTP server setup and routing.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/external"
	"chaseapp.tv/api/internal/handler"
	"chaseapp.tv/api/internal/middleware"
	"chaseapp.tv/api/internal/realtime"
	"chaseapp.tv/api/internal/repository"
	"chaseapp.tv/api/internal/worker"
	"chaseapp.tv/api/pkg/scraper"
)

// Server represents the HTTP server.
type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	router *mux.Router
	http   *http.Server
	pool   *pgxpool.Pool

	// Handlers
	chaseHandler    *handler.ChaseHandler
	aircraftHandler *handler.AircraftHandler
	pushHandler     *handler.PushHandler
	externalHandler *handler.ExternalHandler
	streamHandler   *handler.StreamHandler
	geoHandler      *handler.GeoHandler
	authHandler     *handler.AuthHandler

	// Realtime
	publisher *realtime.Publisher

	// Workers
	aircraftWorker *worker.AircraftSyncWorker
}

// New creates a new Server instance.
func New(cfg *config.Config, logger *slog.Logger, pool *pgxpool.Pool) (*Server, error) {
	// Initialize repositories
	chaseRepo := repository.NewChaseRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	aircraftRepo := repository.NewAircraftRepository(pool)
	pushTokenRepo := repository.NewPushTokenRepository(pool)

	publisher, err := realtime.NewPublisher(cfg.NATS, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	externalClient := external.NewClient(cfg.External, logger)
	streamExtractor := scraper.NewExtractor()
	chatSigner, err := auth.NewChatTokenSigner(cfg.Chat)
	if err != nil {
		return nil, fmt.Errorf("chat token signer init: %w", err)
	}

	s := &Server{
		cfg:       cfg,
		logger:    logger,
		router:    mux.NewRouter(),
		pool:      pool,
		publisher: publisher,

		// Initialize handlers with their dependencies
		chaseHandler:    handler.NewChaseHandler(chaseRepo, publisher, logger),
		aircraftHandler: handler.NewAircraftHandler(aircraftRepo, logger),
		pushHandler:     handler.NewPushHandler(pushTokenRepo, userRepo, cfg.Push, logger),
		externalHandler: handler.NewExternalHandler(externalClient, logger),
		streamHandler:   handler.NewStreamHandler(chaseRepo, streamExtractor, publisher, logger),
		geoHandler:      handler.NewGeoHandler(logger),
		authHandler:     handler.NewAuthHandler(chatSigner, logger),

		// Workers
		aircraftWorker: worker.NewAircraftSyncWorker(aircraftRepo, logger),
	}

	s.setupMiddleware()
	s.setupRoutes()

	s.http = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return s, nil
}

// setupMiddleware configures global middleware.
func (s *Server) setupMiddleware() {
	// Apply middleware in order (outermost first)
	s.router.Use(middleware.Recovery(s.logger))
	s.router.Use(middleware.CORS([]string{"*"})) // TODO: Configure allowed origins
	s.router.Use(middleware.Metrics)
	s.router.Use(middleware.Logging(s.logger))
	s.router.Use(middleware.Auth)
}

// setupRoutes configures all HTTP routes.
func (s *Server) setupRoutes() {
	// Health and metrics endpoints (no auth required)
	s.router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)
	s.router.HandleFunc("/ready", s.readinessCheck).Methods(http.MethodGet)
	s.router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Chases
	api.HandleFunc("/chases", s.chaseHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/chases", s.chaseHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/chases/bundle", s.chaseHandler.GetBundle).Methods(http.MethodGet)
	api.HandleFunc("/chases/{id}", s.chaseHandler.Get).Methods(http.MethodGet)
	api.HandleFunc("/chases/{id}", s.chaseHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/chases/{id}", s.chaseHandler.Delete).Methods(http.MethodDelete)

	// Aircraft
	api.HandleFunc("/aircraft", s.aircraftHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/aircraft/cluster", s.aircraftHandler.Cluster).Methods(http.MethodPost)

	// External data
	api.HandleFunc("/quakes", s.externalHandler.GetQuakes).Methods(http.MethodGet)
	api.HandleFunc("/boats", s.externalHandler.GetBoats).Methods(http.MethodGet)
	api.HandleFunc("/launches", s.externalHandler.GetLaunches).Methods(http.MethodGet)
	api.HandleFunc("/weather/alerts", s.externalHandler.GetWeatherAlerts).Methods(http.MethodGet)

	// Streams
	api.HandleFunc("/streams/extract", s.streamHandler.ExtractStreamURLs).Methods(http.MethodPost)

	// Geo utilities
	api.HandleFunc("/geo/bounding-rect", s.geoHandler.GetBoundingRectangle).Methods(http.MethodPost)

	// Auth
	api.HandleFunc("/auth/chat-token", s.authHandler.GetChatToken).Methods(http.MethodPost)

	// Push notifications
	api.HandleFunc("/push/subscribe", s.pushHandler.Subscribe).Methods(http.MethodPost)
	api.HandleFunc("/push/unsubscribe", s.pushHandler.Unsubscribe).Methods(http.MethodPost)
	api.HandleFunc("/push/safari-package", s.pushHandler.GetSafariPushPackage).Methods(http.MethodGet)

	// Webhooks
	api.HandleFunc("/webhooks/discord", handler.SendDiscordWebhook).Methods(http.MethodPost)
}

// readinessCheck verifies database connectivity.
func (s *Server) readinessCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.pool.Ping(r.Context()); err != nil {
		handler.Error(w, http.StatusServiceUnavailable, "database not ready")
		return
	}
	if s.publisher == nil || !s.publisher.IsConnected() {
		handler.Error(w, http.StatusServiceUnavailable, "nats not ready")
		return
	}
	handler.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	s.logger.Info("starting server",
		slog.String("addr", s.http.Addr),
	)

	if s.aircraftWorker != nil {
		s.logger.Info("starting aircraft sync worker")
		s.aircraftWorker.Start()
	}

	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}
	if s.publisher != nil {
		s.publisher.Close()
	}
	if s.aircraftWorker != nil {
		s.aircraftWorker.Stop(ctx)
	}
	return nil
}

// Router returns the underlying mux router for testing.
func (s *Server) Router() *mux.Router {
	return s.router
}
