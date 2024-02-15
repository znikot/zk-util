package misc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffleSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ShuffleSlice(&slice)

	PrintJSON(slice, false)
}

func TestSearchStrings(t *testing.T) {
	strs := []string{"abc", "b", "hello", "world"}

	assert.Equal(t, ExistsString("hello", true, strs...), true)
	assert.Equal(t, ExistsString("hello1", false, strs...), false)
	assert.Equal(t, ExistsString("worlds", false, strs...), false)
}

func TestSearchInts(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 18, 9, 10}

	assert.Equal(t, ExistsInt(18, false, ints...), true)
	assert.Equal(t, ExistsInt(181, false, ints...), false)
	assert.Equal(t, ExistsInt(10, false, ints...), true)
}

func TestSearchInt64s(t *testing.T) {
	ints := []int64{1, 2, 3, 4, 5, 780, 25, 1123, 765}

	assert.Equal(t, ExistsInt64(25, false, ints...), true)
	assert.Equal(t, ExistsInt64(251, false, ints...), false)
	assert.Equal(t, ExistsInt64(765, false, ints...), true)
	assert.Equal(t, ExistsInt64(7651, false, ints...), false)
	assert.Equal(t, ExistsInt64(1123, false, ints...), true)
}
