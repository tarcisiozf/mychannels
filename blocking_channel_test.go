package mychannels

import (
	"testing"
)

const n = 4096

func BenchmarkNativeBlockingChannel(b *testing.B) {
	ch := make(chan int)
	done := make(chan struct{})
	var sum, expected int

	for i := 1; i <= n; i++ {
		expected += i
	}

	go func() {
		for x := range ch {
			sum += x
		}

		if sum != expected {
			b.Errorf("expected: %d, got: %d\n", expected, sum)
		}

		done <- struct{}{}
	}()

	for i := 1; i <= n; i++ {
		ch <- i
	}

	close(ch)
	<-done
}

func BenchmarkMyBlockingChannel(b *testing.B) {
	done := make(chan struct{})
	var ch Channel[int] = NewBlockingChannel[int]()
	var sum, expected int

	for i := 1; i <= n; i++ {
		expected += i
	}

	go func() {
		for {
			x, ok := ch.Receive()
			if ok {
				sum += x
			} else {
				break
			}
		}

		if sum != expected {
			b.Errorf("expected: %d, got: %d\n", expected, sum)
		}

		done <- struct{}{}
	}()

	for i := 1; i <= n; i++ {
		ch.Send(i)
	}

	ch.Close()
	<-done
}
