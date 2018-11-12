package config

type IOC interface {
	New(interface{}) IOC
}
