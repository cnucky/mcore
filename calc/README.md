95th calc
=========================
Calculate 95th from a set of metrics for billing. Basically 95%
is taking all measurepoints, removing the highest 5% of metric points and
returning the highest value from the remaining measurepoints.

```
import "github.com/xsnews/mcore/calc"

n := NewNineFifth()
nineFifthDayOne := n.Add(dayOne, false)

nineFifthDayTwo := n.Add(dayTwo, true)

nineFifthTotal := n.Total95th()
```

> For a descent example please have a look at the unit test.