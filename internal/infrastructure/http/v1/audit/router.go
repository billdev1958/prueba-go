package audit

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *AuditHandler) {
	audit := r.Group("/audit")
	{
		audit.GET("", handler.GetAll)
	}
}
