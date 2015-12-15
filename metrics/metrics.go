// Metric library abstraction.
// Abstract away counter and gauge implementations
// for creating statistics.
// @author mdroog/tomarus
package metrics

import (
	"github.com/itshosted/go-metrics" // diff rcrowley/go-metrics
	"net"
	"time"
	"errors"
	"fmt"
	"log"
)

/*
  WARN: Library is updated, it returns errors..
*/

var verbose bool
var enabled bool
var stopper chan bool
var counters map[string]metrics.Counter
var gauges map[string]metrics.Gauge
var gaugesf map[string]metrics.GaugeFloat64
var l log.Logger

// Prepare memory
func init() {
	counters = make(map[string]metrics.Counter)
	gauges = make(map[string]metrics.Gauge)
	gaugesf = make(map[string]metrics.GaugeFloat64)
	stopper = make(chan bool)
}

func Health() *string {
	if !enabled {
		msg := "metrics - thread disabled"
		return &msg
	}
	return nil
}

// Resolve DNS and start go-routine
func Start(tag, ip string, isverbose bool, logger log.Logger) error {
	if enabled {
		return errors.New("Already started")
	}
	verbose = isverbose
	l = logger
	addr, e := net.ResolveTCPAddr("tcp", ip)
	if e != nil {
		return e
	}
	go updater(metrics.DefaultRegistry, tag, addr, stopper)
	enabled = true
	return nil
}

// Stop Go-routine
func Stop() {
	stopper <- true
}

// Period (every minute) update Graphite
func updater(r metrics.Registry, prefix string, addr *net.TCPAddr, exit chan bool) {
	for {
		select {
		case <-time.After(time.Second * 60):
			if e := metrics.GraphiteOnce(r, prefix, addr); e != nil {
				if verbose {
					fmt.Println("Lost graphite stats: " + e.Error())
				}
				l.Printf("Lost Graphite stats: " + e.Error())
			}
		case <-exit:
			return
		}
	}
	enabled = false
}

// key += value for this minute
func AddCounter(key string, value int64) error {
	if !enabled {
		panic(errors.New("Not connected to Graphite"))
	}

	if _, x := counters[key]; !x {
		counters[key] = metrics.NewCounter()
		if err := metrics.Register(key, counters[key]); err != nil {
			return err
		}
	}

	if verbose {
		fmt.Println(fmt.Sprintf("Counter %s=%d", key, value))
	}
	counters[key].Inc(value)
	return nil
}

// key=value for this minute
func SetGauge(key string, value int64) error {
	if !enabled {
		panic(errors.New("Not connected to Graphite"))
	}

	if _, x := gauges[key]; !x {
		gauges[key] = metrics.NewGauge()
		if err := metrics.Register(key, gauges[key]); err != nil {
			return err
		}
	}

	if verbose {
		fmt.Println(fmt.Sprintf("Gauge %s=%d", key, value))
	}
	gauges[key].Update(value)
	return nil
}

// key=value for this minute
func SetGaugeFloat(key string, value float64) error {
	if !enabled {
		panic(errors.New("Not connected to Graphite"))
	}

	if _, x := gauges[key]; !x {
		gaugesf[key] = metrics.NewGaugeFloat64()
		if err := metrics.Register(key, gaugesf[key]); err != nil {
			return err
		}
	}

	if verbose {
		fmt.Println(fmt.Sprintf("GaugeFloat %s=%d", key, value))
	}
	gaugesf[key].Update(value)
	return nil
}
