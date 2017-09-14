package agent

// todo(fs): func TestOperator_RaftConfiguration(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/operator/raft/configuration", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.OperatorRaftConfiguration(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 	}
// todo(fs): 	out, ok := obj.(structs.RaftConfigurationResponse)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("unexpected: %T", obj)
// todo(fs): 	}
// todo(fs): 	if len(out.Servers) != 1 ||
// todo(fs): 		!out.Servers[0].Leader ||
// todo(fs): 		!out.Servers[0].Voter {
// todo(fs): 		t.Fatalf("bad: %v", out)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_RaftPeer(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("DELETE", "/v1/operator/raft/peer?address=nope", body)
// todo(fs): 		// If we get this error, it proves we sent the address all the
// todo(fs): 		// way through.
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		_, err := a.srv.OperatorRaftPeer(resp, req)
// todo(fs): 		if err == nil || !strings.Contains(err.Error(),
// todo(fs): 			"address \"nope\" was not found in the Raft configuration") {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs):
// todo(fs): 	t.Run("", func(t *testing.T) {
// todo(fs): 		a := NewTestAgent(t.Name(), nil)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		body := bytes.NewBuffer(nil)
// todo(fs): 		req, _ := http.NewRequest("DELETE", "/v1/operator/raft/peer?id=nope", body)
// todo(fs): 		// If we get this error, it proves we sent the ID all the
// todo(fs): 		// way through.
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		_, err := a.srv.OperatorRaftPeer(resp, req)
// todo(fs): 		if err == nil || !strings.Contains(err.Error(),
// todo(fs): 			"id \"nope\" was not found in the Raft configuration") {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_KeyringInstall(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	oldKey := "H3/9gBxcKKRf45CaI2DlRg=="
// todo(fs): 	newKey := "z90lFx3sZZLtTOkutXcwYg=="
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EncryptKey = oldKey
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBufferString(fmt.Sprintf("{\"Key\":\"%s\"}", newKey))
// todo(fs): 	req, _ := http.NewRequest("POST", "/v1/operator/keyring", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	_, err := a.srv.OperatorKeyringEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	listResponse, err := a.ListKeys("", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if len(listResponse.Responses) != 2 {
// todo(fs): 		t.Fatalf("bad: %d", len(listResponse.Responses))
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, response := range listResponse.Responses {
// todo(fs): 		count, ok := response.Keys[newKey]
// todo(fs): 		if !ok {
// todo(fs): 			t.Fatalf("bad: %v", response.Keys)
// todo(fs): 		}
// todo(fs): 		if count != response.NumNodes {
// todo(fs): 			t.Fatalf("bad: %d, %d", count, response.NumNodes)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_KeyringList(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	key := "H3/9gBxcKKRf45CaI2DlRg=="
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EncryptKey = key
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/operator/keyring", nil)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	r, err := a.srv.OperatorKeyringEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	responses, ok := r.([]*structs.KeyringResponse)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("err: %v", !ok)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check that we get both a LAN and WAN response, and that they both only
// todo(fs): 	// contain the original key
// todo(fs): 	if len(responses) != 2 {
// todo(fs): 		t.Fatalf("bad: %d", len(responses))
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// WAN
// todo(fs): 	if len(responses[0].Keys) != 1 {
// todo(fs): 		t.Fatalf("bad: %d", len(responses[0].Keys))
// todo(fs): 	}
// todo(fs): 	if !responses[0].WAN {
// todo(fs): 		t.Fatalf("bad: %v", responses[0].WAN)
// todo(fs): 	}
// todo(fs): 	if _, ok := responses[0].Keys[key]; !ok {
// todo(fs): 		t.Fatalf("bad: %v", ok)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// LAN
// todo(fs): 	if len(responses[1].Keys) != 1 {
// todo(fs): 		t.Fatalf("bad: %d", len(responses[1].Keys))
// todo(fs): 	}
// todo(fs): 	if responses[1].WAN {
// todo(fs): 		t.Fatalf("bad: %v", responses[1].WAN)
// todo(fs): 	}
// todo(fs): 	if _, ok := responses[1].Keys[key]; !ok {
// todo(fs): 		t.Fatalf("bad: %v", ok)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_KeyringRemove(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	key := "H3/9gBxcKKRf45CaI2DlRg=="
// todo(fs): 	tempKey := "z90lFx3sZZLtTOkutXcwYg=="
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EncryptKey = key
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	_, err := a.InstallKey(tempKey, "", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure the temp key is installed
// todo(fs): 	list, err := a.ListKeys("", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	responses := list.Responses
// todo(fs): 	if len(responses) != 2 {
// todo(fs): 		t.Fatalf("bad: %d", len(responses))
// todo(fs): 	}
// todo(fs): 	for _, response := range responses {
// todo(fs): 		if len(response.Keys) != 2 {
// todo(fs): 			t.Fatalf("bad: %d", len(response.Keys))
// todo(fs): 		}
// todo(fs): 		if _, ok := response.Keys[tempKey]; !ok {
// todo(fs): 			t.Fatalf("bad: %v", ok)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	body := bytes.NewBufferString(fmt.Sprintf("{\"Key\":\"%s\"}", tempKey))
// todo(fs): 	req, _ := http.NewRequest("DELETE", "/v1/operator/keyring", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.OperatorKeyringEndpoint(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure the temp key has been removed
// todo(fs): 	list, err = a.ListKeys("", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	responses = list.Responses
// todo(fs): 	if len(responses) != 2 {
// todo(fs): 		t.Fatalf("bad: %d", len(responses))
// todo(fs): 	}
// todo(fs): 	for _, response := range responses {
// todo(fs): 		if len(response.Keys) != 1 {
// todo(fs): 			t.Fatalf("bad: %d", len(response.Keys))
// todo(fs): 		}
// todo(fs): 		if _, ok := response.Keys[tempKey]; ok {
// todo(fs): 			t.Fatalf("bad: %v", ok)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_KeyringUse(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	oldKey := "H3/9gBxcKKRf45CaI2DlRg=="
// todo(fs): 	newKey := "z90lFx3sZZLtTOkutXcwYg=="
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EncryptKey = oldKey
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	if _, err := a.InstallKey(newKey, "", 0); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	body := bytes.NewBufferString(fmt.Sprintf("{\"Key\":\"%s\"}", newKey))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/operator/keyring", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	_, err := a.srv.OperatorKeyringEndpoint(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if _, err := a.RemoveKey(oldKey, "", 0); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Make sure only the new key remains
// todo(fs): 	list, err := a.ListKeys("", 0)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	responses := list.Responses
// todo(fs): 	if len(responses) != 2 {
// todo(fs): 		t.Fatalf("bad: %d", len(responses))
// todo(fs): 	}
// todo(fs): 	for _, response := range responses {
// todo(fs): 		if len(response.Keys) != 1 {
// todo(fs): 			t.Fatalf("bad: %d", len(response.Keys))
// todo(fs): 		}
// todo(fs): 		if _, ok := response.Keys[newKey]; !ok {
// todo(fs): 			t.Fatalf("bad: %v", ok)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_Keyring_InvalidRelayFactor(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	key := "H3/9gBxcKKRf45CaI2DlRg=="
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EncryptKey = key
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	cases := map[string]string{
// todo(fs): 		"999":  "Relay factor must be in range",
// todo(fs): 		"asdf": "Error parsing relay factor",
// todo(fs): 	}
// todo(fs): 	for relayFactor, errString := range cases {
// todo(fs): 		req, _ := http.NewRequest("GET", "/v1/operator/keyring?relay-factor="+relayFactor, nil)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		_, err := a.srv.OperatorKeyringEndpoint(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		body := resp.Body.String()
// todo(fs): 		if !strings.Contains(body, errString) {
// todo(fs): 			t.Fatalf("bad: %v", body)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_AutopilotGetConfiguration(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/operator/autopilot/configuration", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	obj, err := a.srv.OperatorAutopilotConfiguration(resp, req)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 	}
// todo(fs): 	out, ok := obj.(api.AutopilotConfiguration)
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("unexpected: %T", obj)
// todo(fs): 	}
// todo(fs): 	if !out.CleanupDeadServers {
// todo(fs): 		t.Fatalf("bad: %#v", out)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_AutopilotSetConfiguration(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer([]byte(`{"CleanupDeadServers": false}`))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/operator/autopilot/configuration", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.OperatorAutopilotConfiguration(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args := structs.DCSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var reply structs.AutopilotConfig
// todo(fs): 	if err := a.RPC("Operator.AutopilotGetConfiguration", &args, &reply); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if reply.CleanupDeadServers {
// todo(fs): 		t.Fatalf("bad: %#v", reply)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_AutopilotCASConfiguration(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer([]byte(`{"CleanupDeadServers": false}`))
// todo(fs): 	req, _ := http.NewRequest("PUT", "/v1/operator/autopilot/configuration", body)
// todo(fs): 	resp := httptest.NewRecorder()
// todo(fs): 	if _, err := a.srv.OperatorAutopilotConfiguration(resp, req); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if resp.Code != 200 {
// todo(fs): 		t.Fatalf("bad code: %d", resp.Code)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	args := structs.DCSpecificRequest{
// todo(fs): 		Datacenter: "dc1",
// todo(fs): 	}
// todo(fs):
// todo(fs): 	var reply structs.AutopilotConfig
// todo(fs): 	if err := a.RPC("Operator.AutopilotGetConfiguration", &args, &reply); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	if reply.CleanupDeadServers {
// todo(fs): 		t.Fatalf("bad: %#v", reply)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a CAS request, bad index
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte(`{"CleanupDeadServers": true}`))
// todo(fs): 		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/operator/autopilot/configuration?cas=%d", reply.ModifyIndex-1), buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.OperatorAutopilotConfiguration(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); res {
// todo(fs): 			t.Fatalf("should NOT work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Create a CAS request, good index
// todo(fs): 	{
// todo(fs): 		buf := bytes.NewBuffer([]byte(`{"CleanupDeadServers": true}`))
// todo(fs): 		req, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/operator/autopilot/configuration?cas=%d", reply.ModifyIndex), buf)
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.OperatorAutopilotConfiguration(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if res := obj.(bool); !res {
// todo(fs): 			t.Fatalf("should work")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Verify the update
// todo(fs): 	if err := a.RPC("Operator.AutopilotGetConfiguration", &args, &reply); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if !reply.CleanupDeadServers {
// todo(fs): 		t.Fatalf("bad: %#v", reply)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_ServerHealth(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.RaftProtocol = 3
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/operator/autopilot/health", body)
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.OperatorServerHealth(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 200 {
// todo(fs): 			r.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		out, ok := obj.(*api.OperatorHealthReply)
// todo(fs): 		if !ok {
// todo(fs): 			r.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if len(out.Servers) != 1 ||
// todo(fs): 			!out.Servers[0].Healthy ||
// todo(fs): 			out.Servers[0].Name != a.Config.NodeName ||
// todo(fs): 			out.Servers[0].SerfStatus != "alive" ||
// todo(fs): 			out.FailureTolerance != 0 {
// todo(fs): 			r.Fatalf("bad: %v", out)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
// todo(fs):
// todo(fs): func TestOperator_ServerHealth_Unhealthy(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.RaftProtocol = 3
// todo(fs): 	threshold := time.Duration(-1)
// todo(fs): 	cfg.AutopilotLastContactThreshold = threshold
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	body := bytes.NewBuffer(nil)
// todo(fs): 	req, _ := http.NewRequest("GET", "/v1/operator/autopilot/health", body)
// todo(fs): 	retry.Run(t, func(r *retry.R) {
// todo(fs): 		resp := httptest.NewRecorder()
// todo(fs): 		obj, err := a.srv.OperatorServerHealth(resp, req)
// todo(fs): 		if err != nil {
// todo(fs): 			r.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs): 		if resp.Code != 429 {
// todo(fs): 			r.Fatalf("bad code: %d", resp.Code)
// todo(fs): 		}
// todo(fs): 		out, ok := obj.(*api.OperatorHealthReply)
// todo(fs): 		if !ok {
// todo(fs): 			r.Fatalf("unexpected: %T", obj)
// todo(fs): 		}
// todo(fs): 		if len(out.Servers) != 1 ||
// todo(fs): 			out.Healthy ||
// todo(fs): 			out.Servers[0].Name != a.Config.NodeName {
// todo(fs): 			r.Fatalf("bad: %#v", out.Servers)
// todo(fs): 		}
// todo(fs): 	})
// todo(fs): }
