package util

import (
	"fmt"
	"time"
)

func ValidateDate(input string) (time.Time, error) {
	var resp time.Time
	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, input)
	if err != nil {
		return resp, fmt.Errorf("invalid date format, must be YYYY-MM-DD")
	}
	return parsedDate, nil
}
