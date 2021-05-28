package dao

import (
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
)

type URL struct {
	ID        string    `json:"id"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type UrlDAO interface {
	Create(originalURL string) (*URL, *business.Error)
	Get(id string) (*URL, *business.Error)
	Delete(id string) *business.Error
	Expire(num int) ([]string, *business.Error)
}
