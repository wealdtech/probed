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

package postgresql_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/probed/services/probedb"
	"github.com/wealdtech/probed/services/probedb/postgresql"
)

func atoi(input string) int32 {
	val, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		val = -1
	}
	return int32(val)
}

func TestService(t *testing.T) {
	tests := []struct {
		name     string
		server   string
		port     int32
		user     string
		password string
		err      string
	}{
		{
			name:     "ServerMissing",
			err:      "problem with parameters: no server specified",
			port:     atoi(os.Getenv("PROBEDB_PORT")),
			user:     os.Getenv("PROBEDB_USER"),
			password: os.Getenv("PROBEDB_PASSWORD"),
		},
		{
			name:     "PortMissing",
			err:      "problem with parameters: no port specified",
			server:   os.Getenv("PROBEDB_SERVER"),
			user:     os.Getenv("PROBEDB_USER"),
			password: os.Getenv("PROBEDB_PASSWORD"),
		},
		{
			name:     "UserMissing",
			err:      "problem with parameters: no user specified",
			server:   os.Getenv("PROBEDB_SERVER"),
			port:     atoi(os.Getenv("PROBEDB_PORT")),
			password: os.Getenv("PROBEDB_PASSWORD"),
		},
		{
			name:     "Good",
			server:   os.Getenv("PROBEDB_SERVER"),
			port:     atoi(os.Getenv("PROBEDB_PORT")),
			user:     os.Getenv("PROBEDB_USER"),
			password: os.Getenv("PROBEDB_PASSWORD"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := postgresql.New(ctx,
				postgresql.WithServer(test.server),
				postgresql.WithPort(test.port),
				postgresql.WithUser(test.user),
				postgresql.WithPassword(test.password),
			)
			if test.err != "" {
				assert.EqualError(t, err, test.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInterfaces(t *testing.T) {
	ctx := context.Background()

	s, err := postgresql.New(ctx,
		postgresql.WithServer(os.Getenv("PROBEDB_SERVER")),
		postgresql.WithPort(atoi(os.Getenv("PROBEDB_PORT"))),
		postgresql.WithUser(os.Getenv("PROBEDB_USER")),
		postgresql.WithPassword(os.Getenv("PROBEDB_PASSWORD")),
	)
	require.NoError(t, err)

	require.Implements(t, (*probedb.Service)(nil), s)
	require.Implements(t, (*probedb.BlockDelaysSetter)(nil), s)
}
