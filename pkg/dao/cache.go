package dao

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
)

type CacheDAO interface {
	GetOriginalURL(name string) (string, *business.Error)
	SetOriginalURL(name string, originalURL string) *business.Error
	DeleteOriginalURL(name string) *business.Error
	DeleteMultiOriginalURL(names []string) *business.Error
	AddOriginalURLIDInFilters(originalURL string) *business.Error
	ExistOriginalURLIDInFilters(originalURL string) (bool, *business.Error)
	DeleteOriginalURLIDInFilters(originalURL string) (bool, *business.Error)
	DeleteMultiOriginalURLIDInFilters(originalURLIDs []string) (bool, *business.Error)
}
