package utils

import "time"

func FormatTimeFromUnix(unix int64) string {
	timeObj := time.Unix(unix, 0)
	return timeObj.Format("2006-01-02 15:04:05")
}
