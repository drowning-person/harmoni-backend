package iface

type Executor interface {
	Start() error
	Shutdown() error
}
