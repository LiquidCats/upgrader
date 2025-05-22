package prometheus

import (
	"github.com/LiquidCats/graceful"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GerHandler() graceful.Runner {
	mux := gin.New()
	gin.SetMode(gin.ReleaseMode)

	mux.Any("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	return graceful.ServerRunner(mux, graceful.HttpConfig{Port: "9090"})
}
