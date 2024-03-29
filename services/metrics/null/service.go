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

// Package null is a metrics service that discards metrics.
package null

import "github.com/wealdtech/probed/services/metrics"

// Service is a metrics service that drops metrics.
type Service struct{}

// New creates a new metrics service that drops metrics.
func New() metrics.Service {
	return &Service{}
}

// Presenter provides the presenter for this service.
func (s *Service) Presenter() string {
	return "null"
}
