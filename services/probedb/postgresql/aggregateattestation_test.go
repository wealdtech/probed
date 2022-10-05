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

func parseIP(input string) net.IP {
	ipAddr := net.ParseIP(input)
	ip := ipAddr.To4()
	if ip == nil {
		ip = ipAddr
	}
	return ip
}

func TestAggregateAttestations(t *testing.T) {
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

	aggregateAttestations := []*probedb.AggregateAttestation{
		{IPAddr: parseIP("1.2.3.4"), Source: "Source 1", Method: "Method 1", Slot: 12345, CommitteeIndex: 1, AggregationBits: []byte{0x01, 0x10}, BeaconBlockRoot: []byte{0x01}, SourceRoot: []byte{0x02}, TargetRoot: []byte{0x03}, DelayMS: 1123},
		{IPAddr: parseIP("1.2.3.4"), Source: "Source 1", Method: "Method 1", Slot: 12346, CommitteeIndex: 1, AggregationBits: []byte{0x01, 0x10}, BeaconBlockRoot: []byte{0x01}, SourceRoot: []byte{0x02}, TargetRoot: []byte{0x03}, DelayMS: 1345},
	}

	// Set the head delays.
	for _, aggregateAttestation := range aggregateAttestations {
		require.NoError(t, s.SetAggregateAttestation(ctx, aggregateAttestation))
	}

	tests := []struct {
		name   string
		filter *probedb.AggregateAttestationFilter
		res    []*probedb.AggregateAttestation
	}{
		{
			name: "SingleSource",
			filter: &probedb.AggregateAttestationFilter{
				Sources: []string{"Source 1"},
			},
			res: []*probedb.AggregateAttestation{
				aggregateAttestations[0],
				aggregateAttestations[1],
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := s.AggregateAttestations(ctx, test.filter)
			require.NoError(t, err)
			require.Equal(t, len(test.res), len(res))
			for i := range test.res {
				require.Equal(t, test.res[i], res[i])
			}
		})
	}
}
