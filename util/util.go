package util

import (
	"os"
)

func CheckFileExists(path string) (result bool) {
	f, err := os.Stat(path)
	if err == nil && f != nil {
		result = true
	} else {
		result = false
	}

	return
}
