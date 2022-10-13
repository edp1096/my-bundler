package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var chFinish chan bool

func globDIR(root string) (dirs []string, err error) {
	err = filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	return
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println(event.Op.String(), event.Name)

				buildJS()
				buildSCSS()
				buildHTML()

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)

			}

			chReload <- true
		}
	}()

	dirs, err := globDIR(sourceRoot)
	if err != nil {
		log.Fatal(err)
	}

	// Add cwd and all sub directories
	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-chFinish
}
