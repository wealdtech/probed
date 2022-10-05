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

// Package postgresql provides a postgresql implementation of the probe database.
package postgresql

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/wealdtech/probed/services/probedb"
)

// SetAggregateAttestation sets an aggregate attestation.
// If an aggregate attestation already exists then ignore it.
func (s *Service) SetAggregateAttestation(ctx context.Context, aggregateAttestation *probedb.AggregateAttestation) error {
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
	ip := aggregateAttestation.IPAddr.To4()
	if ip == nil {
		ip = aggregateAttestation.IPAddr
	}

	_, err := tx.Exec(ctx, `
INSERT INTO t_aggregate_attestations(f_ip_addr
                                    ,f_source
                                    ,f_method
                                    ,f_slot
                                    ,f_committee_index
                                    ,f_aggregation_bits
                                    ,f_beacon_block_root
                                    ,f_source_root
                                    ,f_target_root
                                    ,f_delay
                          )
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT (f_ip_addr, f_source, f_method, f_slot, f_committee_index, f_aggregation_bits) DO NOTHING
`,
		ip,
		aggregateAttestation.Source,
		aggregateAttestation.Method,
		aggregateAttestation.Slot,
		aggregateAttestation.CommitteeIndex,
		aggregateAttestation.AggregationBits,
		aggregateAttestation.BeaconBlockRoot,
		aggregateAttestation.SourceRoot,
		aggregateAttestation.TargetRoot,
		aggregateAttestation.DelayMS,
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

// AggregateAttestations obtains the aggregate attestations for a given filter.
func (s *Service) AggregateAttestations(
	ctx context.Context,
	filter *probedb.AggregateAttestationFilter,
) (
	[]*probedb.AggregateAttestation,
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

	queryBuilder.WriteString(`
SELECT f_ip_addr
      ,f_source
      ,f_method
      ,f_slot
      ,f_committee_index
      ,f_aggregation_bits
      ,f_beacon_block_root
      ,f_source_root
      ,f_target_root
      ,f_delay
FROM t_aggregate_attestations`)

	wherestr := "WHERE"

	if filter.IPAddr != "" {
		// Force the IP address to be a V4 if possible
		ipAddr := net.ParseIP(filter.IPAddr)
		ip := ipAddr.To4()
		if ip == nil {
			ip = ipAddr
		}
		queryVals = append(queryVals, ip)
		queryBuilder.WriteString(fmt.Sprintf(`
%s f_ip_addr = $%d`, wherestr, len(queryVals)))
		wherestr = "  AND"
	}

	if len(filter.Sources) > 0 {
		queryVals = append(queryVals, filter.Sources)
		queryBuilder.WriteString(fmt.Sprintf(`
%s f_source = ANY($%d)`, wherestr, len(queryVals)))
		wherestr = "  AND"
	}

	if len(filter.Methods) > 0 {
		queryVals = append(queryVals, filter.Methods)
		queryBuilder.WriteString(fmt.Sprintf(`
%s f_method = $%d`, wherestr, len(queryVals)))
		wherestr = "  AND"
	}

	if filter.From != nil {
		queryVals = append(queryVals, *filter.From)
		queryBuilder.WriteString(fmt.Sprintf(`
%s f_slot >= $%d`, wherestr, len(queryVals)))
		wherestr = "  AND"
	}

	if filter.To != nil {
		queryVals = append(queryVals, *filter.To)
		queryBuilder.WriteString(fmt.Sprintf(`
%s f_slot <= $%d`, wherestr, len(queryVals)))
		// wherestr = "  AND"
	}

	switch filter.Order {
	case probedb.OrderEarliest:
		queryBuilder.WriteString(`
ORDER BY f_slot`)
	case probedb.OrderLatest:
		queryBuilder.WriteString(`
ORDER BY f_slot DESC`)
	default:
		return nil, errors.New("no order specified")
	}

	if filter.Limit != 0 {
		queryVals = append(queryVals, filter.Limit)
		queryBuilder.WriteString(fmt.Sprintf(`
LIMIT $%d`, len(queryVals)))
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

	aggregateAttestations := make([]*probedb.AggregateAttestation, 0)
	for rows.Next() {
		aggregateAttestation := &probedb.AggregateAttestation{}
		err := rows.Scan(
			&aggregateAttestation.IPAddr,
			&aggregateAttestation.Source,
			&aggregateAttestation.Method,
			&aggregateAttestation.Slot,
			&aggregateAttestation.CommitteeIndex,
			&aggregateAttestation.AggregationBits,
			&aggregateAttestation.BeaconBlockRoot,
			&aggregateAttestation.SourceRoot,
			&aggregateAttestation.TargetRoot,
			&aggregateAttestation.DelayMS,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		ip := aggregateAttestation.IPAddr.To4()
		if ip == nil {
			ip = aggregateAttestation.IPAddr
		}
		aggregateAttestation.IPAddr = ip
		aggregateAttestations = append(aggregateAttestations, aggregateAttestation)
	}
	return aggregateAttestations, nil
}
