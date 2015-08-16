package watcher

import (
	"os"

	l "sisyphe/log"
)

// Add a file or directory to the watcher
func (w *WatcherType) Add(file string) {
	var err error

	// check if file exists
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		l.Logger.Println(err)
		return
	}

	// add the watcher
	err = w.watcher.Add(file)
	if err != nil {
		l.Logger.Println(err)
		return
	}

	if w.dir {
		l.Logger.Println("added directory:", file)
	} else {

		// get hash of file
		hash, err := getHash(file)
		if err != nil {
			l.Logger.Println(err)
			return
		}

		// lock file list map
		w.mu.Lock()

		List[file] = &FileType{
			hash:  hash,
			perms: FILE_DEFAULT,
		}

		w.mu.Unlock()

		l.Logger.Printf("%s", List[file])
		l.Logger.Println("added file:", file)
	}

	return
}

// Remove a file or directory from the watcher
func (w *WatcherType) Remove(file string) {
	var err error

	err = w.watcher.Remove(file)
	if err != nil {
		l.Logger.Println(err)
		return
	}

	if w.dir {
		l.Logger.Println("removed directory:", file)
	} else {
		l.Logger.Println("removed file:", file)
	}

	return
}

// Close the watcher and channels
func (w *WatcherType) Close() {
	var err error

	// close the channels watcher
	w.exit <- true

	// close the watcher
	err = w.watcher.Close()
	if err != nil {
		l.Logger.Println(err)
		return
	}

	return
}
