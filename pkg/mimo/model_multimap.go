package mimo

type Multimap map[string]map[string]int

// Add a key/value pair to the multimap.
func (m Multimap) Add(key string, value string) {
	set, exists := m[key]
	if !exists {
		set = make(map[string]int)
	}

	set[value]++

	m[key] = set
}

// Count the number of values associated to key.
func (m Multimap) Count(key string) int {
	set, exists := m[key]
	if !exists {
		return 0
	}

	return len(set)
}

// Rate return the percentage of keys that have a count of 1.
func (m Multimap) Rate() float64 {
	cnt := 0

	for _, set := range m {
		if len(set) == 1 {
			cnt++
		}
	}

	return float64(cnt) / float64(len(m))
}
