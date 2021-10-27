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

func (configPointer *PublishConfig) ToMap() map[string]string {
	return map[string]string{
		"queueName":    configPointer.QueueName,
		"exChangeName": configPointer.ExChangeName,
		"bindingName":  configPointer.BindingName,
		"routeName":    configPointer.RouteName,
		"exchangeKind": configPointer.ExChangeKind,
	}
}

func rabbitPublishConfigFromMap(rabbitMapPointer *map[string]string) PublishConfig {
	var rabbitMap = *rabbitMapPointer
	return PublishConfig{
		QueueName:    rabbitMap["queueName"],
		ExChangeName: rabbitMap["exChangeName"],
		BindingName:  rabbitMap["bindingName"],
		RouteName:    rabbitMap["routeName"],
		ExChangeKind: rabbitMap["exchangeKind"],
	}
}

func (configPointer *ConsumeConfig) ToMap() map[string]string {
	return map[string]string{
		"queueName": configPointer.QueueName,
	}
}

func rabbitConsumeConfigFromMap(rabbitMapPointer *map[string]string) ConsumeConfig {
	var rabbitMap = *rabbitMapPointer
	return ConsumeConfig{
		QueueName: rabbitMap["queueName"],
	}
}
