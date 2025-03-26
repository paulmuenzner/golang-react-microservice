package date

import "time"

func GetCurrentUTCTime() time.Time {
	return time.Now().UTC()
}

func FormatDate(t time.Time) string {
	return t.Format("02. January 2006 15:04:05")
}
