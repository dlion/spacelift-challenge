package hash

import "hash/fnv"

func hashId(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}

func GetInstanceFromKey(key string, numberOfInstances int) int {
	hash := hashId(key)
	return int(hash) % int(numberOfInstances)
}
