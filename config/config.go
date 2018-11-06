package config

import (
	"encoding/json"
	"fmt"
	"../actor"
	"../common/maybe"
	"../message"
	"../router"
	"../topo"
	"io/ioutil"
)

type RemoteEntry struct {
	FromOffset int64  `json:"FromOffset"`
	ToOffset   int64  `json:"ToOffset""`
	Address    string `json:"Address"`
}

type Attribute struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type Actor struct {
	Class      string      `json:"Class"`
	ActorNum   int         `json:ActorNum`
	Attributes interface{} `json:"Attributes"`
}

type Router struct {
	Class      string      `json:"Class"`
	ActorClass string      `json:"ActorClass"`
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
		ClassName  string      `json:"Topology"`
		Attributes interface{} `json:"Attributes"`
		//LocalHostClass  string `json:"LocalHost"`
		//RemoteHostClass string `json:"LocalHost"`
		//LocalNum    int64         `json:"LocalNum"`
		//LocalOffset int32  `json:"LocalOffset`
		//RemoteTable []RemoteEntry `json:"RemoteTable>Entry"`
	} `json:"Topo"`
	actors   []*Actor `json:"Actors"`
	Actors   map[string]*Actor
	routers  []*Router `json:"Routers"`
	Routers  map[string]*Router
	messages []*Message `json:"Messages"`
	Messages map[int]*Message
	hosts    []*Host `json:"Hosts"`
	Hosts    map[string]*Host
}

var (
	GlobalConfig Config
)

func init() {
	configFile, err := ioutil.ReadFile("conf/main.xml")
	if err != nil {
		panic(err)
	}
	GlobalConfig = Config{
		Actors:   make(map[string]*Actor),
		Routers:  make(map[string]*Router),
		Messages: make(map[int]*Message),
	}
	err = json.Unmarshal(configFile, &GlobalConfig)
	if err != nil {
		panic(err)
	}

	GlobalConfig.Process().Test()
}

func (this Config) Process() (err maybe.MaybeError) {
	for i, a := range this.actors {
		if a.Class == "" {
			err.Error(fmt.Errorf("actor class name not set: index[%d]", i))
		}
		this.Actors[a.Class] = a
		actor.AddActors(this, a.Class, a.ActorNum).Test()
	}

	for i, r := range this.routers {
		if r.Class == "" {
			err.Error(fmt.Errorf("router class name not set: index[%d]", i))
		}
		this.Routers[r.Class] = r
		router.AddRouter(this, r.Class).Test()
	}

	for i, m := range this.messages {
		if m.Type <= 0 {
			err.Error(fmt.Errorf("illegal message type: %d, index[%d]", m.Type, i))
		}
		this.Messages[m.Type] = m
		message.RegisterMessageCanonical(this, m.Type)
	}

	for i, h := range this.hosts {
		if h.Class == "" {
			err.Error(fmt.Errorf("host class name not set: index[%d]", i))
		}
		this.Hosts[h.Class] = h
	}

	topo.SetGlobalTopo(this)

	return
}
