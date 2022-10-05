// Copyright © 2021, 2022 Weald Technology Trading.
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

// AggregateAttestationsSetter defines functions to create and update aggregate attestations.
type AggregateAttestationsSetter interface {
	Service

	// SetAggregateAttestation sets an aggregate attestation.
	SetAggregateAttestation(ctx context.Context, aggregateAttestation *AggregateAttestation) error
}

// AggregateAttestationsProvider defines functions to obtain aggregate attestations.
type AggregateAttestationsProvider interface {
	// AggregateAttestations obtains the aggregate attestations for a filter.
	AggregateAttestations(ctx context.Context, filter *AggregateAttestationFilter) ([]*AggregateAttestation, error)
}

// AttestationSummariesSetter defines functions to create and update attestation summaries.
type AttestationSummariesSetter interface {
	Service

	// SetAttestationSummary sets an attestation summary.
	SetAttestationSummary(ctx context.Context, summary *AttestationSummary) error
}

// AttestationSummariesProvider defines functions to obtain attestation summaries.
type AttestationSummariesProvider interface {
	// AttestationSummaries obtains the attestation summaries for a filter.
	AttestationSummaries(ctx context.Context, filter *AttestationSummaryFilter) ([]*AttestationSummary, error)
}

// BlockDelaysSetter defines functions to create and update block delays.
type BlockDelaysSetter interface {
	Service

	// SetBlockDelay sets a block delay.
	SetBlockDelay(ctx context.Context, delay *Delay) error
}

// BlockDelaysProvider defines functions to obtain block delays.
type BlockDelaysProvider interface {
	// BlockDelays obtains the block delays for a range of slots.
	BlockDelays(ctx context.Context, filter *DelayFilter) ([]*DelayValue, error)
}

// HeadDelaysSetter defines functions to create and update head delays.
type HeadDelaysSetter interface {
	Service

	// SetHeadDelay sets a head delay.
	SetHeadDelay(ctx context.Context, delay *Delay) error
}

// HeadDelaysProvider defines functions to obtain head delays.
type HeadDelaysProvider interface {
	// HeadDelays obtains the minimum head delays for a range of slots.
	HeadDelays(ctx context.Context, filter *DelayFilter) ([]*DelayValue, error)
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
