package util

import (
	"errors"
	"fmt"
)

func LogInfo(log ...interface{}) {
	fmt.Println(ToStringTrim(log, ""))
}

func LogErr(log ...interface{}) error {
	return errors.New(ToStringTrim(log, ""))
}
