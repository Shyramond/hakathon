package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/Shyramond/hakathon/backend/pkg"
)

func init() {
    gin.SetMode(gin.TestMode)
}

func newDevXsollaClient() *xsolla.Client {
    return xsolla.NewClient(xsolla.Config{
        ProjectID:      "XXXXXX",
        OAuth2ClientID: "XXXXXX",
    })
}

func TestXsollaAuth_NoHeader(t *testing.T) {
    client := newDevXsollaClient()

    w := httptest.NewRecorder()
    c, r := gin.CreateTestContext(w)

    r.Use(XsollaAuth(client))
    r.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
    r.ServeHTTP(w, c.Request)

    if w.Code != http.StatusUnauthorized {
        t.Errorf("expected 401, got %d", w.Code)
    }
}

func TestXsollaAuth_InvalidFormat(t *testing.T) {
    client := newDevXsollaClient()

    tests := []struct {
        name   string
        header string
    }{
        {"no bearer prefix", "Token abc123"},
        {"only bearer", "Bearer "},
        {"empty", ""},
        {"basic auth", "Basic dXNlcjpwYXNz"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            _, r := gin.CreateTestContext(w)

            r.Use(XsollaAuth(client))
            r.GET("/test", func(c *gin.Context) {
                c.JSON(200, gin.H{"ok": true})
            })

            req := httptest.NewRequest(http.MethodGet, "/test", nil)
            if tt.header != "" {
                req.Header.Set("Authorization", tt.header)
            }
            r.ServeHTTP(w, req)

            if w.Code != http.StatusUnauthorized {
                t.Errorf("expected 401, got %d", w.Code)
            }
        })
    }
}

func TestXsollaAuth_ValidDevToken(t *testing.T) {
    client := newDevXsollaClient()

    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)

    var capturedXsollaID string

    r.Use(XsollaAuth(client))
    r.GET("/test", func(c *gin.Context) {
        val, exists := c.Get(ContextKeyXsollaUserID)
        if exists {
            capturedXsollaID = val.(string)
        }
        c.JSON(200, gin.H{"ok": true})
    })

    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    req.Header.Set("Authorization", "Bearer test_user_42")
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }

    if capturedXsollaID != "test_user_42" {
        t.Errorf("expected xsolla user ID 'test_user_42', got '%s'", capturedXsollaID)
    }
}

func TestOptionalXsollaAuth_NoHeader(t *testing.T) {
    client := newDevXsollaClient()

    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)

    var hasUser bool

    r.Use(OptionalXsollaAuth(client))
    r.GET("/test", func(c *gin.Context) {
        _, hasUser = c.Get(ContextKeyXsollaUserID)
        c.JSON(200, gin.H{"ok": true})
    })

    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200 (optional), got %d", w.Code)
    }

    if hasUser {
        t.Error("expected no user in context when no header")
    }
}

func TestOptionalXsollaAuth_WithHeader(t *testing.T) {
    client := newDevXsollaClient()

    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)

    var capturedID string

    r.Use(OptionalXsollaAuth(client))
    r.GET("/test", func(c *gin.Context) {
        val, _ := c.Get(ContextKeyXsollaUserID)
        if val != nil {
            capturedID = val.(string)
        }
        c.JSON(200, gin.H{"ok": true})
    })

    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    req.Header.Set("Authorization", "Bearer opt_user")
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }

    if capturedID != "opt_user" {
        t.Errorf("expected 'opt_user', got '%s'", capturedID)
    }
}

func TestExtractBearerToken_TableDriven(t *testing.T) {
    tests := []struct {
        name      string
        header    string
        wantToken string
        wantErr   bool
    }{
        {"valid", "Bearer abc123", "abc123", false},
        {"case insensitive", "bearer XYZ", "XYZ", false},
        {"BEARER", "BEARER token", "token", false},
        {"missing", "", "", true},
        {"no bearer", "Token abc", "", true},
        {"empty token", "Bearer ", "", true},
        {"bearer only", "Bearer", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request = httptest.NewRequest("GET", "/", nil)
            if tt.header != "" {
                c.Request.Header.Set("Authorization", tt.header)
            }

            token, err := extractBearerToken(c)

            if tt.wantErr {
                if err == nil {
                    t.Error("expected error")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if token != tt.wantToken {
                    t.Errorf("expected '%s', got '%s'", tt.wantToken, token)
                }
            }
        })
    }
}