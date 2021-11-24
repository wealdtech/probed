// Copyright © 2021 Weald Technology Trading.
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

package probedb

import (
	"context"
)

// BlockDelaysSetter defines functions to create and update block delays.
type BlockDelaysSetter interface {
	// SetBlockDelay sets a block delay.
	SetBlockDelay(ctx context.Context, delay *Delay) error
}

// BlockDelaysProvider defines functions to obtain block delays.
type BlockDelaysProvider interface {

	// MedianBlockDelays obtains the median block delays for a range of slots.
	MedianBlockDelays(ctx context.Context,
		locationID uint16,
		sourceID uint16,
		method string,
		fromSlot uint32,
		toSlot uint32,
	) (
		[]*DelayValue,
		error,
	)

	// MinimumBlockDelays obtains the minimum block delays for a range of slots.
	MinimumBlockDelays(ctx context.Context,
		locationID uint16,
		sourceID uint16,
		method string,
		fromSlot uint32,
		toSlot uint32,
	) (
		[]*DelayValue,
		error,
	)
}

// HeadDelaysSetter defines functions to create and update head delays.
type HeadDelaysSetter interface {
	// SetHeadDelay sets a head delay.
	SetHeadDelay(ctx context.Context, delay *Delay) error
}

// HeadDelaysProvider defines functions to obtain head delays.
type HeadDelaysProvider interface {

	// MedianHeadDelays obtains the median head delays for a range of slots.
	MedianHeadDelays(ctx context.Context,
		locationID uint16,
		sourceID uint16,
		method string,
		fromSlot uint32,
		toSlot uint32,
	) (
		[]*DelayValue,
		error,
	)

	// MinimumHeadDelays obtains the minimum head delays for a range of slots.
	MinimumHeadDelays(ctx context.Context,
		locationID uint16,
		sourceID uint16,
		method string,
		fromSlot uint32,
		toSlot uint32,
	) (
		[]*DelayValue,
		error,
	)
}

// Service defines a minimal probe database service.
type Service interface {
	// BeginTx begins a transaction.
	BeginTx(ctx context.Context) (context.Context, context.CancelFunc, error)

	// CommitTx commits a transaction.
	CommitTx(ctx context.Context) error

	// SetMetadata sets a metadata key to a JSON value.
	SetMetadata(ctx context.Context, key string, value []byte) error

	// Metadata obtains the JSON value from a metadata key.
	Metadata(ctx context.Context, key string) ([]byte, error)
}
