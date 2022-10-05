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

// Package types provides types used in the REST API.
package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// AggregateAttestation holds information about a aggregateAttestation.
type AggregateAttestation struct {
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

// aggregateAttestationJSON is a raw representation of the struct.
type aggregateAttestationJSON struct {
	Source          string `json:"source"`
	Method          string `json:"method"`
	Slot            string `json:"slot"`
	CommitteeIndex  string `json:"committee_index"`
	AggregationBits string `json:"aggregation_bits"`
	BeaconBlockRoot string `json:"beacon_block_root"`
	SourceRoot      string `json:"source_root"`
	TargetRoot      string `json:"target_root"`
	DelayMS         string `json:"delay_ms"`
}

// MarshalJSON implements json.Marshaler.
func (d *AggregateAttestation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&aggregateAttestationJSON{
		Source:          d.Source,
		Method:          d.Method,
		Slot:            fmt.Sprintf("%d", d.Slot),
		CommitteeIndex:  fmt.Sprintf("%d", d.CommitteeIndex),
		AggregationBits: fmt.Sprintf("%#x", d.AggregationBits),
		BeaconBlockRoot: fmt.Sprintf("%#x", d.BeaconBlockRoot),
		SourceRoot:      fmt.Sprintf("%#x", d.SourceRoot),
		TargetRoot:      fmt.Sprintf("%#x", d.TargetRoot),
		DelayMS:         fmt.Sprintf("%d", d.DelayMS),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *AggregateAttestation) UnmarshalJSON(input []byte) error {
	var data aggregateAttestationJSON
	err := json.Unmarshal(input, &data)
	if err != nil {
		return err
	}

	if data.Source == "" {
		return errors.New("source missing")
	}
	d.Source = data.Source

	if data.Method == "" {
		return errors.New("method missing")
	}
	d.Method = data.Method

	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 32)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	d.Slot = uint32(slot)

	if data.CommitteeIndex == "" {
		return errors.New("committee index missing")
	}
	committeeIndex, err := strconv.ParseUint(data.CommitteeIndex, 10, 16)
	if err != nil {
		return errors.Wrap(err, "invalid value for committee index")
	}
	d.CommitteeIndex = uint16(committeeIndex)

	if data.AggregationBits == "" {
		return errors.New("aggregation bits missing")
	}
	d.AggregationBits, err = hex.DecodeString(strings.TrimPrefix(data.AggregationBits, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for aggregation bits")
	}

	if data.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	d.BeaconBlockRoot, err = hex.DecodeString(strings.TrimPrefix(data.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}

	if data.SourceRoot == "" {
		return errors.New("source root root missing")
	}
	d.SourceRoot, err = hex.DecodeString(strings.TrimPrefix(data.SourceRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for source root root")
	}

	if data.TargetRoot == "" {
		return errors.New("target root root missing")
	}
	d.TargetRoot, err = hex.DecodeString(strings.TrimPrefix(data.TargetRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for target root root")
	}

	if data.DelayMS == "" {
		return errors.New("delay_ms missing")
	}
	delayMS, err := strconv.ParseUint(data.DelayMS, 10, 32)
	if err != nil {
		return errors.Wrap(err, "invalid value for delay_ms")
	}
	d.DelayMS = uint32(delayMS)

	return nil
}
