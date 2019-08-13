package ghost

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"syscall"
)

func watchProject(){
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("project watcher event:", event)
				syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("project watcher error:", err)
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}
}