package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"
)

// todo(fs): func TestHTTPServer_UnixSocket(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	if runtime.GOOS == "windows" {
// todo(fs): 		t.SkipNow()
// todo(fs): 	}
// todo(fs):
// todo(fs): 	tempDir := testutil.TempDir(t, "consul")
// todo(fs): 	defer os.RemoveAll(tempDir)
// todo(fs): 	socket := filepath.Join(tempDir, "test.sock")
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.Addresses.HTTP = "unix://" + socket
// todo(fs):
// todo(fs): 	// Only testing mode, since uid/gid might not be settable
// todo(fs): 	// from test environment.
// todo(fs): 	cfg.UnixSockets = UnixSocketConfig{}
// todo(fs): 	cfg.UnixSockets.Perms = "0777"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Ensure the socket was created
// todo(fs): 	if _, err := os.Stat(socket); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the mode was set properly
// todo(fs): 	fi, err := os.Stat(socket)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if fi.Mode().String() != "Srwxrwxrwx" {
// todo(fs): 		t.Fatalf("bad permissions: %s", fi.Mode())
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we can get a response from the socket.
// todo(fs): 	path := socketPath(a.Config.Addresses.HTTP)
// todo(fs): 	trans := cleanhttp.DefaultTransport()
// todo(fs): 	trans.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
// todo(fs): 		return net.Dial("unix", path)
// todo(fs): 	}
// todo(fs): 	client := &http.Client{
// todo(fs): 		Transport: trans,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// This URL doesn't look like it makes sense, but the scheme (http://) and
// todo(fs): 	// the host (127.0.0.1) are required by the HTTP client library. In reality
// todo(fs): 	// this will just use the custom dialer and talk to the socket.
// todo(fs): 	resp, err := client.Get("http://127.0.0.1/v1/agent/self")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	defer resp.Body.Close()
// todo(fs):
// todo(fs): 	if body, err := ioutil.ReadAll(resp.Body); err != nil || len(body) == 0 {
// todo(fs): 		t.Fatalf("bad: %s %v", body, err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHTTPServer_UnixSocket_FileExists(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	if runtime.GOOS == "windows" {
// todo(fs): 		t.SkipNow()
// todo(fs): 	}
// todo(fs):
// todo(fs): 	tempDir := testutil.TempDir(t, "consul")
// todo(fs): 	defer os.RemoveAll(tempDir)
// todo(fs): 	socket := filepath.Join(tempDir, "test.sock")
// todo(fs):
// todo(fs): 	// Create a regular file at the socket path
// todo(fs): 	if err := ioutil.WriteFile(socket, []byte("hello world"), 0644); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	fi, err := os.Stat(socket)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if !fi.Mode().IsRegular() {
// todo(fs): 		t.Fatalf("not a regular file: %s", socket)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.Addresses.HTTP = "unix://" + socket
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Ensure the file was replaced by the socket
// todo(fs): 	fi, err = os.Stat(socket)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if fi.Mode()&os.ModeSocket == 0 {
// todo(fs): 		t.Fatalf("expected socket to replace file")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSetIndex(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	setIndex(resp, 1000)
// todo(fs): 	header := resp.Header().Get("X-Consul-Index")
// todo(fs): 	if header != "1000" {
// todo(fs): 		t.Fatalf("Bad: %v", header)
// todo(fs): 	}
// todo(fs): 	setIndex(resp, 2000)
// todo(fs): 	if v := resp.Header()["X-Consul-Index"]; len(v) != 1 {
// todo(fs): 		t.Fatalf("bad: %#v", v)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSetKnownLeader(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	setKnownLeader(resp, true)
// todo(fs): 	header := resp.Header().Get("X-Consul-KnownLeader")
// todo(fs): 	if header != "true" {
// todo(fs): 		t.Fatalf("Bad: %v", header)
// todo(fs): 	}
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	setKnownLeader(resp, false)
// todo(fs): 	header = resp.Header().Get("X-Consul-KnownLeader")
// todo(fs): 	if header != "false" {
// todo(fs): 		t.Fatalf("Bad: %v", header)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSetLastContact(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	tests := []struct {
// todo(fs): 		desc string
// todo(fs): 		d    time.Duration
// todo(fs): 		h    string
// todo(fs): 	}{
// todo(fs): 		{"neg", -1, "0"},
// todo(fs): 		{"zero", 0, "0"},
// todo(fs): 		{"pos", 123 * time.Millisecond, "123"},
// todo(fs): 		{"pos ms only", 123456 * time.Microsecond, "123"},
// todo(fs): 	}
// todo(fs): 	for _, tt := range tests {
// todo(fs): 		t.Run(tt.desc, func(t *testing.T) {
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			setLastContact(resp, tt.d)
// todo(fs): 			header := resp.Header().Get("X-Consul-LastContact")
// todo(fs): 			if got, want := header, tt.h; got != want {
// todo(fs): 				t.Fatalf("got X-Consul-LastContact header %q want %q", got, want)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSetMeta(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	meta := structs.QueryMeta{
// todo(fs): 		Index:       1000,
// todo(fs): 		KnownLeader: true,
// todo(fs): 		LastContact: 123456 * time.Microsecond,
// todo(fs): 	}
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	setMeta(resp, &meta)
// todo(fs): 	header := resp.Header().Get("X-Consul-Index")
// todo(fs): 	if header != "1000" {
// todo(fs): 		t.Fatalf("Bad: %v", header)
// todo(fs): 	}
// todo(fs): 	header = resp.Header().Get("X-Consul-KnownLeader")
// todo(fs): 	if header != "true" {
// todo(fs): 		t.Fatalf("Bad: %v", header)
// todo(fs): 	}
// todo(fs): 	header = resp.Header().Get("X-Consul-LastContact")
// todo(fs): 	if header != "123" {
// todo(fs): 		t.Fatalf("Bad: %v", header)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHTTPAPI_BlockEndpoints(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs):
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.HTTPConfig.BlockEndpoints = []string{
// todo(fs): 		"/v1/agent/self",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	handler := func(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
// todo(fs): 		return nil, nil
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try a blocked endpoint, which should get a 403.
// todo(fs): 	{
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/self", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		a.srv.wrap(handler)(resp, req)
// todo(fs): 		if got, want := resp.Code, http.StatusForbidden; got != want {
// todo(fs): 			t.Fatalf("bad response code got %d want %d", got, want)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure some other endpoint still works.
// todo(fs): 	{
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/checks", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		a.srv.wrap(handler)(resp, req)
// todo(fs): 		if got, want := resp.Code, http.StatusOK; got != want {
// todo(fs): 			t.Fatalf("bad response code got %d want %d", got, want)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHTTPAPI_TranslateAddrHeader(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	// Header should not be present if address translation is off.
// todo(fs): 	{
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		handler := func(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
// todo(fs): 			return nil, nil
// todo(fs): 		}
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/self", nil)
// todo(fs): 		a.srv.wrap(handler)(resp, req)
// todo(fs):
// todo(fs): 		translate := resp.Header().Get("X-Consul-Translate-Addresses")
// todo(fs): 		if translate != "" {
// todo(fs): 			t.Fatalf("bad: expected %q, got %q", "", translate)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Header should be set to true if it's turned on.
// todo(fs): 	{
// todo(fs): 		cfg := TestConfig()
// todo(fs): 		cfg.TranslateWanAddrs = true
// todo(fs): 		a := NewTestAgent(t.Name(), cfg)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		handler := func(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
// todo(fs): 			return nil, nil
// todo(fs): 		}
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/self", nil)
// todo(fs): 		a.srv.wrap(handler)(resp, req)
// todo(fs):
// todo(fs): 		translate := resp.Header().Get("X-Consul-Translate-Addresses")
// todo(fs): 		if translate != "true" {
// todo(fs): 			t.Fatalf("bad: expected %q, got %q", "true", translate)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHTTPAPIResponseHeaders(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.HTTPConfig.ResponseHeaders = map[string]string{
// todo(fs): 		"Access-Control-Allow-Origin": "*",
// todo(fs): 		"X-XSS-Protection":            "1; mode=block",
// todo(fs): 	}
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	handler := func(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
// todo(fs): 		return nil, nil
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/self", nil)
// todo(fs): 	a.srv.wrap(handler)(resp, req)
// todo(fs):
// todo(fs): 	origin := resp.Header().Get("Access-Control-Allow-Origin")
// todo(fs): 	if origin != "*" {
// todo(fs): 		t.Fatalf("bad Access-Control-Allow-Origin: expected %q, got %q", "*", origin)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	xss := resp.Header().Get("X-XSS-Protection")
// todo(fs): 	if xss != "1; mode=block" {
// todo(fs): 		t.Fatalf("bad X-XSS-Protection header: expected %q, got %q", "1; mode=block", xss)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestContentTypeIsJSON(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	handler := func(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
// todo(fs): 		// stub out a DirEntry so that it will be encoded as JSON
// todo(fs): 		return &structs.DirEntry{Key: "key"}, nil
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/kv/key", nil)
// todo(fs): 	a.srv.wrap(handler)(resp, req)
// todo(fs):
// todo(fs): 	contentType := resp.Header().Get("Content-Type")
// todo(fs):
// todo(fs): 	if contentType != "application/json" {
// todo(fs): 		t.Fatalf("Content-Type header was not 'application/json'")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHTTP_wrap_obfuscateLog(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	buf := new(bytes.Buffer)
// todo(fs): 	a := &TestAgent{Name: t.Name(), LogOutput: buf}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	handler := func(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
// todo(fs): 		return nil, nil
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, pair := range [][]string{
// todo(fs): 		{
// todo(fs): 			"/some/url?token=secret1&token=secret2",
// todo(fs): 			"/some/url?token=<hidden>&token=<hidden>",
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			"/v1/acl/clone/secret1",
// todo(fs): 			"/v1/acl/clone/<hidden>",
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			"/v1/acl/clone/secret1?token=secret2",
// todo(fs): 			"/v1/acl/clone/<hidden>?token=<hidden>",
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			"/v1/acl/destroy/secret1",
// todo(fs): 			"/v1/acl/destroy/<hidden>",
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			"/v1/acl/destroy/secret1?token=secret2",
// todo(fs): 			"/v1/acl/destroy/<hidden>?token=<hidden>",
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			"/v1/acl/info/secret1",
// todo(fs): 			"/v1/acl/info/<hidden>",
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			"/v1/acl/info/secret1?token=secret2",
// todo(fs): 			"/v1/acl/info/<hidden>?token=<hidden>",
// todo(fs): 		},
// todo(fs): 	} {
// todo(fs): 		url, want := pair[0], pair[1]
// todo(fs): 		t.Run(url, func(t *testing.T) {
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			req, _ := http.NewRequest("GET", url, nil)
// todo(fs): 			a.srv.wrap(handler)(resp, req)
// todo(fs):
// todo(fs): 			if got := buf.String(); !strings.Contains(got, want) {
// todo(fs): 				t.Fatalf("got %s want %s", got, want)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestPrettyPrint(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	testPrettyPrint("pretty=1", t)
// todo(fs): }
// todo(fs):
// todo(fs): func TestPrettyPrintBare(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	testPrettyPrint("pretty", t)
// todo(fs): }
// todo(fs):
// todo(fs): func testPrettyPrint(pretty string, t *testing.T) {
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	r := &structs.DirEntry{Key: "key"}
// todo(fs):
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	handler := func(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
// todo(fs): 		return r, nil
// todo(fs): 	}
// todo(fs):
// todo(fs): 	urlStr := "/v1/kv/key?" + pretty
// todo(fs): 	req, _ := http.NewRequest("GET", urlStr, nil)
// todo(fs): 	a.srv.wrap(handler)(resp, req)
// todo(fs):
// todo(fs): 	expected, _ := json.MarshalIndent(r, "", "    ")
// todo(fs): 	expected = append(expected, "\n"...)
// todo(fs): 	actual, err := ioutil.ReadAll(resp.Body)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if !bytes.Equal(expected, actual) {
// todo(fs): 		t.Fatalf("bad: %q", string(actual))
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestParseSource(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Default is agent's DC and no node (since the user didn't care, then
// todo(fs): 	// just give them the cheapest possible query).
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes", nil)
// todo(fs): 	source := structs.QuerySource{}
// todo(fs): 	a.srv.parseSource(req, &source)
// todo(fs): 	if source.Datacenter != "dc1" || source.Node != "" {
// todo(fs): 		t.Fatalf("bad: %v", source)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Adding the source parameter should set that node.
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/catalog/nodes?near=bob", nil)
// todo(fs): 	source = structs.QuerySource{}
// todo(fs): 	a.srv.parseSource(req, &source)
// todo(fs): 	if source.Datacenter != "dc1" || source.Node != "bob" {
// todo(fs): 		t.Fatalf("bad: %v", source)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// We should follow whatever dc parameter was given so that the node is
// todo(fs): 	// looked up correctly on the receiving end.
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/catalog/nodes?near=bob&dc=foo", nil)
// todo(fs): 	source = structs.QuerySource{}
// todo(fs): 	a.srv.parseSource(req, &source)
// todo(fs): 	if source.Datacenter != "foo" || source.Node != "bob" {
// todo(fs): 		t.Fatalf("bad: %v", source)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// The magic "_agent" node name will use the agent's local node name.
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/catalog/nodes?near=_agent", nil)
// todo(fs): 	source = structs.QuerySource{}
// todo(fs): 	a.srv.parseSource(req, &source)
// todo(fs): 	if source.Datacenter != "dc1" || source.Node != a.Config.NodeName {
// todo(fs): 		t.Fatalf("bad: %v", source)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestParseWait(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	var b structs.QueryOptions
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?wait=60s&index=1000", nil)
// todo(fs): 	if d := parseWait(resp, req, &b); d {
// todo(fs): 		t.Fatalf("unexpected done")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if b.MinQueryIndex != 1000 {
// todo(fs): 		t.Fatalf("Bad: %v", b)
// todo(fs): 	}
// todo(fs): 	if b.MaxQueryTime != 60*time.Second {
// todo(fs): 		t.Fatalf("Bad: %v", b)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestParseWait_InvalidTime(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	var b structs.QueryOptions
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?wait=60foo&index=1000", nil)
// todo(fs): 	if d := parseWait(resp, req, &b); !d {
// todo(fs): 		t.Fatalf("expected done")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("bad code: %v", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestParseWait_InvalidIndex(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	var b structs.QueryOptions
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?wait=60s&index=foo", nil)
// todo(fs): 	if d := parseWait(resp, req, &b); !d {
// todo(fs): 		t.Fatalf("expected done")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("bad code: %v", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestParseConsistency(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	var b structs.QueryOptions
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?stale", nil)
// todo(fs): 	if d := parseConsistency(resp, req, &b); d {
// todo(fs): 		t.Fatalf("unexpected done")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if !b.AllowStale {
// todo(fs): 		t.Fatalf("Bad: %v", b)
// todo(fs): 	}
// todo(fs): 	if b.RequireConsistent {
// todo(fs): 		t.Fatalf("Bad: %v", b)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	b = structs.QueryOptions{}
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/catalog/nodes?consistent", nil)
// todo(fs): 	if d := parseConsistency(resp, req, &b); d {
// todo(fs): 		t.Fatalf("unexpected done")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if b.AllowStale {
// todo(fs): 		t.Fatalf("Bad: %v", b)
// todo(fs): 	}
// todo(fs): 	if !b.RequireConsistent {
// todo(fs): 		t.Fatalf("Bad: %v", b)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestParseConsistency_Invalid(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	var b structs.QueryOptions
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?stale&consistent", nil)
// todo(fs): 	if d := parseConsistency(resp, req, &b); !d {
// todo(fs): 		t.Fatalf("expected done")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("bad code: %v", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): // Test ACL token is resolved in correct order
// todo(fs): func TestACLResolution(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	var token string
// todo(fs): 	// Request without token
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes", nil)
// todo(fs): 	// Request with explicit token
// todo(fs): 	reqToken, _ := http.NewRequest("GET", "/v1/catalog/nodes?token=foo", nil)
// todo(fs): 	// Request with header token only
// todo(fs): 	reqHeaderToken, _ := http.NewRequest("GET", "/v1/catalog/nodes", nil)
// todo(fs): 	reqHeaderToken.Header.Add("X-Consul-Token", "bar")
// todo(fs):
// todo(fs): 	// Request with header and querystring tokens
// todo(fs): 	reqBothTokens, _ := http.NewRequest("GET", "/v1/catalog/nodes?token=baz", nil)
// todo(fs): 	reqBothTokens.Header.Add("X-Consul-Token", "zap")
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Check when no token is set
// todo(fs): 	a.tokens.UpdateUserToken("")
// todo(fs): 	a.srv.parseToken(req, &token)
// todo(fs): 	if token != "" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check when ACLToken set
// todo(fs): 	a.tokens.UpdateUserToken("agent")
// todo(fs): 	a.srv.parseToken(req, &token)
// todo(fs): 	if token != "agent" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Explicit token has highest precedence
// todo(fs): 	a.srv.parseToken(reqToken, &token)
// todo(fs): 	if token != "foo" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Header token has precedence over agent token
// todo(fs): 	a.srv.parseToken(reqHeaderToken, &token)
// todo(fs): 	if token != "bar" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Querystring token has precedence over header and agent tokens
// todo(fs): 	a.srv.parseToken(reqBothTokens, &token)
// todo(fs): 	if token != "baz" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestEnableWebUI(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EnableUI = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/ui/", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	a.srv.Handler.ServeHTTP(resp, req)
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("should handle ui")
// todo(fs): 	}
// todo(fs): }

// assertIndex tests that X-Consul-Index is set and non-zero
func assertIndex(t *testing.T, resp *httptest.ResponseRecorder) {
	header := resp.Header().Get("X-Consul-Index")
	if header == "" || header == "0" {
		t.Fatalf("Bad: %v", header)
	}
}

// checkIndex is like assertIndex but returns an error
func checkIndex(resp *httptest.ResponseRecorder) error {
	header := resp.Header().Get("X-Consul-Index")
	if header == "" || header == "0" {
		return fmt.Errorf("Bad: %v", header)
	}
	return nil
}

// getIndex parses X-Consul-Index
func getIndex(t *testing.T, resp *httptest.ResponseRecorder) uint64 {
	header := resp.Header().Get("X-Consul-Index")
	if header == "" {
		t.Fatalf("Bad: %v", header)
	}
	val, err := strconv.Atoi(header)
	if err != nil {
		t.Fatalf("Bad: %v", header)
	}
	return uint64(val)
}

func jsonReader(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		panic(err)
	}
	return b
}
