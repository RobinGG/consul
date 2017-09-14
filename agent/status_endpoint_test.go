package agent

// todo(fs): func TestStatusLeader(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	obj, err := a.srv.StatusLeader(nil, nil)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs): 	val := obj.(string)
// todo(fs): 	if val == "" {
// todo(fs): 		t.Fatalf("bad addr: %v", obj)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestStatusPeers(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), nil)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	obj, err := a.srv.StatusPeers(nil, nil)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("Err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	peers := obj.([]string)
// todo(fs): 	if len(peers) != 1 {
// todo(fs): 		t.Fatalf("bad peers: %v", peers)
// todo(fs): 	}
// todo(fs): }
