package agent

// todo(fs): func makeReadOnlyAgentACL(t *testing.T, srv *HTTPServer) string {
// todo(fs): 	args := map[string]interface{}{
// todo(fs): 		"Name":  "User Token",
// todo(fs): 		"Type":  "client",
// todo(fs): 		"Rules": `agent "" { policy = "read" }`,
// todo(fs): 	}
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/acl/create?token=root", jsonReader(args))
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := srv.ACLCreate(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	aclResp := obj.(aclCreateResponse)
// todo(fs): 	return aclResp.ID
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Services(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:      "mysql",
// todo(fs): 		Service: "mysql",
// todo(fs): 		Tags:    []string{"master"},
// todo(fs): 		Port:    5000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, "")
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/services", nil)
// todo(fs): 	obj, err := a.srv.AgentServices(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	val := obj.(map[string]*structs.NodeService)
// todo(fs): 	if len(val) != 1 {
// todo(fs): 		t.Fatalf("bad services: %v", obj)
// todo(fs): 	}
// todo(fs): 	if val["mysql"].Port != 5000 {
// todo(fs): 		t.Fatalf("bad service: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Services_ACLFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:      "mysql",
// todo(fs): 		Service: "mysql",
// todo(fs): 		Tags:    []string{"master"},
// todo(fs): 		Port:    5000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, "")
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/services", nil)
// todo(fs): 		obj, err := a.srv.AgentServices(nil, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("Err: %v", err)
// todo(fs): 		}
// todo(fs): 		val := obj.(map[string]*structs.NodeService)
// todo(fs): 		if len(val) != 0 {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/services?token=root", nil)
// todo(fs): 		obj, err := a.srv.AgentServices(nil, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("Err: %v", err)
// todo(fs): 		}
// todo(fs): 		val := obj.(map[string]*structs.NodeService)
// todo(fs): 		if len(val) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Checks(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk1 := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "mysql",
// todo(fs): 		Name:    "mysql",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk1, "")
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/checks", nil)
// todo(fs): 	obj, err := a.srv.AgentChecks(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	val := obj.(map[types.CheckID]*structs.HealthCheck)
// todo(fs): 	if len(val) != 1 {
// todo(fs): 		t.Fatalf("bad checks: %v", obj)
// todo(fs): 	}
// todo(fs): 	if val["mysql"].Status != api.HealthPassing {
// todo(fs): 		t.Fatalf("bad check: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Checks_ACLFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk1 := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "mysql",
// todo(fs): 		Name:    "mysql",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk1, "")
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/checks", nil)
// todo(fs): 		obj, err := a.srv.AgentChecks(nil, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("Err: %v", err)
// todo(fs): 		}
// todo(fs): 		val := obj.(map[types.CheckID]*structs.HealthCheck)
// todo(fs): 		if len(val) != 0 {
// todo(fs): 			t.Fatalf("bad checks: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/checks?token=root", nil)
// todo(fs): 		obj, err := a.srv.AgentChecks(nil, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("Err: %v", err)
// todo(fs): 		}
// todo(fs): 		val := obj.(map[types.CheckID]*structs.HealthCheck)
// todo(fs): 		if len(val) != 1 {
// todo(fs): 			t.Fatalf("bad checks: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Self(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeMeta = map[string]string{"somekey": "somevalue"}
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/self", nil)
// todo(fs): 	obj, err := a.srv.AgentSelf(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	val := obj.(Self)
// todo(fs): 	if int(val.Member.Port) != a.Config.SerfPortLAN {
// todo(fs): 		t.Fatalf("incorrect port: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if int(val.Config.SerfPortLAN) != a.Config.SerfPortLAN {
// todo(fs): 		t.Fatalf("incorrect port: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	cs, err := a.GetLANCoordinate()
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if c := cs[cfg.SegmentName]; !reflect.DeepEqual(c, val.Coord) {
// todo(fs): 		t.Fatalf("coordinates are not equal: %v != %v", c, val.Coord)
// todo(fs): 	}
// todo(fs): 	delete(val.Meta, structs.MetaSegmentKey) // Added later, not in config.
// todo(fs): 	if !reflect.DeepEqual(cfg.NodeMeta, val.Meta) {
// todo(fs): 		t.Fatalf("meta fields are not equal: %v != %v", cfg.NodeMeta, val.Meta)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure there's nothing called "token" that's leaked.
// todo(fs): 	raw, err := a.srv.marshalJSON(req, obj)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if bytes.Contains(bytes.ToLower(raw), []byte("token")) {
// todo(fs): 		t.Fatalf("bad: %s", raw)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Self_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/self", nil)
// todo(fs): 		if _, err := a.srv.AgentSelf(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("agent master token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/self?token=towel", nil)
// todo(fs): 		if _, err := a.srv.AgentSelf(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("read-only token", func(t *testing.T) {
// todo(fs): 		ro := makeReadOnlyAgentACL(t, a.srv)
// todo(fs): 		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/self?token=%s", ro), nil)
// todo(fs): 		if _, err := a.srv.AgentSelf(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Metrics_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/metrics", nil)
// todo(fs): 		if _, err := a.srv.AgentSelf(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("agent master token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/metrics?token=towel", nil)
// todo(fs): 		if _, err := a.srv.AgentSelf(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("read-only token", func(t *testing.T) {
// todo(fs): 		ro := makeReadOnlyAgentACL(t, a.srv)
// todo(fs): 		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/metrics?token=%s", ro), nil)
// todo(fs): 		if _, err := a.srv.AgentSelf(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Reload(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLEnforceVersion8 = false
// todo(fs): 	cfg.Services = []*structs.ServiceDefinition{
// todo(fs): 		&structs.ServiceDefinition{Name: "redis"},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	params := map[string]interface{}{
// todo(fs): 		"datacenter": "dc1",
// todo(fs): 		"type":       "key",
// todo(fs): 		"key":        "test",
// todo(fs): 		"handler":    "true",
// todo(fs): 	}
// todo(fs): 	wp, err := watch.ParseExempt(params, []string{"handler"})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Expected watch.Parse to succeed %v", err)
// todo(fs): 	}
// todo(fs): 	cfg.Watches = append(cfg.Watches, params)
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	if _, ok := a.state.services["redis"]; !ok {
// todo(fs): 		t.Fatalf("missing redis service")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	cfg2 := TestConfig()
// todo(fs): 	cfg2.ACLEnforceVersion8 = false
// todo(fs): 	cfg2.Services = []*structs.ServiceDefinition{
// todo(fs): 		&structs.ServiceDefinition{Name: "redis-reloaded"},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if err := a.ReloadConfig(cfg2); err != nil {
// todo(fs): 		t.Fatalf("got error %v want nil", err)
// todo(fs): 	}
// todo(fs): 	if _, ok := a.state.services["redis-reloaded"]; !ok {
// todo(fs): 		t.Fatalf("missing redis-reloaded service")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, wp := range a.watchPlans {
// todo(fs): 		if !wp.IsStopped() {
// todo(fs): 			t.Fatalf("Reloading configs should stop watch plans of the previous configuration")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Reload_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/reload", nil)
// todo(fs): 		if _, err := a.srv.AgentReload(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("read-only token", func(t *testing.T) {
// todo(fs): 		ro := makeReadOnlyAgentACL(t, a.srv)
// todo(fs): 		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/agent/reload?token=%s", ro), nil)
// todo(fs): 		if _, err := a.srv.AgentReload(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// This proves we call the ACL function, and we've got the other reload
// todo(fs): 	// test to prove we do the reload, which should be sufficient.
// todo(fs): 	// The reload logic is a little complex to set up so isn't worth
// todo(fs): 	// repeating again here.
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Members(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/members", nil)
// todo(fs): 	obj, err := a.srv.AgentMembers(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	val := obj.([]serf.Member)
// todo(fs): 	if len(val) == 0 {
// todo(fs): 		t.Fatalf("bad members: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if int(val[0].Port) != a.Config.SerfPortLAN {
// todo(fs): 		t.Fatalf("not lan: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Members_WAN(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/members?wan=true", nil)
// todo(fs): 	obj, err := a.srv.AgentMembers(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	val := obj.([]serf.Member)
// todo(fs): 	if len(val) == 0 {
// todo(fs): 		t.Fatalf("bad members: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if int(val[0].Port) != a.Config.SerfPortWAN {
// todo(fs): 		t.Fatalf("not wan: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Members_ACLFilter(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/members", nil)
// todo(fs): 		obj, err := a.srv.AgentMembers(nil, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("Err: %v", err)
// todo(fs): 		}
// todo(fs): 		val := obj.([]serf.Member)
// todo(fs): 		if len(val) != 0 {
// todo(fs): 			t.Fatalf("bad members: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/members?token=root", nil)
// todo(fs): 		obj, err := a.srv.AgentMembers(nil, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("Err: %v", err)
// todo(fs): 		}
// todo(fs): 		val := obj.([]serf.Member)
// todo(fs): 		if len(val) != 1 {
// todo(fs): 			t.Fatalf("bad members: %v", obj)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Join(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a1 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a1.Shutdown()
// todo(fs): 	a2 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a2.Config.SerfPortLAN)
// todo(fs): 	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/join/%s", addr), nil)
// todo(fs): 	obj, err := a1.srv.AgentJoin(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("Err: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(a1.LANMembers()) != 2 {
// todo(fs): 		t.Fatalf("should have 2 members")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if got, want := len(a2.LANMembers()), 2; got != want {
// todo(fs): 			r.Fatalf("got %d LAN members want %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Join_WAN(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a1 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a1.Shutdown()
// todo(fs): 	a2 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a2.Config.SerfPortWAN)
// todo(fs): 	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/join/%s?wan=true", addr), nil)
// todo(fs): 	obj, err := a1.srv.AgentJoin(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("Err: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(a1.WANMembers()) != 2 {
// todo(fs): 		t.Fatalf("should have 2 members")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if got, want := len(a2.WANMembers()), 2; got != want {
// todo(fs): 			r.Fatalf("got %d WAN members want %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Join_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a1 := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a1.Shutdown()
// todo(fs): 	a2 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a2.Config.SerfPortLAN)
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/join/%s", addr), nil)
// todo(fs): 		if _, err := a1.srv.AgentJoin(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("agent master token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/join/%s?token=towel", addr), nil)
// todo(fs): 		_, err := a1.srv.AgentJoin(nil, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("read-only token", func(t *testing.T) {
// todo(fs): 		ro := makeReadOnlyAgentACL(t, a1.srv)
// todo(fs): 		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/join/%s?token=%s", addr, ro), nil)
// todo(fs): 		if _, err := a1.srv.AgentJoin(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): type mockNotifier struct{ s string }
// todo(fs):
// todo(fs): func (n *mockNotifier) Notify(state string) error {
// todo(fs): 	n.s = state
// todo(fs): 	return nil
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_JoinLANNotify(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a1 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a1.Shutdown()
// todo(fs):
// todo(fs): 	cfg2 := TestConfig()
// todo(fs): 	cfg2.ServerMode = false
// todo(fs): 	cfg2.Bootstrap = false
// todo(fs): 	a2 := NewTestAgent(t.Name(), cfg2)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	notif := &mockNotifier{}
// todo(fs): 	a1.joinLANNotifier = notif
// todo(fs):
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a2.Config.SerfPortLAN)
// todo(fs): 	_, err := a1.JoinLAN([]string{addr})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if got, want := notif.s, "READY=1"; got != want {
// todo(fs): 		t.Fatalf("got joinLAN notification %q want %q", got, want)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Leave(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a1 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a1.Shutdown()
// todo(fs):
// todo(fs): 	cfg2 := TestConfig()
// todo(fs): 	cfg2.ServerMode = false
// todo(fs): 	cfg2.Bootstrap = false
// todo(fs): 	a2 := NewTestAgent(t.Name(), cfg2)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	// Join first
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a2.Config.SerfPortLAN)
// todo(fs): 	_, err := a1.JoinLAN([]string{addr})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Graceful leave now
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/agent/leave", nil)
// todo(fs): 	obj, err := a2.srv.AgentLeave(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("Err: %v", obj)
// todo(fs): 	}
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		m := a1.LANMembers()
// todo(fs): 		if got, want := m[1].Status, serf.StatusLeft; got != want {
// todo(fs): 			r.Fatalf("got status %q want %q", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Leave_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/leave", nil)
// todo(fs): 		if _, err := a.srv.AgentLeave(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("read-only token", func(t *testing.T) {
// todo(fs): 		ro := makeReadOnlyAgentACL(t, a.srv)
// todo(fs): 		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/agent/leave?token=%s", ro), nil)
// todo(fs): 		if _, err := a.srv.AgentLeave(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// this sub-test will change the state so that there is no leader.
// todo(fs): 	// it must therefore be the last one in this list.
// todo(fs): 	t.Run("agent master token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/leave?token=towel", nil)
// todo(fs): 		if _, err := a.srv.AgentLeave(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_ForceLeave(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a1 := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a1.Shutdown()
// todo(fs): 	a2 := NewTestAgent(t.Name(), nil)
// todo(fs):
// todo(fs): 	// Join first
// todo(fs): 	addr := fmt.Sprintf("127.0.0.1:%d", a2.Config.SerfPortLAN)
// todo(fs): 	_, err := a1.JoinLAN([]string{addr})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// todo(fs): this test probably needs work
// todo(fs): 	a2.Shutdown()
// todo(fs):
// todo(fs): 	// Force leave now
// todo(fs): 	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/force-leave/%s", a2.Config.NodeName), nil)
// todo(fs): 	obj, err := a1.srv.AgentForceLeave(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("Err: %v", obj)
// todo(fs): 	}
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		m := a1.LANMembers()
// todo(fs): 		if got, want := m[1].Status, serf.StatusLeft; got != want {
// todo(fs): 			r.Fatalf("got status %q want %q", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_ForceLeave_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/force-leave/nope", nil)
// todo(fs): 		if _, err := a.srv.AgentForceLeave(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("agent master token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/force-leave/nope?token=towel", nil)
// todo(fs): 		if _, err := a.srv.AgentForceLeave(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("read-only token", func(t *testing.T) {
// todo(fs): 		ro := makeReadOnlyAgentACL(t, a.srv)
// todo(fs): 		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/agent/force-leave/nope?token=%s", ro), nil)
// todo(fs): 		if _, err := a.srv.AgentForceLeave(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.CheckDefinition{
// todo(fs): 		Name: "test",
// todo(fs): 		TTL:  15 * time.Second,
// todo(fs): 	}
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/check/register?token=abc123", jsonReader(args))
// todo(fs): 	obj, err := a.srv.AgentRegisterCheck(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	checkID := types.CheckID("test")
// todo(fs): 	if _, ok := a.state.Checks()[checkID]; !ok {
// todo(fs): 		t.Fatalf("missing test check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if _, ok := a.checkTTLs[checkID]; !ok {
// todo(fs): 		t.Fatalf("missing test check ttl")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the token was configured
// todo(fs): 	if token := a.state.CheckToken(checkID); token == "" {
// todo(fs): 		t.Fatalf("missing token")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// By default, checks start in critical state.
// todo(fs): 	state := a.state.Checks()[checkID]
// todo(fs): 	if state.Status != api.HealthCritical {
// todo(fs): 		t.Fatalf("bad: %v", state)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterCheck_Passing(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.CheckDefinition{
// todo(fs): 		Name:   "test",
// todo(fs): 		TTL:    15 * time.Second,
// todo(fs): 		Status: api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/check/register", jsonReader(args))
// todo(fs): 	obj, err := a.srv.AgentRegisterCheck(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	checkID := types.CheckID("test")
// todo(fs): 	if _, ok := a.state.Checks()[checkID]; !ok {
// todo(fs): 		t.Fatalf("missing test check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if _, ok := a.checkTTLs[checkID]; !ok {
// todo(fs): 		t.Fatalf("missing test check ttl")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	state := a.state.Checks()[checkID]
// todo(fs): 	if state.Status != api.HealthPassing {
// todo(fs): 		t.Fatalf("bad: %v", state)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterCheck_BadStatus(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	args := &structs.CheckDefinition{
// todo(fs): 		Name:   "test",
// todo(fs): 		TTL:    15 * time.Second,
// todo(fs): 		Status: "fluffy",
// todo(fs): 	}
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/check/register", jsonReader(args))
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.AgentRegisterCheck(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("accepted bad status")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterCheck_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	args := &structs.CheckDefinition{
// todo(fs): 		Name: "test",
// todo(fs): 		TTL:  15 * time.Second,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/register", jsonReader(args))
// todo(fs): 		if _, err := a.srv.AgentRegisterCheck(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/register?token=root", jsonReader(args))
// todo(fs): 		if _, err := a.srv.AgentRegisterCheck(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_DeregisterCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	if err := a.AddCheck(chk, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Register node
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/check/deregister/test", nil)
// todo(fs): 	obj, err := a.srv.AgentDeregisterCheck(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	if _, ok := a.state.Checks()["test"]; ok {
// todo(fs): 		t.Fatalf("have test check")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_DeregisterCheckACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	if err := a.AddCheck(chk, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/deregister/test", nil)
// todo(fs): 		if _, err := a.srv.AgentDeregisterCheck(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/deregister/test?token=root", nil)
// todo(fs): 		if _, err := a.srv.AgentDeregisterCheck(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_PassCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/check/pass/test", nil)
// todo(fs): 	obj, err := a.srv.AgentCheckPass(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	state := a.state.Checks()["test"]
// todo(fs): 	if state.Status != api.HealthPassing {
// todo(fs): 		t.Fatalf("bad: %v", state)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_PassCheck_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/pass/test", nil)
// todo(fs): 		if _, err := a.srv.AgentCheckPass(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/pass/test?token=root", nil)
// todo(fs): 		if _, err := a.srv.AgentCheckPass(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_WarnCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/check/warn/test", nil)
// todo(fs): 	obj, err := a.srv.AgentCheckWarn(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	state := a.state.Checks()["test"]
// todo(fs): 	if state.Status != api.HealthWarning {
// todo(fs): 		t.Fatalf("bad: %v", state)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_WarnCheck_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/warn/test", nil)
// todo(fs): 		if _, err := a.srv.AgentCheckWarn(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/warn/test?token=root", nil)
// todo(fs): 		if _, err := a.srv.AgentCheckWarn(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_FailCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/check/fail/test", nil)
// todo(fs): 	obj, err := a.srv.AgentCheckFail(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	state := a.state.Checks()["test"]
// todo(fs): 	if state.Status != api.HealthCritical {
// todo(fs): 		t.Fatalf("bad: %v", state)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_FailCheck_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/fail/test", nil)
// todo(fs): 		if _, err := a.srv.AgentCheckFail(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/check/fail/test?token=root", nil)
// todo(fs): 		if _, err := a.srv.AgentCheckFail(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_UpdateCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	cases := []checkUpdate{
// todo(fs): 		checkUpdate{api.HealthPassing, "hello-passing"},
// todo(fs): 		checkUpdate{api.HealthCritical, "hello-critical"},
// todo(fs): 		checkUpdate{api.HealthWarning, "hello-warning"},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, c := range cases {
// todo(fs): 		t.Run(c.Status, func(t *testing.T) {
// todo(fs): 			req, _ := http.NewRequest("PUT", "/v1/agent/check/update/test", jsonReader(c))
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			obj, err := a.srv.AgentCheckUpdate(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 			if obj != nil {
// todo(fs): 				t.Fatalf("bad: %v", obj)
// todo(fs): 			}
// todo(fs): 			if resp.Code != 200 {
// todo(fs): 				t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 			}
// todo(fs):
// todo(fs): 			state := a.state.Checks()["test"]
// todo(fs): 			if state.Status != c.Status || state.Output != c.Output {
// todo(fs): 				t.Fatalf("bad: %v", state)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("log output limit", func(t *testing.T) {
// todo(fs): 		args := checkUpdate{
// todo(fs): 			Status: api.HealthPassing,
// todo(fs): 			Output: strings.Repeat("-= bad -=", 5*CheckBufSize),
// todo(fs): 		}
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/check/update/test", jsonReader(args))
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.AgentCheckUpdate(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if obj != nil {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Since we append some notes about truncating, we just do a
// todo(fs): 		// rough check that the output buffer was cut down so this test
// todo(fs): 		// isn't super brittle.
// todo(fs): 		state := a.state.Checks()["test"]
// todo(fs): 		if state.Status != api.HealthPassing || len(state.Output) > 2*CheckBufSize {
// todo(fs): 			t.Fatalf("bad: %v", state)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("bogus status", func(t *testing.T) {
// todo(fs): 		args := checkUpdate{Status: "itscomplicated"}
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/check/update/test", jsonReader(args))
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.AgentCheckUpdate(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if obj != nil {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 400 {
// todo(fs): 			t.Fatalf("expected 400, got %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("bogus verb", func(t *testing.T) {
// todo(fs): 		args := checkUpdate{Status: api.HealthPassing}
// todo(fs): 		req, _ := http.NewRequest("POST", "/v1/agent/check/update/test", jsonReader(args))
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.AgentCheckUpdate(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if obj != nil {
// todo(fs): 			t.Fatalf("bad: %v", obj)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 405 {
// todo(fs): 			t.Fatalf("expected 405, got %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_UpdateCheck_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	chk := &structs.HealthCheck{Name: "test", CheckID: "test"}
// todo(fs): 	chkType := &structs.CheckType{TTL: 15 * time.Second}
// todo(fs): 	if err := a.AddCheck(chk, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		args := checkUpdate{api.HealthPassing, "hello-passing"}
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/check/update/test", jsonReader(args))
// todo(fs): 		if _, err := a.srv.AgentCheckUpdate(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		args := checkUpdate{api.HealthPassing, "hello-passing"}
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/check/update/test?token=root", jsonReader(args))
// todo(fs): 		if _, err := a.srv.AgentCheckUpdate(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterService(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	args := &structs.ServiceDefinition{
// todo(fs): 		Name: "test",
// todo(fs): 		Tags: []string{"master"},
// todo(fs): 		Port: 8000,
// todo(fs): 		Check: structs.CheckType{
// todo(fs): 			TTL: 15 * time.Second,
// todo(fs): 		},
// todo(fs): 		Checks: []*structs.CheckType{
// todo(fs): 			&structs.CheckType{
// todo(fs): 				TTL: 20 * time.Second,
// todo(fs): 			},
// todo(fs): 			&structs.CheckType{
// todo(fs): 				TTL: 30 * time.Second,
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/service/register?token=abc123", jsonReader(args))
// todo(fs):
// todo(fs): 	obj, err := a.srv.AgentRegisterService(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the servie
// todo(fs): 	if _, ok := a.state.Services()["test"]; !ok {
// todo(fs): 		t.Fatalf("missing test service")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	checks := a.state.Checks()
// todo(fs): 	if len(checks) != 3 {
// todo(fs): 		t.Fatalf("bad: %v", checks)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if len(a.checkTTLs) != 3 {
// todo(fs): 		t.Fatalf("missing test check ttls: %v", a.checkTTLs)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the token was configured
// todo(fs): 	if token := a.state.ServiceToken("test"); token == "" {
// todo(fs): 		t.Fatalf("missing token")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterService_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	args := &structs.ServiceDefinition{
// todo(fs): 		Name: "test",
// todo(fs): 		Tags: []string{"master"},
// todo(fs): 		Port: 8000,
// todo(fs): 		Check: structs.CheckType{
// todo(fs): 			TTL: 15 * time.Second,
// todo(fs): 		},
// todo(fs): 		Checks: []*structs.CheckType{
// todo(fs): 			&structs.CheckType{
// todo(fs): 				TTL: 20 * time.Second,
// todo(fs): 			},
// todo(fs): 			&structs.CheckType{
// todo(fs): 				TTL: 30 * time.Second,
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/service/register", jsonReader(args))
// todo(fs): 		if _, err := a.srv.AgentRegisterService(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/service/register?token=root", jsonReader(args))
// todo(fs): 		if _, err := a.srv.AgentRegisterService(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterService_InvalidAddress(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	for _, addr := range []string{"0.0.0.0", "::", "[::]"} {
// todo(fs): 		t.Run("addr "+addr, func(t *testing.T) {
// todo(fs): 			args := &structs.ServiceDefinition{
// todo(fs): 				Name:    "test",
// todo(fs): 				Address: addr,
// todo(fs): 				Port:    8000,
// todo(fs): 			}
// todo(fs): 			req, _ := http.NewRequest("GET", "/v1/agent/service/register?token=abc123", jsonReader(args))
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			_, err := a.srv.AgentRegisterService(resp, req)
// todo(fs): 			if err != nil {
// todo(fs): 				t.Fatalf("got error %v want nil", err)
// todo(fs): 			}
// todo(fs): 			if got, want := resp.Code, 400; got != want {
// todo(fs): 				t.Fatalf("got code %d want %d", got, want)
// todo(fs): 			}
// todo(fs): 			if got, want := resp.Body.String(), "Invalid service address"; got != want {
// todo(fs): 				t.Fatalf("got body %q want %q", got, want)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_DeregisterService(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	service := &structs.NodeService{
// todo(fs): 		ID:      "test",
// todo(fs): 		Service: "test",
// todo(fs): 	}
// todo(fs): 	if err := a.AddService(service, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/service/deregister/test", nil)
// todo(fs): 	obj, err := a.srv.AgentDeregisterService(nil, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if obj != nil {
// todo(fs): 		t.Fatalf("bad: %v", obj)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	if _, ok := a.state.Services()["test"]; ok {
// todo(fs): 		t.Fatalf("have test service")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if _, ok := a.state.Checks()["test"]; ok {
// todo(fs): 		t.Fatalf("have test check")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_DeregisterService_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	service := &structs.NodeService{
// todo(fs): 		ID:      "test",
// todo(fs): 		Service: "test",
// todo(fs): 	}
// todo(fs): 	if err := a.AddService(service, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/service/deregister/test", nil)
// todo(fs): 		if _, err := a.srv.AgentDeregisterService(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/service/deregister/test?token=root", nil)
// todo(fs): 		if _, err := a.srv.AgentDeregisterService(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_ServiceMaintenance_BadRequest(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("not PUT", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/agent/service/maintenance/test?enable=true", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.AgentServiceMaintenance(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 405 {
// todo(fs): 			t.Fatalf("expected 405, got %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("not enabled", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/test", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.AgentServiceMaintenance(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 400 {
// todo(fs): 			t.Fatalf("expected 400, got %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("no service id", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/?enable=true", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.AgentServiceMaintenance(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 400 {
// todo(fs): 			t.Fatalf("expected 400, got %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("bad service id", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/_nope_?enable=true", nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		if _, err := a.srv.AgentServiceMaintenance(resp, req); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 404 {
// todo(fs): 			t.Fatalf("expected 404, got %d", resp.Code)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_ServiceMaintenance_Enable(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register the service
// todo(fs): 	service := &structs.NodeService{
// todo(fs): 		ID:      "test",
// todo(fs): 		Service: "test",
// todo(fs): 	}
// todo(fs): 	if err := a.AddService(service, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Force the service into maintenance mode
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/test?enable=true&reason=broken&token=mytoken", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.AgentServiceMaintenance(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the maintenance check was registered
// todo(fs): 	checkID := serviceMaintCheckID("test")
// todo(fs): 	check, ok := a.state.Checks()[checkID]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should have registered maintenance check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the token was added
// todo(fs): 	if token := a.state.CheckToken(checkID); token != "mytoken" {
// todo(fs): 		t.Fatalf("expected 'mytoken', got '%s'", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the reason was set in notes
// todo(fs): 	if check.Notes != "broken" {
// todo(fs): 		t.Fatalf("bad: %#v", check)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_ServiceMaintenance_Disable(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register the service
// todo(fs): 	service := &structs.NodeService{
// todo(fs): 		ID:      "test",
// todo(fs): 		Service: "test",
// todo(fs): 	}
// todo(fs): 	if err := a.AddService(service, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Force the service into maintenance mode
// todo(fs): 	if err := a.EnableServiceMaintenance("test", "", ""); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Leave maintenance mode
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/test?enable=false", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.AgentServiceMaintenance(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the maintenance check was removed
// todo(fs): 	checkID := serviceMaintCheckID("test")
// todo(fs): 	if _, ok := a.state.Checks()[checkID]; ok {
// todo(fs): 		t.Fatalf("should have removed maintenance check")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_ServiceMaintenance_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register the service.
// todo(fs): 	service := &structs.NodeService{
// todo(fs): 		ID:      "test",
// todo(fs): 		Service: "test",
// todo(fs): 	}
// todo(fs): 	if err := a.AddService(service, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/test?enable=true&reason=broken", nil)
// todo(fs): 		if _, err := a.srv.AgentServiceMaintenance(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/service/maintenance/test?enable=true&reason=broken&token=root", nil)
// todo(fs): 		if _, err := a.srv.AgentServiceMaintenance(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_NodeMaintenance_BadRequest(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Fails on non-PUT
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/self/maintenance?enable=true", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.AgentNodeMaintenance(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 405 {
// todo(fs): 		t.Fatalf("expected 405, got %d", resp.Code)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Fails when no enable flag provided
// todo(fs): 	req, _ = http.NewRequest("PUT", "/v1/agent/self/maintenance", nil)
// todo(fs): 	resp = httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.AgentNodeMaintenance(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("expected 400, got %d", resp.Code)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_NodeMaintenance_Enable(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Force the node into maintenance mode
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/agent/self/maintenance?enable=true&reason=broken&token=mytoken", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.AgentNodeMaintenance(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the maintenance check was registered
// todo(fs): 	check, ok := a.state.Checks()[structs.NodeMaint]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("should have registered maintenance check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check that the token was used
// todo(fs): 	if token := a.state.CheckToken(structs.NodeMaint); token != "mytoken" {
// todo(fs): 		t.Fatalf("expected 'mytoken', got '%s'", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the reason was set in notes
// todo(fs): 	if check.Notes != "broken" {
// todo(fs): 		t.Fatalf("bad: %#v", check)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_NodeMaintenance_Disable(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Force the node into maintenance mode
// todo(fs): 	a.EnableNodeMaintenance("", "")
// todo(fs):
// todo(fs): 	// Leave maintenance mode
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/agent/self/maintenance?enable=false", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.AgentNodeMaintenance(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("expected 200, got %d", resp.Code)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the maintenance check was removed
// todo(fs): 	if _, ok := a.state.Checks()[structs.NodeMaint]; ok {
// todo(fs): 		t.Fatalf("should have removed maintenance check")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_NodeMaintenance_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	t.Run("no token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/self/maintenance?enable=true&reason=broken", nil)
// todo(fs): 		if _, err := a.srv.AgentNodeMaintenance(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("root token", func(t *testing.T) {
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/self/maintenance?enable=true&reason=broken&token=root", nil)
// todo(fs): 		if _, err := a.srv.AgentNodeMaintenance(nil, req); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RegisterCheck_Service(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	args := &structs.ServiceDefinition{
// todo(fs): 		Name: "memcache",
// todo(fs): 		Port: 8000,
// todo(fs): 		Check: structs.CheckType{
// todo(fs): 			TTL: 15 * time.Second,
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// First register the service
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/service/register", jsonReader(args))
// todo(fs): 	if _, err := a.srv.AgentRegisterService(nil, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Now register an additional check
// todo(fs): 	checkArgs := &structs.CheckDefinition{
// todo(fs): 		Name:      "memcache_check2",
// todo(fs): 		ServiceID: "memcache",
// todo(fs): 		TTL:       15 * time.Second,
// todo(fs): 	}
// todo(fs): 	req, _ = http.NewRequest("GET", "/v1/agent/check/register", jsonReader(checkArgs))
// todo(fs): 	if _, err := a.srv.AgentRegisterCheck(nil, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	result := a.state.Checks()
// todo(fs): 	if _, ok := result["service:memcache"]; !ok {
// todo(fs): 		t.Fatalf("missing memcached check")
// todo(fs): 	}
// todo(fs): 	if _, ok := result["memcache_check2"]; !ok {
// todo(fs): 		t.Fatalf("missing memcache_check2 check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure the new check is associated with the service
// todo(fs): 	if result["memcache_check2"].ServiceID != "memcache" {
// todo(fs): 		t.Fatalf("bad: %#v", result["memcached_check2"])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Monitor(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	logWriter := logger.NewLogWriter(512)
// todo(fs): 	a := &TestAgent{
// todo(fs): 		Name:      t.Name(),
// todo(fs): 		LogWriter: logWriter,
// todo(fs): 		LogOutput: io.MultiWriter(os.Stderr, logWriter),
// todo(fs): 	}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Try passing an invalid log level
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/monitor?loglevel=invalid", nil)
// todo(fs): 	resp := newClosableRecorder()
// todo(fs): 	if _, err := a.srv.AgentMonitor(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 400 {
// todo(fs): 		t.Fatalf("bad: %v", resp.Code)
// todo(fs): 	}
// todo(fs): 	body, _ := ioutil.ReadAll(resp.Body)
// todo(fs): 	if !strings.Contains(string(body), "Unknown log level") {
// todo(fs): 		t.Fatalf("bad: %s", body)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try to stream logs until we see the expected log line
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		req, _ = http.NewRequest("GET", "/v1/agent/monitor?loglevel=debug", nil)
// todo(fs): 		resp = newClosableRecorder()
// todo(fs): 		done := make(chan struct{})
// todo(fs): 		go func() {
// todo(fs): 			if _, err := a.srv.AgentMonitor(resp, req); err != nil {
// todo(fs): 				t.Fatalf("err: %s", err)
// todo(fs): 			}
// todo(fs): 			close(done)
// todo(fs): 		}()
// todo(fs):
// todo(fs): 		resp.Close()
// todo(fs): 		<-done
// todo(fs):
// todo(fs): 		got := resp.Body.Bytes()
// todo(fs): 		want := []byte("raft: Initial configuration (index=1)")
// todo(fs): 		if !bytes.Contains(got, want) {
// todo(fs): 			r.Fatalf("got %q and did not find %q", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): type closableRecorder struct {
// todo(fs): 	*httptest.ResponseRecorder
// todo(fs): 	closer chan bool
// todo(fs): }
// todo(fs):
// todo(fs): func newClosableRecorder() *closableRecorder {
// todo(fs): 	r := httptest.NewRecorder()
// todo(fs): 	closer := make(chan bool)
// todo(fs): 	return &closableRecorder{r, closer}
// todo(fs): }
// todo(fs):
// todo(fs): func (r *closableRecorder) Close() {
// todo(fs): 	close(r.closer)
// todo(fs): }
// todo(fs):
// todo(fs): func (r *closableRecorder) CloseNotify() <-chan bool {
// todo(fs): 	return r.closer
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Monitor_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), TestACLConfig())
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Try without a token.
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/agent/monitor", nil)
// todo(fs): 	if _, err := a.srv.AgentMonitor(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// This proves we call the ACL function, and we've got the other monitor
// todo(fs): 	// test to prove monitor works, which should be sufficient. The monitor
// todo(fs): 	// logic is a little complex to set up so isn't worth repeating again
// todo(fs): 	// here.
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_Token(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestACLConfig()
// todo(fs): 	cfg.ACLToken = ""
// todo(fs): 	cfg.ACLAgentToken = ""
// todo(fs): 	cfg.ACLAgentMasterToken = ""
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	type tokens struct {
// todo(fs): 		user, agent, master, repl string
// todo(fs): 	}
// todo(fs):
// todo(fs): 	resetTokens := func(got tokens) {
// todo(fs): 		a.tokens.UpdateUserToken(got.user)
// todo(fs): 		a.tokens.UpdateAgentToken(got.agent)
// todo(fs): 		a.tokens.UpdateAgentMasterToken(got.master)
// todo(fs): 		a.tokens.UpdateACLReplicationToken(got.repl)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	body := func(token string) io.Reader {
// todo(fs): 		return jsonReader(&api.AgentToken{Token: token})
// todo(fs): 	}
// todo(fs):
// todo(fs): 	badJSON := func() io.Reader {
// todo(fs): 		return jsonReader(false)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	tests := []struct {
// todo(fs): 		name        string
// todo(fs): 		method, url string
// todo(fs): 		body        io.Reader
// todo(fs): 		code        int
// todo(fs): 		got, want   tokens
// todo(fs): 	}{
// todo(fs): 		{
// todo(fs): 			name:   "bad method",
// todo(fs): 			method: "GET",
// todo(fs): 			url:    "acl_token",
// todo(fs): 			code:   http.StatusMethodNotAllowed,
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "bad token name",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "nope?token=root",
// todo(fs): 			body:   body("X"),
// todo(fs): 			code:   http.StatusNotFound,
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "bad JSON",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_token?token=root",
// todo(fs): 			body:   badJSON(),
// todo(fs): 			code:   http.StatusBadRequest,
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "set user",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_token?token=root",
// todo(fs): 			body:   body("U"),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			want:   tokens{user: "U", agent: "U"},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "set agent",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_agent_token?token=root",
// todo(fs): 			body:   body("A"),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			got:    tokens{user: "U", agent: "U"},
// todo(fs): 			want:   tokens{user: "U", agent: "A"},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "set master",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_agent_master_token?token=root",
// todo(fs): 			body:   body("M"),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			want:   tokens{master: "M"},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "set repl",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_replication_token?token=root",
// todo(fs): 			body:   body("R"),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			want:   tokens{repl: "R"},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "clear user",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_token?token=root",
// todo(fs): 			body:   body(""),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			got:    tokens{user: "U"},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "clear agent",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_agent_token?token=root",
// todo(fs): 			body:   body(""),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			got:    tokens{agent: "A"},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "clear master",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_agent_master_token?token=root",
// todo(fs): 			body:   body(""),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			got:    tokens{master: "M"},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			name:   "clear repl",
// todo(fs): 			method: "PUT",
// todo(fs): 			url:    "acl_replication_token?token=root",
// todo(fs): 			body:   body(""),
// todo(fs): 			code:   http.StatusOK,
// todo(fs): 			got:    tokens{repl: "R"},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	for _, tt := range tests {
// todo(fs): 		t.Run(tt.name, func(t *testing.T) {
// todo(fs): 			resetTokens(tt.got)
// todo(fs): 			url := fmt.Sprintf("/v1/agent/token/%s", tt.url)
// todo(fs): 			resp := httptest.NewRecorder()
// todo(fs): 			req, _ := http.NewRequest(tt.method, url, tt.body)
// todo(fs): 			if _, err := a.srv.AgentToken(resp, req); err != nil {
// todo(fs): 				t.Fatalf("err: %v", err)
// todo(fs): 			}
// todo(fs): 			if got, want := resp.Code, tt.code; got != want {
// todo(fs): 				t.Fatalf("got %d want %d", got, want)
// todo(fs): 			}
// todo(fs): 			if got, want := a.tokens.UserToken(), tt.want.user; got != want {
// todo(fs): 				t.Fatalf("got %q want %q", got, want)
// todo(fs): 			}
// todo(fs): 			if got, want := a.tokens.AgentToken(), tt.want.agent; got != want {
// todo(fs): 				t.Fatalf("got %q want %q", got, want)
// todo(fs): 			}
// todo(fs): 			if tt.want.master != "" && !a.tokens.IsAgentMasterToken(tt.want.master) {
// todo(fs): 				t.Fatalf("%q should be the master token", tt.want.master)
// todo(fs): 			}
// todo(fs): 			if got, want := a.tokens.ACLReplicationToken(), tt.want.repl; got != want {
// todo(fs): 				t.Fatalf("got %q want %q", got, want)
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// This one returns an error that is interpreted by the HTTP wrapper, so
// todo(fs): 	// doesn't fit into our table above.
// todo(fs): 	t.Run("permission denied", func(t *testing.T) {
// todo(fs): 		resetTokens(tokens{})
// todo(fs): 		req, _ := http.NewRequest("PUT", "/v1/agent/token/acl_token", body("X"))
// todo(fs): 		if _, err := a.srv.AgentToken(nil, req); !acl.IsErrPermissionDenied(err) {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if got, want := a.tokens.UserToken(), ""; got != want {
// todo(fs): 			t.Fatalf("got %q want %q", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
