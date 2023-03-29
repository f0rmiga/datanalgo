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

package concurrent

import "sync"

// Transform applies the given transformer to each item in the input slice and
// returns the results in a new slice. The transformer is applied concurrently
// using the given number of workers. The order of the results is the same as
// the order of the input items.
func Transform[Input any, Output any](items []Input, transformer Transformer[Input, Output], workers int) ([]Output, error) {
	var wg sync.WaitGroup
	inputCh := make(chan indexedItem[Input])
	outputCh := make(chan indexedResult[Output], len(items))

	// Start the worker goroutines.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(inputCh, outputCh, transformer, &wg)
	}

	// Send the items to the workers along with their indices.
	go func() {
		for i, item := range items {
			inputCh <- indexedItem[Input]{Index: i, Item: item}
		}
		close(inputCh)
	}()

	// Wait for all workers to finish.
	go func() {
		wg.Wait()
		close(outputCh)
	}()

	// Collect the results and maintain the input order.
	results := make([]Output, len(items))
	for indexedResult := range outputCh {
		if indexedResult.Err != nil {
			return nil, indexedResult.Err
		}
		results[indexedResult.Index] = indexedResult.Item
	}

	return results, nil
}

// Transformer is a function that transforms an input into an output.
type Transformer[Input any, Output any] func(Input) (Output, error)

type indexedItem[Item any] struct {
	Index int
	Item  Item
}

type indexedResult[Item any] struct {
	Index int
	Item  Item
	Err   error
}

func worker[Input any, Output any](inputCh <-chan indexedItem[Input], outputCh chan<- indexedResult[Output], transformer Transformer[Input, Output], wg *sync.WaitGroup) {
	defer wg.Done()
	for indexedInput := range inputCh {
		output, err := transformer(indexedInput.Item)
		indexedResult := indexedResult[Output]{Index: indexedInput.Index, Item: output, Err: err}
		outputCh <- indexedResult
	}
}
