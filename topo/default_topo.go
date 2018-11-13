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
	commonTopo

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

func (this *defaultTopo) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeTopo{}
	topo := &defaultTopo{
		localHosts: make([]host.Host, 0, 0),
		remoteHosts: make([]host.Host, 0, 0),
	}
	attrsMap, ok := cfg.Topo.Attributes.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal cfg type when new topo %s", defaultActorClassName))
		return ret
	}
	totalHostNum, ok := attrsMap["TotalHostNum"]
	if !ok {
		ret.Error(errors.New("attribute TotalHostNum not found"))
		return ret
	}
	localNumInt, ok := totalHostNum.(int64)
	if !ok {
		ret.Error(fmt.Errorf("total host num cfg type error(expecting int): %+v", totalHostNum))
		return ret
	}
	if localNumInt <= 0 {
		ret.Error(fmt.Errorf("illegal total host num : %d", localNumInt))
		return ret
	}
	topo.totalHostNum = localNumInt

	localHostMod, ok := attrsMap["LocalHostMod"]
	if !ok {
		ret.Error(errors.New("attribute LocalHostMod not found"))
		return ret
	}
	localHostModInt, ok := localHostMod.(int32)
	if !ok {
		ret.Error(fmt.Errorf("local host mod cfg type error(expecting int): %+v", localHostMod))
		return ret
	}
	if localHostModInt <= 0 {
		ret.Error(fmt.Errorf("illegal local host mod : %d", localHostModInt))
		return ret
	}
	topo.localHostMod = localHostModInt

	localHostClass, ok := attrsMap["LocalHostClass"]
	if !ok {
		ret.Error(errors.New("attribute LocalHostClass not found"))
		return ret
	}
	localHostClassStr, ok := localHostClass.(string)
	if !ok {
		ret.Error(fmt.Errorf("local host class cfg type error(expecting string): %+v", localHostClass))
		return ret
	}
	if localHostClassStr == ""{
		ret.Error(errors.New("empty LocalHostClass"))
		return ret
	}
	topo.localHostClass = localHostClassStr

	remoteHostClass, ok := attrsMap["RemoteHostClass"]
	if !ok {
		ret.Error(errors.New("attribute RemoteHostClass not found"))
		return ret
	}
	remoteHostClassStr, ok := remoteHostClass.(string)
	if !ok {
		ret.Error(fmt.Errorf("remote host class cfg type error(expecting string): %+v", remoteHostClass))
		return ret
	}
	if remoteHostClassStr == ""{
		ret.Error(errors.New("empty RemoteHostClass"))
		return ret
	}
	topo.remoteHostClass = remoteHostClassStr

	remoteEntries, ok := attrsMap["RemoteEntries"]
	if !ok {
		ret.Error(errors.New("attribute RemoteEntries not found"))
		return ret
	}
	remoteEntriesMap, ok := remoteEntries.([]map[string]string)
	if !ok {
		ret.Error(fmt.Errorf("attribute RemoteEntries has illegal type(expecting []map[string]string: %+v", remoteEntries))
		return ret
	}
	entryNum := len(remoteEntriesMap)

	if topo.localHostMod != int32(entryNum) {
		ret.Error(fmt.Errorf("local offset(%d) != total entry num - 1(%d)", topo.localHostMod, entryNum))
		return ret
	}

	for i := 0; i < entryNum; i++ {
		remoteHostCfg, ok := cfg.Actors[topo.remoteHostClass]
		if !ok {
			ret.Error(fmt.Errorf("no remote host cfg found: %s", topo.remoteHostClass))
			return ret
		}
		topo.remoteHosts = append(
			topo.remoteHosts, host.GetHostPrototype(topo.remoteHostClass).Right().(config.IOC).New(remoteHostCfg.Attributes, cfg).(host.MaybeHost).Right())
	}

	for i:=0;int64(i)<topo.totalHostNum;i++ {
		if int32(i%entryNum) == topo.localHostMod {
			localHostCfg, ok := cfg.Actors[topo.localHostClass]
			if !ok {
				ret.Error(fmt.Errorf("no local host cfg found: %s", topo.localHostClass))
				return ret
			}
			localHost := host.GetHostPrototype(topo.localHostClass).Right().New(localHostCfg.Attributes, cfg).(host.MaybeHost).Right()
			localHost.SetId(int64(i))
			topo.localHosts = append(topo.localHosts, localHost)
		}
	}

	ret.Value(topo)
	return ret
}
