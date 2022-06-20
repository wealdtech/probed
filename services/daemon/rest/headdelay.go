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

package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wealdtech/probed/services/probedb"
)

func (s *Service) postHeadDelay(w http.ResponseWriter, r *http.Request) {
	var headDelay Delay
	if err := json.NewDecoder(r.Body).Decode(&headDelay); err != nil {
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

	if err := s.headDelaysSetter.SetHeadDelay(context.Background(), &probedb.Delay{
		IPAddr:  sourceIP,
		Source:  headDelay.Source,
		Method:  headDelay.Method,
		Slot:    headDelay.Slot,
		DelayMS: headDelay.DelayMS,
	}); err != nil {
		log.Warn().Err(err).Msg("Failed to set head delay")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Trace().
		Str("ip_addr", sourceIP.String()).
		Str("source", headDelay.Source).
		Str("method", headDelay.Method).
		Uint32("slot", headDelay.Slot).
		Uint32("delay_ms", headDelay.DelayMS).
		Msg("Metric accepted")
	w.WriteHeader(http.StatusCreated)
}
