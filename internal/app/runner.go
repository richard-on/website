package app

import (
	"os"
	"os/signal"
	"syscall"
)

func (a *App) Run() {
	idleConn := make(chan struct{})

	go func() {
		// Waiting for quit signal on exit
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

		err := a.app.Shutdown()
		if err != nil {
			a.log.Fatalf(err, "could not shutdown server")
		}

		close(idleConn)
		<-quit
	}()

	if err := a.app.Listen(":80"); err != nil {
		a.log.Fatalf(err, "error while listening at port 80")
	}

	<-idleConn
}
