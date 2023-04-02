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

// Package concurrent provides functionality to concurrently apply a series of transformations
// on a list of input items while preserving the order of input items in the output. It supports
// transformations with and without error handling and can be used with any input and output types.
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

	// TransformChannels takes a channel of input items and applies the provided actions
	// concurrently, not guaranteeing the order of input items in the output channel. Each action
	// is a function that transforms an input item into an output item. This method doesn't
	// handle errors and assumes that the actions will not return an error.
	TransformChannels(items <-chan Input, actions ...TransformAction[Input, Output]) <-chan Output

	// TransformChannelsWithError takes a channel of input items and applies the provided
	// actions concurrently, not guaranteeing the order of input items in the output channel. Each
	// action is a function that transforms an input item into an output item and may return
	// an error. If an action returns an error, the processing is halted, and the error is
	// sent to the error channel.
	TransformChannelsWithError(items <-chan Input, actions ...TransformActionWithError[Input, Output]) (<-chan Output, <-chan error)
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
	// Send the items to the first channel along with their indices.
	itemsCh := make(chan IndexedItem[any], len(items))
	go func() {
		defer close(itemsCh)
		for i, item := range items {
			itemsCh <- IndexedItem[any]{Index: i, Item: item}
		}
	}()

	transformedItemsCh := process[Input, Output](itemsCh, actions, t.workers, func(
		inputChan <-chan IndexedItem[any],
		outputChan chan<- IndexedItem[any],
		action TransformAction[Input, Output],
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		for indexedInput := range inputChan {
			output := action(indexedInput.Item.(Input))
			outputChan <- IndexedItem[any]{Index: indexedInput.Index, Item: output}
		}
	})

	// Collect the results and maintain the input order.
	transformedItems := make([]Output, len(items))
	for indexedItem := range transformedItemsCh {
		transformedItems[indexedItem.Index] = indexedItem.Item.(Output)
	}

	return transformedItems
}

func (t *transformer[Input, Output]) TransformWithError(items []Input, actions ...TransformActionWithError[Input, Output]) ([]Output, error) {
	// Send the items to the first channel along with their indices.
	itemsCh := make(chan IndexedItemWithError[any], len(items))
	go func() {
		defer close(itemsCh)
		for i, item := range items {
			itemsCh <- IndexedItemWithError[any]{Index: i, Item: item}
		}
	}()

	transformedItemsCh := process[Input, Output](itemsCh, actions, t.workers, func(
		inputChan <-chan IndexedItemWithError[any],
		outputChan chan<- IndexedItemWithError[any],
		action TransformActionWithError[Input, Output],
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		for indexedInput := range inputChan {
			output, err := action(indexedInput.Item.(Input))
			outputChan <- IndexedItemWithError[any]{Index: indexedInput.Index, Item: output, Err: err}
		}
	})

	// Collect the results and maintain the input order.
	transformedItems := make([]Output, len(items))
	for indexedItem := range transformedItemsCh {
		if indexedItem.Err != nil {
			return nil, indexedItem.Err
		}
		transformedItems[indexedItem.Index] = indexedItem.Item.(Output)
	}

	return transformedItems, nil
}

func (t *transformer[Input, Output]) TransformChannels(items <-chan Input, actions ...TransformAction[Input, Output]) <-chan Output {
	itemsCh := make(chan any, len(items))
	go func() {
		defer close(itemsCh)
		for item := range items {
			itemsCh <- item
		}
	}()

	transformedItemsCh := process[Input, Output](itemsCh, actions, t.workers, func(
		inputChan <-chan any,
		outputChan chan<- any,
		action TransformAction[Input, Output],
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		for indexedInput := range inputChan {
			output := action(indexedInput.(Input))
			outputChan <- output
		}
	})

	transformedItems := make(chan Output, len(items))
	go func() {
		defer close(transformedItems)
		for item := range transformedItemsCh {
			transformedItems <- item.(Output)
		}
	}()

	return transformedItems
}

func (t *transformer[Input, Output]) TransformChannelsWithError(items <-chan Input, actions ...TransformActionWithError[Input, Output]) (<-chan Output, <-chan error) {
	itemsCh := make(chan ItemWithError[any], len(items))
	go func() {
		defer close(itemsCh)
		for item := range items {
			itemsCh <- ItemWithError[any]{Item: item}
		}
	}()

	transformedItemsCh := process[Input, Output](itemsCh, actions, t.workers, func(
		inputChan <-chan ItemWithError[any],
		outputChan chan<- ItemWithError[any],
		action TransformActionWithError[Input, Output],
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		for indexedInput := range inputChan {
			output, err := action(indexedInput.Item.(Input))
			outputChan <- ItemWithError[any]{Item: output, Err: err}
		}
	})

	transformedItems := make(chan Output, len(items))
	errors := make(chan error, 1)
	go func() {
		defer close(transformedItems)
		defer close(errors)
		for item := range transformedItemsCh {
			if item.Err != nil {
				errors <- item.Err
				return
			}
			transformedItems <- item.Item.(Output)
		}
	}()

	return transformedItems, errors
}

func process[
	Input any,
	Output any,
	Item itemType,
	Action actionType[Input, Output],
	Worker func(inputChan <-chan Item, outputChan chan<- Item, action Action, wg *sync.WaitGroup),
](
	items <-chan Item,
	actions []Action,
	workers int,
	worker Worker,
) <-chan Item {
	channels := make([]chan Item, len(actions))

	// Create a pipeline of worker functions connected by channels.
	for i, action := range actions {
		var inputChan <-chan Item
		if i == 0 {
			inputChan = items
		} else {
			inputChan = channels[i-1]
		}
		outputChan := make(chan Item, len(items))
		channels[i] = outputChan
		var wg sync.WaitGroup
		for j := 0; j < workers; j++ {
			wg.Add(1)
			go worker(inputChan, outputChan, action, &wg)
		}
		go func() {
			wg.Wait()
			close(outputChan)
		}()
	}

	return channels[len(channels)-1]
}

// TransformAction is a function that takes an input item and transforms it into an output item.
// This function is used with the Transform method and assumes that the transformation will not
// return an error.
type TransformAction[Input any, Output any] func(Input) Output

// TransformActionWithError is a function that takes an input item and transforms it into an output item.
// It may return an error if the transformation fails. This function is used with the
// TransformWithError method to handle errors during the transformation process.
type TransformActionWithError[Input any, Output any] func(Input) (Output, error)

// ItemWithError is a struct that holds an item and an associated error. It is used to
// represent the result of a transformation that may return an error.
type ItemWithError[Item any] struct {
	Item Item
	Err  error
}

// IndexedItem is a struct that holds an item and its index. It is used to
// maintain the order of input items during concurrent transformations.
type IndexedItem[Item any] struct {
	Index int
	Item  Item
}

// IndexedItemWithError is a struct that holds an item, its index, and an associated error.
// It is used to maintain the order of input items during concurrent transformations that
// may return errors.
type IndexedItemWithError[Item any] struct {
	Index int
	Item  Item
	Err   error
}

type itemType interface {
	ItemWithError[any] | IndexedItem[any] | IndexedItemWithError[any] | any
}

type actionType[Input any, Output any] interface {
	TransformAction[Input, Output] | TransformActionWithError[Input, Output]
}
