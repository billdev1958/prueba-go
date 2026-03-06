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
