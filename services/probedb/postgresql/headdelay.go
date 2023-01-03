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

package postgresql

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/wealdtech/probed/services/probedb"
)

// SetHeadDelay sets a head delay.
// If a delay already exists for this head then ignore it.
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

	// Force the IP address to be a V4 if possible
	ip := delay.IPAddr.To4()
	if ip == nil {
		ip = delay.IPAddr
	}

	_, err := tx.Exec(ctx, `
INSERT INTO t_head_delays(f_ip_addr
                         ,f_source
                         ,f_method
                         ,f_slot
                         ,f_delay
                         )
VALUES($1,$2,$3,$4,$5)
ON CONFLICT (f_ip_addr, f_source, f_method, f_slot) DO NOTHING
`,
		ip,
		delay.Source,
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

// HeadDelays obtains the head delays for a range of slots.
func (s *Service) HeadDelays(
	ctx context.Context,
	filter *probedb.DelayFilter,
) (
	[]*probedb.Delay,
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
	queryVals := make([]interface{}, 0)

	switch filter.Selection {
	case probedb.SelectionMinimum:
		queryBuilder.WriteString(`
SELECT f_slot
      ,MIN(f_delay)`)
	case probedb.SelectionMaximum:
		queryBuilder.WriteString(`
SELECT f_slot
      ,MAX(f_delay)`)
	case probedb.SelectionMedian:
		queryBuilder.WriteString(`
SELECT f_slot
      ,(PERCENTILE_CONT(0.5) WITHIN GROUP(ORDER BY f_delay))::INT`)
	case probedb.SelectionAll:
		queryBuilder.WriteString(`
SELECT f_ip_addr
      ,f_source
      ,f_method
      ,f_slot
      ,f_delay`)
	default:
		return nil, errors.New("unhandled selection criteria")
	}

	queryBuilder.WriteString(`
FROM t_head_delays`)

	conditions := make([]string, 0)

	if filter.IPAddr != "" {
		// Force the IP address to be a V4 if possible
		ipAddr := net.ParseIP(filter.IPAddr)
		ip := ipAddr.To4()
		if ip == nil {
			ip = ipAddr
		}
		queryVals = append(queryVals, ip)
		conditions = append(conditions, fmt.Sprintf(`f_ip_addr = $%d`, len(queryVals)))
	}

	if len(filter.Sources) > 0 {
		queryVals = append(queryVals, filter.Sources)
		conditions = append(conditions, fmt.Sprintf(`f_source = ANY($%d)`, len(queryVals)))
	}

	if len(filter.Methods) > 0 {
		queryVals = append(queryVals, filter.Methods)
		conditions = append(conditions, fmt.Sprintf(`f_method = ANY($%d)`, len(queryVals)))
	}

	if filter.From != nil {
		queryVals = append(queryVals, *filter.From)
		conditions = append(conditions, fmt.Sprintf(`f_slot >= $%d`, len(queryVals)))
	}

	if filter.To != nil {
		queryVals = append(queryVals, *filter.To)
		conditions = append(conditions, fmt.Sprintf(`f_slot <= $%d`, len(queryVals)))
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString("\nWHERE ")
		queryBuilder.WriteString(strings.Join(conditions, "\n  AND "))
	}

	if filter.Selection == probedb.SelectionAll {
		queryBuilder.WriteString(`
ORDER BY f_slot
        ,f_method
        ,f_ip_addr
        ,f_source`)
	} else {
		queryBuilder.WriteString(`
GROUP BY f_slot
ORDER BY f_slot
`)
	}

	if e := log.Trace(); e.Enabled() {
		params := make([]string, len(queryVals))
		for i := range queryVals {
			params[i] = fmt.Sprintf("%v", queryVals[i])
		}
		log.Trace().Str("query", strings.ReplaceAll(queryBuilder.String(), "\n", " ")).Strs("params", params).Msg("SQL query")
	}

	rows, err := tx.Query(ctx,
		queryBuilder.String(),
		queryVals...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	delays := make([]*probedb.Delay, 0)
	for rows.Next() {
		delay := &probedb.Delay{}
		if filter.Selection == probedb.SelectionAll {
			err = rows.Scan(
				&delay.IPAddr,
				&delay.Source,
				&delay.Method,
				&delay.Slot,
				&delay.DelayMS,
			)
		} else {
			err = rows.Scan(
				&delay.Slot,
				&delay.DelayMS,
			)
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		if len(delay.IPAddr) > 0 {
			ip := delay.IPAddr.To4()
			if ip != nil {
				delay.IPAddr = ip
			}
		}
		delays = append(delays, delay)
	}
	return delays, nil
}
