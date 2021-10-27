package rabbit

const (
	TopicExchange  = "topic"
	DirectExchange = "direct"
	FanoutExchange = "fanout"
	HeaderExchange = "header"
)

type Config struct {
	QueueName    string
	ExChangeName string
	BindingName  string
	RouteName    string
	ExChangeKind string
}

func (configPointer *Config) ConfigToMap() map[string]string {
	return map[string]string{
		"queueName":    configPointer.QueueName,
		"exChangeName": configPointer.ExChangeName,
		"bindingName":  configPointer.BindingName,
		"routeName":    configPointer.RouteName,
		"exchangeKind": configPointer.ExChangeKind,
	}
}

func rabbitConfigFromMap(rabbitMapPointer *map[string]string) Config {
	var rabbitMap = *rabbitMapPointer
	return Config{
		QueueName:    rabbitMap["queueName"],
		ExChangeName: rabbitMap["exChangeName"],
		BindingName:  rabbitMap["bindingName"],
		RouteName:    rabbitMap["routeName"],
		ExChangeKind: rabbitMap["exchangeKind"],
	}
}
