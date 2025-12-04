// Package config provides configuration loading from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the API server.
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	NATS          NATSConfig
	Search        SearchConfig
	Push          PushConfig
	Chat          ChatConfig
	External      ExternalConfig
	Observability ObservabilityConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds PostgreSQL configuration.
type DatabaseConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	SSLMode     string
	MaxConns    int
	MinConns    int
	MaxConnLife time.Duration
	MaxConnIdle time.Duration
}

// ConnectionString returns the PostgreSQL connection string.
func (d DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Database, d.SSLMode,
	)
}

// NATSConfig holds NATS configuration.
type NATSConfig struct {
	URL           string
	ClusterID     string
	ClientID      string
	MaxReconnects int
	ReconnectWait time.Duration
}

// SearchConfig holds Typesense configuration.
type SearchConfig struct {
	Host     string
	Port     int
	Protocol string
	APIKey   string
}

// URL returns the Typesense server URL.
func (s SearchConfig) URL() string {
	return fmt.Sprintf("%s://%s:%d", s.Protocol, s.Host, s.Port)
}

// PushConfig holds push notification configuration.
type PushConfig struct {
	NtfyURL      string
	NtfyToken    string
	APNsKeyID    string
	APNsTeamID   string
	APNsKeyPath  string
	APNsBundleID string
	FCMProjectID string
	FCMKeyPath   string
	SafariPushID string
	SafariWebURL string
}

// ChatConfig holds chat token signing configuration.
type ChatConfig struct {
	SigningKey string
	Issuer     string
	Audience   string
	TokenTTL   time.Duration
}

// ExternalConfig holds external API configuration.
type ExternalConfig struct {
	USGSBaseURL          string
	AISHubBaseURL        string
	AISHubAPIKey         string
	NOAABaseURL          string
	LaunchLibraryBaseURL string
	DiscordWebhook       string
}

// ObservabilityConfig holds tracing/metrics settings.
type ObservabilityConfig struct {
	ServiceName  string
	OTLPEndpoint string
	OTLPHeaders  string
	OTLPInsecure bool
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnvInt("DB_PORT", 5432),
			User:        getEnv("DB_USER", "chaseapp"),
			Password:    getEnv("DB_PASSWORD", ""),
			Database:    getEnv("DB_NAME", "chaseapp"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			MaxConns:    getEnvInt("DB_MAX_CONNS", 25),
			MinConns:    getEnvInt("DB_MIN_CONNS", 5),
			MaxConnLife: getEnvDuration("DB_MAX_CONN_LIFE", time.Hour),
			MaxConnIdle: getEnvDuration("DB_MAX_CONN_IDLE", 30*time.Minute),
		},
		NATS: NATSConfig{
			URL:           getEnv("NATS_URL", "nats://localhost:4222"),
			ClusterID:     getEnv("NATS_CLUSTER_ID", "chaseapp"),
			ClientID:      getEnv("NATS_CLIENT_ID", "api-server"),
			MaxReconnects: getEnvInt("NATS_MAX_RECONNECTS", 10),
			ReconnectWait: getEnvDuration("NATS_RECONNECT_WAIT", 2*time.Second),
		},
		Search: SearchConfig{
			Host:     getEnv("TYPESENSE_HOST", "localhost"),
			Port:     getEnvInt("TYPESENSE_PORT", 8108),
			Protocol: getEnv("TYPESENSE_PROTOCOL", "http"),
			APIKey:   getEnv("TYPESENSE_API_KEY", ""),
		},
		Push: PushConfig{
			NtfyURL:      getEnv("NTFY_URL", ""),
			NtfyToken:    getEnv("NTFY_TOKEN", ""),
			APNsKeyID:    getEnv("APNS_KEY_ID", ""),
			APNsTeamID:   getEnv("APNS_TEAM_ID", ""),
			APNsKeyPath:  getEnv("APNS_KEY_PATH", ""),
			APNsBundleID: getEnv("APNS_BUNDLE_ID", ""),
			FCMProjectID: getEnv("FCM_PROJECT_ID", ""),
			FCMKeyPath:   getEnv("FCM_KEY_PATH", ""),
			SafariPushID: getEnv("SAFARI_PUSH_ID", ""),
			SafariWebURL: getEnv("SAFARI_WEB_SERVICE_URL", ""),
		},
		Chat: ChatConfig{
			SigningKey: getEnv("CHAT_SIGNING_KEY", ""),
			Issuer:     getEnv("CHAT_TOKEN_ISSUER", "chaseapp"),
			Audience:   getEnv("CHAT_TOKEN_AUDIENCE", "chat"),
			TokenTTL:   getEnvDuration("CHAT_TOKEN_TTL", 15*time.Minute),
		},
		External: ExternalConfig{
			USGSBaseURL:          getEnv("USGS_BASE_URL", "https://earthquake.usgs.gov"),
			AISHubBaseURL:        getEnv("AISHUB_BASE_URL", "https://data.aishub.net/ws.php"),
			AISHubAPIKey:         getEnv("AISHUB_API_KEY", ""),
			NOAABaseURL:          getEnv("NOAA_BASE_URL", "https://api.weather.gov"),
			LaunchLibraryBaseURL: getEnv("LAUNCH_LIBRARY_BASE_URL", "https://ll.thespacedevs.com/2.2.0"),
			DiscordWebhook:       getEnv("DISCORD_WEBHOOK_URL", ""),
		},
		Observability: ObservabilityConfig{
			ServiceName:  getEnv("OTEL_SERVICE_NAME", "chaseapp-api"),
			OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
			OTLPHeaders:  getEnv("OTEL_EXPORTER_OTLP_HEADERS", ""),
			OTLPInsecure: getEnvBool("OTEL_EXPORTER_OTLP_INSECURE", false),
		},
	}

	return cfg, nil
}

// getEnv returns an environment variable or a default value.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvInt returns an environment variable as int or a default value.
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

// getEnvDuration returns an environment variable as duration or a default value.
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

// getEnvBool returns an environment variable as bool or a default value.
func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		switch val {
		case "1", "true", "TRUE", "True", "yes", "YES", "Yes":
			return true
		case "0", "false", "FALSE", "False", "no", "NO", "No":
			return false
		}
	}
	return defaultVal
}
