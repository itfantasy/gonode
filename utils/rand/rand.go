package rand

import (
	"math/rand"
	"time"
)

func Random(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max+1-min)
}
