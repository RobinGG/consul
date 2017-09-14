package agent

// todo(fs): func TestSnapshot(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	var snap io.Reader
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/snapshot?token=root", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.Snapshot(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		snap = resp.Body
// todo(fs):
// todo(fs): 		header := resp.Header().Get("X-Consul-Index")
// todo(fs): 		if header == "" {
// todo(fs): 			t.Fatalf("bad: %v", header)
// todo(fs): 		}
// todo(fs): 		header = resp.Header().Get("X-Consul-KnownLeader")
// todo(fs): 		if header != "true" {
// todo(fs): 			t.Fatalf("bad: %v", header)
// todo(fs): 		}
// todo(fs): 		header = resp.Header().Get("X-Consul-LastContact")
// todo(fs): 		if header != "0" {
// todo(fs): 			t.Fatalf("bad: %v", header)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/snapshot?token=root", snap)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.Snapshot(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestSnapshot_Options(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	for _, method := range []string{"GET", "PUT"} {
// todo(fs): 		t.Run(method, func(t *testing.T) {
// todo(fs): 			a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 			defer a.Shutdown()
// todo(fs):
// todo(fs): 			body := bytes.NewBuffer(nil)
// todo(fs): 			req, _ := http.NewRequest(method, "/v1/snapshot?token=anonymous", body)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			_, err := a.srv.Snapshot(resp, req)
// todo(fs): 			if !acl.IsErrPermissionDenied(err) {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs):
// todo(fs): 		t.Run(method, func(t *testing.T) {
// todo(fs): 			a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 			defer a.Shutdown()
// todo(fs):
// todo(fs): 			body := bytes.NewBuffer(nil)
// todo(fs): 			req, _ := http.NewRequest(method, "/v1/snapshot?dc=nope", body)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			_, err := a.srv.Snapshot(resp, req)
// todo(fs): 			if err == nil || !strings.Contains(err.Error(), "No path to datacenter") {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs):
// todo(fs): 		t.Run(method, func(t *testing.T) {
// todo(fs): 			a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 			defer a.Shutdown()
// todo(fs):
// todo(fs): 			body := bytes.NewBuffer(nil)
// todo(fs): 			req, _ := http.NewRequest(method, "/v1/snapshot?token=root&stale", body)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			_, err := a.srv.Snapshot(resp, req)
// todo(fs): 			if method == "GET" {
// todo(fs): 				if err != nil {
// todo(fs): 					t.Fatalf("err: %v", err)
// todo(fs): 				}
// todo(fs): 			} else {
// todo(fs): 				if err == nil || !strings.Contains(err.Error(), "stale not allowed") {
// todo(fs): 					t.Fatalf("err: %v", err)
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestSnapshot_BadMethods(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("POST", "/v1/snapshot", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		_, err := a.srv.Snapshot(resp, req)
// todo(fs): 		if err != nil {
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
// todo(fs): 		req, _ := http.NewRequest("DELETE", "/v1/snapshot", body)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		_, err := a.srv.Snapshot(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 405 {
// todo(fs): 			t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
