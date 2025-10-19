package date

import "time"

func GetCurrentUTCTime() time.Time {
	return time.Now().UTC()
}

func FormatDate(t time.Time) string {
	return t.Format("02. January 2006 15:04:05")
}

// Returns: A string formatted according to the provided layout.
func FormatCurrentUTCTime(formatLayout string) string {
	// 1. Get the current time in UTC using the dedicated function GetCurrentUTCTime().
	nowUTC := GetCurrentUTCTime()

	// 2. Format the time using the defined layout passed as an argument.
	return nowUTC.Format(formatLayout)
}
