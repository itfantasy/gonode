package logger

import (
	"fmt"

	"github.com/itfantasy/gonode/utils/io"
	log "github.com/jeanphorn/log4go"
)

var configed bool = false

func LoadConfig(filePath string) {
	if filePath != "" {
		b, _ := io.FileExists(filePath)
		if b {
			log.LoadConfiguration(filePath)
		} else {
			autoConfig()
		}
	} else {
		autoConfig()
	}
	configed = true
}

func autoConfig() {
	fmt.Println("[log4]::begin using a default config...")
	txt := `{
    "console": {
        "enable": true,
        "level": "DEBUG"
    },  
    "files": [{
        "enable": true,
        "level": "INFO",
        "filename":"./test.log",
        "category": "Test",
        "pattern": "[%D %T] [%C] [%L] (%S) %M"
    },{ 
        "enable": false,
        "level": "DEBUG",
        "filename":"rotate_test.log",
        "category": "TestRotate",
        "pattern": "[%D %T] [%C] [%L] (%S) %M",
        "rotate": true,
        "maxsize": "500M",
        "maxlines": "10K",
        "daily": true,
        "sanitize": true
    }], 
    "sockets": [{
        "enable": false,
        "level": "DEBUG",
        "category": "TestSocket",
        "pattern": "[%D %T] [%C] [%L] (%S) %M",
        "addr": "127.0.0.1:12124",
        "protocol":"udp"
    }]
}`

	filePath := "./logger.json"
	io.SaveFile(filePath, txt)
	log.LoadConfiguration(filePath)
}

func NewLogger(filter string) *log.Filter {
	if !configed {
		autoConfig()
	}

	return log.LOGGER(filter)
}
