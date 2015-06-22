# filesmatch - Recursive read dir by glob pattern
http://golang.org/pkg/path/filepath/#Match

```go
package main

import (
        "fmt"
        "github.com/xsnews/mcore/filesmatch"
)

func main() {
        f, _ := filesmatch.Match("/data01/*")
        fmt.Println(f)
}
```

Example output
```
# ./t
map[golang:/data01/golang var:/data01/var]

```
