package main

import (
	"fmt"

	"github.com/nsinitsyn/go-utils/partitioner"
)

type TestItem struct {
	value int
}

func (t TestItem) GetId() int {
	return t.value
}

func main() {
	// partitioner example
	const degreeOfConcurrency int = 5

	input := make(chan partitioner.Item)

	go func() {
		for i := 1; i <= 20; i++ {
			input <- TestItem{value: i}
		}
		close(input)
	}()

	partitioner.StartReading(degreeOfConcurrency, 100, input, func(channelIndex int, item partitioner.Item) {
		fmt.Printf("%d - %d\n", channelIndex, item.GetId())
	})
}
