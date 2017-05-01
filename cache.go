package kcache

type Cache interface {
	CacheReader
	Put(metav1.Object) error
	Import([]metav1.Object) error
	Stop()
}
