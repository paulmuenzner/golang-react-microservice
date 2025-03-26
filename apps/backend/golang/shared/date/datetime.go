package date

import "time"

func GetCurrentUTCTime() time.Time {
	return time.Now().UTC()
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
