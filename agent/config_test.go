package agent

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/agent/structs"
	"github.com/pascaldekloe/goe/verify"
)

// done(fs): func TestConfigEncryptBytes(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs): 	// Test with some input
// done(fs): 	src := []byte("abc")
// done(fs): 	c := &Config{
// done(fs): 		EncryptKey: base64.StdEncoding.EncodeToString(src),
// done(fs): 	}
// done(fs):
// done(fs): 	result, err := c.EncryptBytes()
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	if !bytes.Equal(src, result) {
// done(fs): 		t.Fatalf("bad: %#v", result)
// done(fs): 	}
// done(fs):
// done(fs): 	// Test with no input
// done(fs): 	c = &Config{}
// done(fs): 	result, err = c.EncryptBytes()
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	if len(result) > 0 {
// done(fs): 		t.Fatalf("bad: %#v", result)
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestDecodeConfig(t *testing.T) {
// done(fs): 	tests := []struct {
// done(fs): 		desc             string
// done(fs): 		in               string
// done(fs): 		c                *Config
// done(fs): 		err              error
// done(fs): 		parseTemplateErr error
// done(fs): 	}{
// done(fs): 		// special flows
// done(fs): 		{
// done(fs): 			in:  `{"bad": "no way jose"}`,
// done(fs): 			err: errors.New("Config has invalid keys: bad"),
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:               `{"advertise_addr":"unix:///path/to/file"}`,
// done(fs): 			parseTemplateErr: errors.New("Failed to parse Advertise address: unix:///path/to/file"),
// done(fs): 			c:                &Config{AdvertiseAddr: "unix:///path/to/file"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:               `{"advertise_addr_wan":"unix:///path/to/file"}`,
// done(fs): 			parseTemplateErr: errors.New("Failed to parse Advertise WAN address: unix:///path/to/file"),
// done(fs): 			c:                &Config{AdvertiseAddrWan: "unix:///path/to/file"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:               `{"addresses":{"http":"notunix://blah"}}`,
// done(fs): 			parseTemplateErr: errors.New("Failed to parse HTTP address, \"notunix://blah\" is not a valid IP address or socket"),
// done(fs): 			c:                &Config{Addresses: AddressConfig{HTTP: "notunix://blah"}},
// done(fs): 		},
// done(fs):
// done(fs): 		// happy flows in alphabetical order
// done(fs): 		{
// done(fs): 			in: `{"acl_agent_master_token":"a"}`,
// done(fs): 			c:  &Config{ACLAgentMasterToken: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_agent_token":"a"}`,
// done(fs): 			c:  &Config{ACLAgentToken: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_datacenter":"a"}`,
// done(fs): 			c:  &Config{ACLDatacenter: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_default_policy":"a"}`,
// done(fs): 			c:  &Config{ACLDefaultPolicy: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_down_policy":"a"}`,
// done(fs): 			c:  &Config{ACLDownPolicy: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_enforce_version_8":true}`,
// done(fs): 			c:  &Config{ACLEnforceVersion8: Bool(true)},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_master_token":"a"}`,
// done(fs): 			c:  &Config{ACLMasterToken: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_replication_token":"a"}`,
// done(fs): 			c: &Config{
// done(fs): 				EnableACLReplication: true,
// done(fs): 				ACLReplicationToken:  "a",
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_token":"a"}`,
// done(fs): 			c:  &Config{ACLToken: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"acl_ttl":"2s"}`,
// done(fs): 			c:  &Config{ACLTTL: 2 * time.Second, ACLTTLRaw: "2s"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"dns":"1.2.3.4"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{DNS: "1.2.3.4"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"dns":"{{\"1.2.3.4\"}}"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{DNS: "1.2.3.4"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"http":"1.2.3.4"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{HTTP: "1.2.3.4"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"http":"unix:///var/foo/bar"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{HTTP: "unix:///var/foo/bar"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"http":"{{\"1.2.3.4\"}}"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{HTTP: "1.2.3.4"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"https":"1.2.3.4"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{HTTPS: "1.2.3.4"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"https":"unix:///var/foo/bar"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{HTTPS: "unix:///var/foo/bar"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"https":"{{\"1.2.3.4\"}}"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{HTTPS: "1.2.3.4"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"addresses":{"rpc":"a"}}`,
// done(fs): 			c:  &Config{Addresses: AddressConfig{RPC: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"advertise_addr":"1.2.3.4"}`,
// done(fs): 			c:  &Config{AdvertiseAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"advertise_addr":"{{\"1.2.3.4\"}}"}`,
// done(fs): 			c:  &Config{AdvertiseAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"advertise_addr_wan":"1.2.3.4"}`,
// done(fs): 			c:  &Config{AdvertiseAddrWan: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"advertise_addr_wan":"{{\"1.2.3.4\"}}"}`,
// done(fs): 			c:  &Config{AdvertiseAddrWan: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"advertise_addrs":{"rpc":"1.2.3.4:5678"}}`,
// done(fs): 			c: &Config{
// done(fs): 				AdvertiseAddrs: AdvertiseAddrsConfig{
// done(fs): 					RPC:    &net.TCPAddr{IP: net.ParseIP("1.2.3.4"), Port: 5678},
// done(fs): 					RPCRaw: "1.2.3.4:5678",
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"advertise_addrs":{"serf_lan":"1.2.3.4:5678"}}`,
// done(fs): 			c: &Config{
// done(fs): 				AdvertiseAddrs: AdvertiseAddrsConfig{
// done(fs): 					SerfLan:    &net.TCPAddr{IP: net.ParseIP("1.2.3.4"), Port: 5678},
// done(fs): 					SerfLanRaw: "1.2.3.4:5678",
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"advertise_addrs":{"serf_wan":"1.2.3.4:5678"}}`,
// done(fs): 			c: &Config{
// done(fs): 				AdvertiseAddrs: AdvertiseAddrsConfig{
// done(fs): 					SerfWan:    &net.TCPAddr{IP: net.ParseIP("1.2.3.4"), Port: 5678},
// done(fs): 					SerfWanRaw: "1.2.3.4:5678",
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"atlas_acl_token":"a"}`,
// done(fs): 			c:  &Config{DeprecatedAtlasACLToken: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"atlas_endpoint":"a"}`,
// done(fs): 			c:  &Config{DeprecatedAtlasEndpoint: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"atlas_infrastructure":"a"}`,
// done(fs): 			c:  &Config{DeprecatedAtlasInfrastructure: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"atlas_join":true}`,
// done(fs): 			c:  &Config{DeprecatedAtlasJoin: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"atlas_token":"a"}`,
// done(fs): 			c:  &Config{DeprecatedAtlasToken: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"autopilot":{"cleanup_dead_servers":true}}`,
// done(fs): 			c:  &Config{Autopilot: Autopilot{CleanupDeadServers: Bool(true)}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"autopilot":{"disable_upgrade_migration":true}}`,
// done(fs): 			c:  &Config{Autopilot: Autopilot{DisableUpgradeMigration: Bool(true)}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"autopilot":{"upgrade_version_tag":"rev"}}`,
// done(fs): 			c:  &Config{Autopilot: Autopilot{UpgradeVersionTag: "rev"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"autopilot":{"last_contact_threshold":"2s"}}`,
// done(fs): 			c:  &Config{Autopilot: Autopilot{LastContactThreshold: Duration(2 * time.Second), LastContactThresholdRaw: "2s"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"autopilot":{"max_trailing_logs":10}}`,
// done(fs): 			c:  &Config{Autopilot: Autopilot{MaxTrailingLogs: Uint64(10)}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"autopilot":{"server_stabilization_time":"2s"}}`,
// done(fs): 			c:  &Config{Autopilot: Autopilot{ServerStabilizationTime: Duration(2 * time.Second), ServerStabilizationTimeRaw: "2s"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"autopilot":{"cleanup_dead_servers":true}}`,
// done(fs): 			c:  &Config{Autopilot: Autopilot{CleanupDeadServers: Bool(true)}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"bind_addr":"1.2.3.4"}`,
// done(fs): 			c:  &Config{BindAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"bind_addr":"{{\"1.2.3.4\"}}"}`,
// done(fs): 			c:  &Config{BindAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"bootstrap":true}`,
// done(fs): 			c:  &Config{Bootstrap: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"bootstrap_expect":3}`,
// done(fs): 			c:  &Config{BootstrapExpect: 3},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ca_file":"a"}`,
// done(fs): 			c:  &Config{CAFile: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ca_path":"a"}`,
// done(fs): 			c:  &Config{CAPath: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"check_update_interval":"2s"}`,
// done(fs): 			c:  &Config{CheckUpdateInterval: 2 * time.Second, CheckUpdateIntervalRaw: "2s"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"cert_file":"a"}`,
// done(fs): 			c:  &Config{CertFile: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"client_addr":"1.2.3.4"}`,
// done(fs): 			c:  &Config{ClientAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"client_addr":"{{\"1.2.3.4\"}}"}`,
// done(fs): 			c:  &Config{ClientAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"data_dir":"a"}`,
// done(fs): 			c:  &Config{DataDir: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"datacenter":"a"}`,
// done(fs): 			c:  &Config{Datacenter: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"disable_coordinates":true}`,
// done(fs): 			c:  &Config{DisableCoordinates: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"disable_host_node_id":false}`,
// done(fs): 			c:  &Config{DisableHostNodeID: Bool(false)},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"allow_stale":true}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{AllowStale: Bool(true)}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"disable_compression":true}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{DisableCompression: true}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"enable_truncate":true}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{EnableTruncate: true}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"max_stale":"2s"}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{MaxStale: 2 * time.Second, MaxStaleRaw: "2s"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"node_ttl":"2s"}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{NodeTTL: 2 * time.Second, NodeTTLRaw: "2s"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"only_passing":true}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{OnlyPassing: true}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"recursor_timeout":"2s"}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{RecursorTimeout: 2 * time.Second, RecursorTimeoutRaw: "2s"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"service_ttl":{"*":"2s","a":"456s"}}}`,
// done(fs): 			c: &Config{
// done(fs): 				DNSConfig: DNSConfig{
// done(fs): 					ServiceTTL:    map[string]time.Duration{"*": 2 * time.Second, "a": 456 * time.Second},
// done(fs): 					ServiceTTLRaw: map[string]string{"*": "2s", "a": "456s"},
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dns_config":{"udp_answer_limit":123}}`,
// done(fs): 			c:  &Config{DNSConfig: DNSConfig{UDPAnswerLimit: 123}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"disable_anonymous_signature":true}`,
// done(fs): 			c:  &Config{DisableAnonymousSignature: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"disable_remote_exec":false}`,
// done(fs): 			c:  &Config{DisableRemoteExec: Bool(false)},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"disable_update_check":true}`,
// done(fs): 			c:  &Config{DisableUpdateCheck: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dogstatsd_addr":"a"}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{DogStatsdAddr: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"dogstatsd_tags":["a:b","c:d"]}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{DogStatsdTags: []string{"a:b", "c:d"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"domain":"a"}`,
// done(fs): 			c:  &Config{Domain: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"enable_acl_replication":true}`,
// done(fs): 			c:  &Config{EnableACLReplication: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"enable_debug":true}`,
// done(fs): 			c:  &Config{EnableDebug: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"enable_syslog":true}`,
// done(fs): 			c:  &Config{EnableSyslog: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"disable_keyring_file":true}`,
// done(fs): 			c:  &Config{DisableKeyringFile: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"enable_script_checks":true}`,
// done(fs): 			c:  &Config{EnableScriptChecks: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"encrypt_verify_incoming":true}`,
// done(fs): 			c:  &Config{EncryptVerifyIncoming: Bool(true)},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"encrypt_verify_outgoing":true}`,
// done(fs): 			c:  &Config{EncryptVerifyOutgoing: Bool(true)},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"http_config":{"block_endpoints":["a","b","c","d"]}}`,
// done(fs): 			c:  &Config{HTTPConfig: HTTPConfig{BlockEndpoints: []string{"a", "b", "c", "d"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"http_api_response_headers":{"a":"b","c":"d"}}`,
// done(fs): 			c:  &Config{HTTPConfig: HTTPConfig{ResponseHeaders: map[string]string{"a": "b", "c": "d"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"http_config":{"response_headers":{"a":"b","c":"d"}}}`,
// done(fs): 			c:  &Config{HTTPConfig: HTTPConfig{ResponseHeaders: map[string]string{"a": "b", "c": "d"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"key_file":"a"}`,
// done(fs): 			c:  &Config{KeyFile: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"leave_on_terminate":true}`,
// done(fs): 			c:  &Config{LeaveOnTerm: Bool(true)},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"limits": {"rpc_rate": 100, "rpc_max_burst": 50}}}`,
// done(fs): 			c:  &Config{Limits: Limits{RPCRate: 100, RPCMaxBurst: 50}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"log_level":"a"}`,
// done(fs): 			c:  &Config{LogLevel: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"node_id":"a"}`,
// done(fs): 			c:  &Config{NodeID: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"node_meta":{"a":"b","c":"d"}}`,
// done(fs): 			c:  &Config{Meta: map[string]string{"a": "b", "c": "d"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"node_name":"a"}`,
// done(fs): 			c:  &Config{NodeName: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"performance": { "raft_multiplier": 3 }}`,
// done(fs): 			c:  &Config{Performance: Performance{RaftMultiplier: 3}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:  `{"performance": { "raft_multiplier": 11 }}`,
// done(fs): 			err: errors.New("Performance.RaftMultiplier must be <= 10"),
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"pid_file":"a"}`,
// done(fs): 			c:  &Config{PidFile: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ports":{"dns":1234}}`,
// done(fs): 			c:  &Config{Ports: PortConfig{DNS: 1234}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ports":{"http":1234}}`,
// done(fs): 			c:  &Config{Ports: PortConfig{HTTP: 1234}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ports":{"https":1234}}`,
// done(fs): 			c:  &Config{Ports: PortConfig{HTTPS: 1234}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ports":{"serf_lan":1234}}`,
// done(fs): 			c:  &Config{Ports: PortConfig{SerfLan: 1234}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ports":{"serf_wan":1234}}`,
// done(fs): 			c:  &Config{Ports: PortConfig{SerfWan: 1234}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ports":{"server":1234}}`,
// done(fs): 			c:  &Config{Ports: PortConfig{Server: 1234}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ports":{"rpc":1234}}`,
// done(fs): 			c:  &Config{Ports: PortConfig{RPC: 1234}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"raft_protocol":3}`,
// done(fs): 			c:  &Config{RaftProtocol: 3},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:  `{"reconnect_timeout":"4h"}`,
// done(fs): 			err: errors.New("ReconnectTimeoutLan must be >= 8h0m0s"),
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"reconnect_timeout":"8h"}`,
// done(fs): 			c:  &Config{ReconnectTimeoutLan: 8 * time.Hour, ReconnectTimeoutLanRaw: "8h"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:  `{"reconnect_timeout_wan":"4h"}`,
// done(fs): 			err: errors.New("ReconnectTimeoutWan must be >= 8h0m0s"),
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"reconnect_timeout_wan":"8h"}`,
// done(fs): 			c:  &Config{ReconnectTimeoutWan: 8 * time.Hour, ReconnectTimeoutWanRaw: "8h"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"recursor":"a"}`,
// done(fs): 			c:  &Config{DNSRecursor: "a", DNSRecursors: []string{"a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"recursors":["a","b"]}`,
// done(fs): 			c:  &Config{DNSRecursors: []string{"a", "b"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"rejoin_after_leave":true}`,
// done(fs): 			c:  &Config{RejoinAfterLeave: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"retry_interval":"2s"}`,
// done(fs): 			c:  &Config{RetryInterval: 2 * time.Second, RetryIntervalRaw: "2s"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"retry_interval_wan":"2s"}`,
// done(fs): 			c:  &Config{RetryIntervalWan: 2 * time.Second, RetryIntervalWanRaw: "2s"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"retry_join":["a","b"]}`,
// done(fs): 			c:  &Config{RetryJoin: []string{"a", "b"}},
// done(fs): 		},
// done(fs): 		// todo(fs): temporarily disabling tests after moving the code
// done(fs): 		// todo(fs): to patch the deprecated retry-join flags to command/agent.go
// done(fs): 		// todo(fs): where it cannot be tested.
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_azure":{"client_id":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=azure client_id=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_azure":{"tag_name":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=azure tag_name=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_azure":{"tag_value":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=azure tag_value=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_azure":{"secret_access_key":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=azure secret_access_key=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_azure":{"subscription_id":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=azure subscription_id=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_azure":{"tenant_id":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=azure tenant_id=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_ec2":{"access_key_id":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=aws access_key_id=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_ec2":{"region":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=aws region=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_ec2":{"tag_key":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=aws tag_key=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_ec2":{"tag_value":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=aws tag_value=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_ec2":{"secret_access_key":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=aws secret_access_key=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_gce":{"credentials_file":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=gce credentials_file=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_gce":{"project_name":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=gce project_name=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_gce":{"tag_value":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=gce tag_value=a"}},
// done(fs): 		//		},
// done(fs): 		//		{
// done(fs): 		//			in: `{"retry_join_gce":{"zone_pattern":"a"}}`,
// done(fs): 		//			c:  &Config{RetryJoin: []string{"provider=gce zone_pattern=a"}},
// done(fs): 		//		},
// done(fs): 		{
// done(fs): 			in: `{"retry_join_wan":["a","b"]}`,
// done(fs): 			c:  &Config{RetryJoinWan: []string{"a", "b"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"retry_max":123}`,
// done(fs): 			c:  &Config{RetryMaxAttempts: 123},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"retry_max_wan":123}`,
// done(fs): 			c:  &Config{RetryMaxAttemptsWan: 123},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"serf_lan_bind":"1.2.3.4"}`,
// done(fs): 			c:  &Config{SerfLanBindAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:               `{"serf_lan_bind":"unix:///var/foo/bar"}`,
// done(fs): 			c:                &Config{SerfLanBindAddr: "unix:///var/foo/bar"},
// done(fs): 			parseTemplateErr: errors.New("Failed to parse Serf LAN address: unix:///var/foo/bar"),
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"serf_lan_bind":"{{\"1.2.3.4\"}}"}`,
// done(fs): 			c:  &Config{SerfLanBindAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"serf_wan_bind":"1.2.3.4"}`,
// done(fs): 			c:  &Config{SerfWanBindAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in:               `{"serf_wan_bind":"unix:///var/foo/bar"}`,
// done(fs): 			c:                &Config{SerfWanBindAddr: "unix:///var/foo/bar"},
// done(fs): 			parseTemplateErr: errors.New("Failed to parse Serf WAN address: unix:///var/foo/bar"),
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"serf_wan_bind":"{{\"1.2.3.4\"}}"}`,
// done(fs): 			c:  &Config{SerfWanBindAddr: "1.2.3.4"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"server":true}`,
// done(fs): 			c:  &Config{Server: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"server_name":"a"}`,
// done(fs): 			c:  &Config{ServerName: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"session_ttl_min":"2s"}`,
// done(fs): 			c:  &Config{SessionTTLMin: 2 * time.Second, SessionTTLMinRaw: "2s"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"skip_leave_on_interrupt":true}`,
// done(fs): 			c:  &Config{SkipLeaveOnInt: Bool(true)},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"start_join":["a","b"]}`,
// done(fs): 			c:  &Config{StartJoin: []string{"a", "b"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"start_join_wan":["a","b"]}`,
// done(fs): 			c:  &Config{StartJoinWan: []string{"a", "b"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"statsd_addr":"a"}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{StatsdAddr: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"statsite_addr":"a"}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{StatsiteAddr: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"statsite_prefix":"a"}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{StatsitePrefix: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"syslog_facility":"a"}`,
// done(fs): 			c:  &Config{SyslogFacility: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_api_app":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusAPIApp: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_api_token":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusAPIToken: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_api_url":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusAPIURL: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_broker_id":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusBrokerID: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_broker_select_tag":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusBrokerSelectTag: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_check_display_name":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusCheckDisplayName: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_check_force_metric_activation":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusCheckForceMetricActivation: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_check_id":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusCheckID: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_check_instance_id":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusCheckInstanceID: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_check_search_tag":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusCheckSearchTag: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_check_tags":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusCheckTags: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_submission_interval":"2s"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusSubmissionInterval: "2s"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"circonus_submission_url":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{CirconusCheckSubmissionURL: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"disable_hostname":true}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{DisableHostname: true}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"dogstatsd_addr":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{DogStatsdAddr: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"dogstatsd_tags":["a","b"]}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{DogStatsdTags: []string{"a", "b"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"filter_default":true}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{FilterDefault: Bool(true)}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"prefix_filter":["+consul.metric","-consul.othermetric"]}}`,
// done(fs): 			c: &Config{Telemetry: Telemetry{
// done(fs): 				PrefixFilter:    []string{"+consul.metric", "-consul.othermetric"},
// done(fs): 				AllowedPrefixes: []string{"consul.metric"},
// done(fs): 				BlockedPrefixes: []string{"consul.othermetric"},
// done(fs): 			}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"statsd_address":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{StatsdAddr: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"statsite_address":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{StatsiteAddr: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"telemetry":{"statsite_prefix":"a"}}`,
// done(fs): 			c:  &Config{Telemetry: Telemetry{StatsitePrefix: "a"}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"tls_cipher_suites":"TLS_RSA_WITH_AES_256_CBC_SHA"}`,
// done(fs): 			c: &Config{
// done(fs): 				TLSCipherSuites:    []uint16{tls.TLS_RSA_WITH_AES_256_CBC_SHA},
// done(fs): 				TLSCipherSuitesRaw: "TLS_RSA_WITH_AES_256_CBC_SHA",
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"tls_min_version":"a"}`,
// done(fs): 			c:  &Config{TLSMinVersion: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"tls_prefer_server_cipher_suites":true}`,
// done(fs): 			c:  &Config{TLSPreferServerCipherSuites: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"translate_wan_addrs":true}`,
// done(fs): 			c:  &Config{TranslateWanAddrs: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ui":true}`,
// done(fs): 			c:  &Config{EnableUI: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"ui_dir":"a"}`,
// done(fs): 			c:  &Config{UIDir: "a"},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"unix_sockets":{"user":"a"}}`,
// done(fs): 			c:  &Config{UnixSockets: UnixSocketConfig{UnixSocketPermissions{Usr: "a"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"unix_sockets":{"group":"a"}}`,
// done(fs): 			c:  &Config{UnixSockets: UnixSocketConfig{UnixSocketPermissions{Grp: "a"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"unix_sockets":{"mode":"a"}}`,
// done(fs): 			c:  &Config{UnixSockets: UnixSocketConfig{UnixSocketPermissions{Perms: "a"}}},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"verify_incoming":true}`,
// done(fs): 			c:  &Config{VerifyIncoming: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"verify_incoming_https":true}`,
// done(fs): 			c:  &Config{VerifyIncomingHTTPS: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"verify_incoming_rpc":true}`,
// done(fs): 			c:  &Config{VerifyIncomingRPC: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"verify_outgoing":true}`,
// done(fs): 			c:  &Config{VerifyOutgoing: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"verify_server_hostname":true}`,
// done(fs): 			c:  &Config{VerifyServerHostname: true},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			in: `{"watches":[{"type":"a","prefix":"b","handler":"c"}]}`,
// done(fs): 			c: &Config{
// done(fs): 				Watches: []map[string]interface{}{
// done(fs): 					map[string]interface{}{
// done(fs): 						"type":    "a",
// done(fs): 						"prefix":  "b",
// done(fs): 						"handler": "c",
// done(fs): 					},
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs):
// done(fs): 		// complex flows
// done(fs): 		{
// done(fs): 			desc: "single service with check",
// done(fs): 			in: `{
// done(fs): 					"service": {
// done(fs): 						"ID": "a",
// done(fs): 						"Name": "b",
// done(fs): 						"Tags": ["c", "d"],
// done(fs): 						"Address": "e",
// done(fs): 						"Token": "f",
// done(fs): 						"Port": 123,
// done(fs): 						"EnableTagOverride": true,
// done(fs): 						"Check": {
// done(fs): 							"CheckID": "g",
// done(fs): 							"Name": "h",
// done(fs): 							"Status": "i",
// done(fs): 							"Notes": "j",
// done(fs): 							"Script": "k",
// done(fs): 							"HTTP": "l",
// done(fs): 							"Header": {"a":["b"], "c":["d", "e"]},
// done(fs): 							"Method": "x",
// done(fs): 							"TCP": "m",
// done(fs): 							"DockerContainerID": "n",
// done(fs): 							"Shell": "o",
// done(fs): 							"TLSSkipVerify": true,
// done(fs): 							"Interval": "2s",
// done(fs): 							"Timeout": "3s",
// done(fs): 							"TTL": "4s",
// done(fs): 							"DeregisterCriticalServiceAfter": "5s"
// done(fs): 						}
// done(fs): 					}
// done(fs): 				}`,
// done(fs): 			c: &Config{
// done(fs): 				Services: []*structs.ServiceDefinition{
// done(fs): 					&structs.ServiceDefinition{
// done(fs): 						ID:                "a",
// done(fs): 						Name:              "b",
// done(fs): 						Tags:              []string{"c", "d"},
// done(fs): 						Address:           "e",
// done(fs): 						Port:              123,
// done(fs): 						Token:             "f",
// done(fs): 						EnableTagOverride: true,
// done(fs): 						Check: structs.CheckType{
// done(fs): 							CheckID:           "g",
// done(fs): 							Name:              "h",
// done(fs): 							Status:            "i",
// done(fs): 							Notes:             "j",
// done(fs): 							Script:            "k",
// done(fs): 							HTTP:              "l",
// done(fs): 							Header:            map[string][]string{"a": []string{"b"}, "c": []string{"d", "e"}},
// done(fs): 							Method:            "x",
// done(fs): 							TCP:               "m",
// done(fs): 							DockerContainerID: "n",
// done(fs): 							Shell:             "o",
// done(fs): 							TLSSkipVerify:     true,
// done(fs): 							Interval:          2 * time.Second,
// done(fs): 							Timeout:           3 * time.Second,
// done(fs): 							TTL:               4 * time.Second,
// done(fs): 							DeregisterCriticalServiceAfter: 5 * time.Second,
// done(fs): 						},
// done(fs): 					},
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			desc: "single service with multiple checks",
// done(fs): 			in: `{
// done(fs): 					"service": {
// done(fs): 						"ID": "a",
// done(fs): 						"Name": "b",
// done(fs): 						"Tags": ["c", "d"],
// done(fs): 						"Address": "e",
// done(fs): 						"Token": "f",
// done(fs): 						"Port": 123,
// done(fs): 						"EnableTagOverride": true,
// done(fs): 						"Checks": [
// done(fs): 							{
// done(fs): 								"CheckID": "g",
// done(fs): 								"Name": "h",
// done(fs): 								"Status": "i",
// done(fs): 								"Notes": "j",
// done(fs): 								"Script": "k",
// done(fs): 								"HTTP": "l",
// done(fs): 								"Header": {"a":["b"], "c":["d", "e"]},
// done(fs): 								"Method": "x",
// done(fs): 								"TCP": "m",
// done(fs): 								"DockerContainerID": "n",
// done(fs): 								"Shell": "o",
// done(fs): 								"TLSSkipVerify": true,
// done(fs): 								"Interval": "2s",
// done(fs): 								"Timeout": "3s",
// done(fs): 								"TTL": "4s",
// done(fs): 								"DeregisterCriticalServiceAfter": "5s"
// done(fs): 							},
// done(fs): 							{
// done(fs): 								"CheckID": "gg",
// done(fs): 								"Name": "hh",
// done(fs): 								"Status": "ii",
// done(fs): 								"Notes": "jj",
// done(fs): 								"Script": "kk",
// done(fs): 								"HTTP": "ll",
// done(fs): 								"Header": {"aa":["bb"], "cc":["dd", "ee"]},
// done(fs): 								"Method": "xx",
// done(fs): 								"TCP": "mm",
// done(fs): 								"DockerContainerID": "nn",
// done(fs): 								"Shell": "oo",
// done(fs): 								"TLSSkipVerify": false,
// done(fs): 								"Interval": "22s",
// done(fs): 								"Timeout": "33s",
// done(fs): 								"TTL": "44s",
// done(fs): 								"DeregisterCriticalServiceAfter": "55s"
// done(fs): 							}
// done(fs): 						]
// done(fs): 					}
// done(fs): 				}`,
// done(fs): 			c: &Config{
// done(fs): 				Services: []*structs.ServiceDefinition{
// done(fs): 					&structs.ServiceDefinition{
// done(fs): 						ID:                "a",
// done(fs): 						Name:              "b",
// done(fs): 						Tags:              []string{"c", "d"},
// done(fs): 						Address:           "e",
// done(fs): 						Port:              123,
// done(fs): 						Token:             "f",
// done(fs): 						EnableTagOverride: true,
// done(fs): 						Checks: []*structs.CheckType{
// done(fs): 							{
// done(fs): 								CheckID:           "g",
// done(fs): 								Name:              "h",
// done(fs): 								Status:            "i",
// done(fs): 								Notes:             "j",
// done(fs): 								Script:            "k",
// done(fs): 								HTTP:              "l",
// done(fs): 								Header:            map[string][]string{"a": []string{"b"}, "c": []string{"d", "e"}},
// done(fs): 								Method:            "x",
// done(fs): 								TCP:               "m",
// done(fs): 								DockerContainerID: "n",
// done(fs): 								Shell:             "o",
// done(fs): 								TLSSkipVerify:     true,
// done(fs): 								Interval:          2 * time.Second,
// done(fs): 								Timeout:           3 * time.Second,
// done(fs): 								TTL:               4 * time.Second,
// done(fs): 								DeregisterCriticalServiceAfter: 5 * time.Second,
// done(fs): 							},
// done(fs): 							{
// done(fs): 								CheckID:           "gg",
// done(fs): 								Name:              "hh",
// done(fs): 								Status:            "ii",
// done(fs): 								Notes:             "jj",
// done(fs): 								Script:            "kk",
// done(fs): 								HTTP:              "ll",
// done(fs): 								Header:            map[string][]string{"aa": []string{"bb"}, "cc": []string{"dd", "ee"}},
// done(fs): 								Method:            "xx",
// done(fs): 								TCP:               "mm",
// done(fs): 								DockerContainerID: "nn",
// done(fs): 								Shell:             "oo",
// done(fs): 								TLSSkipVerify:     false,
// done(fs): 								Interval:          22 * time.Second,
// done(fs): 								Timeout:           33 * time.Second,
// done(fs): 								TTL:               44 * time.Second,
// done(fs): 								DeregisterCriticalServiceAfter: 55 * time.Second,
// done(fs): 							},
// done(fs): 						},
// done(fs): 					},
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			desc: "multiple services with check",
// done(fs): 			in: `{
// done(fs): 					"services": [
// done(fs): 						{
// done(fs): 							"ID": "a",
// done(fs): 							"Name": "b",
// done(fs): 							"Tags": ["c", "d"],
// done(fs): 							"Address": "e",
// done(fs): 							"Token": "f",
// done(fs): 							"Port": 123,
// done(fs): 							"EnableTagOverride": true,
// done(fs): 							"Check": {
// done(fs): 								"CheckID": "g",
// done(fs): 								"Name": "h",
// done(fs): 								"Status": "i",
// done(fs): 								"Notes": "j",
// done(fs): 								"Script": "k",
// done(fs): 								"HTTP": "l",
// done(fs): 								"Header": {"a":["b"], "c":["d", "e"]},
// done(fs): 								"Method": "x",
// done(fs): 								"TCP": "m",
// done(fs): 								"DockerContainerID": "n",
// done(fs): 								"Shell": "o",
// done(fs): 								"TLSSkipVerify": true,
// done(fs): 								"Interval": "2s",
// done(fs): 								"Timeout": "3s",
// done(fs): 								"TTL": "4s",
// done(fs): 								"DeregisterCriticalServiceAfter": "5s"
// done(fs): 							}
// done(fs): 						},
// done(fs): 						{
// done(fs): 							"ID": "aa",
// done(fs): 							"Name": "bb",
// done(fs): 							"Tags": ["cc", "dd"],
// done(fs): 							"Address": "ee",
// done(fs): 							"Token": "ff",
// done(fs): 							"Port": 246,
// done(fs): 							"EnableTagOverride": false,
// done(fs): 							"Check": {
// done(fs): 								"CheckID": "gg",
// done(fs): 								"Name": "hh",
// done(fs): 								"Status": "ii",
// done(fs): 								"Notes": "jj",
// done(fs): 								"Script": "kk",
// done(fs): 								"HTTP": "ll",
// done(fs): 								"Header": {"aa":["bb"], "cc":["dd", "ee"]},
// done(fs): 								"Method": "xx",
// done(fs): 								"TCP": "mm",
// done(fs): 								"DockerContainerID": "nn",
// done(fs): 								"Shell": "oo",
// done(fs): 								"TLSSkipVerify": false,
// done(fs): 								"Interval": "22s",
// done(fs): 								"Timeout": "33s",
// done(fs): 								"TTL": "44s",
// done(fs): 								"DeregisterCriticalServiceAfter": "55s"
// done(fs): 							}
// done(fs): 						}
// done(fs): 					]
// done(fs): 				}`,
// done(fs): 			c: &Config{
// done(fs): 				Services: []*structs.ServiceDefinition{
// done(fs): 					&structs.ServiceDefinition{
// done(fs): 						ID:                "a",
// done(fs): 						Name:              "b",
// done(fs): 						Tags:              []string{"c", "d"},
// done(fs): 						Address:           "e",
// done(fs): 						Port:              123,
// done(fs): 						Token:             "f",
// done(fs): 						EnableTagOverride: true,
// done(fs): 						Check: structs.CheckType{
// done(fs): 							CheckID:           "g",
// done(fs): 							Name:              "h",
// done(fs): 							Status:            "i",
// done(fs): 							Notes:             "j",
// done(fs): 							Script:            "k",
// done(fs): 							HTTP:              "l",
// done(fs): 							Header:            map[string][]string{"a": []string{"b"}, "c": []string{"d", "e"}},
// done(fs): 							Method:            "x",
// done(fs): 							TCP:               "m",
// done(fs): 							DockerContainerID: "n",
// done(fs): 							Shell:             "o",
// done(fs): 							TLSSkipVerify:     true,
// done(fs): 							Interval:          2 * time.Second,
// done(fs): 							Timeout:           3 * time.Second,
// done(fs): 							TTL:               4 * time.Second,
// done(fs): 							DeregisterCriticalServiceAfter: 5 * time.Second,
// done(fs): 						},
// done(fs): 					},
// done(fs): 					&structs.ServiceDefinition{
// done(fs): 						ID:                "aa",
// done(fs): 						Name:              "bb",
// done(fs): 						Tags:              []string{"cc", "dd"},
// done(fs): 						Address:           "ee",
// done(fs): 						Port:              246,
// done(fs): 						Token:             "ff",
// done(fs): 						EnableTagOverride: false,
// done(fs): 						Check: structs.CheckType{
// done(fs): 							CheckID:           "gg",
// done(fs): 							Name:              "hh",
// done(fs): 							Status:            "ii",
// done(fs): 							Notes:             "jj",
// done(fs): 							Script:            "kk",
// done(fs): 							HTTP:              "ll",
// done(fs): 							Header:            map[string][]string{"aa": []string{"bb"}, "cc": []string{"dd", "ee"}},
// done(fs): 							Method:            "xx",
// done(fs): 							TCP:               "mm",
// done(fs): 							DockerContainerID: "nn",
// done(fs): 							Shell:             "oo",
// done(fs): 							TLSSkipVerify:     false,
// done(fs): 							Interval:          22 * time.Second,
// done(fs): 							Timeout:           33 * time.Second,
// done(fs): 							TTL:               44 * time.Second,
// done(fs): 							DeregisterCriticalServiceAfter: 55 * time.Second,
// done(fs): 						},
// done(fs): 					},
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs):
// done(fs): 		{
// done(fs): 			desc: "single check",
// done(fs): 			in: `{
// done(fs): 					"check": {
// done(fs): 						"id": "a",
// done(fs): 						"name": "b",
// done(fs): 						"notes": "c",
// done(fs): 						"service_id": "x",
// done(fs): 						"token": "y",
// done(fs): 						"status": "z",
// done(fs): 						"script": "d",
// done(fs): 						"shell": "e",
// done(fs): 						"http": "f",
// done(fs): 						"Header": {"a":["b"], "c":["d", "e"]},
// done(fs): 						"Method": "x",
// done(fs): 						"tcp": "g",
// done(fs): 						"docker_container_id": "h",
// done(fs): 						"tls_skip_verify": true,
// done(fs): 						"interval": "2s",
// done(fs): 						"timeout": "3s",
// done(fs): 						"ttl": "4s",
// done(fs): 						"deregister_critical_service_after": "5s"
// done(fs): 					}
// done(fs): 				}`,
// done(fs): 			c: &Config{
// done(fs): 				Checks: []*structs.CheckDefinition{
// done(fs): 					&structs.CheckDefinition{
// done(fs): 						ID:                "a",
// done(fs): 						Name:              "b",
// done(fs): 						Notes:             "c",
// done(fs): 						ServiceID:         "x",
// done(fs): 						Token:             "y",
// done(fs): 						Status:            "z",
// done(fs): 						Script:            "d",
// done(fs): 						Shell:             "e",
// done(fs): 						HTTP:              "f",
// done(fs): 						Header:            map[string][]string{"a": []string{"b"}, "c": []string{"d", "e"}},
// done(fs): 						Method:            "x",
// done(fs): 						TCP:               "g",
// done(fs): 						DockerContainerID: "h",
// done(fs): 						TLSSkipVerify:     true,
// done(fs): 						Interval:          2 * time.Second,
// done(fs): 						Timeout:           3 * time.Second,
// done(fs): 						TTL:               4 * time.Second,
// done(fs): 						DeregisterCriticalServiceAfter: 5 * time.Second,
// done(fs): 					},
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 		{
// done(fs): 			desc: "multiple checks",
// done(fs): 			in: `{
// done(fs): 					"checks": [
// done(fs): 						{
// done(fs): 							"id": "a",
// done(fs): 							"name": "b",
// done(fs): 							"notes": "c",
// done(fs): 							"service_id": "d",
// done(fs): 							"token": "e",
// done(fs): 							"status": "f",
// done(fs): 							"script": "g",
// done(fs): 							"shell": "h",
// done(fs): 							"http": "i",
// done(fs): 							"Header": {"a":["b"], "c":["d", "e"]},
// done(fs): 							"Method": "x",
// done(fs): 							"tcp": "j",
// done(fs): 							"docker_container_id": "k",
// done(fs): 							"tls_skip_verify": true,
// done(fs): 							"interval": "2s",
// done(fs): 							"timeout": "3s",
// done(fs): 							"ttl": "4s",
// done(fs): 							"deregister_critical_service_after": "5s"
// done(fs): 						},
// done(fs): 						{
// done(fs): 							"id": "aa",
// done(fs): 							"name": "bb",
// done(fs): 							"notes": "cc",
// done(fs): 							"service_id": "dd",
// done(fs): 							"token": "ee",
// done(fs): 							"status": "ff",
// done(fs): 							"script": "gg",
// done(fs): 							"shell": "hh",
// done(fs): 							"http": "ii",
// done(fs): 							"Header": {"aa":["bb"], "cc":["dd", "ee"]},
// done(fs): 							"Method": "xx",
// done(fs): 							"tcp": "jj",
// done(fs): 							"docker_container_id": "kk",
// done(fs): 							"tls_skip_verify": false,
// done(fs): 							"interval": "22s",
// done(fs): 							"timeout": "33s",
// done(fs): 							"ttl": "44s",
// done(fs): 							"deregister_critical_service_after": "55s"
// done(fs): 						}
// done(fs): 					]
// done(fs): 				}`,
// done(fs): 			c: &Config{
// done(fs): 				Checks: []*structs.CheckDefinition{
// done(fs): 					&structs.CheckDefinition{
// done(fs): 						ID:                "a",
// done(fs): 						Name:              "b",
// done(fs): 						Notes:             "c",
// done(fs): 						ServiceID:         "d",
// done(fs): 						Token:             "e",
// done(fs): 						Status:            "f",
// done(fs): 						Script:            "g",
// done(fs): 						Shell:             "h",
// done(fs): 						HTTP:              "i",
// done(fs): 						Header:            map[string][]string{"a": []string{"b"}, "c": []string{"d", "e"}},
// done(fs): 						Method:            "x",
// done(fs): 						TCP:               "j",
// done(fs): 						DockerContainerID: "k",
// done(fs): 						TLSSkipVerify:     true,
// done(fs): 						Interval:          2 * time.Second,
// done(fs): 						Timeout:           3 * time.Second,
// done(fs): 						TTL:               4 * time.Second,
// done(fs): 						DeregisterCriticalServiceAfter: 5 * time.Second,
// done(fs): 					},
// done(fs): 					&structs.CheckDefinition{
// done(fs): 						ID:                "aa",
// done(fs): 						Name:              "bb",
// done(fs): 						Notes:             "cc",
// done(fs): 						ServiceID:         "dd",
// done(fs): 						Token:             "ee",
// done(fs): 						Status:            "ff",
// done(fs): 						Script:            "gg",
// done(fs): 						Shell:             "hh",
// done(fs): 						HTTP:              "ii",
// done(fs): 						Header:            map[string][]string{"aa": []string{"bb"}, "cc": []string{"dd", "ee"}},
// done(fs): 						Method:            "xx",
// done(fs): 						TCP:               "jj",
// done(fs): 						DockerContainerID: "kk",
// done(fs): 						TLSSkipVerify:     false,
// done(fs): 						Interval:          22 * time.Second,
// done(fs): 						Timeout:           33 * time.Second,
// done(fs): 						TTL:               44 * time.Second,
// done(fs): 						DeregisterCriticalServiceAfter: 55 * time.Second,
// done(fs): 					},
// done(fs): 				},
// done(fs): 			},
// done(fs): 		},
// done(fs): 	}
// done(fs):
// done(fs): 	for _, tt := range tests {
// done(fs): 		desc := tt.desc
// done(fs): 		if desc == "" {
// done(fs): 			desc = tt.in
// done(fs): 		}
// done(fs): 		t.Run(desc, func(t *testing.T) {
// done(fs): 			c, err := DecodeConfig(strings.NewReader(tt.in))
// done(fs): 			if got, want := err, tt.err; !reflect.DeepEqual(got, want) {
// done(fs): 				t.Fatalf("got error %v want %v", got, want)
// done(fs): 			}
// done(fs): 			err = c.ResolveTmplAddrs()
// done(fs): 			if got, want := err, tt.parseTemplateErr; !reflect.DeepEqual(got, want) {
// done(fs): 				t.Fatalf("got error %v on ResolveTmplAddrs, expected %v", err, want)
// done(fs): 			}
// done(fs): 			got, want := c, tt.c
// done(fs): 			verify.Values(t, "", got, want)
// done(fs): 		})
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestDecodeConfig_VerifyUniqueListeners(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs): 	tests := []struct {
// done(fs): 		desc string
// done(fs): 		in   string
// done(fs): 		err  error
// done(fs): 	}{
// done(fs): 		{
// done(fs): 			"http_dns1",
// done(fs): 			`{"addresses": {"http": "0.0.0.0", "dns": "127.0.0.1"}, "ports": {"dns": 8000}}`,
// done(fs): 			nil,
// done(fs): 		},
// done(fs): 		{
// done(fs): 			"http_dns IP identical",
// done(fs): 			`{"addresses": {"http": "0.0.0.0", "dns": "0.0.0.0"}, "ports": {"http": 8000, "dns": 8000}}`,
// done(fs): 			errors.New("HTTP address already configured for DNS"),
// done(fs): 		},
// done(fs): 	}
// done(fs):
// done(fs): 	for _, tt := range tests {
// done(fs): 		t.Run(tt.desc, func(t *testing.T) {
// done(fs): 			c, err := DecodeConfig(strings.NewReader(tt.in))
// done(fs): 			if err != nil {
// done(fs): 				t.Fatalf("got error %v want nil", err)
// done(fs): 			}
// done(fs):
// done(fs): 			err = c.VerifyUniqueListeners()
// done(fs): 			if got, want := err, tt.err; !reflect.DeepEqual(got, want) {
// done(fs): 				t.Fatalf("got error %v want %v", got, want)
// done(fs): 			}
// done(fs): 		})
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestDefaultConfig(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs):
// done(fs): 	// ACL flag for Consul version 0.8 features (broken out since we will
// done(fs): 	// eventually remove this).
// done(fs): 	config := DefaultConfig()
// done(fs): 	if *config.ACLEnforceVersion8 != true {
// done(fs): 		t.Fatalf("bad: %#v", config)
// done(fs): 	}
// done(fs):
// done(fs): 	// Remote exec is disabled by default.
// done(fs): 	if *config.DisableRemoteExec != true {
// done(fs): 		t.Fatalf("bad: %#v", config)
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestMergeConfig(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs): 	a := &Config{
// done(fs): 		Bootstrap:              false,
// done(fs): 		BootstrapExpect:        0,
// done(fs): 		Datacenter:             "dc1",
// done(fs): 		DataDir:                "/tmp/foo",
// done(fs): 		Domain:                 "basic",
// done(fs): 		LogLevel:               "debug",
// done(fs): 		NodeID:                 "bar",
// done(fs): 		NodeName:               "foo",
// done(fs): 		ClientAddr:             "127.0.0.1",
// done(fs): 		BindAddr:               "127.0.0.1",
// done(fs): 		AdvertiseAddr:          "127.0.0.1",
// done(fs): 		Server:                 false,
// done(fs): 		LeaveOnTerm:            new(bool),
// done(fs): 		SkipLeaveOnInt:         new(bool),
// done(fs): 		EnableDebug:            false,
// done(fs): 		CheckUpdateIntervalRaw: "8m",
// done(fs): 		RetryIntervalRaw:       "10s",
// done(fs): 		RetryIntervalWanRaw:    "10s",
// done(fs): 		DeprecatedRetryJoinEC2: RetryJoinEC2{
// done(fs): 			Region:          "us-east-1",
// done(fs): 			TagKey:          "Key1",
// done(fs): 			TagValue:        "Value1",
// done(fs): 			AccessKeyID:     "nope",
// done(fs): 			SecretAccessKey: "nope",
// done(fs): 		},
// done(fs): 		Telemetry: Telemetry{
// done(fs): 			DisableHostname: false,
// done(fs): 			StatsdAddr:      "nope",
// done(fs): 			StatsiteAddr:    "nope",
// done(fs): 			StatsitePrefix:  "nope",
// done(fs): 			DogStatsdAddr:   "nope",
// done(fs): 			DogStatsdTags:   []string{"nope"},
// done(fs): 		},
// done(fs): 		Meta: map[string]string{
// done(fs): 			"key": "value1",
// done(fs): 		},
// done(fs): 	}
// done(fs):
// done(fs): 	b := &Config{
// done(fs): 		Limits: Limits{
// done(fs): 			RPCRate:     100,
// done(fs): 			RPCMaxBurst: 50,
// done(fs): 		},
// done(fs): 		Performance: Performance{
// done(fs): 			RaftMultiplier: 99,
// done(fs): 		},
// done(fs): 		Bootstrap:       true,
// done(fs): 		BootstrapExpect: 3,
// done(fs): 		Datacenter:      "dc2",
// done(fs): 		DataDir:         "/tmp/bar",
// done(fs): 		DNSRecursors:    []string{"127.0.0.2:1001"},
// done(fs): 		DNSConfig: DNSConfig{
// done(fs): 			AllowStale:         Bool(false),
// done(fs): 			EnableTruncate:     true,
// done(fs): 			DisableCompression: true,
// done(fs): 			MaxStale:           30 * time.Second,
// done(fs): 			NodeTTL:            10 * time.Second,
// done(fs): 			ServiceTTL: map[string]time.Duration{
// done(fs): 				"api": 10 * time.Second,
// done(fs): 			},
// done(fs): 			UDPAnswerLimit:  4,
// done(fs): 			RecursorTimeout: 30 * time.Second,
// done(fs): 		},
// done(fs): 		Domain:            "other",
// done(fs): 		LogLevel:          "info",
// done(fs): 		NodeID:            "bar",
// done(fs): 		DisableHostNodeID: Bool(false),
// done(fs): 		NodeName:          "baz",
// done(fs): 		ClientAddr:        "127.0.0.2",
// done(fs): 		BindAddr:          "127.0.0.2",
// done(fs): 		AdvertiseAddr:     "127.0.0.2",
// done(fs): 		AdvertiseAddrWan:  "127.0.0.2",
// done(fs): 		Ports: PortConfig{
// done(fs): 			DNS:     1,
// done(fs): 			HTTP:    2,
// done(fs): 			SerfLan: 4,
// done(fs): 			SerfWan: 5,
// done(fs): 			Server:  6,
// done(fs): 			HTTPS:   7,
// done(fs): 		},
// done(fs): 		Addresses: AddressConfig{
// done(fs): 			DNS:   "127.0.0.1",
// done(fs): 			HTTP:  "127.0.0.2",
// done(fs): 			HTTPS: "127.0.0.4",
// done(fs): 		},
// done(fs): 		Segment: "alpha",
// done(fs): 		Segments: []NetworkSegment{
// done(fs): 			{
// done(fs): 				Name:      "alpha",
// done(fs): 				Bind:      "127.0.0.1",
// done(fs): 				Port:      1234,
// done(fs): 				Advertise: "127.0.0.2",
// done(fs): 			},
// done(fs): 		},
// done(fs): 		Server:         true,
// done(fs): 		LeaveOnTerm:    Bool(true),
// done(fs): 		SkipLeaveOnInt: Bool(true),
// done(fs): 		RaftProtocol:   3,
// done(fs): 		Autopilot: Autopilot{
// done(fs): 			CleanupDeadServers:      Bool(true),
// done(fs): 			LastContactThreshold:    Duration(time.Duration(10)),
// done(fs): 			MaxTrailingLogs:         Uint64(10),
// done(fs): 			ServerStabilizationTime: Duration(time.Duration(100)),
// done(fs): 		},
// done(fs): 		EnableDebug:            true,
// done(fs): 		VerifyIncoming:         true,
// done(fs): 		VerifyOutgoing:         true,
// done(fs): 		CAFile:                 "test/ca.pem",
// done(fs): 		CertFile:               "test/cert.pem",
// done(fs): 		KeyFile:                "test/key.pem",
// done(fs): 		TLSMinVersion:          "tls12",
// done(fs): 		Checks:                 []*structs.CheckDefinition{nil},
// done(fs): 		Services:               []*structs.ServiceDefinition{nil},
// done(fs): 		StartJoin:              []string{"1.1.1.1"},
// done(fs): 		StartJoinWan:           []string{"1.1.1.1"},
// done(fs): 		EnableUI:               true,
// done(fs): 		UIDir:                  "/opt/consul-ui",
// done(fs): 		EnableSyslog:           true,
// done(fs): 		RejoinAfterLeave:       true,
// done(fs): 		RetryJoin:              []string{"1.1.1.1"},
// done(fs): 		RetryIntervalRaw:       "10s",
// done(fs): 		RetryInterval:          10 * time.Second,
// done(fs): 		RetryJoinWan:           []string{"1.1.1.1"},
// done(fs): 		RetryIntervalWanRaw:    "10s",
// done(fs): 		RetryIntervalWan:       10 * time.Second,
// done(fs): 		ReconnectTimeoutLanRaw: "24h",
// done(fs): 		ReconnectTimeoutLan:    24 * time.Hour,
// done(fs): 		ReconnectTimeoutWanRaw: "36h",
// done(fs): 		ReconnectTimeoutWan:    36 * time.Hour,
// done(fs): 		EnableScriptChecks:     true,
// done(fs): 		CheckUpdateInterval:    8 * time.Minute,
// done(fs): 		CheckUpdateIntervalRaw: "8m",
// done(fs): 		ACLToken:               "1111",
// done(fs): 		ACLAgentMasterToken:    "2222",
// done(fs): 		ACLAgentToken:          "3333",
// done(fs): 		ACLMasterToken:         "4444",
// done(fs): 		ACLDatacenter:          "dc2",
// done(fs): 		ACLTTL:                 15 * time.Second,
// done(fs): 		ACLTTLRaw:              "15s",
// done(fs): 		ACLDownPolicy:          "deny",
// done(fs): 		ACLDefaultPolicy:       "deny",
// done(fs): 		ACLReplicationToken:    "8765309",
// done(fs): 		ACLEnforceVersion8:     Bool(true),
// done(fs): 		Watches: []map[string]interface{}{
// done(fs): 			map[string]interface{}{
// done(fs): 				"type":    "keyprefix",
// done(fs): 				"prefix":  "foo/",
// done(fs): 				"handler": "foobar",
// done(fs): 			},
// done(fs): 		},
// done(fs): 		DisableRemoteExec: Bool(true),
// done(fs): 		Telemetry: Telemetry{
// done(fs): 			StatsiteAddr:    "127.0.0.1:7250",
// done(fs): 			StatsitePrefix:  "stats_prefix",
// done(fs): 			StatsdAddr:      "127.0.0.1:7251",
// done(fs): 			DisableHostname: true,
// done(fs): 			DogStatsdAddr:   "127.0.0.1:7254",
// done(fs): 			DogStatsdTags:   []string{"tag_1:val_1", "tag_2:val_2"},
// done(fs): 		},
// done(fs): 		Meta: map[string]string{
// done(fs): 			"key": "value2",
// done(fs): 		},
// done(fs): 		DisableUpdateCheck:        true,
// done(fs): 		DisableAnonymousSignature: true,
// done(fs): 		HTTPConfig: HTTPConfig{
// done(fs): 			BlockEndpoints: []string{
// done(fs): 				"/v1/agent/self",
// done(fs): 				"/v1/acl",
// done(fs): 			},
// done(fs): 			ResponseHeaders: map[string]string{
// done(fs): 				"Access-Control-Allow-Origin": "*",
// done(fs): 			},
// done(fs): 		},
// done(fs): 		UnixSockets: UnixSocketConfig{
// done(fs): 			UnixSocketPermissions{
// done(fs): 				Usr:   "500",
// done(fs): 				Grp:   "500",
// done(fs): 				Perms: "0700",
// done(fs): 			},
// done(fs): 		},
// done(fs): 		DeprecatedRetryJoinEC2: RetryJoinEC2{
// done(fs): 			Region:          "us-east-2",
// done(fs): 			TagKey:          "Key2",
// done(fs): 			TagValue:        "Value2",
// done(fs): 			AccessKeyID:     "foo",
// done(fs): 			SecretAccessKey: "bar",
// done(fs): 		},
// done(fs): 		SessionTTLMinRaw: "1000s",
// done(fs): 		SessionTTLMin:    1000 * time.Second,
// done(fs): 		AdvertiseAddrs: AdvertiseAddrsConfig{
// done(fs): 			SerfLan:    &net.TCPAddr{},
// done(fs): 			SerfLanRaw: "127.0.0.5:1231",
// done(fs): 			SerfWan:    &net.TCPAddr{},
// done(fs): 			SerfWanRaw: "127.0.0.5:1232",
// done(fs): 			RPC:        &net.TCPAddr{},
// done(fs): 			RPCRaw:     "127.0.0.5:1233",
// done(fs): 		},
// done(fs): 	}
// done(fs):
// done(fs): 	c := MergeConfig(a, b)
// done(fs):
// done(fs): 	if !reflect.DeepEqual(c, b) {
// done(fs): 		t.Fatalf("should be equal %#v %#v", c, b)
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestReadConfigPaths_badPath(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs): 	_, err := ReadConfigPaths([]string{"/i/shouldnt/exist/ever/rainbows"})
// done(fs): 	if err == nil {
// done(fs): 		t.Fatal("should have err")
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestReadConfigPaths_file(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs): 	tf := testutil.TempFile(t, "consul")
// done(fs): 	tf.Write([]byte(`{"node_name":"bar"}`))
// done(fs): 	tf.Close()
// done(fs): 	defer os.Remove(tf.Name())
// done(fs):
// done(fs): 	config, err := ReadConfigPaths([]string{tf.Name()})
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	if config.NodeName != "bar" {
// done(fs): 		t.Fatalf("bad: %#v", config)
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestReadConfigPaths_dir(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs): 	td := testutil.TempDir(t, "consul")
// done(fs): 	defer os.RemoveAll(td)
// done(fs):
// done(fs): 	err := ioutil.WriteFile(filepath.Join(td, "a.json"),
// done(fs): 		[]byte(`{"node_name": "bar"}`), 0644)
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	err = ioutil.WriteFile(filepath.Join(td, "b.json"),
// done(fs): 		[]byte(`{"node_name": "baz"}`), 0644)
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	// A non-json file, shouldn't be read
// done(fs): 	err = ioutil.WriteFile(filepath.Join(td, "c"),
// done(fs): 		[]byte(`{"node_name": "bad"}`), 0644)
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	// An empty file shouldn't be read
// done(fs): 	err = ioutil.WriteFile(filepath.Join(td, "d.json"),
// done(fs): 		[]byte{}, 0664)
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	config, err := ReadConfigPaths([]string{td})
// done(fs): 	if err != nil {
// done(fs): 		t.Fatalf("err: %s", err)
// done(fs): 	}
// done(fs):
// done(fs): 	if config.NodeName != "baz" {
// done(fs): 		t.Fatalf("bad: %#v", config)
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): func TestUnixSockets(t *testing.T) {
// done(fs): 	t.Parallel()
// done(fs): 	if p := socketPath("unix:///path/to/socket"); p != "/path/to/socket" {
// done(fs): 		t.Fatalf("bad: %q", p)
// done(fs): 	}
// done(fs): 	if p := socketPath("notunix://blah"); p != "" {
// done(fs): 		t.Fatalf("bad: %q", p)
// done(fs): 	}
// done(fs): }

func TestCheckDefinitionToCheckType(t *testing.T) {
	t.Parallel()
	got := &structs.CheckDefinition{
		ID:     "id",
		Name:   "name",
		Status: "green",
		Notes:  "notes",

		ServiceID:         "svcid",
		Token:             "tok",
		Script:            "/bin/foo",
		HTTP:              "someurl",
		TCP:               "host:port",
		Interval:          1 * time.Second,
		DockerContainerID: "abc123",
		Shell:             "/bin/ksh",
		TLSSkipVerify:     true,
		Timeout:           2 * time.Second,
		TTL:               3 * time.Second,
		DeregisterCriticalServiceAfter: 4 * time.Second,
	}
	want := &structs.CheckType{
		CheckID: "id",
		Name:    "name",
		Status:  "green",
		Notes:   "notes",

		Script:            "/bin/foo",
		HTTP:              "someurl",
		TCP:               "host:port",
		Interval:          1 * time.Second,
		DockerContainerID: "abc123",
		Shell:             "/bin/ksh",
		TLSSkipVerify:     true,
		Timeout:           2 * time.Second,
		TTL:               3 * time.Second,
		DeregisterCriticalServiceAfter: 4 * time.Second,
	}
	verify.Values(t, "", got.CheckType(), want)
}
