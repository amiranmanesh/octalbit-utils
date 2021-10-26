package queues

import "time"

type Item struct {
	Body      []byte
	Timestamp time.Time
}

type Queue interface {
	Create(exChangeName, exchangeKind, queueName, bindingName string) error
	Publish(exChangeName, routingName string, body []byte) error
	Consume(queueName string, autoAck bool, itemHook func(Item)) error
}
