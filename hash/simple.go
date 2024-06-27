package hash

import "hash/fnv"

type HashManager struct {
}

func (m *HashManager) GetInstanceFromKey(key string, numberOfInstances int) int {
	hash := m.HashId(key)
	return int(hash) % int(numberOfInstances)
}

func (m *HashManager) HashId(id string) int {
	h := fnv.New32a()
	h.Write([]byte(id))
	return int(h.Sum32())
}
