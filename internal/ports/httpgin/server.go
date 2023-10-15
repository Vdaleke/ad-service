package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"homework10/internal/app"
)

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	api := router.Group("/api/v1")
	AppRouter(api, a)

	httpServer := http.Server{
		Addr:    port,
		Handler: router,
	}

	return &httpServer
}
