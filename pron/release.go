// +build !debug

package pron

import (
	"os"
)

func debug(args ...interface{}) {
}

func removeFromDisk(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		// What to do here?
	}
}
