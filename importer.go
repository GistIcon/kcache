package kcache

type Importer interface {
	Initialized()
	Stop()
}
