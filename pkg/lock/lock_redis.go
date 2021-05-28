package lock

import (
	"context"
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"

	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/redislib"
)

func NewRedisLocker(logger *loglib.Logger, client *redislib.GORedisClient) *RedisLocker {
	return &RedisLocker{logger: logger, client: client}
}

type RedisLocker struct {
	logger *loglib.Logger
	client *redislib.GORedisClient
}

func (r *RedisLocker) AcquireLock(name string, lockDuration, waitTime time.Duration) (bool, *business.Error) {
	for _, t := range waitTimeSeries(waitTime) {
		ok, err := r.client.SetNX(context.Background(), name, "", lockDuration).Result()
		if err != nil {
			return false, redisErrorHandle(r.logger, err)
		}

		if ok {
			return true, nil
		} else {
			if t >= 0 {
				time.Sleep(t)
			} else {
				break
			}
		}
	}
	return false, nil
}

func (r *RedisLocker) ReleaseLock(name string) *business.Error {
	_, err := r.client.Del(context.Background(), name).Result()
	if err != nil {
		return redisErrorHandle(r.logger, err)
	}
	return nil
}
