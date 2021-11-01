package rabbit

const (
	TopicExchange  = "topic"
	DirectExchange = "direct"
	FanoutExchange = "fanout"
	HeaderExchange = "header"
)

type QueueConfig struct {
	QueueName    string
	ExChangeName string
	BindingName  string
	RouteName    string
	ExChangeKind string
}
