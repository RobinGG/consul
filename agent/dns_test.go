package agent

import (
	"net"
	"testing"

	"github.com/hashicorp/consul/agent/structs"
	"github.com/miekg/dns"
)

const (
	configUDPAnswerLimit   = 4
	defaultNumUDPResponses = 3
	testUDPTruncateLimit   = 8

	pctNodesWithIPv6 = 0.5

	// generateNumNodes is the upper bounds for the number of hosts used
	// in testing below.  Generate an arbitrarily large number of hosts.
	generateNumNodes = testUDPTruncateLimit * defaultNumUDPResponses * configUDPAnswerLimit
)

// makeRecursor creates a generic DNS server which always returns
// the provided reply. This is useful for mocking a DNS recursor with
// an expected result.
func makeRecursor(t *testing.T, answer dns.Msg) *dns.Server {
	mux := dns.NewServeMux()
	mux.HandleFunc(".", func(resp dns.ResponseWriter, msg *dns.Msg) {
		answer.SetReply(msg)
		if err := resp.WriteMsg(&answer); err != nil {
			t.Fatalf("err: %s", err)
		}
	})
	up := make(chan struct{})
	server := &dns.Server{
		Addr:              "127.0.0.1:0",
		Net:               "udp",
		Handler:           mux,
		NotifyStartedFunc: func() { close(up) },
	}
	go server.ListenAndServe()
	<-up
	server.Addr = server.PacketConn.LocalAddr().String()
	return server
}

// dnsCNAME returns a DNS CNAME record struct
func dnsCNAME(src, dest string) *dns.CNAME {
	return &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   dns.Fqdn(src),
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
		},
		Target: dns.Fqdn(dest),
	}
}

// dnsA returns a DNS A record struct
func dnsA(src, dest string) *dns.A {
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   dns.Fqdn(src),
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
		},
		A: net.ParseIP(dest),
	}
}

func TestRecursorAddr(t *testing.T) {
	t.Parallel()
	addr, err := recursorAddr("8.8.8.8")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if addr != "8.8.8.8:53" {
		t.Fatalf("bad: %v", addr)
	}
}

func TestDNS_NodeLookup(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// Register node
	args := &structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		TaggedAddresses: map[string]string{
			"wan": "127.0.0.2",
		},
	}

	var out struct{}
	if err := a.RPC("Catalog.Register", args, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	m := new(dns.Msg)
	m.SetQuestion("foo.node.consul.", dns.TypeANY)

	c := new(dns.Client)
	in, _, err := c.Exchange(m, a.DNSAddr())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(in.Answer) != 1 {
		t.Fatalf("Bad: %#v", in)
	}

	aRec, ok := in.Answer[0].(*dns.A)
	if !ok {
		t.Fatalf("Bad: %#v", in.Answer[0])
	}
	if aRec.A.String() != "127.0.0.1" {
		t.Fatalf("Bad: %#v", in.Answer[0])
	}
	if aRec.Hdr.Ttl != 0 {
		t.Fatalf("Bad: %#v", in.Answer[0])
	}

	// Re-do the query, but specify the DC
	m = new(dns.Msg)
	m.SetQuestion("foo.node.dc1.consul.", dns.TypeANY)

	c = new(dns.Client)
	in, _, err = c.Exchange(m, a.DNSAddr())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(in.Answer) != 1 {
		t.Fatalf("Bad: %#v", in)
	}

	aRec, ok = in.Answer[0].(*dns.A)
	if !ok {
		t.Fatalf("Bad: %#v", in.Answer[0])
	}
	if aRec.A.String() != "127.0.0.1" {
		t.Fatalf("Bad: %#v", in.Answer[0])
	}
	if aRec.Hdr.Ttl != 0 {
		t.Fatalf("Bad: %#v", in.Answer[0])
	}

	// lookup a non-existing node, we should receive a SOA
	m = new(dns.Msg)
	m.SetQuestion("nofoo.node.dc1.consul.", dns.TypeANY)

	c = new(dns.Client)
	in, _, err = c.Exchange(m, a.DNSAddr())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(in.Ns) != 1 {
		t.Fatalf("Bad: %#v %#v", in, len(in.Answer))
	}

	soaRec, ok := in.Ns[0].(*dns.SOA)
	if !ok {
		t.Fatalf("Bad: %#v", in.Ns[0])
	}
	if soaRec.Hdr.Ttl != 0 {
		t.Fatalf("Bad: %#v", in.Ns[0])
	}

}

