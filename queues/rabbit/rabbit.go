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
