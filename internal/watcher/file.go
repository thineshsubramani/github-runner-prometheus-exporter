// Utilizing File notify for track runner state by tracking event.json
package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchLogDir(dir string, onChange func(path string, event string)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	log.Printf("Watcher started on directory: %s", dir)

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	// Initial state check
	eventPath := filepath.Join(dir, "event.json")
	if fileExists(eventPath) {
		log.Println("event.json already exists — assuming runner is busy : ", eventPath)
		onChange(eventPath, "created")
	} else {
		log.Println("No event.json on start — runner idle")
		onChange(eventPath, "deleted")
	}

	for {
		select {
		case event := <-watcher.Events:
			// Just be verbose for debugging
			log.Printf("FS Event: %s on %s", event.Op.String(), event.Name)

			if filepath.Base(event.Name) != "event.json" {
				continue
			}

			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				log.Println("event.json created — runner active")
				onChange(event.Name, "created")

			case event.Op&fsnotify.Remove == fsnotify.Remove,
				event.Op&fsnotify.Rename == fsnotify.Rename:
				log.Println("event.json deleted/renamed — runner idle")
				onChange(event.Name, "deleted")

			case event.Op&fsnotify.Write == fsnotify.Write:
				log.Println("event.json modified")
				onChange(event.Name, "modified")
			}

		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
