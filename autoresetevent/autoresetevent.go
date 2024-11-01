package autoresetevent

type AutoResetEvent struct {
	ch chan struct{}
}

func New() *AutoResetEvent {
	return &AutoResetEvent{make(chan struct{})}
}

func (ase *AutoResetEvent) WaitOne() {
	ase.ch <- struct{}{}
}

func (ase *AutoResetEvent) Set() {
	<-ase.ch
}
