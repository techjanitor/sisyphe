package watcher

import (
	"crypto/sha1"
	"fmt"
	"github.com/techjanitor/fsnotify"
	"io"
	"os"

	l "sisyphe/log"
)

const (
	FILE_DEFAULT os.FileMode = 0644
)

// Watch events from watcher
func (w *WatcherType) watch() {

	if w.dir {
		l.Logger.Println("directory watcher initialized")
	} else {
		l.Logger.Println("file watcher initialized")
	}

WatchLoop:
	for {
		select {
		case <-w.exit:
			if w.dir {
				l.Logger.Println("directory watcher closed")
			} else {
				l.Logger.Println("file watcher closed")
			}
			break WatchLoop
		case event := <-w.watcher.Events:
			if w.dir {
				l.Logger.Println("directory modified:", event)
			} else {
				switch event.Op {
				case fsnotify.Write:
					w.checkHash(event)
				case fsnotify.Chmod:
					w.checkPerms(event)
				case fsnotify.Remove:
					l.Logger.Println("file removed:", event.Name, event.Op.String())
				case fsnotify.Rename:
					l.Logger.Println("file renamed:", event.Name, event.Op.String())
				}
			}
		case err := <-w.watcher.Errors:
			l.Logger.Println("error:", err)
		}
	}

	return
}

// Check file hash
func (w *WatcherType) checkHash(event fsnotify.Event) {

	// get hash of file
	hash, err := getHash(event.Name)
	if err != nil {
		l.Logger.Println(err)
		return
	}

	// lock file list map
	w.mu.Lock()

	info := List[event.Name]

	// if the perms are different change state and log
	if hash != info.hash {
		info.changed = true
		info.event = event.Op.String()
		l.Logger.Println("file modified:", event.Name, event.Op.String())
	} else {
		info.changed = false
		info.event = ""
	}

	w.mu.Unlock()

}

// Check file permissions
func (w *WatcherType) checkPerms(event fsnotify.Event) {

	// get file stat
	fi, err := os.Stat(event.Name)
	if err != nil {
		l.Logger.Println(err)
		return
	}

	// lock file list map
	w.mu.Lock()

	info := List[event.Name]

	// if the perms are different change state and log
	if fi.Mode() != info.perms {
		info.changed = true
		info.event = event.Op.String()
		l.Logger.Println("file permissions changed:", event.Name, event.Op.String())
	} else {
		info.changed = false
		info.event = ""
	}

	w.mu.Unlock()

}

// GetHash creates a hash of a file
func getHash(file string) (hash string, err error) {

	// new sha1 hasher
	hasher := sha1.New()

	// open file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	// copy into hasher
	_, err = io.Copy(hasher, f)
	if err != nil {
		return
	}

	// return hash
	hash = fmt.Sprintf("%x", hasher.Sum(nil))

	return
}
