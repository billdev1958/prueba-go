package reports

import (
	"context"
	"net/http"
	"prueba-go/internal/infrastructure/http/v1/reports/models"
	usecases "prueba-go/internal/usecases/report"
	"prueba-go/pkg/types"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	usecase usecases.ReportUsecase
}

func NewReportHandler(usecase usecases.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		usecase: usecase,
	}
}

// GetGlobalEarnings godoc
// @Summary Get total global earnings
// @Description Calculate the sum of all net earnings across all transactions
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {object} models.EarningsResponse
// @Failure 500 {object} common.ErrorResponse
// @Security UserID
// @Router /reports/earnings/global [get]
func (h *ReportHandler) GetGlobalEarnings(c *gin.Context) {
	actor := c.GetHeader("X-User-Id")
	ctx := context.WithValue(c.Request.Context(), "actor", actor)

	amount, err := h.usecase.GetGlobalEarnings(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.EarningsResponse{
		Amount: amount.AmmountToString(),
	})
}

// GetEarningsByMerchant godoc
// @Summary Get earnings for a specific merchant
// @Description Calculate the sum of net earnings for a merchant by their UUID
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "UUID del comercio" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 200 {object} models.EarningsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Security UserID
// @Router /reports/earnings/merchant/{id} [get]
func (h *ReportHandler) GetEarningsByMerchant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "merchant id is required"})
		return
	}

	actor := c.GetHeader("X-User-Id")
	ctx := context.WithValue(c.Request.Context(), "actor", actor)

	amount, err := h.usecase.GetEarningsByMerchant(ctx, types.UID(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.EarningsResponse{
		Amount: amount.AmmountToString(),
	})
}
