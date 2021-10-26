package queue

import "time"

type Item struct {
	Body      []byte
	Timestamp time.Time
}

type Queue interface {
	Create(exChangeName, queueName, bindingName string) error
	Publish(queueName string, body []byte) error
	Consume(queueName string, itemHook func(Item)) error
}
