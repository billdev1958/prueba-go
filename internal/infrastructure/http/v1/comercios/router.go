package comercios

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *ComercioHandler) {
	comercios := r.Group("/comercios")
	{
		comercios.POST("", handler.Create)
		comercios.GET("", handler.GetAll)
		comercios.GET("/:id", handler.GetByID)
		comercios.PUT("/:id", handler.Update)
		comercios.DELETE("/:id", handler.Delete)
	}
}
