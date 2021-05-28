package dao

import "time"

const randStrLength = 6
const expiredDuration = 1 * time.Hour

const prefixHotOriginalURL = "HOT-ORIGINAL-URL"
const hotOriginalURLBaseTTL = 30 * time.Minute
const randomOriginalURLTTLNumber = 60
