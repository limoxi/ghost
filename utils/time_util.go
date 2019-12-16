package utils

import "time"

const DEFAULT_TIME_LAYOUT = "2006-01-02 15:04:05"

func FormatDatetime(dt time.Time, args ...string) string{
	layout := DEFAULT_TIME_LAYOUT
	switch len(args) {
	case 1:
		layout = args[0]
	}

	return dt.Format(layout)
}