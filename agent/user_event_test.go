package agent

import (
	"strings"
	"testing"
)

func TestValidateUserEventParams(t *testing.T) {
	t.Parallel()
	p := &UserEvent{}
	err := validateUserEventParams(p)
	if err == nil || err.Error() != "User event missing name" {
		t.Fatalf("err: %v", err)
	}
	p.Name = "foo"

	p.NodeFilter = "("
	err = validateUserEventParams(p)
	if err == nil || !strings.Contains(err.Error(), "Invalid node filter") {
		t.Fatalf("err: %v", err)
	}

	p.NodeFilter = ""
	p.ServiceFilter = "("
	err = validateUserEventParams(p)
	if err == nil || !strings.Contains(err.Error(), "Invalid service filter") {
		t.Fatalf("err: %v", err)
	}

	p.ServiceFilter = "foo"
	p.TagFilter = "("
	err = validateUserEventParams(p)
	if err == nil || !strings.Contains(err.Error(), "Invalid tag filter") {
		t.Fatalf("err: %v", err)
	}

	p.ServiceFilter = ""
	p.TagFilter = "foo"
	err = validateUserEventParams(p)
	if err == nil || !strings.Contains(err.Error(), "tag filter without service") {
		t.Fatalf("err: %v", err)
	}
}

// todo(fs): func TestShouldProcessUserEvent(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:      "mysql",
// todo(fs): 		Service: "mysql",
// todo(fs): 		Tags:    []string{"test", "foo", "bar", "master"},
// todo(fs): 		Port:    5000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, "")
// todo(fs):
// todo(fs): 	p := &UserEvent{}
// todo(fs): 	if !a.shouldProcessUserEvent(p) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Bad node name
// todo(fs): 	p = &UserEvent{
// todo(fs): 		NodeFilter: "foobar",
// todo(fs): 	}
// todo(fs): 	if a.shouldProcessUserEvent(p) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Good node name
// todo(fs): 	p = &UserEvent{
// todo(fs): 		NodeFilter: "^Node",
// todo(fs): 	}
// todo(fs): 	if !a.shouldProcessUserEvent(p) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Bad service name
// todo(fs): 	p = &UserEvent{
// todo(fs): 		ServiceFilter: "foobar",
// todo(fs): 	}
// todo(fs): 	if a.shouldProcessUserEvent(p) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Good service name
// todo(fs): 	p = &UserEvent{
// todo(fs): 		ServiceFilter: ".*sql",
// todo(fs): 	}
// todo(fs): 	if !a.shouldProcessUserEvent(p) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Bad tag name
// todo(fs): 	p = &UserEvent{
// todo(fs): 		ServiceFilter: ".*sql",
// todo(fs): 		TagFilter:     "slave",
// todo(fs): 	}
// todo(fs): 	if a.shouldProcessUserEvent(p) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Good service name
// todo(fs): 	p = &UserEvent{
// todo(fs): 		ServiceFilter: ".*sql",
// todo(fs): 		TagFilter:     "master",
// todo(fs): 	}
// todo(fs): 	if !a.shouldProcessUserEvent(p) {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestIngestUserEvent(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	for i := 0; i < 512; i++ {
// todo(fs): 		msg := &UserEvent{LTime: uint64(i), Name: "test"}
// todo(fs): 		a.ingestUserEvent(msg)
// todo(fs): 		if a.LastUserEvent() != msg {
// todo(fs): 			t.Fatalf("bad: %#v", msg)
// todo(fs): 		}
// todo(fs): 		events := a.UserEvents()
// todo(fs):
// todo(fs): 		expectLen := 256
// todo(fs): 		if i < 256 {
// todo(fs): 			expectLen = i + 1
// todo(fs): 		}
// todo(fs): 		if len(events) != expectLen {
// todo(fs): 			t.Fatalf("bad: %d %d %d", i, expectLen, len(events))
// todo(fs): 		}
// todo(fs):
// todo(fs): 		counter := i
// todo(fs): 		for j := len(events) - 1; j >= 0; j-- {
// todo(fs): 			if events[j].LTime != uint64(counter) {
// todo(fs): 				t.Fatalf("bad: %#v", events)
// todo(fs): 			}
// todo(fs): 			counter--
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestFireReceiveEvent(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:      "mysql",
// todo(fs): 		Service: "mysql",
// todo(fs): 		Tags:    []string{"test", "foo", "bar", "master"},
// todo(fs): 		Port:    5000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, "")
// todo(fs):
// todo(fs): 	p1 := &UserEvent{Name: "deploy", ServiceFilter: "web"}
// todo(fs): 	err := a.UserEvent("dc1", "root", p1)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	p2 := &UserEvent{Name: "deploy"}
// todo(fs): 	err = a.UserEvent("dc1", "root", p2)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if got, want := len(a.UserEvents()), 1; got != want {
// todo(fs): 			r.Fatalf("got %d events want %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	last := a.LastUserEvent()
// todo(fs): 	if last.ID != p2.ID {
// todo(fs): 		t.Fatalf("bad: %#v", last)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestUserEventToken(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLDefaultPolicy = "deny" // Set the default policies to deny
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create an ACL token
// todo(fs): 	args := structs.ACLRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Op:         structs.ACLSet,
// todo(fs): 		ACL: structs.ACL{
// todo(fs): 			Name:  "User token",
// todo(fs): 			Type:  structs.ACLTypeClient,
// todo(fs): 			Rules: testEventPolicy,
// todo(fs): 		},
// todo(fs): 		WriteRequest: structs.WriteRequest{Token: "root"},
// todo(fs): 	}
// todo(fs): 	var token string
// todo(fs): 	if err := a.RPC("ACL.Apply", &args, &token); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	type tcase struct {
// todo(fs): 		name   string
// todo(fs): 		expect bool
// todo(fs): 	}
// todo(fs): 	cases := []tcase{
// todo(fs): 		{"foo", false},
// todo(fs): 		{"bar", false},
// todo(fs): 		{"baz", true},
// todo(fs): 		{"zip", false},
// todo(fs): 	}
// todo(fs): 	for _, c := range cases {
// todo(fs): 		event := &UserEvent{Name: c.name}
// todo(fs): 		err := a.UserEvent("dc1", token, event)
// todo(fs): 		allowed := !acl.IsErrPermissionDenied(err)
// todo(fs): 		if allowed != c.expect {
// todo(fs): 			t.Fatalf("bad: %#v result: %v", c, allowed)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
const testEventPolicy = `
event "foo" {
	policy = "deny"
}
event "bar" {
	policy = "read"
}
event "baz" {
	policy = "write"
}
`
