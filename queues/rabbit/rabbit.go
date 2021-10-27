package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
)

var connection *amqp.Connection
var lazyInit func() error

func Setup(username, password, port, address string) {
	if connection != nil {
		_ = connection.Close()
		connection = nil
	}
	lazyInit = func() error {
		err := connect(
			username,
			password,
			port,
			address,
		)
		return err
	}
}

func connect(username string, password string, address, port string) error {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s",
			username,
			password,
			address,
			port,
		),
	)
	connection = conn
	return err
}

//
//import (
//	"fmt"
//	"git.aasal.co/octal/octalbit/backend/utils/queues"
//	"github.com/streadway/amqp"
//	"reflect"
//	"sync"
//)

//type Configs struct {
//	Address  string
//	Port     string
//	Username string
//	Password string
//}
//
//func NewQueue(address, port, username, password string) (queues.Queue, error) {
//	return dialRabbitQueue(address, port, username, password)
//}
//
//func NewQueueWithConfigStruct(conf Configs) (queues.Queue, error) {
//	return dialRabbitQueue(conf.Address, conf.Port, conf.Username, conf.Password)
//}
//
//func (r *rabbit) Create(exChangeName, exchangeKind, queueName, bindingName string) error {
//	if r.Queues.isExist(exChangeName, queueName, bindingName) {
//		return fmt.Errorf("this queue is already exist")
//	}
//	return r.createQueue(
//		queueInfo{
//			exChangeName,
//			exchangeKind,
//			queueName,
//			bindingName,
//		},
//	)
//}
//
//func (r *rabbit) Publish(exChangeName, routingName string, body []byte) error {
//	if r.Chan == nil {
//		return fmt.Errorf("create must be called first")
//	}
//	//TODO: check exChangeName and routingName are exist or not
//	return r.Chan.Publish(
//		exChangeName,
//		routingName,
//		false,
//		false,
//		amqp.Publishing{
//			ContentType:  "application/json",
//			Body:         body,
//			DeliveryMode: amqp.Persistent,
//		})
//}
//
//func (r *rabbit) Consume(queueName string, autoAck bool, itemHook func(queues.Item)) error {
//	if r.Chan == nil {
//		return fmt.Errorf("create must be called first")
//	}
//	//TODO: check queueName is exist or not
//	consumer, err := r.Chan.Consume(
//		queueName,
//		"",
//		autoAck,
//		false,
//		false,
//		false,
//		nil,
//	)
//	if err != nil {
//		return err
//	}
//	go func() {
//		for consumeItem := range consumer {
//			if autoAck == false {
//				consumeItem.Ack(false)
//			}
//			item := queues.Item{
//				Body:      consumeItem.Body,
//				Timestamp: consumeItem.Timestamp,
//			}
//			itemHook(item)
//		}
//	}()
//	return nil
//}
//
//type rabbit struct {
//	Conn   *amqp.Connection
//	Chan   *amqp.Channel
//	Queues queuesInfo
//}
//
//func dialRabbitQueue(address, port, username, password string) (queues.Queue, error) {
//	conn, err := amqp.Dial(
//		fmt.Sprintf("amqp://%s:%s@%s:%s",
//			username,
//			password,
//			address,
//			port,
//		),
//	)
//	if err != nil {
//		return nil, err
//	}
//	return &rabbit{Conn: conn, Queues: []queueInfo{}}, nil
//}
//
//type queueInfo struct {
//	ExChangeName string
//	ExChangeKind string
//	QueueName    string
//	BindingName  string
//}
//
//type queuesInfo []queueInfo
//
//func (qs *queuesInfo) isExist(exChangeName, queueName, bindingName string) bool {
//	flag := false
//	switch reflect.TypeOf(*qs).Kind() {
//	case reflect.Slice:
//		{
//			s := reflect.ValueOf(*qs)
//			fmt.Println(s.Len())
//			for i := 0; i < s.Len(); i++ {
//				itemI := s.Index(i).Interface().(queueInfo)
//				if itemI.ExChangeName == exChangeName &&
//					itemI.QueueName == queueName &&
//					itemI.BindingName == bindingName {
//					flag = true
//				}
//			}
//		}
//	}
//	return flag
//}
//
//var onlyOnce sync.Once
//
//func (r *rabbit) createQueue(info queueInfo) error {
//	onlyOnce.Do(func() {
//		//initialize the channel
//		r.Chan, _ = r.Conn.Channel()
//	})
//	err := r.Chan.ExchangeDeclare(
//		info.ExChangeName,
//		info.ExChangeKind,
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//	if err != nil {
//		return err
//	}
//	_, err = r.Chan.QueueDeclare(
//		info.QueueName,
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//	if err != nil {
//		return err
//	}
//	err = r.Chan.QueueBind(
//		info.QueueName,
//		info.BindingName,
//		info.ExChangeName,
//		false,
//		nil,
//	)
//	if err != nil {
//		return err
//	}
//	r.Queues = append(r.Queues, info)
//	return nil
//}
