package dao

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/pglib"
	"github.com/KennyChenFight/randstr"
	"time"
)

func NewPGUrlDAO(logger *loglib.Logger, client *pglib.GOPGClient, randomStrGenerator randstr.RandomStrGenerator) *PGUrlDAO {
	return &PGUrlDAO{logger: logger, client: client, randomStrGenerator: randomStrGenerator}
}

type PGUrlDAO struct {
	logger             *loglib.Logger
	client             *pglib.GOPGClient
	randomStrGenerator randstr.RandomStrGenerator
}

const randStrLength = 6
const expiredDuration = 10 * time.Second

func (p *PGUrlDAO) Create(originalURL string) (*URL, *business.Error) {
	now := time.Now()
	url := &URL{
		ID:        p.randomStrGenerator.GenerateRandomStr(randStrLength),
		Original:  originalURL,
		CreatedAt: now,
		ExpiredAt: now.Add(expiredDuration),
	}
	_, err := p.client.Model(url).Insert()
	if err != nil {
		return nil, pgErrorHandle(p.logger, err)
	}
	return url, nil
}

func (p *PGUrlDAO) BatchCreate(originalURL string) *business.Error {
	now := time.Now()
	var urls []URL
	for i := 0; i < 1000000; i++ {
		urls = append(urls, URL{
			ID:        p.randomStrGenerator.GenerateRandomStr(randStrLength),
			Original:  originalURL,
			CreatedAt: now,
			ExpiredAt: now.Add(expiredDuration),
		})
	}
	_, err := p.client.Model(&urls).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return pgErrorHandle(p.logger, err)
	}
	return nil
}

func (p *PGUrlDAO) Get(id string) (*URL, *business.Error) {
	url := &URL{
		ID: id,
	}
	err := p.client.Model(url).
		WherePK().
		Where("expired_at > ?", time.Now()).Select()
	if err != nil {
		return nil, pgErrorHandle(p.logger, err)
	}
	return url, nil
}

func (p *PGUrlDAO) Delete(id string) *business.Error {
	url := &URL{
		ID: id,
	}
	_, err := p.client.Model(url).WherePK().Delete()
	if err != nil {
		return pgErrorHandle(p.logger, err)
	}
	return nil
}
