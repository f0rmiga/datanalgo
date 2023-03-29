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

// Transformer is an interface that provides methods to concurrently apply a series
// of transformations on a list of input items. It preserves the order of input items
// in the output and can handle both transformations with and without error handling.
// It can be used with any input and output types.
type Transformer[Input any, Output any] interface {
	// Transform applies the provided actions on the input items concurrently,
	// preserving the order of input items in the output. Each action is a function
	// that transforms an input item into an output item. This method doesn't handle
	// errors and assumes that the actions will not return an error.
	Transform(items []Input, actions ...TransformAction[Input, Output]) []Output

	// TransformWithError applies the provided actions on the input items concurrently,
	// preserving the order of input items in the output. Each action is a function
	// that transforms an input item into an output item and may return an error.
	// If an action returns an error, the processing is halted, and the error is returned.
	TransformWithError(items []Input, actions ...TransformActionWithError[Input, Output]) ([]Output, error)
}

type transformer[Input any, Output any] struct {
	workers int
}

// NewTransformer returns a new Transformer instance with the specified number of workers.
// The workers parameter determines the concurrency level of the Transformer.
func NewTransformer[Input any, Output any](workers int) Transformer[Input, Output] {
	return &transformer[Input, Output]{
		workers: workers,
	}
}

func (t *transformer[Input, Output]) Transform(items []Input, actions ...TransformAction[Input, Output]) []Output {
	channels := make([]chan indexedItem[any], len(actions)+1)
	channels[0] = make(chan indexedItem[any], len(items))

	// Send the items to the first channel along with their indices.
	go func() {
		defer close(channels[0])
		for i, item := range items {
			channels[0] <- indexedItem[any]{Index: i, Item: item}
		}
	}()

	// Create a pipeline of worker functions connected by channels.
	for i, action := range actions {
		outputCh := make(chan indexedItem[any], len(items))
		channels[i+1] = outputCh
		var wg sync.WaitGroup
		for j := 0; j < t.workers; j++ {
			wg.Add(1)
			go worker(channels[i], outputCh, action, &wg)
		}
		go func() {
			wg.Wait()
			close(outputCh)
		}()
	}

	// Collect the results and maintain the input order.
	results := make([]Output, len(items))
	for indexedItem := range channels[len(channels)-1] {
		results[indexedItem.Index] = indexedItem.Item.(Output)
	}

	return results
}

func (t *transformer[Input, Output]) TransformWithError(items []Input, actions ...TransformActionWithError[Input, Output]) ([]Output, error) {
	channels := make([]chan indexedItem[any], len(actions)+1)
	channels[0] = make(chan indexedItem[any], len(items))

	// Send the items to the first channel along with their indices.
	go func() {
		defer close(channels[0])
		for i, item := range items {
			channels[0] <- indexedItem[any]{Index: i, Item: item}
		}
	}()

	// Create a pipeline of worker functions connected by channels.
	for i, action := range actions {
		outputCh := make(chan indexedItem[any], len(items))
		channels[i+1] = outputCh
		var wg sync.WaitGroup
		for j := 0; j < t.workers; j++ {
			wg.Add(1)
			go workerWithError(channels[i], outputCh, action, &wg)
		}
		go func() {
			wg.Wait()
			close(outputCh)
		}()
	}

	// Collect the results and maintain the input order.
	results := make([]Output, len(items))
	for indexedItem := range channels[len(channels)-1] {
		if indexedItem.Err != nil {
			return nil, indexedItem.Err
		}
		results[indexedItem.Index] = indexedItem.Item.(Output)
	}

	return results, nil
}

// TransformAction is a function that takes an input item and transforms it into an output item.
// This function is used with the Transform method and assumes that the transformation will not
// return an error.
type TransformAction[Input any, Output any] func(Input) Output

// TransformActionWithError is a function that takes an input item and transforms it into an output item.
// It may return an error if the transformation fails. This function is used with the
// TransformWithError method to handle errors during the transformation process.
type TransformActionWithError[Input any, Output any] func(Input) (Output, error)

type indexedItem[Item any] struct {
	Index int
	Item  Item
	Err   error
}

func worker[Input any, Output any](inputCh <-chan indexedItem[any], outputCh chan<- indexedItem[any], action TransformAction[Input, Output], wg *sync.WaitGroup) {
	defer wg.Done()
	for indexedInput := range inputCh {
		output := action(indexedInput.Item.(Input))
		indexedItem := indexedItem[any]{Index: indexedInput.Index, Item: output}
		outputCh <- indexedItem
	}
}

func workerWithError[Input any, Output any](inputCh <-chan indexedItem[any], outputCh chan<- indexedItem[any], action TransformActionWithError[Input, Output], wg *sync.WaitGroup) {
	defer wg.Done()
	for indexedInput := range inputCh {
		output, err := action(indexedInput.Item.(Input))
		indexedItem := indexedItem[any]{Index: indexedInput.Index, Item: output, Err: err}
		outputCh <- indexedItem
	}
}
