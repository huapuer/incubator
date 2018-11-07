package config

type IOC interface {
	New(Config) IOC
}
