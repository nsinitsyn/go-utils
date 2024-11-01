package manualresetevent

type ManualResetEvent struct {
	ch       chan struct{}
	signaled bool
	reset    chan struct{}
}

func New() *ManualResetEvent {
	return &ManualResetEvent{make(chan struct{}), false, make(chan struct{})}
}

func (mre *ManualResetEvent) Wait() {
	if mre.signaled {
		return
	}
	mre.ch <- struct{}{}
}

func (mre *ManualResetEvent) Set() {
	mre.signaled = true
	go func() { // goroutine for unblocking Set
		for {
			if !mre.signaled {
				break
			}
			select {
			case <-mre.ch:
				{

				}
			case <-mre.reset: // if many Wait calls occurred during Reset, we break with !mre.signaled condition. But if we blocked on case <-mre.ch, then break will be with case <-mre.reset
				{
					break
				}
			}
		}
		mre.reset = make(chan struct{})
	}()
}

func (mre *ManualResetEvent) Reset() {
	mre.signaled = false
	close(mre.reset)
}

func (mre *ManualResetEvent) IsSet() bool {
	return mre.signaled
}
