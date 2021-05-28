package server

import (
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/middleware"

	"github.com/KennyChenFight/Shortening-URL/pkg/service"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer(engine *gin.Engine, port string, mwe *middleware.BaseMiddleware, svc *service.BaseService) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: registerRoutingRule(engine, mwe, svc),
	}
}

func registerRoutingRule(engine *gin.Engine, mwe *middleware.BaseMiddleware, svc *service.BaseService) *gin.Engine {
	engine.Use(mwe.GlobalErrorHandle())
	engine.NoMethod(svc.HandleMethodNotAllowed)
	engine.NoRoute(svc.HandlePathNotFound)
	v1APIGroup := engine.Group("/api/v1")
	{
		v1APIGroup.POST("/urls", svc.CreateShorteningURL)
		v1APIGroup.DELETE("/urls/:id", svc.DeleteShorteningURL)
		// for local test, need to remove in production
		v1APIGroup.POST("/_internal/keys", svc.BatchCreateKeys)
	}

	// for redirect
	engine.GET("/:id", mwe.SlideWindowRateLimiter(), svc.GetOriginalURL)
	return engine
}
