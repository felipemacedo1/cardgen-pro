package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/felipemacedo/cardgen-pro/internal/generator"
	"github.com/felipemacedo/cardgen-pro/internal/iso"
	"github.com/felipemacedo/cardgen-pro/internal/models"
)

// Server represents the HTTP API server for fixtures
type Server struct {
	token      string
	port       int
	rateLimiter *RateLimiter
}

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request from IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	
	// Clean old requests
	requests := rl.requests[ip]
	validRequests := []time.Time{}
	for _, t := range requests {
		if now.Sub(t) < rl.window {
			validRequests = append(validRequests, t)
		}
	}

	// Check limit
	if len(validRequests) >= rl.limit {
		return false
	}

	// Add new request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests

	return true
}

// NewServer creates a new API server
func NewServer(token string, port int) *Server {
	return &Server{
		token:       token,
		port:        port,
		rateLimiter: NewRateLimiter(100, time.Minute), // 100 requests per minute
	}
}

// authMiddleware validates the bearer token
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing Authorization header", http.StatusUnauthorized)
			return
		}

		expectedAuth := "Bearer " + s.token
		if authHeader != expectedAuth {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// rateLimitMiddleware implements rate limiting
func (s *Server) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		
		if !s.rateLimiter.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}

// handleGenerateCards handles GET /v1/cards
func (s *Server) handleGenerateCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	brand := r.URL.Query().Get("brand")
	if brand == "" {
		brand = "visa"
	}

	bin := r.URL.Query().Get("bin")
	
	countStr := r.URL.Query().Get("count")
	count := 10
	if countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil && c > 0 && c <= 100 {
			count = c
		}
	}

	secret := r.URL.Query().Get("secret")

	// Generate cards
	opts := models.GenerateOptions{
		BIN:           bin,
		Brand:         brand,
		Count:         count,
		Secret:        secret,
		IncludeISO:    true,
		IncludeTrack2: true,
	}

	cards := []*models.Card{}
	for i := 0; i < count; i++ {
		card, err := generator.GenerateCard(opts)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate card: %v", err), http.StatusInternalServerError)
			return
		}

		// Add ISO fields
		if opts.IncludeISO {
			isoFields := iso.GenerateISO8583Fields(card, 10000, "986")
			card.ISOFields = isoFields
		}

		cards = append(cards, card)
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cards": cards,
		"count": len(cards),
	})
}

// handleScenarios handles GET /v1/scenarios
func (s *Server) handleScenarios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	scenarios := GetScenarios()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenarios)
}

// handleHealth handles GET /health
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("/health", s.handleHealth)

	// Protected endpoints
	mux.HandleFunc("/v1/cards", s.rateLimitMiddleware(s.authMiddleware(s.handleGenerateCards)))
	mux.HandleFunc("/v1/scenarios", s.rateLimitMiddleware(s.authMiddleware(s.handleScenarios)))

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting API server on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  GET /health")
	log.Printf("  GET /v1/cards (protected)")
	log.Printf("  GET /v1/scenarios (protected)")
	
	return http.ListenAndServe(addr, mux)
}
