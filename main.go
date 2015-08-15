package main

import (
	"os"
	"os/signal"
	"syscall"

	w "sisyphe/watcher"
)

var (
	exit    chan bool
	signals chan os.Signal
)

func init() {
	exit = make(chan bool)
	signals = make(chan os.Signal, 10)

	// watch for shutdown signals to shutdown cleanly
	signal.Notify(signals, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-signals
		w.Watcher.Exit <- true
		w.Watcher.Close()
		exit <- true
	}()
}

func main() {

	w.Watcher.Add("/Users/puntme/heh")

	<-exit

}
