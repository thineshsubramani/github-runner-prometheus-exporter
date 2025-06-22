package watcher

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchLogDir(dir string, onChange func(path string)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("📡 File updated:", event.Name)
				if filepath.Ext(event.Name) == ".log" {
					onChange(event.Name)
				}
			}
		case err := <-watcher.Errors:
			log.Println("❌ Watcher error:", err)
		}
	}
}
