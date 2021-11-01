package rabbit

const (
	TopicExchange  = "topic"
	DirectExchange = "direct"
	FanoutExchange = "fanout"
	HeaderExchange = "header"
)

type PublishConfig struct {
	QueueName    string
	ExChangeName string
	BindingName  string
	RouteName    string
	ExChangeKind string
}

type ConsumeConfig struct {
	QueueName string
}

func (config PublishConfig) ToMap() map[string]string {
	return map[string]string{
		"queueName":    config.QueueName,
		"exChangeName": config.ExChangeName,
		"bindingName":  config.BindingName,
		"routeName":    config.RouteName,
		"exchangeKind": config.ExChangeKind,
	}
}

func rabbitPublishConfigFromMap(rabbitMap map[string]string) PublishConfig {
	return PublishConfig{
		QueueName:    rabbitMap["queueName"],
		ExChangeName: rabbitMap["exChangeName"],
		BindingName:  rabbitMap["bindingName"],
		RouteName:    rabbitMap["routeName"],
		ExChangeKind: rabbitMap["exchangeKind"],
	}
}

func (config ConsumeConfig) ToMap() map[string]string {
	return map[string]string{
		"queueName": config.QueueName,
	}
}

func rabbitConsumeConfigFromMap(rabbitMap map[string]string) ConsumeConfig {
	return ConsumeConfig{
		QueueName: rabbitMap["queueName"],
	}
}
