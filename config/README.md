Config loader
================
Remove duplicate code from multiple projects into one.

```go
package config

import (
	"github.com/xsnews/mcore/config"
)

type Config struct {
	// TODO: body here
}

var C Config

func Init(filename string) error {
	return config.Load(filename, &C)
}
```

Load files from a conf.d

```go
package main

import (
  "github.com/xsnews/mcore/config"
)

type Host struct {
  Ip   string
  Port int
}

func main() {
  hosts := make(map[string]Host)

  err := config.LoadJsonD("./conf.d/", &hosts)
  if err != nil {
          panic(err)
  }
}
```
