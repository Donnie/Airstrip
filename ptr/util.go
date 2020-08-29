package ptr

import "time"

// Bool returns a pointer to the specified bool value.
func Bool(b bool) (p *bool) {
	p = &b
	return
}

// Int64 returns a pointer to the specified int value.
func Int64(i int64) (p *int64) {
	p = &i
	return
}

// Int returns a pointer to the specified int value.
func Int(i int) (p *int) {
	p = &i
	return
}

// Float returns a pointer to the specified float value.
func Float(f float64) (p *float64) {
	p = &f
	return
}

// String returns a pointer to the specified string value.
func String(s string) (p *string) {
	p = &s
	return
}

// Time returns a pointer to the specified time.Time value.
func Time(t time.Time) (p *time.Time) {
	p = &t
	return
}
