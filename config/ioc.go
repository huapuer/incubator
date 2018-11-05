package config

type IOC interface {
	New(Config, ...int32) IOC
}
