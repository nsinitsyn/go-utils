package partitioner

import (
	"fmt"
	"sync/atomic"
)

// Read input channel and distribute items by partitions by ID. Function guarantee that items with the same identifier handle in one partition (channel). Read input channel until close.
func StartReading(degreeOfConcurrency int, channelBufferLength int, in <-chan Item, handler func(int, Item)) {
	channels, completed := registerHandlers(degreeOfConcurrency, 100, func(channelIndex int, item Item) {
		fmt.Printf("%d - %d\n", channelIndex, item.GetId())
	})

	// todo: if some partitioner-channel reach the maximum capacity, it will be slower input channel reading!
	for item := range in {
		channelIndex := item.GetId() % degreeOfConcurrency
		channels[channelIndex] <- item
	}

	for _, channel := range channels {
		close(channel)
	}

	<-completed
}

func registerHandlers(degreeOfConcurrency int, channelBufferLength int, handler func(int, Item)) ([]chan Item, chan struct{}) {
	channels := make([]chan Item, 0, degreeOfConcurrency)

	completed := make(chan struct{})

	var completedChannelsCount atomic.Uint32

	for i := 0; i < degreeOfConcurrency; i++ {
		i := i
		ch := make(chan Item, channelBufferLength)
		channels = append(channels, ch)

		go func() {
			for item := range ch {
				handler(i, item)
			}

			if completedChannelsCount.Add(1) == uint32(degreeOfConcurrency) {
				close(completed)
			}
		}()
	}

	return channels, completed
}

type Item interface {
	GetId() int
}
