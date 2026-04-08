package middleware

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)

    r.Use(SecurityHeaders())
    r.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    r.ServeHTTP(w, req)

    expectedHeaders := map[string]string{
        "X-Content-Type-Options": "nosniff",
        "X-Frame-Options":       "DENY",
        "X-XSS-Protection":      "1; mode=block",
    }

    for key, val := range expectedHeaders {
        actual := w.Header().Get(key)
        if actual != val {
            t.Errorf("header %s: expected '%s', got '%s'", key, val, actual)
        }
    }
}

func TestSQLInjectionGuard(t *testing.T) {
    tests := []struct {
        name       string
        path       string
        blocked    bool
    }{
        {"clean path", "/api/v1/benefits", false},
        {"UUID path", "/api/v1/benefits/550e8400-e29b-41d4-a716-446655440000", false},
        {"SQL union", "/api/v1/benefits?type=UNION+SELECT", true},
        {"SQL drop", "/api/v1/benefits?name=DROP+TABLE", true},
        {"SQL semicolon", "/api/v1/benefits?name=;DELETE", true},
        {"double dash", "/api/v1/benefits?name=admin--", true},
        {"normal query", "/api/v1/benefits?type=skin&min_price=10", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            _, r := gin.CreateTestContext(w)

            r.Use(SQLInjectionGuard())
            r.GET("/api/v1/benefits", func(c *gin.Context) {
                c.JSON(200, gin.H{"ok": true})
            })

            req := httptest.NewRequest(http.MethodGet, tt.path, nil)
            r.ServeHTTP(w, req)

            if tt.blocked && w.Code != http.StatusForbidden {
                t.Errorf("expected 403 (blocked), got %d", w.Code)
            }
            if !tt.blocked && w.Code != http.StatusOK {
                t.Errorf("expected 200 (allowed), got %d", w.Code)
            }
        })
    }
}

func TestMaxBodySize(t *testing.T) {
    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)

    r.Use(MaxBodySize(100)) // 100 байт
    r.POST("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    // Большой body
    bigBody := strings.Repeat("a", 200)
    req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(bigBody))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)

    if w.Code != http.StatusRequestEntityTooLarge {
        t.Errorf("expected 413, got %d", w.Code)
    }
}

func TestMaxBodySize_SmallBody(t *testing.T) {
    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)

    r.Use(MaxBodySize(1024))
    r.POST("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"key":"value"}`))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}