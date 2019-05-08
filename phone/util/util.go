package util

import (
	"regexp"
)

func Normalize(phoneNumber string) (string, error) {

	rgx, err := regexp.Compile(`\D+`)

	if err != nil {
		return "", err
	}

	ret := rgx.ReplaceAll([]byte(phoneNumber), []byte(""))

	return string(ret), nil

}
