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

package postgresql_test

import (
	"context"
	"net"
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

	headDelay := &probedb.Delay{
		IPAddr:  net.ParseIP("1.2.3.4"),
		Source:  "Dummy client",
		Method:  "test",
		Slot:    12345,
		DelayMS: 234,
	}

	// Set the head delay.
	require.NoError(t, s.SetHeadDelay(ctx, headDelay))

	// Attempt to overwrite; should be ignored but no error.
	headDelay.DelayMS = 345
	require.NoError(t, s.SetHeadDelay(ctx, headDelay))
}

func TestHeadDelays(t *testing.T) {
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

	headDelays := []*probedb.Delay{
		{IPAddr: net.ParseIP("1.2.3.4"), Source: "Source 1", Method: "Method 1", Slot: 12345, DelayMS: 1123},
		{IPAddr: net.ParseIP("1.2.3.4"), Source: "Source 2", Method: "Method 1", Slot: 12345, DelayMS: 1234},
		{IPAddr: net.ParseIP("1.2.3.4"), Source: "Source 3", Method: "Method 1", Slot: 12345, DelayMS: 1345},
		{IPAddr: net.ParseIP("2.3.4.5"), Source: "Source 1", Method: "Method 2", Slot: 12345, DelayMS: 1456},
		{IPAddr: net.ParseIP("2.3.4.5"), Source: "Source 2", Method: "Method 2", Slot: 12345, DelayMS: 1567},
		{IPAddr: net.ParseIP("2.3.4.5"), Source: "Source 3", Method: "Method 2", Slot: 12345, DelayMS: 1678},
		{IPAddr: net.ParseIP("1.2.3.4"), Source: "Source 1", Method: "Method 1", Slot: 12346, DelayMS: 2123},
		{IPAddr: net.ParseIP("1.2.3.4"), Source: "Source 2", Method: "Method 1", Slot: 12346, DelayMS: 2234},
		{IPAddr: net.ParseIP("1.2.3.4"), Source: "Source 3", Method: "Method 1", Slot: 12346, DelayMS: 2345},
		{IPAddr: net.ParseIP("2.3.4.5"), Source: "Source 1", Method: "Method 2", Slot: 12346, DelayMS: 2456},
		{IPAddr: net.ParseIP("2.3.4.5"), Source: "Source 2", Method: "Method 2", Slot: 12346, DelayMS: 2567},
		{IPAddr: net.ParseIP("2.3.4.5"), Source: "Source 3", Method: "Method 2", Slot: 12346, DelayMS: 2678},
	}

	// Set the head delays.
	for _, headDelay := range headDelays {
		require.NoError(t, s.SetHeadDelay(ctx, headDelay))
	}

	tests := []struct {
		name   string
		filter *probedb.DelayFilter
		res    []*probedb.DelayValue
	}{
		{
			name:   "Default",
			filter: &probedb.DelayFilter{},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1123},
				{Slot: 12346, DelayMS: 2123},
			},
		},
		{
			name: "SingleSlot",
			filter: &probedb.DelayFilter{
				From: slotPtr(12345),
				To:   slotPtr(12345),
			},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1123},
			},
		},
		{
			name: "MinSlot",
			filter: &probedb.DelayFilter{
				From: slotPtr(12346),
			},
			res: []*probedb.DelayValue{
				{Slot: 12346, DelayMS: 2123},
			},
		},
		{
			name: "MaxSlot",
			filter: &probedb.DelayFilter{
				To: slotPtr(12345),
			},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1123},
			},
		},
		{
			name: "Median",
			filter: &probedb.DelayFilter{
				Selection: probedb.SelectionMedian,
			},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1400},
				{Slot: 12346, DelayMS: 2400},
			},
		},
		{
			name: "Maximum",
			filter: &probedb.DelayFilter{
				Selection: probedb.SelectionMaximum,
			},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1678},
				{Slot: 12346, DelayMS: 2678},
			},
		},
		{
			name: "IPAddrFilter",
			filter: &probedb.DelayFilter{
				IPAddr: "2.3.4.5",
			},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1456},
				{Slot: 12346, DelayMS: 2456},
			},
		},
		{
			name: "SourceFilter",
			filter: &probedb.DelayFilter{
				Source: "Source 2",
			},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1234},
				{Slot: 12346, DelayMS: 2234},
			},
		},
		{
			name: "MethodFilter",
			filter: &probedb.DelayFilter{
				Method: "Method 2",
			},
			res: []*probedb.DelayValue{
				{Slot: 12345, DelayMS: 1456},
				{Slot: 12346, DelayMS: 2456},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := s.HeadDelays(ctx, test.filter)
			require.NoError(t, err)
			require.Equal(t, len(test.res), len(res))
			for i := range test.res {
				require.Equal(t, test.res[i], res[i])
			}
		})
	}
}
