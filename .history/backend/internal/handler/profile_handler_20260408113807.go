package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/Shyramond/hakathon/backend/internal/middleware"
    "github.com/Shyramond/hakathon/backend/internal/repository"
    "github.com/Shyramond/hakathon/backend/internal/service"
)

type ProfileHandler struct {
    shopService    *service.ShopService
    tokenService   *service.TokenService
    transactionRepo repository.TransactionRepository
}

func NewProfileHandler(
    shopService *service.ShopService,
    tokenService *service.TokenService,
    transactionRepo repository.TransactionRepository,
) *ProfileHandler {
    return &ProfileHandler{
        shopService:     shopService,
        tokenService:    tokenService,
        transactionRepo: transactionRepo,
    }
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
    user := middleware.MustGetUser(c)
    userID := middleware.MustGetUserID(c)

    inventory, err := h.shopService.GetInventory(c.Request.Context(), userID)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    response := ProfileResponse{
        User:      toUserDTO(user),
        Inventory: toInventoryItemDTOs(inventory),
    }

    respondSuccess(c, http.StatusOK, response)
}

func (h *ProfileHandler) GetTransactions(c *gin.Context) {
    userID := middleware.MustGetUserID(c)

    var pagination PaginationQuery
    if err := c.ShouldBindQuery(&pagination); err != nil {
        respondError(c, http.StatusBadRequest, "invalid_params", "Invalid pagination parameters")
        return
    }
    pagination.Validate()

    transactions, total, err := h.transactionRepo.GetByUserID(
        c.Request.Context(), userID,
        pagination.Limit, pagination.Offset,
    )
    if err != nil {
        handleServiceError(c, err)
        return
    }

    respondSuccessWithMeta(c, http.StatusOK, toTransactionDTOs(transactions), &APIMeta{
        Total:  total,
        Limit:  pagination.Limit,
        Offset: pagination.Offset,
    })
}