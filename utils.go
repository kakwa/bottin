package bottin

import (
	"github.com/miekg/dns"
	"strings"
	"time"
)

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
