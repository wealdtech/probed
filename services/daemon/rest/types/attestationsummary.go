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

package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// AttestationSummary holds summary information about attestations with a particular vote.
type AttestationSummary struct {
	Method       string
	Slot         uint32
	Attestations []*Attestation
}

// attestationSummaryJSON is a raw representation of the struct.
type attestationSummaryJSON struct {
	Method       string         `json:"method"`
	Slot         string         `json:"slot"`
	Attestations []*Attestation `json:"attestations"`
}

// MarshalJSON implements json.Marshaler.
func (a *AttestationSummary) MarshalJSON() ([]byte, error) {
	return json.Marshal(&attestationSummaryJSON{
		Method:       a.Method,
		Slot:         fmt.Sprintf("%d", a.Slot),
		Attestations: a.Attestations,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AttestationSummary) UnmarshalJSON(input []byte) error {
	var data attestationSummaryJSON
	err := json.Unmarshal(input, &data)
	if err != nil {
		return err
	}

	if data.Method == "" {
		return errors.New("method missing")
	}
	a.Method = data.Method

	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 32)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	a.Slot = uint32(slot)

	if data.Attestations == nil {
		return errors.New("attestations missing")
	}
	a.Attestations = data.Attestations

	return nil
}

// Attestation holds information about an attestation with a given vote.
type Attestation struct {
	CommitteeIndex  uint16
	BeaconBlockRoot []byte
	SourceRoot      []byte
	TargetRoot      []byte
	Buckets         map[string]*[120]bitfield.Bitlist
}

// attestationJSON is a raw representation of the struct.
type attestationJSON struct {
	CommitteeIndex  string              `json:"committee_index"`
	BeaconBlockRoot string              `json:"beacon_block_root"`
	SourceRoot      string              `json:"source_root"`
	TargetRoot      string              `json:"target_root"`
	Buckets         map[string][]string `json:"buckets"`
}

// MarshalJSON implements json.Marshaler.
func (a *Attestation) MarshalJSON() ([]byte, error) {
	bucketsStr := make(map[string][]string)
	for source, buckets := range a.Buckets {
		bucketsStr[source] = make([]string, 120)
		for i, bucket := range buckets {
			bucketsStr[source][i] = fmt.Sprintf("%#x", bucket)
		}
	}

	return json.Marshal(&attestationJSON{
		CommitteeIndex:  fmt.Sprintf("%d", a.CommitteeIndex),
		BeaconBlockRoot: fmt.Sprintf("%#x", a.BeaconBlockRoot),
		SourceRoot:      fmt.Sprintf("%#x", a.SourceRoot),
		TargetRoot:      fmt.Sprintf("%#x", a.TargetRoot),
		Buckets:         bucketsStr,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Attestation) UnmarshalJSON(input []byte) error {
	var data attestationJSON
	err := json.Unmarshal(input, &data)
	if err != nil {
		return err
	}

	if data.CommitteeIndex == "" {
		return errors.New("committee index missing")
	}
	committeeIndex, err := strconv.ParseUint(data.CommitteeIndex, 10, 16)
	if err != nil {
		return errors.Wrap(err, "invalid value for committee index")
	}
	a.CommitteeIndex = uint16(committeeIndex)

	if data.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	a.BeaconBlockRoot, err = hex.DecodeString(strings.TrimPrefix(data.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}

	if data.SourceRoot == "" {
		return errors.New("source root root missing")
	}
	a.SourceRoot, err = hex.DecodeString(strings.TrimPrefix(data.SourceRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for source root root")
	}

	if data.TargetRoot == "" {
		return errors.New("target root root missing")
	}
	a.TargetRoot, err = hex.DecodeString(strings.TrimPrefix(data.TargetRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for target root root")
	}

	if len(data.Buckets) == 0 {
		return errors.New("buckets missing")
	}

	a.Buckets = make(map[string]*[120]bitfield.Bitlist)
	for source, buckets := range data.Buckets {
		a.Buckets[source] = &[120]bitfield.Bitlist{}
		for i, bucket := range buckets {
			if bucket != "" {
				a.Buckets[source][i], err = hex.DecodeString(strings.TrimPrefix(bucket, "0x"))
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("invalid value for bucket %s index %d", source, i))
				}
			}
		}
	}

	return nil
}
