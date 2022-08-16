package utils

import (
	"errors"
	"time"
)

func ToTrueTime(date string) (time.Time, error) {
	timeString := "20060102"
	loc, _ := time.LoadLocation("Local")
	if len(date) == 8 {
		return time.ParseInLocation(timeString, string(date), loc)
	}
	if len(date) == 6 {
		return time.ParseInLocation(timeString[:6], string(date), loc)
	}
	if len(date) == 4 {
		return time.ParseInLocation(timeString[:4], string(date), loc)
	}
	return time.Time{}, errors.New("time transcode error")
}

func UnixFormat(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02T15:04:05Z")
}
