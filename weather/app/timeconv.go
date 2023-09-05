package app

import (
	"time"
)

func dayMultiplier(currentDate time.Time) int {
	if currentDate.Weekday() == time.Sunday {
		return 6
	}

	return int(currentDate.Weekday()) - 1
}
