package calc

import (
	"github.com/itshosted/mcore/sort"
	"testing"
)

// Using a.b.c.d keys to re-use sortedmap
func TestNineFifth(t *testing.T) {
	dayOne := sort.NewValSorter(map[string]int64{
		"a": 100, "b": 30, "c": 33, "d": 3000, "e": 301, "f": 222, "g": 25,
		"h": 88, "i": 99, "j": 100, "k": 123, "l": 100, "m": 123, "n": 104,
		"o": 22, "p": 300, "q": 102, "r": 178, "s": 109, "t": 200,
	})
	dayOne.Sort()

	dayTwo := sort.NewValSorter(map[string]int64{
		"a": 200, "b": 90, "c": 22, "d": 11, "e": 89, "f": 6700, "g": 1200,
		"h": 22, "i": 4, "j": 100, "k": 22, "l": 200, "m": 304, "n": 166,
		"o": 106, "p": 1790, "q": 162, "r": 185, "s": 188, "t": 2000,
	})
	dayTwo.Sort()

	n := NewNineFifth()
	one := n.Add(dayOne.Vals, false)
	if one != 301 {
		t.Errorf("dayOne 95th value wrong, should be %d, received=%d", 301, one)
	}
	two := n.Add(dayTwo.Vals, true)
	if two != 2000 {
		t.Errorf("dayTwo 95th value wrong, should be %d, received=%d", 2000, two)
	}

	total := n.Total95th()
	if total != 2000 {
		t.Errorf("total 95th value wrong, should be %d, received=%d", 2000, total)
	}
}
