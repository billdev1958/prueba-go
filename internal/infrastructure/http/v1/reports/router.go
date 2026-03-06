package reports

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *ReportHandler) {
	reports := r.Group("/reports")
	{
		earnings := reports.Group("/earnings")
		{
			earnings.GET("/global", handler.GetGlobalEarnings)
			earnings.GET("/merchant/:id", handler.GetEarningsByMerchant)
		}
	}
}
