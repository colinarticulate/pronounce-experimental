//go:build !debug
// +build !debug

package pron

import (
	"os"
)

func _debug(args ...interface{}) {
}

func removeFromDisk(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		// What to do here?
	}
}
