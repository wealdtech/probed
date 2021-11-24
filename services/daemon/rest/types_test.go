// Copyright © 2021 Attestant Limited.
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

package rest_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wealdtech/probed/services/daemon/rest"
	"gotest.tools/assert"
)

func TestDelayJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		res   *rest.Delay
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "json: cannot unmarshal array into Go value of type rest.delayJSON",
		},
		{
			name:  "LocationIDMissing",
			input: []byte(`{"source_id":"2","method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "location_id missing",
		},
		{
			name:  "LocationIDWrongType",
			input: []byte(`{"location_id":true,"source_id":"2","method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.location_id of type string",
		},
		{
			name:  "LocationIDInvalid",
			input: []byte(`{"location_id":"-1","source_id":"2","method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "invalid value for location_id: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SourceIDMissing",
			input: []byte(`{"location_id":"1","method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "source_id missing",
		},
		{
			name:  "SourceIDWrongType",
			input: []byte(`{"location_id":"1","source_id":true,"method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.source_id of type string",
		},
		{
			name:  "SourceIDInvalid",
			input: []byte(`{"location_id":"1","source_id":"-2","method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "invalid value for source_id: strconv.ParseUint: parsing \"-2\": invalid syntax",
		},
		{
			name:  "MethodMissing",
			input: []byte(`{"location_id":"1","source_id":"2","slot":"123","delay_ms":"12345"}`),
			err:   "method missing",
		},
		{
			name:  "MethodWrongType",
			input: []byte(`{"location_id":"1","source_id":"2","method":true,"slot":"123","delay_ms":"12345"}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.method of type string",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"location_id":"1","source_id":"2","method":"head event","delay_ms":"12345"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"location_id":"1","source_id":"2","method":"head event","slot":true,"delay_ms":"12345"}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.slot of type string",
		},
		{
			name:  "SlotWnvalid",
			input: []byte(`{"location_id":"1","source_id":"2","method":"head event","slot":"-1","delay_ms":"12345"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "DelayMSMissing",
			input: []byte(`{"location_id":"1","source_id":"2","method":"head event","slot":"123"}`),
			err:   "delay_ms missing",
		},
		{
			name:  "DelayMSWrongType",
			input: []byte(`{"location_id":"1","source_id":"2","method":"head event","slot":"123","delay_ms":true}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.delay_ms of type string",
		},
		{
			name:  "DelayMSInvalid",
			input: []byte(`{"location_id":"1","source_id":"2","method":"head event","slot":"123","delay_ms":"-1"}`),
			err:   "invalid value for delay_ms: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"location_id":"1","source_id":"2","method":"head event","slot":"123","delay_ms":"12345"}`),
			res: &rest.Delay{
				LocationID: 1,
				SourceID:   2,
				Method:     "head event",
				Slot:       123,
				DelayMS:    12345,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res rest.Delay
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				require.Equal(t, test.res.LocationID, res.LocationID)
				require.Equal(t, test.res.SourceID, res.SourceID)
				require.Equal(t, test.res.Method, res.Method)
				require.Equal(t, test.res.Slot, res.Slot)
				require.Equal(t, test.res.DelayMS, res.DelayMS)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
