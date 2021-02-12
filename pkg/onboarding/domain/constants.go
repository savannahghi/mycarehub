package domain

import "time"

var (
	// TimeLocation ...
	TimeLocation, _ = time.LoadLocation("Africa/Nairobi")

	// TimeFormatStr date time string format
	TimeFormatStr = "2006-01-02T15:04:05+03:00"
)
