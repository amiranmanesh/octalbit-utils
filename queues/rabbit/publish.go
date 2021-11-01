package rabbit

import (
	"encoding/json"

	"github.com/streadway/amqp"

	"github.com/amiranmanesh/octalbit-utils/queues"
)

var publisher = make(map[string]queues.Publisher)

func GetPublisher(config map[string]string) (queues.Publisher, error) {
	if connection == nil {
		err := lazyInit()
		if err != nil {
			return nil, err
		}
	}
	publishConfig := rabbitQueueConfigFromMap(config)
	jsonKey, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	key := string(jsonKey)
	if publisher[key] == nil {
		err := setPublisherValue(key, publishConfig)
		if err != nil {
			return nil, err
		}
	}
	return publisher[key], nil
}

func setPublisherValue(key string, config QueueConfig) error {
	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	err = createQueue(channel, config)
	if err != nil {
		return err
	}
	publisher[key] = &publish{channel: channel, config: config}
	return nil
}

type publish struct {
	channel *amqp.Channel
	config  QueueConfig
}

func (p *publish) Publish(body []byte) error {
	return p.channel.Publish(
		p.config.ExChangeName,
		p.config.RouteName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
}

func (p *publish) Dispose() error {
	return p.channel.Close()
}
