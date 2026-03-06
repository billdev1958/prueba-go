package comercios

import (
	"context"
	"net/http"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/internal/infrastructure/http/v1/comercios/models"
	usecases "prueba-go/internal/usecases/commerce"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"

	"github.com/gin-gonic/gin"
)

type ComercioHandler struct {
	usecase usecases.ComercioUsecase
}

func NewComercioHandler(usecase usecases.ComercioUsecase) *ComercioHandler {
	return &ComercioHandler{
		usecase: usecase,
	}
}

// Create godoc
// @Summary Create a new commerce
// @Description Create a new commerce with the given name and commission rate
// @Tags comercios
// @Accept json
// @Produce json
// @Param comercio body models.CreateComercioRequest true "Commerce to create"
// @Success 201 {object} models.ComercioResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Security UserID
// @Router /comercios [post]
func (h *ComercioHandler) Create(c *gin.Context) {
	var req models.CreateComercioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rate, err := money.NewRate(req.ComissionRate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comission rate format"})
		return
	}

	commerceEntity := &comercio.Comercio{
		Name:          req.Name,
		ComissionRate: rate,
	}

	actor := c.GetHeader("X-User-Id")
	ctx := context.WithValue(c.Request.Context(), "actor", actor)

	created, err := h.usecase.Create(ctx, commerceEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toResponse(created))
}

// GetByID godoc
// @Summary Get a commerce by ID
// @Description Get detailed information about a commerce using its UUID
// @Tags comercios
// @Accept json
// @Produce json
// @Param id path string true "UUID del comercio" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 200 {object} models.ComercioResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /comercios/{id} [get]
func (h *ComercioHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	res, err := h.usecase.GetByID(c.Request.Context(), types.UID(id))
	if err != nil {
		if err == comercio.ErrComercioNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toResponse(&res))
}

// GetAll godoc
// @Summary List all comercios
// @Description Get a list of all registered comercios
// @Tags comercios
// @Accept json
// @Produce json
// @Success 200 {array} models.ComercioResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /comercios [get]
func (h *ComercioHandler) GetAll(c *gin.Context) {
	list, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.ComercioResponse, len(list))
	for i, item := range list {
		responses[i] = *toResponse(&item)
	}

	c.JSON(http.StatusOK, responses)
}

// Update godoc
// @Summary Update an existing commerce
// @Description Update the name or commission rate of a commerce
// @Tags comercios
// @Accept json
// @Produce json
// @Param id path string true "UUID del comercio" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Param comercio body models.UpdateComercioRequest true "Updated commerce data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Security UserID
// @Router /comercios/{id} [put]
func (h *ComercioHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req models.UpdateComercioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commerceEntity := &comercio.Comercio{
		ID: types.UID(id),
	}

	if req.Name != "" {
		commerceEntity.Name = req.Name
	}

	if req.ComissionRate != "" {
		rate, err := money.NewRate(req.ComissionRate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comission rate format"})
			return
		}
		commerceEntity.ComissionRate = rate
	}

	actor := c.GetHeader("X-User-Id")
	ctx := context.WithValue(c.Request.Context(), "actor", actor)

	err := h.usecase.Update(ctx, commerceEntity)
	if err != nil {
		if err == comercio.ErrComercioNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete godoc
// @Summary Delete a commerce
// @Description Remove a commerce by its UUID
// @Tags comercios
// @Accept json
// @Produce json
// @Param id path string true "UUID del comercio" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 204 "No Content"
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Security UserID
// @Router /comercios/{id} [delete]
func (h *ComercioHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	actor := c.GetHeader("X-User-Id")
	ctx := context.WithValue(c.Request.Context(), "actor", actor)

	err := h.usecase.Delete(ctx, types.UID(id))
	if err != nil {
		if err == comercio.ErrComercioNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func toResponse(c *comercio.Comercio) *models.ComercioResponse {
	return &models.ComercioResponse{
		ID:            string(c.ID),
		Name:          c.Name,
		ComissionRate: c.ComissionRate.RateToString(),
		CreatedAt:     c.CreatedAt,
	}
}
