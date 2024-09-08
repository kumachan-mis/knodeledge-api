package testutil

import "time"

func Date() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}
