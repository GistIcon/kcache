package kcache

type Index interface {
	IndexReader
	Put(metav1.Object) error
	Delete(metav1.Object) error
	Import([]metav1.Object) error
}

type ImportConsumer interface {
	Sync([]metav1.Object) error
}

type Indexer interface {
}
