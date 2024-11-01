package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/nsinitsyn/go-utils/autoresetevent"
	"github.com/nsinitsyn/go-utils/manualresetevent"
	"github.com/nsinitsyn/go-utils/partitioner"
	"github.com/nsinitsyn/go-utils/semaphore"
)

type TestItem struct {
	value int
}

func (t TestItem) GetId() int {
	return t.value
}

func partitionerExample() {
	fmt.Println("Partitioner example")
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

func semaphoreExample() {
	fmt.Println("Semaphore example")
	sem := semaphore.New(3)

	wg := sync.WaitGroup{}

	for i := 0; i < 7; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem.Wait()
			defer sem.Release()

			fmt.Println(i)
			time.Sleep(time.Second * 2)
		}()
	}

	wg.Wait()
}

func autoResetEventExample() {
	fmt.Println("AutoResetEvent example")
	wg := sync.WaitGroup{}
	are := autoresetevent.New()
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			are.WaitOne()
			fmt.Println(i)
		}()
	}

	for i := 0; i < 5; i++ {
		go func() {
			are.Set()
		}()
		time.Sleep(time.Second * 1)
	}
	wg.Wait()
}

func manualResetEventExample1() {
	fmt.Println("ManualResetEvent example 1")
	wg := sync.WaitGroup{}
	mre := manualresetevent.New()
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			mre.Wait()
			fmt.Println(i)
		}()
	}

	time.Sleep(time.Second * 3)
	mre.Set()

	// Test Reset during Wait calls
	go func() {
		time.Sleep(time.Second * 6)
		mre.Reset()
	}()

	// Only some goroutines will be finished until mre.Reset() called
	for i := 0; i < 10; i++ {
		i := i
		time.Sleep(time.Second * 1)

		wg.Add(1)
		go func() {
			defer wg.Done()

			mre.Wait()
			fmt.Println(i)
		}()
	}

	// Finish remaining goroutines
	go func() {
		time.Sleep(time.Second * 3)
		mre.Set()
	}()

	wg.Wait()
}

func manualResetEventExample2() {
	fmt.Println("ManualResetEvent example 2")
	wg := sync.WaitGroup{}
	mre := manualresetevent.New()
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			mre.Wait()
			fmt.Println(i)
		}()
	}

	time.Sleep(time.Second * 3)
	mre.Set()

	// Test Reset without Wait calls
	time.Sleep(time.Second * 4)
	mre.Reset()

	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			mre.Wait()
			fmt.Println(i)
		}()
	}

	time.Sleep(time.Second * 3)
	mre.Set()

	wg.Wait()
}

func main() {
	partitionerExample()
	semaphoreExample()
	autoResetEventExample()
	manualResetEventExample1()
	manualResetEventExample2()
}
