package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/middleware"
    "github.com/Shyramond/hakathon/backend/internal/service"
)

type ShopHandler struct {
    shopService *service.ShopService
}

func NewShopHandler(shopService *service.ShopService) *ShopHandler {
    return &ShopHandler{shopService: shopService}
}

func (h *ShopHandler) Purchase(c *gin.Context) {
    userID := middleware.MustGetUserID(c)

    benefitID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        respondError(c, http.StatusBadRequest, "invalid_id", "Invalid benefit ID format")
        return
    }

    result, err := h.shopService.Purchase(c.Request.Context(), userID, benefitID)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    respondSuccess(c, http.StatusCreated, toPurchaseResponse(result))
}

func (h *ShopHandler) GetInventory(c *gin.Context) {
    userID := middleware.MustGetUserID(c)

    items, err := h.shopService.GetInventory(c.Request.Context(), userID)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    respondSuccess(c, http.StatusOK, toInventoryItemDTOs(items))
}

func (h *ShopHandler) EquipItem(c *gin.Context) {
    userID := middleware.MustGetUserID(c)

    itemID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        respondError(c, http.StatusBadRequest, "invalid_id", "Invalid item ID format")
        return
    }

    item, err := h.shopService.EquipItem(c.Request.Context(), userID, itemID)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    respondSuccess(c, http.StatusOK, toInventoryItemDTO(item))
}

func (h *ShopHandler) UnequipItem(c *gin.Context) {
    userID := middleware.MustGetUserID(c)

    itemID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        respondError(c, http.StatusBadRequest, "invalid_id", "Invalid item ID format")
        return
    }

    item, err := h.shopService.UnequipItem(c.Request.Context(), userID, itemID)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    respondSuccess(c, http.StatusOK, toInventoryItemDTO(item))
}