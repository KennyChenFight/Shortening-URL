package dao

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/KennyChenFight/golib/redislib"
	"github.com/go-redis/redismock/v8"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisCacheDAO", func() {
	var redisCacheDAO *RedisCacheDAO
	var logger *loglib.Logger

	BeforeEach(func() {
		logger = loglib.NewNopLogger()
		redisCacheDAO = &RedisCacheDAO{logger, testRedisClient}
	})

	var _ = Describe("GetOriginalURL", func() {
		var (
			expectOriginalURL string
			getErr            *business.Error
		)

		ctx := context.Background()
		name := "testName"
		key := fmt.Sprintf("%s-%s", prefixHotOriginalURL, name)
		actualOriginalURL := "http://example.com"

		JustBeforeEach(func() {
			expectOriginalURL, getErr = redisCacheDAO.GetOriginalURL(name)
		})

		Context("success", func() {
			BeforeEach(func() {
				testRedisClient.Set(ctx, key, actualOriginalURL, -1)
			})

			AfterEach(func() {
				testRedisClient.Del(ctx, key)
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectOriginalURL).To(Equal(actualOriginalURL))
			})
		})

		Context("redis key not exist", func() {
			It("result", func() {
				Expect(getErr).To(Equal(business.NewError(business.NotFound, http.StatusNotFound, "record not found", RedisErrKeyNotExist)))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})

		Context("redis internal error", func() {
			wrapperClient, redisMock := redismock.NewClientMock()
			internalErr := errors.New("internal error")
			BeforeEach(func() {
				redisCacheDAO.client = &redislib.GORedisClient{Client: wrapperClient}
				redisMock.ExpectGet(key).SetErr(internalErr)
			})

			AfterEach(func() {
				redisMock.ClearExpect()
				wrapperClient.Close()
				redisCacheDAO.client = testRedisClient
			})

			It("result", func() {
				Expect(getErr).To(Equal(business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", internalErr)))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})
	})

	var _ = Describe("SetOriginalURL", func() {
		var (
			setErr *business.Error
		)

		ctx := context.Background()
		name := "testName"
		key := fmt.Sprintf("%s-%s", prefixHotOriginalURL, name)
		originalURL := "http://example.com"

		JustBeforeEach(func() {
			setErr = redisCacheDAO.SetOriginalURL(name, originalURL)
		})

		Context("success", func() {
			AfterEach(func() {
				testRedisClient.Del(ctx, key)
			})

			It("result", func() {
				Expect(setErr).To(BeNil())
				Expect(testRedisClient.Get(ctx, key).Val()).To(Equal(originalURL))
			})
		})

		Context("redis internal error", func() {
			wrapperClient, redisMock := redismock.NewClientMock()
			internalErr := errors.New("internal error")
			BeforeEach(func() {
				redisCacheDAO.client = &redislib.GORedisClient{Client: wrapperClient}
				redisMock.ExpectSet(key, originalURL, hotOriginalURLBaseTTL).SetErr(internalErr)
			})

			AfterEach(func() {
				redisMock.ClearExpect()
				wrapperClient.Close()
				redisCacheDAO.client = testRedisClient
			})

			It("result", func() {
				Expect(setErr).To(Equal(business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", internalErr)))
			})
		})
	})
})
