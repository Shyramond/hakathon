package middleware

import (
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type ipLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rps      rate.Limit
    burst    int
}

func newIPLimiter(rps float64, burst int) *ipLimiter {
    return &ipLimiter{
        limiters: make(map[string]*rate.Limiter),
        rps:      rate.Limit(rps),
        burst:    burst,
    }
}

func (l *ipLimiter) getLimiter(ip string) *rate.Limiter {
    l.mu.RLock()
    limiter, exists := l.limiters[ip]
    l.mu.RUnlock()

    if exists {
        return limiter
    }

    l.mu.Lock()
    defer l.mu.Unlock()

    if limiter, exists = l.limiters[ip]; exists {
        return limiter
    }

    limiter = rate.NewLimiter(l.rps, l.burst)
    l.limiters[ip] = limiter
    return limiter
}

func RateLimit(rps float64, burst int) gin.HandlerFunc {
    lim := newIPLimiter(rps, burst)

    return func(c *gin.Context) {
        ip := c.ClientIP()
        limiter := lim.getLimiter(ip)

        if !limiter.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded, try again later",
            })
            return
        }

        c.Next()
    }
}