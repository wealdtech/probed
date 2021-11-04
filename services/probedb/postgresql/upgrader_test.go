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
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wealdtech/probed/services/probedb/postgresql"
)

func TestUpgrader(t *testing.T) {
	ctx := context.Background()

	s, err := postgresql.New(ctx,
		postgresql.WithServer(os.Getenv("PROBEDB_SERVER")),
		postgresql.WithPort(atoi(os.Getenv("PROBEDB_PORT"))),
		postgresql.WithUser(os.Getenv("PROBEDB_USER")),
		postgresql.WithPassword(os.Getenv("PROBEDB_PASSWORD")),
	)
	require.NoError(t, err)

	// Ensure upgrader runs.
	require.NoError(t, s.Upgrade(ctx))
	// Ensure repeat run does not error.
	require.NoError(t, s.Upgrade(ctx))
}
