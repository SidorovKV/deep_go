package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap[K cmp.Ordered, V any] struct {
	head *node[K, V]
	len  int
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	m.len++

	if m.head == nil {
		m.head = &node[K, V]{
			key:   key,
			value: value,
		}

		return
	}

	next := m.head

	for {
		if key == next.key {
			next.value = value

			return
		}

		child := &next.left
		if key > next.key {
			child = &next.right
		}

		if *child == nil {
			*child = &node[K, V]{key: key, value: value}

			return
		}

		next = *child
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	m.len--

	var parent *node[K, V]
	current := m.head
	isLeftChild := false

	for current != nil && current.key != key {
		parent = current
		if key < current.key {
			current = current.left
			isLeftChild = true
		} else {
			current = current.right
			isLeftChild = false
		}
	}

	if current == nil {
		return
	}

	if current.left == nil || current.right == nil {
		var child *node[K, V]
		if current.left != nil {
			child = current.left
		} else {
			child = current.right
		}

		if parent == nil {
			m.head = nil

			return
		}

		if isLeftChild {
			parent.left = child
		} else {
			parent.right = child
		}

		return
	}

	successorParent := current
	successor := current.right

	for successor.left != nil {
		successorParent = successor
		successor = successor.left
	}

	current.key = successor.key
	current.value = successor.value

	if successorParent == current {
		successorParent.right = successor.right
	} else {
		successorParent.left = successor.right
	}

}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	current := m.head

	for current != nil {
		if key == current.key {
			return true
		}

		if key < current.key {
			current = current.left

			continue
		}

		current = current.right
	}

	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.len
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	stack := make([]*node[K, V], 0, m.len)
	current := m.head

	for current != nil || len(stack) > 0 {
		for current != nil {
			stack = append(stack, current)
			current = current.left
		}

		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		action(current.key, current.value)

		current = current.right
	}
}

type node[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *node[K, V]
	right *node[K, V]
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
