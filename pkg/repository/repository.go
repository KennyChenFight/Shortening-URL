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
	// todo setCache 失敗 則需要回滾database create operation
	// 如果擔心setCache這段太久而導致transaction等太久 也許可以採用context來設定timout較短時間 來終止transaction
	// 當然也可以選擇set cache失敗的話寫log 忽略錯誤 讓下次get的時候再去setCache
	// 但可能導致get loading重 但是get有上lock 以防大量地get進來
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
	// todo 考慮這邊再多加上bloom filter或是Cuckoo Filters(可支援刪除)
	originalURL, err := u.CacheDAO.GetOriginalURL(id)
	if err != nil {
		if err.Reason == dao.RedisErrKeyNotExist {
			ok, err := u.locker.AcquireLock(fmt.Sprintf("%s-%s", prefixLockURLResource, id), lockURLResourceDuration, waitingLockURLResourceDuration)
			if err != nil {
				return "", err
			}
			if !ok {
				return "", business.NewError(business.AcquireLockURLResourceError, http.StatusServiceUnavailable, "server unavailable", errors.New("server unavailable"))
			} else {
				defer u.locker.ReleaseLock(fmt.Sprintf("%s-%s", prefixLockURLResource, id))
			}
			originalURL, err = u.CacheDAO.GetOriginalURL(id)
			if err != nil {
				if err.Reason == dao.RedisErrKeyNotExist {
					url, err := u.UrlDAO.Get(id)
					if err != nil {
						return "", err
					}
					// todo 如果setCache失敗也回傳正確的originalURL出去
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
	// todo 需要expire cache 如果 expire cache操作失敗則應該回滾database delete operation
	err := u.UrlDAO.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
