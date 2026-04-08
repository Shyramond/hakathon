package middleware

import (
    "errors"
    "log/slog"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

const (
    ContextKeyUserID = "user_id"
    ContextKeyUser   = "user"
)

// UserContext загружает пользователя из БД по xsolla_user_id (установленному в XsollaAuth)
// и кладёт user + user_id в context.
// Если пользователь не найден в БД — возвращает 401.
func UserContext(userRepo repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        xsollaUserID, exists := c.Get(ContextKeyXsollaUserID)
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error":   "unauthorized",
                "message": "user not authenticated",
            })
            return
        }

        xsollaID, ok := xsollaUserID.(string)
        if !ok || xsollaID == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error":   "unauthorized",
                "message": "invalid user identity",
            })
            return
        }

        user, err := userRepo.GetByXsollaID(c.Request.Context(), xsollaID)
        if err != nil {
            if errors.Is(err, repository.ErrUserNotFound) {
                // Пользователь не зарегистрирован — нужен сначала POST /auth/login
                c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                    "error":   "unauthorized",
                    "message": "user not registered, please login first via POST /api/v1/auth/login",
                })
                return
            }

            slog.Error("failed to load user from DB", "error", err, "xsolla_user_id", xsollaID)
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "error":   "internal_error",
                "message": "failed to load user profile",
            })
            return
        }

        c.Set(ContextKeyUserID, user.ID)
        c.Set(ContextKeyUser, user)

        c.Next()
    }
}

// MustGetUserID извлекает user_id из gin.Context. Паникует если отсутствует.
func MustGetUserID(c *gin.Context) uuid.UUID {
    val, exists := c.Get(ContextKeyUserID)
    if !exists {
        panic("user_id not found in context — UserContext middleware missing")
    }
    id, ok := val.(uuid.UUID)
    if !ok {
        panic("user_id in context has wrong type")
    }
    return id
}

// MustGetUser извлекает *models.User из gin.Context. Паникует если отсутствует.
func MustGetUser(c *gin.Context) *models.User {
    val, exists := c.Get(ContextKeyUser)
    if !exists {
        panic("user not found in context — UserContext middleware missing")
    }
    user, ok := val.(*models.User)
    if !ok {
        panic("user in context has wrong type")
    }
    return user
}

// GetUserIDOptional извлекает user_id если есть (для optional auth).
func GetUserIDOptional(c *gin.Context) *uuid.UUID {
    val, exists := c.Get(ContextKeyUserID)
    if !exists {
        return nil
    }
    id, ok := val.(uuid.UUID)
    if !ok {
        return nil
    }
    return &id
}