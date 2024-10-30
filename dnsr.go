package bottin

/*
Significant portions of this code are derived from:

dnsr - A DNS resolver for Go.

Copyright (c) 2014 nb.io, LLC

Licensed under the MIT License (MIT).
Original source: https://github.com/domainr/dnsr
*/

import (
	"fmt"
	"io"
	"net"
	"time"
)

var (
	Timeout             = 2000 * time.Millisecond
	TypicalResponseTime = 100 * time.Millisecond
	MaxRecursion        = 10
	MaxNameservers      = 2
	MaxIPs              = 2
)

// Resolver errors.
var (
	NXDOMAIN = fmt.Errorf("NXDOMAIN")

	ErrMaxRecursion = fmt.Errorf("maximum recursion depth reached: %d", MaxRecursion)
	ErrMaxIPs       = fmt.Errorf("maximum name server IPs queried: %d", MaxIPs)
	ErrNoARecords   = fmt.Errorf("no A records found for name server")
	ErrNoResponse   = fmt.Errorf("no responses received")
	ErrTimeout      = fmt.Errorf("timeout expired") // TODO: Timeouter interface? e.g. func (e) Timeout() bool { return true }
)

// Option specifies a configuration option for a Resolver.
type Option func(*Resolver)

// DebugLogger will receive writes of DNS resolution traces if not nil.
var DebugLogger io.Writer

// WithCache specifies a cache with capacity cap.
func WithCache(cap int) Option {
	return func(r *Resolver) {
		return
	}
}

// WithDialer sets a custom dialer for the Resolver.
func WithDialer(dialer *net.Dialer) Option {
	return func(r *Resolver) {
	}
}

// WithExpiry sets an expiry duration for cached responses.
func WithExpiry() Option {
	return func(r *Resolver) {
	}
}

func WithTCPRetry() Option {
	return func(r *Resolver) {
	}
}

// WithTimeout sets a timeout for the Resolver's operations.
func WithTimeout(timeout time.Duration) Option {
	return func(r *Resolver) {
	}
}
