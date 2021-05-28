package lock

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/KennyChenFight/golib/redislib"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisLocker", func() {
	var redisLocker *RedisLocker
	var logger *loglib.Logger

	BeforeEach(func() {
		logger = loglib.NewNopLogger()
		redisLocker = NewRedisLocker(logger, testRedisClient)
	})

	var _ = Describe("AcquireLock", func() {
		var (
			expectGetLock bool
			getErr        *business.Error
		)

		ctx := context.Background()
		name := "resourceName"
		lockDuration := time.Second
		waitTime := time.Second

		JustBeforeEach(func() {
			expectGetLock, getErr = redisLocker.AcquireLock(name, lockDuration, waitTime)
		})

		Context("success get lock", func() {
			AfterEach(func() {
				Expect(testRedisClient.Del(ctx, name).Err()).NotTo(HaveOccurred())
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectGetLock).To(Equal(true))
			})
		})

		Context("fail with redis internal problem", func() {
			var businessSetErr *business.Error
			var wrapper *redis.Client
			var clientMock redismock.ClientMock
			BeforeEach(func() {
				wrapper, clientMock = redismock.NewClientMock()
				redisLocker.client = &redislib.GORedisClient{Client: wrapper}
				setErr := errors.New("internal")
				businessSetErr = business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", setErr)
				clientMock.ExpectSetNX(name, "", lockDuration).SetErr(setErr)
			})
			AfterEach(func() {
				wrapper.Close()
				redisLocker.client = testRedisClient
			})

			It("result", func() {
				Expect(getErr).To(Equal(businessSetErr))
				Expect(expectGetLock).To(Equal(false))
			})
		})

		Context("fail with others hold the lock", func() {
			BeforeEach(func() {
				testRedisClient.Set(context.Background(), name, "", -1)
			})

			AfterEach(func() {
				Expect(testRedisClient.Del(ctx, name).Err()).NotTo(HaveOccurred())
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectGetLock).To(Equal(false))
			})
		})

	})

	var _ = Describe("ReleaseLock", func() {
		var (
			releaseErr *business.Error
		)

		ctx := context.Background()
		name := "resourceName"

		JustBeforeEach(func() {
			releaseErr = redisLocker.ReleaseLock(name)
		})

		Context("success", func() {
			BeforeEach(func() {
				Expect(testRedisClient.Set(ctx, name, "", -1).Err()).NotTo(HaveOccurred())
			})

			It("result", func() {
				Expect(releaseErr).To(BeNil())
				Expect(testRedisClient.Get(ctx, name).Err()).To(Equal(redis.Nil))
			})
		})

		Context("fail with redis internal problem", func() {
			var businessSetErr *business.Error
			var wrapper *redis.Client
			var clientMock redismock.ClientMock
			BeforeEach(func() {
				wrapper, clientMock = redismock.NewClientMock()
				redisLocker.client = &redislib.GORedisClient{Client: wrapper}
				setErr := errors.New("internal")
				businessSetErr = business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", setErr)
				clientMock.ExpectDel(name).SetErr(setErr)
			})
			AfterEach(func() {
				wrapper.Close()
				redisLocker.client = testRedisClient
			})

			It("result", func() {
				Expect(releaseErr).To(Equal(businessSetErr))
			})
		})
	})
})
