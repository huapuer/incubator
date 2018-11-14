package network

import (
	"../common/maybe"
	"../message"
	"net"
)

type defaultClient struct {
	conn net.Conn
}

func (this *defaultClient) Start(addr string) (err maybe.MaybeError) {
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		err.Error(e)
	}
	return
	this.conn = conn
	return
}

func (this *defaultClient) Send(msg message.Message) (err maybe.MaybeError) {
	this.conn.Write(msg.Marshal(msg))
	return
}
