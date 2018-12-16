package config

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"runtime"
	"time"
)

const (
	LAYER_START_MODE_NEW = iota
	LAYER_START_MODE_RECOVER
	LAYER_START_MODE_REBOOT
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

type Server struct {
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
		StartMode  int         `json:"StartMode"`
		Class      string      `json:"Class"`
		SuperLayer int32       `json:"SuperLayer"`
		Attributes interface{} `json:"Attributes"`
		//LocalHostClass  string `json:"LocalHost"`
		//RemoteHostClass string `json:"LocalHost"`
		//TotalHostNum    int64         `json:"TotalHostNum"`
		//LocalHostMod int32  `json:"LocalHostMod`
		//RemoteTable []RemoteEntry `json:"RemoteTable>Entry"`
	} `json:"Layer"`
	Services []struct {
		ServerSchema int32 `json:"ServerSchema"`
		Port         int   `json:"Port"`
	} `json:"Services"`
	IO struct {
		Class      string `json:"Class"`
		Attributes interface{}
	}
	Actors     []*Actor `json:"Actors"`
	ActorMap   map[int32]*Actor
	Routers    []*Router `json:"Routers"`
	RouterMap  map[int32]*Router
	Messages   []*Message `json:"Messages"`
	MessageMap map[int32]*Message
	Hosts      []*Host `json:"Hosts"`
	HostMap    map[int32]*Host
	Links      []*Link `json:"Links"`
	LinkMap    map[int32]*Link
	Topos      []*Topo `json:"Topos"`
	TopoMap    map[int32]*Topo
	Servers    []*Server `json:"Servers"`
	ServerMap  map[int32]*Server
	Clients    []*Client `json:"Clients"`
	ClientMap  map[int32]*Client
}

func (this *Config) init() (err maybe.MaybeError) {
	if this.Layer.Id < 0 {
		err.Error(fmt.Errorf("illegal layer layer: %d", this.Layer.Id))
	}

	this.ActorMap = make(map[int32]*Actor)
	for _, a := range this.Actors {
		if a.Schema < 0 {
			err.Error(fmt.Errorf("illegal actor schema: %d", a.Schema))
			return
		}
		if _, ok := this.ActorMap[a.Schema]; ok {
			err.Error(fmt.Errorf("actor already exists: %d", a.Schema))
			return
		}
		this.ActorMap[a.Schema] = a
	}

	this.RouterMap = make(map[int32]*Router)
	for _, r := range this.Routers {
		if r.Id < 0 {
			err.Error(fmt.Errorf("illegal router id: %d", r.Id))
			return
		}
		if _, ok := this.RouterMap[r.Id]; ok {
			err.Error(fmt.Errorf("router already exists: %d", r.Id))
			return
		}
		this.RouterMap[r.Id] = r
	}

	this.MessageMap = make(map[int32]*Message)
	for _, m := range this.Messages {
		if m.Type < 0 {
			err.Error(fmt.Errorf("illegal message type: %d", m.Type))
			return
		}
		if _, ok := this.MessageMap[m.Type]; ok {
			err.Error(fmt.Errorf("message already exists: %d", m.Type))
			return
		}
		this.MessageMap[m.Type] = m
	}

	this.HostMap = make(map[int32]*Host)
	for _, h := range this.Hosts {
		if h.Schema < 0 {
			err.Error(fmt.Errorf("illegal host schema: %d", h.Schema))
			return
		}
		if _, ok := this.HostMap[h.Schema]; ok {
			err.Error(fmt.Errorf("host already exists: %d", h.Schema))
			return
		}
		this.HostMap[h.Schema] = h
	}

	this.LinkMap = make(map[int32]*Link)
	for _, l := range this.Links {
		if l.Schema < 0 {
			err.Error(fmt.Errorf("illegal link schema: %d", l.Schema))
			return
		}
		if _, ok := this.LinkMap[l.Schema]; ok {
			err.Error(fmt.Errorf("link already exists: %d", l.Schema))
			return
		}
		this.LinkMap[l.Schema] = l
	}

	this.TopoMap = make(map[int32]*Topo)
	for _, t := range this.Topos {
		if t.Schema < 0 {
			err.Error(fmt.Errorf("illegal topo schema: %d", t.Schema))
			return
		}
		if _, ok := this.TopoMap[t.Schema]; ok {
			err.Error(fmt.Errorf("topo already exists: %d", t.Schema))
			return
		}
		this.TopoMap[t.Schema] = t
	}

	this.ServerMap = make(map[int32]*Server)
	for _, s := range this.Servers {
		if s.Schema < 0 {
			err.Error(fmt.Errorf("illegal server schema: %d", s.Schema))
			return
		}
		if _, ok := this.ServerMap[s.Schema]; ok {
			err.Error(fmt.Errorf("server already exists: %d", s.Schema))
			return
		}
		this.ServerMap[s.Schema] = s
	}

	this.ClientMap = make(map[int32]*Client)
	for _, c := range this.Clients {
		if c.Schema < 0 {
			err.Error(fmt.Errorf("illegal client schema: %d", c.Schema))
			return
		}
		if _, ok := this.ClientMap[c.Schema]; ok {
			err.Error(fmt.Errorf("server already exists: %d", c.Schema))
			return
		}
		this.ClientMap[c.Schema] = c
	}

	err.Error(nil)
	return
}

func (this *Config) Process() (err maybe.MaybeError) {
	this.init().Test()

	if this.Layer.StartMode != LAYER_START_MODE_NEW {
		interfaces.DeleteLayer(this.Layer.Id).Test()
		runtime.GC()
		time.Sleep(10 * time.Second)
	}

	interfaces.AddLayer(this).Test()

	err.Error(nil)
	return
}

func (this Config) GetLayerId() int32 {
	return this.Layer.Id
}

func (this Config) GetLayerClass() string {
	return this.Layer.Class
}

func (this Config) GetLayerAttr() interface{} {
	return this.Layer.Attributes
}
