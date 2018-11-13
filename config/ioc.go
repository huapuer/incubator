package config

type IOC interface {
	New(interface{}, Config) IOC
}
