package repository

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
	"github.com/KennyChenFight/Shortening-URL/pkg/lock"
	"github.com/KennyChenFight/golib/loglib"
)

type Repository interface {
	CreateShorteningURL(originalURL string) (*dao.URL, *business.Error)
	GetOriginalURL(id string) (string, *business.Error)
	DeleteShorteningURL(id string) *business.Error
}

func NewURLRepository(logger *loglib.Logger, urlDAO dao.UrlDAO, cacheDAO dao.CacheDAO, locker lock.Locker) *URLRepository {
	return &URLRepository{
		logger:   logger,
		UrlDAO:   urlDAO,
		CacheDAO: cacheDAO,
		locker:   locker,
	}
}

type URLRepository struct {
	logger   *loglib.Logger
	UrlDAO   dao.UrlDAO
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
		return nil, err
	}
	return url, nil
}

func (u *URLRepository) GetOriginalURL(id string) (string, *business.Error) {
	originalURL, err := u.CacheDAO.GetOriginalURL(id)
	if err != nil {
		if err.Reason == dao.RedisErrKeyNotExist {
			ok, err := u.locker.AcquireLock(fmt.Sprintf("%s-%s", prefixLockURLResource, id), lockURLResourceDuration, waitingLockURLResourceDuration)
			if err != nil {
				return "", err
			}
			if !ok {
				return "", business.NewError(business.AcquireLockURLResourceError, http.StatusServiceUnavailable, "server unavailable", errors.New("server unavailable"))
			}
			defer u.locker.ReleaseLock(fmt.Sprintf("%s-%s", prefixLockURLResource, id))
			originalURL, err = u.CacheDAO.GetOriginalURL(id)
			if err != nil {
				if err.Reason == dao.RedisErrKeyNotExist {
					url, err := u.UrlDAO.Get(id)
					if err != nil {
						return "", err
					}
					err = u.CacheDAO.SetOriginalURL(url.ID, url.Original)
					if err != nil {
						return "", err
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
	return u.UrlDAO.Delete(id)
}
