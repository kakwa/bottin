package bottin

import (
	"context"
	"github.com/miekg/dns"
	"strings"
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

// calculateExpiry calculates the expiry time of an RR.
func calculateExpiry(drr dns.RR) (time.Duration, time.Time) {
	ttl := time.Second * time.Duration(drr.Header().Ttl)
	expiry := time.Now().Add(ttl)
	return ttl, expiry
}

func convertRR(drr dns.RR, expire bool) (RR, bool) {
	var ttl time.Duration
	var expiry time.Time
	if expire {
		ttl, expiry = calculateExpiry(drr)
	}
	switch t := drr.(type) {
	case *dns.SOA:
		return RR{toLowerFQDN(t.Hdr.Name), "SOA", toLowerFQDN(t.Ns), ttl, expiry}, true
	case *dns.NS:
		return RR{toLowerFQDN(t.Hdr.Name), "NS", toLowerFQDN(t.Ns), ttl, expiry}, true
	case *dns.CNAME:
		return RR{toLowerFQDN(t.Hdr.Name), "CNAME", toLowerFQDN(t.Target), ttl, expiry}, true
	case *dns.A:
		return RR{toLowerFQDN(t.Hdr.Name), "A", t.A.String(), ttl, expiry}, true
	case *dns.AAAA:
		return RR{toLowerFQDN(t.Hdr.Name), "AAAA", t.AAAA.String(), ttl, expiry}, true
	case *dns.TXT:
		return RR{toLowerFQDN(t.Hdr.Name), "TXT", strings.Join(t.Txt, "\t"), ttl, expiry}, true
	default:
		fields := strings.Fields(drr.String())
		if len(fields) >= 4 {
			return RR{toLowerFQDN(fields[0]), fields[3], strings.Join(fields[4:], "\t"), ttl, expiry}, true
		}
	}
	return RR{}, false
}

func parent(name string) (string, bool) {
	labels := dns.SplitDomainName(name)
	if labels == nil {
		return "", false
	}
	return toLowerFQDN(strings.Join(labels[1:], ".")), true
}

func toLowerFQDN(name string) string {
	return dns.Fqdn(strings.ToLower(name))
}
