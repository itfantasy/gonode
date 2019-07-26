package datacenter

var outOfDates map[string]int

func init() {
	outOfDates = make(map[string]int)
}

func clearOutOfDate(id string) {
	_, exist := outOfDates[id]
	if exist {
		delete(outOfDates, id)
	}
}

func checkOutOfDate(id string) bool {
	_, exist := outOfDates[id]
	if !exist {
		outOfDates[id] = 0
	}
	outOfDates[id]++
	num, _ := outOfDates[id]
	if num >= 3 {
		return true
	}
	return false
}
