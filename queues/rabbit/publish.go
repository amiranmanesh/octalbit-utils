package rabbit

//import (
//	"fmt"
//	"git.aasal.co/octal/octalbit/backend/utils/queues"
//	"github.com/streadway/amqp"
//)
//
//func NewQueueCreator(conf Configs) (queues.Creator, error) {
//	return dialRabbitQueue(conf.Address, conf.Port, conf.Username, conf.Password)
//}
//
//
//func (r *rabbit) Create(queueName string) queues.Publisher {
//	return &publish{};
//}
//
//func (p *publish) Publish(queueName string, body []byte) error {
//
//	panic("implement me")
//}
//
/////*
////	This function will publish @body into the Rabbit queue and create the queue if it doesn't exist.'
////	To generate each queue we need 3 keys: ExChangeName , QueueName , BindingName
////	Also we need a ExChangeKind for the queue like `direct` or `topic`.
////	So, The client should pass these keys all in one string as @queueName with this pattern
////	queueName = "{ExChangeKind}={ExChangeName}={BindingKey}={QueueName}={RoutingKey}
////*/
////func (r *rabbit) Publish(queueName string, body []byte) error {
////	panic("implement me")
////}
//
//func (r *rabbit) getOrCreateChannel(queueName string) error {
//	for i, item := range r.Queues {
//
//	}
//}
//
//type rabbit struct {
//	Conn   *amqp.Connection
//	Queues []*queueInfo
//	//Chan   *amqp.Channel
//	//Queues queuesInfo
//}
//
//type publish struct {
//}
//
//type queueInfo struct {
//	QueueName string
//	Channel   *amqp.Channel
//}
//
//func dialRabbitQueue(address, port, username, password string) (queues.Creator, error) {
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
//	return &rabbit{Conn: conn}, nil
//}
