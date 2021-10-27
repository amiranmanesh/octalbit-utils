package rabbit

import (
	"encoding/json"
	"fmt"
	"git.aasal.co/octal/octalbit/backend/utils/queues"
	"github.com/streadway/amqp"
)

var connection *amqp.Connection
var publisher = make(map[string]queues.Publisher)

func GetPublisher(configPointer *map[string]string) (queues.Publisher, error) {
	if connection == nil {
		err := connect(
			"guest",
			"guest",
			"localhost",
			"5672",
		)
		if err != nil {
			return nil, err
		}
	}
	config := rabbitConfigFromMap(configPointer)
	key := mapToJsonToString(*configPointer)
	if publisher[key] == nil {
		channel, err := connection.Channel()
		if err != nil {
			return nil, err
		}
		publishObject := &publish{channel: channel, config: config}
		err = publishObject.exchangeDeclare()
		if err != nil {
			return nil, err
		}
		_, err = publishObject.queueDeclare()
		if err != nil {
			return nil, err
		}
		err = publishObject.queueBind()
		if err != nil {
			return nil, err
		}
		publisher[key] = publishObject
	}
	return publisher[key], nil
}

type publish struct {
	channel *amqp.Channel
	config  Config
}

func (p *publish) Publish(body []byte) error {
	fmt.Println(p.config.ExChangeName, p.config.RouteName)
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

func connect(username string, password string, address, port string) error {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s",
			username,
			password,
			address,
			port,
		),
	)
	connection = conn
	return err
}

//

func mapToJsonToString(data map[string]string) string {
	key, _ := json.Marshal(data)
	return string(key)
}
