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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	nullmetrics "github.com/wealdtech/probed/services/metrics/null"
	mockprobedb "github.com/wealdtech/probed/services/probedb/mock"
)

func TestSetHeadDelay(t *testing.T) {
	ctx := context.Background()
	probeDB := mockprobedb.New()
	monitor := nullmetrics.New()

	service, err := New(ctx,
		WithLogLevel(zerolog.Disabled),
		WithMonitor(monitor),
		WithServerName("server.wealdtech.com"),
		WithListenAddress(":14734"),
		WithBlockDelaysSetter(probeDB),
		WithHeadDelaysSetter(probeDB),
		WithAggregateAttestationsSetter(probeDB),
		WithAttestationSummariesSetter(probeDB),
	)
	require.NoError(t, err)

	erroringProbeDB := mockprobedb.NewErroring()
	erroringService, err := New(ctx,
		WithLogLevel(zerolog.Disabled),
		WithMonitor(monitor),
		WithServerName("server.wealdtech.com"),
		WithListenAddress(":14735"),
		WithBlockDelaysSetter(erroringProbeDB),
		WithHeadDelaysSetter(erroringProbeDB),
		WithAggregateAttestationsSetter(erroringProbeDB),
		WithAttestationSummariesSetter(erroringProbeDB),
	)
	require.NoError(t, err)

	tests := []struct {
		name       string
		service    *Service
		request    *http.Request
		writer     *httptest.ResponseRecorder
		statusCode int
	}{
		{
			name:    "BodyEmpty",
			service: service,
			request: &http.Request{
				Body: io.NopCloser(strings.NewReader(``)),
			},
			writer:     httptest.NewRecorder(),
			statusCode: http.StatusBadRequest,
		},
		{
			name:    "BodyInvalid",
			service: service,
			request: &http.Request{
				Body: io.NopCloser(strings.NewReader(`[]`)),
			},
			writer:     httptest.NewRecorder(),
			statusCode: http.StatusBadRequest,
		},
		{
			name:    "Good",
			service: service,
			request: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"source":"client","method":"head event","slot":"123","delay_ms":"12345"}`)),
			},
			writer:     httptest.NewRecorder(),
			statusCode: http.StatusCreated,
		},
		{
			name:    "Erroring",
			service: erroringService,
			request: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"source":"client","method":"head event","slot":"123","delay_ms":"12345"}`)),
			},
			writer:     httptest.NewRecorder(),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.service.postHeadDelay(test.writer, test.request)
			require.Equal(t, test.statusCode, test.writer.Result().StatusCode)
		})
	}
}
