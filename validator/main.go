package validator

import (
	"slices"
	"strconv"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

// Valid Intervals: (D)ays, (M)onths,(H)ours
var validIntervalUnits = []string{"D", "M", "H"}

// Interval defaults to 30
// Unit defaults to D
func (v *Validator) TimeRange(interval, unit string) (int, string) {
	var i int
	var u string
	// default to 30 days
	if interval == "" || unit == "" {
		return 30, "D"
	}

	i, err := strconv.Atoi(interval)
	if err != nil {
		i = 30
	}

	if !slices.Contains(validIntervalUnits, unit) {
		u = "D"
	}

	return i, u
}
