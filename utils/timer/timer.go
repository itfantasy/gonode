package timer

import (
	"time"
)

func Sleep(ms int) {
	<-time.After(time.Millisecond * time.Duration(ms))
}
