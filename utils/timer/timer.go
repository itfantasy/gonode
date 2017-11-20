package timer

import (
	"time"
)

func Sleep(ms time.Duration) {
	<-time.After(time.Millisecond * ms)
}
