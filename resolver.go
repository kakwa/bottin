package bottin

import (
	"context"
	"time"
)

type RR struct {
	Name   string
	Type   string
	Value  string
	TTL    time.Duration
	Expiry time.Time
}

type RRs struct {
	AnswerRRs     []RR
	AuthorityRRs  []RR
	AdditionalRRs []RR
}

type Resolver interface {
	Resolve(qname, qtype string) RRs
	ResolveContext(ctx context.Context, qname, qtype string) (RRs, error)
	ResolveCtx(ctx context.Context, qname, qtype string) (RRs, error)
	ResolveErr(qname, qtype string) (RRs, error)
}

type BottinResolver struct {
	cache int
}

func New(cap int) *BottinResolver {
	return nil
}

func NewExpiring(cap int) *BottinResolver {
	return nil
}

func NewExpiringWithTimeout(cap int, timeout time.Duration) *BottinResolver {
	return nil
}

func NewResolver(options ...Option) *BottinResolver {
	return nil
}

func NewWithTimeout(cap int, timeout time.Duration) *BottinResolver {
	return nil
}

func (br *BottinResolver) Resolve(qname, qtype string) RRs {
	return RRs{}
}

func (br *BottinResolver) ResolveContext(ctx context.Context, qname, qtype string) (RRs, error) {
	return RRs{}, nil
}

func (br *BottinResolver) ResolveCtx(ctx context.Context, qname, qtype string) (RRs, error) {
	return RRs{}, nil
}

func (br *BottinResolver) ResolveErr(qname, qtype string) (RRs, error) {
	return RRs{}, nil
}
