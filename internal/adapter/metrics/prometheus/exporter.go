package prometheus

import (
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GerHandler() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.Use(logger.SetLogger(logger.WithUTC(true)))

	mux.Any("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	return mux
}
