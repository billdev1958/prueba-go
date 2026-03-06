package audit

import (
	"net/http"
	"prueba-go/internal/infrastructure/http/v1/audit/models"
	usecases "prueba-go/internal/usecases/audit"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	usecase usecases.AuditUsecase
}

func NewAuditHandler(usecase usecases.AuditUsecase) *AuditHandler {
	return &AuditHandler{
		usecase: usecase,
	}
}

// GetAll godoc
// @Summary List all audit logs
// @Description Get a full history of all actions recorded in the system
// @Tags audit
// @Accept json
// @Produce json
// @Success 200 {array} models.AuditResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /audit [get]
func (h *AuditHandler) GetAll(c *gin.Context) {
	logs, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.AuditResponse, len(logs))
	for i, l := range logs {
		responses[i] = models.AuditResponse{
			LogID:      l.LogID,
			Action:     l.Action,
			Actor:      l.Actor,
			ResourceID: l.ResourceID,
			Timestamp:  l.Timestamp,
		}
	}

	c.JSON(http.StatusOK, responses)
}
