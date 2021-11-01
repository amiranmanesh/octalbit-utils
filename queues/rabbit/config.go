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

func (config QueueConfig) ToMap() map[string]string {
	return map[string]string{
		"queueName":    config.QueueName,
		"exChangeName": config.ExChangeName,
		"bindingName":  config.BindingName,
		"routeName":    config.RouteName,
		"exchangeKind": config.ExChangeKind,
	}
}

func rabbitQueueConfigFromMap(rabbitMap map[string]string) QueueConfig {
	return QueueConfig{
		QueueName:    rabbitMap["queueName"],
		ExChangeName: rabbitMap["exChangeName"],
		BindingName:  rabbitMap["bindingName"],
		RouteName:    rabbitMap["routeName"],
		ExChangeKind: rabbitMap["exchangeKind"],
	}
}
