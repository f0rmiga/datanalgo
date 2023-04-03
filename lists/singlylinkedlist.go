// Copyright 2023 Thulio Ferraz Assis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lists

import (
	"errors"
	"fmt"
	"strings"
)

// Node represents a node in the singly linked list. It contains a value of type T and a pointer to
// the next node.
type Node[T comparable] struct {
	value T
	next  *Node[T]
}

// SinglyLinkedList represents a singly linked list with a head pointer and size.
type SinglyLinkedList[T comparable] struct {
	head *Node[T]
	size int
}

var _ LinkedList[struct{}] = (*SinglyLinkedList[struct{}])(nil)

// InsertFirst inserts a new node with the given value at the beginning of the list.
func (list *SinglyLinkedList[T]) InsertFirst(value T) {
	newNode := &Node[T]{value: value, next: list.head}
	list.head = newNode
	list.size++
}

// InsertLast inserts a new node with the given value at the end of the list.
func (list *SinglyLinkedList[T]) InsertLast(value T) {
	newNode := &Node[T]{value: value, next: nil}
	if list.head == nil {
		list.head = newNode
	} else {
		current := list.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
	list.size++
}

// InsertAt inserts a new node with the given value at the specified index in the list. Returns an
// error if the index is out of range.
func (list *SinglyLinkedList[T]) InsertAt(value T, index int) error {
	if index < 0 || index > list.size {
		return errors.New("index out of range")
	}
	if index == 0 {
		list.InsertFirst(value)
		return nil
	}
	newNode := &Node[T]{value: value, next: nil}
	current := list.head
	for i := 1; i < index; i++ {
		current = current.next
	}
	newNode.next = current.next
	current.next = newNode
	list.size++
	return nil
}

// DeleteFirst deletes the first node in the list and returns its value. Returns an error if the
// list is empty.
func (list *SinglyLinkedList[T]) DeleteFirst() (val T, err error) {
	if list.head == nil {
		return val, errors.New("list is empty")
	}
	value := list.head.value
	list.head = list.head.next
	list.size--
	return value, nil
}

// DeleteLast deletes the last node in the list and returns its value.
func (list *SinglyLinkedList[T]) DeleteLast() (val T, err error) {
	if list.head == nil {
		return val, errors.New("list is empty")
	}
	if list.head.next == nil {
		value := list.head.value
		list.head = nil
		list.size--
		return value, nil
	}
	current := list.head
	for current.next.next != nil {
		current = current.next
	}
	value := current.next.value
	current.next = nil
	list.size--
	return value, nil
}

// DeleteAt deletes the node at the specified index in the list and returns its value. Returns an
// error if the index is out of range.
func (list *SinglyLinkedList[T]) DeleteAt(index int) (val T, err error) {
	if index < 0 || index >= list.size {
		return val, errors.New("index out of range")
	}
	if index == 0 {
		return list.DeleteFirst()
	}
	current := list.head
	for i := 1; i < index; i++ {
		current = current.next
	}
	value := current.next.value
	current.next = current.next.next
	list.size--
	return value, nil
}

// DeleteValue deletes the first occurrence of the given value in the list. Returns true if the
// value was found and deleted, false if the value was not found. Returns an error if the list is
// empty.
func (list *SinglyLinkedList[T]) DeleteValue(value T) (bool, error) {
	if list.head == nil {
		return false, errors.New("list is empty")
	}
	if list.head.value == value {
		list.head = list.head.next
		list.size--
		return true, nil
	}
	current := list.head
	for current.next != nil && current.next.value != value {
		current = current.next
	}
	if current.next == nil {
		return false, nil
	}
	current.next = current.next.next
	list.size--
	return true, nil
}

// Search searches for the given value in the list and returns the index of the first occurrence.
// Returns -1 if the value is not found. Returns an error if the list is empty.
func (list *SinglyLinkedList[T]) Search(value T) (int, error) {
	if list.head == nil {
		return -1, errors.New("list is empty")
	}
	current := list.head
	index := 0
	for current != nil {
		if current.value == value {
			return index, nil
		}
		index++
		current = current.next
	}
	return -1, nil
}

// Traversal traverses the list from the head to the tail, calling the given function for each
// node's value. Returns an error if the function returns an error for any value.
func (list *SinglyLinkedList[T]) Traversal(fn func(T) error) error {
	current := list.head
	for current != nil {
		if err := fn(current.value); err != nil {
			return err
		}
		current = current.next
	}
	return nil
}

// ReverseTraversal is not applicable to singly linked lists and will result in a panic if called.
func (list *SinglyLinkedList[T]) ReverseTraversal(fn func(T) error) error {
	panic("ReverseTraversal is not applicable to Singly Linked Lists")
}

// Size returns the size of the list (number of nodes).
func (list *SinglyLinkedList[T]) Size() int {
	return list.size
}

// IsEmpty returns true if the list is empty, false otherwise.
func (list *SinglyLinkedList[T]) IsEmpty() bool {
	return list.size == 0
}

// String returns a string representation of the list, with each value followed by an arrow ("->")
// pointing to the next value. The last value points to "nil", indicating the end of the list.
func (list *SinglyLinkedList[T]) String() string {
	var sb strings.Builder
	current := list.head
	for current != nil {
		fmt.Fprintf(&sb, "%v -> ", current.value)
		current = current.next
	}
	sb.WriteString("nil")
	return sb.String()
}
