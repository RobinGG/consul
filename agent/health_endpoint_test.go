package agent

// todo(fs): func TestHealthChecksInState(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("warning", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/health/state/warning?dc=dc1", nil)
// todo(fs): 		retry.Run(t, func(r *retry.R) {
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.HealthChecksInState(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				r.Fatal(err)
// todo(fs): 			}
// todo(fs): 			if err := checkIndex(resp); err != nil {
// todo(fs): 				r.Fatal(err)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			// Should be a non-nil empty list
// todo(fs): 			nodes := obj.(structs.HealthChecks)
// todo(fs): 			if nodes == nil || len(nodes) != 0 {
// todo(fs): 				r.Fatalf("bad: %v", obj)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("passing", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/health/state/passing?dc=dc1", nil)
// todo(fs): 		retry.Run(t, func(r *retry.R) {
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.HealthChecksInState(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				r.Fatal(err)
// todo(fs): 			}
// todo(fs): 			if err := checkIndex(resp); err != nil {
// todo(fs): 				r.Fatal(err)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			// Should be 1 health check for the server
// todo(fs): 			nodes := obj.(structs.HealthChecks)
// todo(fs): 			if len(nodes) != 1 {
// todo(fs): 				r.Fatalf("bad: %v", obj)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthChecksInState_NodeMetaFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		NodeMeta:   map[string]string{"somekey": "somevalue"},
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			Node:   "bar",
// todo(fs): 			Name:   "node check",
// todo(fs): 			Status: api.HealthCritical,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/state/critical?node-meta=somekey:somevalue", nil)
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.HealthChecksInState(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs): 		if err := checkIndex(resp); err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Should be 1 health check for the server
// todo(fs): 		nodes := obj.(structs.HealthChecks)
// todo(fs): 		if len(nodes) != 1 {
// todo(fs): 			r.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthChecksInState_DistanceSort(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			Node:   "bar",
// todo(fs): 			Name:   "node check",
// todo(fs): 			Status: api.HealthCritical,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args.Node, args.Check.Node = "foo", "foo"
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/state/critical?dc=dc1&near=foo", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthChecksInState(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs): 	nodes := obj.(structs.HealthChecks)
// todo(fs): 	if len(nodes) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[0].Node != "bar" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[1].Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Send an update for the node and wait for it to get applied.
// todo(fs): 	arg := structs.CoordinateUpdateRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Coord:      coordinate.NewCoordinate(coordinate.DefaultConfig()),
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Coordinate.Update", &arg, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	// Retry until foo moves to the front of the line.
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		resp = httptest.NewRecorder()
// todo(fs): 		obj, err = a.srv.HealthChecksInState(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		assertIndex(t, resp)
// todo(fs): 		nodes = obj.(structs.HealthChecks)
// todo(fs): 		if len(nodes) != 2 {
// todo(fs): 			r.Fatalf("bad: %v", nodes)
// todo(fs): 		}
// todo(fs): 		if nodes[0].Node != "foo" {
// todo(fs): 			r.Fatalf("bad: %v", nodes)
// todo(fs): 		}
// todo(fs): 		if nodes[1].Node != "bar" {
// todo(fs): 			r.Fatalf("bad: %v", nodes)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthNodeChecks(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/node/nope?dc=dc1", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthNodeChecks(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be a non-nil empty list
// todo(fs): 	nodes := obj.(structs.HealthChecks)
// todo(fs): 	if nodes == nil || len(nodes) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/health/node/%s?dc=dc1", a.Config.NodeName), nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.HealthNodeChecks(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be 1 health check for the server
// todo(fs): 	nodes = obj.(structs.HealthChecks)
// todo(fs): 	if len(nodes) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceChecks(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/checks/consul?dc=dc1", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthServiceChecks(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be a non-nil empty list
// todo(fs): 	nodes := obj.(structs.HealthChecks)
// todo(fs): 	if nodes == nil || len(nodes) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a service check
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			Name:      "consul check",
// todo(fs): 			ServiceID: "consul",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err = a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/health/checks/consul?dc=dc1", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.HealthServiceChecks(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be 1 health check for consul
// todo(fs): 	nodes = obj.(structs.HealthChecks)
// todo(fs): 	if len(nodes) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceChecks_NodeMetaFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/checks/consul?dc=dc1&node-meta=somekey:somevalue", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthServiceChecks(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be a non-nil empty list
// todo(fs): 	nodes := obj.(structs.HealthChecks)
// todo(fs): 	if nodes == nil || len(nodes) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a service check
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		NodeMeta:   map[string]string{"somekey": "somevalue"},
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			Name:      "consul check",
// todo(fs): 			ServiceID: "consul",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err = a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/health/checks/consul?dc=dc1&node-meta=somekey:somevalue", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.HealthServiceChecks(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be 1 health check for consul
// todo(fs): 	nodes = obj.(structs.HealthChecks)
// todo(fs): 	if len(nodes) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceChecks_DistanceSort(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create a service check
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			ID:      "test",
// todo(fs): 			Service: "test",
// todo(fs): 		},
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			Node:      "bar",
// todo(fs): 			Name:      "test check",
// todo(fs): 			ServiceID: "test",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args.Node, args.Check.Node = "foo", "foo"
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/checks/test?dc=dc1&near=foo", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthServiceChecks(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs): 	nodes := obj.(structs.HealthChecks)
// todo(fs): 	if len(nodes) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): 	if nodes[0].Node != "bar" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[1].Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Send an update for the node and wait for it to get applied.
// todo(fs): 	arg := structs.CoordinateUpdateRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Coord:      coordinate.NewCoordinate(coordinate.DefaultConfig()),
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Coordinate.Update", &arg, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	// Retry until foo has moved to the front of the line.
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		resp = httptest.NewRecorder()
// todo(fs): 		obj, err = a.srv.HealthServiceChecks(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		assertIndex(t, resp)
// todo(fs): 		nodes = obj.(structs.HealthChecks)
// todo(fs): 		if len(nodes) != 2 {
// todo(fs): 			r.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 		if nodes[0].Node != "foo" {
// todo(fs): 			r.Fatalf("bad: %v", nodes)
// todo(fs): 		}
// todo(fs): 		if nodes[1].Node != "bar" {
// todo(fs): 			r.Fatalf("bad: %v", nodes)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceNodes(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/service/consul?dc=dc1", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be 1 health check for consul
// todo(fs): 	nodes := obj.(structs.CheckServiceNodes)
// todo(fs): 	if len(nodes) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/health/service/nope?dc=dc1", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.HealthServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be a non-nil empty list
// todo(fs): 	nodes = obj.(structs.CheckServiceNodes)
// todo(fs): 	if nodes == nil || len(nodes) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			ID:      "test",
// todo(fs): 			Service: "test",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/health/service/test?dc=dc1", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.HealthServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be a non-nil empty list for checks
// todo(fs): 	nodes = obj.(structs.CheckServiceNodes)
// todo(fs): 	if len(nodes) != 1 || nodes[0].Checks == nil || len(nodes[0].Checks) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceNodes_NodeMetaFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/service/consul?dc=dc1&node-meta=somekey:somevalue", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be a non-nil empty list
// todo(fs): 	nodes := obj.(structs.CheckServiceNodes)
// todo(fs): 	if nodes == nil || len(nodes) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		NodeMeta:   map[string]string{"somekey": "somevalue"},
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			ID:      "test",
// todo(fs): 			Service: "test",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/health/service/test?dc=dc1&node-meta=somekey:somevalue", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.HealthServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Should be a non-nil empty list for checks
// todo(fs): 	nodes = obj.(structs.CheckServiceNodes)
// todo(fs): 	if len(nodes) != 1 || nodes[0].Checks == nil || len(nodes[0].Checks) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceNodes_DistanceSort(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create a service check
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			ID:      "test",
// todo(fs): 			Service: "test",
// todo(fs): 		},
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			Node:      "bar",
// todo(fs): 			Name:      "test check",
// todo(fs): 			ServiceID: "test",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args.Node, args.Check.Node = "foo", "foo"
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/service/test?dc=dc1&near=foo", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.HealthServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs): 	nodes := obj.(structs.CheckServiceNodes)
// todo(fs): 	if len(nodes) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): 	if nodes[0].Node.Node != "bar" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[1].Node.Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Send an update for the node and wait for it to get applied.
// todo(fs): 	arg := structs.CoordinateUpdateRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Coord:      coordinate.NewCoordinate(coordinate.DefaultConfig()),
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Coordinate.Update", &arg, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	// Retry until foo has moved to the front of the line.
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		resp = httptest.NewRecorder()
// todo(fs): 		obj, err = a.srv.HealthServiceNodes(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		assertIndex(t, resp)
// todo(fs): 		nodes = obj.(structs.CheckServiceNodes)
// todo(fs): 		if len(nodes) != 2 {
// todo(fs): 			r.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 		if nodes[0].Node.Node != "foo" {
// todo(fs): 			r.Fatalf("bad: %v", nodes)
// todo(fs): 		}
// todo(fs): 		if nodes[1].Node.Node != "bar" {
// todo(fs): 			r.Fatalf("bad: %v", nodes)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceNodes_PassingFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create a failing service check
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			Name:      "consul check",
// todo(fs): 			ServiceID: "consul",
// todo(fs): 			Status:    api.HealthCritical,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("bc_no_query_value", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/health/service/consul?passing", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.HealthServiceNodes(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		// Should be 0 health check for consul
// todo(fs): 		nodes := obj.(structs.CheckServiceNodes)
// todo(fs): 		if len(nodes) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("passing_true", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/health/service/consul?passing=true", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.HealthServiceNodes(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		// Should be 0 health check for consul
// todo(fs): 		nodes := obj.(structs.CheckServiceNodes)
// todo(fs): 		if len(nodes) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("passing_false", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/health/service/consul?passing=false", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.HealthServiceNodes(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		// Should be 1 consul, it's unhealthy, but we specifically asked for
// todo(fs): 		// everything.
// todo(fs): 		nodes := obj.(structs.CheckServiceNodes)
// todo(fs): 		if len(nodes) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("passing_bad", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/health/service/consul?passing=nope-nope-nope", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		a.srv.HealthServiceNodes(resp, req)
// todo(fs):
// todo(fs): 		if code := resp.Code; code != 400 {
// todo(fs): 			t.Errorf("bad response code %d, expected %d", code, 400)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		body, err := ioutil.ReadAll(resp.Body)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatal(err)
// todo(fs): 		}
// todo(fs): 		if !bytes.Contains(body, []byte("Invalid value for ?passing")) {
// todo(fs): 			t.Errorf("bad %s", body)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestHealthServiceNodes_WanTranslation(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg1 := TestConfig()
// todo(fs): 	cfg1.Datacenter = "dc1"
// todo(fs): 	cfg1.TranslateWanAddrs = true
// todo(fs): 	cfg1.ACLDatacenter = ""
// todo(fs): 	a1 := NewTestAgent(t.Name(), cfg1)
// todo(fs): 	defer a1.Shutdown()
// todo(fs):
// todo(fs): 	cfg2 := TestConfig()
// todo(fs): 	cfg2.Datacenter = "dc2"
// todo(fs): 	cfg2.TranslateWanAddrs = true
// todo(fs): 	cfg2.ACLDatacenter = ""
// todo(fs): 	a2 := NewTestAgent(t.Name(), cfg2)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	// Wait for the WAN join.
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a1.Config.Ports.SerfWan)
// todo(fs): 	if _, err := a2.JoinWAN([]string{addr}); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if got, want := len(a1.WANMembers()), 2; got < want {
// todo(fs): 			r.Fatalf("got %d WAN members want at least %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Register a node with DC2.
// todo(fs): 	{
// todo(fs): 		args := &structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc2",
// todo(fs): 			Node:       "foo",
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 			TaggedAddresses: map[string]string{
// todo(fs): 				"wan": "127.0.0.2",
// todo(fs): 			},
// todo(fs): 			Service: &structs.NodeService{
// todo(fs): 				Service: "http_wan_translation_test",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs):
// todo(fs): 		var out struct{}
// todo(fs): 		if err := a2.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Query for a service in DC2 from DC1.
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/health/service/http_wan_translation_test?dc=dc2", nil)
// todo(fs): 	resp1 := httptest.NewRecorder()
// todo(fs): 	obj1, err1 := a1.srv.HealthServiceNodes(resp1, req)
// todo(fs): 	if err1 != nil {
// todo(fs): 		t.Fatalf("err: %v", err1)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp1)
// todo(fs):
// todo(fs): 	// Expect that DC1 gives us a WAN address (since the node is in DC2).
// todo(fs): 	nodes1 := obj1.(structs.CheckServiceNodes)
// todo(fs): 	if len(nodes1) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj1)
// todo(fs): 	}
// todo(fs): 	node1 := nodes1[0].Node
// todo(fs): 	if node1.Address != "127.0.0.2" {
// todo(fs): 		t.Fatalf("bad: %v", node1)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Query DC2 from DC2.
// todo(fs): 	resp2 := httptest.NewRecorder()
// todo(fs): 	obj2, err2 := a2.srv.HealthServiceNodes(resp2, req)
// todo(fs): 	if err2 != nil {
// todo(fs): 		t.Fatalf("err: %v", err2)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp2)
// todo(fs):
// todo(fs): 	// Expect that DC2 gives us a private address (since the node is in DC2).
// todo(fs): 	nodes2 := obj2.(structs.CheckServiceNodes)
// todo(fs): 	if len(nodes2) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj2)
// todo(fs): 	}
// todo(fs): 	node2 := nodes2[0].Node
// todo(fs): 	if node2.Address != "127.0.0.1" {
// todo(fs): 		t.Fatalf("bad: %v", node2)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestFilterNonPassing(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	nodes := structs.CheckServiceNodes{
// todo(fs): 		structs.CheckServiceNode{
// todo(fs): 			Checks: structs.HealthChecks{
// todo(fs): 				&structs.HealthCheck{
// todo(fs): 					Status: api.HealthCritical,
// todo(fs): 				},
// todo(fs): 				&structs.HealthCheck{
// todo(fs): 					Status: api.HealthCritical,
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 		structs.CheckServiceNode{
// todo(fs): 			Checks: structs.HealthChecks{
// todo(fs): 				&structs.HealthCheck{
// todo(fs): 					Status: api.HealthCritical,
// todo(fs): 				},
// todo(fs): 				&structs.HealthCheck{
// todo(fs): 					Status: api.HealthCritical,
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 		structs.CheckServiceNode{
// todo(fs): 			Checks: structs.HealthChecks{
// todo(fs): 				&structs.HealthCheck{
// todo(fs): 					Status: api.HealthPassing,
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	out := filterNonPassing(nodes)
// todo(fs): 	if len(out) != 1 && reflect.DeepEqual(out[0], nodes[2]) {
// todo(fs): 		t.Fatalf("bad: %v", out)
// todo(fs): 	}
// todo(fs): }
