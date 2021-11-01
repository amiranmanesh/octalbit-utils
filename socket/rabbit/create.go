package rabbit

import "github.com/streadway/amqp"

func createQueue(channel *amqp.Channel, config QueueConfig) error {
	newQueueObject := &newQueue{channel: channel, config: config}
	err := newQueueObject.exchangeDeclare()
	if err != nil {
		return err
	}
	_, err = newQueueObject.queueDeclare()
	if err != nil {
		return err
	}
	err = newQueueObject.queueBind()
	if err != nil {
		return err
	}
	return nil
}

type newQueue struct {
	channel *amqp.Channel
	config  QueueConfig
}

func (q *newQueue) exchangeDeclare() error {
	return q.channel.ExchangeDeclare(
		q.config.ExChangeName,
		q.config.ExChangeKind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (q *newQueue) queueDeclare() (amqp.Queue, error) {
	return q.channel.QueueDeclare(
		q.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (q *newQueue) queueBind() error {
	return q.channel.QueueBind(
		q.config.QueueName,
		q.config.BindingName,
		q.config.ExChangeName,
		false,
		nil,
	)
}
