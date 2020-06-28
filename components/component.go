package components

import (
	"errors"
	"strings"

	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/strs"
)

const (
	redis    string = "redis"
	mysql           = "mysql"
	mongodb         = "mongodb"
	rabbitmq        = "rabbitmq"
	kafka           = "kafka"
	nsq             = "nsq"
	etcd            = "etcd"
	email           = "email"

	urlParasError = "illegal url!! —— componenttype://usr:pass@url/host?op_key=op_val&op_key=op_val...."
)

type IComponent interface {
	Conn(string, string) error
	Close()
	SetAuthor(string, string)
	SetOption(string, interface{})
}

func NewComponent(url string) (IComponent, error) {
	tempInfos := strings.Split(url, "?")
	length := len(tempInfos)
	if length != 1 && length != 2 {
		return nil, errors.New(urlParasError + "[A]")
	}
	tempUrl := tempInfos[0]

	tempInfos2 := strings.Split(tempUrl, "://")
	if len(tempInfos2) != 2 {
		return nil, errors.New(urlParasError + "[B]")
	}

	compType := tempInfos2[0]
	tempUrl2 := strings.Split(tempInfos2[1], "@")
	tempUrl2Len := len(tempUrl2)
	if tempUrl2Len != 1 && tempUrl2Len != 2 {
		return nil, errors.New(urlParasError + "[C]")
	}

	usr := ""
	pass := ""
	bAuthor := false
	if tempUrl2Len == 2 {
		bAuthor = true
		tempInfo3 := tempUrl2[0]
		usrAndPass := strings.Split(tempInfo3, ":")
		length2 := len(usrAndPass)
		if length2 != 1 && length2 != 2 {
			return nil, errors.New(urlParasError + "[D]")
		}
		usr = usrAndPass[0]
		if length2 == 2 {
			pass = usrAndPass[1]
		}
	}

	tempUrl3 := tempUrl2[0]
	if bAuthor {
		tempUrl3 = tempUrl2[1]
	}
	urlAndHost := strings.Split(tempUrl3, "/")
	if len(urlAndHost) != 2 {
		return nil, errors.New(urlParasError + "[E]" + tempUrl3)
	}
	theurl := urlAndHost[0]
	host := urlAndHost[1]

	var comp IComponent = nil
	switch compType {
	case redis:
		comp = NewRedis()
	case mongodb:
		comp = NewMongoDB()
	case mysql:
		comp = NewMySql()
	case rabbitmq:
		comp = NewRabbitMQ()
	case etcd:
		comp = NewEtcd()
	case email:
		comp = NewEmail()
	}

	if comp == nil {
		return nil, errors.New("illegal component type!!")
	}

	comp.SetAuthor(usr, pass)

	if length == 2 {
		tempOpt := tempInfos[1]
		opts := strings.Split(tempOpt, "&")
		if len(opts) > 0 {
			for _, kv := range opts {
				keyAndValue := strings.Split(kv, "=")
				if len(keyAndValue) != 2 {
					return nil, errors.New("illegal component option!!" + kv)
				}
				opKey := keyAndValue[0]
				strVal := keyAndValue[1]
				if strs.StartsWith(strVal, "{") && strs.EndsWith(strVal, "}") {
					var opVal map[string]interface{}
					if err := json.Unmarshal(strVal, &opVal); err != nil {
						return nil, err
					}
					comp.SetOption(opKey, opVal)
				}
				comp.SetOption(opKey, strVal)
			}
		}
	}

	err := comp.Conn(theurl, host)
	if err != nil {
		return nil, err
	}
	return comp, nil
}
