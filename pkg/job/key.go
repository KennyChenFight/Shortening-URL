package job

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
)

func NewGenerateKeyJob(cfg GenerateKeyJobConfig, keyDAO dao.KeyDAO) *GenerateKeyJob {
	return &GenerateKeyJob{cfg: cfg, keyDAO: keyDAO}
}

type GenerateKeyJobConfig struct {
	Name           string
	EveryKeyNumber int
	TimerFormat    string
}

type GenerateKeyJob struct {
	cfg    GenerateKeyJobConfig
	keyDAO dao.KeyDAO
}

func (g *GenerateKeyJob) Name() string {
	return g.cfg.Name
}

func (g *GenerateKeyJob) Work() (map[string]interface{}, *business.Error) {
	var result = make(map[string]interface{})
	length, err := g.keyDAO.BatchCreate(g.cfg.EveryKeyNumber)
	if err == nil {
		result["length"] = length
	}
	return result, err
}

func (g *GenerateKeyJob) TimerFormat() string {
	return g.cfg.TimerFormat
}
