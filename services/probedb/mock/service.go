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

package mock

import (
	"context"

	"github.com/wealdtech/probed/services/probedb"
)

// Service is a mock.
type Service struct{}

// New returns a mock probe database.
func New() *Service {
	return &Service{}
}

// SetBlockDelay sets a block delay.
func (s *Service) SetBlockDelay(ctx context.Context, delay *probedb.Delay) error {
	return nil
}

// SetHeadDelay sets a head delay.
func (s *Service) SetHeadDelay(ctx context.Context, delay *probedb.Delay) error {
	return nil
}

// BeginTx begins a transaction.
func (s *Service) BeginTx(ctx context.Context) (context.Context, context.CancelFunc, error) {
	return nil, nil, nil
}

// CommitTx commits a transaction.
func (s *Service) CommitTx(ctx context.Context) error {
	return nil
}

// SetMetadata sets a metadata key to a JSON value.
func (s *Service) SetMetadata(ctx context.Context, key string, value []byte) error {
	return nil
}

// Metadata obtains the JSON value from a metadata key.
func (s *Service) Metadata(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
