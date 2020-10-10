package main

import (
	"strings"
	"time"

	"github.com/jinzhu/now"
)

func parseDate(input string) (out time.Time) {
	input = strings.ToLower(input)
	now.TimeFormats = append(now.TimeFormats, "2 Jan")
	now.TimeFormats = append(now.TimeFormats, "Jan 2")

	switch input {
	case "now":
		out = time.Now()
	case "today":
		out = now.BeginningOfDay()
	case "yday", "y'day", "yesterday":
		out = now.BeginningOfDay().AddDate(0, 0, -1)
	case "tomm", "tomorrow":
		out = now.BeginningOfDay().AddDate(0, 0, 1)
	default:
		out, _ = now.Parse(input)
	}
	return
}
