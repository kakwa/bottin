package bottin

import (
	"context"
	"github.com/miekg/dns"
	"time"
)

// RR represents a DNS-like record.
type RR struct {
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	Value  string        `json:"value"`
	TTL    time.Duration `json:"ttl"`
	Expiry time.Time     `json:"expiry"`
}

func (rr *RR) Key() string {
	return toLowerFQDN(rr.Name) + "|" + rr.Type
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
	cache  *Cache
	client *dns.Client
}

func New(cap int) *BottinResolver {
	cache := NewCache()

	res := BottinResolver{
		cache: cache,
	}
	res.initRoot()
	return &res
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
