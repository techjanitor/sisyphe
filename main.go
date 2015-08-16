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
		w.Files.Close()
		w.Dirs.Close()
		exit <- true
	}()
}

func main() {

	w.Files.Add("/Users/puntme/heh")
	w.Files.Add("/Users/puntme/test.go")
	w.Files.Add("/Users/puntme/fakeballs")

	w.Dirs.Add("/Users/puntme/")
	w.Dirs.Add("/Users/asfawef/")

	<-exit

}
