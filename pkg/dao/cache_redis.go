package dao

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/redislib"
)

func NewRedisCacheDAO(logger *loglib.Logger, client *redislib.GORedisClient) *RedisCacheDAO {
	return &RedisCacheDAO{logger, client}
}

type RedisCacheDAO struct {
	logger *loglib.Logger
	client *redislib.GORedisClient
}

func (r *RedisCacheDAO) GetOriginalURL(name string) (string, *business.Error) {
	originalURL, err := r.client.Get(context.Background(), fmt.Sprintf("%s-%s", prefixHotOriginalURL, name)).Result()
	if err != nil {
		return "", redisErrorHandle(r.logger, err)
	}
	return originalURL, nil
}

func (r *RedisCacheDAO) SetOriginalURL(name string, originalURL string) *business.Error {
	// 加入random seconds to prevent 大量緩存同時失效的問題
	random := time.Duration(rand.Int63n(int64(randomOriginalURLTTLNumber))+1) * time.Second
	expire := hotOriginalURLBaseTTL + random
	_, err := r.client.Set(context.Background(), fmt.Sprintf("%s-%s", prefixHotOriginalURL, name), originalURL, expire).Result()
	if err != nil {
		return redisErrorHandle(r.logger, err)
	}
	return nil
}
