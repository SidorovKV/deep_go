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
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.head == nil {
		m.head = &node[K, V]{
			key:   key,
			value: value,
		}

		return
	}

	m.head.insert(key, value)
}

func (m *OrderedMap[K, V]) Erase(key K) {
	remove(m.head, key)
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	return m.head.contains(key)
}

func (m *OrderedMap[K, V]) Size() int {
	return m.head.countElements()
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	m.head.forEach(action)
}

type node[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *node[K, V]
	right *node[K, V]
}

func (n *node[K, V]) insert(key K, value V) {
	if n.key == key {
		n.value = value

		return
	}

	if key < n.key {
		if n.left == nil {
			n.left = &node[K, V]{
				key:   key,
				value: value,
			}

			return
		}

		n.left.insert(key, value)

		return
	}

	if n.right == nil {
		n.right = &node[K, V]{
			key:   key,
			value: value,
		}

		return
	}

	n.right.insert(key, value)
}

func (n *node[K, V]) countElements() int {
	if n == nil {
		return 0
	}

	return n.left.countElements() + n.right.countElements() + 1
}

func (n *node[K, V]) contains(key K) bool {
	if n == nil {
		return false
	}

	if n.key == key {
		return true
	}

	if key < n.key {
		return n.left.contains(key)
	}

	return n.right.contains(key)
}

func (n *node[K, V]) forEach(action func(K, V)) {
	if n == nil {
		return
	}

	n.left.forEach(action)
	action(n.key, n.value)
	n.right.forEach(action)
}

func remove[K cmp.Ordered, V any](node *node[K, V], key K) *node[K, V] {
	if node == nil {
		return nil
	}

	if key < node.key {
		node.left = remove(node.left, key)

		return node
	}

	if key > node.key {
		node.right = remove(node.right, key)

		return node
	}

	// key == node.key
	if node.left == nil && node.right == nil {
		node = nil

		return nil
	}

	if node.left == nil {
		node = node.right

		return node
	}

	if node.right == nil {
		node = node.left

		return node
	}
	leftmostrightside := node.right

	for {
		//find smallest value on the right side
		if leftmostrightside != nil && leftmostrightside.left != nil {
			leftmostrightside = leftmostrightside.left

			continue
		}

		break
	}

	node.key, node.value = leftmostrightside.key, leftmostrightside.value
	node.right = remove(node.right, node.key)

	return node
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
