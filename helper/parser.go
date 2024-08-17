package helper

import (
	"fmt"
	"strconv"
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

func ToUint(value string) (uint, error) {
	u, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		fmt.Println("Error to convert string to uint:", err)
		return 0, err
	}
	return uint(u), nil
}

func ToFloat(value string) (float64, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println("Error to convert string to float:", err)
		return 0, err
	}
	return floatValue, nil
}

func ToString(value float64) string {
	strValue := fmt.Sprintf("%.2f", value)
	return strValue
}
