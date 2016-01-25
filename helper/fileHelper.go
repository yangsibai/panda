package helper

import (
	"os"
)

func CreateDirIfNotExists(dir string) error {
	exists, err := IsExists(dir)
	if err != nil {
		return err
	}
	if !exists {
		return os.Mkdir(dir, 0755)
	}
	return nil
}

func IsExists(filename string) (exists bool, err error) {
	_, err = os.Stat(filename)
	if err == nil {
		exists = true
		return
	}
	if os.IsNotExist(err) {
		exists = false
		err = nil
		return
	}
	return false, err
}
