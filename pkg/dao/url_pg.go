package dao

import (
	"context"
	"errors"
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/pglib"
	"github.com/go-pg/pg/v10"
)

func NewPGUrlDAO(logger *loglib.Logger, client *pglib.GOPGClient) *PGUrlDAO {
	return &PGUrlDAO{logger: logger, client: client}
}

type PGUrlDAO struct {
	logger *loglib.Logger
	client *pglib.GOPGClient
}

func (p *PGUrlDAO) Create(originalURL string) (*URL, *business.Error) {
	var url URL
	var key Key
	now := time.Now()
	err := p.client.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		res, err := tx.Model((*URL)(nil)).Query(pg.Scan(&url.ID, &url.Original, &url.CreatedAt, &url.ExpiredAt), "INSERT INTO urls SELECT id, ?, ?, ? FROM keys FOR UPDATE SKIP LOCKED LIMIT 1 RETURNING id, original", originalURL, now, now.Add(expiredDuration))
		if err != nil {
			return err
		}
		if res.RowsReturned() == 0 {
			return errors.New(PGErrMsgNoRowsFound)
		}

		key.ID = url.ID
		_, err = tx.Model(&key).WherePK().Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, pgErrorHandle(p.logger, err)
	}
	return &url, nil
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

func (p *PGUrlDAO) Expire(num int) ([]string, *business.Error) {
	var ids []string
	subQuery := p.client.Model((*URL)(nil)).Column("id").Where("expired_at < ?", time.Now()).Limit(num)
	_, err := p.client.Model((*URL)(nil)).
		Where("id in (?)", subQuery).
		Returning("id").
		Delete(&ids)
	if err != nil {
		return ids, pgErrorHandle(p.logger, err)
	}
	return ids, nil
}
