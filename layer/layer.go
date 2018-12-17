package layer

import (
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"github.com/incubator/router"
	"github.com/incubator/serialization"
	"math/rand"
)

type CommonLayer struct {
	space                    string
	id                       int32
	messageClassToType       map[int]int32
	messageCanonicalFromType map[int32]interfaces.RemoteMessage
	routers                  map[int32]interfaces.Router
	messageRouters           map[int32]interfaces.Router
	services                 []interfaces.Server
	cfg                      *config.Config
	version                  int64
	superLayer               int32
	io                       interfaces.IO
}

func (this *CommonLayer) Init(attrs interface{}, cfg *config.Config) (err maybe.MaybeError) {
	if cfg.Layer.Space == "" {
		err.Error(errors.New("layer space not set"))
		return
	}
	if cfg.Layer.Id < 0 {
		err.Error(fmt.Errorf("illegal layer id: %d", cfg.Layer.Id))
		return
	}
	if cfg.Layer.SuperLayer < 0 {
		err.Error(fmt.Errorf("illegal supervisor layer id: %d", cfg.Layer.SuperLayer))
		return
	}

	this.id = cfg.Layer.Id
	this.space = cfg.Layer.Space

	this.messageRouters = make(map[int32]interfaces.Router)
	this.routers = make(map[int32]interfaces.Router)

	this.messageClassToType = make(map[int]int32)
	this.messageCanonicalFromType = make(map[int32]interfaces.RemoteMessage)

	for _, msgCfg := range cfg.MessageMap {
		if _, ok := this.messageCanonicalFromType[msgCfg.Type]; ok {
			err.Error(fmt.Errorf("message canonical type already exists: %d", msgCfg.Type))
			return
		}

		msgCanon := interfaces.GetMessagePrototype(msgCfg.Class).Right()
		msgCanon.SetLayer(int8(this.id))
		msgCanon.SetType(int8(msgCfg.Type))

		this.messageCanonicalFromType[msgCfg.Type] = msgCanon
		_type := serialization.Eface2TypeInt(msgCanon)
		this.messageClassToType[_type] = msgCfg.Type

		_, ok := this.routers[msgCfg.RouterId]
		if !ok {
			routerCfg, ok := cfg.RouterMap[msgCfg.RouterId]
			if !ok {
				err.Error(fmt.Errorf("router cfg for %d not found when register message type %d",
					msgCfg.RouterId, msgCfg.Type))
				return
			}
			this.routers[msgCfg.RouterId] = router.GetRouterPrototype(routerCfg.Class).
				Right().New(routerCfg.Attributes, cfg).(interfaces.MaybeRouter).Right()
		}
		this.messageRouters[msgCfg.Type] = this.routers[msgCfg.RouterId]
	}

	if cfg.IO.Class != "" {
		this.io = interfaces.GetIOPrototype(cfg.IO.Class).Right().New(cfg.IO.Attributes, cfg).(interfaces.MaybeIO).Right()
	}

	this.superLayer = cfg.Layer.SuperLayer
	this.cfg = cfg
	this.version = rand.Int63()

	err.Error(nil)
	return
}

func (this CommonLayer) GetRouter(id int32) (ret interfaces.MaybeRouter) {
	if val, ok := this.messageRouters[id]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("router not found: %d", id))
	return
}

////go:noescape
func (this CommonLayer) GetMessageType(msg interface{}) (ret maybe.MaybeInt32) {
	_type := serialization.Eface2TypeInt(msg)
	if typ, ok := this.messageClassToType[_type]; ok {
		ret.Value(typ)
		return
	}
	ret.Error(fmt.Errorf("message type from not found: %+v", msg))
	return
}

func (this CommonLayer) GetMessageCanonicalFromType(typ int32) (ret interfaces.MaybeRemoteMessage) {
	if val, ok := this.messageCanonicalFromType[typ]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("message canonical from type not found: %d", typ))
	return
}

func (this CommonLayer) GetConfig() interfaces.Config {
	return this.cfg
}

func (this CommonLayer) GetVersion() int64 {
	return this.version
}

func (this CommonLayer) GetSuperLayer() int32 {
	return this.superLayer
}
