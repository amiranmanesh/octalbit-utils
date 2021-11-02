package rabbit

import (
	"fmt"
	"github.com/amiranmanesh/octalbit-utils/socket"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRabbitPublisherAndConsumerDirect(t *testing.T) {

	Setup("guest", "guest", "localhost", 5672)
	msgBroker, err := GetMessageBroker(QueueConfig{
		QueueName:    "q1",
		ExChangeName: "ex1",
		BindingName:  "b1",
		RouteName:    "b1",
		ExChangeKind: DirectExchange,
	})
	assert.Equal(t, nil, err, "failed to get message broker", err)

	go func() {
		err = msgBroker.Consume(rabbitConsumer)
		assert.Equal(t, nil, err, "failed to consume", err)
	}()

	time.Sleep(5 * time.Second)

	err = msgBroker.Publish([]byte("hello world 1"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	time.Sleep(5 * time.Second)

	err = msgBroker.Publish([]byte("hello world 2"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	time.Sleep(5 * time.Second)

}

func rabbitConsumer(itemHook socket.MessageBrokerItem) error {
	if itemHook.Body == nil {
		return fmt.Errorf("failed to get item hook body")
	}
	fmt.Println("msg received: ", string(itemHook.Body))
	return nil
}

func TestRabbitClearQueue(t *testing.T) {

	Setup("guest", "guest", "localhost", 5672)
	msgBroker, err := GetMessageBroker(QueueConfig{
		QueueName:    "q2",
		ExChangeName: "ex2",
		BindingName:  "b2",
		RouteName:    "b2",
		ExChangeKind: DirectExchange,
	})
	assert.Equal(t, nil, err, "failed to get message broker", err)

	err = msgBroker.Publish([]byte("hello world 1"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	err = msgBroker.Publish([]byte("hello world 2"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	time.Sleep(2 * time.Second)

	err = msgBroker.Clear()
	assert.Equal(t, nil, err, "failed to Publish message", err)

	go func() {
		err = msgBroker.Consume(rabbitConsumer)
		assert.Equal(t, nil, err, "failed to consume", err)
	}()

	time.Sleep(2 * time.Second)

	err = msgBroker.Publish([]byte("hello world 3"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	time.Sleep(5 * time.Second)
}

func TestRabbitMultiPublisherAndConsumerDirect(t *testing.T) {

	Setup("guest", "guest", "localhost", 5672)

	msgBroker, err := GetMessageBroker(QueueConfig{
		QueueName:    "q1",
		ExChangeName: "ex1",
		BindingName:  "b1",
		RouteName:    "b1",
		ExChangeKind: DirectExchange,
	})
	assert.Equal(t, nil, err, "failed to get message broker", err)

	msgBroker2, err := GetMessageBroker(QueueConfig{
		QueueName:    "q3",
		ExChangeName: "ex3",
		BindingName:  "b3",
		RouteName:    "b3",
		ExChangeKind: DirectExchange,
	})
	assert.Equal(t, nil, err, "failed to get message broker", err)

	go func() {
		err = msgBroker.Consume(rabbitConsumer)
		assert.Equal(t, nil, err, "failed to consume", err)
	}()
	go func() {
		err = msgBroker2.Consume(rabbitConsumer)
		assert.Equal(t, nil, err, "failed to consume 2", err)
	}()

	time.Sleep(2 * time.Second)

	err = msgBroker.Publish([]byte("hello world 1"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	time.Sleep(2 * time.Second)

	err = msgBroker2.Publish([]byte("hello world 3"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	time.Sleep(2 * time.Second)

	err = msgBroker.Publish([]byte("hello world 2"))
	assert.Equal(t, nil, err, "failed to Publish message", err)

	time.Sleep(2 * time.Second)

}
