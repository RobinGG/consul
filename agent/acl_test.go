package agent

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/testutil"
)

func TestACL_Bad_Config(t *testing.T) {
	t.Parallel()
	cfg := TestConfig()
	cfg.ACLDownPolicy = "nope"
	cfg.DataDir = testutil.TempDir(t, "agent")

	// do not use TestAgent here since we want
	// the agent to fail during startup.
	_, err := New(cfg)
	if err == nil || !strings.Contains(err.Error(), "invalid ACL down policy") {
		t.Fatalf("err: %v", err)
	}
}

type MockServer struct {
	getPolicyFn func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error
}

func (m *MockServer) GetPolicy(args *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
	if m.getPolicyFn != nil {
		return m.getPolicyFn(args, reply)
	}
	return fmt.Errorf("should not have called GetPolicy")
}

// todo(fs): func TestACL_Version8(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = false
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{
// todo(fs): 		// With version 8 enforcement off, this should not get called.
// todo(fs): 		getPolicyFn: func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 			t.Fatalf("should not have called to server")
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if token, err := a.resolveToken("nope"); token != nil || err != nil {
// todo(fs): 		t.Fatalf("bad: %v err: %v", token, err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Disabled(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLDisabledTTL = 10 * time.Millisecond
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{
// todo(fs): 		// Fetch a token without ACLs enabled and make sure the manager sees it.
// todo(fs): 		getPolicyFn: func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 			return rawacl.ErrDisabled
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if a.acls.isDisabled() {
// todo(fs): 		t.Fatalf("should not be disabled yet")
// todo(fs): 	}
// todo(fs): 	if token, err := a.resolveToken("nope"); token != nil || err != nil {
// todo(fs): 		t.Fatalf("bad: %v err: %v", token, err)
// todo(fs): 	}
// todo(fs): 	if !a.acls.isDisabled() {
// todo(fs): 		t.Fatalf("should be disabled")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Now turn on ACLs and check right away, it should still think ACLs are
// todo(fs): 	// disabled since we don't check again right away.
// todo(fs): 	m.getPolicyFn = func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 		return rawacl.ErrNotFound
// todo(fs): 	}
// todo(fs): 	if token, err := a.resolveToken("nope"); token != nil || err != nil {
// todo(fs): 		t.Fatalf("bad: %v err: %v", token, err)
// todo(fs): 	}
// todo(fs): 	if !a.acls.isDisabled() {
// todo(fs): 		t.Fatalf("should be disabled")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Wait the waiting period and make sure it checks again. Do a few tries
// todo(fs): 	// to make sure we don't think it's disabled.
// todo(fs): 	time.Sleep(2 * cfg.ACLDisabledTTL)
// todo(fs): 	for i := 0; i < 10; i++ {
// todo(fs): 		_, err := a.resolveToken("nope")
// todo(fs): 		if !rawacl.IsErrNotFound(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if a.acls.isDisabled() {
// todo(fs): 			t.Fatalf("should not be disabled")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Special_IDs(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs): 	cfg.ACLAgentMasterToken = "towel"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{
// todo(fs): 		// An empty ID should get mapped to the anonymous token.
// todo(fs): 		getPolicyFn: func(req *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
// todo(fs): 			if req.ACL != "anonymous" {
// todo(fs): 				t.Fatalf("bad: %#v", *req)
// todo(fs): 			}
// todo(fs): 			return rawacl.ErrNotFound
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	_, err := a.resolveToken("")
// todo(fs): 	if !rawacl.IsErrNotFound(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// A root ACL request should get rejected and not call the server.
// todo(fs): 	m.getPolicyFn = func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 		t.Fatalf("should not have called to server")
// todo(fs): 		return nil
// todo(fs): 	}
// todo(fs): 	_, err = a.resolveToken("deny")
// todo(fs): 	if !rawacl.IsErrRootDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// The ACL master token should also not call the server, but should give
// todo(fs): 	// us a working agent token.
// todo(fs): 	acl, err := a.resolveToken("towel")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if acl == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !acl.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should be able to read agent")
// todo(fs): 	}
// todo(fs): 	if !acl.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should be able to write agent")
// todo(fs): 	}
// todo(fs): 	if !acl.NodeRead("hello") {
// todo(fs): 		t.Fatalf("should be able to read any node")
// todo(fs): 	}
// todo(fs): 	if acl.NodeWrite("hello") {
// todo(fs): 		t.Fatalf("should not be able to write any node")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Down_Deny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLDownPolicy = "deny"
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{
// todo(fs): 		// Resolve with ACLs down.
// todo(fs): 		getPolicyFn: func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 			return fmt.Errorf("ACLs are broken")
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	acl, err := a.resolveToken("nope")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if acl == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if acl.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Down_Allow(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLDownPolicy = "allow"
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{
// todo(fs): 		// Resolve with ACLs down.
// todo(fs): 		getPolicyFn: func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 			return fmt.Errorf("ACLs are broken")
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	acl, err := a.resolveToken("nope")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if acl == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !acl.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Down_Extend(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLDownPolicy = "extend-cache"
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{
// todo(fs): 		// Populate the cache for one of the tokens.
// todo(fs): 		getPolicyFn: func(req *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
// todo(fs): 			*reply = structs.ACLPolicy{
// todo(fs): 				Parent: "allow",
// todo(fs): 				Policy: &rawacl.Policy{
// todo(fs): 					Agents: []*rawacl.AgentPolicy{
// todo(fs): 						&rawacl.AgentPolicy{
// todo(fs): 							Node:   cfg.NodeName,
// todo(fs): 							Policy: "read",
// todo(fs): 						},
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	acl, err := a.resolveToken("yep")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if acl == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !acl.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if acl.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Now take down ACLs and make sure a new token fails to resolve.
// todo(fs): 	m.getPolicyFn = func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 		return fmt.Errorf("ACLs are broken")
// todo(fs): 	}
// todo(fs): 	acl, err = a.resolveToken("nope")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if acl == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if acl.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs): 	if acl.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Read the token from the cache while ACLs are broken, which should
// todo(fs): 	// extend.
// todo(fs): 	acl, err = a.resolveToken("yep")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if acl == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !acl.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if acl.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_Cache(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{
// todo(fs): 		// Populate the cache for one of the tokens.
// todo(fs): 		getPolicyFn: func(req *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
// todo(fs): 			*reply = structs.ACLPolicy{
// todo(fs): 				ETag:   "hash1",
// todo(fs): 				Parent: "deny",
// todo(fs): 				Policy: &rawacl.Policy{
// todo(fs): 					Agents: []*rawacl.AgentPolicy{
// todo(fs): 						&rawacl.AgentPolicy{
// todo(fs): 							Node:   cfg.NodeName,
// todo(fs): 							Policy: "read",
// todo(fs): 						},
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 				TTL: 10 * time.Millisecond,
// todo(fs): 			}
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	rule, err := a.resolveToken("yep")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if rule == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !rule.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if rule.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs): 	if rule.NodeRead("nope") {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Fetch right away and make sure it uses the cache.
// todo(fs): 	m.getPolicyFn = func(*structs.ACLPolicyRequest, *structs.ACLPolicy) error {
// todo(fs): 		t.Fatalf("should not have called to server")
// todo(fs): 		return nil
// todo(fs): 	}
// todo(fs): 	rule, err = a.resolveToken("yep")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if rule == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !rule.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if rule.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs): 	if rule.NodeRead("nope") {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Wait for the TTL to expire and try again. This time the token will be
// todo(fs): 	// gone.
// todo(fs): 	time.Sleep(20 * time.Millisecond)
// todo(fs): 	m.getPolicyFn = func(req *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
// todo(fs): 		return rawacl.ErrNotFound
// todo(fs): 	}
// todo(fs): 	_, err = a.resolveToken("yep")
// todo(fs): 	if !rawacl.IsErrNotFound(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Page it back in with a new tag and different policy
// todo(fs): 	m.getPolicyFn = func(req *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
// todo(fs): 		*reply = structs.ACLPolicy{
// todo(fs): 			ETag:   "hash2",
// todo(fs): 			Parent: "deny",
// todo(fs): 			Policy: &rawacl.Policy{
// todo(fs): 				Agents: []*rawacl.AgentPolicy{
// todo(fs): 					&rawacl.AgentPolicy{
// todo(fs): 						Node:   cfg.NodeName,
// todo(fs): 						Policy: "write",
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 			TTL: 10 * time.Millisecond,
// todo(fs): 		}
// todo(fs): 		return nil
// todo(fs): 	}
// todo(fs): 	rule, err = a.resolveToken("yep")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if rule == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !rule.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if !rule.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if rule.NodeRead("nope") {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Wait for the TTL to expire and try again. This will match the tag
// todo(fs): 	// and not send the policy back, but we should have the old token
// todo(fs): 	// behavior.
// todo(fs): 	time.Sleep(20 * time.Millisecond)
// todo(fs): 	var didRefresh bool
// todo(fs): 	m.getPolicyFn = func(req *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
// todo(fs): 		*reply = structs.ACLPolicy{
// todo(fs): 			ETag: "hash2",
// todo(fs): 			TTL:  10 * time.Millisecond,
// todo(fs): 		}
// todo(fs): 		didRefresh = true
// todo(fs): 		return nil
// todo(fs): 	}
// todo(fs): 	rule, err = a.resolveToken("yep")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if rule == nil {
// todo(fs): 		t.Fatalf("should not be nil")
// todo(fs): 	}
// todo(fs): 	if !rule.AgentRead(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if !rule.AgentWrite(cfg.NodeName) {
// todo(fs): 		t.Fatalf("should allow")
// todo(fs): 	}
// todo(fs): 	if rule.NodeRead("nope") {
// todo(fs): 		t.Fatalf("should deny")
// todo(fs): 	}
// todo(fs): 	if !didRefresh {
// todo(fs): 		t.Fatalf("should refresh")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): // catalogPolicy supplies some standard policies to help with testing the
// todo(fs): // catalog-related vet and filter functions.
// todo(fs): func catalogPolicy(req *structs.ACLPolicyRequest, reply *structs.ACLPolicy) error {
// todo(fs): 	reply.Policy = &rawacl.Policy{}
// todo(fs):
// todo(fs): 	switch req.ACL {
// todo(fs):
// todo(fs): 	case "node-ro":
// todo(fs): 		reply.Policy.Nodes = append(reply.Policy.Nodes,
// todo(fs): 			&rawacl.NodePolicy{Name: "Node", Policy: "read"})
// todo(fs):
// todo(fs): 	case "node-rw":
// todo(fs): 		reply.Policy.Nodes = append(reply.Policy.Nodes,
// todo(fs): 			&rawacl.NodePolicy{Name: "Node", Policy: "write"})
// todo(fs):
// todo(fs): 	case "service-ro":
// todo(fs): 		reply.Policy.Services = append(reply.Policy.Services,
// todo(fs): 			&rawacl.ServicePolicy{Name: "service", Policy: "read"})
// todo(fs):
// todo(fs): 	case "service-rw":
// todo(fs): 		reply.Policy.Services = append(reply.Policy.Services,
// todo(fs): 			&rawacl.ServicePolicy{Name: "service", Policy: "write"})
// todo(fs):
// todo(fs): 	case "other-rw":
// todo(fs): 		reply.Policy.Services = append(reply.Policy.Services,
// todo(fs): 			&rawacl.ServicePolicy{Name: "other", Policy: "write"})
// todo(fs):
// todo(fs): 	default:
// todo(fs): 		return fmt.Errorf("unknown token %q", req.ACL)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	return nil
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_vetServiceRegister(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{catalogPolicy}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register a new service, with permission.
// todo(fs): 	err := a.vetServiceRegister("service-rw", &structs.NodeService{
// todo(fs): 		ID:      "my-service",
// todo(fs): 		Service: "service",
// todo(fs): 	})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register a new service without write privs.
// todo(fs): 	err = a.vetServiceRegister("service-ro", &structs.NodeService{
// todo(fs): 		ID:      "my-service",
// todo(fs): 		Service: "service",
// todo(fs): 	})
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try to register over a service without write privs to the existing
// todo(fs): 	// service.
// todo(fs): 	a.state.AddService(&structs.NodeService{
// todo(fs): 		ID:      "my-service",
// todo(fs): 		Service: "other",
// todo(fs): 	}, "")
// todo(fs): 	err = a.vetServiceRegister("service-rw", &structs.NodeService{
// todo(fs): 		ID:      "my-service",
// todo(fs): 		Service: "service",
// todo(fs): 	})
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_vetServiceUpdate(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{catalogPolicy}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update a service that doesn't exist.
// todo(fs): 	err := a.vetServiceUpdate("service-rw", "my-service")
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "Unknown service") {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update with write privs.
// todo(fs): 	a.state.AddService(&structs.NodeService{
// todo(fs): 		ID:      "my-service",
// todo(fs): 		Service: "service",
// todo(fs): 	}, "")
// todo(fs): 	err = a.vetServiceUpdate("service-rw", "my-service")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update without write privs.
// todo(fs): 	err = a.vetServiceUpdate("service-ro", "my-service")
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_vetCheckRegister(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{catalogPolicy}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register a new service check with write privs.
// todo(fs): 	err := a.vetCheckRegister("service-rw", &structs.HealthCheck{
// todo(fs): 		CheckID:     types.CheckID("my-check"),
// todo(fs): 		ServiceID:   "my-service",
// todo(fs): 		ServiceName: "service",
// todo(fs): 	})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register a new service check without write privs.
// todo(fs): 	err = a.vetCheckRegister("service-ro", &structs.HealthCheck{
// todo(fs): 		CheckID:     types.CheckID("my-check"),
// todo(fs): 		ServiceID:   "my-service",
// todo(fs): 		ServiceName: "service",
// todo(fs): 	})
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register a new node check with write privs.
// todo(fs): 	err = a.vetCheckRegister("node-rw", &structs.HealthCheck{
// todo(fs): 		CheckID: types.CheckID("my-check"),
// todo(fs): 	})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register a new node check without write privs.
// todo(fs): 	err = a.vetCheckRegister("node-ro", &structs.HealthCheck{
// todo(fs): 		CheckID: types.CheckID("my-check"),
// todo(fs): 	})
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try to register over a service check without write privs to the
// todo(fs): 	// existing service.
// todo(fs): 	a.state.AddService(&structs.NodeService{
// todo(fs): 		ID:      "my-service",
// todo(fs): 		Service: "service",
// todo(fs): 	}, "")
// todo(fs): 	a.state.AddCheck(&structs.HealthCheck{
// todo(fs): 		CheckID:     types.CheckID("my-check"),
// todo(fs): 		ServiceID:   "my-service",
// todo(fs): 		ServiceName: "other",
// todo(fs): 	}, "")
// todo(fs): 	err = a.vetCheckRegister("service-rw", &structs.HealthCheck{
// todo(fs): 		CheckID:     types.CheckID("my-check"),
// todo(fs): 		ServiceID:   "my-service",
// todo(fs): 		ServiceName: "service",
// todo(fs): 	})
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try to register over a node check without write privs to the node.
// todo(fs): 	a.state.AddCheck(&structs.HealthCheck{
// todo(fs): 		CheckID: types.CheckID("my-node-check"),
// todo(fs): 	}, "")
// todo(fs): 	err = a.vetCheckRegister("service-rw", &structs.HealthCheck{
// todo(fs): 		CheckID:     types.CheckID("my-node-check"),
// todo(fs): 		ServiceID:   "my-service",
// todo(fs): 		ServiceName: "service",
// todo(fs): 	})
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_vetCheckUpdate(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{catalogPolicy}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update a check that doesn't exist.
// todo(fs): 	err := a.vetCheckUpdate("node-rw", "my-check")
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "Unknown check") {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update service check with write privs.
// todo(fs): 	a.state.AddService(&structs.NodeService{
// todo(fs): 		ID:      "my-service",
// todo(fs): 		Service: "service",
// todo(fs): 	}, "")
// todo(fs): 	a.state.AddCheck(&structs.HealthCheck{
// todo(fs): 		CheckID:     types.CheckID("my-service-check"),
// todo(fs): 		ServiceID:   "my-service",
// todo(fs): 		ServiceName: "service",
// todo(fs): 	}, "")
// todo(fs): 	err = a.vetCheckUpdate("service-rw", "my-service-check")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update service check without write privs.
// todo(fs): 	err = a.vetCheckUpdate("service-ro", "my-service-check")
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update node check with write privs.
// todo(fs): 	a.state.AddCheck(&structs.HealthCheck{
// todo(fs): 		CheckID: types.CheckID("my-node-check"),
// todo(fs): 	}, "")
// todo(fs): 	err = a.vetCheckUpdate("node-rw", "my-node-check")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Update without write privs.
// todo(fs): 	err = a.vetCheckUpdate("node-ro", "my-node-check")
// todo(fs): 	if !rawacl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_filterMembers(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{catalogPolicy}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var members []serf.Member
// todo(fs): 	if err := a.filterMembers("node-ro", &members); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if len(members) != 0 {
// todo(fs): 		t.Fatalf("bad: %#v", members)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	members = []serf.Member{
// todo(fs): 		serf.Member{Name: "Node 1"},
// todo(fs): 		serf.Member{Name: "Nope"},
// todo(fs): 		serf.Member{Name: "Node 2"},
// todo(fs): 	}
// todo(fs): 	if err := a.filterMembers("node-ro", &members); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if len(members) != 2 ||
// todo(fs): 		members[0].Name != "Node 1" ||
// todo(fs): 		members[1].Name != "Node 2" {
// todo(fs): 		t.Fatalf("bad: %#v", members)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_filterServices(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{catalogPolicy}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	services := make(map[string]*structs.NodeService)
// todo(fs): 	if err := a.filterServices("node-ro", &services); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	services["my-service"] = &structs.NodeService{ID: "my-service", Service: "service"}
// todo(fs): 	services["my-other"] = &structs.NodeService{ID: "my-other", Service: "other"}
// todo(fs): 	if err := a.filterServices("service-ro", &services); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, ok := services["my-service"]; !ok {
// todo(fs): 		t.Fatalf("bad: %#v", services)
// todo(fs): 	}
// todo(fs): 	if _, ok := services["my-other"]; ok {
// todo(fs): 		t.Fatalf("bad: %#v", services)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestACL_filterChecks(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockServer{catalogPolicy}
// todo(fs): 	if err := a.registerEndpoint("ACL", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	checks := make(map[types.CheckID]*structs.HealthCheck)
// todo(fs): 	if err := a.filterChecks("node-ro", &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	checks["my-node"] = &structs.HealthCheck{}
// todo(fs): 	checks["my-service"] = &structs.HealthCheck{ServiceName: "service"}
// todo(fs): 	checks["my-other"] = &structs.HealthCheck{ServiceName: "other"}
// todo(fs): 	if err := a.filterChecks("service-ro", &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, ok := checks["my-node"]; ok {
// todo(fs): 		t.Fatalf("bad: %#v", checks)
// todo(fs): 	}
// todo(fs): 	if _, ok := checks["my-service"]; !ok {
// todo(fs): 		t.Fatalf("bad: %#v", checks)
// todo(fs): 	}
// todo(fs): 	if _, ok := checks["my-other"]; ok {
// todo(fs): 		t.Fatalf("bad: %#v", checks)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	checks["my-node"] = &structs.HealthCheck{}
// todo(fs): 	checks["my-service"] = &structs.HealthCheck{ServiceName: "service"}
// todo(fs): 	checks["my-other"] = &structs.HealthCheck{ServiceName: "other"}
// todo(fs): 	if err := a.filterChecks("node-ro", &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, ok := checks["my-node"]; !ok {
// todo(fs): 		t.Fatalf("bad: %#v", checks)
// todo(fs): 	}
// todo(fs): 	if _, ok := checks["my-service"]; ok {
// todo(fs): 		t.Fatalf("bad: %#v", checks)
// todo(fs): 	}
// todo(fs): 	if _, ok := checks["my-other"]; ok {
// todo(fs): 		t.Fatalf("bad: %#v", checks)
// todo(fs): 	}
// todo(fs): }
