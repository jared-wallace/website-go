package middleware

import (
	"net"
	"sync"
	"time"
)

// entry holds the fixed-window counter state for a single IP.
type entry struct {
	count     int
	windowEnd time.Time
}

// RateLimiter implements a fixed-window per-IP rate limiter.
// It is safe for concurrent use.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*entry
	limit   int
	window  time.Duration
}

// NewRateLimiter creates a RateLimiter that allows at most limit requests
// per window per IP address.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		entries: make(map[string]*entry),
		limit:   limit,
		window:  window,
	}
}

// Allow reports whether a request from ip should be permitted. ip may be a
// bare address ("1.2.3.4") or an address:port pair ("1.2.3.4:1234") — the
// port, if present, is stripped before counting.
func (rl *RateLimiter) Allow(ip string) bool {
	// Strip port if present; ignore errors (bare IP passes through unchanged).
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	e, ok := rl.entries[ip]
	if !ok || now.After(e.windowEnd) {
		rl.entries[ip] = &entry{count: 1, windowEnd: now.Add(rl.window)}
		return true
	}
	if e.count >= rl.limit {
		return false
	}
	e.count++
	return true
}
