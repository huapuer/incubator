package global

import (
	"fmt"
	"incubator/common/maybe"
	"sync"
)

var (
	listenedPortsLock = sync.Mutex{}
	listenedPorts     = make(map[int]int)
)

func AddListenedPort(port int) (err maybe.MaybeError) {
	listenedPortsLock.Lock()
	if _, ok := listenedPorts[port]; ok {
		err.Error(fmt.Errorf("port already binded: %d", port))
		listenedPortsLock.Unlock()
		return
	}
	listenedPorts[port] = port
	listenedPortsLock.Unlock()

	err.Error(nil)
	return
}

func IsInListenedPorts(port int) bool {
	for _, p := range listenedPorts {
		if p == port {
			return true
		}
	}

	return false
}
