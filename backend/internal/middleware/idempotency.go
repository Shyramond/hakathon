package middleware

import (
    "bytes"
    "encoding/json"
    "io"
    "log/slog"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

const IdempotencyKeyHeader = "Idempotency-Key"


func Idempotency(repo repository.IdempotencyRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Только для мутирующих методов
        if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
            c.Next()
            return
        }

        key := c.GetHeader(IdempotencyKeyHeader)
        if key == "" {
            c.Next()
            return
        }

        if _, err := uuid.Parse(key); err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": "Idempotency-Key must be a valid UUID",
            })
            return
        }

        record, err := repo.Get(c.Request.Context(), key)
        if err != nil {
            slog.Error("idempotency repo get error", "error", err)
            c.Next()
            return
        }

        if record != nil {
            c.Header("X-Idempotent-Replayed", "true")
            var body interface{}
            if err := json.Unmarshal(record.ResponseBody, &body); err == nil {
                c.JSON(record.ResponseCode, body)
            } else {
                c.Data(record.ResponseCode, "application/json", record.ResponseBody)
            }
            c.Abort()
            return
        }

        writer := &responseCapture{
            ResponseWriter: c.Writer,
            body:           &bytes.Buffer{},
        }
        c.Writer = writer

        c.Next()

        if c.Writer.Status() >= 200 && c.Writer.Status() < 500 {
            userIDVal, exists := c.Get("user_id")
            var userID uuid.UUID
            if exists {
                userID, _ = userIDVal.(uuid.UUID)
            }

            idempRecord := &repository.IdempotencyRecord{
                Key:          key,
                UserID:       userID,
                ResponseCode: writer.statusCode,
                ResponseBody: writer.body.Bytes(),
            }

            if err := repo.Set(c.Request.Context(), idempRecord); err != nil {
                slog.Error("idempotency repo set error", "error", err)
            }
        }
    }
}

type responseCapture struct {
    gin.ResponseWriter
    body       *bytes.Buffer
    statusCode int
}

func (w *responseCapture) Write(b []byte) (int, error) {
    w.body.Write(b)
    return w.ResponseWriter.Write(b)
}

func (w *responseCapture) WriteHeader(code int) {
    w.statusCode = code
    w.ResponseWriter.WriteHeader(code)
}

var _ io.Writer = (*responseCapture)(nil)