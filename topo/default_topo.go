package topo

import (
	"errors"
	"fmt"
	"incubator/common/maybe"
	"incubator/config"
	"incubator/host"
)

const (
	defaultActorClassName = "actor.defaultActor"
)

func init() {
	RegisterTopoPrototype(defaultActorClassName, &defaultTopo{})
}

type defaultTopo struct {
	localNum        int64
	localOffset     int32
	localHostClass  string
	remoteHostClass string
	localHosts      []host.Host
	remoteHosts     []host.Host
}

func NewDefaultTopo(localNum int64, offset int32, localHostClass string, remoteHostClass string) (ret MaybeTopo) {
	if localNum <= 0 {
		ret.Error(fmt.Errorf("illegal local num: %d", localNum))
		return
	}
	if offset <= 0 {
		ret.Error(fmt.Errorf("illegal offset: %d", localNum))
		return
	}
	if localHostClass == "" {
		ret.Error(errors.New("empty local host class name"))
		return
	}
	if remoteHostClass == "" {
		ret.Error(errors.New("empty local host class name"))
		return
	}
	return
}

func (this *defaultTopo) Scatter(config config.Config) (err maybe.MaybeError) {
	return
}

func (this *defaultTopo) Lookup(id int64) (ret host.MaybeHost) {
	return
}

func (this *defaultTopo) New(cfg config.Config) config.IOC {
	ret := MaybeTopo{}
	topo := &defaultTopo{}
	attrs := cfg.Topo.Attributes.(map[string]interface{})
	if localNum, ok := attrs["LocalNum"]; ok {
		if localNumInt, ok := localNum.(int64); ok {
			topo.localNum = localNumInt
		} else {
			ret.Error(fmt.Errorf("local host num cfg type error(expecting int): %+v", localNum))
			return ret
		}
	} else {
		ret.Error(errors.New("attribute LocalNum not found"))
		return ret
	}

	if localOffset, ok := attrs["LocalOffset"]; ok {
		if localOffsetInt, ok := localOffset.(int32); ok {
			topo.localOffset = int32(localOffsetInt)
		} else {
			ret.Error(fmt.Errorf("local host offset cfg type error(expecting int): %+v", localOffset))
			return ret
		}
	} else {
		ret.Error(errors.New("attribute LocalOffset not found"))
		return ret
	}

	if localHostClass, ok := attrs["LocalHostClass"]; ok {
		if localHostClassStr, ok := localHostClass.(string); ok {
			if localHostClassStr == "" {
				ret.Error(errors.New("empty LocalHostClass"))
				return ret
			}
			topo.localHostClass = localHostClassStr
		} else {
			ret.Error(fmt.Errorf("local host class cfg type error(expecting string): %+v", localHostClass))
		}
	} else {
		ret.Error(errors.New("attribute LocalHostClass not found"))
		return ret
	}

	if remoteHostClass, ok := attrs["RemoteHostClass"]; ok {
		if remoteHostClassStr, ok := remoteHostClass.(string); ok {
			if remoteHostClassStr == "" {
				ret.Error(errors.New("empty RemoteHostClass"))
				return ret
			}
			topo.remoteHostClass = remoteHostClassStr
		} else {
			ret.Error(fmt.Errorf("remote host class cfg type error(expecting string): %+v", remoteHostClass))
		}
	} else {
		ret.Error(errors.New("attribute RemoteHostClass not found"))
		return ret
	}

	if remoteEntries, ok := attrs["RemoteEntries"]; ok {
		if entries, ok:=remoteEntries.([])
	} else {
		ret.Error(errors.New("attribute RemoteEntries not found"))
		return ret
	}

	for i := 0; int64(i) < topo.localNum; i++ {
		localHost := host.GetHostPrototype(topo.localHostClass).Right()
		//TODO: add shuffle logic here
		localHost.SetId(int64(i))
		topo.localHosts = append(topo.localHosts, localHost)
	}

	return ret
}
