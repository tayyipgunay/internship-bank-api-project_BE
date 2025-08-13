package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register mounts the Prometheus metrics endpoint at /metrics
func Register(router *gin.Engine) {
	println("📊 Prometheus metrics endpoint kaydediliyor: /metrics")
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	println("✅ Prometheus metrics endpoint kaydedildi")
}
