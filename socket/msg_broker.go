package socket

type MessageBrokerItem struct {
	Body []byte
}

type MessageBroker interface {
	Publish(body []byte) error
	Consume(itemHook func(MessageBrokerItem) error) error
	Clear() error
	Dispose() error
}
