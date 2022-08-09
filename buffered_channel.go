package mychannels

import (
	"sync"
)

type BufferedChannel[T any] struct {
	mu      *sync.Mutex
	isOpen  bool
	bufSize int
	head    int
	tail    int
	ring    []T
}

func NewBufferedChannel[T any](bufSize int) *BufferedChannel[T] {
	return &BufferedChannel[T]{
		mu:      &sync.Mutex{},
		isOpen:  true,
		bufSize: bufSize,
		ring:    make([]T, bufSize),
		head:    -1,
		tail:    0,
	}
}

func (ch *BufferedChannel[T]) Send(item T) {
	var pos int

	for {
		ch.mu.Lock()

		if !ch.isOpen {
			panic(errWriteWhenClosed)
		}

		// check if it's full
		if (ch.head - ch.tail + 1) == ch.bufSize {
			ch.mu.Unlock()
			continue
		}

		pos = ch.head + 1
		ch.ring[pos%ch.bufSize] = item
		ch.head = pos

		ch.mu.Unlock()
		return
	}
}

func (ch *BufferedChannel[T]) Receive() (T, bool) {
	var isEmpty bool
	var item T

	for {
		ch.mu.Lock()

		isEmpty = ch.head < ch.tail

		// if closed, wait until is empty
		if !ch.isOpen && isEmpty {
			// item will be empty
			return item, false
		}

		if isEmpty {
			ch.mu.Unlock()
			continue
		}

		item = ch.ring[ch.tail%ch.bufSize]
		ch.tail++

		ch.mu.Unlock()
		return item, true
	}
}

func (ch *BufferedChannel[T]) Close() {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if !ch.isOpen {
		panic(errAlreadyClosed)
	}
	ch.isOpen = false
}
