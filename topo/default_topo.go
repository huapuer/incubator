package topo

import (
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/global"
	"github.com/incubator/host"
	"github.com/incubator/interfaces"
	"github.com/incubator/message"
	"github.com/incubator/persistence"
	"github.com/incubator/serialization"
	"github.com/incubator/storage"
	"math/rand"
	"net"
	"unsafe"
)

const (
	defaultTopoClassName = "topo.defaultTopo"

	LocalHostPersistentClass = "LOCALHOST"
	LinkPersistentClass      = "LINK"

	HostSlotSize = 10
)

func init() {
	interfaces.RegisterTopoPrototype(defaultTopoClassName, &defaultTopo{}).Test()
}

type defaultTopo struct {
	persistence.CommomPersistentable

	layer              int32
	totalHostNum       int64
	totalRemoteHostNum int32
	linkRadius         int64
	localHostMod       int32
	backupFactor       int32
	localHostSchema    int32
	linkSchema         int32
	localHosts         storage.DenseTable
	links              storage.DenseTable
	remoteHosts        []interfaces.Host
	remoteNum          int32
	localHostCanon     interfaces.Host
	linkCanon          interfaces.Host
	addr               string
}

func getIntranetIp() (ret maybe.MaybeString) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		ret.Error(err)
		return
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ret.Value(ipnet.IP.String())
				return
			}
		}
	}

	ret.Error(errors.New("internal ip not found"))
	return
}

