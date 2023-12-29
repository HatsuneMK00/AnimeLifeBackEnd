package routes

import (
	"AnimeLifeBackEnd/metrics"
	"github.com/gin-gonic/gin"
)

type MetricsRouter struct{}

func (r MetricsRouter) AddMetricsRoutes(rg *gin.RouterGroup) {
	rg.GET("/healthz", metrics.HealthProbeHandler)
}
