package xkafka

import (
	"encoding/json"
	"fmt"
	"github.com/amiranmanesh/octalbit-utils/socket"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var topics = make(map[string]socket.MessageBroker)

func GetMessageBroker(config interface{}) (socket.MessageBroker, error) {
	kafkaConfig, ok := config.(Config)
	if !ok {
		return nil, fmt.Errorf("failed to parse kafka config")
	}
	jsonKey, err := json.Marshal(kafkaConfig)
	if err != nil {
		return nil, err
	}
	key := string(jsonKey)
	if topics[key] == nil {
		setTopicsMapValue(key, kafkaConfig)
	}
	return topics[key], nil
}

func setTopicsMapValue(key string, config Config) {
	newProducerLazyFunc := func(config Config) (*kafka.Producer, error) {
		return kafka.NewProducer(
			&kafka.ConfigMap{
				"bootstrap.servers": config.Servers,
				"group.id":          config.GroupID,
			},
		)
	}
	newConsumerLazyFunc := func(config Config) (*kafka.Consumer, error) {
		return kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": config.Servers,
			"group.id":          config.GroupID,
			"auto.offset.reset": "earliest",
		})
	}

	topics[key] = &broker{
		Config:              config,
		lazyMakeNewProducer: newProducerLazyFunc,
		lazyMakeNewConsumer: newConsumerLazyFunc,
	}
}

type broker struct {
	Config              Config
	Producer            *kafka.Producer
	Consumer            *kafka.Consumer
	lazyMakeNewProducer func(config Config) (*kafka.Producer, error)
	lazyMakeNewConsumer func(config Config) (*kafka.Consumer, error)
}

func (b *broker) Publish(body []byte) error {
	if b.Producer == nil {
		var err error
		b.Producer, err = b.lazyMakeNewProducer(b.Config)
		if err != nil {
			return err
		}
	}
	return b.Producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &b.Config.ProduceTopic, Partition: kafka.PartitionAny},
			Value:          body,
		}, nil,
	)
}

func (b *broker) Consume(itemHook func(socket.MessageBrokerItem) error) error {
	if b.Consumer == nil {
		var err error
		b.Consumer, err = b.lazyMakeNewConsumer(b.Config)
		if err != nil {
			return err
		}
	}
	err := b.Consumer.SubscribeTopics(b.Config.ConsumeTopics, nil)
	if err != nil {
		return err
	}
	for {
		msg, err := b.Consumer.ReadMessage(-1)
		if err != nil {
			return err
		}
		_ = itemHook(socket.MessageBrokerItem{Body: msg.Value})
	}
}

func (b *broker) Clear() error {
	//TODO implement
	return nil
}

func (b *broker) Dispose() error {
	if b.Producer != nil {
		b.Producer.Close()
	}
	if b.Consumer != nil {
		err := b.Consumer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
