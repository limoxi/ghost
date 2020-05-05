package utils

import "time"

const DEFAULT_TIME_LAYOUT = "2006-01-02 15:04:05"
const DEFAULT_MINUTE_TIME_LAYOUT = "2006-01-02 15:04"
const DEFAULT_DATE_LAYOUT = "2006-01-02"

var DEFAULT_TIME, _ = time.ParseInLocation(DEFAULT_TIME_LAYOUT, "1900-01-01 00:00:00", time.Local)

func FormatDatetime(dt time.Time, args ...string) string{
	layout := DEFAULT_TIME_LAYOUT
	switch len(args) {
	case 1:
		layout = args[0]
	}

	return dt.Format(layout)
}