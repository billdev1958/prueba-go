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

// Create godoc
// @Summary Create a new transaction
// @Description Record a new transaction for a commerce
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} models.TransactionResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Security UserID
// @Router /transactions [post]
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

// GetByID godoc
// @Summary Get a transaction by ID
// @Description Get detailed information about a transaction using its UUID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "UUID de la transacción" format(uuid) example(f47ac10b-58cc-4372-a567-0e02b2c3d479)
// @Success 200 {object} models.TransactionResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Security UserID
// @Router /transactions/{id} [get]
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

// GetAll godoc
// @Summary List all transactions
// @Description Get a list of all recorded transactions
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {array} models.TransactionResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /transactions [get]
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
