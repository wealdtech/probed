// Copyright © 2022 Weald Technology Trading.
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

// Package probedb defines a probe data store.
package probedb

import "github.com/attestantio/go-eth2-client/spec/phase0"

// Order is the order in which results should be fetched (N.B. fetched, not returned).
type Order uint8

const (
	// OrderEarliest fetches earliest transactions first.
	OrderEarliest Order = iota
	// OrderLatest fetches latest transactions first.
	OrderLatest
)

// Selection is the selection criterion when multiple delays are present.
type Selection uint8

const (
	// SelectionMinimum fetches the minimum delay.
	SelectionMinimum Selection = iota
	// SelectionMaximum fetches the maximum delay.
	SelectionMaximum
	// SelectionMedian fetches the median delay.
	SelectionMedian
)

// DelayFilter defines a filter for fetching delays.
// Filter elements are ANDed together.
// Results are always returned in ascending slot/method/IP address/source order.
type DelayFilter struct {
	// IPAddr is the IP address from which to fetch delays.
	// If empty then there is no IP address filter.
	IPAddr string

	// Source is the beacon node source from which to fetch delays.
	// If empty then there is no source filter.
	Source string

	// Method is the collection method from which to fetch delays.
	// If empty then there is no method filter.
	Method string

	// From is the slot of the earliest delay to fetch.
	// If nil then there is no earliest slot.
	From *phase0.Slot

	// To is the slot of the latest delay to fetch.
	// If nil then there is no latest slot.
	To *phase0.Slot

	// Order is either OrderEarliest, in which case the earliest results
	// that match the filter are returned, or OrderLatest, in which case the
	// latest results that match the filter are returned.
	// The default is OrderEarliest.
	Order Order

	// Selection is the selection of the delay(s).
	// The default is SelectionMinimum.
	Selection Selection
}

// AggregateAttestationFilter defines a filter for fetching aggregate attestations.
// Filter elements are ANDed together.
// Results are always returned in ascending slot/method/IP address/source order.
type AggregateAttestationFilter struct {
	// IPAddr is the IP address from which to fetch results.
	// If empty then there is no IP address filter.
	IPAddr string

	// Sources are the beacon nodes from which to fetch results.
	// If empty then there is no source filter.
	Sources []string

	// Method are the collection methods from which to fetch results.
	// If empty then there is no method filter.
	Methods []string

	// From is the slot of the earliest result to fetch.
	// If nil then there is no earliest slot.
	From *phase0.Slot

	// To is the slot of the latest result to fetch.
	// If nil then there is no latest slot.
	To *phase0.Slot

	// Order is either OrderEarliest, in which case the earliest results
	// that match the filter are returned, or OrderLatest, in which case the
	// latest results that match the filter are returned.
	// The default is OrderEarliest.
	Order Order

	// Limit is the maximum number of results to return.
	// If 0 then there is no limit.
	Limit uint32
}

// AttestationSummaryFilter defines a filter for fetching attestation summaries.
// Filter elements are ANDed together.
// Results are always returned in ascending slot/method/IP address/source order.
type AttestationSummaryFilter struct {
	// IPAddr is the IP address from which to fetch data.
	// If empty then there is no IP address filter.
	IPAddr string

	// Source is the beacon node source from which to fetch data.
	// If empty then there is no source filter.
	Source string

	// Method is the collection method from which to fetch data.
	// If empty then there is no method filter.
	Method string

	// From is the slot of the earliest data to fetch.
	// If nil then there is no earliest slot.
	From *phase0.Slot

	// To is the slot of the latest data to fetch.
	// If nil then there is no latest slot.
	To *phase0.Slot

	// Order is either OrderEarliest, in which case the earliest results
	// that match the filter are returned, or OrderLatest, in which case the
	// latest results that match the filter are returned.
	// The default is OrderEarliest.
	Order Order
}
