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

func (this *ArgParser) AddArg(key string, val string, des string) *ArgParser {
	this.args[key] = flag.String(key, val, des)
	return this
}

func (this *ArgParser) Parse() *ArgParser {
	flag.Parse()
	return this
}

func (this *ArgParser) Get(key string) (string, bool) {
	if !flag.Parsed() {
		flag.Parse()
	}

	val, exist := this.args[key]
	if !exist {
		return "", false
	}
	return *val, true
}

func (this *ArgParser) GetInt(key string) (int, bool) {
	strVal, exist := this.Get(key)
	if !exist {
		return 0, false
	}
	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, false
	}
	return val, true
}
