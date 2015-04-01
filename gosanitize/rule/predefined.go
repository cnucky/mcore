package rule

import (
	"fmt"
)

func Logger(obj interface{}) error {
	fmt.Println(obj)

	return nil
}
