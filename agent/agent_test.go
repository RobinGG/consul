package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/testutil"
	"github.com/pascaldekloe/goe/verify"
)

func externalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("Unable to lookup network interfaces: %v", err)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("Unable to find a non-loopback interface")
}

func TestAgent_MultiStartStop(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			a := NewTestAgent(t.Name(), "")
			time.Sleep(250 * time.Millisecond)
			a.Shutdown()
		})
	}
}

func TestAgent_StartStop(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	// defer a.Shutdown()

	if err := a.Leave(); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := a.Shutdown(); err != nil {
		t.Fatalf("err: %v", err)
	}

	select {
	case <-a.ShutdownCh():
	default:
		t.Fatalf("should be closed")
	}
}

func TestAgent_RPCPing(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	var out struct{}
	if err := a.RPC("Status.Ping", struct{}{}, &out); err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestAgent_TokenStore(t *testing.T) {
	t.Parallel()

	a := NewTestAgent(t.Name(), `
		acl_token = "user"
		acl_agent_token = "agent"
		acl_agent_master_token = "master"`,
	)
	defer a.Shutdown()

	if got, want := a.tokens.UserToken(), "user"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if got, want := a.tokens.AgentToken(), "agent"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if got, want := a.tokens.IsAgentMasterToken("master"), true; got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}

// todo(fs): func TestAgent_CheckPerformanceSettings(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	// Try a default config.
// todo(fs): 	{
// todo(fs): 		cfg := TestConfig()
// todo(fs): 		cfg.Bootstrap = false
// todo(fs): 		cfg.ConsulConfig = nil
// todo(fs): 		a := NewTestAgent(t.Name(), cfg)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		raftMult := time.Duration(consul.DefaultRaftMultiplier)
// todo(fs): 		r := a.consulConfig().RaftConfig
// todo(fs): 		def := raft.DefaultConfig()
// todo(fs): 		if r.HeartbeatTimeout != raftMult*def.HeartbeatTimeout ||
// todo(fs): 			r.ElectionTimeout != raftMult*def.ElectionTimeout ||
// todo(fs): 			r.LeaderLeaseTimeout != raftMult*def.LeaderLeaseTimeout {
// todo(fs): 			t.Fatalf("bad: %#v", *r)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Try a multiplier.
// todo(fs): 	{
// todo(fs): 		cfg := TestConfig()
// todo(fs): 		cfg.Bootstrap = false
// todo(fs): 		cfg.PerformanceRaftMultiplier = 99
// todo(fs): 		a := NewTestAgent(t.Name(), cfg)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		const raftMult time.Duration = 99
// todo(fs): 		r := a.consulConfig().RaftConfig
// todo(fs): 		def := raft.DefaultConfig()
// todo(fs): 		if r.HeartbeatTimeout != raftMult*def.HeartbeatTimeout ||
// todo(fs): 			r.ElectionTimeout != raftMult*def.ElectionTimeout ||
// todo(fs): 			r.LeaderLeaseTimeout != raftMult*def.LeaderLeaseTimeout {
// todo(fs): 			t.Fatalf("bad: %#v", *r)
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_ReconnectConfigSettings(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	func() {
// todo(fs): 		a := NewTestAgent(t.Name(), "")
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		lan := a.consulConfig().SerfLANConfig.ReconnectTimeout
// todo(fs): 		if lan != 3*24*time.Hour {
// todo(fs): 			t.Fatalf("bad: %s", lan.String())
// todo(fs): 		}
// todo(fs):
// todo(fs): 		wan := a.consulConfig().SerfWANConfig.ReconnectTimeout
// todo(fs): 		if wan != 3*24*time.Hour {
// todo(fs): 			t.Fatalf("bad: %s", wan.String())
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs):
// todo(fs): 	func() {
// todo(fs): 		cfg := TestConfig()
// todo(fs): 		cfg.ReconnectTimeoutLAN = 24 * time.Hour
// todo(fs): 		cfg.ReconnectTimeoutWAN = 36 * time.Hour
// todo(fs): 		a := NewTestAgent(t.Name(), cfg)
// todo(fs): 		defer a.Shutdown()
// todo(fs):
// todo(fs): 		lan := a.consulConfig().SerfLANConfig.ReconnectTimeout
// todo(fs): 		if lan != 24*time.Hour {
// todo(fs): 			t.Fatalf("bad: %s", lan.String())
// todo(fs): 		}
// todo(fs):
// todo(fs): 		wan := a.consulConfig().SerfWANConfig.ReconnectTimeout
// todo(fs): 		if wan != 36*time.Hour {
// todo(fs): 			t.Fatalf("bad: %s", wan.String())
// todo(fs): 		}
// todo(fs): 	}()
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_setupNodeID(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeID = ""
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// The auto-assigned ID should be valid.
// todo(fs): 	id := a.consulConfig().NodeID
// todo(fs): 	if _, err := uuid.ParseUUID(string(id)); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Running again should get the same ID (persisted in the file).
// todo(fs): 	cfg.NodeID = ""
// todo(fs): 	if err := a.setupNodeID(cfg); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if newID := a.consulConfig().NodeID; id != newID {
// todo(fs): 		t.Fatalf("bad: %q vs %q", id, newID)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Set an invalid ID via.Config.
// todo(fs): 	cfg.NodeID = types.NodeID("nope")
// todo(fs): 	err := a.setupNodeID(cfg)
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "uuid string is wrong length") {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Set a valid ID via.Config.
// todo(fs): 	newID, err := uuid.GenerateUUID()
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	cfg.NodeID = types.NodeID(strings.ToUpper(newID))
// todo(fs): 	if err := a.setupNodeID(cfg); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if id := a.consulConfig().NodeID; string(id) != newID {
// todo(fs): 		t.Fatalf("bad: %q vs. %q", id, newID)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Set an invalid ID via the file.
// todo(fs): 	fileID := filepath.Join(cfg.DataDir, "node-id")
// todo(fs): 	if err := ioutil.WriteFile(fileID, []byte("adf4238a!882b!9ddc!4a9d!5b6758e4159e"), 0600); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	cfg.NodeID = ""
// todo(fs): 	err = a.setupNodeID(cfg)
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "uuid is improperly formatted") {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Set a valid ID via the file.
// todo(fs): 	if err := ioutil.WriteFile(fileID, []byte("ADF4238a-882b-9ddc-4a9d-5b6758e4159e"), 0600); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	cfg.NodeID = ""
// todo(fs): 	if err := a.setupNodeID(cfg); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if id := a.consulConfig().NodeID; string(id) != "adf4238a-882b-9ddc-4a9d-5b6758e4159e" {
// todo(fs): 		t.Fatalf("bad: %q vs. %q", id, newID)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_makeNodeID(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeID = ""
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// We should get a valid host-based ID initially.
// todo(fs): 	id, err := a.makeNodeID()
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, err := uuid.ParseUUID(string(id)); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Calling again should yield a random ID by default.
// todo(fs): 	another, err := a.makeNodeID()
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if id == another {
// todo(fs): 		t.Fatalf("bad: %s vs %s", id, another)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Turn on host-based IDs and try again. We should get the same ID
// todo(fs): 	// each time (and a different one from the random one above).
// todo(fs): 	a.Config.DisableHostNodeID = false
// todo(fs): 	id, err = a.makeNodeID()
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if id == another {
// todo(fs): 		t.Fatalf("bad: %s vs %s", id, another)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Calling again should yield the host-based ID.
// todo(fs): 	another, err = a.makeNodeID()
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if id != another {
// todo(fs): 		t.Fatalf("bad: %s vs %s", id, another)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddService(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeName = "node1"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	tests := []struct {
// todo(fs): 		desc       string
// todo(fs): 		srv        *structs.NodeService
// todo(fs): 		chkTypes   []*structs.CheckType
// todo(fs): 		healthChks map[string]*structs.HealthCheck
// todo(fs): 	}{
// todo(fs): 		{
// todo(fs): 			"one check",
// todo(fs): 			&structs.NodeService{
// todo(fs): 				ID:      "svcid1",
// todo(fs): 				Service: "svcname1",
// todo(fs): 				Tags:    []string{"tag1"},
// todo(fs): 				Port:    8100,
// todo(fs): 			},
// todo(fs): 			[]*structs.CheckType{
// todo(fs): 				&structs.CheckType{
// todo(fs): 					CheckID: "check1",
// todo(fs): 					Name:    "name1",
// todo(fs): 					TTL:     time.Minute,
// todo(fs): 					Notes:   "note1",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 			map[string]*structs.HealthCheck{
// todo(fs): 				"check1": &structs.HealthCheck{
// todo(fs): 					Node:        "node1",
// todo(fs): 					CheckID:     "check1",
// todo(fs): 					Name:        "name1",
// todo(fs): 					Status:      "critical",
// todo(fs): 					Notes:       "note1",
// todo(fs): 					ServiceID:   "svcid1",
// todo(fs): 					ServiceName: "svcname1",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 		{
// todo(fs): 			"multiple checks",
// todo(fs): 			&structs.NodeService{
// todo(fs): 				ID:      "svcid2",
// todo(fs): 				Service: "svcname2",
// todo(fs): 				Tags:    []string{"tag2"},
// todo(fs): 				Port:    8200,
// todo(fs): 			},
// todo(fs): 			[]*structs.CheckType{
// todo(fs): 				&structs.CheckType{
// todo(fs): 					CheckID: "check1",
// todo(fs): 					Name:    "name1",
// todo(fs): 					TTL:     time.Minute,
// todo(fs): 					Notes:   "note1",
// todo(fs): 				},
// todo(fs): 				&structs.CheckType{
// todo(fs): 					CheckID: "check-noname",
// todo(fs): 					TTL:     time.Minute,
// todo(fs): 				},
// todo(fs): 				&structs.CheckType{
// todo(fs): 					Name: "check-noid",
// todo(fs): 					TTL:  time.Minute,
// todo(fs): 				},
// todo(fs): 				&structs.CheckType{
// todo(fs): 					TTL: time.Minute,
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 			map[string]*structs.HealthCheck{
// todo(fs): 				"check1": &structs.HealthCheck{
// todo(fs): 					Node:        "node1",
// todo(fs): 					CheckID:     "check1",
// todo(fs): 					Name:        "name1",
// todo(fs): 					Status:      "critical",
// todo(fs): 					Notes:       "note1",
// todo(fs): 					ServiceID:   "svcid2",
// todo(fs): 					ServiceName: "svcname2",
// todo(fs): 				},
// todo(fs): 				"check-noname": &structs.HealthCheck{
// todo(fs): 					Node:        "node1",
// todo(fs): 					CheckID:     "check-noname",
// todo(fs): 					Name:        "Service 'svcname2' check",
// todo(fs): 					Status:      "critical",
// todo(fs): 					ServiceID:   "svcid2",
// todo(fs): 					ServiceName: "svcname2",
// todo(fs): 				},
// todo(fs): 				"service:svcid2:3": &structs.HealthCheck{
// todo(fs): 					Node:        "node1",
// todo(fs): 					CheckID:     "service:svcid2:3",
// todo(fs): 					Name:        "check-noid",
// todo(fs): 					Status:      "critical",
// todo(fs): 					ServiceID:   "svcid2",
// todo(fs): 					ServiceName: "svcname2",
// todo(fs): 				},
// todo(fs): 				"service:svcid2:4": &structs.HealthCheck{
// todo(fs): 					Node:        "node1",
// todo(fs): 					CheckID:     "service:svcid2:4",
// todo(fs): 					Name:        "Service 'svcname2' check",
// todo(fs): 					Status:      "critical",
// todo(fs): 					ServiceID:   "svcid2",
// todo(fs): 					ServiceName: "svcname2",
// todo(fs): 				},
// todo(fs): 			},
// todo(fs): 		},
// todo(fs): 	}
// todo(fs):
// todo(fs): 	for _, tt := range tests {
// todo(fs): 		t.Run(tt.desc, func(t *testing.T) {
// todo(fs): 			// check the service registration
// todo(fs): 			t.Run(tt.srv.ID, func(t *testing.T) {
// todo(fs): 				err := a.AddService(tt.srv, tt.chkTypes, false, "")
// todo(fs): 				if err != nil {
// todo(fs): 					t.Fatalf("err: %v", err)
// todo(fs): 				}
// todo(fs):
// todo(fs): 				got, want := a.state.Services()[tt.srv.ID], tt.srv
// todo(fs): 				verify.Values(t, "", got, want)
// todo(fs): 			})
// todo(fs):
// todo(fs): 			// check the health checks
// todo(fs): 			for k, v := range tt.healthChks {
// todo(fs): 				t.Run(k, func(t *testing.T) {
// todo(fs): 					got, want := a.state.Checks()[types.CheckID(k)], v
// todo(fs): 					verify.Values(t, k, got, want)
// todo(fs): 				})
// todo(fs): 			}
// todo(fs):
// todo(fs): 			// check the ttl checks
// todo(fs): 			for k := range tt.healthChks {
// todo(fs): 				t.Run(k+" ttl", func(t *testing.T) {
// todo(fs): 					chk := a.checkTTLs[types.CheckID(k)]
// todo(fs): 					if chk == nil {
// todo(fs): 						t.Fatal("got nil want TTL check")
// todo(fs): 					}
// todo(fs): 					if got, want := string(chk.CheckID), k; got != want {
// todo(fs): 						t.Fatalf("got CheckID %v want %v", got, want)
// todo(fs): 					}
// todo(fs): 					if got, want := chk.TTL, time.Minute; got != want {
// todo(fs): 						t.Fatalf("got TTL %v want %v", got, want)
// todo(fs): 					}
// todo(fs): 				})
// todo(fs): 			}
// todo(fs): 		})
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RemoveService(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), "")
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Remove a service that doesn't exist
// todo(fs): 	if err := a.RemoveService("redis", false); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Remove without an ID
// todo(fs): 	if err := a.RemoveService("", false); err == nil {
// todo(fs): 		t.Fatalf("should have errored")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Removing a service with a single check works
// todo(fs): 	{
// todo(fs): 		srv := &structs.NodeService{
// todo(fs): 			ID:      "memcache",
// todo(fs): 			Service: "memcache",
// todo(fs): 			Port:    8000,
// todo(fs): 		}
// todo(fs): 		chkTypes := []*structs.CheckType{&structs.CheckType{TTL: time.Minute}}
// todo(fs):
// todo(fs): 		if err := a.AddService(srv, chkTypes, false, ""); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Add a check after the fact with a specific check ID
// todo(fs): 		check := &structs.CheckDefinition{
// todo(fs): 			ID:        "check2",
// todo(fs): 			Name:      "check2",
// todo(fs): 			ServiceID: "memcache",
// todo(fs): 			TTL:       time.Minute,
// todo(fs): 		}
// todo(fs): 		hc := check.HealthCheck("node1")
// todo(fs): 		if err := a.AddCheck(hc, check.CheckType(), false, ""); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		if err := a.RemoveService("memcache", false); err != nil {
// todo(fs): 			t.Fatalf("err: %s", err)
// todo(fs): 		}
// todo(fs): 		if _, ok := a.state.Checks()["service:memcache"]; ok {
// todo(fs): 			t.Fatalf("have memcache check")
// todo(fs): 		}
// todo(fs): 		if _, ok := a.state.Checks()["check2"]; ok {
// todo(fs): 			t.Fatalf("have check2 check")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Removing a service with multiple checks works
// todo(fs): 	{
// todo(fs): 		srv := &structs.NodeService{
// todo(fs): 			ID:      "redis",
// todo(fs): 			Service: "redis",
// todo(fs): 			Port:    8000,
// todo(fs): 		}
// todo(fs): 		chkTypes := []*structs.CheckType{
// todo(fs): 			&structs.CheckType{TTL: time.Minute},
// todo(fs): 			&structs.CheckType{TTL: 30 * time.Second},
// todo(fs): 		}
// todo(fs): 		if err := a.AddService(srv, chkTypes, false, ""); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Remove the service
// todo(fs): 		if err := a.RemoveService("redis", false); err != nil {
// todo(fs): 			t.Fatalf("err: %v", err)
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Ensure we have a state mapping
// todo(fs): 		if _, ok := a.state.Services()["redis"]; ok {
// todo(fs): 			t.Fatalf("have redis service")
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Ensure checks were removed
// todo(fs): 		if _, ok := a.state.Checks()["service:redis:1"]; ok {
// todo(fs): 			t.Fatalf("check redis:1 should be removed")
// todo(fs): 		}
// todo(fs): 		if _, ok := a.state.Checks()["service:redis:2"]; ok {
// todo(fs): 			t.Fatalf("check redis:2 should be removed")
// todo(fs): 		}
// todo(fs):
// todo(fs): 		// Ensure a TTL is setup
// todo(fs): 		if _, ok := a.checkTTLs["service:redis:1"]; ok {
// todo(fs): 			t.Fatalf("check ttl for redis:1 should be removed")
// todo(fs): 		}
// todo(fs): 		if _, ok := a.checkTTLs["service:redis:2"]; ok {
// todo(fs): 			t.Fatalf("check ttl for redis:2 should be removed")
// todo(fs): 		}
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RemoveServiceRemovesAllChecks(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.NodeName = "node1"
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	svc := &structs.NodeService{ID: "redis", Service: "redis", Port: 8000}
// todo(fs): 	chk1 := &structs.CheckType{CheckID: "chk1", Name: "chk1", TTL: time.Minute}
// todo(fs): 	chk2 := &structs.CheckType{CheckID: "chk2", Name: "chk2", TTL: 2 * time.Minute}
// todo(fs): 	hchk1 := &structs.HealthCheck{Node: "node1", CheckID: "chk1", Name: "chk1", Status: "critical", ServiceID: "redis", ServiceName: "redis"}
// todo(fs): 	hchk2 := &structs.HealthCheck{Node: "node1", CheckID: "chk2", Name: "chk2", Status: "critical", ServiceID: "redis", ServiceName: "redis"}
// todo(fs):
// todo(fs): 	// register service with chk1
// todo(fs): 	if err := a.AddService(svc, []*structs.CheckType{chk1}, false, ""); err != nil {
// todo(fs): 		t.Fatal("Failed to register service", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// verify chk1 exists
// todo(fs): 	if a.state.Checks()["chk1"] == nil {
// todo(fs): 		t.Fatal("Could not find health check chk1")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// update the service with chk2
// todo(fs): 	if err := a.AddService(svc, []*structs.CheckType{chk2}, false, ""); err != nil {
// todo(fs): 		t.Fatal("Failed to update service", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// check that both checks are there
// todo(fs): 	if got, want := a.state.Checks()["chk1"], hchk1; !verify.Values(t, "", got, want) {
// todo(fs): 		t.FailNow()
// todo(fs): 	}
// todo(fs): 	if got, want := a.state.Checks()["chk2"], hchk2; !verify.Values(t, "", got, want) {
// todo(fs): 		t.FailNow()
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Remove service
// todo(fs): 	if err := a.RemoveService("redis", false); err != nil {
// todo(fs): 		t.Fatal("Failed to remove service", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Check that both checks are gone
// todo(fs): 	if a.state.Checks()["chk1"] != nil {
// todo(fs): 		t.Fatal("Found health check chk1 want nil")
// todo(fs): 	}
// todo(fs): 	if a.state.Checks()["chk2"] != nil {
// todo(fs): 		t.Fatal("Found health check chk2 want nil")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EnableScriptChecks = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:    "foo",
// todo(fs): 		CheckID: "mem",
// todo(fs): 		Name:    "memory util",
// todo(fs): 		Status:  api.HealthCritical,
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		Script:   "exit 0",
// todo(fs): 		Interval: 15 * time.Second,
// todo(fs): 	}
// todo(fs): 	err := a.AddCheck(health, chk, false, "")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	sChk, ok := a.state.Checks()["mem"]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("missing mem check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure our check is in the right state
// todo(fs): 	if sChk.Status != api.HealthCritical {
// todo(fs): 		t.Fatalf("check not critical")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure a TTL is setup
// todo(fs): 	if _, ok := a.checkMonitors["mem"]; !ok {
// todo(fs): 		t.Fatalf("missing mem monitor")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddCheck_StartPassing(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EnableScriptChecks = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:    "foo",
// todo(fs): 		CheckID: "mem",
// todo(fs): 		Name:    "memory util",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		Script:   "exit 0",
// todo(fs): 		Interval: 15 * time.Second,
// todo(fs): 	}
// todo(fs): 	err := a.AddCheck(health, chk, false, "")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	sChk, ok := a.state.Checks()["mem"]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("missing mem check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure our check is in the right state
// todo(fs): 	if sChk.Status != api.HealthPassing {
// todo(fs): 		t.Fatalf("check not passing")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure a TTL is setup
// todo(fs): 	if _, ok := a.checkMonitors["mem"]; !ok {
// todo(fs): 		t.Fatalf("missing mem monitor")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddCheck_MinInterval(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EnableScriptChecks = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:    "foo",
// todo(fs): 		CheckID: "mem",
// todo(fs): 		Name:    "memory util",
// todo(fs): 		Status:  api.HealthCritical,
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		Script:   "exit 0",
// todo(fs): 		Interval: time.Microsecond,
// todo(fs): 	}
// todo(fs): 	err := a.AddCheck(health, chk, false, "")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	if _, ok := a.state.Checks()["mem"]; !ok {
// todo(fs): 		t.Fatalf("missing mem check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure a TTL is setup
// todo(fs): 	if mon, ok := a.checkMonitors["mem"]; !ok {
// todo(fs): 		t.Fatalf("missing mem monitor")
// todo(fs): 	} else if mon.Interval != MinInterval {
// todo(fs): 		t.Fatalf("bad mem monitor interval")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddCheck_MissingService(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EnableScriptChecks = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:      "foo",
// todo(fs): 		CheckID:   "baz",
// todo(fs): 		Name:      "baz check 1",
// todo(fs): 		ServiceID: "baz",
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		Script:   "exit 0",
// todo(fs): 		Interval: time.Microsecond,
// todo(fs): 	}
// todo(fs): 	err := a.AddCheck(health, chk, false, "")
// todo(fs): 	if err == nil || err.Error() != `ServiceID "baz" does not exist` {
// todo(fs): 		t.Fatalf("expected service id error, got: %v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddCheck_RestoreState(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), "")
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Create some state and persist it
// todo(fs): 	ttl := &CheckTTL{
// todo(fs): 		CheckID: "baz",
// todo(fs): 		TTL:     time.Minute,
// todo(fs): 	}
// todo(fs): 	err := a.persistCheckState(ttl, api.HealthPassing, "yup")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Build and register the check definition and initial state
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:    "foo",
// todo(fs): 		CheckID: "baz",
// todo(fs): 		Name:    "baz check 1",
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		TTL: time.Minute,
// todo(fs): 	}
// todo(fs): 	err = a.AddCheck(health, chk, false, "")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the check status was restored during registration
// todo(fs): 	checks := a.state.Checks()
// todo(fs): 	check, ok := checks["baz"]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("missing check")
// todo(fs): 	}
// todo(fs): 	if check.Status != api.HealthPassing {
// todo(fs): 		t.Fatalf("bad: %#v", check)
// todo(fs): 	}
// todo(fs): 	if check.Output != "yup" {
// todo(fs): 		t.Fatalf("bad: %#v", check)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_AddCheck_ExecDisable(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs):
// todo(fs): 	a := NewTestAgent(t.Name(), "")
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:    "foo",
// todo(fs): 		CheckID: "mem",
// todo(fs): 		Name:    "memory util",
// todo(fs): 		Status:  api.HealthCritical,
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		Script:   "exit 0",
// todo(fs): 		Interval: 15 * time.Second,
// todo(fs): 	}
// todo(fs): 	err := a.AddCheck(health, chk, false, "")
// todo(fs): 	if err == nil || !strings.Contains(err.Error(), "Scripts are disabled on this agent") {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we don't have a check mapping
// todo(fs): 	if memChk := a.state.Checks()["mem"]; memChk != nil {
// todo(fs): 		t.Fatalf("should be missing mem check")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_RemoveCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.EnableScriptChecks = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	// Remove check that doesn't exist
// todo(fs): 	if err := a.RemoveCheck("mem", false); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Remove without an ID
// todo(fs): 	if err := a.RemoveCheck("", false); err == nil {
// todo(fs): 		t.Fatalf("should have errored")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:    "foo",
// todo(fs): 		CheckID: "mem",
// todo(fs): 		Name:    "memory util",
// todo(fs): 		Status:  api.HealthCritical,
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		Script:   "exit 0",
// todo(fs): 		Interval: 15 * time.Second,
// todo(fs): 	}
// todo(fs): 	err := a.AddCheck(health, chk, false, "")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Remove check
// todo(fs): 	if err := a.RemoveCheck("mem", false); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping
// todo(fs): 	if _, ok := a.state.Checks()["mem"]; ok {
// todo(fs): 		t.Fatalf("have mem check")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure a TTL is setup
// todo(fs): 	if _, ok := a.checkMonitors["mem"]; ok {
// todo(fs): 		t.Fatalf("have mem monitor")
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_updateTTLCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), "")
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	health := &structs.HealthCheck{
// todo(fs): 		Node:    "foo",
// todo(fs): 		CheckID: "mem",
// todo(fs): 		Name:    "memory util",
// todo(fs): 		Status:  api.HealthCritical,
// todo(fs): 	}
// todo(fs): 	chk := &structs.CheckType{
// todo(fs): 		TTL: 15 * time.Second,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Add check and update it.
// todo(fs): 	err := a.AddCheck(health, chk, false, "")
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if err := a.updateTTLCheck("mem", api.HealthPassing, "foo"); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure we have a check mapping.
// todo(fs): 	status := a.state.Checks()["mem"]
// todo(fs): 	if status.Status != api.HealthPassing {
// todo(fs): 		t.Fatalf("bad: %v", status)
// todo(fs): 	}
// todo(fs): 	if status.Output != "foo" {
// todo(fs): 		t.Fatalf("bad: %v", status)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_PersistService(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ServerMode = false
// todo(fs): 	cfg.DataDir = testutil.TempDir(t, "agent") // we manage the data dir
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer os.RemoveAll(cfg.DataDir)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	svc := &structs.NodeService{
// todo(fs): 		ID:      "redis",
// todo(fs): 		Service: "redis",
// todo(fs): 		Tags:    []string{"foo"},
// todo(fs): 		Port:    8000,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	file := filepath.Join(a.Config.DataDir, servicesDir, stringHash(svc.ID))
// todo(fs):
// todo(fs): 	// Check is not persisted unless requested
// todo(fs): 	if err := a.AddService(svc, nil, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, err := os.Stat(file); err == nil {
// todo(fs): 		t.Fatalf("should not persist")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Persists to file if requested
// todo(fs): 	if err := a.AddService(svc, nil, true, "mytoken"); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, err := os.Stat(file); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	expected, err := json.Marshal(persistedService{
// todo(fs): 		Token:   "mytoken",
// todo(fs): 		Service: svc,
// todo(fs): 	})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	content, err := ioutil.ReadFile(file)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if !bytes.Equal(expected, content) {
// todo(fs): 		t.Fatalf("bad: %s", string(content))
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Updates service definition on disk
// todo(fs): 	svc.Port = 8001
// todo(fs): 	if err := a.AddService(svc, nil, true, "mytoken"); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	expected, err = json.Marshal(persistedService{
// todo(fs): 		Token:   "mytoken",
// todo(fs): 		Service: svc,
// todo(fs): 	})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	content, err = ioutil.ReadFile(file)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if !bytes.Equal(expected, content) {
// todo(fs): 		t.Fatalf("bad: %s", string(content))
// todo(fs): 	}
// todo(fs): 	a.Shutdown()
// todo(fs):
// todo(fs): 	// Should load it back during later start
// todo(fs): 	a2 := NewTestAgent(t.Name()+"-a2", cfg)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	restored, ok := a2.state.services[svc.ID]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("bad: %#v", a2.state.services)
// todo(fs): 	}
// todo(fs): 	if a2.state.serviceTokens[svc.ID] != "mytoken" {
// todo(fs): 		t.Fatalf("bad: %#v", a2.state.services[svc.ID])
// todo(fs): 	}
// todo(fs): 	if restored.Port != 8001 {
// todo(fs): 		t.Fatalf("bad: %#v", restored)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_persistedService_compat(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	// Tests backwards compatibility of persisted services from pre-0.5.1
// todo(fs): 	a := NewTestAgent(t.Name(), "")
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	svc := &structs.NodeService{
// todo(fs): 		ID:      "redis",
// todo(fs): 		Service: "redis",
// todo(fs): 		Tags:    []string{"foo"},
// todo(fs): 		Port:    8000,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Encode the NodeService directly. This is what previous versions
// todo(fs): 	// would serialize to the file (without the wrapper)
// todo(fs): 	encoded, err := json.Marshal(svc)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Write the content to the file
// todo(fs): 	file := filepath.Join(a.Config.DataDir, servicesDir, stringHash(svc.ID))
// todo(fs): 	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if err := ioutil.WriteFile(file, encoded, 0600); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Load the services
// todo(fs): 	if err := a.loadServices(a.Config); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Ensure the service was restored
// todo(fs): 	services := a.state.Services()
// todo(fs): 	result, ok := services["redis"]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("missing service")
// todo(fs): 	}
// todo(fs): 	if !reflect.DeepEqual(result, svc) {
// todo(fs): 		t.Fatalf("bad: %#v", result)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_PurgeService(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	a := NewTestAgent(t.Name(), "")
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	svc := &structs.NodeService{
// todo(fs): 		ID:      "redis",
// todo(fs): 		Service: "redis",
// todo(fs): 		Tags:    []string{"foo"},
// todo(fs): 		Port:    8000,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	file := filepath.Join(a.Config.DataDir, servicesDir, stringHash(svc.ID))
// todo(fs): 	if err := a.AddService(svc, nil, true, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Not removed
// todo(fs): 	if err := a.RemoveService(svc.ID, false); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if _, err := os.Stat(file); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Re-add the service
// todo(fs): 	if err := a.AddService(svc, nil, true, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Removed
// todo(fs): 	if err := a.RemoveService(svc.ID, true); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if _, err := os.Stat(file); !os.IsNotExist(err) {
// todo(fs): 		t.Fatalf("bad: %#v", err)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_PurgeServiceOnDuplicate(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ServerMode = false
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	svc1 := &structs.NodeService{
// todo(fs): 		ID:      "redis",
// todo(fs): 		Service: "redis",
// todo(fs): 		Tags:    []string{"foo"},
// todo(fs): 		Port:    8000,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// First persist the service
// todo(fs): 	if err := a.AddService(svc1, nil, true, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	a.Shutdown()
// todo(fs):
// todo(fs): 	// Try bringing the agent back up with the service already
// todo(fs): 	// existing in the config
// todo(fs): 	svc2 := &structs.ServiceDefinition{
// todo(fs): 		ID:   "redis",
// todo(fs): 		Name: "redis",
// todo(fs): 		Tags: []string{"bar"},
// todo(fs): 		Port: 9000,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	cfg.Services = []*structs.ServiceDefinition{svc2}
// todo(fs): 	a2 := NewTestAgent(t.Name()+"-a2", cfg)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	file := filepath.Join(a.Config.DataDir, servicesDir, stringHash(svc1.ID))
// todo(fs): 	if _, err := os.Stat(file); err == nil {
// todo(fs): 		t.Fatalf("should have removed persisted service")
// todo(fs): 	}
// todo(fs): 	result, ok := a2.state.services[svc2.ID]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("missing service registration")
// todo(fs): 	}
// todo(fs): 	if !reflect.DeepEqual(result.Tags, svc2.Tags) || result.Port != svc2.Port {
// todo(fs): 		t.Fatalf("bad: %#v", result)
// todo(fs): 	}
// todo(fs): }
// todo(fs):
// todo(fs): func TestAgent_PersistCheck(t *testing.T) {
// todo(fs): 	t.Parallel()
// todo(fs): 	cfg := TestConfig()
// todo(fs): 	cfg.ServerMode = false
// todo(fs): 	cfg.DataDir = testutil.TempDir(t, "agent") // we manage the data dir
// todo(fs): 	cfg.EnableScriptChecks = true
// todo(fs): 	a := NewTestAgent(t.Name(), cfg)
// todo(fs): 	defer os.RemoveAll(cfg.DataDir)
// todo(fs): 	defer a.Shutdown()
// todo(fs):
// todo(fs): 	check := &structs.HealthCheck{
// todo(fs): 		Node:    cfg.NodeName,
// todo(fs): 		CheckID: "mem",
// todo(fs): 		Name:    "memory check",
// todo(fs): 		Status:  api.HealthPassing,
// todo(fs): 	}
// todo(fs): 	chkType := &structs.CheckType{
// todo(fs): 		Script:   "/bin/true",
// todo(fs): 		Interval: 10 * time.Second,
// todo(fs): 	}
// todo(fs):
// todo(fs): 	file := filepath.Join(a.Config.DataDir, checksDir, checkIDHash(check.CheckID))
// todo(fs):
// todo(fs): 	// Not persisted if not requested
// todo(fs): 	if err := a.AddCheck(check, chkType, false, ""); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, err := os.Stat(file); err == nil {
// todo(fs): 		t.Fatalf("should not persist")
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Should persist if requested
// todo(fs): 	if err := a.AddCheck(check, chkType, true, "mytoken"); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	if _, err := os.Stat(file); err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	expected, err := json.Marshal(persistedCheck{
// todo(fs): 		Check:   check,
// todo(fs): 		ChkType: chkType,
// todo(fs): 		Token:   "mytoken",
// todo(fs): 	})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	content, err := ioutil.ReadFile(file)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if !bytes.Equal(expected, content) {
// todo(fs): 		t.Fatalf("bad: %s", string(content))
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Updates the check definition on disk
// todo(fs): 	check.Name = "mem1"
// todo(fs): 	if err := a.AddCheck(check, chkType, true, "mytoken"); err != nil {
// todo(fs): 		t.Fatalf("err: %v", err)
// todo(fs): 	}
// todo(fs): 	expected, err = json.Marshal(persistedCheck{
// todo(fs): 		Check:   check,
// todo(fs): 		ChkType: chkType,
// todo(fs): 		Token:   "mytoken",
// todo(fs): 	})
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	content, err = ioutil.ReadFile(file)
// todo(fs): 	if err != nil {
// todo(fs): 		t.Fatalf("err: %s", err)
// todo(fs): 	}
// todo(fs): 	if !bytes.Equal(expected, content) {
// todo(fs): 		t.Fatalf("bad: %s", string(content))
// todo(fs): 	}
// todo(fs): 	a.Shutdown()
// todo(fs):
// todo(fs): 	// Should load it back during later start
// todo(fs): 	a2 := NewTestAgent(t.Name()+"-a2", cfg)
// todo(fs): 	defer a2.Shutdown()
// todo(fs):
// todo(fs): 	result, ok := a2.state.checks[check.CheckID]
// todo(fs): 	if !ok {
// todo(fs): 		t.Fatalf("bad: %#v", a2.state.checks)
// todo(fs): 	}
// todo(fs): 	if result.Status != api.HealthCritical {
// todo(fs): 		t.Fatalf("bad: %#v", result)
// todo(fs): 	}
// todo(fs): 	if result.Name != "mem1" {
// todo(fs): 		t.Fatalf("bad: %#v", result)
// todo(fs): 	}
// todo(fs):
// todo(fs): 	// Should have restored the monitor
// todo(fs): 	if _, ok := a2.checkMonitors[check.CheckID]; !ok {
// todo(fs): 		t.Fatalf("bad: %#v", a2.checkMonitors)
// todo(fs): 	}
// todo(fs): 	if a2.state.checkTokens[check.CheckID] != "mytoken" {
// todo(fs): 		t.Fatalf("bad: %s", a2.state.checkTokens[check.CheckID])
// todo(fs): 	}
// todo(fs): }
// todo(fs):
func TestAgent_PurgeCheck(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	check := &structs.HealthCheck{
		Node:    a.Config.NodeName,
		CheckID: "mem",
		Name:    "memory check",
		Status:  api.HealthPassing,
	}

	file := filepath.Join(a.Config.DataDir, checksDir, checkIDHash(check.CheckID))
	if err := a.AddCheck(check, nil, true, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Not removed
	if err := a.RemoveCheck(check.CheckID, false); err != nil {
		t.Fatalf("err: %s", err)
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Removed
	if err := a.RemoveCheck(check.CheckID, true); err != nil {
		t.Fatalf("err: %s", err)
	}
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Fatalf("bad: %#v", err)
	}
}

func TestAgent_PurgeCheckOnDuplicate(t *testing.T) {
	t.Parallel()
	nodeID := NodeID()
	dataDir := testutil.TempDir(t, "agent")
	a := NewTestAgent(t.Name(), `
	    node_id = "`+nodeID+`"
	    node_name = "Node `+nodeID+`"
		data_dir = "`+dataDir+`"
		server = false
		enable_script_checks = true
	`)
	defer os.RemoveAll(dataDir)
	defer a.Shutdown()

	check1 := &structs.HealthCheck{
		Node:    a.Config.NodeName,
		CheckID: "mem",
		Name:    "memory check",
		Status:  api.HealthPassing,
	}

	// First persist the check
	if err := a.AddCheck(check1, nil, true, ""); err != nil {
		t.Fatalf("err: %v", err)
	}
	a.Shutdown()

	// Start again with the check registered in config
	a2 := NewTestAgent(t.Name()+"-a2", `
	    node_id = "`+nodeID+`"
	    node_name = "Node `+nodeID+`"
		data_dir = "`+dataDir+`"
		server = false
		enable_script_checks = true
		check = {
			id = "mem"
			name = "memory check"
			notes = "my cool notes"
			script = "/bin/check-redis.py"
			interval = "30s"
		}
	`)
	defer a2.Shutdown()

	file := filepath.Join(dataDir, checksDir, checkIDHash(check1.CheckID))
	if _, err := os.Stat(file); err == nil {
		t.Fatalf("should have removed persisted check")
	}
	result, ok := a2.state.checks["mem"]
	if !ok {
		t.Fatalf("missing check registration")
	}
	expected := &structs.HealthCheck{
		Node:    a2.Config.NodeName,
		CheckID: "mem",
		Name:    "memory check",
		Status:  api.HealthCritical,
		Notes:   "my cool notes",
	}
	if got, want := result, expected; !verify.Values(t, "", got, want) {
		t.FailNow()
	}
}

func TestAgent_loadChecks_token(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), `
		check = {
			id = "rabbitmq"
			name = "rabbitmq"
			token = "abc123"
			ttl = "10s"
		}
	`)
	defer a.Shutdown()

	checks := a.state.Checks()
	if _, ok := checks["rabbitmq"]; !ok {
		t.Fatalf("missing check")
	}
	if token := a.state.CheckToken("rabbitmq"); token != "abc123" {
		t.Fatalf("bad: %s", token)
	}
}

func TestAgent_unloadChecks(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// First register a service
	svc := &structs.NodeService{
		ID:      "redis",
		Service: "redis",
		Tags:    []string{"foo"},
		Port:    8000,
	}
	if err := a.AddService(svc, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a check
	check1 := &structs.HealthCheck{
		Node:        a.Config.NodeName,
		CheckID:     "service:redis",
		Name:        "redischeck",
		Status:      api.HealthPassing,
		ServiceID:   "redis",
		ServiceName: "redis",
	}
	if err := a.AddCheck(check1, nil, false, ""); err != nil {
		t.Fatalf("err: %s", err)
	}
	found := false
	for check := range a.state.Checks() {
		if check == check1.CheckID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("check should have been registered")
	}

	// Unload all of the checks
	if err := a.unloadChecks(); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Make sure it was unloaded
	for check := range a.state.Checks() {
		if check == check1.CheckID {
			t.Fatalf("should have unloaded checks")
		}
	}
}

func TestAgent_loadServices_token(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), `
		service = {
			id = "rabbitmq"
			name = "rabbitmq"
			port = 5672
			token = "abc123"
		}
	`)
	defer a.Shutdown()

	services := a.state.Services()
	if _, ok := services["rabbitmq"]; !ok {
		t.Fatalf("missing service")
	}
	if token := a.state.ServiceToken("rabbitmq"); token != "abc123" {
		t.Fatalf("bad: %s", token)
	}
}

func TestAgent_unloadServices(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	svc := &structs.NodeService{
		ID:      "redis",
		Service: "redis",
		Tags:    []string{"foo"},
		Port:    8000,
	}

	// Register the service
	if err := a.AddService(svc, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}
	found := false
	for id := range a.state.Services() {
		if id == svc.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("should have registered service")
	}

	// Unload all services
	if err := a.unloadServices(); err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(a.state.Services()) != 0 {
		t.Fatalf("should have unloaded services")
	}
}

func TestAgent_Service_MaintenanceMode(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	svc := &structs.NodeService{
		ID:      "redis",
		Service: "redis",
		Tags:    []string{"foo"},
		Port:    8000,
	}

	// Register the service
	if err := a.AddService(svc, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Enter maintenance mode for the service
	if err := a.EnableServiceMaintenance("redis", "broken", "mytoken"); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Make sure the critical health check was added
	checkID := serviceMaintCheckID("redis")
	check, ok := a.state.Checks()[checkID]
	if !ok {
		t.Fatalf("should have registered critical maintenance check")
	}

	// Check that the token was used to register the check
	if token := a.state.CheckToken(checkID); token != "mytoken" {
		t.Fatalf("expected 'mytoken', got: '%s'", token)
	}

	// Ensure the reason was set in notes
	if check.Notes != "broken" {
		t.Fatalf("bad: %#v", check)
	}

	// Leave maintenance mode
	if err := a.DisableServiceMaintenance("redis"); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure the check was deregistered
	if _, ok := a.state.Checks()[checkID]; ok {
		t.Fatalf("should have deregistered maintenance check")
	}

	// Enter service maintenance mode without providing a reason
	if err := a.EnableServiceMaintenance("redis", "", ""); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure the check was registered with the default notes
	check, ok = a.state.Checks()[checkID]
	if !ok {
		t.Fatalf("should have registered critical check")
	}
	if check.Notes != defaultServiceMaintReason {
		t.Fatalf("bad: %#v", check)
	}
}

func TestAgent_Service_Reap(t *testing.T) {
	// t.Parallel() // timing test. no parallel
	a := NewTestAgent(t.Name(), `
		check_reap_interval = "50ms"
		check_deregister_interval_min = "0s"
	`)
	defer a.Shutdown()

	svc := &structs.NodeService{
		ID:      "redis",
		Service: "redis",
		Tags:    []string{"foo"},
		Port:    8000,
	}
	chkTypes := []*structs.CheckType{
		&structs.CheckType{
			Status: api.HealthPassing,
			TTL:    25 * time.Millisecond,
			DeregisterCriticalServiceAfter: 200 * time.Millisecond,
		},
	}

	// Register the service.
	if err := a.AddService(svc, chkTypes, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Make sure it's there and there's no critical check yet.
	if _, ok := a.state.Services()["redis"]; !ok {
		t.Fatalf("should have redis service")
	}
	if checks := a.state.CriticalChecks(); len(checks) > 0 {
		t.Fatalf("should not have critical checks")
	}

	// Wait for the check TTL to fail but before the check is reaped.
	time.Sleep(100 * time.Millisecond)
	if _, ok := a.state.Services()["redis"]; !ok {
		t.Fatalf("should have redis service")
	}
	if checks := a.state.CriticalChecks(); len(checks) != 1 {
		t.Fatalf("should have a critical check")
	}

	// Pass the TTL.
	if err := a.updateTTLCheck("service:redis", api.HealthPassing, "foo"); err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, ok := a.state.Services()["redis"]; !ok {
		t.Fatalf("should have redis service")
	}
	if checks := a.state.CriticalChecks(); len(checks) > 0 {
		t.Fatalf("should not have critical checks")
	}

	// Wait for the check TTL to fail again.
	time.Sleep(100 * time.Millisecond)
	if _, ok := a.state.Services()["redis"]; !ok {
		t.Fatalf("should have redis service")
	}
	if checks := a.state.CriticalChecks(); len(checks) != 1 {
		t.Fatalf("should have a critical check")
	}

	// Wait for the reap.
	time.Sleep(400 * time.Millisecond)
	if _, ok := a.state.Services()["redis"]; ok {
		t.Fatalf("redis service should have been reaped")
	}
	if checks := a.state.CriticalChecks(); len(checks) > 0 {
		t.Fatalf("should not have critical checks")
	}
}

func TestAgent_Service_NoReap(t *testing.T) {
	// t.Parallel() // timing test. no parallel
	a := NewTestAgent(t.Name(), `
		check_reap_interval = "50ms"
		check_deregister_interval_min = "0s"
	`)
	defer a.Shutdown()

	svc := &structs.NodeService{
		ID:      "redis",
		Service: "redis",
		Tags:    []string{"foo"},
		Port:    8000,
	}
	chkTypes := []*structs.CheckType{
		&structs.CheckType{
			Status: api.HealthPassing,
			TTL:    25 * time.Millisecond,
		},
	}

	// Register the service.
	if err := a.AddService(svc, chkTypes, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Make sure it's there and there's no critical check yet.
	if _, ok := a.state.Services()["redis"]; !ok {
		t.Fatalf("should have redis service")
	}
	if checks := a.state.CriticalChecks(); len(checks) > 0 {
		t.Fatalf("should not have critical checks")
	}

	// Wait for the check TTL to fail.
	time.Sleep(200 * time.Millisecond)
	if _, ok := a.state.Services()["redis"]; !ok {
		t.Fatalf("should have redis service")
	}
	if checks := a.state.CriticalChecks(); len(checks) != 1 {
		t.Fatalf("should have a critical check")
	}

	// Wait a while and make sure it doesn't reap.
	time.Sleep(200 * time.Millisecond)
	if _, ok := a.state.Services()["redis"]; !ok {
		t.Fatalf("should have redis service")
	}
	if checks := a.state.CriticalChecks(); len(checks) != 1 {
		t.Fatalf("should have a critical check")
	}
}

func TestAgent_addCheck_restoresSnapshot(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// First register a service
	svc := &structs.NodeService{
		ID:      "redis",
		Service: "redis",
		Tags:    []string{"foo"},
		Port:    8000,
	}
	if err := a.AddService(svc, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a check
	check1 := &structs.HealthCheck{
		Node:        a.Config.NodeName,
		CheckID:     "service:redis",
		Name:        "redischeck",
		Status:      api.HealthPassing,
		ServiceID:   "redis",
		ServiceName: "redis",
	}
	if err := a.AddCheck(check1, nil, false, ""); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Re-registering the service preserves the state of the check
	chkTypes := []*structs.CheckType{&structs.CheckType{TTL: 30 * time.Second}}
	if err := a.AddService(svc, chkTypes, false, ""); err != nil {
		t.Fatalf("err: %s", err)
	}
	check, ok := a.state.Checks()["service:redis"]
	if !ok {
		t.Fatalf("missing check")
	}
	if check.Status != api.HealthPassing {
		t.Fatalf("bad: %s", check.Status)
	}
}

func TestAgent_NodeMaintenanceMode(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// Enter maintenance mode for the node
	a.EnableNodeMaintenance("broken", "mytoken")

	// Make sure the critical health check was added
	check, ok := a.state.Checks()[structs.NodeMaint]
	if !ok {
		t.Fatalf("should have registered critical node check")
	}

	// Check that the token was used to register the check
	if token := a.state.CheckToken(structs.NodeMaint); token != "mytoken" {
		t.Fatalf("expected 'mytoken', got: '%s'", token)
	}

	// Ensure the reason was set in notes
	if check.Notes != "broken" {
		t.Fatalf("bad: %#v", check)
	}

	// Leave maintenance mode
	a.DisableNodeMaintenance()

	// Ensure the check was deregistered
	if _, ok := a.state.Checks()[structs.NodeMaint]; ok {
		t.Fatalf("should have deregistered critical node check")
	}

	// Enter maintenance mode without passing a reason
	a.EnableNodeMaintenance("", "")

	// Make sure the check was registered with the default note
	check, ok = a.state.Checks()[structs.NodeMaint]
	if !ok {
		t.Fatalf("should have registered critical node check")
	}
	if check.Notes != defaultNodeMaintReason {
		t.Fatalf("bad: %#v", check)
	}
}

func TestAgent_checkStateSnapshot(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// First register a service
	svc := &structs.NodeService{
		ID:      "redis",
		Service: "redis",
		Tags:    []string{"foo"},
		Port:    8000,
	}
	if err := a.AddService(svc, nil, false, ""); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a check
	check1 := &structs.HealthCheck{
		Node:        a.Config.NodeName,
		CheckID:     "service:redis",
		Name:        "redischeck",
		Status:      api.HealthPassing,
		ServiceID:   "redis",
		ServiceName: "redis",
	}
	if err := a.AddCheck(check1, nil, true, ""); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Snapshot the state
	snap := a.snapshotCheckState()

	// Unload all of the checks
	if err := a.unloadChecks(); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Reload the checks
	if err := a.loadChecks(a.Config); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Restore the state
	a.restoreCheckState(snap)

	// Search for the check
	out, ok := a.state.Checks()[check1.CheckID]
	if !ok {
		t.Fatalf("check should have been registered")
	}

	// Make sure state was restored
	if out.Status != api.HealthPassing {
		t.Fatalf("should have restored check state")
	}
}

func TestAgent_loadChecks_checkFails(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// Persist a health check with an invalid service ID
	check := &structs.HealthCheck{
		Node:      a.Config.NodeName,
		CheckID:   "service:redis",
		Name:      "redischeck",
		Status:    api.HealthPassing,
		ServiceID: "nope",
	}
	if err := a.persistCheck(check, nil); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check to make sure the check was persisted
	checkHash := checkIDHash(check.CheckID)
	checkPath := filepath.Join(a.Config.DataDir, checksDir, checkHash)
	if _, err := os.Stat(checkPath); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Try loading the checks from the persisted files
	if err := a.loadChecks(a.Config); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Ensure the erroneous check was purged
	if _, err := os.Stat(checkPath); err == nil {
		t.Fatalf("should have purged check")
	}
}

func TestAgent_persistCheckState(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// Create the TTL check to persist
	check := &CheckTTL{
		CheckID: "check1",
		TTL:     10 * time.Minute,
	}

	// Persist some check state for the check
	err := a.persistCheckState(check, api.HealthCritical, "nope")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check the persisted file exists and has the content
	file := filepath.Join(a.Config.DataDir, checkStateDir, stringHash("check1"))
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Decode the state
	var p persistedCheckState
	if err := json.Unmarshal(buf, &p); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check the fields
	if p.CheckID != "check1" {
		t.Fatalf("bad: %#v", p)
	}
	if p.Output != "nope" {
		t.Fatalf("bad: %#v", p)
	}
	if p.Status != api.HealthCritical {
		t.Fatalf("bad: %#v", p)
	}

	// Check the expiration time was set
	if p.Expires < time.Now().Unix() {
		t.Fatalf("bad: %#v", p)
	}
}

func TestAgent_loadCheckState(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// Create a check whose state will expire immediately
	check := &CheckTTL{
		CheckID: "check1",
		TTL:     0,
	}

	// Persist the check state
	err := a.persistCheckState(check, api.HealthPassing, "yup")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Try to load the state
	health := &structs.HealthCheck{
		CheckID: "check1",
		Status:  api.HealthCritical,
	}
	if err := a.loadCheckState(health); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Should not have restored the status due to expiration
	if health.Status != api.HealthCritical {
		t.Fatalf("bad: %#v", health)
	}
	if health.Output != "" {
		t.Fatalf("bad: %#v", health)
	}

	// Should have purged the state
	file := filepath.Join(a.Config.DataDir, checksDir, stringHash("check1"))
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Fatalf("should have purged state")
	}

	// Set a TTL which will not expire before we check it
	check.TTL = time.Minute
	err = a.persistCheckState(check, api.HealthPassing, "yup")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Try to load
	if err := a.loadCheckState(health); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Should have restored
	if health.Status != api.HealthPassing {
		t.Fatalf("bad: %#v", health)
	}
	if health.Output != "yup" {
		t.Fatalf("bad: %#v", health)
	}
}

func TestAgent_purgeCheckState(t *testing.T) {
	t.Parallel()
	a := NewTestAgent(t.Name(), "")
	defer a.Shutdown()

	// No error if the state does not exist
	if err := a.purgeCheckState("check1"); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Persist some state to the data dir
	check := &CheckTTL{
		CheckID: "check1",
		TTL:     time.Minute,
	}
	err := a.persistCheckState(check, api.HealthPassing, "yup")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Purge the check state
	if err := a.purgeCheckState("check1"); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Removed the file
	file := filepath.Join(a.Config.DataDir, checkStateDir, stringHash("check1"))
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Fatalf("should have removed file")
	}
}

func TestAgent_GetCoordinate(t *testing.T) {
	t.Parallel()
	check := func(server bool) {
		a := NewTestAgent(t.Name(), `
			server = true
		`)
		defer a.Shutdown()

		// This doesn't verify the returned coordinate, but it makes
		// sure that the agent chooses the correct Serf instance,
		// depending on how it's configured as a client or a server.
		// If it chooses the wrong one, this will crash.
		if _, err := a.GetLANCoordinate(); err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	check(true)
	check(false)
}
