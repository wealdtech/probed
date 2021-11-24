// Copyright © 2021 Weald Technology Limited.
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

package postgresql_test

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/probed/services/probedb"
	"github.com/wealdtech/probed/services/probedb/postgresql"
)

func TestSetHeadDelay(t *testing.T) {
	ctx := context.Background()
	s, err := postgresql.New(ctx,
		postgresql.WithLogLevel(zerolog.Disabled),
		postgresql.WithServer(os.Getenv("PROBEDB_SERVER")),
		postgresql.WithPort(atoi(os.Getenv("PROBEDB_PORT"))),
		postgresql.WithUser(os.Getenv("PROBEDB_USER")),
		postgresql.WithPassword(os.Getenv("PROBEDB_PASSWORD")),
	)
	require.NoError(t, err)

	ctx, cancel, err := s.BeginTx(ctx)
	require.NoError(t, err)
	defer cancel()

	delay := &probedb.Delay{
		LocationID: 1,
		SourceID:   2,
		Method:     "test",
		Slot:       12345,
		DelayMS:    234,
	}

	// Set the block delay.
	require.NoError(t, s.SetHeadDelay(ctx, delay))

	// Overwrite the old values.
	delay.DelayMS = 345
	require.NoError(t, s.SetHeadDelay(ctx, delay))
}

func TestMedianHeadDelay(t *testing.T) {
	ctx := context.Background()
	s, err := postgresql.New(ctx,
		postgresql.WithLogLevel(zerolog.Disabled),
		postgresql.WithServer(os.Getenv("PROBEDB_SERVER")),
		postgresql.WithPort(atoi(os.Getenv("PROBEDB_PORT"))),
		postgresql.WithUser(os.Getenv("PROBEDB_USER")),
		postgresql.WithPassword(os.Getenv("PROBEDB_PASSWORD")),
	)
	require.NoError(t, err)

	ctx, cancel, err := s.BeginTx(ctx)
	require.NoError(t, err)
	defer cancel()

	blockDelays := []*probedb.Delay{
		{LocationID: 1, SourceID: 1, Method: "test", Slot: 12345, DelayMS: 123},
		{LocationID: 1, SourceID: 2, Method: "test", Slot: 12345, DelayMS: 234},
		{LocationID: 1, SourceID: 3, Method: "test", Slot: 12345, DelayMS: 345},
		{LocationID: 2, SourceID: 1, Method: "test", Slot: 12345, DelayMS: 456},
		{LocationID: 2, SourceID: 2, Method: "test", Slot: 12345, DelayMS: 567},
		{LocationID: 2, SourceID: 3, Method: "test", Slot: 12345, DelayMS: 678},
	}

	// Set the block delays.
	for _, blockDelay := range blockDelays {
		require.NoError(t, s.SetHeadDelay(ctx, blockDelay))
	}

	tests := []struct {
		name       string
		locationID uint16
		sourceID   uint16
		method     string
		fromSlot   uint32
		toSlot     uint32
		res        []*probedb.DelayValue
	}{
		{
			name:     "All",
			fromSlot: 1,
			toSlot:   99999,
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 400},
			},
		},
		{
			name:       "Location1",
			locationID: 1,
			fromSlot:   1,
			toSlot:     99999,
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 234},
			},
		},
		{
			name:     "Source1",
			sourceID: 1,
			fromSlot: 1,
			toSlot:   99999,
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 290},
			},
		},
		{
			name:     "MethodTest",
			method:   "test",
			fromSlot: 1,
			toSlot:   99999,
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 400},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := s.MedianHeadDelays(ctx,
				test.locationID,
				test.sourceID,
				test.method,
				test.fromSlot,
				test.toSlot,
			)
			require.NoError(t, err)
			for i := range test.res {
				require.Equal(t, test.res[i], res[i])
			}
		})
	}
}
