package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
)

var addr string
var members string

func init() {
	flag.StringVar(&addr, "port", ":4001", "port for http")
	flag.StringVar(&members, "members", "", "comma seperated list of members")
	flag.Parse()
}

func main() {
	if err := Main(); err != nil {
		panic(err)
	}
	http.HandleFunc("/add", add)
	http.HandleFunc("/del", del)
	http.HandleFunc("/get", get)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func Main() error {

	data = map[string]string{}
	name, _ := os.Hostname()
	c := memberlist.DefaultLANConfig()
	c.Name = name + uuid.NewUUID().String()
	c.BindPort = 0
	c.Delegate = &delegate{}
	c.Events = &delegate{}
	ml, err := memberlist.Create(c)
	if err != nil {
		return err
	}
	if len(members) > 0 {
		ms := strings.Split(members, ",")
		if _, err := ml.Join(ms); err != nil {
			return err
		}
	}
	limitedQueue = memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return ml.NumMembers()
		},
		RetransmitMult: 3,
	}
	return nil
}
