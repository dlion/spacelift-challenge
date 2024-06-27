package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistent(t *testing.T) {
	t.Run("Given a certain ID, it should return the correspondent hash", func(t *testing.T) {
		hashManager := HashManager{}
		assert.Equal(t, hashManager.HashId("tYxdBGXBdjPx4Bus7Nsbvya99JWCreyR"), 1082065427)
	})

	t.Run("Given a certain ID and number of instances, it should return the correspondent instance", func(t *testing.T) {
		hashManager := HashManager{}
		assert.Equal(t, hashManager.GetInstanceFromKey("tYxdBGXBdjPx4Bus7Nsbvya99JWCreyR", 3), 2)
	})
}
