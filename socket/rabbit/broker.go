package rabbit

import (
	"encoding/json"
	"fmt"
	"github.com/amiranmanesh/octalbit-utils/socket"
	"github.com/streadway/amqp"
)

var queues = make(map[string]socket.MessageBroker)

func GetMessageBroker(config interface{}) (socket.MessageBroker, error) {
	queueConfig, ok := config.(QueueConfig)
	if !ok {
		return nil, fmt.Errorf("failed to parse rabbit config")
	}
	if connection == nil {
		err := lazyInit()
		if err != nil {
			return nil, err
		}
	}
	jsonKey, err := json.Marshal(queueConfig)
	if err != nil {
		return nil, err
	}
	key := string(jsonKey)
	if queues[key] == nil {
		err := setQueuesMapValue(key, queueConfig)
		if err != nil {
			return nil, err
		}
	}
	return queues[key], nil
}

func setQueuesMapValue(key string, config QueueConfig) error {
	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	err = createQueue(channel, config)
	if err != nil {
		return err
	}

	queues[key] = &queue{channel: channel, config: config}
	return nil
}

type queue struct {
	channel *amqp.Channel
	config  QueueConfig
}

func (q *queue) Publish(body []byte) error {
	return q.channel.Publish(
		q.config.ExChangeName,
		q.config.RouteName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
}

func (q *queue) Consume(itemHook func(socket.MessageBrokerItem) error) error {
	consumer, err := q.channel.Consume(
		q.config.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for consumeItem := range consumer {
		err := itemHook(
			socket.MessageBrokerItem{
				Body: consumeItem.Body,
			},
		)
		if err != nil {
			return err
		}
		_ = consumeItem.Ack(false)
	}
	return nil
}

func (q *queue) Clear() error {
	_, err := q.channel.QueuePurge(q.config.QueueName, false)
	return err
}

func (q *queue) Dispose() error {
	return q.channel.Close()
}
