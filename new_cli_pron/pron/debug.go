// +build debug

package pron

import (
	"fmt"
)

func debug(args ...interface{}) {
	fmt.Println(args...)
}

func removeFromDisk(filepath string) {
}
