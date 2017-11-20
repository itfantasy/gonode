package ini

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type Conf struct {
	path string
	vals map[string]map[string]string
}

func Load(path string) (*Conf, error) {
	conf := new(Conf)
	conf.path = path

	vals, err := conf.loadVals()
	if err != nil {
		return nil, err
	}
	conf.vals = vals
	return conf, nil
}

func (conf *Conf) Get(section string, name string) string {
	vals, exist := conf.vals[section]
	if !exist {
		return ""
	}
	val, exist := vals[name]
	if !exist {
		return ""
	}
	return val
}

func (conf *Conf) GetInt(section string, name string, defVal int) int {
	strVal := conf.Get(section, name)
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return defVal
	} else {
		return intVal
	}
}

func (conf *Conf) loadVals() (map[string]map[string]string, error) {

	file, err := os.Open(conf.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]map[string]string
	data = make(map[string]map[string]string)
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case string(line[0]) == "#":
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
		}

	}

	return data, nil
}
