package bottin

import (
	"context"
	"time"
	"strings"
	"github.com/miekg/dns"
	"github.com/dgraph-io/ristretto/v2"
)

type RR struct {
	Name   string
	Type   string
	Value  string
	TTL    time.Duration
	Expiry time.Time
}

func (rr *RR)Key() string {
	return rr.Name + "|" + rr.Type
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
	cache  *ristretto.Cache[string, []RR]
	client *dns.Client
}

func New(cap int) *BottinResolver {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []RR]{
	NumCounters: 1 << 17,
	MaxCost:     1 << 16,
	BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}
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
