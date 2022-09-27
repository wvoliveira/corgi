package group

import (
	"errors"
	"regexp"
)

func checkName(name string) (err error) {
	err = errors.New("the name is required and must be alphanumeric, lowercase and 3-30 length")

	if name == "" || len(name) < 3 || len(name) > 30 {
		return
	}

	ok, err := regexp.Match("^[a-z0-9-]+$", []byte(name))
	if !ok || err != nil {
		return
	}

	return nil
}
