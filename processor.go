package kcache

// process incomming events and generate event stream

type Processor interface {
	Sync([]metav1.Object)
	Add(metav1.Object)
	Update(metav1.Object)
	Delete(metav1.Object)
	Stop()
}
