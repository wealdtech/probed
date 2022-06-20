// Copyright © 2022 Weald Technology Trading.
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
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// sourceIP fetches the IP address of the request.
func sourceIP(r *http.Request) (net.IP, error) {
	// Attempt to obtain from the X-REAL-IP header.
	// This is a single address that represents the source of the request.
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return netIP, nil
	}

	// Attempt to obtain from the X-FORWARDED-FOR header.
	// This is multiple addresses that represents the path taken by the request.
	// Get IP from X-FORWARDED-FOR header.
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return netIP, nil
		}
	}

	// Fetch from the internals of the request itself.
	if r.RemoteAddr == "" {
		// This suggests localhost.
		return net.ParseIP("127.0.0.1"), nil

	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to split host/port from remote address")
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return netIP, nil
	}

	return nil, errors.New("No valid ip found")
}
