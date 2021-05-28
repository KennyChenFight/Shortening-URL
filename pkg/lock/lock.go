package lock

import (
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
)

type Locker interface {
	AcquireLock(name string, lockDuration, waitTime time.Duration) (bool, *business.Error)
	ReleaseLock(name string) *business.Error
}

// 可以根據系統的performance來做更改
// 避免等lock花太久時間
func waitTimeSeries(waitTime time.Duration) []time.Duration {
	beginning := 800 * time.Millisecond
	ending := 50 * time.Millisecond

	temp := []time.Duration{0}

	t := ending
	total := 0 * time.Millisecond
	for t < beginning && total+t <= waitTime {
		temp = append(temp, t)
		total = total + t
		t = time.Duration(float64(t) * 1.2)
	}

	for total+beginning < waitTime {
		temp = append(temp, beginning)
		total = total + beginning
	}

	if t = waitTime - total; t > 0 {
		temp = append(temp, t)
	}

	output := make([]time.Duration, 0)
	for i := range temp {
		output = append(output, temp[len(temp)-i-1])
	}
	return output
}
