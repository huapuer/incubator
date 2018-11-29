package topo

import (
	"../common/maybe"
	"../config"
	"../host"
	"../persistence"
	"../serialization"
	"../storage"
	"errors"
	"fmt"
	"../message"
	"unsafe"
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
	linkRadius       int64
	localHostMod     int32
	backupFactor     int32
	localHostSchema  int32
	linkSchema       int32
	remoteHostSchema int32
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

	attrsMap, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal cfg type when new layer %s", defaultTopoClassName))
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

	linkRadius, ok := attrsMap["LinkRadius"]
	if !ok {
		ret.Error(errors.New("attribute LinkRadius not found"))
		return ret
	}
	linkRadiusInt, ok := linkRadius.(int64)
	if !ok {
		ret.Error(fmt.Errorf("link radius cfg type error(expecting int): %+v", linkRadius))
		return ret
	}
	if linkRadiusInt <= 0 {
		ret.Error(fmt.Errorf("illegal link radius: %d", linkRadiusInt))
		return ret
	}
	topo.linkRadius = linkRadiusInt

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
	topo.backupFactor = backupFactoInt

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
	topo.remoteNum = int32(len(remoteEntriesMap))

	if topo.localHostMod != topo.remoteNum {
		ret.Error(fmt.Errorf("local offset(%d) != total entry num - 1(%d)", topo.localHostMod, topo.remoteNum))
		return ret
	}

	for i := int32(0); i < topo.remoteNum; i++ {
		remoteHostCfg, ok := cfg.Hosts[topo.remoteHostSchema]
		if !ok {
			ret.Error(fmt.Errorf("no remote host cfg found: %d", topo.remoteHostSchema))
			return ret
		}
		topo.remoteHosts = append(
			topo.remoteHosts, host.GetHostPrototype(remoteHostCfg.Class).Right().(config.IOC).New(remoteHostCfg.Attributes, cfg).(host.MaybeHost).Right())
	}

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

	linkAttr := linkCfg.Attributes
	if linkAttr == nil {
		ret.Error(errors.New("link attr is nil"))
		return ret
	}
	linkAttrMap, ok := linkAttr.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal link attr type when new layer %s", defaultTopoClassName))
		return ret
	}
	linkSparseEntries, ok := linkAttrMap["SparseEnties"]
	if !ok {
		ret.Error(errors.New("link sparse entry cfg not found"))
		return ret
	}
	linkSparseEntriesArray, ok := linkSparseEntries.([]map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("link sparse entries cfg type error(expecting []map[stirng]interface{}): %+v", linkSparseEntries))
		return ret
	}
	entries := make([]*storage.SparseEntry, 0, 0)
	for i, entryCfg := range linkSparseEntriesArray {
		keyTo, ok := entryCfg["KeyTo"]
		if !ok {
			ret.Error(fmt.Errorf("sparse entry KeyTo attr not found, index: %d", i))
			return ret
		}
		keyToInt, ok := keyTo.(int64)
		if !ok {
			ret.Error(fmt.Errorf("illegal KeyTo type(expecting int64): %v, index: %d", keyTo, i))
			return ret
		}

		offset, ok := entryCfg["Offset"]
		if !ok {
			ret.Error(fmt.Errorf("sparse entry Offset attr not found, index: %d", i))
			return ret
		}
		offsetInt, ok := offset.(int64)
		if !ok {
			ret.Error(fmt.Errorf("illegal Offset type(expecting int64): %v, index: %d", offset, i))
			return ret
		}

		size, ok := entryCfg["Size"]
		if !ok {
			ret.Error(fmt.Errorf("sparse entry Size attr not found, index: %d", i))
			return ret
		}
		sizeInt, ok := size.(int64)
		if !ok {
			ret.Error(fmt.Errorf("illegal Size type(expecting int64): %v, index: %d", size, i))
			return ret
		}

		hashDepth, ok := entryCfg["HashDepth"]
		if !ok {
			ret.Error(fmt.Errorf("sparse entry HashDepth attr not found, index: %d", i))
			return ret
		}
		hashDepthInt, ok := hashDepth.(int32)
		if !ok {
			ret.Error(fmt.Errorf("illegal HashDepth type(expecting int32): %v, index: %d", hashDepth, i))
			return ret
		}

		entries = append(entries, &storage.SparseEntry{
			KeyTo:      keyToInt,
			Offset:     offsetInt,
			Size:       sizeInt,
			HashStride: hashDepthInt,
		})
	}

	linkDenseSize, ok := linkAttrMap["DenseSize"]
	if !ok {
		ret.Error(errors.New("link DenseSize attr not found"))
		return ret
	}
	linkDenseSizeInt, ok := linkDenseSize.(int64)
	if !ok {
		ret.Error(fmt.Errorf("illegal link DenseSize attr type(expecting int64): %v", linkDenseSize))
		return ret
	}

	linkHashDepth, ok := linkAttrMap["HashDepth"]
	if !ok {
		ret.Error(errors.New("link HashDepth attr not found"))
		return ret
	}
	linkHashDepthInt, ok := linkHashDepth.(int32)
	if !ok {
		ret.Error(fmt.Errorf("illegal link HashDepth attr type(expecting int32): %v", linkHashDepth))
		return ret
	}

	if cfg.Layer.Recover {
		topo.localHosts = storage.NewDenseTable(
			topo.localHostCanon.(storage.DenseTableElement),
			topo.totalHostNum*int64(topo.backupFactor),
			topo.totalHostNum/int64(topo.remoteNum),
			entries,
			topo.localHostCanon.(host.Host).GetSize(),
			linkHashDepthInt,
			persistence.FromPersistence(
				cfg.Layer.Space,
				cfg.Layer.Id,
				LinkPersistentClass,
				0).Right()).Right()
	} else {
		topo.links = storage.NewDenseTable(
			topo.linkCanon.(storage.DenseTableElement),
			topo.totalHostNum*int64(topo.backupFactor),
			linkDenseSizeInt,
			entries,
			topo.linkCanon.(host.Host).GetSize(),
			linkHashDepthInt,
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
	idx := int32(id/int64(this.remoteNum)/int64(this.backupFactor+1)) + mod
	if idx > int32(this.localHosts.ElemLen()) {
		err.Error(fmt.Errorf("master id exceeds local host range: %d", id))
		return
	}

	hosts := make([]host.Host, HostSlotSize, 0)

	if mod == this.localHostMod {
		h := this.localHostCanon
		ptr := this.localHosts.Get(0, int64(idx)).Right()
		serialization.Ptr2IFace(&h, ptr)
		hosts = append(hosts, h)
	} else {
		hosts = append(hosts, this.remoteHosts[mod])
	}
	if mod < this.localHostMod+this.backupFactor {
		h := this.localHostCanon
		ptr := this.localHosts.Get(0, int64(idx)).Right()
		serialization.Ptr2IFace(&h, ptr)
		hosts = append(hosts, h)
	}
	for offset := int32(0); offset < this.backupFactor-1; offset++ {
		ridx := (mod + offset) % (this.remoteNum)
		hosts = append(hosts, this.remoteHosts[ridx])
	}

	masterSended := false
	for _, h := range hosts {
		if h.IsValid() {
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
		err.Error(fmt.Errorf("no available master host found for id: %d", id))
		return
	}

	err.Error(nil)
	return
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
		if h.IsValid() {
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

func (this defaultTopo) TraverseLinksOfHost(hid int64, callback func(ptr unsafe.Pointer) bool) (err maybe.MaybeError) {
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
