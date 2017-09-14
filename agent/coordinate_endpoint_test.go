package agent

// todo(fs): func TestCoordinate_Datacenters(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/coordinate/datacenters", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CoordinateDatacenters(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	maps := obj.([]structs.DatacenterMap)
// todo(fs): 	if len(maps) != 1 ||
// todo(fs): 		maps[0].Datacenter != "dc1" ||
// todo(fs): 		len(maps[0].Coordinates) != 1 ||
// todo(fs): 		maps[0].Coordinates[0].Node != a.Config.NodeName {
// todo(fs): 		t.Fatalf("bad: %v", maps)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestCoordinate_Nodes(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Make sure an empty list is non-nil.
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/coordinate/nodes?dc=dc1", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.CoordinateNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	coordinates := obj.(structs.Coordinates)
// todo(fs): 	if coordinates == nil || len(coordinates) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", coordinates)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register the nodes.
// todo(fs): 	nodes := []string{"foo", "bar"}
// todo(fs): 	for _, node := range nodes {
// todo(fs): 		req := structs.RegisterRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       node,
// todo(fs): 			Address:    "127.0.0.1",
// todo(fs): 		}
// todo(fs): 		var reply struct{}
// todo(fs): 		if err := a.RPC("Catalog.Register", &req, &reply); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Send some coordinates for a few nodes, waiting a little while for the
// todo(fs): 	// batch update to run.
// todo(fs): 	arg1 := structs.CoordinateUpdateRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "foo",
// todo(fs): 		Segment:    "alpha",
// todo(fs): 		Coord:      coordinate.NewCoordinate(coordinate.DefaultConfig()),
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Coordinate.Update", &arg1, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	arg2 := structs.CoordinateUpdateRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       "bar",
// todo(fs): 		Coord:      coordinate.NewCoordinate(coordinate.DefaultConfig()),
// todo(fs): 	}
// todo(fs): 	if err := a.RPC("Coordinate.Update", &arg2, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	time.Sleep(300 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Query back and check the nodes are present and sorted correctly.
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/coordinate/nodes?dc=dc1", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.CoordinateNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	coordinates = obj.(structs.Coordinates)
// todo(fs): 	if len(coordinates) != 2 ||
// todo(fs): 		coordinates[0].Node != "bar" ||
// todo(fs): 		coordinates[1].Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", coordinates)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Filter on a nonexistant node segment
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/coordinate/nodes?segment=nope", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.CoordinateNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	coordinates = obj.(structs.Coordinates)
// todo(fs): 	if len(coordinates) != 0 {
// todo(fs): 		t.Fatalf("bad: %v", coordinates)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Filter on a real node segment
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/coordinate/nodes?segment=alpha", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.CoordinateNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	coordinates = obj.(structs.Coordinates)
// todo(fs): 	if len(coordinates) != 1 || coordinates[0].Node != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", coordinates)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure the empty filter works
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/coordinate/nodes?segment=", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.CoordinateNodes(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	coordinates = obj.(structs.Coordinates)
// todo(fs): 	if len(coordinates) != 1 || coordinates[0].Node != "bar" {
// todo(fs): 		t.Fatalf("bad: %v", coordinates)
// todo(fs): 	}
// todo(fs): }
