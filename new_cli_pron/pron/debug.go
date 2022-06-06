//go:build debug
// +build debug

package pron

import (
	"fmt"
)

func _debug(args ...interface{}) {
	fmt.Println(args...)
}

func removeFromDisk(filepath string) {
}
