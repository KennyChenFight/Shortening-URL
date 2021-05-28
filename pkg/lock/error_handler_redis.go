package lock

import (
	"errors"
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	RedisErrKeyNotExist = errors.New("redis key not exist")
)

func redisErrorHandle(logger *loglib.Logger, err error) *business.Error {
	switch {
	case err == redis.Nil:
		return business.NewError(business.NotFound, http.StatusNotFound, "record not found", RedisErrKeyNotExist)
	default:
		logger.Error("redis internal error", zap.Error(err))
		return business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", err)
	}
}
