package common

import (
	"time"
)

// GetTimeElapse call api openweathermap
func GetTimeElapse(timestamp string) time.Duration {
	varTimestamp, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		varTimestamp = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 651387237, time.UTC).AddDate(0, 0, -1)
	}
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 651387237, time.UTC)
	duration := now.Sub(varTimestamp)
	return duration
}
