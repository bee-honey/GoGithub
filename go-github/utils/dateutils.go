package utils

import (
	"time"
)

func FormatToReadableDate(isoDate string) string {
	// Parse the ISO date format
	parsedTime, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return isoDate
	}

	// Format the date into a readable format
	readableDate := parsedTime.Format("January 2, 2006 at 3:04 PM MST")
	return readableDate
}
