package agent

// todo(fs): func TestCatalogRegister(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Node:    "foo",
// todo(fs): 		Address: "127.0.0.1",
// todo(fs): 	}
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/register", jsonReader(args))
// todo(fs): 	obj, err := a.srv.CatalogRegister(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	res := obj.(bool)
// todo(fs): 	if res != true {
// todo(fs): 		t.Fatalf("bad: %v", res)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	func() {
// todo(fs): 		a.state.Lock()
// todo(fs): 		defer a.state.Unlock()
// todo(fs):
// todo(fs): 		// Service should be in sync
// todo(fs): 		if err := a.state.syncService("foo"); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 		if _, ok := a.state.serviceStatus["foo"]; !ok {
// todo(fs): 			t.Fatalf("bad: %#v", a.state.serviceStatus)
// todo(fs): 		}
// todo(fs): 		if !a.state.serviceStatus["foo"].inSync {
// todo(fs): 			t.Fatalf("should be in sync")
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogRegister_Service_InvalidAddress(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	for _, addr := range []string{"0.0.0.0", "::", "[::]"} {
// todo(fs): 		t.Run("addr "+addr, func(t *testing.T) {
// todo(fs): 			args := &structs.RegisterRequest{
// todo(fs): 				Node:    "foo",
// todo(fs): 				Address: "127.0.0.1",
// todo(fs): 				Service: &structs.NodeService{
// todo(fs): 					Service: "test",
// todo(fs): 					Address: addr,
// todo(fs): 					Port:    8080,
// todo(fs): 				},
// todo(fs): 			}
// todo(fs): 			req, _ := http.NewRequest("GET", "/v1/catalog/register", jsonReader(args))
// todo(fs): 			_, err := a.srv.CatalogRegister(nil, req)
// todo(fs): 			if err == nil || err.Error() != "Invalid service address" {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogDeregister(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.DeregisterRequest{Node: "foo"}
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/deregister", jsonReader(args))
// todo(fs): 	obj, err := a.srv.CatalogDeregister(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	res := obj.(bool)
// todo(fs): 	if res != true {
// todo(fs): 		t.Fatalf("bad: %v", res)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogDatacenters(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		obj, err := a.srv.CatalogDatacenters(nil, nil)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		dcs := obj.([]string)
// todo(fs): 		if got, want := len(dcs), 1; got != want {
// todo(fs): 			r.Fatalf("got %d data centers want %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogNodes(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
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
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?dc=dc1", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify an index is set
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	nodes := obj.(structs.Nodes)
// todo(fs): 	if len(nodes) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogNodes_MetaFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register a node with a meta field
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		NodeMeta: map[string]string{
// todo(fs): 			"somekey": "somevalue",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?node-meta=somekey:somevalue", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify an index is set
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	// Verify we only get the node with the correct meta field back
// todo(fs): 	nodes := obj.(structs.Nodes)
// todo(fs): 	if len(nodes) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): 	if v, ok := nodes[0].Meta["somekey"]; !ok || v != "somevalue" {
// todo(fs): 		t.Fatalf("bad: %v", nodes[0].Meta)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogNodes_WanTranslation(t *testing.T) {
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
// todo(fs): 	// Wait for the WAN join.
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a1.Config.SerfPortWAN)
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
// todo(fs): 			Node:       "wan_translation_test",
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
// todo(fs): 	// Query nodes in DC2 from DC1.
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?dc=dc2", nil)
// todo(fs): 	resp1 := httptest.NewRecorder()
// todo(fs): 	obj1, err1 := a1.srv.CatalogNodes(resp1, req)
// todo(fs): 	if err1 != nil {
// todo(fs): 		t.Fatalf("err: %v", err1)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp1)
// todo(fs):
// todo(fs): 	// Expect that DC1 gives us a WAN address (since the node is in DC2).
// todo(fs): 	nodes1 := obj1.(structs.Nodes)
// todo(fs): 	if len(nodes1) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", obj1)
// todo(fs): 	}
// todo(fs): 	var address string
// todo(fs): 	for _, node := range nodes1 {
// todo(fs): 		if node.Node == "wan_translation_test" {
// todo(fs): 			address = node.Address
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): 	if address != "127.0.0.2" {
// todo(fs): 		t.Fatalf("bad: %s", address)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Query DC2 from DC2.
// todo(fs): 	resp2 := httptest.NewRecorder()
// todo(fs): 	obj2, err2 := a2.srv.CatalogNodes(resp2, req)
// todo(fs): 	if err2 != nil {
// todo(fs): 		t.Fatalf("err: %v", err2)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp2)
// todo(fs):
// todo(fs): 	// Expect that DC2 gives us a private address (since the node is in DC2).
// todo(fs): 	nodes2 := obj2.(structs.Nodes)
// todo(fs): 	if len(nodes2) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", obj2)
// todo(fs): 	}
// todo(fs): 	for _, node := range nodes2 {
// todo(fs): 		if node.Node == "wan_translation_test" {
// todo(fs): 			address = node.Address
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): 	if address != "127.0.0.1" {
// todo(fs): 		t.Fatalf("bad: %s", address)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogNodes_Blocking(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.DCSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out structs.IndexedNodes
// todo(fs): 	if err := a.RPC("Catalog.ListNodes", *args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// t.Fatal must be called from the main go routine
// todo(fs): 	// of the test. Because of this we cannot call
// todo(fs): 	// t.Fatal from within the go routines and use
// todo(fs): 	// an error channel instead.
// todo(fs): 	errch := make(chan error, 2)
// todo(fs): 	go func() {
// todo(fs): 		start := time.Now()
// todo(fs):
// todo(fs): 		// register a service after the blocking call
// todo(fs): 		// in order to unblock it.
// todo(fs): 		time.AfterFunc(100*time.Millisecond, func() {
// todo(fs): 			args := &structs.RegisterRequest{
// todo(fs): 				Datacenter: "dc1",
// todo(fs): 				Node:       "foo",
// todo(fs): 				Address:    "127.0.0.1",
// todo(fs): 			}
// todo(fs): 			var out struct{}
// todo(fs): 			errch <- a.RPC("Catalog.Register", args, &out)
// todo(fs): 		})
// todo(fs):
// todo(fs): 		// now block
// todo(fs): 		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/catalog/nodes?wait=3s&index=%d", out.Index+1), nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.CatalogNodes(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			errch <- err
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Should block for a while
// todo(fs): 		if d := time.Now().Sub(start); d < 50*time.Millisecond {
// todo(fs): 			errch <- fmt.Errorf("too fast: %v", d)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if idx := getIndex(t, resp); idx <= out.Index {
// todo(fs): 			errch <- fmt.Errorf("bad: %v", idx)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		nodes := obj.(structs.Nodes)
// todo(fs): 		if len(nodes) != 2 {
// todo(fs): 			errch <- fmt.Errorf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 		errch <- nil
// todo(fs): 	}()
// todo(fs):
// todo(fs): 	// wait for both go routines to return
// todo(fs): 	if err := <-errch; err != nil {
// todo(fs): 		t.Fatal(err)
// todo(fs): 	}
// todo(fs): 	if err := <-errch; err != nil {
// todo(fs): 		t.Fatal(err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogNodes_DistanceSort(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register nodes.
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args = &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.2",
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Nobody has coordinates set so this will still return them in the
// todo(fs): 	// order they are indexed.
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/nodes?dc=dc1&near=foo", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs): 	nodes := obj.(structs.Nodes)
// todo(fs): 	if len(nodes) != 3 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): 	if nodes[0].Node != "bar" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[1].Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[2].Node != a.Config.NodeName {
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
// todo(fs): 	time.Sleep(300 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Query again and now foo should have moved to the front of the line.
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/catalog/nodes?dc=dc1&near=foo", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.CatalogNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs): 	nodes = obj.(structs.Nodes)
// todo(fs): 	if len(nodes) != 3 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): 	if nodes[0].Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[1].Node != "bar" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[2].Node != a.Config.NodeName {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogServices(t *testing.T) {
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
// todo(fs): 			Service: "api",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/services?dc=dc1", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogServices(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	services := obj.(structs.Services)
// todo(fs): 	if len(services) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogServices_NodeMetaFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		NodeMeta: map[string]string{
// todo(fs): 			"somekey": "somevalue",
// todo(fs): 		},
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "api",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/services?node-meta=somekey:somevalue", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogServices(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	services := obj.(structs.Services)
// todo(fs): 	if len(services) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): 	if _, ok := services[args.Service.Service]; !ok {
// todo(fs): 		t.Fatalf("bad: %v", services)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogServiceNodes(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Make sure an empty list is returned, not a nil
// todo(fs): 	{
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/catalog/service/api?tag=a", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.CatalogServiceNodes(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		nodes := obj.(structs.ServiceNodes)
// todo(fs): 		if nodes == nil || len(nodes) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "api",
// todo(fs): 			Tags:    []string{"a"},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/service/api?tag=a", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	nodes := obj.(structs.ServiceNodes)
// todo(fs): 	if len(nodes) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogServiceNodes_NodeMetaFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Make sure an empty list is returned, not a nil
// todo(fs): 	{
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/catalog/service/api?node-meta=somekey:somevalue", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.CatalogServiceNodes(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		assertIndex(t, resp)
// todo(fs):
// todo(fs): 		nodes := obj.(structs.ServiceNodes)
// todo(fs): 		if nodes == nil || len(nodes) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		NodeMeta: map[string]string{
// todo(fs): 			"somekey": "somevalue",
// todo(fs): 		},
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "api",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/service/api?node-meta=somekey:somevalue", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	nodes := obj.(structs.ServiceNodes)
// todo(fs): 	if len(nodes) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogServiceNodes_WanTranslation(t *testing.T) {
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
// todo(fs): 	// Wait for the WAN join.
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a1.Config.SerfPortWAN)
// todo(fs): 	if _, err := a2.srv.agent.JoinWAN([]string{addr}); err != nil {
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
// todo(fs): 	// Query for the node in DC2 from DC1.
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/service/http_wan_translation_test?dc=dc2", nil)
// todo(fs): 	resp1 := httptest.NewRecorder()
// todo(fs): 	obj1, err1 := a1.srv.CatalogServiceNodes(resp1, req)
// todo(fs): 	if err1 != nil {
// todo(fs): 		t.Fatalf("err: %v", err1)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp1)
// todo(fs):
// todo(fs): 	// Expect that DC1 gives us a WAN address (since the node is in DC2).
// todo(fs): 	nodes1 := obj1.(structs.ServiceNodes)
// todo(fs): 	if len(nodes1) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj1)
// todo(fs): 	}
// todo(fs): 	node1 := nodes1[0]
// todo(fs): 	if node1.Address != "127.0.0.2" {
// todo(fs): 		t.Fatalf("bad: %v", node1)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Query DC2 from DC2.
// todo(fs): 	resp2 := httptest.NewRecorder()
// todo(fs): 	obj2, err2 := a2.srv.CatalogServiceNodes(resp2, req)
// todo(fs): 	if err2 != nil {
// todo(fs): 		t.Fatalf("err: %v", err2)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp2)
// todo(fs):
// todo(fs): 	// Expect that DC2 gives us a local address (since the node is in DC2).
// todo(fs): 	nodes2 := obj2.(structs.ServiceNodes)
// todo(fs): 	if len(nodes2) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj2)
// todo(fs): 	}
// todo(fs): 	node2 := nodes2[0]
// todo(fs): 	if node2.Address != "127.0.0.1" {
// todo(fs): 		t.Fatalf("bad: %v", node2)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogServiceNodes_DistanceSort(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register nodes.
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "api",
// todo(fs): 			Tags:    []string{"a"},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/service/api?tag=a", nil)
// todo(fs): 	args = &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Address:    "127.0.0.2",
// todo(fs): 		Service: &structs.NodeService{
// todo(fs): 			Service: "api",
// todo(fs): 			Tags:    []string{"a"},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Nobody has coordinates set so this will still return them in the
// todo(fs): 	// order they are indexed.
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/catalog/service/api?tag=a&near=foo", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs): 	nodes := obj.(structs.ServiceNodes)
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
// todo(fs): 	time.Sleep(300 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Query again and now foo should have moved to the front of the line.
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/catalog/service/api?tag=a&near=foo", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.CatalogServiceNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	assertIndex(t, resp)
// todo(fs): 	nodes = obj.(structs.ServiceNodes)
// todo(fs): 	if len(nodes) != 2 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): 	if nodes[0].Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): 	if nodes[1].Node != "bar" {
// todo(fs): 		t.Fatalf("bad: %v", nodes)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogNodeServices(t *testing.T) {
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
// todo(fs): 			Service: "api",
// todo(fs): 			Tags:    []string{"a"},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/node/foo?dc=dc1", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CatalogNodeServices(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp)
// todo(fs):
// todo(fs): 	services := obj.(*structs.NodeServices)
// todo(fs): 	if len(services.Services) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCatalogNodeServices_WanTranslation(t *testing.T) {
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
// todo(fs): 	// Wait for the WAN join.
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a1.Config.SerfPortWAN)
// todo(fs): 	if _, err := a2.srv.agent.JoinWAN([]string{addr}); err != nil {
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
// todo(fs): 	// Query for the node in DC2 from DC1.
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/catalog/node/foo?dc=dc2", nil)
// todo(fs): 	resp1 := httptest.NewRecorder()
// todo(fs): 	obj1, err1 := a1.srv.CatalogNodeServices(resp1, req)
// todo(fs): 	if err1 != nil {
// todo(fs): 		t.Fatalf("err: %v", err1)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp1)
// todo(fs):
// todo(fs): 	// Expect that DC1 gives us a WAN address (since the node is in DC2).
// todo(fs): 	services1 := obj1.(*structs.NodeServices)
// todo(fs): 	if len(services1.Services) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj1)
// todo(fs): 	}
// todo(fs): 	service1 := services1.Node
// todo(fs): 	if service1.Address != "127.0.0.2" {
// todo(fs): 		t.Fatalf("bad: %v", service1)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Query DC2 from DC2.
// todo(fs): 	resp2 := httptest.NewRecorder()
// todo(fs): 	obj2, err2 := a2.srv.CatalogNodeServices(resp2, req)
// todo(fs): 	if err2 != nil {
// todo(fs): 		t.Fatalf("err: %v", err2)
// todo(fs): 	}
// todo(fs): 	assertIndex(t, resp2)
// todo(fs):
// todo(fs): 	// Expect that DC2 gives us a private address (since the node is in DC2).
// todo(fs): 	services2 := obj2.(*structs.NodeServices)
// todo(fs): 	if len(services2.Services) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", obj2)
// todo(fs): 	}
// todo(fs): 	service2 := services2.Node
// todo(fs): 	if service2.Address != "127.0.0.1" {
// todo(fs): 		t.Fatalf("bad: %v", service2)
// todo(fs): 	}
// todo(fs): }
