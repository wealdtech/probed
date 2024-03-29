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

package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wealdtech/probed/services/daemon/rest/types"
	"gotest.tools/assert"
)

func TestDelayJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		res   *types.Delay
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "json: cannot unmarshal array into Go value of type types.delayJSON",
		},
		{
			name:  "SourceMissing",
			input: []byte(`{"method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "source missing",
		},
		{
			name:  "SourceWrongType",
			input: []byte(`{"source":true,"method":"head event","slot":"123","delay_ms":"12345"}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.source of type string",
		},
		{
			name:  "MethodMissing",
			input: []byte(`{"source":"client","slot":"123","delay_ms":"12345"}`),
			err:   "method missing",
		},
		{
			name:  "MethodWrongType",
			input: []byte(`{"source":"client","method":true,"slot":"123","delay_ms":"12345"}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.method of type string",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"source":"client","method":"head event","delay_ms":"12345"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"source":"client","method":"head event","slot":true,"delay_ms":"12345"}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"source":"client","method":"head event","slot":"-1","delay_ms":"12345"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "DelayMSMissing",
			input: []byte(`{"source":"client","method":"head event","slot":"123"}`),
			err:   "delay_ms missing",
		},
		{
			name:  "DelayMSWrongType",
			input: []byte(`{"source":"client","method":"head event","slot":"123","delay_ms":true}`),
			err:   "json: cannot unmarshal bool into Go struct field delayJSON.delay_ms of type string",
		},
		{
			name:  "DelayMSInvalid",
			input: []byte(`{"source":"client","method":"head event","slot":"123","delay_ms":"-1"}`),
			err:   "invalid value for delay_ms: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"source":"client","method":"head event","slot":"123","delay_ms":"12345"}`),
			res: &types.Delay{
				Source:  "client",
				Method:  "head event",
				Slot:    123,
				DelayMS: 12345,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res types.Delay
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				require.Equal(t, test.res.Source, res.Source)
				require.Equal(t, test.res.Method, res.Method)
				require.Equal(t, test.res.Slot, res.Slot)
				require.Equal(t, test.res.DelayMS, res.DelayMS)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
