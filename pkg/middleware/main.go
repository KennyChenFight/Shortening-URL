package middleware

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/validation"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/ratelimitlib"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewMiddleware(logger *loglib.Logger, validationTranslator validation.Translator, rateLimiter ratelimitlib.RateLimiter) *BaseMiddleware {
	return &BaseMiddleware{logger: logger, validationTranslator: validationTranslator, rateLimiter: rateLimiter}
}

type BaseMiddleware struct {
	logger               *loglib.Logger
	validationTranslator validation.Translator

	rateLimiter ratelimitlib.RateLimiter
}

func (b *BaseMiddleware) sendErrorResponse(c *gin.Context, businessError *business.Error) {
	translated, err := b.validationTranslator.Translate(c.GetHeader("Accept-Language"), businessError.Reason)
	if translated == nil && err == nil {
		c.JSON(businessError.HTTPStatusCode, businessError)
		return
	}

	if err != nil {
		b.logger.Error("fail to translate validation message", zap.Error(err))
		c.JSON(businessError.HTTPStatusCode, businessError)
	} else {
		businessError.ValidationErrors = translated
		c.JSON(businessError.HTTPStatusCode, businessError)
	}
}

func (b *BaseMiddleware) sendSuccessResponse(c *gin.Context, success *business.Success) {
	c.JSON(success.HTTPStatusCode, success.Response)
}

func (b *BaseMiddleware) sendRedirect(c *gin.Context, success *business.Success) {
	c.Redirect(success.HTTPStatusCode, success.Response.(string))
}
