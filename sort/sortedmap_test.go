package sort

import (
	"testing"
)

func TestValSort(t *testing.T) {
	dayCount := make(map[string]int64)
	dayCount["2015-12-01X00:00"] += 12
	dayCount["2015-12-01X00:00"] += 12

	dayCount["2015-12-01X00:01"] += 12

	vs := NewValSorter(dayCount)
	vs.Sort()

	if len(vs.Vals) != 2 {
		t.Errorf("Value array not 2 positions, received=%v", vs.Vals)
	}
	if vs.Vals[0] != 12 {
		t.Errorf("First position should be 12, received=%d", vs.Vals[0])
	}
	if vs.Vals[1] != 24 {
		t.Errorf("Second position should be 24, received=%d", vs.Vals[1])
	}
}

func TestValSortBig(t *testing.T) {
	dayCount := make(map[string]int64)
	dayCount["2015-12-01X00:00"] += 12
	dayCount["2015-12-01X00:00"] += 24

	dayCount["2015-12-01X00:01"] += 108
	dayCount["2015-12-01X00:02"] += 96
	dayCount["2015-12-01X00:03"] += 84
	dayCount["2015-12-01X00:04"] += 72
	dayCount["2015-12-01X00:05"] += 60
	dayCount["2015-12-01X00:06"] += 48

	vs := NewValSorter(dayCount)
	vs.Sort()

	if len(vs.Vals) != 7 {
		t.Errorf("Value array not 7 positions, received=%v", vs.Vals)
	}
	if vs.Vals[0] != 36 {
		t.Errorf("First position should be 36, received=%d", vs.Vals[0])
	}

	last5 := vs.Vals[ len(vs.Vals)-5:len(vs.Vals) ]
	expected := []int64{60,72,84,96,108}
	for i, val := range last5 {
		if val != expected[i] {
			t.Errorf("Pos %d got %d but expected %d", i, val, expected[i])
		}
	}
}