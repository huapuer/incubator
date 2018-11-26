package config

import (
	"../common/maybe"
	"../layer"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Actor struct {
	Schema     int32       `json:"Schema"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Router struct {
	Id         int32       `json:"Id"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Message struct {
	Type     int32  `json:"Type"`
	Class    string `json:"Class"`
	RouterId int32  `json:"RouterId"`
}

type Host struct {
	Schema     int32       `json:"Schema"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Link struct {
	Schema     int32       `json:"Schema"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Topo struct {
	Schema     int32       `json:"Schema"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Client struct {
	Schema     int32       `json:"Schema"`
	Class      string      `json:"Class"`
	Attributes interface{} `json:"Attributes"`
}

type Config struct {
	Layer struct {
		Space      string      `json:"Space"`
		Id         int32       `json:"Id"`
		Recover    bool        `json:"Recover"`
		Class      string      `json:"Topology"`
		Attributes interface{} `json:"Attributes"`
		//LocalHostClass  string `json:"LocalHost"`
		//RemoteHostClass string `json:"LocalHost"`
		//TotalHostNum    int64         `json:"TotalHostNum"`
		//LocalHostMod int32  `json:"LocalHostMod`
		//RemoteTable []RemoteEntry `json:"RemoteTable>Entry"`
	} `json:"Id"`
	actors   []*Actor `json:"Actors"`
	Actors   map[int32]*Actor
	routers  []*Router `json:"Routers"`
	Routers  map[int32]*Router
	messages []*Message `json:"Messages"`
	Messages map[int32]*Message
	hosts    []*Host `json:"Hosts"`
	Hosts    map[int32]*Host
	links    []*Link `json:"Links"`
	Links    map[int32]*Link
	topos    []*Topo `json:"Topos"`
	Topos    map[int32]*Topo
	clients  []*Client `json:"Clients"`
	Clients  map[int32]*Client
}

func init() {
	configFile, err := ioutil.ReadFile("conf/main.xml")
	if err != nil {
		panic(err)
	}
	cfg := &Config{}
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		panic(err)
	}

	cfg.Process().Test()
}

func (this *Config) Process() (err maybe.MaybeError) {
	if this.Layer.Id < 0 {
		err.Error(fmt.Errorf("illegal layer layer: %d", this.Layer.Id))
	}

	this.Actors = make(map[int32]*Actor)
	for _, a := range this.actors {
		if a.Schema < 0 {
			err.Error(fmt.Errorf("illegal actor schema: %d", a.Schema))
			return
		}
		this.Actors[a.Schema] = a
	}

	this.Routers = make(map[int32]*Router)
	for _, r := range this.routers {
		if r.Id < 0 {
			err.Error(fmt.Errorf("illegal router id: %d", r.Id))
			return
		}
		if _, ok := this.Routers[r.Id]; ok {
			err.Error(fmt.Errorf("router already exists: %d", r.Id))
			return
		}
		this.Routers[r.Id] = r
	}

	this.Messages = make(map[int32]*Message)
	for _, m := range this.messages {
		if m.Type < 0 {
			err.Error(fmt.Errorf("illegal message type: %d", m.Type))
			return
		}
		if _, ok := this.Messages[m.Type]; ok {
			err.Error(fmt.Errorf("message already exists: %d", m.Type))
			return
		}
		this.Messages[m.Type] = m
	}

	this.Hosts = make(map[int32]*Host)
	for _, h := range this.hosts {
		if h.Schema < 0 {
			err.Error(fmt.Errorf("illegal host schema: %d", h.Schema))
			return
		}
		this.Hosts[h.Schema] = h
	}

	this.Links = make(map[int32]*Link)
	for _, l := range this.links {
		if l.Schema < 0 {
			err.Error(fmt.Errorf("illegal link schema: %d", l.Schema))
			return
		}
		this.Links[l.Schema] = l
	}

	this.Topos = make(map[int32]*Topo)
	for _, t := range this.topos {
		if t.Schema < 0 {
			err.Error(fmt.Errorf("illegal topo schema: %d", t.Schema))
			return
		}
		this.Topos[t.Schema] = t
	}

	this.Clients = make(map[int32]*Client)
	for _, c := range this.clients {
		if c.Schema < 0 {
			err.Error(fmt.Errorf("illegal client schema: %d", c.Schema))
			return
		}
		this.Clients[c.Schema] = c
	}

	layer.SetLayer(*this).Test()

	return
}
