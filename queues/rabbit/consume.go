package rabbit

import (
	"encoding/json"
	"git.aasal.co/octal/octalbit/backend/utils/queues"
	"github.com/streadway/amqp"
)

var consumer = make(map[string]queues.Consumer)

func GetConsumer(config map[string]string) (queues.Consumer, error) {
	if connection == nil {
		err := lazyInit()
		if err != nil {
			return nil, err
		}
	}
	consumerConfig := rabbitConsumeConfigFromMap(config)
	jsonKey, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	key := string(jsonKey)
	if consumer[key] == nil {
		err := setConsumerValue(key, consumerConfig)
		if err != nil {
			return nil, err
		}
	}
	return consumer[key], nil

}

func setConsumerValue(key string, config ConsumeConfig) error {
	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	consumer[key] = &consume{channel: channel, config: config}
	return nil
}

type consume struct {
	channel *amqp.Channel
	config  ConsumeConfig
}

func (c *consume) Consume(itemHook func(queues.Item) error) error {
	consumer, err := c.channel.Consume(
		c.config.QueueName,
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
	go func() {
		for consumeItem := range consumer {
			err := itemHook(
				queues.Item{
					Body: consumeItem.Body,
				},
			)
			if err != nil {
				_ = c.Dispose()
				return
			}
			_ = consumeItem.Ack(false)
		}
	}()
	return nil
}

func (c *consume) Dispose() error {
	return c.channel.Close()
}
