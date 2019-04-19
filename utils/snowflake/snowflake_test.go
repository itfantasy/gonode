package snowflake

import (
	"fmt"
)

func TestGenerate() {
	seed, err := NewSeed(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch := make(chan ID)
	count := 10000
	for i := 0; i < count; i++ {
		go func() {
			id := seed.Generate()
			ch <- id
		}()
	}
	defer close(ch)
	m := make(map[ID]int)
	for i := 0; i < count; i++ {
		id := <-ch
		_, ok := m[id]
		if ok {
			fmt.Printf("ID is not unique!\n")
			return
		}
		m[id] = i
	}
	fmt.Println("All ", count, " snowflake ID generate successed!\n")
}
