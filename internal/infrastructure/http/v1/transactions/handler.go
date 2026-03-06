package transactions

import (
	"context"
	"net/http"
	"prueba-go/internal/domain/transaction"
	"prueba-go/internal/infrastructure/http/v1/transactions/models"
	usecases "prueba-go/internal/usecases/transaction"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	usecase usecases.TransactionUsecases
}

func NewTransactionHandler(usecase usecases.TransactionUsecases) *TransactionHandler {
	return &TransactionHandler{
		usecase: usecase,
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := money.NewAmmount(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount format"})
		return
	}

	tr := &transaction.Transaction{
		CommercioID: types.UID(req.CommercioID),
		Amount:      amount,
	}

	actor := c.GetHeader("X-User-Id")
	ctx := context.WithValue(c.Request.Context(), "actor", actor)

	res, err := h.usecase.Create(ctx, tr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toResponse(res))
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	actor := c.GetHeader("X-User-Id")
	ctx := context.WithValue(c.Request.Context(), "actor", actor)

	res, err := h.usecase.GetByID(ctx, types.UID(id))
	if err != nil {
		if err == transaction.ErrTransactionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toResponse(&res))
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	list, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.TransactionResponse, len(list))
	for i, item := range list {
		responses[i] = *toResponse(&item)
	}

	c.JSON(http.StatusOK, responses)
}

func toResponse(t *transaction.Transaction) *models.TransactionResponse {
	return &models.TransactionResponse{
		ID:          string(t.ID),
		CommercioID: string(t.CommercioID),
		Amount:      t.Amount.AmmountToString(),
		AppliedRate: t.AppliedRate.RateToString(),
		Commission:  t.Commission.AmmountToString(),
		NetAmount:   t.NetAmount.AmmountToString(),
		CreatedAt:   t.CreatedAt,
	}
}
