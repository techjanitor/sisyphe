package watcher

import (
	"github.com/techjanitor/fsnotify"
	"os"
	"sync"

	l "sisyphe/log"
)

var (
	Files WatcherType
	Dirs  WatcherType
	List  map[string]*FileType
)

type WatcherType struct {
	mu      sync.RWMutex
	watcher *fsnotify.Watcher
	exit    chan bool
	dir     bool
}

type FileType struct {
	hash    string
	perms   os.FileMode
	changed bool
	event   string
}

func init() {
	var err error

	List = make(map[string]*FileType)

	// watcher for files
	Files = WatcherType{
		exit: make(chan bool),
	}

	// watcher for dirs
	Dirs = WatcherType{
		exit: make(chan bool),
		dir:  true,
	}

	// create our new inotify watcher
	Files.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		l.Logger.Fatal(err)
	}

	// create our new inotify watcher
	Dirs.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		l.Logger.Fatal(err)
	}

	// start our watchers
	go Files.watch()
	go Dirs.watch()

}
