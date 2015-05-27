# METRICS

# Copy/Paste
```go
import (
	"github.com/xsnews/mcore/metrics"
)

...

  // Init
  if err := metrics.Start(fmt.Sprintf("ntd-mon.%s", hostname), config.Pref.Graphite, config.Verbose, l); err != nil {
    panic(err)
  }

  // Counter
	if err := metrics.AddCounter(fmt.Sprintf("mail.%s.sent", mailbox), 1); err != nil {
		fmt.Println("WARN: Failed flushing to graphite: " + e.Error())
	}

```
