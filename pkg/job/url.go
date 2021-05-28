package job

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
)

func NewExpiredURLJob(cfg ExpiredURLJobConfig, urlDAO dao.UrlDAO) *ExpiredURLJob {
	return &ExpiredURLJob{cfg: cfg, urlDAO: urlDAO}
}

type ExpiredURLJobConfig struct {
	Name            string
	TimerFormat     string
	ExpireURLNumber int
}

type ExpiredURLJob struct {
	cfg    ExpiredURLJobConfig
	urlDAO dao.UrlDAO
}

func (e *ExpiredURLJob) Name() string {
	return e.cfg.Name
}

func (e *ExpiredURLJob) Work() (map[string]interface{}, *business.Error) {
	var result = make(map[string]interface{})
	length, err := e.urlDAO.Expire(e.cfg.ExpireURLNumber)
	if err == nil {
		result["length"] = length
	}
	return result, err
}

func (e *ExpiredURLJob) TimerFormat() string {
	return e.cfg.TimerFormat
}
