package util

import (
	"os"
	"path/filepath"
	"router-practice/router"
)

func CheckFileExists(path string, isEmbed bool) (result bool) {
	switch isEmbed {
	case true:
		fname := filepath.Base(path)
		dir, _ := router.Content.ReadDir(filepath.Dir(path))
		for _, f := range dir {
			if f.Name() == fname {
				result = true
				break
			}
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
