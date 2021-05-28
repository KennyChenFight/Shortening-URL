package dao

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"time"
)

type URL struct {
	ID        string    `json:"id"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type UrlDAO interface {
	Create(originalURL string) (*URL, *business.Error)
	BatchCreate(originalURL string) *business.Error
	Get(id string) (*URL, *business.Error)
	Delete(id string) *business.Error
}