func (this defaultTopo) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	interalIP := getIntranetIp().Right()

	ret := interfaces.MaybeTopo{}
	topo := &defaultTopo{
		remoteHosts: make([]interfaces.Host, 0, 0),
	}

	topo.totalHostNum = config.GetAttrInt64(attrs, "TotalHostNum", config.CheckInt64GT0).Right()
	topo.linkRadius = config.GetAttrInt64(attrs, "LinkRadius", config.CheckInt64GT0).Right()
	topo.backupFactor = config.GetAttrInt32(attrs, "BackupFactor", config.CheckInt32GT0).Right()
	topo.localHostSchema = config.GetAttrInt32(attrs, "LocalHostSchema", config.CheckInt32GT0).Right()
	topo.linkSchema = config.GetAttrInt32(attrs, "LinkSchema", config.CheckInt32GT0).Right()

	remoteEntries := config.GetAttrMapEfaceArray(attrs, "RemoteEntries").Right().([]map[string]interface{})
	topo.remoteNum = int32(len(remoteEntries))

	if topo.localHostMod != topo.remoteNum {
		ret.Error(fmt.Errorf("local offset(%d) != total entry num - 1(%d)", topo.localHostMod, topo.remoteNum))
		return ret
	}

	topo.localHostMod = -1
	for i, entry := range remoteEntries {
		remoteHostSchema := config.GetAttrInt32(entry, "RemoteHostSchema", config.CheckInt32GT0).Right()
		remoteHostAttr := config.GetAttrMapEface(entry, "Attributes").Right()

		ip := config.GetAttrString(remoteHostAttr, "IP", config.CheckStringNotEmpty).Right()
		port := config.GetAttrInt(remoteHostAttr, "Port", config.CheckIntGT0).Right()
		if ip == interalIP && global.IsInListenedPorts(port) {
			topo.localHostMod = int32(i)
			topo.addr = fmt.Sprint("%s:%d", ip, port)
		}

		remoteHostCfg, ok := cfg.(*config.Config).HostMap[remoteHostSchema]
		if !ok {
			ret.Error(fmt.Errorf("no remote host cfg found: %d", remoteHostSchema))
			return ret
		}

		h := host.GetHostPrototype(remoteHostCfg.Class).Right().New(remoteHostAttr, cfg).(interfaces.MaybeHost).Right()
		h.SetId(int64(i))
		h.SetIP(ip)
		h.SetPort(port)
		h.Start().Test()

		topo.remoteHosts = append(topo.remoteHosts, h)
	}
	if topo.localHostMod == -1 {
		ret.Error(errors.New("internal IP not in remote hosts IPs"))
		return ret
	}
	topo.totalRemoteHostNum = int32(len(remoteEntries))

	localHostCfg, ok := cfg.(*config.Config).HostMap[topo.localHostSchema]
	if !ok {
		ret.Error(fmt.Errorf("no local host cfg found: %d", topo.localHostSchema))
		return ret
	}
	if topo.totalHostNum/int64(topo.remoteNum)*int64(topo.remoteNum) != topo.totalHostNum {
		ret.Error(fmt.Errorf("total host num is not times of remote num: %d / %d", topo.totalHostNum, topo.remoteNum))
		return ret
	}
	topo.localHostCanon = host.GetHostPrototype(localHostCfg.Class).Right()

	switch cfg.(*config.Config).Layer.StartMode {
	case config.LAYER_START_MODE_RECOVER:
		l := interfaces.GetLayer(cfg.(*config.Config).Layer.Id).Right()
		topo.localHosts = storage.NewDenseTable(
			topo.localHostCanon.(interfaces.DenseTableElement),
			1,
			topo.totalHostNum/int64(topo.remoteNum),
			[]*storage.SparseEntry{},
			topo.localHostCanon.(interfaces.Host).GetSize(),
			0,
			persistence.FromPersistence(
				persistence.FROM_PERSISTENCE_MODE_RECOVER,
				topo.GetLoadExpiration(),
				l.GetVersion(),
				cfg.(*config.Config).Layer.Space,
				cfg.(*config.Config).Layer.Id,
				LocalHostPersistentClass).Right()).Right()
	case config.LAYER_START_MODE_REBOOT:
		topo.localHosts = storage.NewDenseTable(
			topo.localHostCanon.(interfaces.DenseTableElement),
			1,
			topo.totalHostNum/int64(topo.remoteNum),
			[]*storage.SparseEntry{},
			topo.localHostCanon.(interfaces.Host).GetSize(),
			0,
			persistence.FromPersistence(
				persistence.FROM_PERSISTENCE_MODE_RECOVER,
				0,
				0,
				cfg.(*config.Config).Layer.Space,
				cfg.(*config.Config).Layer.Id,
				LocalHostPersistentClass).Right()).Right()
	case config.LAYER_START_MODE_NEW:
		topo.localHosts = storage.NewDenseTable(
			topo.localHostCanon.(interfaces.DenseTableElement),
			1,
			topo.totalHostNum/int64(topo.remoteNum),
			[]*storage.SparseEntry{},
			topo.localHostCanon.(interfaces.Host).GetSize(),
			0,
			nil).Right()

		for i := int64(0); i < topo.totalHostNum; i++ {
			mod := int32(i) % topo.remoteNum
			if mod >= topo.localHostMod && mod < topo.localHostMod+topo.backupFactor {
				localHost := topo.localHostCanon.New(localHostCfg.Attributes, cfg).(interfaces.MaybeHost).Right()
				localHost.SetId(int64(i))
				topo.localHosts.Put(0, int64(i), serialization.IFace2Ptr(&localHost))
			}
		}
	default:
		ret.Error(fmt.Errorf("unknown layer start mode: %d", cfg.(*config.Config).Layer.StartMode))
		return ret
	}

	linkCfg, ok := cfg.(*config.Config).LinkMap[topo.linkSchema]
	if !ok {
		ret.Error(fmt.Errorf("no link cfg found: %d", topo.localHostSchema))
		return ret
	}
	topo.linkCanon = host.GetHostPrototype(linkCfg.Class).Right()

	linkSparseEntries := config.GetAttrMapEfaceArray(attrs, "LinkSparseEntries").Right().([]map[string]interface{})

	entries := make([]*storage.SparseEntry, 0, 0)
	for _, entryCfg := range linkSparseEntries {
		keyTo := config.GetAttrInt64(entryCfg, "KeyTo", config.CheckInt64GT0).Right()
		size := config.GetAttrInt64(entryCfg, "Size", config.CheckInt64GT0).Right()
		hashStride := config.GetAttrInt32(entryCfg, "HashStride", config.CheckInt32GET0).Right()

		entries = append(entries, &storage.SparseEntry{
			KeyTo:      keyTo,
			Size:       size,
			HashStride: hashStride,
		})
	}

	linkDenseSize := config.GetAttrInt64(attrs, "LinkDenseSize", config.CheckInt64GT0).Right()
	linkHashDepth := config.GetAttrInt32(attrs, "LinkHashDepth", config.CheckInt32GET0).Right()

	switch cfg.(*config.Config).Layer.StartMode {
	case config.LAYER_START_MODE_RECOVER:
		l := interfaces.GetLayer(cfg.(*config.Config).Layer.Id).Right()
		topo.localHosts = storage.NewDenseTable(
			topo.linkCanon.(interfaces.DenseTableElement),
			topo.totalHostNum*int64(topo.backupFactor),
			topo.totalHostNum/int64(topo.remoteNum),
			entries,
			topo.linkCanon.(interfaces.Host).GetSize(),
			linkHashDepth,
			persistence.FromPersistence(
				persistence.FROM_PERSISTENCE_MODE_RECOVER,
				topo.GetLoadExpiration(),
				l.GetVersion(),
				cfg.(*config.Config).Layer.Space,
				cfg.(*config.Config).Layer.Id,
				LinkPersistentClass).Right()).Right()
	case config.LAYER_START_MODE_REBOOT:
		topo.localHosts = storage.NewDenseTable(
			topo.linkCanon.(interfaces.DenseTableElement),
			topo.totalHostNum*int64(topo.backupFactor),
			topo.totalHostNum/int64(topo.remoteNum),
			entries,
			topo.linkCanon.(interfaces.Host).GetSize(),
			linkHashDepth,
			persistence.FromPersistence(
				persistence.FROM_PERSISTENCE_MODE_REBOOT,
				0,
				0,
				cfg.(*config.Config).Layer.Space,
				cfg.(*config.Config).Layer.Id,
				LinkPersistentClass).Right()).Right()
	case config.LAYER_START_MODE_NEW:
		topo.links = storage.NewDenseTable(
			topo.linkCanon.(interfaces.DenseTableElement),
			topo.totalHostNum*int64(topo.backupFactor),
			linkDenseSize,
			entries,
			topo.linkCanon.(interfaces.Host).GetSize(),
			linkHashDepth,
			nil).Right()
	default:
		ret.Error(fmt.Errorf("unknown layer start mode: %d", cfg.(*config.Config).Layer.StartMode))
		return ret
	}

	//TODO: init links

	//TODO: add potential link

	ret.Value(topo)
	return ret
}

