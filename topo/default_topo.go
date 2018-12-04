package topo

import (
	"../common/maybe"
	"../config"
	"../host"
	"../persistence"
	"../serialization"
	"../storage"
	"fmt"
	"../message"
	"unsafe"
	"math/rand"
)

const (
	defaultTopoClassName = "topo.defaultTopo"

	LocalHostPersistentClass = "LOCALHOST"
	LinkPersistentClass      = "LINK"

	HostSlotSize = 10
)

func init() {
	RegisterTopoPrototype(defaultTopoClassName, &defaultTopo{}).Test()
}

type defaultTopo struct {
	totalHostNum     int64
	totalRemoteHostNum int32
	linkRadius       int64
	localHostMod     int32
	backupFactor     int32
	localHostSchema  int32
	linkSchema       int32
	localHosts       storage.DenseTable
	links            storage.DenseTable
	remoteHosts      []host.Host
	remoteNum        int32
	localHostCanon   host.Host
	linkCanon        host.Host
}

func (this defaultTopo) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeTopo{}
	topo := &defaultTopo{
		remoteHosts: make([]host.Host, 0, 0),
	}

	topo.totalHostNum = config.GetAttrInt64(attrs, "TotalHostNum", config.CheckInt64GT0).Right()
	topo.linkRadius = config.GetAttrInt64(attrs, "LinkRadius", config.CheckInt64GT0).Right()
	topo.localHostMod = config.GetAttrInt32(attrs, "LocalHostMod", config.CheckInt32GT0).Right()
	topo.backupFactor = config.GetAttrInt32(attrs, "BackupFactor", config.CheckInt32GT0).Right()
	topo.localHostSchema = config.GetAttrInt32(attrs, "LocalHostSchema", config.CheckInt32GT0).Right()

	remoteEntries := config.GetAttrMapEfaceArray(attrs, "RemoteEntries").Right().([]map[string]interface{})
	topo.remoteNum = int32(len(remoteEntries))

	if topo.localHostMod != topo.remoteNum {
		ret.Error(fmt.Errorf("local offset(%d) != total entry num - 1(%d)", topo.localHostMod, topo.remoteNum))
		return ret
	}

	for i, entry := range remoteEntries {
		remoteHostClass := config.GetAttrString(entry, "RemoteHostClass", config.CheckStringNotEmpty).Right()
		remoteHostAttr := config.GetAttrMapEface(entry, "Attributes").Right()

		h :=  host.GetHostPrototype(remoteHostClass).Right().New(remoteHostAttr, cfg).(host.MaybeHost).Right()
		h.SetId(int64(i))
		h.(host.HealthManager).Start()

		topo.remoteHosts = append(topo.remoteHosts, h)
	}
	topo.totalRemoteHostNum = int32(len(remoteEntries))

	localHostCfg, ok := cfg.Hosts[topo.localHostSchema]
	if !ok {
		ret.Error(fmt.Errorf("no local host cfg found: %d", topo.localHostSchema))
		return ret
	}
	if topo.totalHostNum/int64(topo.remoteNum)*int64(topo.remoteNum) != topo.totalHostNum {
		ret.Error(fmt.Errorf("total host num is not times of remote num: %d / %d", topo.totalHostNum, topo.remoteNum))
		return ret
	}
	topo.localHostCanon = host.GetHostPrototype(localHostCfg.Class).Right()

	if cfg.Layer.Recover {
		topo.localHosts = storage.NewDenseTable(
			topo.localHostCanon.(storage.DenseTableElement),
			1,
			topo.totalHostNum/int64(topo.remoteNum),
			[]*storage.SparseEntry{},
			topo.localHostCanon.(host.Host).GetSize(),
			0,
			persistence.FromPersistence(
				cfg.Layer.Space,
				cfg.Layer.Id,
				LocalHostPersistentClass,
				0).Right()).Right()
	} else {
		topo.localHosts = storage.NewDenseTable(
			topo.localHostCanon.(storage.DenseTableElement),
			1,
			topo.totalHostNum/int64(topo.remoteNum),
			[]*storage.SparseEntry{},
			topo.localHostCanon.(host.Host).GetSize(),
			0,
			nil).Right()

		for i := int64(0); i < topo.totalHostNum; i++ {
			mod := int32(i) % topo.remoteNum
			if mod >= topo.localHostMod && mod < topo.localHostMod+topo.backupFactor {
				localHost := topo.localHostCanon.New(localHostCfg.Attributes, cfg).(host.MaybeHost).Right()
				localHost.SetId(int64(i))
				topo.localHosts.Put(0, int64(i), serialization.IFace2Ptr(&localHost))
			}
		}
	}

	linkCfg, ok := cfg.Links[topo.linkSchema]
	if !ok {
		ret.Error(fmt.Errorf("no link cfg found: %d", topo.localHostSchema))
		return ret
	}
	topo.linkCanon = host.GetHostPrototype(linkCfg.Class).Right()

	linkSparseEntries := config.GetAttrMapEfaceArray(linkCfg.Attributes, "SparseEnties").Right().([]map[string]interface{})

	entries := make([]*storage.SparseEntry, 0, 0)
	for _, entryCfg := range linkSparseEntries {
		keyTo:=config.GetAttrInt64(entryCfg, "KeyTo", config.CheckInt64GT0).Right()
		offset:=config.GetAttrInt64(entryCfg, "Offset", config.CheckInt64GT0).Right()
		size:=config.GetAttrInt64(entryCfg, "Size", config.CheckInt64GT0).Right()
		hashDepth:=config.GetAttrInt32(entryCfg, "Offset", config.CheckInt32GET0).Right()

		entries = append(entries, &storage.SparseEntry{
			KeyTo:      keyTo,
			Offset:     offset,
			Size:       size,
			HashStride: hashDepth,
		})
	}

	linkDenseSize := config.GetAttrInt64(linkCfg.Attributes, "DenseSize", config.CheckInt64GT0).Right()
	linkHashDepth:=config.GetAttrInt32(linkCfg.Attributes, "HashDepth", config.CheckInt32GET0).Right()

	if cfg.Layer.Recover {
		topo.localHosts = storage.NewDenseTable(
			topo.localHostCanon.(storage.DenseTableElement),
			topo.totalHostNum*int64(topo.backupFactor),
			topo.totalHostNum/int64(topo.remoteNum),
			entries,
			topo.localHostCanon.(host.Host).GetSize(),
			linkHashDepth,
			persistence.FromPersistence(
				cfg.Layer.Space,
				cfg.Layer.Id,
				LinkPersistentClass,
				0).Right()).Right()
	} else {
		topo.links = storage.NewDenseTable(
			topo.linkCanon.(storage.DenseTableElement),
			topo.totalHostNum*int64(topo.backupFactor),
			linkDenseSize,
			entries,
			topo.linkCanon.(host.Host).GetSize(),
			linkHashDepth,
			nil).Right()

		//TODO: init links
	}

	//TODO: add potential link

	ret.Value(topo)
	return ret
}

