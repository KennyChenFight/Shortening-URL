package dao

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prashantv/gostub"

	"github.com/go-redis/redis/v8"

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
			stub := gostub.New()
			actualRandomDuration := int64(1)
			BeforeEach(func() {
				stub.Stub(&randomFunc, func(_ int64) int64 {
					return actualRandomDuration
				})
				redisCacheDAO.client = &redislib.GORedisClient{Client: wrapperClient}
				redisMock.ExpectSet(key, originalURL, hotOriginalURLBaseTTL+(time.Duration(actualRandomDuration+1)*time.Second)).SetErr(internalErr)
			})

			AfterEach(func() {
				stub.Reset()
				redisMock.ClearExpect()
				wrapperClient.Close()
				redisCacheDAO.client = testRedisClient
			})

			It("result", func() {
				Expect(setErr).To(Equal(business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", internalErr)))
			})
		})
	})

	var _ = Describe("DeleteOriginalURL", func() {
		var (
			deleteErr *business.Error
		)

		ctx := context.Background()
		name := "testName"
		key := fmt.Sprintf("%s-%s", prefixHotOriginalURL, name)
		originalURL := "http://example.com"

		JustBeforeEach(func() {
			deleteErr = redisCacheDAO.DeleteOriginalURL(name)
		})

		Context("success", func() {
			BeforeEach(func() {
				testRedisClient.Set(ctx, key, originalURL, -1)
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
				Expect(testRedisClient.Get(ctx, key).Err()).To(Equal(redis.Nil))
			})
		})

		Context("redis internal error", func() {
			wrapperClient, redisMock := redismock.NewClientMock()
			internalErr := errors.New("internal error")
			BeforeEach(func() {
				redisCacheDAO.client = &redislib.GORedisClient{Client: wrapperClient}
				redisMock.ExpectDel(key).SetErr(internalErr)
			})

			AfterEach(func() {
				redisMock.ClearExpect()
				wrapperClient.Close()
				redisCacheDAO.client = testRedisClient
			})

			It("result", func() {
				Expect(deleteErr).To(Equal(business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", internalErr)))
			})
		})
	})

	var _ = Describe("DeleteMultiOriginalURL", func() {
		var (
			deleteErr *business.Error
		)

		ctx := context.Background()
		names := []string{"testName1", "testName2"}
		var keys []string
		for _, name := range names {
			keys = append(keys, fmt.Sprintf("%s-%s", prefixHotOriginalURL, name))
		}
		originalURL := "http://example.com"

		JustBeforeEach(func() {
			deleteErr = redisCacheDAO.DeleteMultiOriginalURL(names)
		})

		Context("success", func() {
			BeforeEach(func() {
				for _, key := range keys {
					testRedisClient.Set(ctx, key, originalURL, -1)
				}
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
				for _, key := range keys {
					Expect(testRedisClient.Get(ctx, key).Err()).To(Equal(redis.Nil))
				}
			})
		})

		Context("redis internal error", func() {
			wrapperClient, redisMock := redismock.NewClientMock()
			internalErr := errors.New("internal error")
			BeforeEach(func() {
				redisCacheDAO.client = &redislib.GORedisClient{Client: wrapperClient}
				redisMock.ExpectDel(keys...).SetErr(internalErr)
			})

			AfterEach(func() {
				redisMock.ClearExpect()
				wrapperClient.Close()
				redisCacheDAO.client = testRedisClient
			})

			It("result", func() {
				Expect(deleteErr).To(Equal(business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", internalErr)))
			})
		})
	})

	var _ = Describe("AddOriginalURLIDInFilters", func() {
		var (
			addErr *business.Error
		)

		ctx := context.Background()
		id := "random"

		JustBeforeEach(func() {
			addErr = redisCacheDAO.AddOriginalURLIDInFilters(id)
		})

		Context("success", func() {
			AfterEach(func() {
				_, err := testRedisClient.Del(ctx, originalURLIDsFilterName).Result()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(addErr).To(BeNil())
				ok, err := testRedisClient.Do(ctx, "CF.EXISTS", originalURLIDsFilterName, id).Bool()
				Expect(err).To(BeNil())
				Expect(ok).To(Equal(true))
			})
		})
	})

	var _ = Describe("ExistOriginalURLIDInFilters", func() {
		var (
			existErr *business.Error
			expectOk bool
		)

		ctx := context.Background()
		id := "random"

		JustBeforeEach(func() {
			expectOk, existErr = redisCacheDAO.ExistOriginalURLIDInFilters(id)
		})

		Context("success", func() {
			BeforeEach(func() {
				_, err := testRedisClient.Do(ctx, "CF.ADD", originalURLIDsFilterName, id).Result()
				Expect(err).To(BeNil())
			})
			AfterEach(func() {
				_, err := testRedisClient.Del(ctx, originalURLIDsFilterName).Result()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(existErr).To(BeNil())
				actualOk, err := testRedisClient.Do(ctx, "CF.EXISTS", originalURLIDsFilterName, id).Bool()
				Expect(err).To(BeNil())
				Expect(actualOk).To(Equal(expectOk))
			})
		})
	})

	var _ = Describe("DeleteOriginalURLIDInFilters", func() {
		var (
			deleteErr *business.Error
			expectOk  bool
		)

		ctx := context.Background()
		id := "random"

		JustBeforeEach(func() {
			expectOk, deleteErr = redisCacheDAO.DeleteOriginalURLIDInFilters(id)
		})

		Context("success", func() {
			BeforeEach(func() {
				_, err := testRedisClient.Do(ctx, "CF.ADD", originalURLIDsFilterName, id).Result()
				Expect(err).To(BeNil())
			})
			AfterEach(func() {
				_, err := testRedisClient.Del(ctx, originalURLIDsFilterName).Result()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
				actualOk, err := testRedisClient.Do(ctx, "CF.EXISTS", originalURLIDsFilterName, id).Bool()
				Expect(err).To(BeNil())
				Expect(!actualOk).To(Equal(expectOk))
			})
		})
	})

	var _ = Describe("DeleteMultiOriginalURLIDInFilters", func() {
		var (
			deleteErr *business.Error
			expectOk  bool
		)

		ctx := context.Background()
		ids := []string{"random1", "random2"}

		JustBeforeEach(func() {
			expectOk, deleteErr = redisCacheDAO.DeleteMultiOriginalURLIDInFilters(ids)
		})

		Context("success", func() {
			BeforeEach(func() {
				for _, id := range ids {
					_, err := testRedisClient.Do(ctx, "CF.ADD", originalURLIDsFilterName, id).Result()
					Expect(err).To(BeNil())
				}
			})
			AfterEach(func() {
				_, err := testRedisClient.Del(ctx, originalURLIDsFilterName).Result()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
				for _, id := range ids {
					actualOk, err := testRedisClient.Do(ctx, "CF.EXISTS", originalURLIDsFilterName, id).Bool()
					Expect(err).To(BeNil())
					Expect(!actualOk).To(Equal(expectOk))
				}
			})
		})
	})
})
