package httpgin

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"

	"homework10/internal/app"
)

func CustomMW(c *gin.Context) {
	t := time.Now()

	c.Next()

	latency := time.Since(t)
	status := c.Writer.Status()

	log.Println("latency", latency, "method", c.Request.Method, "path", c.Request.URL.Path, "status", status)
}

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.Use(CustomMW)

	r.POST("/ads", createAd(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a))
	r.PUT("/ads/:ad_id", updateAd(a))
	r.GET("/ads", listAds(a))
	r.GET("/ads/:ad_id", getAd(a))
	r.GET("/ads/search/:pattern", searchAds(a))
	r.DELETE("/ads/:ad_id", deleteAd(a))

	r.POST("/users", createUser(a))
	r.PUT("/users/:user_id", updateUser(a))
	r.GET("/users/:user_id", getUser(a))
	r.DELETE("/users/:user_id", deleteUser(a))
}
