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

package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wealdtech/probed/services/daemon/rest/types"
	"github.com/wealdtech/probed/services/probedb"
)

func (s *Service) postAttestationSummary(w http.ResponseWriter, r *http.Request) {
	var summary types.AttestationSummary
	if err := json.NewDecoder(r.Body).Decode(&summary); err != nil {
		log.Debug().Err(err).Msg("Supplied with invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sourceIP, err := sourceIP(r)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to obtain source IP")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Need to store attestations on a pre-source basis.
	for _, attestation := range summary.Attestations {
		for source, buckets := range attestation.Buckets {
			dbBuckets := make([][]byte, 0, len(buckets))
			for _, bucket := range buckets {
				dbBuckets = append(dbBuckets, bucket)
			}

			if err := s.attestationSummariesSetter.SetAttestationSummary(context.Background(), &probedb.AttestationSummary{
				IPAddr:          sourceIP,
				Source:          source,
				Method:          summary.Method,
				Slot:            summary.Slot,
				CommitteeIndex:  attestation.CommitteeIndex,
				BeaconBlockRoot: attestation.BeaconBlockRoot,
				SourceRoot:      attestation.SourceRoot,
				TargetRoot:      attestation.TargetRoot,
				AttesterBuckets: dbBuckets,
			}); err != nil {
				log.Warn().Err(err).Msg("Failed to set attestation summary")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
}
