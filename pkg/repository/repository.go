package repository

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
	"github.com/KennyChenFight/Shortening-URL/pkg/lock"
	"github.com/KennyChenFight/golib/loglib"
)

type Repository interface {
	CreateShorteningURL(originalURL string) (*dao.URL, *business.Error)
	GetOriginalURL(id string) (string, *business.Error)
	DeleteShorteningURL(id string) *business.Error
	BatchCreateKeys(num int) (int, *business.Error)
}

func NewURLRepository(logger *loglib.Logger, urlDAO dao.UrlDAO, keyDAO dao.KeyDAO, cacheDAO dao.CacheDAO, locker lock.Locker) *URLRepository {
	return &URLRepository{
		logger:   logger,
		UrlDAO:   urlDAO,
		KeyDAO:   keyDAO,
		CacheDAO: cacheDAO,
		locker:   locker,
	}
}

type URLRepository struct {
	logger   *loglib.Logger
	UrlDAO   dao.UrlDAO
	KeyDAO   dao.KeyDAO
	CacheDAO dao.CacheDAO
	locker   lock.Locker
}

func (u *URLRepository) CreateShorteningURL(originalURL string) (*dao.URL, *business.Error) {
	url, err := u.UrlDAO.Create(originalURL)
	if err != nil {
		return nil, err
	}
	err = u.CacheDAO.SetOriginalURL(url.ID, url.Original)
	if err != nil {
		u.logger.Error("fail to set originalURL in cache", zap.Error(err))
	}
	err = u.CacheDAO.AddOriginalURLIDInFilters(url.ID)
	if err != nil {
		u.logger.Error("fail to set originalURL in filter", zap.Error(err))
	}
	return url, nil
}

func (u *URLRepository) GetOriginalURL(id string) (string, *business.Error) {
	// 避免太多random不存在的key的訪問 可以利用這個先擋著
	exist, err := u.CacheDAO.ExistOriginalURLIDInFilters(id)
	if err != nil {
		u.logger.Error("fail to check originalURLID in filters", zap.Error(err))
	} else {
		if !exist {
			return "", business.NewError(business.NotFound, http.StatusNotFound, "record not found", errors.New("can not found in filters"))
		}
	}

	// 先從first cache 拿
	originalURL, err := u.CacheDAO.GetOriginalURL(id)
	if err != nil {
		if err.Reason == dao.RedisErrKeyNotExist {
			// 如果first cache miss 則 需要獲取lock 避免當cache失效時 太多request過來要更新cache 使用 lock 只能有一個進來訪問database並更新cache
			// 這樣以來其他人就可以透過second cache hit來拿到資料 而不用真的訪問到database
			ok, err := u.locker.AcquireLock(fmt.Sprintf("%s-%s", prefixLockURLResource, id), lockURLResourceDuration, waitingLockURLResourceDuration)
			if err != nil {
				return "", err
			}
			if !ok {
				return "", business.NewError(business.AcquireLockURLResourceError, http.StatusServiceUnavailable, "server unavailable", errors.New("server unavailable"))
			} else {
				defer u.locker.ReleaseLock(fmt.Sprintf("%s-%s", prefixLockURLResource, id))
			}
			// second cache check
			originalURL, err = u.CacheDAO.GetOriginalURL(id)
			if err != nil {
				if err.Reason == dao.RedisErrKeyNotExist {
					url, err := u.UrlDAO.Get(id)
					if err != nil {
						return "", err
					}
					err = u.CacheDAO.SetOriginalURL(url.ID, url.Original)
					if err != nil {
						u.logger.Error("fail to set originalURL cache", zap.Error(err))
					}
					return url.Original, nil
				}
				return "", err
			}
			return originalURL, nil
		}
		return "", err
	}
	return originalURL, nil
}

func (u *URLRepository) DeleteShorteningURL(id string) *business.Error {
	err := u.CacheDAO.DeleteOriginalURL(id)
	if err != nil {
		u.logger.Error("fail to delete originalURL in cache", zap.Error(err))
	}
	err = u.UrlDAO.Delete(id)
	if err != nil {
		return err
	}
	ok, err := u.CacheDAO.DeleteOriginalURLIDInFilters(id)
	if err != nil || !ok {
		u.logger.Error("fail to delete originalURL in filter", zap.Bool("ok", ok), zap.Error(err))
	}
	return nil
}

func (u *URLRepository) BatchCreateKeys(num int) (int, *business.Error) {
	return u.KeyDAO.BatchCreate(num)
}
