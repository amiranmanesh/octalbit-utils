package xkafka

import (
	"fmt"
	"github.com/amiranmanesh/octalbit-utils/socket"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestKafka(t *testing.T) {
	broker, err := GetMessageBroker(Config{
		Servers:       "127.0.0.1:9092",
		GroupID:       "myGroup",
		ProduceTopic:  "myTopic",
		ConsumeTopics: []string{"myTopic"},
	})
	assert.NoError(t, err)

	go func() {
		err = broker.Consume(kafkaConsumer)
		assert.NoError(t, err)
	}()

	time.Sleep(2 * time.Second)

	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		err = broker.Publish([]byte(word))
		assert.NoError(t, err)
		time.Sleep(1 * time.Second)
	}
}

func kafkaConsumer(itemHook socket.MessageBrokerItem) error {
	if itemHook.Body == nil {
		return fmt.Errorf("failed to get item hook body")
	}
	fmt.Println("msg received: ", string(itemHook.Body))
	return nil
}
