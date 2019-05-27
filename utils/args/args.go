package args

import (
	"flag"
	"strconv"
)

type ArgParser struct {
	args map[string]*string
}

func Parser() *ArgParser {
	parser := new(ArgParser)
	parser.args = make(map[string]*string)
	return parser
}

func (a *ArgParser) AddArg(key string, val string, des string) *ArgParser {
	a.args[key] = flag.String(key, val, des)
	return a
}

func (a *ArgParser) Parse() *ArgParser {
	flag.Parse()
	return a
}

func (a *ArgParser) Get(key string) (string, bool) {
	if !flag.Parsed() {
		flag.Parse()
	}

	val, exist := a.args[key]
	if !exist {
		return "", false
	}
	return *val, true
}

func (a *ArgParser) GetInt(key string) (int, bool) {
	strVal, exist := a.Get(key)
	if !exist {
		return 0, false
	}
	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, false
	}
	return val, true
}
