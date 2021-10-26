package rabbit

import (
	"encoding/json"
	"git.aasal.co/octal/octalbit/backend/utils/queues"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	queue, err := NewQueue(
		"localhost",
		"5672",
		"guest",
		"guest",
	)
	assert.Equal(t, nil, err)

	err = queue.Create("ex1", "direct", "q1", "b1")
	assert.Equal(t, nil, err)

	err = queue.Consume("q1", true, func(item queues.Item) {
		var tempData string
		assert.Equal(t, nil, json.Unmarshal(item.Body, &tempData))
		t.Log("From Queue 1: " + tempData)
	})
	assert.Equal(t, nil, err)

	body1, _ := json.Marshal("msg 1")
	err = queue.Publish("ex1", "b1", body1)
	assert.Equal(t, nil, err)

	time.Sleep(2 * time.Second)

	body2, _ := json.Marshal("msg 2")
	err = queue.Publish("ex1", "b1", body2)
	assert.Equal(t, nil, err)

	time.Sleep(2 * time.Second)

	body3, _ := json.Marshal("msg 3")
	err = queue.Publish("ex2", "b1", body3)
	assert.Equal(t, nil, err)

	sigChan := make(chan os.Signal, 1)
	<-sigChan

}
