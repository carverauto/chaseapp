package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiterConfig configures the rate limiter.
type RateLimiterConfig struct {
	// RequestsPerSecond is the number of requests allowed per second.
	RequestsPerSecond float64
	// BurstSize is the maximum number of requests allowed in a burst.
	BurstSize int
	// CleanupInterval is how often to clean up stale entries.
	CleanupInterval time.Duration
}

// DefaultRateLimiterConfig returns default rate limiter configuration.
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute,
	}
}

// tokenBucket implements the token bucket algorithm.
type tokenBucket struct {
	tokens     float64
	lastUpdate time.Time
	mu         sync.Mutex
}

// RateLimiter implements a per-client rate limiter using the token bucket algorithm.
type RateLimiter struct {
	config  RateLimiterConfig
	clients map[string]*tokenBucket
	mu      sync.RWMutex
	stopCh  chan struct{}
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		config:  config,
		clients: make(map[string]*tokenBucket),
		stopCh:  make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given key is allowed.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.RLock()
	bucket, exists := rl.clients[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		bucket, exists = rl.clients[key]
		if !exists {
			bucket = &tokenBucket{
				tokens:     float64(rl.config.BurstSize),
				lastUpdate: time.Now(),
			}
			rl.clients[key] = bucket
		}
		rl.mu.Unlock()
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastUpdate).Seconds()
	bucket.lastUpdate = now

	// Add tokens based on time elapsed
	bucket.tokens += elapsed * rl.config.RequestsPerSecond
	if bucket.tokens > float64(rl.config.BurstSize) {
		bucket.tokens = float64(rl.config.BurstSize)
	}

	// Check if request is allowed
	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup periodically removes stale entries.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			cutoff := time.Now().Add(-rl.config.CleanupInterval * 2)
			for key, bucket := range rl.clients {
				bucket.mu.Lock()
				if bucket.lastUpdate.Before(cutoff) {
					delete(rl.clients, key)
				}
				bucket.mu.Unlock()
			}
			rl.mu.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

// Stop stops the rate limiter cleanup goroutine.
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

// RateLimit returns a middleware that rate limits requests.
func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client identifier (prefer user ID, fall back to IP)
			key := getClientKey(r)

			if !limiter.Allow(key) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientKey returns a unique identifier for the client.
func getClientKey(r *http.Request) string {
	// Prefer authenticated user ID
	if userID, ok := r.Context().Value(UserIDKey).(string); ok && userID != "" {
		return "user:" + userID
	}

	// Fall back to IP address
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	return "ip:" + ip
}

// RateLimitByEndpoint returns a middleware that applies different rate limits per endpoint.
type EndpointRateLimits struct {
	Default  RateLimiterConfig
	Endpoints map[string]RateLimiterConfig
}

// RateLimitByEndpoint creates a rate limiter with per-endpoint configuration.
func RateLimitByEndpoint(limits EndpointRateLimits) func(http.Handler) http.Handler {
	defaultLimiter := NewRateLimiter(limits.Default)
	endpointLimiters := make(map[string]*RateLimiter)

	for endpoint, config := range limits.Endpoints {
		endpointLimiters[endpoint] = NewRateLimiter(config)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := defaultLimiter

			// Check for endpoint-specific limiter
			if el, ok := endpointLimiters[r.URL.Path]; ok {
				limiter = el
			}

			key := getClientKey(r)

			if !limiter.Allow(key) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
