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

// Package util provides utilities for the probe system.
package util

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	majordomo "github.com/wealdtech/go-majordomo"
	"github.com/wealdtech/probed/services/probedb"
	postgresqlprobedb "github.com/wealdtech/probed/services/probedb/postgresql"
)

// InitProbeDB initialises the probe database.
func InitProbeDB(ctx context.Context, majordomo majordomo.Service) (probedb.Service, error) {
	opts := []postgresqlprobedb.Parameter{
		postgresqlprobedb.WithLogLevel(LogLevel("probedb")),
		postgresqlprobedb.WithServer(viper.GetString("probedb.server")),
		postgresqlprobedb.WithUser(viper.GetString("probedb.user")),
		postgresqlprobedb.WithPassword(viper.GetString("probedb.password")),
		postgresqlprobedb.WithPort(viper.GetInt32("probedb.port")),
	}

	if viper.GetString("probedb.client-cert") != "" {
		clientCert, err := majordomo.Fetch(ctx, viper.GetString("probedb.client-cert"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to read client certificate")
		}
		opts = append(opts, postgresqlprobedb.WithClientCert(clientCert))
	}

	if viper.GetString("probedb.client-key") != "" {
		clientKey, err := majordomo.Fetch(ctx, viper.GetString("probedb.client-key"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to read client key")
		}
		opts = append(opts, postgresqlprobedb.WithClientKey(clientKey))
	}

	if viper.GetString("probedb.ca-cert") != "" {
		caCert, err := majordomo.Fetch(ctx, viper.GetString("probedb.ca-cert"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to read certificate authority certificate")
		}
		opts = append(opts, postgresqlprobedb.WithCACert(caCert))
	}

	return postgresqlprobedb.New(ctx, opts...)
}
