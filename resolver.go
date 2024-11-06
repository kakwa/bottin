package bottin

import (
	"context"
	"errors"
	"fmt"
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
	root   *Cache
	cache  *Cache
	client *dns.Client
}

func New(cap int) *BottinResolver {
	root := NewCache()
	cache := NewCache()

	res := BottinResolver{
		cache: cache,
		root:  root,
	}
	res.initRoot()
	return &res
}

func NewExpiring(cap int) *BottinResolver {
	// FIXME
	return New(cap)
}

func NewExpiringWithTimeout(cap int, timeout time.Duration) *BottinResolver {
	// FIXME
	return New(cap)
}

func NewResolver(options ...Option) *BottinResolver {
	// FIXME
	return New(42)
}

func NewWithTimeout(cap int, timeout time.Duration) *BottinResolver {
	// FIXME
	return New(cap)
}

func (br *BottinResolver) Resolve(qname, qtype string) RRs {
	ret, _ := br.ResolveErr(qname, qtype)
	return ret
}

func (br *BottinResolver) ResolveCtx(ctx context.Context, qname, qtype string) (RRs, error) {
	var err error
	var nsRRs RRs
	parentQname, isRoot := parent(qname)
	if isRoot && qtype == "NS" {
		c, _ := br.root.DumpJSON()
		fmt.Printf("%s\n", c)
		nsRRsNS, _ := br.root.Get(".|NS")
		for _, rr := range nsRRsNS {
			nsRRsA, _ := br.root.Get(rr.Value + "|A")
			nsRRs.AnswerRRs = append(nsRRs.AnswerRRs, nsRRsA...)
		}
		// for _, rr := range nsRRsNS {
		// 	nsRRsAAAA, _ := br.root.Get(rr.Value + "|AAAA")
		// 	nsRRs.AnswerRRs = append(nsRRs.AnswerRRs, nsRRsAAAA...)
		// }

		return nsRRs, nil
	} else {
		nsRRs, err = br.ResolveCtx(ctx, parentQname, "NS")
		if err != nil {
			return RRs{}, err
		}
	}

	return br.exchange(ctx, qname, qtype, nsRRs)
}

// perform none recursive query using miekg dns
func (br *BottinResolver) exchange(ctx context.Context, qname, qtype string, nsRRs RRs) (RRs, error) {
	var results RRs
	dnsType, ok := dns.StringToType[qtype]
	if !ok {
		return RRs{}, errors.New("invalid query type")
	}

	for _, nsRR := range nsRRs.AnswerRRs {
		// Create a DNS client and set the timeout
		client := new(dns.Client)
		client.Timeout = 5 * time.Second

		// Create a DNS message for the query
		msg := new(dns.Msg)
		msg.SetQuestion(dns.Fqdn(qname), dnsType)
		msg.RecursionDesired = false // Non-recursive query

		// Get the nameserver address from the RR
		nsAddr := nsRR.Value // Assume `nsRR.Value` contains the IP address of the nameserver

		// Send the query to the nameserver
		fmt.Printf(">>>> trying: %s\n", nsAddr)
		resp, _, err := client.ExchangeContext(ctx, msg, nsAddr+":53")
		if err != nil {
			fmt.Printf("%#v\n", err)
			continue // Try the next nameserver if there's an error
		}

		// Process the response and add records to results
		for _, rr := range resp.Answer {
			crr, _ := convertRR(rr, true)
			results.AnswerRRs = append(results.AnswerRRs, crr)
		}
		if len(results.AnswerRRs) > 0 {
			return results, nil // Return if we got a successful response
		}
	}

	return RRs{}, nil
}

func (br *BottinResolver) ResolveContext(ctx context.Context, qname, qtype string) (RRs, error) {
	return br.ResolveCtx(ctx, qname, qtype)
}

func (br *BottinResolver) ResolveErr(qname, qtype string) (RRs, error) {
	return br.ResolveCtx(context.Background(), qname, qtype)
}
