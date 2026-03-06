package transactions

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *TransactionHandler) {
	transactions := r.Group("/transactions")
	{
		transactions.POST("", handler.Create)
		transactions.GET("", handler.GetAll)
		transactions.GET("/:id", handler.GetByID)
	}
}
