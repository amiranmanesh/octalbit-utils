package queues

type Item struct {
	Body []byte
}

type Publisher interface {
	Publish(body []byte) error
	Dispose() error
}

type Consumer interface {
	Consume(itemHook func(Item) error) error
	Dispose() error
}
