package hash

import "hash/fnv"

func GetInstanceFromKey(key string, numberOfInstances int) int {
	hash := HashId(key)
	return int(hash) % int(numberOfInstances)
}

func HashId(id string) int {
	h := fnv.New32a()
	h.Write([]byte(id))
	return int(h.Sum32())
}
