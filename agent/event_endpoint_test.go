package agent

// todo(fs): func TestEventFire(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer([]byte("test"))
// todo(fs): 	url := "/v1/event/fire/test?node=Node&service=foo&tag=bar"
// todo(fs): 	req, _ := http.NewRequest("PUT", url, body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.EventFire(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	event, ok := obj.(*UserEvent)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("bad: %#v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if event.ID == "" {
// todo(fs): 		t.Fatalf("bad: %#v", event)
// todo(fs): 	}
// todo(fs): 	if event.Name != "test" {
// todo(fs): 		t.Fatalf("bad: %#v", event)
// todo(fs): 	}
// todo(fs): 	if string(event.Payload) != "test" {
// todo(fs): 		t.Fatalf("bad: %#v", event)
// todo(fs): 	}
// todo(fs): 	if event.NodeFilter != "Node" {
// todo(fs): 		t.Fatalf("bad: %#v", event)
// todo(fs): 	}
// todo(fs): 	if event.ServiceFilter != "foo" {
// todo(fs): 		t.Fatalf("bad: %#v", event)
// todo(fs): 	}
// todo(fs): 	if event.TagFilter != "bar" {
// todo(fs): 		t.Fatalf("bad: %#v", event)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestEventFire_token(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
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
// todo(fs): 		event   string
// todo(fs): 		allowed bool
// todo(fs): 	}
// todo(fs): 	tcases := []tcase{
// todo(fs): 		{"foo", false},
// todo(fs): 		{"bar", false},
// todo(fs): 		{"baz", true},
// todo(fs): 	}
// todo(fs): 	for _, c := range tcases {
// todo(fs): 		// Try to fire the event over the HTTP interface
// todo(fs): 		url := fmt.Sprintf("/v1/event/fire/%s?token=%s", c.event, token)
// todo(fs): 		req, _ := http.NewRequest("PUT", url, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.EventFire(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Check the result
// todo(fs): 		body := resp.Body.String()
// todo(fs): 		if c.allowed {
// todo(fs): 			if acl.IsErrPermissionDenied(errors.New(body)) {
// todo(fs): 				t.Fatalf("bad: %s", body)
// todo(fs): 			}
// todo(fs): 			if resp.Code != 200 {
// todo(fs): 				t.Fatalf("bad: %d", resp.Code)
// todo(fs): 			}
// todo(fs): 		} else {
// todo(fs): 			if !acl.IsErrPermissionDenied(errors.New(body)) {
// todo(fs): 				t.Fatalf("bad: %s", body)
// todo(fs): 			}
// todo(fs): 			if resp.Code != 403 {
// todo(fs): 				t.Fatalf("bad: %d", resp.Code)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestEventList(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	p := &UserEvent{Name: "test"}
// todo(fs): 	if err := a.UserEvent("dc1", "root", p); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/event/list", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.EventList(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		list, ok := obj.([]*UserEvent)
// todo(fs): 		if !ok {
// todo(fs): 			r.Fatalf("bad: %#v", obj)
// todo(fs): 		}
// todo(fs): 		if len(list) != 1 || list[0].Name != "test" {
// todo(fs): 			r.Fatalf("bad: %#v", list)
// todo(fs): 		}
// todo(fs): 		header := resp.Header().Get("X-Consul-Index")
// todo(fs): 		if header == "" || header == "0" {
// todo(fs): 			r.Fatalf("bad: %#v", header)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestEventList_Filter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	p := &UserEvent{Name: "test"}
// todo(fs): 	if err := a.UserEvent("dc1", "root", p); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	p = &UserEvent{Name: "foo"}
// todo(fs): 	if err := a.UserEvent("dc1", "root", p); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/event/list?name=foo", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.EventList(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		list, ok := obj.([]*UserEvent)
// todo(fs): 		if !ok {
// todo(fs): 			r.Fatalf("bad: %#v", obj)
// todo(fs): 		}
// todo(fs): 		if len(list) != 1 || list[0].Name != "foo" {
// todo(fs): 			r.Fatalf("bad: %#v", list)
// todo(fs): 		}
// todo(fs): 		header := resp.Header().Get("X-Consul-Index")
// todo(fs): 		if header == "" || header == "0" {
// todo(fs): 			r.Fatalf("bad: %#v", header)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestEventList_ACLFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Fire an event.
// todo(fs): 	p := &UserEvent{Name: "foo"}
// todo(fs): 	if err := a.UserEvent("dc1", "root", p); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		retry.Run(t, func(r *retry.R) {
// todo(fs): 			req, _ := http.NewRequest("GET", "/v1/event/list", nil)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.EventList(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				r.Fatal(err)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			list, ok := obj.([]*UserEvent)
// todo(fs): 			if !ok {
// todo(fs): 				r.Fatalf("bad: %#v", obj)
// todo(fs): 			}
// todo(fs): 			if len(list) != 0 {
// todo(fs): 				r.Fatalf("bad: %#v", list)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		retry.Run(t, func(r *retry.R) {
// todo(fs): 			req, _ := http.NewRequest("GET", "/v1/event/list?token=root", nil)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.EventList(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				r.Fatal(err)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			list, ok := obj.([]*UserEvent)
// todo(fs): 			if !ok {
// todo(fs): 				r.Fatalf("bad: %#v", obj)
// todo(fs): 			}
// todo(fs): 			if len(list) != 1 || list[0].Name != "foo" {
// todo(fs): 				r.Fatalf("bad: %#v", list)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestEventList_Blocking(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	p := &UserEvent{Name: "test"}
// todo(fs): 	if err := a.UserEvent("dc1", "root", p); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var index string
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/event/list", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.EventList(resp, req); err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs): 		header := resp.Header().Get("X-Consul-Index")
// todo(fs): 		if header == "" || header == "0" {
// todo(fs): 			r.Fatalf("bad: %#v", header)
// todo(fs): 		}
// todo(fs): 		index = header
// todo(fs): 	})
// todo(fs):
// todo(fs): 	go func() {
// todo(fs): 		time.Sleep(50 * time.Millisecond)
// todo(fs): 		p := &UserEvent{Name: "second"}
// todo(fs): 		if err := a.UserEvent("dc1", "root", p); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		url := "/v1/event/list?index=" + index
// todo(fs): 		req, _ := http.NewRequest("GET", url, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.EventList(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		list, ok := obj.([]*UserEvent)
// todo(fs): 		if !ok {
// todo(fs): 			r.Fatalf("bad: %#v", obj)
// todo(fs): 		}
// todo(fs): 		if len(list) != 2 || list[1].Name != "second" {
// todo(fs): 			r.Fatalf("bad: %#v", list)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestEventList_EventBufOrder(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Fire some events in a non-sequential order
// todo(fs): 	expected := &UserEvent{Name: "foo"}
// todo(fs):
// todo(fs): 	for _, e := range []*UserEvent{
// todo(fs): 		&UserEvent{Name: "foo"},
// todo(fs): 		&UserEvent{Name: "bar"},
// todo(fs): 		&UserEvent{Name: "foo"},
// todo(fs): 		expected,
// todo(fs): 		&UserEvent{Name: "bar"},
// todo(fs): 	} {
// todo(fs): 		if err := a.UserEvent("dc1", "root", e); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): 	// Test that the event order is preserved when name
// todo(fs): 	// filtering on a list of > 1 matching event.
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		url := "/v1/event/list?name=foo"
// todo(fs): 		req, _ := http.NewRequest("GET", url, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.EventList(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs): 		list, ok := obj.([]*UserEvent)
// todo(fs): 		if !ok {
// todo(fs): 			r.Fatalf("bad: %#v", obj)
// todo(fs): 		}
// todo(fs): 		if len(list) != 3 || list[2].ID != expected.ID {
// todo(fs): 			r.Fatalf("bad: %#v", list)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestUUIDToUint64(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	inp := "cb9a81ad-fff6-52ac-92a7-5f70687805ec"
// todo(fs):
// todo(fs): 	// Output value was computed using python
// todo(fs): 	if uuidToUint64(inp) != 6430540886266763072 {
// todo(fs): 		t.Fatalf("bad")
// todo(fs): 	}
// todo(fs): }
