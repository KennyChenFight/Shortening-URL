package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"

	"github.com/gin-gonic/gin"

	"github.com/KennyChenFight/Shortening-URL/internal/repositorymock"
	"github.com/golang/mock/gomock"

	"github.com/KennyChenFight/golib/loglib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseService", func() {
	var baseService *BaseService
	var mockCtrl *gomock.Controller
	var repositoryMock *repositorymock.MockRepository
	var config *Config

	BeforeEach(func() {
		config = &Config{FQDN: "http://example.com"}
		logger := loglib.NewNopLogger()
		mockCtrl = gomock.NewController(GinkgoT())
		repositoryMock = repositorymock.NewMockRepository(mockCtrl)
		baseService = NewService(config, logger, repositoryMock)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var _ = Describe("CreateShorteningURL", func() {
		ginMockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		JustBeforeEach(func() {
			baseService.CreateShorteningURL(ginMockContext)
		})

		Context("success", func() {
			var mockRequest *http.Request
			var mockRequestBody = make(map[string]interface{}, 0)
			var actualURL string
			var shorteningURL *dao.URL
			BeforeEach(func() {
				actualURL = "http://test.com"
				mockRequestBody["url"] = actualURL
				b, err := json.Marshal(&mockRequestBody)
				Expect(err).To(BeNil())
				mockRequest, err = http.NewRequest("POST", "http://server.com", bytes.NewBuffer(b))
				Expect(err).To(BeNil())
				ginMockContext.Request = mockRequest

				shorteningURL = &dao.URL{
					ID:        "abcdef",
					Original:  actualURL,
					CreatedAt: time.Now(),
					ExpiredAt: time.Now(),
				}
				repositoryMock.EXPECT().CreateShorteningURL(actualURL).Return(shorteningURL, nil)
			})

			It("result", func() {
				actualSuccess := business.NewSuccess(http.StatusCreated, gin.H{"id": shorteningURL.ID, "shortUrl": combineFQDNWithShorteningURLID(baseService.config.FQDN, shorteningURL.ID)})
				expectSuccess, _ := ginMockContext.Get("success")
				Expect(expectSuccess).To(Equal(actualSuccess))
			})
		})

		Context("binding validation fail", func() {
			var mockRequest *http.Request
			var mockRequestBody = make(map[string]interface{}, 0)
			var actualURL string
			BeforeEach(func() {
				actualURL = "http://test.com"
				mockRequestBody["url"] = strings.Repeat(actualURL, 2048)
				b, err := json.Marshal(&mockRequestBody)
				Expect(err).To(BeNil())
				mockRequest, err = http.NewRequest("POST", "http://server.com", bytes.NewBuffer(b))
				Expect(err).To(BeNil())
				ginMockContext.Request = mockRequest
			})

			It("result", func() {
				expectError := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err
				businessError, ok := expectError.(*business.Error)
				Expect(ok).To(Equal(true))

				actualError := business.NewError(business.Validation, http.StatusBadRequest, "invalid url field", businessError.Reason)
				Expect(expectError).To(Equal(actualError))
			})
		})

		Context("create shortening url fail", func() {
			var mockRequest *http.Request
			var mockRequestBody = make(map[string]interface{}, 0)
			var actualURL string
			var createErr *business.Error
			BeforeEach(func() {
				actualURL = "http://test.com"
				mockRequestBody["url"] = actualURL
				b, err := json.Marshal(&mockRequestBody)
				Expect(err).To(BeNil())
				mockRequest, err = http.NewRequest("POST", "http://server.com", bytes.NewBuffer(b))
				Expect(err).To(BeNil())
				ginMockContext.Request = mockRequest

				createErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New(""))
				repositoryMock.EXPECT().CreateShorteningURL(actualURL).Return(nil, createErr)
			})

			It("result", func() {
				expectError := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err
				businessError, ok := expectError.(*business.Error)
				Expect(ok).To(Equal(true))
				Expect(businessError).To(Equal(createErr))
			})
		})
	})

	var _ = Describe("GetOriginalURL", func() {
		ginMockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		JustBeforeEach(func() {
			baseService.GetOriginalURL(ginMockContext)
		})

		Context("success", func() {
			var actualID string
			var originalURL string
			BeforeEach(func() {
				actualID = "random"
				ginMockContext.Params = gin.Params{
					{
						Key:   "id",
						Value: actualID,
					},
				}
				originalURL = "http://example.com"
				repositoryMock.EXPECT().GetOriginalURL(actualID).Return(originalURL, nil)
			})

			It("result", func() {
				actualSuccess := business.NewSuccess(http.StatusTemporaryRedirect, originalURL)
				expectSuccess, _ := ginMockContext.Get("success")
				Expect(expectSuccess).To(Equal(actualSuccess))
			})
		})

		Context("binding validation fail", func() {
			var actualID string
			BeforeEach(func() {
				actualID = "test"
				ginMockContext.Params = gin.Params{
					{
						Key:   "id",
						Value: actualID,
					},
				}
			})
			It("result", func() {
				expectError := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err
				businessError, ok := expectError.(*business.Error)
				Expect(ok).To(Equal(true))

				actualError := business.NewError(business.Validation, http.StatusBadRequest, "invalid id field", businessError.Reason)
				Expect(expectError).To(Equal(actualError))
			})
		})

		Context("get originalURL fail", func() {
			var actualID string
			var getErr *business.Error
			BeforeEach(func() {
				actualID = "random"
				ginMockContext.Params = gin.Params{
					{
						Key:   "id",
						Value: actualID,
					},
				}
				getErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New(""))
				repositoryMock.EXPECT().GetOriginalURL(actualID).Return("", getErr)
			})

			It("result", func() {
				expectError := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err
				businessError, ok := expectError.(*business.Error)
				Expect(ok).To(Equal(true))

				Expect(businessError).To(Equal(getErr))
			})
		})
	})

	var _ = Describe("DeleteShorteningURL", func() {
		ginMockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		JustBeforeEach(func() {
			baseService.DeleteShorteningURL(ginMockContext)
		})

		Context("success", func() {
			var actualID string
			BeforeEach(func() {
				actualID = "random"
				ginMockContext.Params = gin.Params{
					{
						Key:   "id",
						Value: actualID,
					},
				}
				repositoryMock.EXPECT().DeleteShorteningURL(actualID).Return(nil)
			})

			It("result", func() {
				actualSuccess := business.NewSuccess(http.StatusNoContent, nil)
				expectSuccess, _ := ginMockContext.Get("success")
				Expect(expectSuccess).To(Equal(actualSuccess))
			})
		})

		Context("binding validation fail", func() {
			var actualID string
			BeforeEach(func() {
				actualID = "test"
				ginMockContext.Params = gin.Params{
					{
						Key:   "id",
						Value: actualID,
					},
				}
			})
			It("result", func() {
				expectError := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err
				businessError, ok := expectError.(*business.Error)
				Expect(ok).To(Equal(true))

				actualError := business.NewError(business.Validation, http.StatusBadRequest, "invalid id field", businessError.Reason)
				Expect(expectError).To(Equal(actualError))
			})
		})

		Context("delete shorteningURL fail", func() {
			var actualID string
			var deleteErr *business.Error
			BeforeEach(func() {
				actualID = "random"
				ginMockContext.Params = gin.Params{
					{
						Key:   "id",
						Value: actualID,
					},
				}
				deleteErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New(""))
				repositoryMock.EXPECT().DeleteShorteningURL(actualID).Return(deleteErr)
			})

			It("result", func() {
				expectError := ginMockContext.Errors[len(ginMockContext.Errors)-1].Err
				businessError, ok := expectError.(*business.Error)
				Expect(ok).To(Equal(true))

				Expect(businessError).To(Equal(deleteErr))
			})
		})
	})
})
