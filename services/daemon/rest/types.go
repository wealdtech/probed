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

package rest

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// BlockDelay holds information about a block delay.
type BlockDelay struct {
	LocationID uint16
	SourceID   uint16
	Method     string
	Slot       uint32
	DelayMS    uint32
}

// blockDelayJSON is a raw representation of the struct.
type blockDelayJSON struct {
	LocationID string `json:"location_id"`
	SourceID   string `json:"source_id"`
	Method     string `json:"method"`
	Slot       string `json:"slot"`
	DelayMS    string `json:"delay_ms"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlockDelay) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blockDelayJSON{
		LocationID: fmt.Sprintf("%d", b.LocationID),
		SourceID:   fmt.Sprintf("%d", b.SourceID),
		Method:     b.Method,
		Slot:       fmt.Sprintf("%d", b.Slot),
		DelayMS:    fmt.Sprintf("%d", b.DelayMS),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlockDelay) UnmarshalJSON(input []byte) error {
	var blockDelayJSON blockDelayJSON
	err := json.Unmarshal(input, &blockDelayJSON)
	if err != nil {
		return err
	}

	if blockDelayJSON.LocationID == "" {
		return errors.New("location_id missing")
	}
	locationID, err := strconv.ParseUint(blockDelayJSON.LocationID, 10, 16)
	if err != nil {
		return errors.Wrap(err, "invalid value for location_id")
	}
	b.LocationID = uint16(locationID)

	if blockDelayJSON.SourceID == "" {
		return errors.New("source_id missing")
	}
	sourceID, err := strconv.ParseUint(blockDelayJSON.SourceID, 10, 16)
	if err != nil {
		return errors.Wrap(err, "invalid value for source_id")
	}
	b.SourceID = uint16(sourceID)

	if blockDelayJSON.Method == "" {
		return errors.New("method missing")
	}
	b.Method = blockDelayJSON.Method

	if blockDelayJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(blockDelayJSON.Slot, 10, 32)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = uint32(slot)

	if blockDelayJSON.DelayMS == "" {
		return errors.New("delay_ms missing")
	}
	delayMS, err := strconv.ParseUint(blockDelayJSON.DelayMS, 10, 32)
	if err != nil {
		return errors.Wrap(err, "invalid value for delay_ms")
	}
	b.DelayMS = uint32(delayMS)

	return nil
}
