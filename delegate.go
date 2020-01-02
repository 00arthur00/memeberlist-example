package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/memberlist"
)

type delegate struct {
}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(nmsg []byte) {
	fmt.Println("NotifyMsg")
	if nmsg[0] != 'd' {
		return
	}
	var u update
	if err := json.Unmarshal(nmsg[1:], &u); err != nil {
		log.Println(err)
		return
	}
	mtx.Lock()
	defer mtx.Unlock()
	if u.Action == "add" {
		for k, v := range u.Data {
			data[k] = v
		}
		return
	}
	if u.Action == "del" {
		for k, _ := range u.Data {
			delete(data, k)
		}
	}
}
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return limitedQueue.GetBroadcasts(overhead, limit)
}
func (d *delegate) LocalState(join bool) []byte {
	mtx.Lock()
	defer mtx.Unlock()
	stateJSON, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return nil
	}

	return stateJSON
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	log.Println("MergeRemoteState")
	if len(buf) == 0 {
		return
	}
	if false == join {
		return
	}
	var m map[string]string
	if err := json.Unmarshal(buf, &m); err != nil {
		log.Println(err)
		return
	}
	mtx.Lock()
	for k, v := range m {
		data[k] = v
	}
	mtx.Unlock()
}

func (d *delegate) NotifyJoin(n *memberlist.Node) {
	fmt.Println(n.Name, " joined")
}

func (d *delegate) NotifyLeave(n *memberlist.Node) {
	fmt.Println(n.Name, " leave")
}

func (d *delegate) NotifyUpdate(n *memberlist.Node) {
	fmt.Println(n.Name, " update")
}
