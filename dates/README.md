Dates pattern lib
=================
Simple date 'pattern' parser that extends time.ParseDuration with additional units d, M, y.

```go
base := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)
out, e := ParseDuration("-1y5d6h", base)
if e != nil {
	panic(e)
}
fmt.Println(out.String())
// 2003-12-26 18:00:00 +0000 UTC
```

Patterns
=================
```
y - year
M - month
d - day
h - hour
m - minute
s - second
```