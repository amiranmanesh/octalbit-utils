package rabbit

import (
	"encoding/json"

	"github.com/streadway/amqp"

	"github.com/amiranmanesh/octalbit-utils/queues"
)

var publisher = make(map[string]queues.Publisher)

func GetPublisher(configPointer *map[string]string) (queues.Publisher, error) {
	if connection == nil {
		err := lazyInit()
		if err != nil {
			return nil, err
		}
	}
	config := rabbitPublishConfigFromMap(configPointer)
	jsonKey, err := json.Marshal(*configPointer)
	if err != nil {
		return nil, err
	}
	key := string(jsonKey)
	if publisher[key] == nil {
		err := setPublisherValue(key, config)
		if err != nil {
			return nil, err
		}
	}
	return publisher[key], nil
}

func setPublisherValue(key string, config PublishConfig) error {
	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	publishObject := &publish{channel: channel, config: config}
	err = publishObject.exchangeDeclare()
	if err != nil {
		return err
	}
	_, err = publishObject.queueDeclare()
	if err != nil {
		return err
	}
	err = publishObject.queueBind()
	if err != nil {
		return err
	}
	publisher[key] = publishObject
	return nil
}

type publish struct {
	channel *amqp.Channel
	config  PublishConfig
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

func (p *publish) exchangeDeclare() error {
	return p.channel.ExchangeDeclare(
		p.config.ExChangeName,
		p.config.ExChangeKind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (p *publish) queueDeclare() (amqp.Queue, error) {
	return p.channel.QueueDeclare(
		p.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (p *publish) queueBind() error {
	return p.channel.QueueBind(
		p.config.QueueName,
		p.config.BindingName,
		p.config.ExChangeName,
		false,
		nil,
	)
}
