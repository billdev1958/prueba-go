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
