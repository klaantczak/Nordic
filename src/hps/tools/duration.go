package tools

import (
	"fmt"
	"strconv"
	"strings"
)

func parseDigits(str string, i int) (int, bool) {
	j := i
loop:
	for j < len(str) {
		switch str[j] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			j++
			continue
		default:
			break loop
		}
	}
	return j, j != i
}

func parsePoint(str string, i int) (int, bool) {
	if i < len(str) && str[i] == '.' {
		return i + 1, true
	}
	return i, false
}

func parseNumber(str string, i int) (float64, int, bool) {
	var err error

	v := 0.0
	j := i
	ok := false

	j, ok = parseDigits(str, j)
	if ok {
		j, ok = parsePoint(str, j)
		if ok {
			j, ok = parseDigits(str, j)
		}
	}

	if i == j {
		return 0, i, false
	}

	if v, err = strconv.ParseFloat(str[i:j], 64); err == nil {
		return v, j, true
	}

	return 0, i, false
}

// Returns the duration of the interval, in years.
func ParseDuration(duration string) (float64, error) {
	str := strings.TrimSpace(duration)

	v, p, ok := parseNumber(str, 0)
	if ok {
		if p != len(str) {
			if f, err := frequencyPerYear(strings.TrimSpace(str[p:])); err == nil {
				return v / f, nil
			}
			return 0.0, fmt.Errorf("Incorrect duration")
		}
		return v, nil
	}

	if f, err := frequencyPerYear(str); err == nil {
		return 1.0 / f, nil
	}

	return 0.0, fmt.Errorf("Incorrect duration")
}

// Returns number of units in one year.
func frequencyPerYear(units string) (float64, error) {
	value := 0
	switch strings.ToUpper(units) {
	case "S", "SEC", "SECOND", "SECONDS":
		value = 60 * 60 * 24 * 365
	case "M", "MIN", "MINS", "MINUTE", "MINUTES":
		value = 60 * 24 * 365
	case "H", "HOUR", "HOURS":
		value = 24 * 365
	case "D", "DAY", "DAYS":
		value = 365
	case "W", "WEEK", "WEEKS":
		value = 52
	case "MONTH", "MONTHS":
		value = 12
	case "Y", "YEAR", "YEARS":
		value = 1
	}
	if value == 0 {
		return 0.0, fmt.Errorf("Unsupported time unit.")
	}
	return float64(value), nil
}
