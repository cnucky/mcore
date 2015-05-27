# METRICS

# Copy/Paste
```go
import (
  "log"
	"github.com/xsnews/mcore/metrics"
)

...

  // Init
  l := *log.New(os.Stderr, "", log.Ldate|log.Ltime)
  if err := metrics.Start(fmt.Sprintf("service.%s", hostname), GraphiteHostPort, VerboseFlag, l); err != nil {
    panic(err)
  }

  // Counter
	if err := metrics.AddCounter(fmt.Sprintf("mail.%s.sent", mailbox), 1); err != nil {
		fmt.Println("WARN: Failed flushing to graphite: " + err.Error())
	}

```
