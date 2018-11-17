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
	backupFactor    int32
	localHostSchema  int32
	remoteHostSchema int32
	localHosts      []host.Host
	remoteHosts     []host.Host
	remoteNum       int32
}

func (this *defaultTopo) Lookup(id int64) (ret host.MaybeHost) {
	mod := id%(this.remoteNum + 1)
	idx := id/(this.remoteNum + 1)/(this.backupFactor + 1)+mod
	hosts := make([]host.Host, 0, 0)

	if mod == this.localHostMod {
		if idx > len(this.localHosts) {
			ret.Error(fmt.Errorf("master id exceeds local host range: %d", id))
			return
		}
		hosts = append(hosts, this.localHosts[idx])
	}else{
		ridx := mod
		if mod > this.localHostMod {
			ridx--
		}
		hosts = append(hosts, this.remoteHosts[ridx])
	}
	if mod < this.localHostMod + this.backupFactor {
		if idx > len(this.localHosts) {
			ret.Error(fmt.Errorf("slave id exceeds local host range: %d", id))
			return
		}
		hosts = append(hosts, this.localHosts[idx])
	}
	for offset:=0;offset < this.backupFactor - 1;offset++{
		ridx := (mod + offset)% (this.remoteNum + 1)
		if mod > this.localHostMod {
			ridx--
		}
		hosts = append(hosts, this.remoteHosts[ridx])
	}

	var master host.Host
	slaves := make([]host.Host, 0, 0)
	for _, h := range hosts {
		if h.IsValid(){
			if master == nil {
				master = h
			} else {
				slaves = append(slaves, h)
			}
		}
	}

	if master == nil {
		ret.Error(fmt.Errorf("no available master host found for id: %d", id))
		return
	}

	ret.Value(host.NewDuplicatedHost(master, slaves).Right())
	return
}

func (this *defaultTopo) New(attrs interface{}, cfg config.Config) config.IOC {
	this.Init(cfg)

	ret := MaybeTopo{}
	topo := &defaultTopo{
		commonTopo{
			layer: cfg.Topo.Layer,
		},
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

	backupFactor, ok := attrsMap["BackupFactor"]
	if !ok {
		ret.Error(errors.New("attribute BackupFactor not found"))
		return ret
	}
	backupFactoInt, ok := backupFactor.(int32)
	if !ok {
		ret.Error(fmt.Errorf("backup factor cfg type error(expecting int): %+v", backupFactor))
		return ret
	}
	if backupFactoInt < 0 {
		ret.Error(fmt.Errorf("illegal backup factor : %d", localHostModInt))
		return ret
	}
	topo.backupFactor = backupFactor


	localHostSchema, ok := attrsMap["LocalHostSchema"]
	if !ok {
		ret.Error(errors.New("attribute LocalHostSchema not found"))
		return ret
	}
	localHostSchemaInt, ok := localHostSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("local host class cfg type error(expecting int32): %+v", localHostSchema))
		return ret
	}
	if localHostSchemaInt <= 0 {
		ret.Error(fmt.Errorf("illegal LocalHostSchema: %d", localHostSchemaInt))
		return ret
	}
	topo.localHostSchema = localHostSchemaInt

	remoteHostSchema, ok := attrsMap["RemoteHostSchema"]
	if !ok {
		ret.Error(errors.New("attribute RemoteHostSchema not found"))
		return ret
	}
	remoteHostClassInt, ok := remoteHostSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("remote host class cfg type error(expecting int): %+v", remoteHostSchema))
		return ret
	}
	if remoteHostClassInt <= 0 {
		ret.Error(fmt.Errorf("illegal RemoteHostSchema: %d", remoteHostSchema))
		return ret
	}
	topo.remoteHostSchema = remoteHostClassInt

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
	topo.remoteNum = len(remoteEntriesMap)

	if topo.localHostMod != topo.remoteNum {
		ret.Error(fmt.Errorf("local offset(%d) != total entry num - 1(%d)", topo.localHostMod, topo.remoteNum))
		return ret
	}

	for i := 0; i < topo.remoteNum; i++ {
		remoteHostCfg, ok := cfg.Hosts[topo.remoteHostSchema]
		if !ok {
			ret.Error(fmt.Errorf("no remote host cfg found: %d", topo.remoteHostSchema))
			return ret
		}
		topo.remoteHosts = append(
			topo.remoteHosts, host.GetHostPrototype(remoteHostCfg.Class).Right().(config.IOC).New(remoteHostCfg.Attributes, cfg).(host.MaybeHost).Right())
	}

	for i:=0;int64(i)<topo.totalHostNum;i++ {
		mod := i%(topo.remoteNum + 1)
		if mod >= topo.localHostMod && mod < topo.localHostMod + topo.backupFactor {
			localHostCfg, ok := cfg.Hosts[topo.localHostSchema]
			if !ok {
				ret.Error(fmt.Errorf("no local host cfg found: %d", topo.localHostSchema))
				return ret
			}
			localHost := host.GetHostPrototype(localHostCfg.Class).Right().New(localHostCfg.Attributes, cfg).(host.MaybeHost).Right()
			localHost.SetId(int64(i))
			topo.localHosts = append(topo.localHosts, localHost)
		}
	}

	ret.Value(topo)
	return ret
}
