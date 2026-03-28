package middleware_test

import (
	"testing"
	"time"

	"github.com/jared-wallace/website-go/internal/middleware"
)

func TestRateLimiter_AllowsUpToLimit(t *testing.T) {
	rl := middleware.NewRateLimiter(5, time.Minute)

	for i := 0; i < 5; i++ {
		if !rl.Allow("1.2.3.4") {
			t.Errorf("call %d: expected Allow to return true, got false", i+1)
		}
	}
}

func TestRateLimiter_BlocksAfterLimit(t *testing.T) {
	rl := middleware.NewRateLimiter(5, time.Minute)

	for i := 0; i < 5; i++ {
		rl.Allow("1.2.3.4")
	}

	if rl.Allow("1.2.3.4") {
		t.Error("6th call: expected Allow to return false, got true")
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	// Use a very short window so we can actually expire it in a test.
	rl := middleware.NewRateLimiter(2, 50*time.Millisecond)

	rl.Allow("5.6.7.8")
	rl.Allow("5.6.7.8")

	if rl.Allow("5.6.7.8") {
		t.Error("before window: expected false after hitting limit, got true")
	}

	// Wait for the window to expire.
	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("5.6.7.8") {
		t.Error("after window reset: expected true, got false")
	}
}

func TestRateLimiter_IsolatesIPs(t *testing.T) {
	rl := middleware.NewRateLimiter(1, time.Minute)

	if !rl.Allow("10.0.0.1") {
		t.Error("first call for 10.0.0.1: expected true")
	}
	if rl.Allow("10.0.0.1") {
		t.Error("second call for 10.0.0.1: expected false (limit reached)")
	}

	// Different IP should be unaffected.
	if !rl.Allow("10.0.0.2") {
		t.Error("first call for 10.0.0.2: expected true")
	}
}

func TestRateLimiter_StripPort(t *testing.T) {
	rl := middleware.NewRateLimiter(1, time.Minute)

	// IP with port should be treated the same as bare IP.
	if !rl.Allow("192.168.1.1:54321") {
		t.Error("first call with port: expected true")
	}
	if rl.Allow("192.168.1.1:12345") {
		t.Error("second call with different port: expected false (same underlying IP)")
	}
}
