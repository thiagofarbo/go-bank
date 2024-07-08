package util

import (
	"fmt"
	"time"
)

func ToDate(date string) (time.Time, error) {

	layout := "2006-01-02 15:04:05"
	formattedDate, err := time.Parse(layout, date)
	if err != nil {
		fmt.Errorf("error to convert the date: %v", err)
		return time.Time{}, err
	}
	return formattedDate, err
}
