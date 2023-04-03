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

// LinkedList represents a generic linked list that can store elements of any type.
type LinkedList[T any] interface {
	// InsertFirst adds an element to the beginning of the linked list.
	InsertFirst(value T)

	// InsertLast adds an element to the end of the linked list.
	InsertLast(value T)

	// InsertAt inserts an element at the specified index of the linked list.
	InsertAt(value T, index int) error

	// DeleteFirst removes and returns the first element from the linked list.
	DeleteFirst() (T, error)

	// DeleteLast removes and returns the last element from the linked list.
	DeleteLast() (T, error)

	// DeleteAt removes and returns the element at the specified index of the linked list.
	DeleteAt(index int) (T, error)

	// DeleteValue removes the first occurrence of the specified value from the linked list.
	DeleteValue(value T) (bool, error)

	// Search returns the index of the first occurrence of the specified value in the linked list.
	Search(value T) (int, error)

	// Traversal applies the given function to each element of the linked list, in order.
	Traversal(func(T) error) error

	// ReverseTraversal applies the given function to each element of the linked list, in reverse
	// order.
	ReverseTraversal(func(T) error) error

	// Size returns the number of elements in the linked list.
	Size() int

	// IsEmpty returns true if the linked list is empty, false otherwise.
	IsEmpty() bool

	// String returns a string representation of the linked list.
	String() string
}
