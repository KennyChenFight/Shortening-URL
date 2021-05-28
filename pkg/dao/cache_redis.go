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
	random := getRandomOriginalURLTTLSecond()
	expire := hotOriginalURLBaseTTL + random
	_, err := r.client.Set(context.Background(), fmt.Sprintf("%s-%s", prefixHotOriginalURL, name), originalURL, expire).Result()
	if err != nil {
		return redisErrorHandle(r.logger, err)
	}
	return nil
}

func (r *RedisCacheDAO) DeleteOriginalURL(name string) *business.Error {
	err := r.client.Del(context.Background(), fmt.Sprintf("%s-%s", prefixHotOriginalURL, name)).Err()
	if err != nil {
		return redisErrorHandle(r.logger, err)
	}
	return nil
}

func (r *RedisCacheDAO) DeleteMultiOriginalURL(names []string) *business.Error {
	var formatNames []string
	for _, name := range names {
		formatNames = append(formatNames, fmt.Sprintf("%s-%s", prefixHotOriginalURL, name))
	}
	err := r.client.Del(context.Background(), formatNames...).Err()
	if err != nil {
		return redisErrorHandle(r.logger, err)
	}
	return nil
}

func (r *RedisCacheDAO) AddOriginalURLIDInFilters(originalURLID string) *business.Error {
	_, err := r.client.Do(context.Background(), "CF.ADD", originalURLIDsFilterName, originalURLID).Result()
	if err != nil {
		return redisErrorHandle(r.logger, err)
	}
	return nil
}

func (r *RedisCacheDAO) ExistOriginalURLIDInFilters(originalURLID string) (bool, *business.Error) {
	ok, err := r.client.Do(context.Background(), "CF.EXISTS", originalURLIDsFilterName, originalURLID).Bool()
	if err != nil {
		return ok, redisErrorHandle(r.logger, err)
	}
	return ok, nil
}

func (r *RedisCacheDAO) DeleteOriginalURLIDInFilters(originalURLID string) (bool, *business.Error) {
	ok, err := r.client.Do(context.Background(), "CF.DEL", originalURLIDsFilterName, originalURLID).Bool()
	if err != nil {
		return ok, redisErrorHandle(r.logger, err)
	}
	return ok, nil
}

func (r *RedisCacheDAO) DeleteMultiOriginalURLIDInFilters(originalURLIDs []string) (bool, *business.Error) {
	// not good enough :(
	for _, id := range originalURLIDs {
		ok, err := r.client.Do(context.Background(), "CF.DEL", originalURLIDsFilterName, id).Bool()
		if err != nil {
			return ok, redisErrorHandle(r.logger, err)
		}
	}
	return true, nil
}
