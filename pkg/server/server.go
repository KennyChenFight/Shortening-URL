package server

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewHTTPServer(engine *gin.Engine, port string, svc *service.BaseService) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: registerRoutingRule(engine, svc),
	}
}

func registerRoutingRule(engine *gin.Engine, svc *service.BaseService) *gin.Engine {
	v1APIGroup := engine.Group("/api/v1")
	{
		v1APIGroup.POST("/urls", svc.CreateShorteningURL)
		v1APIGroup.POST("/batchUrls", svc.BatchCreateShorteningURL)
		v1APIGroup.DELETE("/urls/:id", svc.DeleteShorteningURL)
	}

	// for redirect
	engine.GET("/:id", svc.GetShorteningURL)
	return engine
}
