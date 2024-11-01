package manualresetevent

import (
	"sync/atomic"
)

type ManualResetEvent struct {
	ch       chan struct{}
	signaled int32
	reset    chan struct{}
}

func New() *ManualResetEvent {
	return &ManualResetEvent{make(chan struct{}), 0, make(chan struct{})}
}

func (mre *ManualResetEvent) Wait() {
	if mre.signaled == 1 {
		return
	}
	mre.ch <- struct{}{}
}

func (mre *ManualResetEvent) Set() {

	if !atomic.CompareAndSwapInt32(&mre.signaled, 0, 1) {
		return
	}

	go func() { // goroutine for unblocking Set
		for {
			if mre.signaled == 0 {
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
	if atomic.CompareAndSwapInt32(&mre.signaled, 1, 0) {
		close(mre.reset)
	}
}

func (mre *ManualResetEvent) IsSet() bool {
	return mre.signaled == 1
}
