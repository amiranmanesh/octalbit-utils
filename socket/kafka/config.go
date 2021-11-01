package xkafka

type Config struct {
	Servers       string
	GroupID       string
	ProduceTopic  string
	ConsumeTopics []string
}
