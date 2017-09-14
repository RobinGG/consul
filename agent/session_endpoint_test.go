package agent

// todo(fs): func TestSessionCreate(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create a health check
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			CheckID:   "consul",
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			Name:      "consul",
// todo(fs): 			ServiceID: "consul",
// todo(fs): 			Status:    api.HealthPassing,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Associate session with node and 2 health checks
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"Name":      "my-cool-session",
// todo(fs): 		"Node":      a.Config.NodeName,
// todo(fs): 		"Checks":    []types.CheckID{structs.SerfCheckID, "consul"},
// todo(fs): 		"LockDelay": "20s",
// todo(fs): 	}
// todo(fs): 	enc.Encode(raw)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/session/create", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.SessionCreate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if _, ok := obj.(sessionCreateResponse); !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionCreateDelete(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create a health check
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 		Check: &structs.HealthCheck{
// todo(fs): 			CheckID:   "consul",
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			Name:      "consul",
// todo(fs): 			ServiceID: "consul",
// todo(fs): 			Status:    api.HealthPassing,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Associate session with node and 2 health checks, and make it delete on session destroy
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"Name":      "my-cool-session",
// todo(fs): 		"Node":      a.Config.NodeName,
// todo(fs): 		"Checks":    []types.CheckID{structs.SerfCheckID, "consul"},
// todo(fs): 		"LockDelay": "20s",
// todo(fs): 		"Behavior":  structs.SessionKeysDelete,
// todo(fs): 	}
// todo(fs): 	enc.Encode(raw)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/session/create", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.SessionCreate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if _, ok := obj.(sessionCreateResponse); !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestFixupLockDelay(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	inp := map[string]interface{}{
// todo(fs): 		"lockdelay": float64(15),
// todo(fs): 	}
// todo(fs): 	if err := FixupLockDelay(inp); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if inp["lockdelay"] != 15*time.Second {
// todo(fs): 		t.Fatalf("bad: %v", inp)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	inp = map[string]interface{}{
// todo(fs): 		"lockDelay": float64(15 * time.Second),
// todo(fs): 	}
// todo(fs): 	if err := FixupLockDelay(inp); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if inp["lockDelay"] != 15*time.Second {
// todo(fs): 		t.Fatalf("bad: %v", inp)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	inp = map[string]interface{}{
// todo(fs): 		"LockDelay": "15s",
// todo(fs): 	}
// todo(fs): 	if err := FixupLockDelay(inp); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if inp["LockDelay"] != 15*time.Second {
// todo(fs): 		t.Fatalf("bad: %v", inp)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func makeTestSession(t *testing.T, srv *HTTPServer) string {
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/session/create", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := srv.SessionCreate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	sessResp := obj.(sessionCreateResponse)
// todo(fs): 	return sessResp.ID
// todo(fs): }
// todo(fs):
// todo(fs): func makeTestSessionDelete(t *testing.T, srv *HTTPServer) string {
// todo(fs): 	// Create Session with delete behavior
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"Behavior": "delete",
// todo(fs): 	}
// todo(fs): 	enc.Encode(raw)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/session/create", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := srv.SessionCreate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	sessResp := obj.(sessionCreateResponse)
// todo(fs): 	return sessResp.ID
// todo(fs): }
// todo(fs):
// todo(fs): func makeTestSessionTTL(t *testing.T, srv *HTTPServer, ttl string) string {
// todo(fs): 	// Create Session with TTL
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	enc := json.NewEncoder(body)
// todo(fs): 	raw := map[string]interface{}{
// todo(fs): 		"TTL": ttl,
// todo(fs): 	}
// todo(fs): 	enc.Encode(raw)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/session/create", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := srv.SessionCreate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	sessResp := obj.(sessionCreateResponse)
// todo(fs): 	return sessResp.ID
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionDestroy(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	id := makeTestSession(t, a.srv)
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/session/destroy/"+id, nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.SessionDestroy(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp := obj.(bool); !resp {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionCustomTTL(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	ttl := 250 * time.Millisecond
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.SessionTTLMin = ttl
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	id := makeTestSessionTTL(t, a.srv, ttl.String())
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/session/info/"+id, nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.SessionGet(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok := obj.(structs.Sessions)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): 	if len(respObj) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", respObj)
// todo(fs): 	}
// todo(fs): 	if respObj[0].TTL != ttl.String() {
// todo(fs): 		t.Fatalf("Incorrect TTL: %s", respObj[0].TTL)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	time.Sleep(ttl*structs.SessionTTLMultiplier + ttl)
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/session/info/"+id, nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.SessionGet(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok = obj.(structs.Sessions)
// todo(fs): 	if len(respObj) != 0 {
// todo(fs): 		t.Fatalf("session '%s' should have been destroyed", id)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionTTLRenew(t *testing.T) {
// todo(fs): 	// t.Parallel() // timing test. no parallel
// todo(fs): 	ttl := 250 * time.Millisecond
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.SessionTTLMin = ttl
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	id := makeTestSessionTTL(t, a.srv, ttl.String())
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/session/info/"+id, nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.SessionGet(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok := obj.(structs.Sessions)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): 	if len(respObj) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", respObj)
// todo(fs): 	}
// todo(fs): 	if respObj[0].TTL != ttl.String() {
// todo(fs): 		t.Fatalf("Incorrect TTL: %s", respObj[0].TTL)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Sleep to consume some time before renew
// todo(fs): 	time.Sleep(ttl * (structs.SessionTTLMultiplier / 2))
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("PUT", "/v1/session/renew/"+id, nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.SessionRenew(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok = obj.(structs.Sessions)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs): 	if len(respObj) != 1 {
// todo(fs): 		t.Fatalf("bad: %v", respObj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Sleep for ttl * TTL Multiplier
// todo(fs): 	time.Sleep(ttl * structs.SessionTTLMultiplier)
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/session/info/"+id, nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.SessionGet(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok = obj.(structs.Sessions)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("session '%s' should have renewed", id)
// todo(fs): 	}
// todo(fs): 	if len(respObj) != 1 {
// todo(fs): 		t.Fatalf("session '%s' should have renewed", id)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// now wait for timeout and expect session to get destroyed
// todo(fs): 	time.Sleep(ttl * structs.SessionTTLMultiplier)
// todo(fs):
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/session/info/"+id, nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.SessionGet(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	respObj, ok = obj.(structs.Sessions)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("session '%s' should have destroyed", id)
// todo(fs): 	}
// todo(fs): 	if len(respObj) != 0 {
// todo(fs): 		t.Fatalf("session '%s' should have destroyed", id)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionGet(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/session/info/adf4238a-882b-9ddc-4a9d-5b6758e4159e", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.SessionGet(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.Sessions)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if respObj == nil || len(respObj) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		id := makeTestSession(t, a.srv)
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/session/info/"+id, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.SessionGet(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.Sessions)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if len(respObj) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionList(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/session/list", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.SessionList(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.Sessions)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if respObj == nil || len(respObj) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		var ids []string
// todo(fs): 		for i := 0; i < 10; i++ {
// todo(fs): 			ids = append(ids, makeTestSession(t, a.srv))
// todo(fs): 		}
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/session/list", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.SessionList(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.Sessions)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if len(respObj) != 10 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionsForNode(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/session/node/"+a.Config.NodeName, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.SessionsForNode(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.Sessions)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if respObj == nil || len(respObj) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		var ids []string
// todo(fs): 		for i := 0; i < 10; i++ {
// todo(fs): 			ids = append(ids, makeTestSession(t, a.srv))
// todo(fs): 		}
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/session/node/"+a.Config.NodeName, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.SessionsForNode(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		respObj, ok := obj.(structs.Sessions)
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 		if len(respObj) != 10 {
// todo(fs): 			t.Fatalf("bad: %v", respObj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestSessionDeleteDestroy(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	id := makeTestSessionDelete(t, a.srv)
// todo(fs):
// todo(fs): 	// now create a new key for the session and acquire it
// todo(fs): 	buf := bytes.NewBuffer([]byte("test"))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/kv/ephemeral?acquire="+id, buf)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.KVSEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if res := obj.(bool); !res {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// now destroy the session, this should delete the key created above
// todo(fs): 	req, _ = http.NewRequest("PUT", "/v1/session/destroy/"+id, nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, err = a.srv.SessionDestroy(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp := obj.(bool); !resp {
// todo(fs): 		t.Fatalf("should work")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify that the key is gone
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/kv/ephemeral", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	obj, _ = a.srv.KVSEndpoint(resp, req)
// todo(fs): 	res, found := obj.(structs.DirEntries)
// todo(fs): 	if found || len(res) != 0 {
// todo(fs): 		t.Fatalf("bad: %v found, should be nothing", res)
// todo(fs): 	}
// todo(fs): }