// todo(fs): func TestDNS_CaseInsensitiveNodeLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "Foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("fOO.node.dc1.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("empty lookup: %#v", in)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NodeLookup_PeriodName(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node with period in name
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo.bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("foo.bar.node.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok := in.Answer[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if aRec.A.String() != "127.0.0.1" {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NodeLookup_AAAA(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "::4242:4242",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("bar.node.consul.", dns.TypeAAAA)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok := in.Answer[0].(*dns.AAAA)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if aRec.AAAA.String() != "::4242:4242" {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Ttl != 0 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NodeLookup_CNAME(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	recursor := makeRecursor(t, dns.Msg{
// todo(fs): 		Answer: []dns.RR{
// todo(fs): 			dnsCNAME("www.google.com", "google.com"),
// todo(fs): 			dnsA("google.com", "1.2.3.4"),
// todo(fs): 		},
// todo(fs): 	})
// todo(fs): 	defer recursor.Shutdown()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSRecursors = []string{recursor.Addr}
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "google",
// todo(fs): 		Address:    "www.google.com",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("google.node.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Should have the service record, CNAME record + A record
// todo(fs): 	if len(in.Answer) != 3 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	cnRec, ok := in.Answer[0].(*dns.CNAME)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if cnRec.Target != "www.google.com." {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if cnRec.Hdr.Ttl != 0 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_EDNS0(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.2",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetEdns0(12345, true)
// todo(fs): 	m.SetQuestion("foo.node.dc1.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("empty lookup: %#v", in)
// todo(fs): 	}
// todo(fs): 	edns := in.IsEdns0()
// todo(fs): 	if edns == nil {
// todo(fs): 		t.Fatalf("empty edns: %#v", in)
// todo(fs): 	}
// todo(fs): 	if edns.UDPSize() != 12345 {
// todo(fs): 		t.Fatalf("bad edns size: %d", edns.UDPSize())
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ReverseLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo2",
// todo(fs): 		Address:    "127.0.0.2",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("2.0.0.127.in-addr.arpa.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	ptrRec, ok := in.Answer[0].(*dns.PTR)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if ptrRec.Ptr != "foo2.node.dc1.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", ptrRec)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ReverseLookup_CustomDomain(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSDomain = "custom"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo2",
// todo(fs): 		Address:    "127.0.0.2",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("2.0.0.127.in-addr.arpa.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	ptrRec, ok := in.Answer[0].(*dns.PTR)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if ptrRec.Ptr != "foo2.node.dc1.custom." {
// todo(fs): 		t.Fatalf("Bad: %#v", ptrRec)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ReverseLookup_IPV6(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "::4242:4242",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("2.4.2.4.2.4.2.4.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	ptrRec, ok := in.Answer[0].(*dns.PTR)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if ptrRec.Ptr != "bar.node.dc1.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", ptrRec)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Lookup a non-existing service/query, we should receive an SOA.
// todo(fs): 	questions = []string{
// todo(fs): 		"nodb.service.consul.",
// todo(fs): 		"nope.query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Ns) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		soaRec, ok := in.Ns[0].(*dns.SOA)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs): 		if soaRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookupWithInternalServiceAddress(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeName = "my.test-node"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service.
// todo(fs): 	// The service is using the consul DNS name as service address
// todo(fs): 	// which triggers a lookup loop and a subsequent stack overflow
// todo(fs): 	// crash.
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "db",
// todo(fs): 			Address: "db.service.consul",
// todo(fs): 			Port:    12345,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Looking up the service should not trigger a loop
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("db.service.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	wantAnswer := []dns.RR{
// todo(fs): 		&dns.SRV{
// todo(fs): 			Hdr:      dns.RR_Header{Name: "db.service.consul.", Rrtype: 0x21, Class: 0x1, Rdlength: 0x15},
// todo(fs): 			Priority: 0x1,
// todo(fs): 			Weight:   0x1,
// todo(fs): 			Port:     12345,
// todo(fs): 			Target:   "foo.node.dc1.consul.",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	verify.Values(t, "answer", in.Answer, wantAnswer)
// todo(fs): 	wantExtra := []dns.RR{
// todo(fs): 		&dns.CNAME{
// todo(fs): 			Hdr:    dns.RR_Header{Name: "foo.node.dc1.consul.", Rrtype: 0x5, Class: 0x1, Rdlength: 0x2},
// todo(fs): 			Target: "db.service.consul.",
// todo(fs): 		},
// todo(fs): 		&dns.A{
// todo(fs): 			Hdr: dns.RR_Header{Name: "db.service.consul.", Rrtype: 0x1, Class: 0x1, Rdlength: 0x4},
// todo(fs): 			A:   []byte{0x7f, 0x0, 0x0, 0x1}, // 127.0.0.1
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	verify.Values(t, "extra", in.Extra, wantExtra)
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ExternalServiceLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with an external service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "www.google.com",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		cnameRec, ok := in.Extra[0].(*dns.CNAME)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Target != "www.google.com." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ExternalServiceToConsulCNAMELookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSDomain = "CONSUL."
// todo(fs): 	cfg.NodeName = "test node"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register the initial node with a service
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "web",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "web",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an external service pointing to the 'web' service
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "alias",
// todo(fs): 			Address:    "web.service.consul",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "alias",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly
// todo(fs): 	questions := []string{
// todo(fs): 		"alias.service.consul.",
// todo(fs): 		"alias.service.CoNsUl.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "alias.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Extra) != 2 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		cnameRec, ok := in.Extra[0].(*dns.CNAME)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Name != "alias.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Target != "web.service.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[1].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "web.service.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs):
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NSRecords(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSDomain = "CONSUL."
// todo(fs): 	cfg.NodeName = "server1"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		TaggedAddresses: map[string]string{
// todo(fs): 			"wan": "127.0.0.2",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("something.node.consul.", dns.TypeNS)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	wantAnswer := []dns.RR{
// todo(fs): 		&dns.NS{
// todo(fs): 			Hdr: dns.RR_Header{Name: "consul.", Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 0, Rdlength: 0x13},
// todo(fs): 			Ns:  "server1.node.dc1.consul.",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	verify.Values(t, "answer", in.Answer, wantAnswer)
// todo(fs): 	wantExtra := []dns.RR{
// todo(fs): 		&dns.A{
// todo(fs): 			Hdr: dns.RR_Header{Name: "server1.node.dc1.consul.", Rrtype: dns.TypeA, Class: dns.ClassINET, Rdlength: 0x4, Ttl: 0},
// todo(fs): 			A:   net.ParseIP("127.0.0.1").To4(),
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	verify.Values(t, "extra", in.Extra, wantExtra)
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NSRecords_IPV6(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig(`
// todo(fs): 		domain = "CONSUL."
// todo(fs): 		node_name = "server1"
// todo(fs): 		advertise_addr = "::1"
// todo(fs): 		advertise_addr_wan = "::1"
// todo(fs): 	`)
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		TaggedAddresses: map[string]string{
// todo(fs): 			"wan": "127.0.0.2",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("server1.node.dc1.consul.", dns.TypeNS)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	wantAnswer := []dns.RR{
// todo(fs): 		&dns.NS{
// todo(fs): 			Hdr: dns.RR_Header{Name: "consul.", Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 0, Rdlength: 0x2},
// todo(fs): 			Ns:  "server1.node.dc1.consul.",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	verify.Values(t, "answer", in.Answer, wantAnswer)
// todo(fs): 	wantExtra := []dns.RR{
// todo(fs): 		&dns.AAAA{
// todo(fs): 			Hdr:  dns.RR_Header{Name: "server1.node.dc1.consul.", Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Rdlength: 0x10, Ttl: 0},
// todo(fs): 			AAAA: net.ParseIP("::1"),
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	verify.Values(t, "extra", in.Extra, wantExtra)
// todo(fs):
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ExternalServiceToConsulCNAMENestedLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeName = "test-node"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register the initial node with a service
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "web",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "web",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an external service pointing to the 'web' service
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "alias",
// todo(fs): 			Address:    "web.service.consul",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "alias",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an external service pointing to the 'alias' service
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "alias2",
// todo(fs): 			Address:    "alias.service.consul",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "alias2",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly
// todo(fs): 	questions := []string{
// todo(fs): 		"alias2.service.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "alias2.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Extra) != 3 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		cnameRec, ok := in.Extra[0].(*dns.CNAME)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Name != "alias2.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Target != "alias.service.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		cnameRec, ok = in.Extra[1].(*dns.CNAME)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Name != "alias.service.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Target != "web.service.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[2].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[2])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "web.service.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[1])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[2])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[2])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_ServiceAddress_A(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Address: "127.0.0.2",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "7f000002.addr.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "7f000002.addr.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.2" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_ServiceAddress_CNAME(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service whose address isn't an IP.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Address: "www.google.com",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		cnameRec, ok := in.Extra[0].(*dns.CNAME)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Target != "www.google.com." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if cnameRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_ServiceAddressIPV6(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Address: "2607:20:4005:808::200e",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "2607002040050808000000000000200e.addr.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.AAAA)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "2607002040050808000000000000200e.addr.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.AAAA.String() != "2607:20:4005:808::200e" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_WanAddress(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg1 := TestConfig()
// todo(fs): 	cfg1.Datacenter = "dc1"
// todo(fs): 	cfg1.TranslateWANAddrs = true
// todo(fs): 	cfg1.ACLDatacenter = ""
// todo(fs): 	a1 := NewTestAgent(t.Name(), cfg1)
// todo(fs): 	defer a1.Shutdown()
// todo(fs):
// todo(fs): 	cfg2 := TestConfig()
// todo(fs): 	cfg2.Datacenter = "dc2"
// todo(fs): 	cfg2.TranslateWANAddrs = true
// todo(fs): 	cfg2.ACLDatacenter = ""
// todo(fs): 	a2 := NewTestAgent(t.Name(), cfg2)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	// Join WAN cluster
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a1.Config.Ports.SerfWan)
// todo(fs): 	if _, err := a2.JoinWAN([]string{addr}); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if got, want := len(a1.WANMembers()), 2; got < want {
// todo(fs): 			r.Fatalf("got %d WAN members want at least %d", got, want)
// todo(fs): 		}
// todo(fs): 		if got, want := len(a2.WANMembers()), 2; got < want {
// todo(fs): 			r.Fatalf("got %d WAN members want at least %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Register a remote node with a service. This is in a retry since we
// todo(fs): 	// need the datacenter to have a route which takes a little more time
// todo(fs): 	// beyond the join, and we don't have direct access to the router here.
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc2",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			TaggedAddresses: map[string]string{
// todo(fs): 				"wan": "127.0.0.2",
// todo(fs): 			},
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a2.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc2",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a2.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the SRV record via service and prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.dc2.consul.",
// todo(fs): 		id + ".query.dc2.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		addr, _ := a1.Config.ClientListener("", a1.Config.Ports.DNS)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "7f000002.addr.dc2.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.2" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Also check the A record directly
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeA)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		addr, _ := a1.Config.ClientListener("", a1.Config.Ports.DNS)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Answer[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != question {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.2" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Now query from the same DC and make sure we get the local address
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		addr, _ := a2.Config.ClientListener("", a2.Config.Ports.DNS)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "foo.node.dc2.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Also check the A record directly from DC2
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeA)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		addr, _ := a2.Config.ClientListener("", a2.Config.Ports.DNS)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Answer[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != question {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_CaseInsensitiveServiceLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "Db",
// todo(fs): 				Tags:    []string{"Master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query, as well as a name.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "somequery",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try some variations to make sure case doesn't matter.
// todo(fs): 	questions := []string{
// todo(fs): 		"master.db.service.consul.",
// todo(fs): 		"mASTER.dB.service.consul.",
// todo(fs): 		"MASTER.dB.service.consul.",
// todo(fs): 		"db.service.consul.",
// todo(fs): 		"DB.service.consul.",
// todo(fs): 		"Db.service.consul.",
// todo(fs): 		"somequery.query.consul.",
// todo(fs): 		"SomeQuery.query.consul.",
// todo(fs): 		"SOMEQUERY.query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("empty lookup: %#v", in)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_TagPeriod(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "db",
// todo(fs): 			Tags:    []string{"v1.master"},
// todo(fs): 			Port:    12345,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("v1.master.db.service.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srvRec.Port != 12345 {
// todo(fs): 		t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 	}
// todo(fs): 	if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.A.String() != "127.0.0.1" {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_PreparedQueryNamePeriod(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register a prepared query with a period in the name.
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "some.query.we.like",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var id string
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("some.query.we.like.query.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srvRec.Port != 12345 {
// todo(fs): 		t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 	}
// todo(fs): 	if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.A.String() != "127.0.0.1" {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_Dedup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a single node with multiple instances of a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				ID:      "db2",
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"slave"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				ID:      "db3",
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"slave"},
// todo(fs): 				Port:    12346,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query, make sure only
// todo(fs): 	// one IP is returned.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeANY)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Answer[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_Dedup_SRV(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a single node with multiple instances of a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				ID:      "db2",
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"slave"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				ID:      "db3",
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"slave"},
// todo(fs): 				Port:    12346,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query, make sure only
// todo(fs): 	// one IP is returned and two unique ports are returned.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 2 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 && srvRec.Port != 12346 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok = in.Answer[1].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[1])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12346 && srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Port == in.Answer[0].(*dns.SRV).Port {
// todo(fs): 			t.Fatalf("should be a different port")
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_Recurse(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	recursor := makeRecursor(t, dns.Msg{
// todo(fs): 		Answer: []dns.RR{dnsA("apple.com", "1.2.3.4")},
// todo(fs): 	})
// todo(fs): 	defer recursor.Shutdown()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSRecursor = recursor.Addr
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("apple.com.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) == 0 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs): 	if in.Rcode != dns.RcodeSuccess {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_Recurse_Truncation(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs):
// todo(fs): 	recursor := makeRecursor(t, dns.Msg{
// todo(fs): 		MsgHdr: dns.MsgHdr{Truncated: true},
// todo(fs): 		Answer: []dns.RR{dnsA("apple.com", "1.2.3.4")},
// todo(fs): 	})
// todo(fs): 	defer recursor.Shutdown()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSRecursor = recursor.Addr
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("apple.com.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != dns.ErrTruncated {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if in.Truncated != true {
// todo(fs): 		t.Fatalf("err: message should have been truncated %v", in)
// todo(fs): 	}
// todo(fs): 	if len(in.Answer) == 0 {
// todo(fs): 		t.Fatalf("Bad: Truncated message ignored, expected some reply %#v", in)
// todo(fs): 	}
// todo(fs): 	if in.Rcode != dns.RcodeSuccess {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_RecursorTimeout(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	serverClientTimeout := 3 * time.Second
// todo(fs): 	testClientTimeout := serverClientTimeout + 5*time.Second
// todo(fs):
// todo(fs): 	resolverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Error(err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	resolver, err := net.ListenUDP("udp", resolverAddr)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Error(err)
// todo(fs): 	}
// todo(fs): 	defer resolver.Close()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSRecursor = resolver.LocalAddr().String() // host must cause a connection|read|write timeout
// todo(fs): 	cfg.DNSConfig.RecursorTimeout = serverClientTimeout
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("apple.com.", dns.TypeANY)
// todo(fs):
// todo(fs): 	// This client calling the server under test must have a longer timeout than the one we set internally
// todo(fs): 	c := &dns.Client{Timeout: testClientTimeout}
// todo(fs):
// todo(fs): 	start := time.Now()
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs):
// todo(fs): 	duration := time.Now().Sub(start)
// todo(fs):
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 0 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs): 	if in.Rcode != dns.RcodeServerFailure {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if duration < serverClientTimeout {
// todo(fs): 		t.Fatalf("Expected the call to return after at least %f seconds but lasted only %f", serverClientTimeout.Seconds(), duration.Seconds())
// todo(fs): 	}
// todo(fs):
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_FilterCritical(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register nodes with health checks in various states.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID: "serf",
// todo(fs): 				Name:    "serf",
// todo(fs): 				Status:  api.HealthCritical,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args2 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "bar",
// todo(fs): 			Address:    "127.0.0.2",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID: "serf",
// todo(fs): 				Name:    "serf",
// todo(fs): 				Status:  api.HealthCritical,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args2, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args3 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "bar",
// todo(fs): 			Address:    "127.0.0.2",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID:   "db",
// todo(fs): 				Name:      "db",
// todo(fs): 				ServiceID: "db",
// todo(fs): 				Status:    api.HealthCritical,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args3, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args4 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "baz",
// todo(fs): 			Address:    "127.0.0.3",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args4, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args5 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "quux",
// todo(fs): 			Address:    "127.0.0.4",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID:   "db",
// todo(fs): 				Name:      "db",
// todo(fs): 				ServiceID: "db",
// todo(fs): 				Status:    api.HealthWarning,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args5, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeANY)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Only 4 and 5 are not failing, so we should get 2 answers
// todo(fs): 		if len(in.Answer) != 2 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		ips := make(map[string]bool)
// todo(fs): 		for _, resp := range in.Answer {
// todo(fs): 			aRec := resp.(*dns.A)
// todo(fs): 			ips[aRec.A.String()] = true
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if !ips["127.0.0.3"] {
// todo(fs): 			t.Fatalf("Bad: %#v should contain 127.0.0.3 (state healthy)", in)
// todo(fs): 		}
// todo(fs): 		if !ips["127.0.0.4"] {
// todo(fs): 			t.Fatalf("Bad: %#v should contain 127.0.0.4 (state warning)", in)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_OnlyFailing(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register nodes with all health checks in a critical state.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID: "serf",
// todo(fs): 				Name:    "serf",
// todo(fs): 				Status:  api.HealthCritical,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args2 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "bar",
// todo(fs): 			Address:    "127.0.0.2",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID: "serf",
// todo(fs): 				Name:    "serf",
// todo(fs): 				Status:  api.HealthCritical,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args2, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args3 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "bar",
// todo(fs): 			Address:    "127.0.0.2",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID:   "db",
// todo(fs): 				Name:      "db",
// todo(fs): 				ServiceID: "db",
// todo(fs): 				Status:    api.HealthCritical,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args3, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeANY)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All 3 are failing, so we should get 0 answers and an NXDOMAIN response
// todo(fs): 		if len(in.Answer) != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if in.Rcode != dns.RcodeNameError {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_OnlyPassing(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSConfig.OnlyPassing = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register nodes with health checks in various states.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID:   "db",
// todo(fs): 				Name:      "db",
// todo(fs): 				ServiceID: "db",
// todo(fs): 				Status:    api.HealthPassing,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args2 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "bar",
// todo(fs): 			Address:    "127.0.0.2",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID:   "db",
// todo(fs): 				Name:      "db",
// todo(fs): 				ServiceID: "db",
// todo(fs): 				Status:    api.HealthWarning,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RPC("Catalog.Register", args2, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args3 := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "baz",
// todo(fs): 			Address:    "127.0.0.3",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 			Check: &structs.HealthCheck{
// todo(fs): 				CheckID:   "db",
// todo(fs): 				Name:      "db",
// todo(fs): 				ServiceID: "db",
// todo(fs): 				Status:    api.HealthCritical,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RPC("Catalog.Register", args3, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service:     "db",
// todo(fs): 					OnlyPassing: true,
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeANY)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Only 1 is passing, so we should only get 1 answer
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		resp := in.Answer[0]
// todo(fs): 		aRec := resp.(*dns.A)
// todo(fs):
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_Randomize(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a large number of nodes.
// todo(fs): 	for i := 0; i < generateNumNodes; i++ {
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       fmt.Sprintf("foo%d", i),
// todo(fs): 			Address:    fmt.Sprintf("127.0.0.%d", i+1),
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "web",
// todo(fs): 				Port:    8000,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "web",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query. Ensure the
// todo(fs): 	// response is randomized each time.
// todo(fs): 	questions := []string{
// todo(fs): 		"web.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		uniques := map[string]struct{}{}
// todo(fs): 		for i := 0; i < 10; i++ {
// todo(fs): 			m := new(dns.Msg)
// todo(fs): 			m.SetQuestion(question, dns.TypeANY)
// todo(fs):
// todo(fs): 			c := &dns.Client{Net: "udp"}
// todo(fs): 			in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			// Response length should be truncated and we should get
// todo(fs): 			// an A record for each response.
// todo(fs): 			if len(in.Answer) != defaultNumUDPResponses {
// todo(fs): 				t.Fatalf("Bad: %#v", len(in.Answer))
// todo(fs): 			}
// todo(fs):
// todo(fs): 			// Collect all the names.
// todo(fs): 			var names []string
// todo(fs): 			for _, rec := range in.Answer {
// todo(fs): 				switch v := rec.(type) {
// todo(fs): 				case *dns.SRV:
// todo(fs): 					names = append(names, v.Target)
// todo(fs): 				case *dns.A:
// todo(fs): 					names = append(names, v.A.String())
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 			nameS := strings.Join(names, "|")
// todo(fs):
// todo(fs): 			// Tally the results.
// todo(fs): 			uniques[nameS] = struct{}{}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Give some wiggle room. Since the responses are randomized and
// todo(fs): 		// there is a finite number of combinations, requiring 0
// todo(fs): 		// duplicates every test run eventually gives us failures.
// todo(fs): 		if len(uniques) < 2 {
// todo(fs): 			t.Fatalf("unique response ratio too low: %d/10\n%v", len(uniques), uniques)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_Truncate(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSConfig.EnableTruncate = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a large number of nodes.
// todo(fs): 	for i := 0; i < generateNumNodes; i++ {
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       fmt.Sprintf("foo%d", i),
// todo(fs): 			Address:    fmt.Sprintf("127.0.0.%d", i+1),
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "web",
// todo(fs): 				Port:    8000,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "web",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query. Ensure the
// todo(fs): 	// response is truncated each time.
// todo(fs): 	questions := []string{
// todo(fs): 		"web.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeANY)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil && err != dns.ErrTruncated {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Check for the truncate bit
// todo(fs): 		if !in.Truncated {
// todo(fs): 			t.Fatalf("should have truncate bit")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_LargeResponses(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSConfig.EnableTruncate = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	longServiceName := "this-is-a-very-very-very-very-very-long-name-for-a-service"
// todo(fs):
// todo(fs): 	// Register a lot of nodes.
// todo(fs): 	for i := 0; i < 4; i++ {
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       fmt.Sprintf("foo%d", i),
// todo(fs): 			Address:    fmt.Sprintf("127.0.0.%d", i+1),
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: longServiceName,
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: longServiceName,
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: longServiceName,
// todo(fs): 					Tags:    []string{"master"},
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var id string
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"_" + longServiceName + "._master.service.consul.",
// todo(fs): 		longServiceName + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil && err != dns.ErrTruncated {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Make sure the response size is RFC 1035-compliant for UDP messages
// todo(fs): 		if in.Len() > 512 {
// todo(fs): 			t.Fatalf("Bad: %d", in.Len())
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should only have two answers now
// todo(fs): 		if len(in.Answer) != 2 {
// todo(fs): 			t.Fatalf("Bad: %d", len(in.Answer))
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Make sure the ADDITIONAL section matches the ANSWER section.
// todo(fs): 		if len(in.Answer) != len(in.Extra) {
// todo(fs): 			t.Fatalf("Bad: %d vs. %d", len(in.Answer), len(in.Extra))
// todo(fs): 		}
// todo(fs): 		for i := 0; i < len(in.Answer); i++ {
// todo(fs): 			srv, ok := in.Answer[i].(*dns.SRV)
// todo(fs): 			if !ok {
// todo(fs): 				t.Fatalf("Bad: %#v", in.Answer[i])
// todo(fs): 			}
// todo(fs):
// todo(fs): 			a, ok := in.Extra[i].(*dns.A)
// todo(fs): 			if !ok {
// todo(fs): 				t.Fatalf("Bad: %#v", in.Extra[i])
// todo(fs): 			}
// todo(fs):
// todo(fs): 			if srv.Target != a.Hdr.Name {
// todo(fs): 				t.Fatalf("Bad: %#v %#v", srv, a)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Check for the truncate bit
// todo(fs): 		if !in.Truncated {
// todo(fs): 			t.Fatalf("should have truncate bit")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func testDNS_ServiceLookup_responseLimits(t *testing.T, answerLimit int, qType uint16,
// todo(fs): 	expectedService, expectedQuery, expectedQueryID int) (bool, error) {
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSConfig.UDPAnswerLimit = answerLimit
// todo(fs): 	cfg.NodeName = "test-node"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	for i := 0; i < generateNumNodes; i++ {
// todo(fs): 		nodeAddress := fmt.Sprintf("127.0.0.%d", i+1)
// todo(fs): 		if rand.Float64() < pctNodesWithIPv6 {
// todo(fs): 			nodeAddress = fmt.Sprintf("fe80::%d", i+1)
// todo(fs): 		}
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       fmt.Sprintf("foo%d", i),
// todo(fs): 			Address:    nodeAddress,
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "api-tier",
// todo(fs): 				Port:    8080,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			return false, fmt.Errorf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "api-tier",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "api-tier",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			return false, fmt.Errorf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"api-tier.service.consul.",
// todo(fs): 		"api-tier.query.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for idx, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, qType)
// todo(fs):
// todo(fs): 		c := &dns.Client{Net: "udp"}
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			return false, fmt.Errorf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		switch idx {
// todo(fs): 		case 0:
// todo(fs): 			if (expectedService > 0 && len(in.Answer) != expectedService) ||
// todo(fs): 				(expectedService < -1 && len(in.Answer) < lib.AbsInt(expectedService)) {
// todo(fs): 				return false, fmt.Errorf("%d/%d answers received for type %v for %s", len(in.Answer), answerLimit, qType, question)
// todo(fs): 			}
// todo(fs): 		case 1:
// todo(fs): 			if (expectedQuery > 0 && len(in.Answer) != expectedQuery) ||
// todo(fs): 				(expectedQuery < -1 && len(in.Answer) < lib.AbsInt(expectedQuery)) {
// todo(fs): 				return false, fmt.Errorf("%d/%d answers received for type %v for %s", len(in.Answer), answerLimit, qType, question)
// todo(fs): 			}
// todo(fs): 		case 2:
// todo(fs): 			if (expectedQueryID > 0 && len(in.Answer) != expectedQueryID) ||
// todo(fs): 				(expectedQueryID < -1 && len(in.Answer) < lib.AbsInt(expectedQueryID)) {
// todo(fs): 				return false, fmt.Errorf("%d/%d answers received for type %v for %s", len(in.Answer), answerLimit, qType, question)
// todo(fs): 			}
// todo(fs): 		default:
// todo(fs): 			panic("abort")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	return true, nil
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_AnswerLimits(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	// Build a matrix of config parameters (udpAnswerLimit), and the
// todo(fs): 	// length of the response per query type and question.  Negative
// todo(fs): 	// values imply the test must return at least the abs(value) number
// todo(fs): 	// of records in the answer section.  This is required because, for
// todo(fs): 	// example, on OS-X and Linux, the number of answers returned in a
// todo(fs): 	// 512B response is different even though both platforms are x86_64
// todo(fs): 	// and using the same version of Go.
// todo(fs): 	//
// todo(fs): 	// TODO(sean@): Why is it not identical everywhere when using the
// todo(fs): 	// same compiler?
// todo(fs): 	tests := []struct {
// todo(fs): 		name                string
// todo(fs): 		udpAnswerLimit      int
// todo(fs): 		expectedAService    int
// todo(fs): 		expectedAQuery      int
// todo(fs): 		expectedAQueryID    int
// todo(fs): 		expectedAAAAService int
// todo(fs): 		expectedAAAAQuery   int
// todo(fs): 		expectedAAAAQueryID int
// todo(fs): 		expectedANYService  int
// todo(fs): 		expectedANYQuery    int
// todo(fs): 		expectedANYQueryID  int
// todo(fs): 	}{
// todo(fs): 		{"0", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
// todo(fs): 		{"1", 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
// todo(fs): 		{"2", 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
// todo(fs): 		{"3", 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
// todo(fs): 		{"4", 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
// todo(fs): 		{"5", 5, 5, 5, 5, 5, 5, 5, 5, 5, 5},
// todo(fs): 		{"6", 6, 6, 6, 6, 6, 6, 5, 6, 6, -5},
// todo(fs): 		{"7", 7, 7, 7, 6, 7, 7, 5, 7, 7, -5},
// todo(fs): 		{"8", 8, 8, 8, 6, 8, 8, 5, 8, 8, -5},
// todo(fs): 		{"9", 9, 8, 8, 6, 8, 8, 5, 8, 8, -5},
// todo(fs): 		{"20", 20, 8, 8, 6, 8, 8, 5, 8, -5, -5},
// todo(fs): 		{"30", 30, 8, 8, 6, 8, 8, 5, 8, -5, -5},
// todo(fs): 	}
// todo(fs): 	for _, test := range tests {
// todo(fs): 		test := test // capture loop var
// todo(fs): 		t.Run("A lookup", func(t *testing.T) {
// todo(fs): 			t.Parallel()
// todo(fs): 			ok, err := testDNS_ServiceLookup_responseLimits(t, test.udpAnswerLimit, dns.TypeA, test.expectedAService, test.expectedAQuery, test.expectedAQueryID)
// todo(fs): 			if !ok {
// todo(fs): 				t.Errorf("Expected service A lookup %s to pass: %v", test.name, err)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs):
// todo(fs): 		t.Run("AAAA lookup", func(t *testing.T) {
// todo(fs): 			t.Parallel()
// todo(fs): 			ok, err := testDNS_ServiceLookup_responseLimits(t, test.udpAnswerLimit, dns.TypeAAAA, test.expectedAAAAService, test.expectedAAAAQuery, test.expectedAAAAQueryID)
// todo(fs): 			if !ok {
// todo(fs): 				t.Errorf("Expected service AAAA lookup %s to pass: %v", test.name, err)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs):
// todo(fs): 		t.Run("ANY lookup", func(t *testing.T) {
// todo(fs): 			t.Parallel()
// todo(fs): 			ok, err := testDNS_ServiceLookup_responseLimits(t, test.udpAnswerLimit, dns.TypeANY, test.expectedANYService, test.expectedANYQuery, test.expectedANYQueryID)
// todo(fs): 			if !ok {
// todo(fs): 				t.Errorf("Expected service ANY lookup %s to pass: %v", test.name, err)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_CNAME(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	recursor := makeRecursor(t, dns.Msg{
// todo(fs): 		Answer: []dns.RR{
// todo(fs): 			dnsCNAME("www.google.com", "google.com"),
// todo(fs): 			dnsA("google.com", "1.2.3.4"),
// todo(fs): 		},
// todo(fs): 	})
// todo(fs): 	defer recursor.Shutdown()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSRecursor = recursor.Addr
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a name for an address.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "google",
// todo(fs): 			Address:    "www.google.com",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "search",
// todo(fs): 				Port:    80,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "search",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"search.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeANY)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Service CNAME, google CNAME, google A record
// todo(fs): 		if len(in.Answer) != 3 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Should have service CNAME
// todo(fs): 		cnRec, ok := in.Answer[0].(*dns.CNAME)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if cnRec.Target != "www.google.com." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Should have google CNAME
// todo(fs): 		cnRec, ok = in.Answer[1].(*dns.CNAME)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[1])
// todo(fs): 		}
// todo(fs): 		if cnRec.Target != "google.com." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[1])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Check we recursively resolve
// todo(fs): 		if _, ok := in.Answer[2].(*dns.A); !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[2])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NodeLookup_TTL(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	recursor := makeRecursor(t, dns.Msg{
// todo(fs): 		Answer: []dns.RR{
// todo(fs): 			dnsCNAME("www.google.com", "google.com"),
// todo(fs): 			dnsA("google.com", "1.2.3.4"),
// todo(fs): 		},
// todo(fs): 	})
// todo(fs): 	defer recursor.Shutdown()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSRecursor = recursor.Addr
// todo(fs): 	cfg.DNSConfig.NodeTTL = 10 * time.Second
// todo(fs): 	cfg.DNSConfig.AllowStale = Bool(true)
// todo(fs): 	cfg.DNSConfig.MaxStale = time.Second
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("foo.node.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok := in.Answer[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if aRec.A.String() != "127.0.0.1" {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Ttl != 10 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register node with IPv6
// todo(fs): 	args = &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "::4242:4242",
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check an IPv6 record
// todo(fs): 	m = new(dns.Msg)
// todo(fs): 	m.SetQuestion("bar.node.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	in, _, err = c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aaaaRec, ok := in.Answer[0].(*dns.AAAA)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if aaaaRec.AAAA.String() != "::4242:4242" {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if aaaaRec.Hdr.Ttl != 10 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register node with CNAME
// todo(fs): 	args = &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "google",
// todo(fs): 		Address:    "www.google.com",
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m = new(dns.Msg)
// todo(fs): 	m.SetQuestion("google.node.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	in, _, err = c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Should have the CNAME record + a few A records
// todo(fs): 	if len(in.Answer) < 2 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	cnRec, ok := in.Answer[0].(*dns.CNAME)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if cnRec.Target != "www.google.com." {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if cnRec.Hdr.Ttl != 10 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_TTL(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSConfig.ServiceTTL = map[string]time.Duration{
// todo(fs): 		"db": 10 * time.Second,
// todo(fs): 		"*":  5 * time.Second,
// todo(fs): 	}
// todo(fs): 	cfg.DNSConfig.AllowStale = Bool(true)
// todo(fs): 	cfg.DNSConfig.MaxStale = time.Second
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node with 2 services
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "db",
// todo(fs): 			Tags:    []string{"master"},
// todo(fs): 			Port:    12345,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args = &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "api",
// todo(fs): 			Port:    2222,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("db.service.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srvRec.Hdr.Ttl != 10 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Ttl != 10 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m = new(dns.Msg)
// todo(fs): 	m.SetQuestion("api.service.consul.", dns.TypeSRV)
// todo(fs): 	in, _, err = c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	srvRec, ok = in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srvRec.Hdr.Ttl != 5 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok = in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Ttl != 5 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_PreparedQuery_TTL(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSConfig.ServiceTTL = map[string]time.Duration{
// todo(fs): 		"db": 10 * time.Second,
// todo(fs): 		"*":  5 * time.Second,
// todo(fs): 	}
// todo(fs): 	cfg.DNSConfig.AllowStale = Bool(true)
// todo(fs): 	cfg.DNSConfig.MaxStale = time.Second
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node and a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "api",
// todo(fs): 				Port:    2222,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register prepared queries with and without a TTL set for "db", as
// todo(fs): 	// well as one for "api".
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "db-ttl",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 				DNS: structs.QueryDNSOptions{
// todo(fs): 					TTL: "18s",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var id string
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "db-nottl",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "api-nottl",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "api",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure the TTL is set when requested, and overrides the agent-
// todo(fs): 	// specific config since the query takes precedence.
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("db-ttl.query.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srvRec.Hdr.Ttl != 18 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Ttl != 18 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// And the TTL should take the service-specific value from the agent's
// todo(fs): 	// config otherwise.
// todo(fs): 	m = new(dns.Msg)
// todo(fs): 	m.SetQuestion("db-nottl.query.consul.", dns.TypeSRV)
// todo(fs): 	in, _, err = c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	srvRec, ok = in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srvRec.Hdr.Ttl != 10 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok = in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Ttl != 10 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// If there's no query TTL and no service-specific value then the wild
// todo(fs): 	// card value should be used.
// todo(fs): 	m = new(dns.Msg)
// todo(fs): 	m.SetQuestion("api-nottl.query.consul.", dns.TypeSRV)
// todo(fs): 	in, _, err = c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	srvRec, ok = in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srvRec.Hdr.Ttl != 5 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	aRec, ok = in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if aRec.Hdr.Ttl != 5 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_PreparedQuery_Failover(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg1 := TestConfig()
// todo(fs): 	cfg1.Datacenter = "dc1"
// todo(fs): 	cfg1.TranslateWANAddrs = true
// todo(fs): 	cfg1.ACLDatacenter = ""
// todo(fs): 	a1 := NewTestAgent(t.Name(), cfg1)
// todo(fs): 	defer a1.Shutdown()
// todo(fs):
// todo(fs): 	cfg2 := TestConfig()
// todo(fs): 	cfg2.Datacenter = "dc2"
// todo(fs): 	cfg2.TranslateWANAddrs = true
// todo(fs): 	cfg2.ACLDatacenter = ""
// todo(fs): 	a2 := NewTestAgent(t.Name(), cfg2)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	// Join WAN cluster.
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a1.Config.Ports.SerfWan)
// todo(fs): 	if _, err := a2.JoinWAN([]string{addr}); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if got, want := len(a1.WANMembers()), 2; got < want {
// todo(fs): 			r.Fatalf("got %d WAN members want at least %d", got, want)
// todo(fs): 		}
// todo(fs): 		if got, want := len(a2.WANMembers()), 2; got < want {
// todo(fs): 			r.Fatalf("got %d WAN members want at least %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Register a remote node with a service. This is in a retry since we
// todo(fs): 	// need the datacenter to have a route which takes a little more time
// todo(fs): 	// beyond the join, and we don't have direct access to the router here.
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc2",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			TaggedAddresses: map[string]string{
// todo(fs): 				"wan": "127.0.0.2",
// todo(fs): 			},
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a2.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Register a local prepared query.
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "my-query",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 					Failover: structs.QueryDatacenterOptions{
// todo(fs): 						Datacenters: []string{"dc2"},
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var id string
// todo(fs): 		if err := a1.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the SRV record via the query.
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("my-query.query.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	cl_addr, _ := a1.Config.ClientListener("", a1.Config.Ports.DNS)
// todo(fs): 	in, _, err := c.Exchange(m, cl_a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure we see the remote DC and that the address gets
// todo(fs): 	// translated.
// todo(fs): 	if len(in.Answer) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v", in)
// todo(fs): 	}
// todo(fs): 	if in.Answer[0].Header().Name != "my-query.query.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	srv, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs): 	if srv.Target != "7f000002.addr.dc2.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 	}
// todo(fs):
// todo(fs): 	a, ok := in.Extra[0].(*dns.A)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if a.Hdr.Name != "7f000002.addr.dc2.consul." {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): 	if a.A.String() != "127.0.0.2" {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_SRV_RFC(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "db",
// todo(fs): 			Tags:    []string{"master"},
// todo(fs): 			Port:    12345,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	questions := []string{
// todo(fs): 		"_db._master.service.dc1.consul.",
// todo(fs): 		"_db._master.service.consul.",
// todo(fs): 		"_db._master.dc1.consul.",
// todo(fs): 		"_db._master.consul.",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_SRV_RFC_TCP_Default(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "db",
// todo(fs): 			Tags:    []string{"master"},
// todo(fs): 			Port:    12345,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	questions := []string{
// todo(fs): 		"_db._tcp.service.dc1.consul.",
// todo(fs): 		"_db._tcp.service.consul.",
// todo(fs): 		"_db._tcp.dc1.consul.",
// todo(fs): 		"_db._tcp.consul.",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		srvRec, ok := in.Answer[0].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if srvRec.Port != 12345 {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Target != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", srvRec)
// todo(fs): 		}
// todo(fs): 		if srvRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Extra[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Name != "foo.node.dc1.consul." {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.String() != "127.0.0.1" {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Extra[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_ServiceLookup_FilterACL(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	tests := []struct {
// todo(fs): 		token   string
// todo(fs): 		results int
// todo(fs): 	}{
// todo(fs): 		{"root", 1},
// todo(fs): 		{"anonymous", 0},
// todo(fs): 	}
// todo(fs): 	for _, tt := range tests {
// todo(fs): 		t.Run("ACLToken == "+tt.token, func(t *testing.T) {
// todo(fs): 			cfg := TestConfig()
// todo(fs): 			cfg.ACLToken = tt.token
// todo(fs): 			cfg.ACLMasterToken = "root"
// todo(fs): 			cfg.ACLDatacenter = "dc1"
// todo(fs): 			cfg.ACLDownPolicy = "deny"
// todo(fs): 			cfg.ACLDefaultPolicy = "deny"
// todo(fs): 			a := NewTestAgent(t.Name(), cfg)
// todo(fs): 			defer a.Shutdown()
// todo(fs):
// todo(fs): 			// Register a service
// todo(fs): 			args := &structs.RegisterRequest{
// todo(fs): 				Datacenter: "dc1",
// todo(fs): 				Node:       "foo",
// todo(fs): 				Address:    "127.0.0.1",
// todo(fs): 				Service: &structs.NodeService{
// todo(fs): 					Service: "foo",
// todo(fs): 					Port:    12345,
// todo(fs): 				},
// todo(fs): 				WriteRequest: structs.WriteRequest{Token: "root"},
// todo(fs): 			}
// todo(fs): 			var out struct{}
// todo(fs): 			if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			// Set up the DNS query
// todo(fs): 			c := new(dns.Client)
// todo(fs): 			m := new(dns.Msg)
// todo(fs): 			m.SetQuestion("foo.service.consul.", dns.TypeA)
// todo(fs):
// todo(fs): 			in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 			if len(in.Answer) != tt.results {
// todo(fs): 				t.Fatalf("Bad: %#v", in)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_AddressLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Look up the addresses
// todo(fs): 	cases := map[string]string{
// todo(fs): 		"7f000001.addr.dc1.consul.": "127.0.0.1",
// todo(fs): 	}
// todo(fs): 	for question, answer := range cases {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aRec, ok := in.Answer[0].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if aRec.A.To4().String() != answer {
// todo(fs): 			t.Fatalf("Bad: %#v", aRec)
// todo(fs): 		}
// todo(fs): 		if aRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_AddressLookupIPV6(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Look up the addresses
// todo(fs): 	cases := map[string]string{
// todo(fs): 		"2607002040050808000000000000200e.addr.consul.": "2607:20:4005:808::200e",
// todo(fs): 		"2607112040051808ffffffffffff200e.addr.consul.": "2607:1120:4005:1808:ffff:ffff:ffff:200e",
// todo(fs): 	}
// todo(fs): 	for question, answer := range cases {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Answer) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		aaaaRec, ok := in.Answer[0].(*dns.AAAA)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 		if aaaaRec.AAAA.To16().String() != answer {
// todo(fs): 			t.Fatalf("Bad: %#v", aaaaRec)
// todo(fs): 		}
// todo(fs): 		if aaaaRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Answer[0])
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NonExistingLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// lookup a non-existing node, we should receive a SOA
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("nonexisting.consul.", dns.TypeANY)
// todo(fs):
// todo(fs): 	c := new(dns.Client)
// todo(fs): 	in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(in.Ns) != 1 {
// todo(fs): 		t.Fatalf("Bad: %#v %#v", in, len(in.Answer))
// todo(fs): 	}
// todo(fs):
// todo(fs): 	soaRec, ok := in.Ns[0].(*dns.SOA)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 	}
// todo(fs): 	if soaRec.Hdr.Ttl != 0 {
// todo(fs): 		t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_NonExistingLookupEmptyAorAAAA(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a v6-only service and a v4-only service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foov6",
// todo(fs): 			Address:    "fe80::1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "webv6",
// todo(fs): 				Port:    8000,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foov4",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "webv4",
// todo(fs): 				Port:    8000,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register equivalent prepared queries.
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "webv4",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "webv4",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var id string
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		args = &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "webv6",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "webv6",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check for ipv6 records on ipv4-only service directly and via the
// todo(fs): 	// prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"webv4.service.consul.",
// todo(fs): 		"webv4.query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeAAAA)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Ns) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs): 		soaRec, ok := in.Ns[0].(*dns.SOA)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs): 		if soaRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if in.Rcode != dns.RcodeSuccess {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check for ipv4 records on ipv6-only service directly and via the
// todo(fs): 	// prepared query.
// todo(fs): 	questions = []string{
// todo(fs): 		"webv6.service.consul.",
// todo(fs): 		"webv6.query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeA)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Ns) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		soaRec, ok := in.Ns[0].(*dns.SOA)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs): 		if soaRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if in.Rcode != dns.RcodeSuccess {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_PreparedQuery_AllowStale(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSConfig.AllowStale = Bool(true)
// todo(fs): 	cfg.DNSConfig.MaxStale = time.Second
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockPreparedQuery{
// todo(fs): 		executeFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExecuteResponse) error {
// todo(fs): 			// Return a response that's perpetually too stale.
// todo(fs): 			reply.LastContact = 2 * time.Second
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure that the lookup terminates and results in an SOA since
// todo(fs): 	// the query doesn't exist.
// todo(fs): 	{
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion("nope.query.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Ns) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		soaRec, ok := in.Ns[0].(*dns.SOA)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs): 		if soaRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_InvalidQueries(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Try invalid forms of queries that should hit the special invalid case
// todo(fs): 	// of our query parser.
// todo(fs): 	questions := []string{
// todo(fs): 		"consul.",
// todo(fs): 		"node.consul.",
// todo(fs): 		"service.consul.",
// todo(fs): 		"query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		in, _, err := c.Exchange(m, a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if len(in.Ns) != 1 {
// todo(fs): 			t.Fatalf("Bad: %#v", in)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		soaRec, ok := in.Ns[0].(*dns.SOA)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs): 		if soaRec.Hdr.Ttl != 0 {
// todo(fs): 			t.Fatalf("Bad: %#v", in.Ns[0])
// todo(fs): 		}
// todo(fs):
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_PreparedQuery_AgentSource(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockPreparedQuery{
// todo(fs): 		executeFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExecuteResponse) error {
// todo(fs): 			// Check that the agent inserted its self-name and datacenter to
// todo(fs): 			// the RPC request body.
// todo(fs): 			if args.Agent.Datacenter != a.Config.Datacenter ||
// todo(fs): 				args.Agent.Node != a.Config.NodeName {
// todo(fs): 				t.Fatalf("bad: %#v", args.Agent)
// todo(fs): 			}
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	{
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion("foo.query.consul.", dns.TypeSRV)
// todo(fs):
// todo(fs): 		c := new(dns.Client)
// todo(fs): 		if _, _, err := c.Exchange(m, a.DNSAddr()); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_trimUDPResponse_NoTrim(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	req := &dns.Msg{}
// todo(fs): 	resp := &dns.Msg{
// todo(fs): 		Answer: []dns.RR{
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 		Extra: []dns.RR{
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("10.0.1.185"),
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	config := &DefaultConfig().DNSConfig
// todo(fs): 	if trimmed := trimUDPResponse(config, req, resp); trimmed {
// todo(fs): 		t.Fatalf("Bad %#v", *resp)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	expected := &dns.Msg{
// todo(fs): 		Answer: []dns.RR{
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 		Extra: []dns.RR{
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("10.0.1.185"),
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if !reflect.DeepEqual(resp, expected) {
// todo(fs): 		t.Fatalf("Bad %#v vs. %#v", *resp, *expected)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_trimUDPResponse_TrimLimit(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	config := &DefaultConfig().DNSConfig
// todo(fs):
// todo(fs): 	req, resp, expected := &dns.Msg{}, &dns.Msg{}, &dns.Msg{}
// todo(fs): 	for i := 0; i < config.UDPAnswerLimit+1; i++ {
// todo(fs): 		target := fmt.Sprintf("ip-10-0-1-%d.node.dc1.consul.", 185+i)
// todo(fs): 		srv := &dns.SRV{
// todo(fs): 			Hdr: dns.RR_Header{
// todo(fs): 				Name:   "redis-cache-redis.service.consul.",
// todo(fs): 				Rrtype: dns.TypeSRV,
// todo(fs): 				Class:  dns.ClassINET,
// todo(fs): 			},
// todo(fs): 			Target: target,
// todo(fs): 		}
// todo(fs): 		a := &dns.A{
// todo(fs): 			Hdr: dns.RR_Header{
// todo(fs): 				Name:   target,
// todo(fs): 				Rrtype: dns.TypeA,
// todo(fs): 				Class:  dns.ClassINET,
// todo(fs): 			},
// todo(fs): 			A: net.ParseIP(fmt.Sprintf("10.0.1.%d", 185+i)),
// todo(fs): 		}
// todo(fs):
// todo(fs): 		resp.Answer = append(resp.Answer, srv)
// todo(fs): 		resp.Extra = append(resp.Extra, a)
// todo(fs): 		if i < config.UDPAnswerLimit {
// todo(fs): 			expected.Answer = append(expected.Answer, srv)
// todo(fs): 			expected.Extra = append(expected.Extra, a)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if trimmed := trimUDPResponse(config, req, resp); !trimmed {
// todo(fs): 		t.Fatalf("Bad %#v", *resp)
// todo(fs): 	}
// todo(fs): 	if !reflect.DeepEqual(resp, expected) {
// todo(fs): 		t.Fatalf("Bad %#v vs. %#v", *resp, *expected)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_trimUDPResponse_TrimSize(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	config := &DefaultConfig().DNSConfig
// todo(fs):
// todo(fs): 	req, resp := &dns.Msg{}, &dns.Msg{}
// todo(fs): 	for i := 0; i < 100; i++ {
// todo(fs): 		target := fmt.Sprintf("ip-10-0-1-%d.node.dc1.consul.", 185+i)
// todo(fs): 		srv := &dns.SRV{
// todo(fs): 			Hdr: dns.RR_Header{
// todo(fs): 				Name:   "redis-cache-redis.service.consul.",
// todo(fs): 				Rrtype: dns.TypeSRV,
// todo(fs): 				Class:  dns.ClassINET,
// todo(fs): 			},
// todo(fs): 			Target: target,
// todo(fs): 		}
// todo(fs): 		a := &dns.A{
// todo(fs): 			Hdr: dns.RR_Header{
// todo(fs): 				Name:   target,
// todo(fs): 				Rrtype: dns.TypeA,
// todo(fs): 				Class:  dns.ClassINET,
// todo(fs): 			},
// todo(fs): 			A: net.ParseIP(fmt.Sprintf("10.0.1.%d", 185+i)),
// todo(fs): 		}
// todo(fs):
// todo(fs): 		resp.Answer = append(resp.Answer, srv)
// todo(fs): 		resp.Extra = append(resp.Extra, a)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// We don't know the exact trim, but we know the resulting answer
// todo(fs): 	// data should match its extra data.
// todo(fs): 	if trimmed := trimUDPResponse(config, req, resp); !trimmed {
// todo(fs): 		t.Fatalf("Bad %#v", *resp)
// todo(fs): 	}
// todo(fs): 	if len(resp.Answer) == 0 || len(resp.Answer) != len(resp.Extra) {
// todo(fs): 		t.Fatalf("Bad %#v", *resp)
// todo(fs): 	}
// todo(fs): 	for i := range resp.Answer {
// todo(fs): 		srv, ok := resp.Answer[i].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should be SRV")
// todo(fs): 		}
// todo(fs):
// todo(fs): 		a, ok := resp.Extra[i].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should be A")
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if srv.Target != a.Header().Name {
// todo(fs): 			t.Fatalf("Bad %#v vs. %#v", *srv, *a)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_trimUDPResponse_TrimSizeEDNS(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	config := &DefaultConfig().DNSConfig
// todo(fs):
// todo(fs): 	req, resp := &dns.Msg{}, &dns.Msg{}
// todo(fs):
// todo(fs): 	for i := 0; i < 100; i++ {
// todo(fs): 		target := fmt.Sprintf("ip-10-0-1-%d.node.dc1.consul.", 150+i)
// todo(fs): 		srv := &dns.SRV{
// todo(fs): 			Hdr: dns.RR_Header{
// todo(fs): 				Name:   "redis-cache-redis.service.consul.",
// todo(fs): 				Rrtype: dns.TypeSRV,
// todo(fs): 				Class:  dns.ClassINET,
// todo(fs): 			},
// todo(fs): 			Target: target,
// todo(fs): 		}
// todo(fs): 		a := &dns.A{
// todo(fs): 			Hdr: dns.RR_Header{
// todo(fs): 				Name:   target,
// todo(fs): 				Rrtype: dns.TypeA,
// todo(fs): 				Class:  dns.ClassINET,
// todo(fs): 			},
// todo(fs): 			A: net.ParseIP(fmt.Sprintf("10.0.1.%d", 150+i)),
// todo(fs): 		}
// todo(fs):
// todo(fs): 		resp.Answer = append(resp.Answer, srv)
// todo(fs): 		resp.Extra = append(resp.Extra, a)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Copy over to a new slice since we are trimming both.
// todo(fs): 	reqEDNS, respEDNS := &dns.Msg{}, &dns.Msg{}
// todo(fs): 	reqEDNS.SetEdns0(2048, true)
// todo(fs): 	respEDNS.Answer = append(respEDNS.Answer, resp.Answer...)
// todo(fs): 	respEDNS.Extra = append(respEDNS.Extra, resp.Extra...)
// todo(fs):
// todo(fs): 	// Trim each response
// todo(fs): 	if trimmed := trimUDPResponse(config, req, resp); !trimmed {
// todo(fs): 		t.Errorf("expected response to be trimmed: %#v", resp)
// todo(fs): 	}
// todo(fs): 	if trimmed := trimUDPResponse(config, reqEDNS, respEDNS); !trimmed {
// todo(fs): 		t.Errorf("expected edns to be trimmed: %#v", resp)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check answer lengths
// todo(fs): 	if len(resp.Answer) == 0 || len(resp.Answer) != len(resp.Extra) {
// todo(fs): 		t.Errorf("bad response answer length: %#v", resp)
// todo(fs): 	}
// todo(fs): 	if len(respEDNS.Answer) == 0 || len(respEDNS.Answer) != len(respEDNS.Extra) {
// todo(fs): 		t.Errorf("bad edns answer length: %#v", resp)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Due to the compression, we can't check exact equality of sizes, but we can
// todo(fs): 	// make two requests and ensure that the edns one returns a larger payload
// todo(fs): 	// than the non-edns0 one.
// todo(fs): 	if len(resp.Answer) >= len(respEDNS.Answer) {
// todo(fs): 		t.Errorf("expected edns have larger answer: %#v\n%#v", resp, respEDNS)
// todo(fs): 	}
// todo(fs): 	if len(resp.Extra) >= len(respEDNS.Extra) {
// todo(fs): 		t.Errorf("expected edns have larger extra: %#v\n%#v", resp, respEDNS)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify that the things point where they should
// todo(fs): 	for i := range resp.Answer {
// todo(fs): 		srv, ok := resp.Answer[i].(*dns.SRV)
// todo(fs): 		if !ok {
// todo(fs): 			t.Errorf("%d should be an SRV", i)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		a, ok := resp.Extra[i].(*dns.A)
// todo(fs): 		if !ok {
// todo(fs): 			t.Errorf("%d should be an A", i)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if srv.Target != a.Header().Name {
// todo(fs): 			t.Errorf("%d: bad %#v vs. %#v", i, srv, a)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_syncExtra(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := &dns.Msg{
// todo(fs): 		Answer: []dns.RR{
// todo(fs): 			// These two are on the same host so the redundant extra
// todo(fs): 			// records should get deduplicated.
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Port:   1001,
// todo(fs): 				Target: "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 			},
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Port:   1002,
// todo(fs): 				Target: "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 			},
// todo(fs): 			// This one isn't in the Consul domain so it will get a
// todo(fs): 			// CNAME and then an A record from the recursor.
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Port:   1003,
// todo(fs): 				Target: "demo.consul.io.",
// todo(fs): 			},
// todo(fs): 			// This one isn't in the Consul domain and it will get
// todo(fs): 			// a CNAME and A record from a recursor that alters the
// todo(fs): 			// case of the name. This proves we look up in the index
// todo(fs): 			// in a case-insensitive way.
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Port:   1001,
// todo(fs): 				Target: "insensitive.consul.io.",
// todo(fs): 			},
// todo(fs): 			// This is also a CNAME, but it'll be set up to loop to
// todo(fs): 			// make sure we don't crash.
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Port:   1001,
// todo(fs): 				Target: "deadly.consul.io.",
// todo(fs): 			},
// todo(fs): 			// This is also a CNAME, but it won't have another record.
// todo(fs): 			&dns.SRV{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "redis-cache-redis.service.consul.",
// todo(fs): 					Rrtype: dns.TypeSRV,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Port:   1001,
// todo(fs): 				Target: "nope.consul.io.",
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 		Extra: []dns.RR{
// todo(fs): 			// These should get deduplicated.
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("10.0.1.185"),
// todo(fs): 			},
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("10.0.1.185"),
// todo(fs): 			},
// todo(fs): 			// This is a normal CNAME followed by an A record but we
// todo(fs): 			// have flipped the order. The algorithm should emit them
// todo(fs): 			// in the opposite order.
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "fakeserver.consul.io.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("127.0.0.1"),
// todo(fs): 			},
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "demo.consul.io.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "fakeserver.consul.io.",
// todo(fs): 			},
// todo(fs): 			// These differ in case to test case insensitivity.
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "INSENSITIVE.CONSUL.IO.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "Another.Server.Com.",
// todo(fs): 			},
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "another.server.com.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("127.0.0.1"),
// todo(fs): 			},
// todo(fs): 			// This doesn't appear in the answer, so should get
// todo(fs): 			// dropped.
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "ip-10-0-1-186.node.dc1.consul.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("10.0.1.186"),
// todo(fs): 			},
// todo(fs): 			// These two test edge cases with CNAME handling.
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "deadly.consul.io.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "deadly.consul.io.",
// todo(fs): 			},
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "nope.consul.io.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "notthere.consul.io.",
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	index := make(map[string]dns.RR)
// todo(fs): 	indexRRs(resp.Extra, index)
// todo(fs): 	syncExtra(index, resp)
// todo(fs):
// todo(fs): 	expected := &dns.Msg{
// todo(fs): 		Answer: resp.Answer,
// todo(fs): 		Extra: []dns.RR{
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "ip-10-0-1-185.node.dc1.consul.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("10.0.1.185"),
// todo(fs): 			},
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "demo.consul.io.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "fakeserver.consul.io.",
// todo(fs): 			},
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "fakeserver.consul.io.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("127.0.0.1"),
// todo(fs): 			},
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "INSENSITIVE.CONSUL.IO.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "Another.Server.Com.",
// todo(fs): 			},
// todo(fs): 			&dns.A{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "another.server.com.",
// todo(fs): 					Rrtype: dns.TypeA,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				A: net.ParseIP("127.0.0.1"),
// todo(fs): 			},
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "deadly.consul.io.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "deadly.consul.io.",
// todo(fs): 			},
// todo(fs): 			&dns.CNAME{
// todo(fs): 				Hdr: dns.RR_Header{
// todo(fs): 					Name:   "nope.consul.io.",
// todo(fs): 					Rrtype: dns.TypeCNAME,
// todo(fs): 					Class:  dns.ClassINET,
// todo(fs): 				},
// todo(fs): 				Target: "notthere.consul.io.",
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if !reflect.DeepEqual(resp, expected) {
// todo(fs): 		t.Fatalf("Bad %#v vs. %#v", *resp, *expected)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_Compression_trimUDPResponse(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	config := &DefaultConfig().DNSConfig
// todo(fs):
// todo(fs): 	req, m := dns.Msg{}, dns.Msg{}
// todo(fs): 	trimUDPResponse(config, &req, &m)
// todo(fs): 	if m.Compress {
// todo(fs): 		t.Fatalf("compression should be off")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// The trim function temporarily turns off compression, so we need to
// todo(fs): 	// make sure the setting gets restored properly.
// todo(fs): 	m.Compress = true
// todo(fs): 	trimUDPResponse(config, &req, &m)
// todo(fs): 	if !m.Compress {
// todo(fs): 		t.Fatalf("compression should be on")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_Compression_Query(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "db",
// todo(fs): 				Tags:    []string{"master"},
// todo(fs): 				Port:    12345,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register an equivalent prepared query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		args := &structs.PreparedQueryRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Op:         structs.PreparedQueryCreate,
// todo(fs): 			Query: &structs.PreparedQuery{
// todo(fs): 				Name: "test",
// todo(fs): 				Service: structs.ServiceQuery{
// todo(fs): 					Service: "db",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.RPC("PreparedQuery.Apply", args, &id); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Look up the service directly and via prepared query.
// todo(fs): 	questions := []string{
// todo(fs): 		"db.service.consul.",
// todo(fs): 		id + ".query.consul.",
// todo(fs): 	}
// todo(fs): 	for _, question := range questions {
// todo(fs): 		m := new(dns.Msg)
// todo(fs): 		m.SetQuestion(question, dns.TypeSRV)
// todo(fs):
// todo(fs): 		conn, err := dns.Dial("udp", a.DNSAddr())
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Do a manual exchange with compression on (the default).
// todo(fs): 		a.DNSDisableCompression(false)
// todo(fs): 		if err := conn.WriteMsg(m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		p := make([]byte, dns.MaxMsgSize)
// todo(fs): 		compressed, err := conn.Read(p)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Disable compression and try again.
// todo(fs): 		a.DNSDisableCompression(true)
// todo(fs): 		if err := conn.WriteMsg(m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		unc, err := conn.Read(p)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We can't see the compressed status given the DNS API, so we
// todo(fs): 		// just make sure the message is smaller to see if it's
// todo(fs): 		// respecting the flag.
// todo(fs): 		if compressed == 0 || unc == 0 || compressed >= unc {
// todo(fs): 			t.Fatalf("'%s' doesn't look compressed: %d vs. %d", question, compressed, unc)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_Compression_ReverseLookup(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node.
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo2",
// todo(fs): 		Address:    "127.0.0.2",
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("2.0.0.127.in-addr.arpa.", dns.TypeANY)
// todo(fs):
// todo(fs): 	conn, err := dns.Dial("udp", a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Do a manual exchange with compression on (the default).
// todo(fs): 	if err := conn.WriteMsg(m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	p := make([]byte, dns.MaxMsgSize)
// todo(fs): 	compressed, err := conn.Read(p)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Disable compression and try again.
// todo(fs): 	a.DNSDisableCompression(true)
// todo(fs): 	if err := conn.WriteMsg(m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	unc, err := conn.Read(p)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// We can't see the compressed status given the DNS API, so we just make
// todo(fs): 	// sure the message is smaller to see if it's respecting the flag.
// todo(fs): 	if compressed == 0 || unc == 0 || compressed >= unc {
// todo(fs): 		t.Fatalf("doesn't look compressed: %d vs. %d", compressed, unc)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNS_Compression_Recurse(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	recursor := makeRecursor(t, dns.Msg{
// todo(fs): 		Answer: []dns.RR{dnsA("apple.com", "1.2.3.4")},
// todo(fs): 	})
// todo(fs): 	defer recursor.Shutdown()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.DNSRecursor = recursor.Addr
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := new(dns.Msg)
// todo(fs): 	m.SetQuestion("apple.com.", dns.TypeANY)
// todo(fs):
// todo(fs): 	conn, err := dns.Dial("udp", a.DNSAddr())
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Do a manual exchange with compression on (the default).
// todo(fs): 	if err := conn.WriteMsg(m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	p := make([]byte, dns.MaxMsgSize)
// todo(fs): 	compressed, err := conn.Read(p)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Disable compression and try again.
// todo(fs): 	a.DNSDisableCompression(true)
// todo(fs): 	if err := conn.WriteMsg(m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	unc, err := conn.Read(p)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// We can't see the compressed status given the DNS API, so we just make
// todo(fs): 	// sure the message is smaller to see if it's respecting the flag.
// todo(fs): 	if compressed == 0 || unc == 0 || compressed >= unc {
// todo(fs): 		t.Fatalf("doesn't look compressed: %d vs. %d", compressed, unc)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestDNSInvalidRegex(t *testing.T) {
// todo(fs): 	tests := []struct {
// todo(fs): 		desc    string
// todo(fs): 		in      string
// todo(fs): 		invalid bool
// todo(fs): 	}{
// todo(fs): 		{"Valid Hostname", "testnode", false},
// todo(fs): 		{"Valid Hostname", "test-node", false},
// todo(fs): 		{"Invalid Hostname with special chars", "test#$$!node", true},
// todo(fs): 		{"Invalid Hostname with special chars in the end", "testnode%^", true},
// todo(fs): 		{"Whitespace", "  ", true},
// todo(fs): 		{"Only special chars", "./$", true},
// todo(fs): 	}
// todo(fs): 	for _, test := range tests {
// todo(fs): 		t.Run(test.desc, func(t *testing.T) {
// todo(fs): 			if got, want := InvalidDnsRe.MatchString(test.in), test.invalid; got != want {
// todo(fs): 				t.Fatalf("Expected %v to return %v", test.in, want)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs):
// todo(fs): 	}
// todo(fs): }
