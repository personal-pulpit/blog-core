package utils

import "time"

func ParasTime(str string) (*time.Time, error) {
	layout := "2006-01-02 15:04:05.999999999-07"
	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}
