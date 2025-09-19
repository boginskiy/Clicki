package repository

import "time"

func convertTimeToStr(tm time.Time, pattern string) string {
	return tm.Format(pattern)
}

func convertStrToTime(tm string, pattern string) (time.Time, error) {
	return time.Parse(pattern, tm)
}
