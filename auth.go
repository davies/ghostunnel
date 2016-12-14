/*-
 * Copyright 2015 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"crypto/tls"
)

func authorized(conn tls.ConnectionState) bool {
	// First up: check if we have a valid client certificate.
	if !*serverNoClientCert && len(conn.VerifiedChains) == 0 {
		return false
	}

	// If --allow-all has been set, a valid cert is sufficient to connect.
	if *serverAllowAll {
		return true
	}

	cert := conn.VerifiedChains[0][0]

	// Check CN against --allow-cn flag(s).
	for _, expectedCN := range *serverAllowedCNs {
		if cert.Subject.CommonName == expectedCN {
			return true
		}
	}

	// Check OUs against --allow-ou flag(s).
	for _, expectedOU := range *serverAllowedOUs {
		for _, clientOU := range cert.Subject.OrganizationalUnit {
			if clientOU == expectedOU {
				return true
			}
		}
	}

	for _, expectedDNS := range *serverAllowedDNSs {
		for _, clientDNS := range cert.DNSNames {
			if clientDNS == expectedDNS {
				return true
			}
		}
	}

	for _, expectedIP := range *serverAllowedIPs {
		for _, clientIP := range cert.IPAddresses {
			if expectedIP.Equal(clientIP) {
				return true
			}
		}
	}

	return false
}
