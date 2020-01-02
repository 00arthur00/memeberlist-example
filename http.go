package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/memberlist"
)

type update struct {
	Action string            `json:"action"`
	Data   map[string]string `json:"data"`
}
type broadcast struct {
	msg    []byte
	notify chan struct{}
}

// Invalidates checks if enqueuing the current broadcast
// invalidates a previous broadcast
func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

// Returns a byte form of the message
func (b *broadcast) Message() []byte {
	return b.msg
}

// Finished is invoked when the message will no longer
// be broadcast, either due to invalidation or to the
// transmit limit being reached
func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := r.Form.Get("key")
	value := r.Form.Get("val")
	fmt.Println("key:", key, "val:", value)
	mtx.Lock()
	data[key] = value
	mtx.Unlock()
	b, err := json.Marshal(&update{

		Action: "add",
		Data: map[string]string{
			key: value,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	limitedQueue.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
}
func del(w http.ResponseWriter, r *http.Request) {
	key := r.Form.Get("key")
	mtx.Lock()
	delete(data, key)
	mtx.Unlock()
	b, err := json.Marshal(&update{

		Action: "del",
		Data: map[string]string{
			key: "",
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	limitedQueue.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
}

func get(w http.ResponseWriter, r *http.Request) {
	mtx.Lock()
	defer mtx.Unlock()
	buf, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(buf)
}
