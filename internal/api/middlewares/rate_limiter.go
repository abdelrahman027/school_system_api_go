package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	// Add fields for rate limiting, e.g., a map to track requests per IP
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *RateLimiter {
	r1 := &RateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	go r1.ResetVisitorsCount()
	return r1
}

func (r1 *RateLimiter) ResetVisitorsCount() {
	time.Sleep(r1.resetTime)
	r1.mu.Lock()
	defer r1.mu.Unlock()
	r1.visitors = make(map[string]int)

}

func (r1 *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		r1.mu.Lock()
		r1.visitors[ip]++
		fmt.Println("Visitor IP:", ip, "Request Count:", r1.visitors[ip], "Limit:", r1.limit)
		if r1.visitors[ip] > r1.limit {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			r1.mu.Unlock()
			return
		}

		defer r1.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
