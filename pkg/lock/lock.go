package lock

import (
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
)

type Locker interface {
	AcquireLock(name string, lockDuration, waitTime time.Duration) (bool, *business.Error)
	ReleaseLock(name string) *business.Error
}

func waitTimeSeries(waitTime time.Duration) []time.Duration {
	beginingMs := 800 * time.Millisecond
	endingMs := 50 * time.Millisecond

	temp := []time.Duration{0}

	t := endingMs
	total := 0 * time.Millisecond
	for t < beginingMs && total+t <= waitTime {
		temp = append(temp, t)
		total = total + t
		t = time.Duration(float64(t) * 1.2)
	}

	for total+beginingMs < waitTime {
		temp = append(temp, beginingMs)
		total = total + beginingMs
	}

	if t := waitTime - total; t > 0 {
		temp = append(temp, t)
	}

	//reverse the series
	output := []time.Duration{}
	for i, _ := range temp {
		output = append(output, temp[len(temp)-i-1])
	}

	return output
}
