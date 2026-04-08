package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/service"
)

type BenefitHandler struct {
    benefitService *service.BenefitService
}

func NewBenefitHandler(benefitService *service.BenefitService) *BenefitHandler {
    return &BenefitHandler{benefitService: benefitService}
}

func (h *BenefitHandler) ListBenefits(c *gin.Context) {
    var filter models.BenefitFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        respondError(c, http.StatusBadRequest, "invalid_params", "Invalid query parameters")
        return
    }

    benefits, err := h.benefitService.GetAll(c.Request.Context(), filter)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    respondSuccess(c, http.StatusOK, toBenefitDTOs(benefits))
}

func (h *BenefitHandler) GetBenefit(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        respondError(c, http.StatusBadRequest, "invalid_id", "Invalid benefit ID format")
        return
    }

    benefit, err := h.benefitService.GetByID(c.Request.Context(), id)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    respondSuccess(c, http.StatusOK, toBenefitDTO(benefit))
}