package wtime

import (
	"errors"
	"time"
)

func ParseLayout(value string, layouts ...string) (t time.Time, err error) {
	return ParseLayoutInLocation(value, time.Local, layouts...)
}

func ParseLayoutInLocation(value string, loc *time.Location, layouts ...string) (t time.Time, err error) {
	for _, v := range layouts {
		t, err = time.ParseInLocation(v, value, loc)
		if err == nil {
			return
		}
	}
	if err == nil {
		err = errors.New("no layouts")
	}
	return
}
