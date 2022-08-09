package mychannels

import "testing"

const m = 16

func BenchmarkNativeBufferedChannel(b *testing.B) {
	done := make(chan struct{})
	ch := make(chan int, m)
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

func BenchmarkMyBufferedChannel(b *testing.B) {
	done := make(chan struct{})
	var ch Channel[int] = NewBufferedChannel[int](m)
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
