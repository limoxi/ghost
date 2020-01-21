package ghost

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var fsWatcher *fsnotify.Watcher

func watchDir(watcher *fsnotify.Watcher, dirPath string){
	err := watcher.Add(dirPath)
	if err != nil{
		log.Fatal(err)
	}
	log.Println("[fsnotify] watch dir: ", dirPath)
	fis, _ := ioutil.ReadDir(dirPath)
	for _, fi := range fis{
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") && !strings.HasPrefix(fi.Name(), "__"){
			pp, _ := filepath.Abs(path.Join(dirPath, fi.Name()))
			watchDir(watcher, pp)
		}
	}
}

func watchProject(){
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case _, ok := <-fsWatcher.Events:
				if !ok {
					return
				}
				log.Println("[fsnotify] project file changed, shutting down server in 5s...")
				time.AfterFunc(time.Second * 10, func() {
					os.Exit(0)
				})
			case err, ok := <-fsWatcher.Errors:
				if !ok {
					return
				}
				log.Println("project watcher error:", err)
			}
		}
	}()
	curDir, _ := os.Getwd()
	watchDir(fsWatcher, curDir)
}