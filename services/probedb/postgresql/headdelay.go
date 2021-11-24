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
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/wealdtech/probed/services/probedb"
)

// SetHeadDelay sets a head delay.
func (s *Service) SetHeadDelay(ctx context.Context, delay *probedb.Delay) error {
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
INSERT INTO t_head_delay(f_location_id
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
		delay.LocationID,
		delay.SourceID,
		delay.Method,
		delay.Slot,
		delay.DelayMS,
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

// MedianHeadDelays obtains the median head delays for a range of slots.
func (s *Service) MedianHeadDelays(ctx context.Context,
	locationID uint16,
	sourceID uint16,
	method string,
	fromSlot uint32,
	toSlot uint32,
) (
	[]*probedb.DelayValue,
	error,
) {
	tx := s.tx(ctx)
	if tx == nil {
		ctx, cancel, err := s.BeginTx(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to begin transaction")
		}
		tx = s.tx(ctx)
		defer cancel()
	}

	// Build the query.
	queryBuilder := strings.Builder{}
	queryVals := make([]interface{}, 2)

	queryVals[0] = fromSlot
	queryVals[1] = toSlot
	queryBuilder.WriteString(`
SELECT f_slot
      ,(PERCENTILE_CONT(0.5) WITHIN GROUP(ORDER BY f_delay))::INT
FROM t_head_delay
WHERE f_slot >= $1
  AND f_slot < $2`)

	if locationID != 0 {
		queryVals = append(queryVals, locationID)
		queryBuilder.WriteString(fmt.Sprintf(`
  AND f_location_id = $%d`, len(queryVals)))
	}

	if sourceID != 0 {
		queryVals = append(queryVals, sourceID)
		queryBuilder.WriteString(fmt.Sprintf(`
  AND f_source_id = $%d`, len(queryVals)))
	}

	if method != "" {
		queryVals = append(queryVals, method)
		queryBuilder.WriteString(fmt.Sprintf(`
  AND f_method = $%d`, len(queryVals)))
	}

	queryBuilder.WriteString(`
GROUP BY f_slot
ORDER BY f_slot`)

	rows, err := tx.Query(ctx,
		queryBuilder.String(),
		queryVals...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	delays := make([]*probedb.DelayValue, 0)
	for rows.Next() {
		delay := &probedb.DelayValue{}
		err := rows.Scan(&delay.Slot, &delay.DelayMS)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		delays = append(delays, delay)
	}
	return delays, nil
}

// MinimumHeadDelays obtains the minimum head delays for a range of slots.
func (s *Service) MinimumHeadDelays(ctx context.Context,
	locationID uint16,
	sourceID uint16,
	method string,
	fromSlot uint32,
	toSlot uint32,
) (
	[]*probedb.DelayValue,
	error,
) {
	tx := s.tx(ctx)
	if tx == nil {
		ctx, cancel, err := s.BeginTx(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to begin transaction")
		}
		tx = s.tx(ctx)
		defer cancel()
	}

	// Build the query.
	queryBuilder := strings.Builder{}
	queryVals := make([]interface{}, 2)

	queryVals[0] = fromSlot
	queryVals[1] = toSlot
	queryBuilder.WriteString(`
SELECT f_slot
      ,MIN(f_delay)
FROM t_head_delay
WHERE f_slot >= $1
  AND f_slot < $2`)

	if locationID != 0 {
		queryVals = append(queryVals, locationID)
		queryBuilder.WriteString(fmt.Sprintf(`
  AND f_location_id = $%d`, len(queryVals)))
	}

	if sourceID != 0 {
		queryVals = append(queryVals, sourceID)
		queryBuilder.WriteString(fmt.Sprintf(`
  AND f_source_id = $%d`, len(queryVals)))
	}

	if method != "" {
		queryVals = append(queryVals, method)
		queryBuilder.WriteString(fmt.Sprintf(`
  AND f_method = $%d`, len(queryVals)))
	}

	queryBuilder.WriteString(`
GROUP BY f_slot
ORDER BY f_slot`)

	rows, err := tx.Query(ctx,
		queryBuilder.String(),
		queryVals...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	delays := make([]*probedb.DelayValue, 0)
	for rows.Next() {
		delay := &probedb.DelayValue{}
		err := rows.Scan(&delay.Slot, &delay.DelayMS)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		delays = append(delays, delay)
	}
	return delays, nil
}
