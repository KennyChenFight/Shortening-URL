package middleware

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/KennyChenFight/Shortening-URL/internal/validationtranslatormock"
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
)

var _ = Describe("GlobalErrorHandle", func() {
	var baseMiddleware *BaseMiddleware
	var mockCtrl *gomock.Controller
	var mockTranslator *validationtranslatormock.MockTranslator

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger := loglib.NewNopLogger()
		mockTranslator = validationtranslatormock.NewMockTranslator(mockCtrl)
		baseMiddleware = NewMiddleware(logger, mockTranslator, nil)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var _ = Describe("GlobalErrorHandle", func() {
		mockWriter := httptest.NewRecorder()
		gin.SetMode("release")
		ginMockContext, _ := gin.CreateTestContext(mockWriter)
		JustBeforeEach(func() {
			baseMiddleware.GlobalErrorHandle()(ginMockContext)
		})

		AfterEach(func() {
			mockWriter = httptest.NewRecorder()
			ginMockContext, _ = gin.CreateTestContext(mockWriter)
		})

		Context("send success response with default statusCode", func() {
			var actualSuccess *business.Success
			var actualResponseBody gin.H
			BeforeEach(func() {
				actualResponseBody = gin.H{"message": "ok"}
				actualSuccess = business.NewSuccess(http.StatusOK, actualResponseBody)
				ginMockContext.Set("success", actualSuccess)
			})

			It("result", func() {
				responseBodyBytes := mockWriter.Body.Bytes()
				var expectResponseBody gin.H
				err := json.Unmarshal(responseBodyBytes, &expectResponseBody)
				Expect(err).To(BeNil())
				Expect(expectResponseBody).To(Equal(actualResponseBody))
			})
		})

		Context("send success response with redirect code", func() {
			var actualSuccess *business.Success
			var actualLocation string
			var request *http.Request
			BeforeEach(func() {
				request = &http.Request{}
				ginMockContext.Request = request
				actualLocation = "http://example.com"
				actualSuccess = business.NewSuccess(http.StatusMovedPermanently, actualLocation)
				ginMockContext.Set("success", actualSuccess)
			})

			It("result", func() {
				expectLocation := mockWriter.Header().Get("location")
				Expect(expectLocation).To(Equal(actualLocation))
			})
		})

		Context("send error response with business error type", func() {
			var actualBusinessErr *business.Error
			var request *http.Request

			BeforeEach(func() {
				request = &http.Request{Header: make(map[string][]string, 0)}
				ginMockContext.Request = request
				unknownErr := errors.New("unknown")
				actualBusinessErr = business.NewError(business.Unknown, http.StatusInternalServerError, "unknown", unknownErr)
				ginMockContext.Error(actualBusinessErr)
				mockTranslator.EXPECT().Translate(ginMockContext.GetHeader("Accept-Language"), actualBusinessErr.Reason).Return(nil, nil)
			})

			It("result", func() {
				result := mockWriter.Result()
				Expect(result.StatusCode).To(Equal(actualBusinessErr.HTTPStatusCode))
				var expectResponseBody business.Error
				body, err := ioutil.ReadAll(result.Body)
				err = json.Unmarshal(body, &expectResponseBody)
				Expect(err).To(BeNil())
				Expect(expectResponseBody.Message).To(Equal(actualBusinessErr.Message))
			})
		})

		Context("send error response with unknown error type", func() {
			var actualErr error
			var request *http.Request

			BeforeEach(func() {
				request = &http.Request{Header: make(map[string][]string, 0)}
				ginMockContext.Request = request
				actualErr = errors.New("unknown")
				ginMockContext.Error(actualErr)
				mockTranslator.EXPECT().Translate(ginMockContext.GetHeader("Accept-Language"), actualErr).Return(nil, nil)
			})

			It("result", func() {
				actualBusinessErr := business.NewError(business.Unknown, http.StatusInternalServerError, actualErr.Error(), actualErr)
				result := mockWriter.Result()
				Expect(result.StatusCode).To(Equal(actualBusinessErr.HTTPStatusCode))
				var expectResponseBody business.Error
				body, err := ioutil.ReadAll(result.Body)
				err = json.Unmarshal(body, &expectResponseBody)
				Expect(err).To(BeNil())
				Expect(expectResponseBody.Message).To(Equal(actualBusinessErr.Message))
			})
		})
	})
})
