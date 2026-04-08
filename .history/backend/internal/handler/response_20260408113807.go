package handler

import (
    "github.com/gin-gonic/gin"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
    Meta    *APIMeta    `json:"meta,omitempty"`
}

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type APIMeta struct {
    Total  int `json:"total,omitempty"`
    Limit  int `json:"limit,omitempty"`
    Offset int `json:"offset,omitempty"`
}

func respondSuccess(c *gin.Context, status int, data interface{}) {
    c.JSON(status, APIResponse{
        Success: true,
        Data:    data,
    })
}

func respondSuccessWithMeta(c *gin.Context, status int, data interface{}, meta *APIMeta) {
    c.JSON(status, APIResponse{
        Success: true,
        Data:    data,
        Meta:    meta,
    })
}

func respondError(c *gin.Context, status int, code, message string) {
    c.JSON(status, APIResponse{
        Success: false,
        Error: &APIError{
            Code:    code,
            Message: message,
        },
    })
}