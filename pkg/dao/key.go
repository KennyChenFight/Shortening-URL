package dao

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"time"
)

type Key struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type KeyDAO interface {
	BatchCreate(num int) (int, *business.Error)
}
