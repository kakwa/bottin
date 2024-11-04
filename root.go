package bottin

import (
	"strings"

	_ "embed"

	"github.com/miekg/dns"
)

//go:generate curl -O https://www.internic.net/domain/named.root

//go:embed named.root
var root string


func (br *BottinResolver)initRoot() {
	zp := dns.NewZoneParser(strings.NewReader(root), "", "")

	for drr, ok := zp.Next(); ok; drr, ok = zp.Next() {
		rr, ok := convertRR(drr, false)
		if ok {
			br.cache.Set(rr.Key(), []RR{rr}, 99)
		}
	}

	if err := zp.Err(); err != nil {
		panic(err)
	}
}
