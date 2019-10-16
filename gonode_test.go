package gonode

import (
	"testing"
	logger2 "github.com/itfantasy/gonode/core/logger"
	"errors"
)

func TestLogger(t *testing.T){
	logger, err := logger2.NewLogger("MyNode", "DEBUG", "", "")
	if err != nil{
		t.Error(err)
	}
	logger.Debug("AAA")
	logger.Debug(errors.New("A ERROR!"))
}