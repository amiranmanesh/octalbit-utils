package shutdown

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//RealGraceFul can garcefully shutdown a http server
//and lunch a function for additional jobs such as
//closing connections, job queues and ...,
//Remember to ignore the ErrServerClosed error from http listeners
func RealGraceFul(server *http.Server, f func()) {

	//init OS signal chan
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigChan

	if server != nil {
		//init context with timeout to make sure of shutdown return
		//in case of zombie connections
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		//release resources if shutdown returned early
		defer cancel()
		server.Shutdown(ctx)
	}

	//lunch additional func by priority order
	if f != nil {
		f()
	}
}
