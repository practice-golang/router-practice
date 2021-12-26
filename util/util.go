package util

import "os"

func CheckFileExists(path string) (result bool) {
	_, err := os.Stat(path)

	if err == nil {
		result = true
	} else {
		result = false
	}

	return
}
