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

//type Publisher interface {
//	Publish(body []byte) error
//	Dispose() error
//}
//
//type Consumer interface {
//	Consume(itemHook func(Item) error) error
//	Dispose() error
//}
