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
	var mockKeyDAO *daomock.MockKeyDAO
	var mockLocker *lockmock.MockLocker
	var logger *loglib.Logger
	var urlRepository *URLRepository

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger = loglib.NewNopLogger()
		mockUrlDAO = daomock.NewMockUrlDAO(mockCtrl)
		mockKeyDAO = daomock.NewMockKeyDAO(mockCtrl)
		mockCacheDAO = daomock.NewMockCacheDAO(mockCtrl)
		mockLocker = lockmock.NewMockLocker(mockCtrl)
		urlRepository = NewURLRepository(logger, mockUrlDAO, mockKeyDAO, mockCacheDAO, mockLocker)
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
				mockCacheDAO.EXPECT().AddOriginalURLIDInFilters(actualURL.ID).Return(nil)
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
				mockCacheDAO.EXPECT().AddOriginalURLIDInFilters(actualURL.ID).Return(nil)
			})

			It("result", func() {
				Expect(createErr).To(BeNil())
				Expect(expectURL).To(Equal(actualURL))
			})
		})

		Context("addOriginalURLIDInFilters cache fail", func() {
			var actualURL *dao.URL
			var addOriginalURLIDInFiltersErr *business.Error
			BeforeEach(func() {
				actualURL = &dao.URL{ID: "random", Original: actualOriginalURL}
				addOriginalURLIDInFiltersErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", nil)
				mockUrlDAO.EXPECT().Create(actualOriginalURL).Return(actualURL, nil)
				mockCacheDAO.EXPECT().SetOriginalURL(actualURL.ID, actualURL.Original).Return(nil)
				mockCacheDAO.EXPECT().AddOriginalURLIDInFilters(actualURL.ID).Return(addOriginalURLIDInFiltersErr)
			})

			It("result", func() {
				Expect(createErr).To(BeNil())
				Expect(expectURL).To(Equal(actualURL))
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
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
				mockCacheDAO.EXPECT().GetOriginalURL(actualID).Return(actualOriginalURL, nil)
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectOriginalURL).To(Equal(actualOriginalURL))
			})
		})

		Context("success with ignore filter internal problem", func() {
			var existOriginalURLIDInFiltersErr *business.Error
			var getOriginalURLErr *business.Error
			var lockName string
			BeforeEach(func() {
				existOriginalURLIDInFiltersErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", nil)
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(false, existOriginalURLIDInFiltersErr)
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

		Context("success with first cache miss but second cache hit", func() {
			var getOriginalURLErr *business.Error
			var lockName string
			BeforeEach(func() {
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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

		Context("fail with not exist in filter", func() {
			var exist bool
			var notExistErr *business.Error
			BeforeEach(func() {
				exist = false
				notExistErr = business.NewError(business.NotFound, http.StatusNotFound, "record not found", errors.New("can not found in filters"))
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(exist, nil)
			})

			It("result", func() {
				Expect(getErr).To(Equal(notExistErr))
				Expect(expectOriginalURL).To(Equal(""))
			})
		})

		Context("fail with first cache internal problem", func() {
			var getOriginalURLErr *business.Error
			BeforeEach(func() {
				getOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", errors.New("unknown"))
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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

		Context("success with first cache miss and second cache miss and success database get URL and set cache fail", func() {
			var firstGetOriginalURLErr *business.Error
			var secondGetOriginalURLErr *business.Error
			var setOriginalURLErr *business.Error
			var actualURL *dao.URL
			var lockName string
			BeforeEach(func() {
				firstGetOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", dao.RedisErrKeyNotExist)
				mockCacheDAO.EXPECT().ExistOriginalURLIDInFilters(actualID).Return(true, nil)
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
				Expect(getErr).To(BeNil())
				Expect(expectOriginalURL).To(Equal(actualOriginalURL))
			})
		})
	})

	var _ = Describe("DeleteShorteningURL", func() {
		var (
			deleteErr *business.Error
		)

		actualID := "random"

		JustBeforeEach(func() {
			deleteErr = urlRepository.DeleteShorteningURL(actualID)
		})

		Context("success", func() {
			BeforeEach(func() {
				mockCacheDAO.EXPECT().DeleteOriginalURL(actualID).Return(nil)
				mockUrlDAO.EXPECT().Delete(actualID).Return(nil)
				mockCacheDAO.EXPECT().DeleteOriginalURLIDInFilters(actualID).Return(true, nil)
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
			})
		})

		Context("success with fail to delete originalURL in cache", func() {
			var deleteOriginalURLErr *business.Error
			BeforeEach(func() {
				deleteOriginalURLErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", nil)
				mockCacheDAO.EXPECT().DeleteOriginalURL(actualID).Return(deleteOriginalURLErr)
				mockUrlDAO.EXPECT().Delete(actualID).Return(nil)
				mockCacheDAO.EXPECT().DeleteOriginalURLIDInFilters(actualID).Return(true, nil)
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
			})
		})

		Context("success with fail to delete originalURLID in filters because internal problem", func() {
			var deleteOriginalURLInFilterErr *business.Error
			BeforeEach(func() {
				deleteOriginalURLInFilterErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", nil)
				mockCacheDAO.EXPECT().DeleteOriginalURL(actualID).Return(nil)
				mockUrlDAO.EXPECT().Delete(actualID).Return(nil)
				mockCacheDAO.EXPECT().DeleteOriginalURLIDInFilters(actualID).Return(false, deleteOriginalURLInFilterErr)
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
			})
		})

		Context("success with fail to delete originalURLID in filters because not exist", func() {
			BeforeEach(func() {
				mockCacheDAO.EXPECT().DeleteOriginalURL(actualID).Return(nil)
				mockUrlDAO.EXPECT().Delete(actualID).Return(nil)
				mockCacheDAO.EXPECT().DeleteOriginalURLIDInFilters(actualID).Return(false, nil)
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
			})
		})

		Context("fail with delete in database", func() {
			var deleteInDBErr *business.Error
			BeforeEach(func() {
				deleteInDBErr = business.NewError(business.Unknown, http.StatusInternalServerError, "", nil)
				mockCacheDAO.EXPECT().DeleteOriginalURL(actualID).Return(nil)
				mockUrlDAO.EXPECT().Delete(actualID).Return(deleteInDBErr)
			})

			It("result", func() {
				Expect(deleteErr).To(Equal(deleteInDBErr))
			})
		})
	})

})
