package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistent(t *testing.T) {
	t.Run("Given a certain ID, it should return the correspondent hash", func(t *testing.T) {
		assert.Equal(t, hashId("tYxdBGXBdjPx4Bus7Nsbvya99JWCreyR"), uint32(1082065427))
	})

	t.Run("Given a certain ID and minio instances, it should return the correspondent minio instance", func(t *testing.T) {
		assert.Equal(t, GetInstanceFromKey("tYxdBGXBdjPx4Bus7Nsbvya99JWCreyR", 3), 2)
	})
}
