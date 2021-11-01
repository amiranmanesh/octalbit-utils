package rabbit

import "github.com/streadway/amqp"

func createQueue(channel *amqp.Channel, config QueueConfig) error {

	queueObject := &queue{channel: channel, config: config}
	err := queueObject.exchangeDeclare()
	if err != nil {
		return err
	}
	_, err = queueObject.queueDeclare()
	if err != nil {
		return err
	}
	err = queueObject.queueBind()
	if err != nil {
		return err
	}
	return nil
}

type queue struct {
	channel *amqp.Channel
	config  QueueConfig
}

func (q *queue) exchangeDeclare() error {
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

func (q *queue) queueDeclare() (amqp.Queue, error) {
	return q.channel.QueueDeclare(
		q.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (q *queue) queueBind() error {
	return q.channel.QueueBind(
		q.config.QueueName,
		q.config.BindingName,
		q.config.ExChangeName,
		false,
		nil,
	)
}
