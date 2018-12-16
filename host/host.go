package host

import (
	"github.com/incubator/common/maybe"
)

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
