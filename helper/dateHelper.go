package helper

import (
	"time"
)

func GetFullDate() string {
	t := time.Now()
	return t.Format("20060102")
}
