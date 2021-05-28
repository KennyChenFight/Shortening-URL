package dao

import (
	"math/rand"
	"time"
)

const randStrLength = 6
const expiredDuration = 1 * time.Hour

const prefixHotOriginalURL = "ORIGINAL-URL-ID"
const hotOriginalURLBaseTTL = 30 * time.Minute
const randomOriginalURLTTLNumber = 60

var randomFunc = rand.Int63n

func getRandomOriginalURLTTLSecond() time.Duration {
	return time.Duration(randomFunc(int64(randomOriginalURLTTLNumber))+1) * time.Second
}

const originalURLIDsFilterName = "FILTER-ORIGINAL-URL-IDs"
