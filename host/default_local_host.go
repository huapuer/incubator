package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"unsafe"
	"../persistence"
	"../context"
	"../serialization"
	"fmt"
)

const (
	defaultLocalHostClassName  = "actor.defaultLocalHost"
)

func init() {
	RegisterHostPrototype(defaultLocalHostClassName, &defaultLocalHost{}).Test()
}

type defaultLocalHost struct {
	commonHost
}

func (this *defaultLocalHost) Receive(msg message.Message) (err maybe.MaybeError) {
	message.Route(msg).Test()
	return
}

func (this defaultLocalHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	//TODO: real logic
	ret.Value(&defaultLocalHost{
		commonHost{
			valid:true,
		},
	})
	return ret
}

func (this *defaultLocalHost) GetJsonBytes() (ret maybe.MaybeBytes) {
	ret.Value([]byte{})
	return
}

func (this *defaultLocalHost) SetJsonField(data []byte) (err maybe.MaybeError) {
	err.Error(nil)
	return
}

func (this *defaultLocalHost) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}

func (this defaultLocalHost) Duplicated() (ret MaybeHost) {
	ret.Value(&defaultLocalHost{
		commonHost{
			id: this.id,
			valid:true,
		},
	})
	return ret
}

func (this *defaultLocalHost) FromPersistenceAync(ctx context.HostRecoverContext, space string, layer int32, id int64) {
	go func(){
		select {
		case <-ctx.Ctx.Done():
			return
		}
		maybe.TryCatch(
			func(){
				content := persistence.FromPersistence(space, layer, defaultLocalHostClassName, id).Right()
				ret := MaybeHost{}
				new := this.Duplicated().Right()
				serialization.Unmarshal(content, new).Test()
				new.SetId(id).Test()
				ret.Value(new)
				ctx.Ret <- ret
			},
			func(err error){
				ret := MaybeHost{}
				ret.Error(err)
				ctx.Ret <- ret

			})
	}()
}

func (this *defaultLocalHost) ToPersistence(space string, layer int32) (err maybe.MaybeError) {
	if this.id <= 0 {
		err.Error(fmt.Errorf("illegal host id: %d", this.id))
		return
	}
	content := serialization.Marshal(this)
	return persistence.ToPersistence(space, layer, defaultLocalHostClassName, this.id, content)
}

func (this *defaultLocalHost) ToPersistenceAync(ctx context.AyncErrorContext, space string, layer int32) {
	go func(){
		select {
		case <-ctx.Ctx.Done():
			return
		}
		if this.id <= 0 {
			err := maybe.MaybeError{}
			err.Error(fmt.Errorf("illegal host id: %d", this.id))
			ctx.Err <- err
			return
		}
		content := serialization.Marshal(this)
		ctx.Err <- persistence.ToPersistence(space, layer, defaultLocalHostClassName, this.id, content)
	}()
}


