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

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wealdtech/probed/services/metrics"
)

var metricsNamespace = "probed_daemon"

var requests *prometheus.GaugeVec

func registerMetrics(ctx context.Context, monitor metrics.Service) error {
	if requests != nil {
		// Already registered.
		return nil
	}
	if monitor == nil {
		// No monitor.
		return nil
	}
	if monitor.Presenter() == "prometheus" {
		return registerPrometheusMetrics(ctx)
	}
	return nil
}

func registerPrometheusMetrics(ctx context.Context) error {
	requests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "requests_total",
		Help:      "Requests",
	}, []string{"request", "result"})
	if err := prometheus.Register(requests); err != nil {
		return errors.Wrap(err, "failed to register requests_total")
	}

	return nil
}

func requestHandled(request string, result string) {
	if requests != nil {
		requests.WithLabelValues(request, result).Inc()
	}
}
