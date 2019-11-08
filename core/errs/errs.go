package errs

import (
	"errors"
	"strconv"
	"strings"
)

func New(errcode int, errmsg string) error {
	if errcode == 0 {
		return errors.New(errmsg)
	}
	return errors.New(strconv.Itoa(errcode) + "##" + errmsg)
}
func Info(err error) (int, string) {
	infos := strings.Split(err.Error(), "##")
	if len(infos) != 2 {
		return 0, err.Error()
	}
	i, err := strconv.Atoi(infos[0])
	if err != nil {
		return 0, err.Error()
	}
	return i, infos[1]
}
