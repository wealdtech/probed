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

package postgresql

import (
	"context"

	"github.com/wealdtech/probed/services/probedb"
)

// SetBlockDelay sets a block delay.
func (s *Service) SetBlockDelay(ctx context.Context, blockDelay *probedb.BlockDelay) error {
	localTx := false
	tx := s.tx(ctx)
	if tx == nil {
		var err error
		tx, err = s.pool.Begin(ctx)
		if err != nil {
			return err
		}
		localTx = true
	}

	_, err := tx.Exec(ctx, `
INSERT INTO t_block_delay(f_location_id
                         ,f_source_id
                         ,f_method
                         ,f_slot
                         ,f_delay
                         )
VALUES($1,$2,$3,$4,$5)
ON CONFLICT (f_location_id,f_source_id,f_method,f_slot) DO
UPDATE
SET f_delay = excluded.f_delay
  `,
		blockDelay.LocationID,
		blockDelay.SourceID,
		blockDelay.Method,
		blockDelay.Slot,
		blockDelay.DelayMS,
	)

	if localTx {
		if err == nil {
			if err := tx.Commit(ctx); err != nil {
				log.Warn().Err(err).Msg("Failed to commit transaction")
			}
		} else {
			if err := tx.Rollback(ctx); err != nil {
				log.Warn().Err(err).Msg("Failed to rollback transaction")
			}
		}
	}

	return err
}
