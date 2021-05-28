package repository

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/internal/daomock"
	"github.com/KennyChenFight/Shortening-URL/internal/lockmock"
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("URLRepository", func() {
	var mockCtrl *gomock.Controller
	var mockUrlDAO *daomock.MockUrlDAO
	var mockCacheDAO *daomock.MockCacheDAO
	var mockLocker *lockmock.MockLocker
	var logger *loglib.Logger
	var urlRepository *URLRepository

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger = loglib.NewNopLogger()
		mockUrlDAO = daomock.NewMockUrlDAO(mockCtrl)
		mockCacheDAO = daomock.NewMockCacheDAO(mockCtrl)
		mockLocker = lockmock.NewMockLocker(mockCtrl)
		urlRepository = NewURLRepository(logger, mockUrlDAO, mockCacheDAO, mockLocker)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var _ = Describe("CreateShorteningURL", func() {
		var (
			expectURL *dao.URL
			createErr *business.Error
		)

		actualOriginalURL := "http://example.com"

		JustBeforeEach(func() {
			expectURL, createErr = urlRepository.CreateShorteningURL(actualOriginalURL)
		})

		Context("success", func() {
			var actualURL *dao.URL
			BeforeEach(func() {
				actualURL = &dao.URL{ID: "random", Original: actualOriginalURL}
				mockUrlDAO.EXPECT().Create(actualOriginalURL).Return(actualURL, nil)
				mockCacheDAO.EXPECT().SetOriginalURL(actualURL.ID, actualURL.Original).Return(nil)
			})

			It("result", func() {
				Expect(createErr).To(BeNil())
				Expect(expectURL).To(Equal(actualURL))
			})
		})

		Context("create url fail", func() {
			var createURLErr *business.Error
			BeforeEach(func() {
				createURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", nil)
				mockUrlDAO.EXPECT().Create(actualOriginalURL).Return(nil, createURLErr)
			})

			It("result", func() {
				Expect(createErr).To(Equal(createURLErr))
				Expect(expectURL).To(BeNil())
			})
		})

		Context("setOriginalURL cache fail", func() {
			var actualURL *dao.URL
			var setOriginalURLErr *business.Error
			BeforeEach(func() {
				actualURL = &dao.URL{ID: "random", Original: actualOriginalURL}
				setOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", nil)
				mockUrlDAO.EXPECT().Create(actualOriginalURL).Return(actualURL, nil)
				mockCacheDAO.EXPECT().SetOriginalURL(actualURL.ID, actualURL.Original).Return(setOriginalURLErr)
			})

			It("result", func() {
				Expect(createErr).To(Equal(setOriginalURLErr))
				Expect(expectURL).To(BeNil())
			})
		})
	})

	var _ = Describe("GetOriginalURL", func() {
		var (
			expectOriginalURL string
			getErr            *business.Error
		)

		actualID := "random"
		actualOriginalURL := "http://example.com"

		JustBeforeEach(func() {
			expectOriginalURL, getErr = urlRepository.GetOriginalURL(actualID)
		})

		Context("success with cache hit", func() {
			BeforeEach(func() {
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return(actualOriginalURL, nil)
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectOriginalURL).To(Equal(actualOriginalURL))
			})
		})

		Context("success with first cache miss but second cache hit", func() {
			var getOriginalURLErr *business.Error
			var lockName string
			BeforeEach(func() {
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", getOriginalURLErr)
				lockName = fmt.Sprintf("%s-%s", prefixLockURLResource, actualID)
				mockLocker.EXPECT().AcquireLock(lockName, lockURLResourceDuration, waitingLockURLResourceDuration).Return(true, nil)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return(actualOriginalURL, nil)
				mockLocker.EXPECT().ReleaseLock(lockName).Return(nil)
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectOriginalURL).To(Equal(actualOriginalURL))
			})
		})

		Context("success with first cache miss and second cache miss, so hit the database", func() {
			var getOriginalURLErr *business.Error
			var lockName string
			var actualURL *dao.URL
			BeforeEach(func() {
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", getOriginalURLErr)
				lockName = fmt.Sprintf("%s-%s", prefixLockURLResource, actualID)
				mockLocker.EXPECT().AcquireLock(lockName, lockURLResourceDuration, waitingLockURLResourceDuration).Return(true, nil)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", getOriginalURLErr)
				actualURL = &dao.URL{ID: actualID, Original: actualOriginalURL}
				mockUrlDAO.EXPECT().Get(actualID).Return(actualURL, nil)
				mockCacheDAO.EXPECT().SetOriginalURL(actualURL.ID, actualURL.Original).Return(nil)
				mockLocker.EXPECT().ReleaseLock(lockName).Return(nil)
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectOriginalURL).To(Equal(actualOriginalURL))
			})
		})

		Context("fail with first cache internal problem", func() {
			var getOriginalURLErr *business.Error
			BeforeEach(func() {
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New("unknown"))
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", getOriginalURLErr)
			})

			It("result", func() {
				Expect(getErr).To(Equal(getOriginalURLErr))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})

		Context("fail with first cache miss but can not acquire lock with internal problem", func() {
			var getOriginalURLErr *business.Error
			var acquireLockErr *business.Error
			var lockName string
			BeforeEach(func() {
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", getOriginalURLErr)
				acquireLockErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New("unknown"))
				lockName = fmt.Sprintf("%s-%s", prefixLockURLResource, actualID)
				mockLocker.EXPECT().AcquireLock(lockName, lockURLResourceDuration, waitingLockURLResourceDuration).Return(false, acquireLockErr)
			})

			It("result", func() {
				Expect(getErr).To(Equal(acquireLockErr))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})

		Context("fail with first cache miss but can not acquire lock with server so busy", func() {
			var getOriginalURLErr *business.Error
			var lockName string
			BeforeEach(func() {
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", getOriginalURLErr)
				lockName = fmt.Sprintf("%s-%s", prefixLockURLResource, actualID)
				mockLocker.EXPECT().AcquireLock(lockName, lockURLResourceDuration, waitingLockURLResourceDuration).Return(false, nil)
			})

			It("result", func() {
				Expect(getErr).To(Equal(business.NewError(business.AcquireLockURLResourceError, http.StatusServiceUnavailable, "server unavailable", errors.New("server unavailable"))))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})

		Context("fail with first cache miss and second cache miss with internal problem", func() {
			var firstGetOriginalURLErr *business.Error
			var secondGetOriginalURLErr *business.Error
			var lockName string
			BeforeEach(func() {
				firstGetOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", firstGetOriginalURLErr)
				lockName = fmt.Sprintf("%s-%s", prefixLockURLResource, actualID)
				mockLocker.EXPECT().AcquireLock(lockName, lockURLResourceDuration, waitingLockURLResourceDuration).Return(true, nil)
				secondGetOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New("unknown"))
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", secondGetOriginalURLErr)
				mockLocker.EXPECT().ReleaseLock(lockName).Return(nil)
			})

			It("result", func() {
				Expect(getErr).To(Equal(secondGetOriginalURLErr))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})

		Context("fail with first cache miss and second cache miss and database get URL fail", func() {
			var firstGetOriginalURLErr *business.Error
			var secondGetOriginalURLErr *business.Error
			var getURLErr *business.Error
			var lockName string
			BeforeEach(func() {
				firstGetOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", firstGetOriginalURLErr)
				lockName = fmt.Sprintf("%s-%s", prefixLockURLResource, actualID)
				mockLocker.EXPECT().AcquireLock(lockName, lockURLResourceDuration, waitingLockURLResourceDuration).Return(true, nil)
				secondGetOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", secondGetOriginalURLErr)
				getURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New("unknown"))
				mockUrlDAO.EXPECT().Get(actualID).Return(nil, getURLErr)
				mockLocker.EXPECT().ReleaseLock(lockName).Return(nil)
			})

			It("result", func() {
				Expect(getErr).To(Equal(getURLErr))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})

		Context("fail with first cache miss and second cache miss and success database get URL and set cache fail", func() {
			var firstGetOriginalURLErr *business.Error
			var secondGetOriginalURLErr *business.Error
			var setOriginalURLErr *business.Error
			var actualURL *dao.URL
			var lockName string
			BeforeEach(func() {
				firstGetOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", firstGetOriginalURLErr)
				lockName = fmt.Sprintf("%s-%s", prefixLockURLResource, actualID)
				mockLocker.EXPECT().AcquireLock(lockName, lockURLResourceDuration, waitingLockURLResourceDuration).Return(true, nil)
				secondGetOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return("", secondGetOriginalURLErr)
				actualURL = &dao.URL{ID: actualID, Original: actualOriginalURL}
				mockUrlDAO.EXPECT().Get(actualID).Return(actualURL, nil)
				setOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New("unknown"))
				mockCacheDAO.EXPECT().SetOriginalURL(actualURL.ID, actualURL.Original).Return(setOriginalURLErr)
				mockLocker.EXPECT().ReleaseLock(lockName).Return(nil)
			})

			It("result", func() {
				Expect(getErr).To(Equal(setOriginalURLErr))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})
	})
})
