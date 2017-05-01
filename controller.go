package kcache

import (
	"context"

	"github.com/boz/kcache/lifecycle"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CacheReader interface {
	GetObject(obj metav1.Object) (metav1.Object, error)
	List() ([]metav1.Object, error)
}

type Publisher interface {
	Subscribe() Subscription
}

type CacheController interface {
	Cache() CacheReader
	Initialized() <-chan struct{}
}

type Controller interface {
	CacheController
	Publisher
	Stop()
}

type Subscription interface {
	Events() <-chan Event
	Close()
}

type Event interface {
	Type() EventType
}

type EventType string

func NewController() Controller {
	c := &controller{
		initializedch: make(chan struct{}),
		lifecycle:     lifecycle.NewLifecycle(),
		ctx:           context.Background(),
	}
	go c.lifecycle.WatchContext(c.ctx)

	go c.run()

	return c
}

type controller struct {

	// closed when initialization complete
	initializedch chan struct{}

	lifecycle lifecycle.Lifecycle

	ctx context.Context
}

func (c *controller) Initialized() <-chan struct{} {
	return c.initializedch
}

func (c *controller) Stop() {
	c.lifecycle.Shutdown()
}

func (c *controller) Cache() CacheReader {
	return nil
}

func (c *controller) Subscribe() Subscription {
	return nil
}

func (c *controller) run() {
	defer c.lifecycle.ShutdownCompleted()
	for {
		select {
		case <-c.lifecycle.ShutdownRequest():
			c.lifecycle.ShutdownInitiated()
			return
		}
	}
}
