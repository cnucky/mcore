package dates

import (
	"testing"
	"time"
)

type cmp struct {
	pattern string
	compare string
	error bool
}

func TestParseDuration(t *testing.T) {
	base := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []cmp{
		cmp{"1y", "2006-01-01 00:00:00 +0000 UTC", false},
		cmp{"-1y", "2004-01-01 00:00:00 +0000 UTC", false},
		cmp{"1y5d6h", "2006-01-06 06:00:00 +0000 UTC", false},
		cmp{"-1y5d6h", "2003-12-26 18:00:00 +0000 UTC", false},
		cmp{"88Z", "", true},
		cmp{"88", "", true},
	}

	for _, task := range tests {
		out, e := ParseDuration(task.pattern, base)
		if e != nil && task.error {
			// We got an error like expected
			continue
		}

		if e != nil {
			t.Errorf("Unexpected error=" + e.Error())
			continue
		}

		if out.String() != task.compare {
			t.Errorf("ParseDuration broken. Expected=[" + task.compare + "] Received=[" + out.String() + "]")
		}
	}
}