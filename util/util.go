package util

import (
	"io/fs"
	"os"

	"github.com/practice-golang/router-practice/router"
)

func CheckFileExists(path string, isEmbed bool) (result bool) {
	switch isEmbed {
	case true:
		ef, err := fs.Stat(router.Content, path)
		if err == nil && ef != nil {
			result = true
		} else {
			result = false
		}
		break
	case false:
		f, err := os.Stat(path)
		if err == nil && f != nil {
			result = true
		} else {
			result = false
		}
		break
	}

	return
}
