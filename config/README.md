Config loader
================
Remove duplicate code from multiple projects into one.

```go
package config

import (
	"github.com/xsnews/microservice-core/config"
)

type Config struct {
	// TODO: body here
}

var C Config

func Init(filename string) error {
	return config.Load(filename, &C)
}
```
