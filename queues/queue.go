package queues

type Item struct {
	Body []byte
}

//type Queue interface {
//	Create(exChangeName, exchangeKind, queueName, bindingName string) error
//	Publish(exChangeName, routingName string, body []byte) error
//	Consume(queueName string, autoAck bool, itemHook func(Item)) error
//}

type Creator interface {
	Create(queueName string) Publisher
}

type Publisher interface {
	Publish(queueName string, body []byte) error
}

type Consumer interface {
	Consume(queueName string, itemHook func(Item) error) error
	Dispose(queueName string) error
}
