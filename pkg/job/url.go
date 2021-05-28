package job

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
)

func NewExpiredURLJob(cfg ExpiredURLJobConfig, urlDAO dao.UrlDAO, cacheDAO dao.CacheDAO) *ExpiredURLJob {
	return &ExpiredURLJob{cfg: cfg, urlDAO: urlDAO, cacheDAO: cacheDAO}
}

type ExpiredURLJobConfig struct {
	Name            string
	TimerFormat     string
	ExpireURLNumber int
}

type ExpiredURLJob struct {
	cfg      ExpiredURLJobConfig
	urlDAO   dao.UrlDAO
	cacheDAO dao.CacheDAO
}

func (e *ExpiredURLJob) Name() string {
	return e.cfg.Name
}

func (e *ExpiredURLJob) Work() (map[string]interface{}, *business.Error) {
	var result = make(map[string]interface{})
	ids, err := e.urlDAO.Expire(e.cfg.ExpireURLNumber)
	if err != nil {
		return nil, err
	}
	err = e.cacheDAO.DeleteMultiOriginalURL(ids)
	if err != nil {
		return nil, err
	}

	_, err = e.cacheDAO.DeleteMultiOriginalURLIDInFilters(ids)
	if err != nil {
		return nil, err
	}

	result["length"] = len(ids)
	return result, nil
}

func (e *ExpiredURLJob) TimerFormat() string {
	return e.cfg.TimerFormat
}
