package agent

import (
	"fmt"

	"github.com/hashicorp/consul/agent/structs"
)

// MockPreparedQuery is a fake endpoint that we inject into the Consul server
// in order to observe the RPC calls made by these HTTP endpoints. This lets
// us make sure that the request is being formed properly without having to
// set up a realistic environment for prepared queries, which is a huge task and
// already done in detail inside the prepared query endpoint's unit tests. If we
// can prove this formats proper requests into that then we should be good to
// go. We will do a single set of end-to-end tests in here to make sure that the
// server is wired up to the right endpoint when not "injected".
type MockPreparedQuery struct {
	applyFn   func(*structs.PreparedQueryRequest, *string) error
	getFn     func(*structs.PreparedQuerySpecificRequest, *structs.IndexedPreparedQueries) error
	listFn    func(*structs.DCSpecificRequest, *structs.IndexedPreparedQueries) error
	executeFn func(*structs.PreparedQueryExecuteRequest, *structs.PreparedQueryExecuteResponse) error
	explainFn func(*structs.PreparedQueryExecuteRequest, *structs.PreparedQueryExplainResponse) error
}

func (m *MockPreparedQuery) Apply(args *structs.PreparedQueryRequest,
	reply *string) (err error) {
	if m.applyFn != nil {
		return m.applyFn(args, reply)
	}
	return fmt.Errorf("should not have called Apply")
}

func (m *MockPreparedQuery) Get(args *structs.PreparedQuerySpecificRequest,
	reply *structs.IndexedPreparedQueries) error {
	if m.getFn != nil {
		return m.getFn(args, reply)
	}
	return fmt.Errorf("should not have called Get")
}

func (m *MockPreparedQuery) List(args *structs.DCSpecificRequest,
	reply *structs.IndexedPreparedQueries) error {
	if m.listFn != nil {
		return m.listFn(args, reply)
	}
	return fmt.Errorf("should not have called List")
}

func (m *MockPreparedQuery) Execute(args *structs.PreparedQueryExecuteRequest,
	reply *structs.PreparedQueryExecuteResponse) error {
	if m.executeFn != nil {
		return m.executeFn(args, reply)
	}
	return fmt.Errorf("should not have called Execute")
}

func (m *MockPreparedQuery) Explain(args *structs.PreparedQueryExecuteRequest,
	reply *structs.PreparedQueryExplainResponse) error {
	if m.explainFn != nil {
		return m.explainFn(args, reply)
	}
	return fmt.Errorf("should not have called Explain")
}