//go:noescape
func (this defaultTopo) SendToHost(id int64, msg message.RemoteMessage) (err maybe.MaybeError) {
	mod := int32(id % (int64(this.remoteNum)))

	hosts := make([]host.Host, HostSlotSize, 0)
	for offset := int32(0); offset < this.backupFactor-1; offset++ {
		ridx := (mod + offset) % (this.remoteNum)
		hosts = append(hosts, this.remoteHosts[ridx])
	}

	masterSended := false
	for i, h := range hosts {
		if !h.IsHealth() {
			continue
		}

		realhost := h
		if mod + int32(i) == this.localHostMod || mod + int32(i) < this.localHostMod+this.backupFactor {
			idx := int32(id/int64(this.remoteNum)/int64(this.backupFactor+1)) + mod
			if idx > int32(this.localHosts.ElemLen()) {
				err.Error(fmt.Errorf("master id exceeds local host range: %d", id))
				return
			}

			h := this.localHostCanon
			ptr := this.localHosts.Get(0, int64(idx)).Right()
			serialization.Ptr2IFace(&h, ptr)

			realhost = h
		}
		if !masterSended {
			msg.Master(message.MASTER_YES)
			masterSended = true
		} else {
			msg.Master(message.MASTER_NO)
		}
		realhost.Receive(msg).Test()
	}

	if !masterSended {
		err.Error(fmt.Errorf("no available master host found for id: %d", id))
		return
	}

	err.Error(nil)
	return
}

