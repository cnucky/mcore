package dates

import (
	"time"
	"errors"
)

const MAXITER = 50

// Determine future days in same month from given date.
func DaysInMonth(t time.Time) ([]time.Time, error) {
	var (
		output []time.Time
		e error
	)

	i := 0
	initialMonth := t.Month()
	for {
		i++
		if i == MAXITER {
			return nil, errors.New("DaysInMonth: Broken for-loop")
		}
		if initialMonth != t.Month() {
			// And we're done
			break
		}
		output = append(output, t)

		t, e = ParseDuration("1d", t)
		if e != nil {
			return nil, e
		}
	}

	return output, nil
}

func MonthsInYear(t time.Time) ([]time.Time, error) {
	var (
		output []time.Time
		e error
	)

	i := 0
	initialYear := t.Year()
	for {
		i++
		if i == MAXITER {
			return nil, errors.New("MonthsInYear: Broken for-loop")
		}
		if initialYear != t.Year() {
			// And we're done
			break
		}
		output = append(output, t)

		t, e = ParseDuration("1M", t)
		if e != nil {
			return nil, e
		}
	}

	return output, nil	
}