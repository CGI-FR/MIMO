package mimo_test

import (
	"testing"

	"github.com/adrienaury/mimo/pkg/mimo"
	"github.com/stretchr/testify/assert"
)

func TestMultimap(t *testing.T) {
	t.Parallel()

	multimap := mimo.Multimap{}

	multimap.Add("A", "X")
	multimap.Add("A", "Y")
	multimap.Add("B", "Z")
	multimap.Add("B", "Z")
	multimap.Add("C", "Z")
	multimap.Add("C", "Z")
	multimap.Add("D", "W")
	multimap.Add("D", "W")
	multimap.Add("E", "V")
	multimap.Add("E", "V")

	assert.Equal(t, 2, multimap.Count("A"))
	assert.Equal(t, 1, multimap.Count("B"))
	assert.Equal(t, 1, multimap.Count("C"))
	assert.Equal(t, 1, multimap.Count("D"))
	assert.Equal(t, 1, multimap.Count("E"))
	assert.Equal(t, 0, multimap.Count("F"))

	assert.Equal(t, 0.8, multimap.Rate())
}
