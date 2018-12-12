package host

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"net"
)

var (
	hostPrototypes = make(map[string]interfaces.Host)
)

func RegisterHostPrototype(name string, val interfaces.Host) (err maybe.MaybeError) {
	if _, ok := hostPrototypes[name]; ok {
		err.Error(fmt.Errorf("host prototype redefined: %s", name))
		return
	}
	hostPrototypes[name] = val
	return
}

func GetHostPrototype(name string) (ret interfaces.MaybeHost) {
	if prototype, ok := hostPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("host prototype for class not found: %s", name))
	return
}

type commonHost struct {
	id int64
}

func (this commonHost) GetId() int64 {
	return this.id
}

func (this *commonHost) SetId(id int64) {
	this.id = id
	return
}

func (this commonHost) SetIP(string) {
	panic("not implemented")
}

func (this commonHost) SetPort(int) {
	panic("not implemented")
}

func (this commonHost) Start() maybe.MaybeError {
	panic("not implemented")
}

type commonLinkHost struct {
	guestId int64
}

func (this *commonLinkHost) GetGuestId() int64 {
	return this.guestId
}

func (this *commonLinkHost) SetGuestId(id int64) {
	this.guestId = id
	return
}

type SessionHost interface {
	interfaces.Host

	SetPeer(net.Conn)
	Replicate() interfaces.MaybeHost
}
