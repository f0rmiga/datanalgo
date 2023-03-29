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
	"fmt"
	"strings"
	"testing"

	"github.com/f0rmiga/datanalgo/concurrent"
)

func TestTransform(t *testing.T) {
	errorFunc := errorFunc[string, string] // Workaround for type inference.
	testCases := []struct {
		name        string
		inputItems  interface{}
		transformer interface{}
		workers     int
		expected    interface{}
		expectErr   bool
	}{
		{
			name:        "To upper case",
			inputItems:  []string{"a", "b", "c"},
			transformer: concurrent.Transformer[string, string](upperCase),
			workers:     3,
			expected:    []string{"A", "B", "C"},
			expectErr:   false,
		},
		{
			name:        "To upper case with 1 worker",
			inputItems:  []string{"a", "b", "c"},
			transformer: concurrent.Transformer[string, string](upperCase),
			workers:     1,
			expected:    []string{"A", "B", "C"},
			expectErr:   false,
		},
		{
			name:        "To upper case with 10 workers",
			inputItems:  []string{"a", "b", "c"},
			transformer: concurrent.Transformer[string, string](upperCase),
			workers:     10,
			expected:    []string{"A", "B", "C"},
			expectErr:   false,
		},
		{
			name:        "Length of strings",
			inputItems:  []string{"a", "bb", "ccc"},
			transformer: concurrent.Transformer[string, int](length),
			workers:     3,
			expected:    []int{1, 2, 3},
			expectErr:   false,
		},
		{
			name:        "Error handling",
			inputItems:  []string{"a", "b", "c"},
			transformer: concurrent.Transformer[string, string](errorFunc),
			workers:     3,
			expected:    []string{"", "", ""},
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch items := tc.inputItems.(type) {
			case []string:
				switch exp := tc.expected.(type) {
				case []string:
					tr := tc.transformer.(concurrent.Transformer[string, string])
					result, err := concurrent.Transform(items, tr, tc.workers)
					if tc.expectErr {
						if err == nil {
							t.Errorf("Expected error, but got none")
						}
					} else {
						if err != nil {
							t.Errorf("Unexpected error: %v", err)
						}
						for i, item := range result {
							if item != exp[i] {
								t.Errorf("Expected item %v at index %d, got %v", exp[i], i, item)
							}
						}
					}
				case []int:
					tr := tc.transformer.(concurrent.Transformer[string, int])
					result, err := concurrent.Transform(items, tr, tc.workers)
					if tc.expectErr {
						if err == nil {
							t.Errorf("Expected error, but got none")
						}
					} else {
						if err != nil {
							t.Errorf("Unexpected error: %v", err)
						}
						for i, item := range result {
							if item != exp[i] {
								t.Errorf("Expected item %v at index %d, got %v", exp[i], i, item)
							}
						}
					}
				default:
					panic(fmt.Errorf("Unexpected type %T", exp))
				}
			}
		})
	}
}

// Test helper: strings.ToUpper
func upperCase(item string) (string, error) {
	return strings.ToUpper(item), nil
}

// Test helper: length function for strings
func length(item string) (int, error) {
	return len(item), nil
}

// Test helper: error function
func errorFunc[T any, R any](item T) (string, error) {
	return "", errors.New("error")
}
