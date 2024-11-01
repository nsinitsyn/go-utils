package semaphore

type Semaphore struct {
	ch chan struct{}
}

func New(degree int) *Semaphore {
	return &Semaphore{make(chan struct{}, degree)}
}

func (sem *Semaphore) Wait() {
	sem.ch <- struct{}{}
}

func (sem *Semaphore) Release() {
	<-sem.ch
}
