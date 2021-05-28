package dao

import (
	"context"
	"fmt"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/redislib"
)

func NewRedisCacheDAO(logger *loglib.Logger, client *redislib.GORedisClient) *RedisCacheDAO {
	return &RedisCacheDAO{logger, client}
}

type RedisCacheDAO struct {
	logger *loglib.Logger
	*redislib.GORedisClient
}

func (r *RedisCacheDAO) GetOriginalURL(name string) (string, *business.Error) {
	originalURL, err := r.GORedisClient.Get(context.Background(), fmt.Sprintf("%s-%s", prefixHotOriginalURL, name)).Result()
	if err != nil {
		return originalURL, redisErrorHandle(r.logger, err)
	}
	return originalURL, nil
}

func (r *RedisCacheDAO) SetOriginalURL(name string, originalURL string) *business.Error {
	_, err := r.GORedisClient.Set(context.Background(), fmt.Sprintf("%s-%s", prefixHotOriginalURL, name), originalURL, hotOriginalURLTTL).Result()
	if err != nil {
		return redisErrorHandle(r.logger, err)
	}
	return nil
}
