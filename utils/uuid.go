package utils

import (
	"strconv"
	"time"
)

func UUID() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
