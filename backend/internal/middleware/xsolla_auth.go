package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Shyramond/hakathon/backend/pkg/xsolla"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyXsollaUserID = "xsolla_user_id"
	ContextKeyXsollaToken  = "xsolla_token"
	ContextKeyUserInfo     = "xsolla_user_info"
)

// XsollaAuth — middleware для валидации JWT-токена Xsolla Login.
// Извлекает Bearer token → валидирует через JWKS → кладёт user info в context.
func XsollaAuth(xsollaClient *xsolla.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractBearerToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": err.Error(),
			})
			return
		}

		userInfo, err := xsollaClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			slog.Warn("token validation failed",
				"error", err,
				"client_ip", c.ClientIP(),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid or expired token",
			})
			return
		}

		if userInfo.ID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "token does not contain user ID",
			})
			return
		}

		// Кладём в context
		c.Set(ContextKeyXsollaUserID, userInfo.ID)
		c.Set(ContextKeyXsollaToken, token)
		c.Set(ContextKeyUserInfo, userInfo)

		c.Next()
	}
}

// OptionalXsollaAuth — не блокирует запрос при отсутствии токена,
// но если токен есть — валидирует и кладёт в context.
func OptionalXsollaAuth(xsollaClient *xsolla.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractBearerToken(c)
		if err != nil {
			c.Next()
			return
		}

		userInfo, err := xsollaClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Next()
			return
		}

		if userInfo.ID != "" {
			c.Set(ContextKeyXsollaUserID, userInfo.ID)
			c.Set(ContextKeyXsollaToken, token)
			c.Set(ContextKeyUserInfo, userInfo)
		}

		c.Next()
	}
}

func extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", fmt.Errorf("invalid Authorization format, expected: Bearer <token>")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, nil
}
