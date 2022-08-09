package mychannels

type CloseableChannel interface {
	Close()
}

type ReadableChannel[T any] interface {
	CloseableChannel
	Receive() (T, bool)
}

type WriteableChannel[T any] interface {
	CloseableChannel
	Send(item T)
}

type Channel[T any] interface {
	ReadableChannel[T]
	WriteableChannel[T]
}
