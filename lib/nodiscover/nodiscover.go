// Copyright (c) Solace Corporation
// SPDX-License-Identifier: BUSL-1.1

// Package nodiscover replaces github.com/hashicorp/go-discover in the Solace
// build of Consul. The go-discover module links the SDK of every cloud
// provider, and some of these SDKs are on the NDAA Sec. 1260H list. This
// package supports no provider, so a "provider=..." entry in retry_join,
// retry_join_wan or auto_config.server_addresses causes an error. Plain
// addresses and host names continue to work.
package nodiscover

import (
	"errors"
	"log"
)

// ErrNotSupported is the error that Addrs returns for every input.
var ErrNotSupported = errors.New("cloud auto-join (go-discover) is not supported in this build of Consul")

// Discover has the same methods as the go-discover Discover type that Consul uses.
type Discover struct{}

// New returns a Discover that does not support any provider.
func New() (*Discover, error) { return &Discover{}, nil }

// Addrs returns ErrNotSupported.
func (*Discover) Addrs(string, *log.Logger) ([]string, error) { return nil, ErrNotSupported }

// Names returns an empty list, because no provider is available.
func (*Discover) Names() []string { return nil }
