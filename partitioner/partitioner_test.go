package partitioner

import "testing"

func Test_registerHandlers(t *testing.T) {
	channels, _ := registerHandlers(10, 1, func(index int, item Item) {})
	if len(channels) != 10 {
		t.Error("incorrect result: expected 10, got", len(channels))
	}
}
