package topo

import (
	"errors"
	"fmt"
	"../config"
	"../host"
)

const (
	defaultActorClassName = "actor.defaultActor"
)

func init() {
	RegisterTopoPrototype(defaultActorClassName, &defaultTopo{})
}

type defaultTopo struct {
	totalHostNum    int64
	localHostMod    int32
	localHostClass  string
	remoteHostClass string
	localHosts      []host.Host
	remoteHosts     []host.Host
}

func (this *defaultTopo) Lookup(id int64) (ret host.MaybeHost) {
	return
}

func (this *defaultTopo) New(cfg config.Config) config.IOC {
	ret := MaybeTopo{}
	topo := &defaultTopo{}
	attrs := cfg.Topo.Attributes.(map[string]interface{})
	if totalHostNum, ok := attrs["TotalHostNum"]; ok {
		if localNumInt, ok := totalHostNum.(int64); ok {
			topo.totalHostNum = localNumInt
		} else {
			ret.Error(fmt.Errorf("total host num cfg type error(expecting int): %+v", totalHostNum))
			return ret
		}
	} else {
		ret.Error(errors.New("attribute LocalNum not found"))
		return ret
	}

	if localHostMod, ok := attrs["LocalHostMod"]; ok {
		if localHostMod, ok := localHostMod.(int32); ok {
			topo.localHostMod = int32(localHostMod)
		} else {
			ret.Error(fmt.Errorf("local host mod cfg type error(expecting int): %+v", localHostMod))
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

	entryNum := 0
	if remoteEntries, ok := attrs["RemoteEntries"]; ok {
		if entries, ok := remoteEntries.([]map[string]string); ok {
			entryNum = len(entries)
			for i := 0; i < entryNum; i++ {
				topo.remoteHosts = append(
					topo.remoteHosts, host.GetHostPrototype(topo.remoteHostClass).Right().(config.IOC).New(cfg, int32(i)).(host.MaybeHost).Right())
			}
		} else {
			ret.Error(errors.New("attribute RemoteEntries has illegal type(expecting []map[string]string"))
			return ret
		}
	} else {
		ret.Error(errors.New("attribute RemoteEntries not found"))
		return ret
	}

	// init localHosts
	if topo.localHostMod != int32(entryNum) {
		ret.Error(fmt.Errorf("local offset(%d) != total entry num - 1(%d)", topo.localHostMod, entryNum))
		return ret
	}

	for i:=0;i<topo.totalHostNum;i++ {
		if int32(i%entryNum) == topo.localHostMod {
			localHost := host.GetHostPrototype(topo.localHostClass).Right()
			localHost.SetId(int64(i))
			topo.localHosts = append(topo.localHosts, localHost)
		}
	}

	return ret
}
