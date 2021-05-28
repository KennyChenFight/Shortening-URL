package dao

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
)

type CacheDAO interface {
	GetOriginalURL(name string) (string, *business.Error)
	SetOriginalURL(name string, originalURL string) *business.Error
}
