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
	"chaseapp.tv/api/internal/handler"
	"chaseapp.tv/api/internal/middleware"
	"chaseapp.tv/api/internal/repository"
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
}

// New creates a new Server instance.
func New(cfg *config.Config, logger *slog.Logger, pool *pgxpool.Pool) *Server {
	// Initialize repositories
	chaseRepo := repository.NewChaseRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	aircraftRepo := repository.NewAircraftRepository(pool)
	pushTokenRepo := repository.NewPushTokenRepository(pool)

	s := &Server{
		cfg:    cfg,
		logger: logger,
		router: mux.NewRouter(),
		pool:   pool,

		// Initialize handlers with their dependencies
		chaseHandler:    handler.NewChaseHandler(chaseRepo, logger),
		aircraftHandler: handler.NewAircraftHandler(aircraftRepo, logger),
		pushHandler:     handler.NewPushHandler(pushTokenRepo, userRepo, logger),
	}

	s.setupMiddleware()
	s.setupRoutes()

	s.http = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return s
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

	// External data (stub handlers for now)
	api.HandleFunc("/quakes", handler.GetQuakes).Methods(http.MethodGet)
	api.HandleFunc("/boats", handler.GetBoats).Methods(http.MethodGet)
	api.HandleFunc("/launches", handler.GetLaunches).Methods(http.MethodGet)
	api.HandleFunc("/weather/alerts", handler.GetWeatherAlerts).Methods(http.MethodGet)

	// Streams
	api.HandleFunc("/streams/extract", handler.ExtractStreamURLs).Methods(http.MethodPost)

	// Geo utilities
	api.HandleFunc("/geo/bounding-rect", handler.GetBoundingRectangle).Methods(http.MethodPost)

	// Auth
	api.HandleFunc("/auth/chat-token", handler.GetChatToken).Methods(http.MethodPost)

	// Push notifications
	api.HandleFunc("/push/subscribe", s.pushHandler.Subscribe).Methods(http.MethodPost)
	api.HandleFunc("/push/unsubscribe", s.pushHandler.Unsubscribe).Methods(http.MethodPost)
	api.HandleFunc("/push/safari-package", handler.GetSafariPushPackage).Methods(http.MethodGet)

	// Webhooks
	api.HandleFunc("/webhooks/discord", handler.SendDiscordWebhook).Methods(http.MethodPost)
}

// readinessCheck verifies database connectivity.
func (s *Server) readinessCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.pool.Ping(r.Context()); err != nil {
		handler.Error(w, http.StatusServiceUnavailable, "database not ready")
		return
	}
	handler.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	s.logger.Info("starting server",
		slog.String("addr", s.http.Addr),
	)
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.http.Shutdown(ctx)
}

// Router returns the underlying mux router for testing.
func (s *Server) Router() *mux.Router {
	return s.router
}
