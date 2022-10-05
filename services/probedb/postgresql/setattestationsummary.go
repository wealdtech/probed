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

package postgresql

import (
	"context"

	"github.com/wealdtech/probed/services/probedb"
)

// SetAttestationSummary sets an attestation summary.
func (s *Service) SetAttestationSummary(ctx context.Context, summary *probedb.AttestationSummary) error {
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

	// Force the IP address to be a V4 if possible
	ip := summary.IPAddr.To4()
	if ip == nil {
		ip = summary.IPAddr
	}

	_, err := tx.Exec(ctx, `
INSERT INTO t_attestation_summaries(f_ip_addr
                                   ,f_source
                                   ,f_method
                                   ,f_slot
                                   ,f_committee_index
                                   ,f_beacon_block_root
                                   ,f_source_root
                                   ,f_target_root
                                   ,f_attester_buckets
                                   )
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
ON CONFLICT (f_ip_addr, f_source, f_method, f_slot, f_committee_index, f_beacon_block_root, f_source_root, f_target_root) DO
NOTHING
-- UPDATE
-- SET f_attester_buckets = excluded.f_attester_buckets
`,
		ip,
		summary.Source,
		summary.Method,
		summary.Slot,
		summary.CommitteeIndex,
		summary.BeaconBlockRoot,
		summary.SourceRoot,
		summary.TargetRoot,
		summary.AttesterBuckets,
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
