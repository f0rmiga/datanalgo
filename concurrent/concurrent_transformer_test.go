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

package concurrent_test

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/f0rmiga/datanalgo/concurrent"
)

func TestTransform(t *testing.T) {
	testCasesStringToString := []testCase[string, string]{
		{
			name:     "To upper case",
			input:    []string{"a", "b", "c"},
			actions:  []concurrent.TransformAction[string, string]{strings.ToUpper},
			workers:  3,
			expected: []string{"A", "B", "C"},
		},
		{
			name:     "To upper case with 1 worker",
			input:    []string{"a", "b", "c"},
			actions:  []concurrent.TransformAction[string, string]{strings.ToUpper},
			workers:  1,
			expected: []string{"A", "B", "C"},
		},
		{
			name:     "To upper case with 10 workers",
			input:    []string{"a", "b", "c"},
			actions:  []concurrent.TransformAction[string, string]{strings.ToUpper},
			workers:  10,
			expected: []string{"A", "B", "C"},
		},
	}

	for _, tc := range testCasesStringToString {
		t.Run(tc.name, runTest(tc))
	}

	testCasesStringToInt := []testCase[string, int]{
		{
			name:     "Length of strings",
			input:    []string{"a", "bb", "ccc"},
			actions:  []concurrent.TransformAction[string, int]{func(s string) int { return len(s) }},
			workers:  3,
			expected: []int{1, 2, 3},
		},
	}

	for _, tc := range testCasesStringToInt {
		t.Run(tc.name, runTest(tc))
	}

	testCasesMultipleTransformers := []testCase[string, string]{
		{
			name:  "Multiple transformers",
			input: []string{"a", "b", "c"},
			actions: []concurrent.TransformAction[string, string]{
				strings.ToUpper,
				strings.ToLower,
			},
			workers:  3,
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tc := range testCasesMultipleTransformers {
		t.Run(tc.name, runTest(tc))
	}
}

func TestTransformWithError(t *testing.T) {
	errorFunc := errorFunc[string] // Workaround for type inference.
	testCasesStringToString := []testCase[string, string]{
		{
			name:      "To upper case",
			input:     []string{"a", "b", "c"},
			actions:   []concurrent.TransformActionWithError[string, string]{upperCase},
			workers:   3,
			expected:  []string{"A", "B", "C"},
			expectErr: false,
		},
		{
			name:      "To upper case with 1 worker",
			input:     []string{"a", "b", "c"},
			actions:   []concurrent.TransformActionWithError[string, string]{upperCase},
			workers:   1,
			expected:  []string{"A", "B", "C"},
			expectErr: false,
		},
		{
			name:      "To upper case with 10 workers",
			input:     []string{"a", "b", "c"},
			actions:   []concurrent.TransformActionWithError[string, string]{upperCase},
			workers:   10,
			expected:  []string{"A", "B", "C"},
			expectErr: false,
		},
		{
			name:      "Error handling",
			input:     []string{"a", "b", "c"},
			actions:   []concurrent.TransformActionWithError[string, string]{errorFunc},
			workers:   3,
			expected:  []string{"", "", ""},
			expectErr: true,
		},
	}

	for _, tc := range testCasesStringToString {
		t.Run(tc.name, runTestWithError(tc))
	}

	testCasesStringToInt := []testCase[string, int]{
		{
			name:      "Length of strings",
			input:     []string{"a", "bb", "ccc"},
			actions:   []concurrent.TransformActionWithError[string, int]{length},
			workers:   3,
			expected:  []int{1, 2, 3},
			expectErr: false,
		},
	}

	for _, tc := range testCasesStringToInt {
		t.Run(tc.name, runTestWithError(tc))
	}

	testCasesMultipleTransformers := []testCase[string, string]{
		{
			name:  "Multiple transformers",
			input: []string{"a", "b", "c"},
			actions: []concurrent.TransformActionWithError[string, string]{
				upperCase,
				lowerCase,
			},
			workers:   3,
			expected:  []string{"a", "b", "c"},
			expectErr: false,
		},
	}

	for _, tc := range testCasesMultipleTransformers {
		t.Run(tc.name, runTestWithError(tc))
	}
}

func TestTransformChannels(t *testing.T) {
	workers := 4
	transformer := concurrent.NewTransformer[int, int](workers)

	inputItems := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	inputChan := make(chan int, len(inputItems))
	for _, item := range inputItems {
		inputChan <- item
	}
	close(inputChan)

	outputChan := transformer.TransformChannels(inputChan, func(item int) int {
		return item * 2
	})

	expectedOutput := make(map[int]bool)
	for _, item := range inputItems {
		expectedOutput[item*2] = false
	}

	for output := range outputChan {
		if _, ok := expectedOutput[output]; !ok {
			t.Errorf("Unexpected output: %d", output)
		} else {
			expectedOutput[output] = true
		}
	}

	for item, seen := range expectedOutput {
		if !seen {
			t.Errorf("Expected output not seen: %d", item)
		}
	}
}

func TestTransformChannelsWithError(t *testing.T) {
	workers := 4
	transformer := concurrent.NewTransformer[int, int](workers)

	inputItems := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	inputChan := make(chan int, len(inputItems))
	for _, item := range inputItems {
		inputChan <- item
	}
	close(inputChan)

	outputChan, errChan := transformer.TransformChannelsWithError(inputChan, func(item int) (int, error) {
		if item%5 == 0 {
			return 0, errors.New("item is divisible by 5")
		}
		return item * 2, nil
	})

	expectedOutput := make(map[int]bool)
	for _, item := range inputItems {
		if item%5 != 0 {
			expectedOutput[item*2] = false
		}
	}

	var outputError error
	var outputs []int
	for {
		select {
		case output, ok := <-outputChan:
			if !ok {
				outputChan = nil
			} else {
				outputs = append(outputs, output)
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				outputError = err
			}
		}

		if outputChan == nil && errChan == nil {
			break
		}
	}

	if outputError == nil {
		t.Errorf("Expected an error but did not receive one")
	}

	sort.Ints(outputs)

	for i, output := range outputs {
		if output != inputItems[i]*2 && inputItems[i]%5 != 0 {
			t.Errorf("Expected output %d, but got %d", inputItems[i]*2, output)
		}
	}
}

func runTest[Input any, Output comparable](tc testCase[Input, Output]) func(t *testing.T) {
	return func(t *testing.T) {
		transformer := concurrent.NewTransformer[Input, Output](tc.workers)
		actions := tc.actions.([]concurrent.TransformAction[Input, Output])
		result := transformer.Transform(tc.input, actions...)
		assertTestCase(t, tc, result)
	}
}

func runTestWithError[Input any, Output comparable](tc testCase[Input, Output]) func(t *testing.T) {
	return func(t *testing.T) {
		transformer := concurrent.NewTransformer[Input, Output](tc.workers)
		actions := tc.actions.([]concurrent.TransformActionWithError[Input, Output])
		result, err := transformer.TransformWithError(tc.input, actions...)
		if tc.expectErr {
			if err == nil {
				t.Errorf("Expected error, but got none")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			assertTestCase(t, tc, result)
		}
	}
}

func assertTestCase[Input any, Output comparable](t *testing.T, tc testCase[Input, Output], result []Output) {
	if len(result) != len(tc.expected) {
		t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
	}
	for i, item := range result {
		if item != tc.expected[i] {
			t.Errorf("Expected item %v at index %d, got %v", tc.expected[i], i, item)
		}
	}
}

type testCase[Input any, Output any] struct {
	name      string
	input     []Input
	actions   any
	workers   int
	expected  []Output
	expectErr bool
}

// Test helper: error function
func errorFunc[T any](item T) (T, error) {
	return item, errors.New("error")
}

// Test helper: strings.ToUpper
func upperCase(item string) (string, error) {
	return strings.ToUpper(item), nil
}

// Test helper: strings.ToLower
func lowerCase(item string) (string, error) {
	return strings.ToLower(item), nil
}

// Test helper: length function for strings
func length(item string) (int, error) {
	return len(item), nil
}
