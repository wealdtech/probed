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
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/wealdtech/probed/services/probedb"
)

// AttestationSummaries obtains the attestation summaries for a filter.
func (s *Service) AttestationSummaries(ctx context.Context,
	filter *probedb.AttestationSummaryFilter,
) (
	[]*probedb.AttestationSummary,
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
      ,f_beacon_block_root
      ,f_source_root
      ,f_target_root
      ,f_attester_buckets
FROM t_attestation_summaries`)

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
		conditions = append(conditions, fmt.Sprintf(`f_method = $%d`, len(queryVals)))
	}

	if filter.From != nil {
		queryVals = append(queryVals, *filter.From)
		conditions = append(conditions, fmt.Sprintf(`f_slot >= $%d`, len(queryVals)))
	}

	if filter.To != nil {
		queryVals = append(queryVals, *filter.To)
		queryBuilder.WriteString(fmt.Sprintf(` f_slot <= $%d`, len(queryVals)))
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString("\nWHERE ")
		queryBuilder.WriteString(strings.Join(conditions, "\n  AND "))
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

	attestationSummaries := make([]*probedb.AttestationSummary, 0)
	for rows.Next() {
		attestationSummary := &probedb.AttestationSummary{}
		err := rows.Scan(
			&attestationSummary.IPAddr,
			&attestationSummary.Source,
			&attestationSummary.Method,
			&attestationSummary.Slot,
			&attestationSummary.CommitteeIndex,
			&attestationSummary.BeaconBlockRoot,
			&attestationSummary.SourceRoot,
			&attestationSummary.TargetRoot,
			&attestationSummary.AttesterBuckets,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		ip := attestationSummary.IPAddr.To4()
		if ip != nil {
			attestationSummary.IPAddr = ip
		}
		attestationSummaries = append(attestationSummaries, attestationSummary)
	}
	return attestationSummaries, nil
}
