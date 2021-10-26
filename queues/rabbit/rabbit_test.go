package rabbit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsExists(t *testing.T) {
	items := queuesInfo{}
	items = append(items, queueInfo{"ex1", "q1", "b1"})
	items = append(items, queueInfo{"ex2", "q1", "b1"})
	items = append(items, queueInfo{"ex3", "q1", "b1"})

	assert.Equal(t, 3, len(items))
	assert.Equal(t, true, items.IsExist("ex3", "q1", "b1"))
}

func TestRepetitiousCreate(t *testing.T) {
	rabbit := rabbit{}
	assert.Equal(t, nil, rabbit.Create("ex1", "q1", "b1"))
	assert.Equal(t, nil, rabbit.Create("ex2", "q1", "b1"))
	assert.NotEqual(t, nil, rabbit.Create("ex1", "q1", "b1"))
}
