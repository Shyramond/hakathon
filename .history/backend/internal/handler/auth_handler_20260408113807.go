package handler

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/Shyramond/hakathon/backend/internal/service"
)

// Ре-экспорт ошибок для тестов и handleServiceError
var (
    ErrInsufficientBalance = service.ErrInsufficientBalance
    ErrBenefitNotFound     = service.ErrBenefitNotFound
    ErrBenefitNotActive    = service.ErrBenefitNotActive
    ErrAlreadyOwned        = service.ErrAlreadyOwned
    ErrItemNotFound        = service.ErrItemNotFound
    ErrNotItemOwner        = service.ErrNotItemOwner
    ErrAlreadyEquipped     = service.ErrAlreadyEquipped
    ErrUserNotFound        = service.ErrUserNotFound
    ErrLoginCooldown       = service.ErrLoginCooldown
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// Login обрабатывает POST /api/v1/auth/login
// Принимает Xsolla JWT токен, создаёт/находит пользователя, начисляет награду за вход.
//
// Request body: { "token": "<xsolla_jwt>" }
// Или: Authorization: Bearer <xsolla_jwt> (токен берётся из хедера)
func (h *AuthHandler) Login(c *gin.Context) {
    var token string

    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err == nil && req.Token != "" {
        token = req.Token
    }

    if token == "" {
        authHeader := c.GetHeader("Authorization")
        if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
            token = authHeader[7:]
        }
    }

    if token == "" {
        respondError(c, http.StatusBadRequest, "missing_token",
            "Token is required. Send as JSON body {\"token\": \"...\"} or Authorization header")
        return
    }

    result, err := h.authService.Login(c.Request.Context(), token)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    status := http.StatusOK
    if result.IsNewUser {
        status = http.StatusCreated
    }

    respondSuccess(c, status, toLoginResponse(result))
}

func handleServiceError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrLoginCooldown):
        respondError(c, http.StatusTooManyRequests, "cooldown", err.Error())

    case errors.Is(err, service.ErrInsufficientBalance):
        respondError(c, http.StatusPaymentRequired, "insufficient_balance", err.Error())

    case errors.Is(err, service.ErrBenefitNotFound):
        respondError(c, http.StatusNotFound, "not_found", err.Error())

    case errors.Is(err, service.ErrBenefitNotActive):
        respondError(c, http.StatusGone, "not_active", err.Error())

    case errors.Is(err, service.ErrAlreadyOwned):
        respondError(c, http.StatusConflict, "already_owned", err.Error())

    case errors.Is(err, service.ErrItemNotFound):
        respondError(c, http.StatusNotFound, "not_found", err.Error())

    case errors.Is(err, service.ErrNotItemOwner):
        respondError(c, http.StatusForbidden, "forbidden", err.Error())

    case errors.Is(err, service.ErrAlreadyEquipped):
        respondError(c, http.StatusConflict, "already_equipped", err.Error())

    case errors.Is(err, service.ErrUserNotFound):
        respondError(c, http.StatusNotFound, "user_not_found", err.Error())

    default:
        respondError(c, http.StatusInternalServerError, "internal_error", "Something went wrong")
    }
}