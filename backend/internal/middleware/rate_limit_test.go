package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestRateLimit_AllowsNormal(t *testing.T) {
    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)

    r.Use(RateLimit(100, 10))
    r.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}

func TestRateLimit_BlocksExcess(t *testing.T) {
    _, r := gin.CreateTestContext(httptest.NewRecorder())

    r.Use(RateLimit(1, 1))
    r.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    w1 := httptest.NewRecorder()
    r.ServeHTTP(w1, httptest.NewRequest(http.MethodGet, "/test", nil))
    if w1.Code != http.StatusOK {
        t.Errorf("first request: expected 200, got %d", w1.Code)
    }

    blocked := false
    for i := 0; i < 10; i++ {
        w := httptest.NewRecorder()
        r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))
        if w.Code == http.StatusTooManyRequests {
            blocked = true
            break
        }
    }

    if !blocked {
        t.Error("expected rate limiter to block excess requests")
    }
}