//go:noescape
func (this defaultTopo) LookupHost(id int64) (ret host.MaybeHost) {
	mod := int32(id % (int64(this.remoteNum)))

	if mod == this.localHostMod {
		idx := int32(id/int64(this.remoteNum)/int64(this.backupFactor+1)) + mod
		if idx > int32(this.localHosts.ElemLen()) {
			ret.Error(fmt.Errorf("master id exceeds local host range: %d", id))
			return
		}

		h := this.localHostCanon
		ptr := this.localHosts.Get(0, int64(idx)).Right()
		serialization.Ptr2IFace(&h, ptr)
		ret.Value(h)
		return
	} else {
		ret.Value(this.remoteHosts[mod])
		return
	}
}

//go:noescape
func (this defaultTopo) SendToLink(hid int64, gid int64, msg message.RemoteMessage) (err maybe.MaybeError) {
	mod := int32(hid % (int64(this.remoteNum)))
	blk := int32(hid/int64(this.remoteNum)/int64(this.backupFactor+1)) + mod

	if blk > int32(this.links.BlockLen()) {
		err.Error(fmt.Errorf("master host id exceeds local host range: %d", hid))
		return
	}

	idx := gid - hid - 1
	if idx < -this.linkRadius || idx > this.linkRadius {
		err.Error(fmt.Errorf("master link id exceeds link range: %d", idx))
		return
	}
	idx += this.linkRadius

	hosts := make([]host.Host, HostSlotSize, 0)

	if mod == this.localHostMod {
		h := this.linkCanon
		ptr := this.links.Get(int64(blk), int64(blk)).Right()
		serialization.Ptr2IFace(&h, ptr)
		hosts = append(hosts, h)
	} else {
		hosts = append(hosts, this.remoteHosts[mod])
	}
	if mod < this.localHostMod+this.backupFactor {
		h := this.linkCanon
		ptr := this.links.Get(int64(blk), int64(blk)).Right()
		serialization.Ptr2IFace(&h, ptr)
		hosts = append(hosts, h)
	}
	for offset := int32(0); offset < this.backupFactor-1; offset++ {
		ridx := (mod + offset) % (this.remoteNum)
		hosts = append(hosts, this.remoteHosts[ridx])
	}

	masterSended := false
	for _, h := range hosts {
		if h.IsHealth() {
			if !masterSended {
				msg.Master(message.MASTER_YES)
				masterSended = true
			} else {
				msg.Master(message.MASTER_NO)
			}
			h.Receive(msg).Test()
		}
	}

	if !masterSended {
		err.Error(fmt.Errorf("no available master link found for id: %d->%d", hid, gid))
		return
	}

	err.Error(nil)
	return
}

func (this defaultTopo) TraverseOutLinksOfHost(hid int64, callback func(ptr unsafe.Pointer) bool) (err maybe.MaybeError) {
	mod := int32(hid % (int64(this.remoteNum)))
	if mod == this.localHostMod {
		blk := int32(hid/int64(this.remoteNum)/int64(this.backupFactor+1)) + mod
		if blk > int32(this.links.BlockLen()) {
			err.Error(fmt.Errorf("master host id exceeds local host range: %d", hid))
			return
		}

		this.links.TraverseBlock(int64(blk), callback)
	} else {
		err.Error(fmt.Errorf("host id exceeds local host range: %d", hid))
		return
	}

	err.Error(nil)
	return
}

func (this defaultTopo) GetRemoteHosts() []host.Host {
	return this.remoteHosts
}

func (this *defaultTopo) AddHost(host.Host) maybe.MaybeError {
	panic("not implemented")
}

func (this defaultTopo) GetRemoteHostId(idx int32) int64 {
	return int64(idx) + int64(this.totalRemoteHostNum)*rand.Int63n(this.totalHostNum / int64(this.totalRemoteHostNum))
}
