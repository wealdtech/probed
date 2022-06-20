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
	"errors"

	"github.com/rs/zerolog"
	"github.com/wealdtech/probed/services/metrics"
	nullmetrics "github.com/wealdtech/probed/services/metrics/null"
	"github.com/wealdtech/probed/services/probedb"
)

type parameters struct {
	logLevel          zerolog.Level
	monitor           metrics.Service
	serverName        string
	listenAddress     string
	blockDelaysSetter probedb.BlockDelaysSetter
	headDelaysSetter  probedb.HeadDelaysSetter
}

// Parameter is the interface for service parameters.
type Parameter interface {
	apply(*parameters)
}

type parameterFunc func(*parameters)

func (f parameterFunc) apply(p *parameters) {
	f(p)
}

// WithLogLevel sets the log level for the module.
func WithLogLevel(logLevel zerolog.Level) Parameter {
	return parameterFunc(func(p *parameters) {
		p.logLevel = logLevel
	})
}

// WithMonitor sets the monitor for the module.
func WithMonitor(monitor metrics.Service) Parameter {
	return parameterFunc(func(p *parameters) {
		p.monitor = monitor
	})
}

// WithServerName sets the server name for this module.
func WithServerName(name string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.serverName = name
	})
}

// WithListenAddress sets the listen address for this module.
func WithListenAddress(listenAddress string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.listenAddress = listenAddress
	})
}

// WithBlockDelaysSetter sets the block delays setter for this module.
func WithBlockDelaysSetter(setter probedb.BlockDelaysSetter) Parameter {
	return parameterFunc(func(p *parameters) {
		p.blockDelaysSetter = setter
	})
}

// WithHeadDelaysSetter sets the head delays setter for this module.
func WithHeadDelaysSetter(setter probedb.HeadDelaysSetter) Parameter {
	return parameterFunc(func(p *parameters) {
		p.headDelaysSetter = setter
	})
}

// parseAndCheckParameters parses and checks parameters to ensure that mandatory parameters are present and correct.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	parameters := parameters{
		logLevel: zerolog.GlobalLevel(),
		monitor:  nullmetrics.New(),
	}
	for _, p := range params {
		if params != nil {
			p.apply(&parameters)
		}
	}

	if parameters.monitor == nil {
		return nil, errors.New("no monitor specified")
	}
	if parameters.serverName == "" {
		return nil, errors.New("no server name specified")
	}
	if parameters.listenAddress == "" {
		return nil, errors.New("no listen address specified")
	}
	if parameters.blockDelaysSetter == nil {
		return nil, errors.New("no block delays setter specified")
	}
	if parameters.headDelaysSetter == nil {
		return nil, errors.New("no head delays setter specified")
	}

	return &parameters, nil
}
