package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

// SecurityHeaders добавляет базовые заголовки безопасности
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Защита от кликджекинга
        c.Header("X-Frame-Options", "DENY")
        // Защита от MIME-sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        // XSS-фильтр
        c.Header("X-XSS-Protection", "1; mode=block")
        // Referrer policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        // Content Security Policy
        c.Header("Content-Security-Policy", "default-src 'self'")
        // Запрет кеширования приватных данных
        c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
        c.Header("Pragma", "no-cache")

        c.Next()
    }
}

// SQLInjectionGuard — базовая защита от SQL-инъекций в query параметрах
func SQLInjectionGuard() gin.HandlerFunc {
    dangerous := []string{
        "--", ";--", "/*", "*/", "@@", "@",
        "char(", "nchar(", "varchar(", "nvarchar(",
        "alter ", "begin ", "cast(", "create ",
        "cursor ", "declare ", "delete ", "drop ",
        "exec(", "execute(", "fetch ", "insert ",
        "kill ", "select ", "sys.", "sysobjects",
        "syscolumns", "table ", "update ", "union ",
    }

    return func(c *gin.Context) {
        // Проверяем query string
        query := strings.ToLower(c.Request.URL.RawQuery)
        for _, pattern := range dangerous {
            if strings.Contains(query, pattern) {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                    "error": "potentially dangerous input detected",
                })
                return
            }
        }

        // Проверяем path параметры
        for _, param := range c.Params {
            val := strings.ToLower(param.Value)
            for _, pattern := range dangerous {
                if strings.Contains(val, pattern) {
                    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                        "error": "potentially dangerous input detected",
                    })
                    return
                }
            }
        }

        c.Next()
    }
}

// MaxBodySize ограничивает размер тела запроса
func MaxBodySize(maxBytes int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Body != nil {
            c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
        }
        c.Next()
    }
}