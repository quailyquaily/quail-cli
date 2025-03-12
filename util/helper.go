package util

import "time"

func AnyDatetimeToRFC3339(datetimeStr string) string {
	if datetimeStr == "" {
		return time.Now().Format(time.RFC3339)
	}
	// try to parse the datetimeStr in common formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02 15",
		"2006-01-02",
	}
	for _, format := range formats {
		tt, err := time.Parse(format, datetimeStr)
		if err == nil {
			return tt.Format(time.RFC3339)
		}
	}
	return time.Now().Format(time.RFC3339)
}
