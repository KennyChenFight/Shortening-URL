package dao

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/pglib"
	"github.com/KennyChenFight/randstr"
	"time"
)

func NewPGKeyDAO(logger *loglib.Logger, client *pglib.GOPGClient, randomStrGenerator randstr.RandomStrGenerator) *PGKeyDAO {
	return &PGKeyDAO{logger: logger, client: client, randomStrGenerator: randomStrGenerator}
}

type PGKeyDAO struct {
	logger             *loglib.Logger
	client             *pglib.GOPGClient
	randomStrGenerator randstr.RandomStrGenerator
}

func (p *PGKeyDAO) BatchCreate(num int) (int, *business.Error) {
	var keys []Key
	now := time.Now()
	for i := 0; i < num; i++ {
		keys = append(keys, Key{
			ID:        p.randomStrGenerator.GenerateRandomStr(randStrLength),
			CreatedAt: now,
		})
	}
	res, err := p.client.Model(&keys).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return 0, pgErrorHandle(p.logger, err)
	}
	return res.RowsAffected(), nil
}
