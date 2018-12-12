package interfaces

type IOC interface {
	New(interface{}, Config) IOC
}
