package watcher

import (
	"github.com/go-fsnotify/fsnotify"

	l "sisyphe/log"
)

var Watcher WatcherType

type WatcherType struct {
	watcher *fsnotify.Watcher
	Exit    chan bool
}

func init() {
	var err error

	Watcher = WatcherType{
		Exit: make(chan bool),
	}

	// create our new inotify watcher
	Watcher.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		l.Logger.Fatal(err)
	}

	l.Logger.Println("watcher initialized")

	// start our watcher
	go Watcher.watch()

	l.Logger.Println("watching kernel inotify messages")
}

// Watch events from watcher
func (w *WatcherType) watch() {
WatchLoop:
	for {
		select {
		case <-w.Exit:
			l.Logger.Println("inotify watcher exiting")
			break WatchLoop
		case event := <-w.watcher.Events:
			switch event.Op {
			case fsnotify.Write:
				l.Logger.Println("file modified:", event.Name)
			case fsnotify.Chmod:
				l.Logger.Println("file permissions changed:", event.Name)
			case fsnotify.Remove:
				l.Logger.Println("file removed:", event.Name)
			case fsnotify.Rename:
				l.Logger.Println("file renamed:", event.Name)
			case fsnotify.Create:
				l.Logger.Println("file created:", event.Name)
			}
		case err := <-w.watcher.Errors:
			l.Logger.Println("error:", err)
		}

	}
}

// Close the watcher and channels
func (w *WatcherType) Close() {
	var err error

	err = w.watcher.Close()
	if err != nil {
		l.Logger.Fatal(err)
	}

	l.Logger.Println("closed watcher")

	return
}

// Add a file or directory to the watcher
func (w *WatcherType) Add(file string) {
	var err error

	err = w.watcher.Add(file)
	if err != nil {
		l.Logger.Fatal(err)
	}

	l.Logger.Println("added file:", file)

	return
}

// Remove a file or directory from the watcher
func (w *WatcherType) Remove(file string) {
	var err error

	err = w.watcher.Remove(file)
	if err != nil {
		l.Logger.Fatal(err)
	}

	l.Logger.Println("removed file:", file)

	return
}
