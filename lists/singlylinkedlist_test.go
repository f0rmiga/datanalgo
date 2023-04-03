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

package lists_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/f0rmiga/datanalgo/lists"
)

func TestInsertFirst(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	if list.String() != "3 -> 2 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestInsertLast(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertLast(1)
	list.InsertLast(2)
	list.InsertLast(3)

	if list.String() != "1 -> 2 -> 3 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestInsertAt(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(3)
	list.InsertAt(2, 1)

	if list.String() != "3 -> 2 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestInsertAtMiddle(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)
	err := list.InsertAt(4, 2)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if list.String() != "3 -> 2 -> 4 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestInsertAtEnd(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)
	err := list.InsertAt(4, 3)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if list.String() != "3 -> 2 -> 1 -> 4 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestInsertAtZeroIndex(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	err := list.InsertAt(4, 0)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if list.String() != "4 -> 3 -> 2 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestInsertAtOutOfRange(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	err := list.InsertAt(4, 5)

	if err == nil {
		t.Error("Expected error for index out of range")
	}
}

func TestDeleteFirst(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)
	list.DeleteFirst()

	if list.String() != "2 -> 3 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestDeleteFirstEmptyList(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	_, err := list.DeleteFirst()

	if err == nil {
		t.Error("Expected error for empty list")
	}
}

func TestDeleteLast(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)
	list.DeleteLast()

	if list.String() != "1 -> 2 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestDeleteLastEmptyList(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	_, err := list.DeleteLast()

	if err == nil {
		t.Error("Expected error for empty list")
	}
}

func TestDeleteLastSingleElementList(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	value, err := list.DeleteLast()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != 1 {
		t.Errorf("Unexpected value: %d", value)
	}

	if !list.IsEmpty() {
		t.Error("Expected empty list after deleting last element")
	}
}

func TestDeleteAt(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)
	list.DeleteAt(1)

	if list.String() != "1 -> 3 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestDeleteAtOutOfRange(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	_, err := list.DeleteAt(5)

	if err == nil {
		t.Error("Expected error for index out of range")
	}
}

func TestDeleteAtZeroIndex(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	value, err := list.DeleteAt(0)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != 3 {
		t.Errorf("Unexpected value: %d", value)
	}

	if list.String() != "2 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestDeleteAtMiddle(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	value, err := list.DeleteAt(1)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != 2 {
		t.Errorf("Unexpected value: %d", value)
	}

	if list.String() != "3 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestDeleteAtEnd(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	value, err := list.DeleteAt(2)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if value != 1 {
		t.Errorf("Unexpected value: %d", value)
	}

	if list.String() != "3 -> 2 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestDeleteValue(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)
	list.DeleteValue(2)

	if list.String() != "1 -> 3 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestSearch(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)

	index, _ := list.Search(2)
	if index != 1 {
		t.Errorf("Unexpected index: %d", index)
	}
}

func TestSearchEmptyList(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}

	_, err := list.Search(1)

	if err == nil {
		t.Error("Expected error for empty list")
	}
}

func TestSearchNotFound(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	index, err := list.Search(4)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if index != -1 {
		t.Error("Expected index -1 for not found")
	}
}

func TestTraversal(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)

	expected := "1 -> 2 -> 3 -> nil"
	actual := ""
	err := list.Traversal(func(value int) error {
		actual += fmt.Sprintf("%d -> ", value)
		return nil
	})
	actual += "nil"

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Errorf("Unexpected traversal: %s", actual)
	}
}

func TestTraversalWithError(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	err := list.Traversal(func(value int) error {
		if value == 2 {
			return errors.New("failed")
		}
		return nil
	})

	if err == nil {
		t.Error("Expected error")
	}
}

func TestSinglyLinkedListReverseTraversal(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for ReverseTraversal method")
		}
	}()

	list := &lists.SinglyLinkedList[int]{}
	list.ReverseTraversal(func(value int) error {
		return nil
	})
}

func TestSize(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)

	if list.Size() != 3 {
		t.Errorf("Unexpected size: %d", list.Size())
	}
}

func TestDeleteValueEmptyList(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}

	found, err := list.DeleteValue(1)

	if err == nil {
		t.Error("Expected error for empty list")
	}

	if found {
		t.Error("Expected value not found")
	}
}

func TestDeleteValueNotFound(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	found, err := list.DeleteValue(4)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if found {
		t.Error("Expected value not found")
	}

	if list.String() != "3 -> 2 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestDeleteValueAtHead(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(1)
	list.InsertFirst(2)
	list.InsertFirst(3)

	found, err := list.DeleteValue(3)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !found {
		t.Error("Expected value found")
	}

	if list.String() != "2 -> 1 -> nil" {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

func TestIsEmpty(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}

	if !list.IsEmpty() {
		t.Errorf("Unexpected state: not empty")
	}

	list.InsertFirst(1)

	if list.IsEmpty() {
		t.Errorf("Unexpected state: empty")
	}
}

func TestString(t *testing.T) {
	list := &lists.SinglyLinkedList[int]{}
	list.InsertFirst(3)
	list.InsertFirst(2)
	list.InsertFirst(1)

	expected := "1 -> 2 -> 3 -> nil"
	actual := list.String()

	if actual != expected {
		t.Errorf("Unexpected string representation: %s", actual)
	}
}

func TestSinglyLinkedListWithCustomStruct(t *testing.T) {
	list := &lists.SinglyLinkedList[Person]{}

	list.InsertFirst(Person{Name: "John", Age: 30})
	list.InsertFirst(Person{Name: "Jane", Age: 25})
	list.InsertLast(Person{Name: "Bob", Age: 40})

	if list.Size() != 3 {
		t.Errorf("Unexpected list size: %d", list.Size())
	}

	expected := "{Jane 25} -> {John 30} -> {Bob 40} -> nil"
	if list.String() != expected {
		t.Errorf("Unexpected list state: %s", list.String())
	}

	index, err := list.Search(Person{Name: "John", Age: 30})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if index != 1 {
		t.Errorf("Unexpected index: %d", index)
	}

	deleted, err := list.DeleteValue(Person{Name: "Jane", Age: 25})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !deleted {
		t.Errorf("Expected Person value to be deleted")
	}
	expected = "{John 30} -> {Bob 40} -> nil"
	if list.String() != expected {
		t.Errorf("Unexpected list state: %s", list.String())
	}
}

type Person struct {
	Name string
	Age  int
}
