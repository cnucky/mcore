package dates

import (
	"time"
	"errors"
	"github.com/xsnews/mcore/log"
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
			return nil, errors.New("Broken for-loop")
		}
		if initialMonth != t.Month() {
			// And we're done
			log.Debug("Finished opening %d files", i)
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