package agent

// todo(fs): func TestAgentAntiEntropy_Services(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := &TestAgent{Name: t.Name(), NoInitialSync: true}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register info
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists both, same (noop)
// todo(fs): 	var out struct{}
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:      "mysql",
// todo(fs): 		Service: "mysql",
// todo(fs): 		Tags:    []string{"master"},
// todo(fs): 		Port:    5000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, "")
// todo(fs): 	args.Service = srv1
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists both, different (update)
// todo(fs): 	srv2 := &structs.NodeService{
// todo(fs): 		ID:      "redis",
// todo(fs): 		Service: "redis",
// todo(fs): 		Tags:    []string{},
// todo(fs): 		Port:    8000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv2, "")
// todo(fs):
// todo(fs): 	srv2_mod := new(structs.NodeService)
// todo(fs): 	*srv2_mod = *srv2
// todo(fs): 	srv2_mod.Port = 9000
// todo(fs): 	args.Service = srv2_mod
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists local (create)
// todo(fs): 	srv3 := &structs.NodeService{
// todo(fs): 		ID:      "web",
// todo(fs): 		Service: "web",
// todo(fs): 		Tags:    []string{},
// todo(fs): 		Port:    80,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv3, "")
// todo(fs):
// todo(fs): 	// Exists remote (delete)
// todo(fs): 	srv4 := &structs.NodeService{
// todo(fs): 		ID:      "lb",
// todo(fs): 		Service: "lb",
// todo(fs): 		Tags:    []string{},
// todo(fs): 		Port:    443,
// todo(fs): 	}
// todo(fs): 	args.Service = srv4
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists both, different address (update)
// todo(fs): 	srv5 := &structs.NodeService{
// todo(fs): 		ID:      "api",
// todo(fs): 		Service: "api",
// todo(fs): 		Tags:    []string{},
// todo(fs): 		Address: "127.0.0.10",
// todo(fs): 		Port:    8000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv5, "")
// todo(fs):
// todo(fs): 	srv5_mod := new(structs.NodeService)
// todo(fs): 	*srv5_mod = *srv5
// todo(fs): 	srv5_mod.Address = "127.0.0.1"
// todo(fs): 	args.Service = srv5_mod
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists local, in sync, remote missing (create)
// todo(fs): 	srv6 := &structs.NodeService{
// todo(fs): 		ID:      "cache",
// todo(fs): 		Service: "cache",
// todo(fs): 		Tags:    []string{},
// todo(fs): 		Port:    11211,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv6, "")
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	a.state.Lock()
// todo(fs): 	a.state.serviceStatus["cache"] = syncStatus{inSync: true}
// todo(fs): 	a.state.Unlock()
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs):
// todo(fs): 	var services structs.IndexedNodeServices
// todo(fs): 	req := structs.NodeSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Make sure we sent along our node info when we synced.
// todo(fs): 		id := services.NodeServices.Node.ID
// todo(fs): 		addrs := services.NodeServices.Node.TaggedAddresses
// todo(fs): 		meta := services.NodeServices.Node.Meta
// todo(fs): 		delete(meta, structs.MetaSegmentKey) // Added later, not in config.
// todo(fs): 		if id != a.Config.NodeID ||
// todo(fs): 			!reflect.DeepEqual(addrs, a.Config.TaggedAddresses) ||
// todo(fs): 			!reflect.DeepEqual(meta, a.Config.NodeMeta) {
// todo(fs): 			r.Fatalf("bad: %v", services.NodeServices.Node)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 6 services (consul included)
// todo(fs): 		if len(services.NodeServices.Services) != 6 {
// todo(fs): 			r.Fatalf("bad: %v", services.NodeServices.Services)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the services should match
// todo(fs): 		for id, serv := range services.NodeServices.Services {
// todo(fs): 			serv.CreateIndex, serv.ModifyIndex = 0, 0
// todo(fs): 			switch id {
// todo(fs): 			case "mysql":
// todo(fs): 				if !reflect.DeepEqual(serv, srv1) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv1)
// todo(fs): 				}
// todo(fs): 			case "redis":
// todo(fs): 				if !reflect.DeepEqual(serv, srv2) {
// todo(fs): 					r.Fatalf("bad: %#v %#v", serv, srv2)
// todo(fs): 				}
// todo(fs): 			case "web":
// todo(fs): 				if !reflect.DeepEqual(serv, srv3) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv3)
// todo(fs): 				}
// todo(fs): 			case "api":
// todo(fs): 				if !reflect.DeepEqual(serv, srv5) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv5)
// todo(fs): 				}
// todo(fs): 			case "cache":
// todo(fs): 				if !reflect.DeepEqual(serv, srv6) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv6)
// todo(fs): 				}
// todo(fs): 			case structs.ConsulServiceID:
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				r.Fatalf("unexpected service: %v", id)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		// Check the local state
// todo(fs): 		if len(a.state.services) != 5 {
// todo(fs): 			r.Fatalf("bad: %v", a.state.services)
// todo(fs): 		}
// todo(fs): 		if len(a.state.serviceStatus) != 5 {
// todo(fs): 			r.Fatalf("bad: %v", a.state.serviceStatus)
// todo(fs): 		}
// todo(fs): 		for name, status := range a.state.serviceStatus {
// todo(fs): 			if !status.inSync {
// todo(fs): 				r.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Remove one of the services
// todo(fs): 	a.state.RemoveService("api")
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 5 services (consul included)
// todo(fs): 		if len(services.NodeServices.Services) != 5 {
// todo(fs): 			r.Fatalf("bad: %v", services.NodeServices.Services)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the services should match
// todo(fs): 		for id, serv := range services.NodeServices.Services {
// todo(fs): 			serv.CreateIndex, serv.ModifyIndex = 0, 0
// todo(fs): 			switch id {
// todo(fs): 			case "mysql":
// todo(fs): 				if !reflect.DeepEqual(serv, srv1) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv1)
// todo(fs): 				}
// todo(fs): 			case "redis":
// todo(fs): 				if !reflect.DeepEqual(serv, srv2) {
// todo(fs): 					r.Fatalf("bad: %#v %#v", serv, srv2)
// todo(fs): 				}
// todo(fs): 			case "web":
// todo(fs): 				if !reflect.DeepEqual(serv, srv3) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv3)
// todo(fs): 				}
// todo(fs): 			case "cache":
// todo(fs): 				if !reflect.DeepEqual(serv, srv6) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv6)
// todo(fs): 				}
// todo(fs): 			case structs.ConsulServiceID:
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				r.Fatalf("unexpected service: %v", id)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		// Check the local state
// todo(fs): 		if len(a.state.services) != 4 {
// todo(fs): 			r.Fatalf("bad: %v", a.state.services)
// todo(fs): 		}
// todo(fs): 		if len(a.state.serviceStatus) != 4 {
// todo(fs): 			r.Fatalf("bad: %v", a.state.serviceStatus)
// todo(fs): 		}
// todo(fs): 		for name, status := range a.state.serviceStatus {
// todo(fs): 			if !status.inSync {
// todo(fs): 				r.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_EnableTagOverride(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := &TestAgent{Name: t.Name(), NoInitialSync: true}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs):
// todo(fs): 	// EnableTagOverride = true
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:                "svc_id1",
// todo(fs): 		Service:           "svc1",
// todo(fs): 		Tags:              []string{"tag1"},
// todo(fs): 		Port:              6100,
// todo(fs): 		EnableTagOverride: true,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, "")
// todo(fs): 	srv1_mod := new(structs.NodeService)
// todo(fs): 	*srv1_mod = *srv1
// todo(fs): 	srv1_mod.Port = 7100
// todo(fs): 	srv1_mod.Tags = []string{"tag1_mod"}
// todo(fs): 	args.Service = srv1_mod
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// EnableTagOverride = false
// todo(fs): 	srv2 := &structs.NodeService{
// todo(fs): 		ID:                "svc_id2",
// todo(fs): 		Service:           "svc2",
// todo(fs): 		Tags:              []string{"tag2"},
// todo(fs): 		Port:              6200,
// todo(fs): 		EnableTagOverride: false,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv2, "")
// todo(fs): 	srv2_mod := new(structs.NodeService)
// todo(fs): 	*srv2_mod = *srv2
// todo(fs): 	srv2_mod.Port = 7200
// todo(fs): 	srv2_mod.Tags = []string{"tag2_mod"}
// todo(fs): 	args.Service = srv2_mod
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs):
// todo(fs): 	req := structs.NodeSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 	}
// todo(fs): 	var services structs.IndexedNodeServices
// todo(fs):
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		//	runtime.Gosched()
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		// All the services should match
// todo(fs): 		for id, serv := range services.NodeServices.Services {
// todo(fs): 			serv.CreateIndex, serv.ModifyIndex = 0, 0
// todo(fs): 			switch id {
// todo(fs): 			case "svc_id1":
// todo(fs): 				if serv.ID != "svc_id1" ||
// todo(fs): 					serv.Service != "svc1" ||
// todo(fs): 					serv.Port != 6100 ||
// todo(fs): 					!reflect.DeepEqual(serv.Tags, []string{"tag1_mod"}) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv1)
// todo(fs): 				}
// todo(fs): 			case "svc_id2":
// todo(fs): 				if serv.ID != "svc_id2" ||
// todo(fs): 					serv.Service != "svc2" ||
// todo(fs): 					serv.Port != 6200 ||
// todo(fs): 					!reflect.DeepEqual(serv.Tags, []string{"tag2"}) {
// todo(fs): 					r.Fatalf("bad: %v %v", serv, srv2)
// todo(fs): 				}
// todo(fs): 			case structs.ConsulServiceID:
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				r.Fatalf("unexpected service: %v", id)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		for name, status := range a.state.serviceStatus {
// todo(fs): 			if !status.inSync {
// todo(fs): 				r.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_Services_WithChecks(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	{
// todo(fs): 		// Single check
// todo(fs): 		srv := &structs.NodeService{
// todo(fs): 			ID:      "mysql",
// todo(fs): 			Service: "mysql",
// todo(fs): 			Tags:    []string{"master"},
// todo(fs): 			Port:    5000,
// todo(fs): 		}
// todo(fs): 		a.state.AddService(srv, "")
// todo(fs):
// todo(fs): 		chk := &structs.HealthCheck{
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			CheckID:   "mysql",
// todo(fs): 			Name:      "mysql",
// todo(fs): 			ServiceID: "mysql",
// todo(fs): 			Status:    api.HealthPassing,
// todo(fs): 		}
// todo(fs): 		a.state.AddCheck(chk, "")
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		func() {
// todo(fs): 			a.state.RLock()
// todo(fs): 			defer a.state.RUnlock()
// todo(fs):
// todo(fs): 			// Sync the service once
// todo(fs): 			if err := a.state.syncService("mysql"); err != nil {
// todo(fs): 				t.Fatalf("err: %s", err)
// todo(fs): 			}
// todo(fs): 		}()
// todo(fs):
// todo(fs): 		// We should have 2 services (consul included)
// todo(fs): 		svcReq := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 		}
// todo(fs): 		var services structs.IndexedNodeServices
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &svcReq, &services); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if len(services.NodeServices.Services) != 2 {
// todo(fs): 			t.Fatalf("bad: %v", services.NodeServices.Services)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have one health check
// todo(fs): 		chkReq := structs.ServiceSpecificRequest{
// todo(fs): 			Datacenter:  "dc1",
// todo(fs): 			ServiceName: "mysql",
// todo(fs): 		}
// todo(fs): 		var checks structs.IndexedHealthChecks
// todo(fs): 		if err := a.RPC("Health.ServiceChecks", &chkReq, &checks); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if len(checks.HealthChecks) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", checks)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	{
// todo(fs): 		// Multiple checks
// todo(fs): 		srv := &structs.NodeService{
// todo(fs): 			ID:      "redis",
// todo(fs): 			Service: "redis",
// todo(fs): 			Tags:    []string{"master"},
// todo(fs): 			Port:    5000,
// todo(fs): 		}
// todo(fs): 		a.state.AddService(srv, "")
// todo(fs):
// todo(fs): 		chk1 := &structs.HealthCheck{
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			CheckID:   "redis:1",
// todo(fs): 			Name:      "redis:1",
// todo(fs): 			ServiceID: "redis",
// todo(fs): 			Status:    api.HealthPassing,
// todo(fs): 		}
// todo(fs): 		a.state.AddCheck(chk1, "")
// todo(fs):
// todo(fs): 		chk2 := &structs.HealthCheck{
// todo(fs): 			Node:      a.Config.NodeName,
// todo(fs): 			CheckID:   "redis:2",
// todo(fs): 			Name:      "redis:2",
// todo(fs): 			ServiceID: "redis",
// todo(fs): 			Status:    api.HealthPassing,
// todo(fs): 		}
// todo(fs): 		a.state.AddCheck(chk2, "")
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		func() {
// todo(fs): 			a.state.RLock()
// todo(fs): 			defer a.state.RUnlock()
// todo(fs):
// todo(fs): 			// Sync the service once
// todo(fs): 			if err := a.state.syncService("redis"); err != nil {
// todo(fs): 				t.Fatalf("err: %s", err)
// todo(fs): 			}
// todo(fs): 		}()
// todo(fs):
// todo(fs): 		// We should have 3 services (consul included)
// todo(fs): 		svcReq := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 		}
// todo(fs): 		var services structs.IndexedNodeServices
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &svcReq, &services); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if len(services.NodeServices.Services) != 3 {
// todo(fs): 			t.Fatalf("bad: %v", services.NodeServices.Services)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have two health checks
// todo(fs): 		chkReq := structs.ServiceSpecificRequest{
// todo(fs): 			Datacenter:  "dc1",
// todo(fs): 			ServiceName: "redis",
// todo(fs): 		}
// todo(fs): 		var checks structs.IndexedHealthChecks
// todo(fs): 		if err := a.RPC("Health.ServiceChecks", &chkReq, &checks); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if len(checks.HealthChecks) != 2 {
// todo(fs): 			t.Fatalf("bad: %v", checks)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): var testRegisterRules = `
// todo(fs): node "" {
// todo(fs): 	policy = "write"
// todo(fs): }
// todo(fs):
// todo(fs): service "api" {
// todo(fs): 	policy = "write"
// todo(fs): }
// todo(fs):
// todo(fs): service "consul" {
// todo(fs): 	policy = "write"
// todo(fs): }
// todo(fs): `
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_Services_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs): 	a := &TestAgent{Name: t.Name(), Config: cfg, NoInitialSync: true}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create the ACL
// todo(fs): 	arg := structs.ACLRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Op:         structs.ACLSet,
// todo(fs): 		ACL: structs.ACL{
// todo(fs): 			Name:  "User token",
// todo(fs): 			Type:  structs.ACLTypeClient,
// todo(fs): 			Rules: testRegisterRules,
// todo(fs): 		},
// todo(fs): 		WriteRequest: structs.WriteRequest{
// todo(fs): 			Token: "root",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var token string
// todo(fs): 	if err := a.RPC("ACL.Apply", &arg, &token); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create service (disallowed)
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:      "mysql",
// todo(fs): 		Service: "mysql",
// todo(fs): 		Tags:    []string{"master"},
// todo(fs): 		Port:    5000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, token)
// todo(fs):
// todo(fs): 	// Create service (allowed)
// todo(fs): 	srv2 := &structs.NodeService{
// todo(fs): 		ID:      "api",
// todo(fs): 		Service: "api",
// todo(fs): 		Tags:    []string{"foo"},
// todo(fs): 		Port:    5001,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv2, token)
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs): 	time.Sleep(200 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	{
// todo(fs): 		req := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 			QueryOptions: structs.QueryOptions{
// todo(fs): 				Token: "root",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var services structs.IndexedNodeServices
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 2 services (consul included)
// todo(fs): 		if len(services.NodeServices.Services) != 2 {
// todo(fs): 			t.Fatalf("bad: %v", services.NodeServices.Services)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the services should match
// todo(fs): 		for id, serv := range services.NodeServices.Services {
// todo(fs): 			serv.CreateIndex, serv.ModifyIndex = 0, 0
// todo(fs): 			switch id {
// todo(fs): 			case "mysql":
// todo(fs): 				t.Fatalf("should not be permitted")
// todo(fs): 			case "api":
// todo(fs): 				if !reflect.DeepEqual(serv, srv2) {
// todo(fs): 					t.Fatalf("bad: %#v %#v", serv, srv2)
// todo(fs): 				}
// todo(fs): 			case structs.ConsulServiceID:
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				t.Fatalf("unexpected service: %v", id)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		func() {
// todo(fs): 			a.state.RLock()
// todo(fs): 			defer a.state.RUnlock()
// todo(fs):
// todo(fs): 			// Check the local state
// todo(fs): 			if len(a.state.services) != 2 {
// todo(fs): 				t.Fatalf("bad: %v", a.state.services)
// todo(fs): 			}
// todo(fs): 			if len(a.state.serviceStatus) != 2 {
// todo(fs): 				t.Fatalf("bad: %v", a.state.serviceStatus)
// todo(fs): 			}
// todo(fs): 			for name, status := range a.state.serviceStatus {
// todo(fs): 				if !status.inSync {
// todo(fs): 					t.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 		}()
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Now remove the service and re-sync
// todo(fs): 	a.state.RemoveService("api")
// todo(fs): 	a.StartSync()
// todo(fs): 	time.Sleep(200 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	{
// todo(fs): 		req := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 			QueryOptions: structs.QueryOptions{
// todo(fs): 				Token: "root",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var services structs.IndexedNodeServices
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 1 service (just consul)
// todo(fs): 		if len(services.NodeServices.Services) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", services.NodeServices.Services)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the services should match
// todo(fs): 		for id, serv := range services.NodeServices.Services {
// todo(fs): 			serv.CreateIndex, serv.ModifyIndex = 0, 0
// todo(fs): 			switch id {
// todo(fs): 			case "mysql":
// todo(fs): 				t.Fatalf("should not be permitted")
// todo(fs): 			case "api":
// todo(fs): 				t.Fatalf("should be deleted")
// todo(fs): 			case structs.ConsulServiceID:
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				t.Fatalf("unexpected service: %v", id)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		func() {
// todo(fs): 			a.state.RLock()
// todo(fs): 			defer a.state.RUnlock()
// todo(fs):
// todo(fs): 			// Check the local state
// todo(fs): 			if len(a.state.services) != 1 {
// todo(fs): 				t.Fatalf("bad: %v", a.state.services)
// todo(fs): 			}
// todo(fs): 			if len(a.state.serviceStatus) != 1 {
// todo(fs): 				t.Fatalf("bad: %v", a.state.serviceStatus)
// todo(fs): 			}
// todo(fs): 			for name, status := range a.state.serviceStatus {
// todo(fs): 				if !status.inSync {
// todo(fs): 					t.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 		}()
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure the token got cleaned up.
// todo(fs): 	if token := a.state.ServiceToken("api"); token != "" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_Checks(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := &TestAgent{Name: t.Name(), NoInitialSync: true}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register info
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists both, same (noop)
// todo(fs): 	var out struct{}
// todo(fs): 	chk1 := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "mysql",
// todo(fs): 		Name:    "mysql",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk1, "")
// todo(fs): 	args.Check = chk1
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists both, different (update)
// todo(fs): 	chk2 := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "redis",
// todo(fs): 		Name:    "redis",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk2, "")
// todo(fs):
// todo(fs): 	chk2_mod := new(structs.HealthCheck)
// todo(fs): 	*chk2_mod = *chk2
// todo(fs): 	chk2_mod.Status = api.HealthCritical
// todo(fs): 	args.Check = chk2_mod
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists local (create)
// todo(fs): 	chk3 := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "web",
// todo(fs): 		Name:    "web",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk3, "")
// todo(fs):
// todo(fs): 	// Exists remote (delete)
// todo(fs): 	chk4 := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "lb",
// todo(fs): 		Name:    "lb",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	args.Check = chk4
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Exists local, in sync, remote missing (create)
// todo(fs): 	chk5 := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "cache",
// todo(fs): 		Name:    "cache",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk5, "")
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	a.state.Lock()
// todo(fs): 	a.state.checkStatus["cache"] = syncStatus{inSync: true}
// todo(fs): 	a.state.Unlock()
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs):
// todo(fs): 	req := structs.NodeSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 	}
// todo(fs): 	var checks structs.IndexedHealthChecks
// todo(fs):
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 5 checks (serf included)
// todo(fs): 		if len(checks.HealthChecks) != 5 {
// todo(fs): 			r.Fatalf("bad: %v", checks)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the checks should match
// todo(fs): 		for _, chk := range checks.HealthChecks {
// todo(fs): 			chk.CreateIndex, chk.ModifyIndex = 0, 0
// todo(fs): 			switch chk.CheckID {
// todo(fs): 			case "mysql":
// todo(fs): 				if !reflect.DeepEqual(chk, chk1) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk1)
// todo(fs): 				}
// todo(fs): 			case "redis":
// todo(fs): 				if !reflect.DeepEqual(chk, chk2) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk2)
// todo(fs): 				}
// todo(fs): 			case "web":
// todo(fs): 				if !reflect.DeepEqual(chk, chk3) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk3)
// todo(fs): 				}
// todo(fs): 			case "cache":
// todo(fs): 				if !reflect.DeepEqual(chk, chk5) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk5)
// todo(fs): 				}
// todo(fs): 			case "serfHealth":
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				r.Fatalf("unexpected check: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	func() {
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		// Check the local state
// todo(fs): 		if len(a.state.checks) != 4 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checks)
// todo(fs): 		}
// todo(fs): 		if len(a.state.checkStatus) != 4 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checkStatus)
// todo(fs): 		}
// todo(fs): 		for name, status := range a.state.checkStatus {
// todo(fs): 			if !status.inSync {
// todo(fs): 				t.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs):
// todo(fs): 	// Make sure we sent along our node info addresses when we synced.
// todo(fs): 	{
// todo(fs): 		req := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 		}
// todo(fs): 		var services structs.IndexedNodeServices
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		id := services.NodeServices.Node.ID
// todo(fs): 		addrs := services.NodeServices.Node.TaggedAddresses
// todo(fs): 		meta := services.NodeServices.Node.Meta
// todo(fs): 		delete(meta, structs.MetaSegmentKey) // Added later, not in config.
// todo(fs): 		if id != a.Config.NodeID ||
// todo(fs): 			!reflect.DeepEqual(addrs, a.Config.TaggedAddresses) ||
// todo(fs): 			!reflect.DeepEqual(meta, a.Config.NodeMeta) {
// todo(fs): 			t.Fatalf("bad: %v", services.NodeServices.Node)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Remove one of the checks
// todo(fs): 	a.state.RemoveCheck("redis")
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs):
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 5 checks (serf included)
// todo(fs): 		if len(checks.HealthChecks) != 4 {
// todo(fs): 			r.Fatalf("bad: %v", checks)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the checks should match
// todo(fs): 		for _, chk := range checks.HealthChecks {
// todo(fs): 			chk.CreateIndex, chk.ModifyIndex = 0, 0
// todo(fs): 			switch chk.CheckID {
// todo(fs): 			case "mysql":
// todo(fs): 				if !reflect.DeepEqual(chk, chk1) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk1)
// todo(fs): 				}
// todo(fs): 			case "web":
// todo(fs): 				if !reflect.DeepEqual(chk, chk3) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk3)
// todo(fs): 				}
// todo(fs): 			case "cache":
// todo(fs): 				if !reflect.DeepEqual(chk, chk5) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk5)
// todo(fs): 				}
// todo(fs): 			case "serfHealth":
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				r.Fatalf("unexpected check: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	func() {
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		// Check the local state
// todo(fs): 		if len(a.state.checks) != 3 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checks)
// todo(fs): 		}
// todo(fs): 		if len(a.state.checkStatus) != 3 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checkStatus)
// todo(fs): 		}
// todo(fs): 		for name, status := range a.state.checkStatus {
// todo(fs): 			if !status.inSync {
// todo(fs): 				t.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_Checks_ACLDeny(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ACLDatacenter = "dc1"
// todo(fs): 	cfg.ACLMasterToken = "root"
// todo(fs): 	cfg.ACLDefaultPolicy = "deny"
// todo(fs): 	cfg.ACLEnforceVersion8 = true
// todo(fs): 	a := &TestAgent{Name: t.Name(), Config: cfg, NoInitialSync: true}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create the ACL
// todo(fs): 	arg := structs.ACLRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Op:         structs.ACLSet,
// todo(fs): 		ACL: structs.ACL{
// todo(fs): 			Name:  "User token",
// todo(fs): 			Type:  structs.ACLTypeClient,
// todo(fs): 			Rules: testRegisterRules,
// todo(fs): 		},
// todo(fs): 		WriteRequest: structs.WriteRequest{
// todo(fs): 			Token: "root",
// todo(fs): 		},
// todo(fs): 	}
// todo(fs): 	var token string
// todo(fs): 	if err := a.RPC("ACL.Apply", &arg, &token); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create services using the root token
// todo(fs): 	srv1 := &structs.NodeService{
// todo(fs): 		ID:      "mysql",
// todo(fs): 		Service: "mysql",
// todo(fs): 		Tags:    []string{"master"},
// todo(fs): 		Port:    5000,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv1, "root")
// todo(fs): 	srv2 := &structs.NodeService{
// todo(fs): 		ID:      "api",
// todo(fs): 		Service: "api",
// todo(fs): 		Tags:    []string{"foo"},
// todo(fs): 		Port:    5001,
// todo(fs): 	}
// todo(fs): 	a.state.AddService(srv2, "root")
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs): 	time.Sleep(200 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	{
// todo(fs): 		req := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 			QueryOptions: structs.QueryOptions{
// todo(fs): 				Token: "root",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var services structs.IndexedNodeServices
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 3 services (consul included)
// todo(fs): 		if len(services.NodeServices.Services) != 3 {
// todo(fs): 			t.Fatalf("bad: %v", services.NodeServices.Services)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the services should match
// todo(fs): 		for id, serv := range services.NodeServices.Services {
// todo(fs): 			serv.CreateIndex, serv.ModifyIndex = 0, 0
// todo(fs): 			switch id {
// todo(fs): 			case "mysql":
// todo(fs): 				if !reflect.DeepEqual(serv, srv1) {
// todo(fs): 					t.Fatalf("bad: %#v %#v", serv, srv1)
// todo(fs): 				}
// todo(fs): 			case "api":
// todo(fs): 				if !reflect.DeepEqual(serv, srv2) {
// todo(fs): 					t.Fatalf("bad: %#v %#v", serv, srv2)
// todo(fs): 				}
// todo(fs): 			case structs.ConsulServiceID:
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				t.Fatalf("unexpected service: %v", id)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// todo(fs): data race
// todo(fs): 		func() {
// todo(fs): 			a.state.RLock()
// todo(fs): 			defer a.state.RUnlock()
// todo(fs):
// todo(fs): 			// Check the local state
// todo(fs): 			if len(a.state.services) != 2 {
// todo(fs): 				t.Fatalf("bad: %v", a.state.services)
// todo(fs): 			}
// todo(fs): 			if len(a.state.serviceStatus) != 2 {
// todo(fs): 				t.Fatalf("bad: %v", a.state.serviceStatus)
// todo(fs): 			}
// todo(fs): 			for name, status := range a.state.serviceStatus {
// todo(fs): 				if !status.inSync {
// todo(fs): 					t.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 		}()
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// This check won't be allowed.
// todo(fs): 	chk1 := &structs.HealthCheck{
// todo(fs): 		Node:        a.Config.NodeName,
// todo(fs): 		ServiceID:   "mysql",
// todo(fs): 		ServiceName: "mysql",
// todo(fs): 		ServiceTags: []string{"master"},
// todo(fs): 		CheckID:     "mysql-check",
// todo(fs): 		Name:        "mysql",
// todo(fs): 		Status:      api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk1, token)
// todo(fs):
// todo(fs): 	// This one will be allowed.
// todo(fs): 	chk2 := &structs.HealthCheck{
// todo(fs): 		Node:        a.Config.NodeName,
// todo(fs): 		ServiceID:   "api",
// todo(fs): 		ServiceName: "api",
// todo(fs): 		ServiceTags: []string{"foo"},
// todo(fs): 		CheckID:     "api-check",
// todo(fs): 		Name:        "api",
// todo(fs): 		Status:      api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(chk2, token)
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait.
// todo(fs): 	a.StartSync()
// todo(fs): 	time.Sleep(200 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		req := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 			QueryOptions: structs.QueryOptions{
// todo(fs): 				Token: "root",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var checks structs.IndexedHealthChecks
// todo(fs): 		if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 2 checks (serf included)
// todo(fs): 		if len(checks.HealthChecks) != 2 {
// todo(fs): 			r.Fatalf("bad: %v", checks)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the checks should match
// todo(fs): 		for _, chk := range checks.HealthChecks {
// todo(fs): 			chk.CreateIndex, chk.ModifyIndex = 0, 0
// todo(fs): 			switch chk.CheckID {
// todo(fs): 			case "mysql-check":
// todo(fs): 				t.Fatalf("should not be permitted")
// todo(fs): 			case "api-check":
// todo(fs): 				if !reflect.DeepEqual(chk, chk2) {
// todo(fs): 					r.Fatalf("bad: %v %v", chk, chk2)
// todo(fs): 				}
// todo(fs): 			case "serfHealth":
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				r.Fatalf("unexpected check: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	func() {
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		// Check the local state.
// todo(fs): 		if len(a.state.checks) != 2 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checks)
// todo(fs): 		}
// todo(fs): 		if len(a.state.checkStatus) != 2 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checkStatus)
// todo(fs): 		}
// todo(fs): 		for name, status := range a.state.checkStatus {
// todo(fs): 			if !status.inSync {
// todo(fs): 				t.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs):
// todo(fs): 	// Now delete the check and wait for sync.
// todo(fs): 	a.state.RemoveCheck("api-check")
// todo(fs): 	a.StartSync()
// todo(fs): 	time.Sleep(200 * time.Millisecond)
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		req := structs.NodeSpecificRequest{
// todo(fs): 			Datacenter: "dc1",
// todo(fs): 			Node:       a.Config.NodeName,
// todo(fs): 			QueryOptions: structs.QueryOptions{
// todo(fs): 				Token: "root",
// todo(fs): 			},
// todo(fs): 		}
// todo(fs): 		var checks structs.IndexedHealthChecks
// todo(fs): 		if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// We should have 1 check (just serf)
// todo(fs): 		if len(checks.HealthChecks) != 1 {
// todo(fs): 			r.Fatalf("bad: %v", checks)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// All the checks should match
// todo(fs): 		for _, chk := range checks.HealthChecks {
// todo(fs): 			chk.CreateIndex, chk.ModifyIndex = 0, 0
// todo(fs): 			switch chk.CheckID {
// todo(fs): 			case "mysql-check":
// todo(fs): 				r.Fatalf("should not be permitted")
// todo(fs): 			case "api-check":
// todo(fs): 				r.Fatalf("should be deleted")
// todo(fs): 			case "serfHealth":
// todo(fs): 				// ignore
// todo(fs): 			default:
// todo(fs): 				r.Fatalf("unexpected check: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	func() {
// todo(fs): 		a.state.RLock()
// todo(fs): 		defer a.state.RUnlock()
// todo(fs):
// todo(fs): 		// Check the local state.
// todo(fs): 		if len(a.state.checks) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checks)
// todo(fs): 		}
// todo(fs): 		if len(a.state.checkStatus) != 1 {
// todo(fs): 			t.Fatalf("bad: %v", a.state.checkStatus)
// todo(fs): 		}
// todo(fs): 		for name, status := range a.state.checkStatus {
// todo(fs): 			if !status.inSync {
// todo(fs): 				t.Fatalf("should be in sync: %v %v", name, status)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs):
// todo(fs): 	// Make sure the token got cleaned up.
// todo(fs): 	if token := a.state.CheckToken("api-check"); token != "" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_Check_DeferSync(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.CheckUpdateInterval = 500 * time.Millisecond
// todo(fs): 	a := &TestAgent{Name: t.Name(), Config: cfg, NoInitialSync: true}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create a check
// todo(fs): 	check := &structs.HealthCheck{
// todo(fs): 		Node:    a.Config.NodeName,
// todo(fs): 		CheckID: "web",
// todo(fs): 		Name:    "web",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 		Output:  "",
// todo(fs): 	}
// todo(fs): 	a.state.AddCheck(check, "")
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs):
// todo(fs): 	// Verify that we are in sync
// todo(fs): 	req := structs.NodeSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 	}
// todo(fs): 	var checks structs.IndexedHealthChecks
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if got, want := len(checks.HealthChecks), 2; got != want {
// todo(fs): 			r.Fatalf("got %d health checks want %d", got, want)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Update the check output! Should be deferred
// todo(fs): 	a.state.UpdateCheck("web", api.HealthPassing, "output")
// todo(fs):
// todo(fs): 	// Should not update for 500 milliseconds
// todo(fs): 	time.Sleep(250 * time.Millisecond)
// todo(fs): 	if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify not updated
// todo(fs): 	for _, chk := range checks.HealthChecks {
// todo(fs): 		switch chk.CheckID {
// todo(fs): 		case "web":
// todo(fs): 			if chk.Output != "" {
// todo(fs): 				t.Fatalf("early update: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): 	// Wait for a deferred update
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Verify updated
// todo(fs): 		for _, chk := range checks.HealthChecks {
// todo(fs): 			switch chk.CheckID {
// todo(fs): 			case "web":
// todo(fs): 				if chk.Output != "output" {
// todo(fs): 					r.Fatalf("no update: %v", chk)
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Change the output in the catalog to force it out of sync.
// todo(fs): 	eCopy := check.Clone()
// todo(fs): 	eCopy.Output = "changed"
// todo(fs): 	reg := structs.RegisterRequest{
// todo(fs): 		Datacenter:      a.Config.Datacenter,
// todo(fs): 		Node:            a.Config.NodeName,
// todo(fs): 		Address:         a.Config.AdvertiseAddr,
// todo(fs): 		TaggedAddresses: a.Config.TaggedAddresses,
// todo(fs): 		Check:           eCopy,
// todo(fs): 		WriteRequest:    structs.WriteRequest{},
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", &reg, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify that the output is out of sync.
// todo(fs): 	if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	for _, chk := range checks.HealthChecks {
// todo(fs): 		switch chk.CheckID {
// todo(fs): 		case "web":
// todo(fs): 			if chk.Output != "changed" {
// todo(fs): 				t.Fatalf("unexpected update: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait.
// todo(fs): 	a.StartSync()
// todo(fs): 	time.Sleep(200 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Verify that the output was synced back to the agent's value.
// todo(fs): 	if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	for _, chk := range checks.HealthChecks {
// todo(fs): 		switch chk.CheckID {
// todo(fs): 		case "web":
// todo(fs): 			if chk.Output != "output" {
// todo(fs): 				t.Fatalf("missed update: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Reset the catalog again.
// todo(fs): 	if err := a.RPC("Catalog.Register", &reg, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify that the output is out of sync.
// todo(fs): 	if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	for _, chk := range checks.HealthChecks {
// todo(fs): 		switch chk.CheckID {
// todo(fs): 		case "web":
// todo(fs): 			if chk.Output != "changed" {
// todo(fs): 				t.Fatalf("unexpected update: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Now make an update that should be deferred.
// todo(fs): 	a.state.UpdateCheck("web", api.HealthPassing, "deferred")
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait.
// todo(fs): 	a.StartSync()
// todo(fs): 	time.Sleep(200 * time.Millisecond)
// todo(fs):
// todo(fs): 	// Verify that the output is still out of sync since there's a deferred
// todo(fs): 	// update pending.
// todo(fs): 	if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	for _, chk := range checks.HealthChecks {
// todo(fs): 		switch chk.CheckID {
// todo(fs): 		case "web":
// todo(fs): 			if chk.Output != "changed" {
// todo(fs): 				t.Fatalf("unexpected update: %v", chk)
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): 	// Wait for the deferred update.
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Health.NodeChecks", &req, &checks); err != nil {
// todo(fs): 			r.Fatal(err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Verify updated
// todo(fs): 		for _, chk := range checks.HealthChecks {
// todo(fs): 			switch chk.CheckID {
// todo(fs): 			case "web":
// todo(fs): 				if chk.Output != "deferred" {
// todo(fs): 					r.Fatalf("no update: %v", chk)
// todo(fs): 				}
// todo(fs): 			}
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_NodeInfo(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeID = types.NodeID("40e4a748-2192-161a-0510-9bf59fe950b5")
// todo(fs): 	cfg.Meta["somekey"] = "somevalue"
// todo(fs): 	a := &TestAgent{Name: t.Name(), Config: cfg, NoInitialSync: true}
// todo(fs): 	a.Start()
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Register info
// todo(fs): 	args := &structs.RegisterRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 		Address:    "127.0.0.1",
// todo(fs): 	}
// todo(fs): 	var out struct{}
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs):
// todo(fs): 	req := structs.NodeSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 		Node:       a.Config.NodeName,
// todo(fs): 	}
// todo(fs): 	var services structs.IndexedNodeServices
// todo(fs): 	// Wait for the sync
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Make sure we synced our node info - this should have ridden on the
// todo(fs): 		// "consul" service sync
// todo(fs): 		id := services.NodeServices.Node.ID
// todo(fs): 		addrs := services.NodeServices.Node.TaggedAddresses
// todo(fs): 		meta := services.NodeServices.Node.Meta
// todo(fs): 		delete(meta, structs.MetaSegmentKey) // Added later, not in config.
// todo(fs): 		if id != cfg.NodeID ||
// todo(fs): 			!reflect.DeepEqual(addrs, cfg.TaggedAddresses) ||
// todo(fs): 			!reflect.DeepEqual(meta, cfg.Meta) {
// todo(fs): 			r.Fatalf("bad: %v", services.NodeServices.Node)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	// Blow away the catalog version of the node info
// todo(fs): 	if err := a.RPC("Catalog.Register", args, &out); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Trigger anti-entropy run and wait
// todo(fs): 	a.StartSync()
// todo(fs): 	// Wait for the sync - this should have been a sync of just the node info
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Catalog.NodeServices", &req, &services); err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		id := services.NodeServices.Node.ID
// todo(fs): 		addrs := services.NodeServices.Node.TaggedAddresses
// todo(fs): 		meta := services.NodeServices.Node.Meta
// todo(fs): 		delete(meta, structs.MetaSegmentKey) // Added later, not in config.
// todo(fs): 		if id != cfg.NodeID ||
// todo(fs): 			!reflect.DeepEqual(addrs, cfg.TaggedAddresses) ||
// todo(fs): 			!reflect.DeepEqual(meta, cfg.Meta) {
// todo(fs): 			r.Fatalf("bad: %v", services.NodeServices.Node)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_deleteService_fails(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	l := new(localState)
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	l.Lock()
// todo(fs): 	defer l.Unlock()
// todo(fs): 	if err := l.deleteService(""); err == nil {
// todo(fs): 		t.Fatalf("should have failed")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgentAntiEntropy_deleteCheck_fails(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	l := new(localState)
// todo(fs):
// todo(fs): 	// todo(fs): data race
// todo(fs): 	l.Lock()
// todo(fs): 	defer l.Unlock()
// todo(fs): 	if err := l.deleteCheck(""); err == nil {
// todo(fs): 		t.Fatalf("should have errored")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_serviceTokens(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs):
// todo(fs): 	tokens := new(token.Store)
// todo(fs): 	tokens.UpdateUserToken("default")
// todo(fs): 	l := NewLocalState(TestConfig(), nil, tokens)
// todo(fs):
// todo(fs): 	l.AddService(&structs.NodeService{
// todo(fs): 		ID: "redis",
// todo(fs): 	}, "")
// todo(fs):
// todo(fs): 	// Returns default when no token is set
// todo(fs): 	if token := l.ServiceToken("redis"); token != "default" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Returns configured token
// todo(fs): 	l.serviceTokens["redis"] = "abc123"
// todo(fs): 	if token := l.ServiceToken("redis"); token != "abc123" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Keeps token around for the delete
// todo(fs): 	l.RemoveService("redis")
// todo(fs): 	if token := l.ServiceToken("redis"); token != "abc123" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_checkTokens(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs):
// todo(fs): 	tokens := new(token.Store)
// todo(fs): 	tokens.UpdateUserToken("default")
// todo(fs): 	l := NewLocalState(TestConfig(), nil, tokens)
// todo(fs):
// todo(fs): 	// Returns default when no token is set
// todo(fs): 	if token := l.CheckToken("mem"); token != "default" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Returns configured token
// todo(fs): 	l.checkTokens["mem"] = "abc123"
// todo(fs): 	if token := l.CheckToken("mem"); token != "abc123" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Keeps token around for the delete
// todo(fs): 	l.RemoveCheck("mem")
// todo(fs): 	if token := l.CheckToken("mem"); token != "abc123" {
// todo(fs): 		t.Fatalf("bad: %s", token)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_checkCriticalTime(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	l := NewLocalState(cfg, nil, new(token.Store))
// todo(fs):
// todo(fs): 	svc := &structs.NodeService{ID: "redis", Service: "redis", Port: 8000}
// todo(fs): 	l.AddService(svc, "")
// todo(fs):
// todo(fs): 	// Add a passing check and make sure it's not critical.
// todo(fs): 	checkID := types.CheckID("redis:1")
// todo(fs): 	chk := &structs.HealthCheck{
// todo(fs): 		Node:      "node",
// todo(fs): 		CheckID:   checkID,
// todo(fs): 		Name:      "redis:1",
// todo(fs): 		ServiceID: "redis",
// todo(fs): 		Status:    api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	l.AddCheck(chk, "")
// todo(fs): 	if checks := l.CriticalChecks(); len(checks) > 0 {
// todo(fs): 		t.Fatalf("should not have any critical checks")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Set it to warning and make sure that doesn't show up as critical.
// todo(fs): 	l.UpdateCheck(checkID, api.HealthWarning, "")
// todo(fs): 	if checks := l.CriticalChecks(); len(checks) > 0 {
// todo(fs): 		t.Fatalf("should not have any critical checks")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Fail the check and make sure the time looks reasonable.
// todo(fs): 	l.UpdateCheck(checkID, api.HealthCritical, "")
// todo(fs): 	if crit, ok := l.CriticalChecks()[checkID]; !ok {
// todo(fs): 		t.Fatalf("should have a critical check")
// todo(fs): 	} else if crit.CriticalFor > time.Millisecond {
// todo(fs): 		t.Fatalf("bad: %#v", crit)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Wait a while, then fail it again and make sure the time keeps track
// todo(fs): 	// of the initial failure, and doesn't reset here.
// todo(fs): 	time.Sleep(50 * time.Millisecond)
// todo(fs): 	l.UpdateCheck(chk.CheckID, api.HealthCritical, "")
// todo(fs): 	if crit, ok := l.CriticalChecks()[checkID]; !ok {
// todo(fs): 		t.Fatalf("should have a critical check")
// todo(fs): 	} else if crit.CriticalFor < 25*time.Millisecond ||
// todo(fs): 		crit.CriticalFor > 75*time.Millisecond {
// todo(fs): 		t.Fatalf("bad: %#v", crit)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Set it passing again.
// todo(fs): 	l.UpdateCheck(checkID, api.HealthPassing, "")
// todo(fs): 	if checks := l.CriticalChecks(); len(checks) > 0 {
// todo(fs): 		t.Fatalf("should not have any critical checks")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Fail the check and make sure the time looks like it started again
// todo(fs): 	// from the latest failure, not the original one.
// todo(fs): 	l.UpdateCheck(checkID, api.HealthCritical, "")
// todo(fs): 	if crit, ok := l.CriticalChecks()[checkID]; !ok {
// todo(fs): 		t.Fatalf("should have a critical check")
// todo(fs): 	} else if crit.CriticalFor > time.Millisecond {
// todo(fs): 		t.Fatalf("bad: %#v", crit)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddCheckFailure(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	l := NewLocalState(cfg, nil, new(token.Store))
// todo(fs):
// todo(fs): 	// Add a check for a service that does not exist and verify that it fails
// todo(fs): 	checkID := types.CheckID("redis:1")
// todo(fs): 	chk := &structs.HealthCheck{
// todo(fs): 		Node:      "node",
// todo(fs): 		CheckID:   checkID,
// todo(fs): 		Name:      "redis:1",
// todo(fs): 		ServiceID: "redis",
// todo(fs): 		Status:    api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	expectedErr := "ServiceID \"redis\" does not exist"
// todo(fs): 	if err := l.AddCheck(chk, ""); err == nil || expectedErr != err.Error() {
// todo(fs): 		t.Fatalf("Expected error when adding a check for a non-existent service but got %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_nestedPauseResume(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	l := new(localState)
// todo(fs): 	if l.isPaused() != false {
// todo(fs): 		t.Fatal("localState should be unPaused after init")
// todo(fs): 	}
// todo(fs): 	l.Pause()
// todo(fs): 	if l.isPaused() != true {
// todo(fs): 		t.Fatal("localState should be Paused after first call to Pause()")
// todo(fs): 	}
// todo(fs): 	l.Pause()
// todo(fs): 	if l.isPaused() != true {
// todo(fs): 		t.Fatal("localState should STILL be Paused after second call to Pause()")
// todo(fs): 	}
// todo(fs): 	l.Resume()
// todo(fs): 	if l.isPaused() != true {
// todo(fs): 		t.Fatal("localState should STILL be Paused after FIRST call to Resume()")
// todo(fs): 	}
// todo(fs): 	l.Resume()
// todo(fs): 	if l.isPaused() != false {
// todo(fs): 		t.Fatal("localState should NOT be Paused after SECOND call to Resume()")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	defer func() {
// todo(fs): 		err := recover()
// todo(fs): 		if err == nil {
// todo(fs): 			t.Fatal("unbalanced Resume() should cause a panic()")
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs): 	l.Resume()
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_sendCoordinate(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.SyncCoordinateRateTarget = 10.0 // updates/sec
// todo(fs): 	cfg.SyncCoordinateIntervalMin = 1 * time.Millisecond
// todo(fs): 	cfg.ConsulConfig.CoordinateUpdatePeriod = 100 * time.Millisecond
// todo(fs): 	cfg.ConsulConfig.CoordinateUpdateBatchSize = 10
// todo(fs): 	cfg.ConsulConfig.CoordinateUpdateMaxBatches = 1
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Make sure the coordinate is present.
// todo(fs): 	req := structs.DCSpecificRequest{
// todo(fs): 		Datacenter: a.Config.Datacenter,
// todo(fs): 	}
// todo(fs): 	var reply structs.IndexedCoordinates
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		if err := a.RPC("Coordinate.ListNodes", &req, &reply); err != nil {
// todo(fs): 			r.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 		if len(reply.Coordinates) != 1 {
// todo(fs): 			r.Fatalf("expected a coordinate: %v", reply)
// todo(fs): 		}
// todo(fs): 		coord := reply.Coordinates[0]
// todo(fs): 		if coord.Node != a.Config.NodeName || coord.Coord == nil {
// todo(fs): 			r.Fatalf("bad: %v", coord)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
