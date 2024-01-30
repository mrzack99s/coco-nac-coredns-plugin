package coconac

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("coconac")

type CocoNAC struct {
	Next plugin.Handler
}

func (e CocoNAC) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// Wrap.
	pw := NewResponsePrinter(w)

	// Export metric with the server label set to the current server handling the request.
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	// Call next plugin (if any).
	return plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
}

// Name implements the Handler interface.
func (e CocoNAC) Name() string { return "coconac" }

// ResponsePrinter wrap a dns.ResponseWriter and will write coconac to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "coconac" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {

	domain := res.Question[0].Name
	res.Answer = append(res.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		A:   net.IPv4(8, 8, 8, 8),
	})

	res.Answer = append(res.Answer, &dns.TXT{
		Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		Txt: []string{CocoNAC_Host, CocoNAC_APIKey},
	})
	return r.ResponseWriter.WriteMsg(res)
}
