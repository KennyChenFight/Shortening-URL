package service

import (
	"fmt"
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/repository"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/gin-gonic/gin"
)

type Config struct {
	FQDN string
}

type BaseService struct {
	config        *Config
	logger        *loglib.Logger
	urlRepository repository.Repository
}

func NewService(config *Config, logger *loglib.Logger, urlRepository repository.Repository) *BaseService {
	return &BaseService{config: config, logger: logger, urlRepository: urlRepository}
}

func (s *BaseService) HandleMethodNotAllowed(c *gin.Context) {
	s.responseWithError(c, business.NewError(business.MethodNowAllowed, http.StatusMethodNotAllowed, "http method not allowed", nil))
}

func (s *BaseService) HandlePathNotFound(c *gin.Context) {
	s.responseWithError(c, business.NewError(business.PathNotFound, http.StatusNotFound, "http path not found", nil))
}

func (s *BaseService) responseWithError(c *gin.Context, businessError *business.Error) {
	c.Abort()
	c.Error(businessError)
}

func (s *BaseService) responseWithSuccess(c *gin.Context, businessSuccess *business.Success) {
	c.Set("success", businessSuccess)
}

func combineFQDNWithShorteningURLID(fqdn, id string) string {
	return fmt.Sprintf("%s/%s", fqdn, id)
}
