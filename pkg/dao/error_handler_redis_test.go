package dao

import (
	"errors"
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("redisErrorHandle", func() {
	var originalError error
	var businessError *business.Error
	JustBeforeEach(func() {
		businessError = redisErrorHandle(loglib.NewNopLogger(), originalError)
	})

	Context("when err == redis.Nil", func() {
		BeforeEach(func() {
			originalError = redis.Nil
		})
		AfterEach(func() {
			originalError = nil
		})
		It("result", func() {
			Expect(businessError).To(Equal(business.NewError(business.NotFound, http.StatusNotFound, "record not found", RedisErrKeyNotExist)))
		})
	})

	Context("internal error", func() {
		internalError := errors.New("internal error")
		BeforeEach(func() {
			originalError = internalError
		})
		AfterEach(func() {
			originalError = nil
		})
		It("result", func() {
			Expect(businessError).To(Equal(business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", internalError)))
		})
	})
})
