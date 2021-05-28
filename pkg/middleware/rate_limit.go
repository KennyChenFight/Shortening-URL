package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"

	"github.com/gin-gonic/gin"
)

func (b *BaseMiddleware) SlideWindowRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.URL.Path + "-" + c.Request.Method + "-" + c.ClientIP()
		total, err := b.rateLimiter.Incr(context.Background(), key, time.Now().UnixNano())
		if err != nil {
			c.Error(business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", err))
			c.Abort()
			return
		}
		if total == -1 {
			c.Error(business.NewError(business.TooManyRequest, http.StatusTooManyRequests, "reach limit for this endpoint", errors.New("reach limit for this endpoint")))
			c.Abort()
			return
		}
		c.Header("X-RateLimit-Total", strconv.FormatInt(total, 10))
	}
}
