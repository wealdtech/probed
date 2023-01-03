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

package rest_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	restdaemon "github.com/wealdtech/probed/services/daemon/rest"
	nullmetrics "github.com/wealdtech/probed/services/metrics/null"
	mockprobedb "github.com/wealdtech/probed/services/probedb/mock"
)

func TestService(t *testing.T) {
	ctx := context.Background()
	probeDB := mockprobedb.New()
	monitor := nullmetrics.New()

	tests := []struct {
		name   string
		params []restdaemon.Parameter
		err    string
	}{
		{
			name: "MonitorMissing",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(nil),
				restdaemon.WithServerName("server.wealdtech.com"),
				restdaemon.WithListenAddress(":14734"),
				restdaemon.WithBlockDelaysSetter(probeDB),
				restdaemon.WithHeadDelaysSetter(probeDB),
				restdaemon.WithAggregateAttestationsSetter(probeDB),
				restdaemon.WithAttestationSummariesSetter(probeDB),
			},
			err: "problem with parameters: no monitor specified",
		},
		{
			name: "ServerNameMissing",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(monitor),
				restdaemon.WithListenAddress(":14734"),
				restdaemon.WithBlockDelaysSetter(probeDB),
				restdaemon.WithHeadDelaysSetter(probeDB),
				restdaemon.WithAggregateAttestationsSetter(probeDB),
				restdaemon.WithAttestationSummariesSetter(probeDB),
			},
			err: "problem with parameters: no server name specified",
		},
		{
			name: "ListenAddressMissing",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(monitor),
				restdaemon.WithServerName("server.wealdtech.com"),
				restdaemon.WithBlockDelaysSetter(probeDB),
				restdaemon.WithHeadDelaysSetter(probeDB),
				restdaemon.WithAggregateAttestationsSetter(probeDB),
				restdaemon.WithAttestationSummariesSetter(probeDB),
			},
			err: "problem with parameters: no listen address specified",
		},
		{
			name: "BlockDelaysSetterMissing",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(monitor),
				restdaemon.WithServerName("server.wealdtech.com"),
				restdaemon.WithListenAddress(":14734"),
				restdaemon.WithHeadDelaysSetter(probeDB),
				restdaemon.WithAggregateAttestationsSetter(probeDB),
				restdaemon.WithAttestationSummariesSetter(probeDB),
			},
			err: "problem with parameters: no block delays setter specified",
		},
		{
			name: "HeadDelaysSetterMissing",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(monitor),
				restdaemon.WithServerName("server.wealdtech.com"),
				restdaemon.WithListenAddress(":14734"),
				restdaemon.WithBlockDelaysSetter(probeDB),
				restdaemon.WithAggregateAttestationsSetter(probeDB),
				restdaemon.WithAttestationSummariesSetter(probeDB),
			},
			err: "problem with parameters: no head delays setter specified",
		},
		{
			name: "AggregateAttestationsSetterMissing",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(monitor),
				restdaemon.WithServerName("server.wealdtech.com"),
				restdaemon.WithListenAddress(":14734"),
				restdaemon.WithBlockDelaysSetter(probeDB),
				restdaemon.WithHeadDelaysSetter(probeDB),
				restdaemon.WithAttestationSummariesSetter(probeDB),
			},
			err: "problem with parameters: no aggregate attestations setter specified",
		},
		{
			name: "AttestationSummariesSetterMissing",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(monitor),
				restdaemon.WithServerName("server.wealdtech.com"),
				restdaemon.WithListenAddress(":14734"),
				restdaemon.WithBlockDelaysSetter(probeDB),
				restdaemon.WithHeadDelaysSetter(probeDB),
				restdaemon.WithAggregateAttestationsSetter(probeDB),
			},
			err: "problem with parameters: no attestation summaries setter specified",
		},
		{
			name: "Good",
			params: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithMonitor(monitor),
				restdaemon.WithServerName("server.wealdtech.com"),
				restdaemon.WithListenAddress(":14734"),
				restdaemon.WithBlockDelaysSetter(probeDB),
				restdaemon.WithHeadDelaysSetter(probeDB),
				restdaemon.WithAggregateAttestationsSetter(probeDB),
				restdaemon.WithAttestationSummariesSetter(probeDB),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := restdaemon.New(ctx, test.params...)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
