package mychannels

import (
	"errors"
	"sync"
)

var (
	errWriteWhenClosed = errors.New("cannot write to closed channel")
	errAlreadyClosed   = errors.New("channel is already closed")
)

type BlockingChannel[T any] struct {
	mu     *sync.Mutex
	isOpen bool
	isFull bool
	item   T
}

func NewBlockingChannel[T any]() *BlockingChannel[T] {
	return &BlockingChannel[T]{
		mu:     &sync.Mutex{},
		isOpen: true,
	}
}

func (ch *BlockingChannel[T]) Send(item T) {
	for {
		ch.mu.Lock()

		if !ch.isOpen {
			panic(errWriteWhenClosed)
		}

		if ch.isFull {
			ch.mu.Unlock()
			continue
		}

		ch.item = item
		ch.isFull = true
		ch.mu.Unlock()
		return
	}
}

func (ch *BlockingChannel[T]) Receive() (T, bool) {
	for {
		ch.mu.Lock()

		if !ch.isOpen && !ch.isFull {
			var none T
			return none, false
		}

		if !ch.isFull {
			ch.mu.Unlock()
			continue
		}

		item := ch.item
		ch.isFull = false
		ch.mu.Unlock()
		return item, true
	}
}

func (ch *BlockingChannel[T]) Close() {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if !ch.isOpen {
		panic(errAlreadyClosed)
	}
	ch.isOpen = false
}
