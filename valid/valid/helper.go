package valid

import (
	"strconv"
)

func FnGetInt(i interface{}) (int64, error) {
	min, err := strconv.ParseInt(i.(string), 0, 0)
	if err != nil {
		return 0, err
	}

	return min, nil
}

func FnGetStr(i interface{}) (string, error) {
	return i.(string), nil
}

func FnGetStrSlice(i interface{}) ([]string, error) {
	return i.([]string), nil
}
