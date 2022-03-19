package gcounter

type GCounter struct {
	Id      string         `json:"id"`
	Counter map[string]int `json:"values"`
}

func Initial(id string) GCounter {
	return GCounter{
		Id: id,
		// Id:      os.Getenv("NODE"),
		Counter: make(map[string]int),
	}

}

func Value(counter GCounter) int {
	var result = 0
	for _, value := range counter.Counter {
		result = result + value
	}
	return result
}

func Inc(counter GCounter) GCounter {
	counter.Counter[counter.Id] += 1
	return counter
}

/*
Read all keys in b, get the max between those and add all missing keys
*/
func Merge(a GCounter, b GCounter) GCounter {
	for id, val := range b.Counter {
		v, ok := a.Counter[id]
		if ok {
			a.Counter[id] = max(v, val)
		} else {
			a.Counter[id] = val
		}
	}
	return a
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
