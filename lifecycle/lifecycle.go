package lifecycle

import "context"

type Lifecycle interface {
	ShutdownRequest() <-chan struct{}
	ShutdownInitiated()
	ShuttingDown() <-chan struct{}
	ShutdownCompleted()
	Done() <-chan struct{}
	WatchContext(context.Context)
	Shutdown()
}

type lifecycle struct {
	stopch     chan struct{}
	stoppingch chan struct{}
	stoppedch  chan struct{}
}

func NewLifecycle() Lifecycle {
	return &lifecycle{
		stopch:     make(chan struct{}),
		stoppingch: make(chan struct{}),
		stoppedch:  make(chan struct{}),
	}
}

func (l *lifecycle) ShutdownRequest() <-chan struct{} {
	return l.stopch
}

func (l *lifecycle) ShutdownInitiated() {
	close(l.stoppingch)
}

func (l *lifecycle) ShuttingDown() <-chan struct{} {
	return l.stoppingch
}

func (l *lifecycle) ShutdownCompleted() {
	close(l.stoppedch)
}

func (l *lifecycle) Done() <-chan struct{} {
	return l.stoppedch
}

func (l *lifecycle) Shutdown() {
	for {
		select {
		case l.stopch <- struct{}{}:
		case <-l.stoppingch:
			<-l.stoppedch
			return
		}
	}
}

func (l *lifecycle) WatchContext(ctx context.Context) {
	var stopch chan struct{}
	donech := ctx.Done()
	for {
		select {
		case <-l.stoppingch:
			return
		case <-donech:
			donech = nil
			stopch = l.stopch
		case stopch <- struct{}{}:
			return
		}
	}
}