// todo(fs): func TestPreparedQuery_Create(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockPreparedQuery{
// todo(fs): 		applyFn: func(args *structs.PreparedQueryRequest, reply *string) error {
// todo(fs): 			expected := &structs.PreparedQueryRequest{
// todo(fs): 				Datacenter: "dc1",
// todo(fs): 				Op:         structs.PreparedQueryCreate,
// todo(fs): 				Query: &structs.PreparedQuery{
// todo(fs): 					Name:    "my-query",
// todo(fs): 					Session: "my-session",
// todo(fs): 					Service: structs.ServiceQuery{
// todo(fs): 						Service: "my-service",
// todo(fs): 						Failover: structs.QueryDatacenterOptions{
// todo(fs): 							NearestN:    4,
// todo(fs): 							Datacenters: []string{"dc1", "dc2"},
// todo(fs): 						},
// todo(fs): 						OnlyPassing: true,
// todo(fs): 						Tags:        []string{"foo", "bar"},
// todo(fs): 						NodeMeta:    map[string]string{"somekey": "somevalue"},
// todo(fs): 					},
// todo(fs): 					DNS: structs.QueryDNSOptions{
// todo(fs): 						TTL: "10s",
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 				WriteRequest: structs.WriteRequest{
// todo(fs): 					Token: "my-token",
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			if !reflect.DeepEqual(args, expected) {
// todo(fs): 				t.Fatalf("bad: %v", args)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			*reply = "my-id"
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"Name":    "my-query",
// todo(fs): 		"Session": "my-session",
// todo(fs): 		"Service": map[string]interface{}{
// todo(fs): 			"Service": "my-service",
// todo(fs): 			"Failover": map[string]interface{}{
// todo(fs): 				"NearestN":    4,
// todo(fs): 				"Datacenters": []string{"dc1", "dc2"},
// todo(fs): 			},
// todo(fs): 			"OnlyPassing": true,
// todo(fs): 			"Tags":        []string{"foo", "bar"},
// todo(fs): 			"NodeMeta":    map[string]string{"somekey": "somevalue"},
// todo(fs): 		},
// todo(fs): 		"DNS": map[string]interface{}{
// todo(fs): 			"TTL": "10s",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := enc.Encode(raw); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("POST", "/v1/query?token=my-token", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.PreparedQueryGeneral(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 	}
// todo(fs): 	r, ok := obj.(preparedQueryCreateResponse)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("unexpected: %T", obj)
// todo(fs): 	}
// todo(fs): 	if r.ID != "my-id" {
// todo(fs): 		t.Fatalf("bad ID: %s", r.ID)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_List(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			listFn: func(args *structs.DCSpecificRequest, reply *structs.IndexedPreparedQueries) error {
// todo(fs): 				// Return an empty response.
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQueryGeneral(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueries)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if r == nil || len(r) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			listFn: func(args *structs.DCSpecificRequest, reply *structs.IndexedPreparedQueries) error {
// todo(fs): 				expected := &structs.DCSpecificRequest{
// todo(fs): 					Datacenter: "dc1",
// todo(fs): 					QueryOptions: structs.QueryOptions{
// todo(fs): 						Token:             "my-token",
// todo(fs): 						RequireConsistent: true,
// todo(fs): 					},
// todo(fs): 				}
// todo(fs): 				if !reflect.DeepEqual(args, expected) {
// todo(fs): 					t.Fatalf("bad: %v", args)
// todo(fs): 				}
// todo(fs):
// todo(fs): 				query := &structs.PreparedQuery{
// todo(fs): 					ID: "my-id",
// todo(fs): 				}
// todo(fs): 				reply.Queries = append(reply.Queries, query)
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query?token=my-token&consistent=true", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQueryGeneral(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueries)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if len(r) != 1 || r[0].ID != "my-id" {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_Execute(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			executeFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExecuteResponse) error {
// todo(fs): 				// Just return an empty response.
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/my-id/execute", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueryExecuteResponse)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if r.Nodes == nil || len(r.Nodes) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			executeFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExecuteResponse) error {
// todo(fs): 				expected := &structs.PreparedQueryExecuteRequest{
// todo(fs): 					Datacenter:    "dc1",
// todo(fs): 					QueryIDOrName: "my-id",
// todo(fs): 					Limit:         5,
// todo(fs): 					Source: structs.QuerySource{
// todo(fs): 						Datacenter: "dc1",
// todo(fs): 						Node:       "my-node",
// todo(fs): 					},
// todo(fs): 					Agent: structs.QuerySource{
// todo(fs): 						Datacenter: a.Config.Datacenter,
// todo(fs): 						Node:       a.Config.NodeName,
// todo(fs): 					},
// todo(fs): 					QueryOptions: structs.QueryOptions{
// todo(fs): 						Token:             "my-token",
// todo(fs): 						RequireConsistent: true,
// todo(fs): 					},
// todo(fs): 				}
// todo(fs): 				if !reflect.DeepEqual(args, expected) {
// todo(fs): 					t.Fatalf("bad: %v", args)
// todo(fs): 				}
// todo(fs):
// todo(fs): 				// Just set something so we can tell this is returned.
// todo(fs): 				reply.Failovers = 99
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/my-id/execute?token=my-token&consistent=true&near=my-node&limit=5", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueryExecuteResponse)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if r.Failovers != 99 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Ensure the proper params are set when no special args are passed
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			executeFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExecuteResponse) error {
// todo(fs): 				if args.Source.Node != "" {
// todo(fs): 					t.Fatalf("expect node to be empty, got %q", args.Source.Node)
// todo(fs): 				}
// todo(fs): 				expect := structs.QuerySource{
// todo(fs): 					Datacenter: a.Config.Datacenter,
// todo(fs): 					Node:       a.Config.NodeName,
// todo(fs): 				}
// todo(fs): 				if !reflect.DeepEqual(args.Agent, expect) {
// todo(fs): 					t.Fatalf("expect: %#v\nactual: %#v", expect, args.Agent)
// todo(fs): 				}
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/my-id/execute", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Ensure WAN translation occurs for a response outside of the local DC.
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		cfg := TestConfig()
// todo(fs): 		cfg.Datacenter = "dc1"
// todo(fs): 		cfg.TranslateWANAddrs = true
// todo(fs): 		a := NewTestAgent(t.Name(), cfg)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			executeFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExecuteResponse) error {
// todo(fs): 				nodesResponse := make(structs.CheckServiceNodes, 1)
// todo(fs): 				nodesResponse[0].Node = &structs.Node{
// todo(fs): 					Node: "foo", Address: "127.0.0.1",
// todo(fs): 					TaggedAddresses: map[string]string{
// todo(fs): 						"wan": "127.0.0.2",
// todo(fs): 					},
// todo(fs): 				}
// todo(fs): 				reply.Nodes = nodesResponse
// todo(fs): 				reply.Datacenter = "dc2"
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/my-id/execute?dc=dc2", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueryExecuteResponse)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if r.Nodes == nil || len(r.Nodes) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		node := r.Nodes[0]
// todo(fs): 		if node.Node.Address != "127.0.0.2" {
// todo(fs): 			t.Fatalf("bad: %v", node.Node)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Ensure WAN translation doesn't occur for the local DC.
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		cfg := TestConfig()
// todo(fs): 		cfg.Datacenter = "dc1"
// todo(fs): 		cfg.TranslateWANAddrs = true
// todo(fs): 		a := NewTestAgent(t.Name(), cfg)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			executeFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExecuteResponse) error {
// todo(fs): 				nodesResponse := make(structs.CheckServiceNodes, 1)
// todo(fs): 				nodesResponse[0].Node = &structs.Node{
// todo(fs): 					Node: "foo", Address: "127.0.0.1",
// todo(fs): 					TaggedAddresses: map[string]string{
// todo(fs): 						"wan": "127.0.0.2",
// todo(fs): 					},
// todo(fs): 				}
// todo(fs): 				reply.Nodes = nodesResponse
// todo(fs): 				reply.Datacenter = "dc1"
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/my-id/execute?dc=dc2", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueryExecuteResponse)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if r.Nodes == nil || len(r.Nodes) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		node := r.Nodes[0]
// todo(fs): 		if node.Node.Address != "127.0.0.1" {
// todo(fs): 			t.Fatalf("bad: %v", node.Node)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/not-there/execute", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 404 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_Explain(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			explainFn: func(args *structs.PreparedQueryExecuteRequest, reply *structs.PreparedQueryExplainResponse) error {
// todo(fs): 				expected := &structs.PreparedQueryExecuteRequest{
// todo(fs): 					Datacenter:    "dc1",
// todo(fs): 					QueryIDOrName: "my-id",
// todo(fs): 					Limit:         5,
// todo(fs): 					Source: structs.QuerySource{
// todo(fs): 						Datacenter: "dc1",
// todo(fs): 						Node:       "my-node",
// todo(fs): 					},
// todo(fs): 					Agent: structs.QuerySource{
// todo(fs): 						Datacenter: a.Config.Datacenter,
// todo(fs): 						Node:       a.Config.NodeName,
// todo(fs): 					},
// todo(fs): 					QueryOptions: structs.QueryOptions{
// todo(fs): 						Token:             "my-token",
// todo(fs): 						RequireConsistent: true,
// todo(fs): 					},
// todo(fs): 				}
// todo(fs): 				if !reflect.DeepEqual(args, expected) {
// todo(fs): 					t.Fatalf("bad: %v", args)
// todo(fs): 				}
// todo(fs):
// todo(fs): 				// Just set something so we can tell this is returned.
// todo(fs): 				reply.Query.Name = "hello"
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/my-id/explain?token=my-token&consistent=true&near=my-node&limit=5", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueryExplainResponse)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if r.Query.Name != "hello" {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/not-there/explain", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 404 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_Get(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		m := MockPreparedQuery{
// todo(fs): 			getFn: func(args *structs.PreparedQuerySpecificRequest, reply *structs.IndexedPreparedQueries) error {
// todo(fs): 				expected := &structs.PreparedQuerySpecificRequest{
// todo(fs): 					Datacenter: "dc1",
// todo(fs): 					QueryID:    "my-id",
// todo(fs): 					QueryOptions: structs.QueryOptions{
// todo(fs): 						Token:             "my-token",
// todo(fs): 						RequireConsistent: true,
// todo(fs): 					},
// todo(fs): 				}
// todo(fs): 				if !reflect.DeepEqual(args, expected) {
// todo(fs): 					t.Fatalf("bad: %v", args)
// todo(fs): 				}
// todo(fs):
// todo(fs): 				query := &structs.PreparedQuery{
// todo(fs): 					ID: "my-id",
// todo(fs): 				}
// todo(fs): 				reply.Queries = append(reply.Queries, query)
// todo(fs): 				return nil
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/my-id?token=my-token&consistent=true", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueries)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if len(r) != 1 || r[0].ID != "my-id" {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/f004177f-2c28-83b7-4229-eacc25fe55d1", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 404 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_Update(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockPreparedQuery{
// todo(fs): 		applyFn: func(args *structs.PreparedQueryRequest, reply *string) error {
// todo(fs): 			expected := &structs.PreparedQueryRequest{
// todo(fs): 				Datacenter: "dc1",
// todo(fs): 				Op:         structs.PreparedQueryUpdate,
// todo(fs): 				Query: &structs.PreparedQuery{
// todo(fs): 					ID:      "my-id",
// todo(fs): 					Name:    "my-query",
// todo(fs): 					Session: "my-session",
// todo(fs): 					Service: structs.ServiceQuery{
// todo(fs): 						Service: "my-service",
// todo(fs): 						Failover: structs.QueryDatacenterOptions{
// todo(fs): 							NearestN:    4,
// todo(fs): 							Datacenters: []string{"dc1", "dc2"},
// todo(fs): 						},
// todo(fs): 						OnlyPassing: true,
// todo(fs): 						Tags:        []string{"foo", "bar"},
// todo(fs): 						NodeMeta:    map[string]string{"somekey": "somevalue"},
// todo(fs): 					},
// todo(fs): 					DNS: structs.QueryDNSOptions{
// todo(fs): 						TTL: "10s",
// todo(fs): 					},
// todo(fs): 				},
// todo(fs): 				WriteRequest: structs.WriteRequest{
// todo(fs): 					Token: "my-token",
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			if !reflect.DeepEqual(args, expected) {
// todo(fs): 				t.Fatalf("bad: %v", args)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			*reply = "don't care"
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"ID":      "this should get ignored",
// todo(fs): 		"Name":    "my-query",
// todo(fs): 		"Session": "my-session",
// todo(fs): 		"Service": map[string]interface{}{
// todo(fs): 			"Service": "my-service",
// todo(fs): 			"Failover": map[string]interface{}{
// todo(fs): 				"NearestN":    4,
// todo(fs): 				"Datacenters": []string{"dc1", "dc2"},
// todo(fs): 			},
// todo(fs): 			"OnlyPassing": true,
// todo(fs): 			"Tags":        []string{"foo", "bar"},
// todo(fs): 			"NodeMeta":    map[string]string{"somekey": "somevalue"},
// todo(fs): 		},
// todo(fs): 		"DNS": map[string]interface{}{
// todo(fs): 			"TTL": "10s",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := enc.Encode(raw); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/query/my-id?token=my-token", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_Delete(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	m := MockPreparedQuery{
// todo(fs): 		applyFn: func(args *structs.PreparedQueryRequest, reply *string) error {
// todo(fs): 			expected := &structs.PreparedQueryRequest{
// todo(fs): 				Datacenter: "dc1",
// todo(fs): 				Op:         structs.PreparedQueryDelete,
// todo(fs): 				Query: &structs.PreparedQuery{
// todo(fs): 					ID: "my-id",
// todo(fs): 				},
// todo(fs): 				WriteRequest: structs.WriteRequest{
// todo(fs): 					Token: "my-token",
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			if !reflect.DeepEqual(args, expected) {
// todo(fs): 				t.Fatalf("bad: %v", args)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			*reply = "don't care"
// todo(fs): 			return nil
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.registerEndpoint("PreparedQuery", &m); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"ID": "this should get ignored",
// todo(fs): 	}
// todo(fs): 	if err := enc.Encode(raw); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("DELETE", "/v1/query/my-id?token=my-token", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_BadMethods(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("DELETE", "/v1/query", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQueryGeneral(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 405 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("POST", "/v1/query/my-id", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 405 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestPreparedQuery_parseLimit(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/query", body)
// todo(fs): 	limit := 99
// todo(fs): 	if err := parseLimit(req, &limit); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if limit != 0 {
// todo(fs): 		t.Fatalf("bad limit: %d", limit)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/query?limit=11", body)
// todo(fs): 	if err := parseLimit(req, &limit); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if limit != 11 {
// todo(fs): 		t.Fatalf("bad limit: %d", limit)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/query?limit=bob", body)
// todo(fs): 	if err := parseLimit(req, &limit); err == nil {
// todo(fs): 		t.Fatalf("bad: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): // Since we've done exhaustive testing of the calls into the endpoints above
// todo(fs): // this is just a basic end-to-end sanity check to make sure things are wired
// todo(fs): // correctly when calling through to the real endpoints.
// todo(fs): func TestPreparedQuery_Integration(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node and a service.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "my-service",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a query.
// todo(fs): 	var id string
// todo(fs): 	{
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		enc := json.NewEncoder(body)
// todo(fs): 		raw := map[string]interface{}{
// todo(fs): 			"Name": "my-query",
// todo(fs): 			"Service": map[string]interface{}{
// todo(fs): 				"Service": "my-service",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := enc.Encode(raw); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("POST", "/v1/query", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQueryGeneral(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(preparedQueryCreateResponse)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		id = r.ID
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// List them all.
// todo(fs): 	{
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query?token=root", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQueryGeneral(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueries)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if len(r) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Execute it.
// todo(fs): 	{
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/"+id+"/execute", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueryExecuteResponse)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if len(r.Nodes) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Read it back.
// todo(fs): 	{
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/query/"+id, body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.PreparedQuerySpecific(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		r, ok := obj.(structs.PreparedQueries)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if len(r) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", r)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make an update to it.
// todo(fs): 	{
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		enc := json.NewEncoder(body)
// todo(fs): 		raw := map[string]interface{}{
// todo(fs): 			"Name": "my-query",
// todo(fs): 			"Service": map[string]interface{}{
// todo(fs): 				"Service":     "my-service",
// todo(fs): 				"OnlyPassing": true,
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		if err := enc.Encode(raw); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/query/"+id, body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Delete it.
// todo(fs): 	{
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("DELETE", "/v1/query/"+id, body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.PreparedQuerySpecific(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
