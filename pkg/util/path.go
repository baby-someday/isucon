package util

import "os"

func PWD() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
