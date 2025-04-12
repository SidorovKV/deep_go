package main

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type integerable interface {
	int8 | int16 | int32 | int64 | int
}

type CircularQueue[T integerable] struct {
	size   int
	first  int
	last   int
	len    int
	values []T
	sync.Mutex
}

func NewCircularQueue[T integerable](size int) *CircularQueue[T] {
	return &CircularQueue[T]{
		size:   size,
		first:  0,
		last:   0,
		len:    0,
		values: make([]T, size),
		Mutex:  sync.Mutex{},
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	q.Lock()
	defer q.Unlock()

	if q.len == 0 {
		q.values[q.last] = value
		q.len++

		return true
	}

	if q.len != q.size {
		q.last = q.nextPos(q.last)
		q.values[q.last] = value
		q.len++

		return true
	}

	return false
}

func (q *CircularQueue[T]) Pop() bool {
	q.Lock()
	defer q.Unlock()

	if q.len == 0 {
		return false
	}

	q.first = q.nextPos(q.first)
	q.len--

	return true
}

func (q *CircularQueue[T]) Front() T {
	q.Lock()
	defer q.Unlock()

	if q.len == 0 {
		return -1
	}

	return q.values[q.first]
}

func (q *CircularQueue[T]) Back() T {
	q.Lock()
	defer q.Unlock()

	if q.len == 0 {
		return -1
	}

	return q.values[q.last]
}

func (q *CircularQueue[T]) Empty() bool {
	q.Lock()
	defer q.Unlock()

	if q.len == 0 {
		return true
	}

	return false
}

func (q *CircularQueue[T]) Full() bool {
	q.Lock()
	defer q.Unlock()

	if q.len == q.size {
		return true
	}

	return false
}

func (q *CircularQueue[T]) nextPos(pointer int) int {
	if pointer+1 >= q.size {
		return 0
	}

	return pointer + 1
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueueInt8(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int8](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int8(-1), queue.Front())
	assert.Equal(t, int8(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int8{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int8(1), queue.Front())
	assert.Equal(t, int8(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int8{4, 2, 3}, queue.values))

	assert.Equal(t, int8(2), queue.Front())
	assert.Equal(t, int8(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueueInt16(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int16](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int16(-1), queue.Front())
	assert.Equal(t, int16(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int16{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int16(1), queue.Front())
	assert.Equal(t, int16(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int16{4, 2, 3}, queue.values))

	assert.Equal(t, int16(2), queue.Front())
	assert.Equal(t, int16(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueueInt32(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int32](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int32(-1), queue.Front())
	assert.Equal(t, int32(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int32{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int32(1), queue.Front())
	assert.Equal(t, int32(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int32{4, 2, 3}, queue.values))

	assert.Equal(t, int32(2), queue.Front())
	assert.Equal(t, int32(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueueInt64(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int64](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int64(-1), queue.Front())
	assert.Equal(t, int64(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int64{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int64(1), queue.Front())
	assert.Equal(t, int64(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int64{4, 2, 3}, queue.values))

	assert.Equal(t, int64(2), queue.Front())
	assert.Equal(t, int64(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
