package service

import (
	"fmt"
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
	"github.com/KennyChenFight/Shortening-URL/pkg/validation"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Config struct {
	FQDN string
}

type BaseService struct {
	config               *Config
	logger               *loglib.Logger
	urlDAO               dao.UrlDAO
	validationTranslator *validation.ValidationTranslator
}

func NewService(config *Config, logger *loglib.Logger, urlDAO dao.UrlDAO, validatorTranslator *validation.ValidationTranslator) *BaseService {
	return &BaseService{config: config, logger: logger, urlDAO: urlDAO, validationTranslator: validatorTranslator}
}

func sendErrorResponse(c *gin.Context, service *BaseService, businessError *business.Error) {
	translated, err := service.validationTranslator.Translate(c.GetHeader("Accept-Language"), businessError.Reason)
	if translated == nil && err == nil {
		c.JSON(businessError.HTTPStatusCode, businessError)
		return
	}

	if err != nil {
		service.logger.Error("validation translate fail", zap.Error(err))
		c.JSON(businessError.HTTPStatusCode, businessError)
		return
	}
	businessError.ValidationErrors = translated
	c.JSON(businessError.HTTPStatusCode, businessError)
}

func sendSuccessResponse(c *gin.Context, statusCode int, response interface{}) {
	c.JSON(statusCode, response)
}

func combineFQDNAndShorteningURLID(fqdn, id string) string {
	return fmt.Sprintf("%s/%s", fqdn, id)
}
