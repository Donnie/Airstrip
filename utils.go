package main

import (
	"strings"
	"time"

	"github.com/jinzhu/now"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func parseDate(input string) (out time.Time) {
	input = strings.ToLower(input)
	now.TimeFormats = append(now.TimeFormats, "2 Jan 2006")
	now.TimeFormats = append(now.TimeFormats, "2 Jan")
	now.TimeFormats = append(now.TimeFormats, "Jan 2 2006")
	now.TimeFormats = append(now.TimeFormats, "Jan 2")

	switch input {
	case "now":
		out = time.Now()
	case "today":
		out = now.BeginningOfDay()
	case "yday", "y'day", "yesterday":
		out = now.BeginningOfDay().AddDate(0, 0, -1)
	case "tom", "tomorrow":
		out = now.BeginningOfDay().AddDate(0, 0, 1)
	default:
		out, _ = now.Parse(input)
	}
	return
}

func getMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getMonthLastDate(month time.Time) time.Time {
	currentYear, currentMonth, _ := month.Date()
	currentLocation := month.Location()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).AddDate(0, 1, -1)
}

func getLastMonthLastDate() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).AddDate(0, 0, -1)
}
