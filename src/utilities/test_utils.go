package utilities

import (
	"strings"
	"time"
)

func getDateString(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func ConvertStringToDate(s string, subString string) time.Time {

	s = getDateString(s, subString)
	s = strings.ReplaceAll(s, ":", "")
	date, _ := time.Parse("150405.000.", s)
	return date

}
