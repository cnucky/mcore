package log

import (
	"fmt"
)

var (
	Verbose bool
)

func Init(verbose bool) {
	Verbose = verbose
}

func Debug(format string, v ...interface{}) {
	if Verbose {
		fmt.Println(fmt.Sprintf(format, v...))
	}
}

func Println(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v...))
}
