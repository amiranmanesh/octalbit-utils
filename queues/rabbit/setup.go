package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
)

var connection *amqp.Connection
var lazyInit func() error

func Setup(username, password, address string, port int) {
	if connection != nil {
		_ = connection.Close()
		connection = nil
	}
	lazyInit = func() error {
		err := connect(
			username,
			password,
			address,
			port,
		)
		return err
	}
}

func connect(username, password, address string, port int) error {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%d",
			username,
			password,
			address,
			port,
		),
	)
	connection = conn
	return err
}
