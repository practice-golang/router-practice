package util

import (
	"io/fs"
	"os"

	"router-practice/router"
)

func CheckFileExists(path string, isEmbed bool) (result bool) {
	result = false

	switch isEmbed {
	case true:
		ef, err := fs.Stat(router.Content, path)
		if err == nil && ef != nil {
			result = true
		}
	case false:
		f, err := os.Stat(path)
		if err == nil && f != nil {
			result = true
		}
	}

	return
}
