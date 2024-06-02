package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	c := NewCache()

	key := "abc"
	err := c.Set(key, 123)
	assert.NoError(t, err, nil)
	val, _ := c.Get(key)
	assert.Equal(t, val, 123)
}
