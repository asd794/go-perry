package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	w    *fsnotify.Watcher
	send func(event fsnotify.Event)
}

func New(send func(event fsnotify.Event)) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		w:    w,
		send: send,
	}

	go watcher.loop()
	return watcher, nil
}

func (w *Watcher) Watch(path string) error { return w.w.Add(path) }

func (w *Watcher) Unwatch(path string) error { return w.w.Remove(path) }

func (w *Watcher) Close() { w.w.Close() }

func (w *Watcher) loop() {
	for {
		select {
		case event, ok := <-w.w.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				// Handle event
				w.send(event)
			}
		case err, ok := <-w.w.Errors:
			if !ok {
				return
			}
			// Handle error
			log.Default().Println("error:", err)
		}
	}
}

// example
func main() {
	Watcher, err := New(func(event fsnotify.Event) {
		log.Default().Println("event:", event)
	})
	if err != nil {
		log.Fatal(err)
	}
	defer Watcher.Close()

	err = Watcher.Watch(".")
	if err != nil {
		log.Fatal(err)
	}

	// Block forever
	select {}
}
