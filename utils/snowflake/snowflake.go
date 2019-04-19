package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	seedBits  uint8 = 10
	stepBits  uint8 = 12
	seedMax   int64 = -1 ^ (-1 << seedBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = seedBits + stepBits
	seedShift uint8 = stepBits
)

var epoch int64 = 1288834974657

type Seed struct {
	mu        sync.Mutex
	timestamp int64
	seed      int64
	step      int64
}

var _seed *Seed

func seed() *Seed {
	if _seed == nil {
		now := time.Now().UnixNano() / 1e6
		tick := now % 1024
		_seed, _ = NewSeed(tick)
	}
	return _seed
}

func Generate() string {
	return fmt.Sprint(GenerateRaw())
}

func GenerateRaw() int64 {
	return seed().Generate()
}

func (n *Seed) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if n.timestamp == now {
		n.step++
		if n.step > stepMax {
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		n.step = 0
	}
	n.timestamp = now
	result := int64((now-epoch)<<timeShift | (n.seed << seedShift) | (n.step))
	return result
}

func NewSeed(seed int64) (*Seed, error) {
	if seed < 0 || seed > seedMax {
		return nil, errors.New("Seed number must be between 0 and 1023")
	}
	return &Seed{
		timestamp: 0,
		seed:      seed,
		step:      0,
	}, nil
}
