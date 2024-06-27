package hash

import "hash/fnv"

func GetInstanceFromKey(key string, numberOfInstances int) int {
	hash := hashId(key)
	return int(hash) % int(numberOfInstances)
}

func hashId(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}
