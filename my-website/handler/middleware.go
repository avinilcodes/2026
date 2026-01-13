package handler

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// simple per-key limiter map (keyed by remote IP). Not distributed.
type clientLimiter struct {
	mu      sync.Mutex
	clients map[string]*rate.Limiter
	r       rate.Limit
	b       int
}

// NewClientLimiter creates a client limiter with given rate and burst.
func NewClientLimiter(r rate.Limit, burst int) *clientLimiter {
	cl := &clientLimiter{
		clients: make(map[string]*rate.Limiter),
		r:       r,
		b:       burst,
	}
	go cl.cleanupLoop()
	return cl
}

func (c *clientLimiter) get(ip string) *rate.Limiter {
	c.mu.Lock()
	defer c.mu.Unlock()
	l, ok := c.clients[ip]
	if !ok {
		l = rate.NewLimiter(c.r, c.b)
		c.clients[ip] = l
	}
	return l
}

func (c *clientLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.clients {
			// best-effort: remove idle limiters
			if v.AllowN(time.Now(), 0) {
				delete(c.clients, k)
			}
		}
		c.mu.Unlock()
	}
}

// RateLimitMiddleware returns middleware enforcing per-client rate limits keyed by remote IP.
func RateLimitMiddleware(cl *clientLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			l := cl.get(ip)
			if !l.Allow() {
				WriteError(w, http.StatusTooManyRequests, "Too many requests")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