func (this defaultTopo) Persistent() (err maybe.MaybeError) {
	l := interfaces.GetLayer(this.layer).Right()

	persistence.ToPersistence(
		this.GetStoreExpiration(),
		l.GetVersion(),
		l.GetConfig().(*config.Config).Layer.Space,
		this.layer,
		LocalHostPersistentClass,
		this.localHosts.GetBytes()).Test()

	persistence.ToPersistence(
		this.GetStoreExpiration(),
		l.GetVersion(),
		l.GetConfig().(*config.Config).Layer.Space,
		this.layer,
		LinkPersistentClass,
		this.links.GetBytes()).Test()

	err.Error(nil)
	return
}

func (this defaultTopo) GetRemoteHostId(idx int32) int64 {
	return int64(idx) + int64(this.totalRemoteHostNum)*rand.Int63n(this.totalHostNum/int64(this.totalRemoteHostNum))
}

////go:noescape
func (this defaultTopo) SendToHost(id int64, msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	mod := int32(id % (int64(this.remoteNum)))

	hosts := make([]interfaces.Host, HostSlotSize, 0)
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
		if mod+int32(i) == this.localHostMod || mod+int32(i) < this.localHostMod+this.backupFactor {
			idx := int32(id/int64(this.remoteNum)/int64(this.backupFactor+1)) + mod
			if idx > int32(this.localHosts.ElemLen()) {
				err.Error(fmt.Errorf("master id exceeds local host range: %d", id))
				return
			}

			h := this.localHostCanon
			ptr := this.localHosts.Get(0, int64(idx)).Right()
			serialization.Ptr2IFace(&h, ptr)

			realhost = h
			msg.SetHostId(h.GetId())
		} else {
			msg.SetHostId(this.GetRemoteHostId(int32(realhost.GetId())))
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

////go:noescape
func (this defaultTopo) LookupHost(id int64) (ret interfaces.MaybeHost) {
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

////go:noescape
func (this defaultTopo) SendToLink(hid int64, gid int64, msg interfaces.RemoteMessage) (err maybe.MaybeError) {
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

	hosts := make([]interfaces.Host, HostSlotSize, 0)

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

func (this defaultTopo) GetRemoteHosts() []interfaces.Host {
	return this.remoteHosts
}

func (this *defaultTopo) AddHost(interfaces.Host) maybe.MaybeError {
	panic("not implemented")
}

func (this defaultTopo) Start() {
	for _, h := range this.remoteHosts {
		if hm, ok := h.(interfaces.HealthManager); ok {
			hm.SetLayer(this.layer)
			hm.Start().Test()
		}
	}
}

func (this defaultTopo) GetLayer() int32 {
	return this.layer
}

func (this *defaultTopo) SetLayer(layer int32) {
	this.layer = layer
}

func (this defaultTopo) GetAddr() string {
	return this.addr
}
