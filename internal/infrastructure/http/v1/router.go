package v1

import (
	_ "prueba-go/internal/docs"
	"prueba-go/internal/infrastructure/http/v1/audit"
	"prueba-go/internal/infrastructure/http/v1/comercios"
	"prueba-go/internal/infrastructure/http/v1/middleware"
	"prueba-go/internal/infrastructure/http/v1/reports"
	"prueba-go/internal/infrastructure/http/v1/transactions"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	CommerceHandler    *comercios.ComercioHandler
	ReportHandler      *reports.ReportHandler
	TransactionHandler *transactions.TransactionHandler
	AuditHandler       *audit.AuditHandler
}

func RegisterRoutes(r *gin.Engine, config RouterConfig) {
	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	v1.Use(middleware.Auth())
	{
		comercios.RegisterRoutes(v1, config.CommerceHandler)
		reports.RegisterRoutes(v1, config.ReportHandler)
		transactions.RegisterRoutes(v1, config.TransactionHandler)
		audit.RegisterRoutes(v1, config.AuditHandler)
	}
}
