package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentile(t *testing.T) {
	n := NewPercentile()
	n.Append(1)
	n.Append(2)
	n.Append(3)
	n.Append(4)
	assert.Equal(t, 3, n.Percentile(50))
}
