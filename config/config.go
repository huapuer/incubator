package config

import (
	"encoding/json"
	"fmt"
	"../common/maybe"
	"../message"
	"../router"
	"../topo"
	"io/ioutil"
)

var(
	topoIndex = 0
)

const(
	spacePerTopo = 100
)

type Actor struct {
	Schema       int32       `json:"Schema"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Router struct {
	Id         int32       `json:"Id"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Message struct {
	Type        int    `json:"Type"`
	Class       string `json:"Class"`
	RouterClass string `json:"Class"`
}

type Host struct {
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Config struct {
	Topo struct {
	     Class      string      `json:"Topology"`
	     Attributes interface{} `json:"Attributes"`
		//LocalHostClass  string `json:"LocalHost"`
		//RemoteHostClass string `json:"LocalHost"`
		//TotalHostNum    int64         `json:"TotalHostNum"`
		//LocalHostMod int32  `json:"LocalHostMod`
		//RemoteTable []RemoteEntry `json:"RemoteTable>Entry"`
	} `json:"Topo"`
	actors   []*Actor `json:"Actors"`
	Actors   map[int32]*Actor
	routers  []*Router `json:"Routers"`
	Routers  map[string]*Router
	messages []*Message `json:"Messages"`
	Messages map[int]*Message
	hosts    []*Host `json:"Hosts"`
	Hosts    map[string]*Host
}

func init() {
	configFile, err := ioutil.ReadFile("conf/main.xml")
	if err != nil {
		panic(err)
	}
	cfg := Config{
		Actors:   make(map[int32]*Actor),
		Routers:  make(map[string]*Router),
		Messages: make(map[int]*Message),
	}
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		panic(err)
	}

	cfg.Process().Test()
}

func (this Config) Process() (err maybe.MaybeError) {
	for i, a := range this.actors {
		if a.Class == "" {
			err.Error(fmt.Errorf("actor class name not set: index[%d]", i))
		}
		this.Actors[a.Schema] = a
	}

	for i, r := range this.routers {
		if r.Class == "" {
			err.Error(fmt.Errorf("router class name not set: index[%d]", i))
		}
		this.Routers[r.Class] = r
		id := int32(topoIndex * spacePerTopo) + r.Id
		router.AddRouter(id, r.Class, this).Test()
	}

	for i, m := range this.messages {
		if m.Type <= 0 {
			err.Error(fmt.Errorf("illegal message type: %d, index[%d]", m.Type, i))
		}
		typ := topoIndex * spacePerTopo + m.Type
		this.Messages[typ] = m
		message.RegisterMessageCanonical(this, typ).Test()
	}

	for i, h := range this.hosts {
		if h.Class == "" {
			err.Error(fmt.Errorf("host class name not set: index[%d]", i))
		}
		this.Hosts[h.Class] = h
	}

	topo.SetTopo(int32(topoIndex), this)
	topoIndex++

	return
}
