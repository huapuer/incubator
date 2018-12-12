package network

import (
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"net"
	"time"
)

type connectionPool struct {
	addr    string
	maxIdle int32
	maxBusy int32
	timeout time.Duration
	idle    chan net.Conn
	busy    chan struct{}
}

type MaybeConnectionPool struct {
	maybe.MaybeError

	value connectionPool
}

func (this MaybeConnectionPool) Value(value connectionPool) {
	this.Error(nil)
	this.value = value
}

func (this MaybeConnectionPool) Right() connectionPool {
	this.Test()
	return this.value
}

func NewConnectionPool(addr string, maxIdle int32, maxBusy int32, timeout time.Duration) (ret MaybeConnectionPool) {
	if addr == "" {
		ret.Error(errors.New("addr is empty"))
		return
	}
	if maxIdle < 0 {
		ret.Error(fmt.Errorf("illegal max spare(expecting >=0): %d", maxIdle))
		return
	}
	if maxBusy <= 0 {
		ret.Error(fmt.Errorf("illegal max busy(expecting >0): %d", maxBusy))
		return
	}

	ret.Value(connectionPool{
		addr:    addr,
		maxIdle: maxIdle,
		maxBusy: maxBusy,
		timeout: timeout,
		idle:    make(chan net.Conn, maxIdle),
		busy:    make(chan struct{}, maxBusy),
	})
	return
}

type MaybeConnection struct {
	maybe.MaybeError

	value net.Conn
}

func (this MaybeConnection) Value(value net.Conn) {
	this.Error(nil)
	this.value = value
}

func (this MaybeConnection) Right() net.Conn {
	this.Test()
	return this.value
}

func (this connectionPool) GetConnection() (ret MaybeConnection) {
	select {
	case c := <-this.idle:
		ret.Value(c)
		return
	case this.busy <- struct{}{}:
		conn, e := net.Dial("tcp", this.addr)
		if e != nil {
			ret.Error(e)
			return
		}
		ret.Value(conn)
		return
	case <-time.After(this.timeout):
		ret.Error(errors.New("get conn time out"))
		return
	}
}

func (this connectionPool) ReleaseConnection(conn net.Conn) (err maybe.MaybeError) {
	select {
	case this.idle <- conn:
		return
	case <-this.busy:
		conn.Close()
		return
	case <-time.After(this.timeout):
		err.Error(errors.New("release conn time out"))
		return
	}
}
