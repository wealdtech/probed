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
	"errors"

	"github.com/wealdtech/probed/services/probedb"
)

// ErroringService is a mock that errors.
type ErroringService struct{}

// NewErroring returns a mock probe database that errors.
func NewErroring() *ErroringService {
	return &ErroringService{}
}

// SetBlockDelay sets a block delay.
func (s *ErroringService) SetBlockDelay(ctx context.Context, blockDelay *probedb.BlockDelay) error {
	return errors.New("mock")
}

// BeginTx begins a transaction.
func (s *ErroringService) BeginTx(ctx context.Context) (context.Context, context.CancelFunc, error) {
	return nil, nil, errors.New("mock")
}

// CommitTx commits a transaction.
func (s *ErroringService) CommitTx(ctx context.Context) error {
	return errors.New("mock")
}

// SetMetadata sets a metadata key to a JSON value.
func (s *ErroringService) SetMetadata(ctx context.Context, key string, value []byte) error {
	return errors.New("mock")
}

// Metadata obtains the JSON value from a metadata key.
func (s *ErroringService) Metadata(ctx context.Context, key string) ([]byte, error) {
	return nil, errors.New("mock")
}
