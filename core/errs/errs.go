package errs

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

func CustomError(errcode int, errmsg string) error {
	if errcode == 0 {
		return errors.New(errmsg)
	}
	return errors.New(strconv.Itoa(errcode) + "##" + errmsg)
}

func ErrorInfo(err error) (int, string) {
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

type ErrorDigester interface {
	OnDigestError(interface{})
}

var _errDigester ErrorDigester

func BindDigester(errDigester ErrorDigester) {
	_errDigester = errDigester
}

func AutoRecover() {
	if err := recover(); err != nil {
		if _errDigester != nil {
			_errDigester.OnDigestError(err)
		} else {
			content := "!!! Auto Recovering...  " + fmt.Sprint(err) +
				"\r=============== - CallStackInfo - =============== \r" + string(debug.Stack())
			fmt.Println(content)
		}
	}
}
