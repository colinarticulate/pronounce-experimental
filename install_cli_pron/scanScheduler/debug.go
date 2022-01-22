// +build debug 

package scanScheduler

import (
  "fmt"
)

func debug(args ...interface{}) {
  fmt.Println(args...)
}