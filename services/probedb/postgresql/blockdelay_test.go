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

func TestSetBlockDelay(t *testing.T) {
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

	blockDelay := &probedb.BlockDelay{
		LocationID: 1,
		SourceID:   2,
		Method:     "test",
		Slot:       12345,
		DelayMS:    234,
	}

	// Set the block delay.
	require.NoError(t, s.SetBlockDelay(ctx, blockDelay))

	// Overwrite the old values.
	blockDelay.DelayMS = 345
	require.NoError(t, s.SetBlockDelay(ctx, blockDelay))
}
