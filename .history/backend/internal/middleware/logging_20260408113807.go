package middleware

import (
    "log/slog"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

func RequestLogging() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Header("X-Request-ID", requestID)
        c.Set("request_id", requestID)

        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()

        attrs := []any{
            "request_id", requestID,
            "method", c.Request.Method,
            "path", path,
            "status", status,
            "latency_ms", latency.Milliseconds(),
            "client_ip", c.ClientIP(),
        }

        if status >= 500 {
            slog.Error("request completed", attrs...)
        } else if status >= 400 {
            slog.Warn("request completed", attrs...)
        } else {
            slog.Info("request completed", attrs...)
        }
    }
}