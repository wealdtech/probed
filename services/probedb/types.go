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
	"net"
)

// Delay holds information about a delay.
type Delay struct {
	IPAddr  net.IP
	Source  string
	Method  string
	Slot    uint32
	DelayMS uint32
}

// AttestationSummary holds summary information about an attestation.
type AttestationSummary struct {
	IPAddr          net.IP
	Source          string
	Method          string
	Slot            uint32
	CommitteeIndex  uint16
	BeaconBlockRoot []byte
	SourceRoot      []byte
	TargetRoot      []byte
	// AttesterBuckets contains the information about when specific indices
	// were first seen.
	// This is a raw representation of a github.com/prysmaticlabs/go-bitfield.Bitlist
	AttesterBuckets [][]byte
}

// AggregateAttestation holds information about an aggregate attestation.
type AggregateAttestation struct {
	IPAddr          net.IP
	Source          string
	Method          string
	Slot            uint32
	CommitteeIndex  uint16
	AggregationBits []byte
	BeaconBlockRoot []byte
	SourceRoot      []byte
	TargetRoot      []byte
	DelayMS         uint32
}
