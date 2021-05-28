package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"

	"github.com/KennyChenFight/Shortening-URL/internal/ratelimitermock"

	"github.com/KennyChenFight/golib/loglib"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
)

var _ = Describe("GlobalErrorHandle", func() {
	var baseMiddleware *BaseMiddleware
	var mockCtrl *gomock.Controller
	var mockRateLimiter *ratelimitermock.MockRateLimiter

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger := loglib.NewNopLogger()
		mockRateLimiter = ratelimitermock.NewMockRateLimiter(mockCtrl)
		baseMiddleware = NewMiddleware(logger, nil, mockRateLimiter)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var _ = Describe("SlideWindowRateLimiter", func() {
		gin.SetMode("release")
		mockWriter := httptest.NewRecorder()
		ginMockContext, _ := gin.CreateTestContext(mockWriter)
		RateLimiterHeaderKey := "X-RateLimit-Total"
		JustBeforeEach(func() {
			baseMiddleware.SlideWindowRateLimiter()(ginMockContext)
		})

		AfterEach(func() {
			mockWriter = httptest.NewRecorder()
			ginMockContext, _ = gin.CreateTestContext(mockWriter)
		})

		Context("success incr and pass", func() {
			var actualTotal int64
			var request *http.Request
			BeforeEach(func() {
				actualTotal = 1
				request = httptest.NewRequest("GET", "http://example.com", nil)
				request.Header.Set("X-Forwarded-For", "127.0.0.1")
				ginMockContext.Request = request
				key := request.URL.Path + "-" + request.Method + "-" + "127.0.0.1"
				mockRateLimiter.EXPECT().Incr(gomock.Any(), key, gomock.Any()).Return(actualTotal, nil)
			})

			It("result", func() {
				expectTotal := mockWriter.Header().Get(RateLimiterHeaderKey)
				Expect(expectTotal).To(Equal(strconv.FormatInt(actualTotal, 10)))
			})
		})

		Context("fail incr with internal problem", func() {
			var request *http.Request
			var incrErr error
			BeforeEach(func() {
				request = httptest.NewRequest("GET", "http://example.com", nil)
				request.Header.Set("X-Forwarded-For", "127.0.0.1")
				ginMockContext.Request = request
				key := request.URL.Path + "-" + request.Method + "-" + "127.0.0.1"
				incrErr = errors.New("internal")
				mockRateLimiter.EXPECT().Incr(gomock.Any(), key, gomock.Any()).Return(int64(0), incrErr)
			})

			It("result", func() {
				actualBusinessErr := business.NewError(business.RedisInternalError, http.StatusInternalServerError, "internal error", incrErr)
				expectBusinessErr := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err.(*business.Error)
				Expect(expectBusinessErr).To(Equal(actualBusinessErr))
			})
		})

		Context("fail incr with exceed rate limit", func() {
			var request *http.Request
			BeforeEach(func() {
				request = httptest.NewRequest("GET", "http://example.com", nil)
				request.Header.Set("X-Forwarded-For", "127.0.0.1")
				ginMockContext.Request = request
				key := request.URL.Path + "-" + request.Method + "-" + "127.0.0.1"
				mockRateLimiter.EXPECT().Incr(gomock.Any(), key, gomock.Any()).Return(int64(-1), nil)
			})

			It("result", func() {
				actualBusinessErr := business.NewError(business.TooManyRequest, http.StatusTooManyRequests, "reach limit for this endpoint", errors.New("reach limit for this endpoint"))
				expectBusinessErr := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err.(*business.Error)
				Expect(expectBusinessErr).To(Equal(actualBusinessErr))
			})
		})
	})
})
