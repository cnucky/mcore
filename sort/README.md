Sorting
===================

Utility methods for easy sorting.

Int64 sort
```
import "github.com/xsnews/mcore/sort"
...
vals := []int64{16,8,4}
sort.Ints(vals)
fmt.Println(vals)
```

map[string]int64 sort
```
import "github.com/xsnews/mcore/sort"
...
vals := make(map[string]int64{"A": 1, "B": 12, "C": 3300}
vs := sort.NewValSorter(vals)
vs.Sort()
fmt.Println(vs.Vals)
```