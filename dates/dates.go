package dates

import (
	"errors"
	"time"
	"strconv"
	"fmt"
)

// Extended time.ParseDuration
// Additional supported units ['d', 'm', 'y']
func ParseDuration(str string, initial time.Time) (time.Time, error) {
	now := initial
	if initial.Unix() < time.Date(2005, time.January, 1, 0, 0, 0, 0, now.Location()).Unix() {
		panic("DevErr: Please supply a valid now=" + now.String())
	}

	val := ""
	s := str
	
	minus := '0'
	if len(s) > 0 && str[0] == '-' {
		minus = '-'
		s = s[1:]
	}

	for len(s) > 0 {
		c := s[0]

		// [0-9]+
		if c >= '0' && c <= '9' {
			val += string(c)
		} else if c == 'y' || c == 'm' || c == 'd' || c == 'M' || c == 's' || c == 'h' {
			if val == "" {
				return now, errors.New("Duration: No value before unit, input=" + str)
			}
			if minus != '0' {
				val = string(minus) + val
			}
			i, e := strconv.ParseInt(val, 10, 64)
			if e != nil {
				return now, e
			}
			if i == 0 {
				return now, errors.New("Duration: Failed parsing val=" + val + ", input=" + str)
			}

			if c == 'd' {
				// day
				now = now.AddDate(0, 0, int(i))
			} else if c == 'M' {
				// month
				now = now.AddDate(0, int(i), 0)
			} else if c == 'y' {
				// year
				now = now.AddDate(int(i), 0, 0)
			} else {
				if c == 'm' && s[2] == 's' {
					// ms
					d, e := time.ParseDuration(val + "ms")
					if e != nil {
						return now, e
					}
					now = now.Add(d)

				} else if (c == 'm' && s[2] == 's') || c == 's' || c == 'm' || c == 'h' {
					d, e := time.ParseDuration(val + string(c))
					if e != nil {
						return now, e
					}
					now = now.Add(d)
				} else {
					return now, fmt.Errorf("Invalid unit=%c for input=%s", c, str)
				}
			}
			val = ""

		} else {
			return now, errors.New("Duration: Invalid char=" + string(c))
		}

		s = s[1:]
	}

	if now.String() == initial.String() {
		return now, fmt.Errorf("Date not changed for str=" + str)
	}
	return now, nil
}