package config

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"bytes"

	"github.com/hashicorp/consul/agent/consul"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/types"
	"github.com/pascaldekloe/goe/verify"
)

// TestConfigFlagsAndEdgecases tests the command line flags and
// edgecases for the config parsing. It provides a test structure which
// checks for warnings on deprecated fields and flags.  These tests
// should check one option at a time if possible and should use generic
// values, e.g. 'a' or 1 instead of 'servicex' or 3306.

func splitIPPort(hostport string) (net.IP, int) {
	h, p, err := net.SplitHostPort(hostport)
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		panic(err)
	}
	return net.ParseIP(h), port
}

func ipAddr(addr string) *net.IPAddr {
	return &net.IPAddr{IP: net.ParseIP(addr)}
}

func tcpAddr(addr string) *net.TCPAddr {
	ip, port := splitIPPort(addr)
	return &net.TCPAddr{IP: ip, Port: port}
}

func udpAddr(addr string) *net.UDPAddr {
	ip, port := splitIPPort(addr)
	return &net.UDPAddr{IP: ip, Port: port}
}

func unixAddr(addr string) *net.UnixAddr {
	if !strings.HasPrefix(addr, "unix://") {
		panic("not a unix socket addr: " + addr)
	}
	return &net.UnixAddr{Net: "unix", Name: addr[len("unix://"):]}
}

func TestConfigFlagsAndEdgecases(t *testing.T) {
	randomString := func(n int) string {
		s := ""
		for ; n > 0; n-- {
			s += "x"
		}
		return s
	}

	metaPairs := func(n int, format string) string {
		var s []string
		for i := 0; i < n; i++ {
			switch format {
			case "json":
				s = append(s, fmt.Sprintf(`"%d":"%d"`, i, i))
			case "hcl":
				s = append(s, fmt.Sprintf(`"%d"="%d"`, i, i))
			default:
				panic("invalid format: " + format)
			}
		}
		switch format {
		case "json":
			return strings.Join(s, ",")
		case "hcl":
			return strings.Join(s, " ")
		default:
			panic("invalid format: " + format)
		}
	}

	tests := []struct {
		desc           string
		flags          []string
		json, jsontail []string
		hcl, hcltail   []string
		patch          func(rt *RuntimeConfig)
		err            string
		warns          []string
		hostname       func() (string, error)
	}{
		// ------------------------------------------------------------
		// cmd line flags
		//

		{
			desc:  "-advertise",
			flags: []string{`-advertise`, `1.2.3.4`},
			patch: func(rt *RuntimeConfig) {
				rt.AdvertiseAddrLAN = tcpAddr("1.2.3.4:8300")
				rt.AdvertiseAddrWAN = tcpAddr("1.2.3.4:8300")
				rt.TaggedAddresses = map[string]string{
					"lan": "1.2.3.4",
					"wan": "1.2.3.4",
				}
			},
		},
		{
			desc:  "-advertise-wan",
			flags: []string{`-advertise-wan`, `1.2.3.4`},
			patch: func(rt *RuntimeConfig) {
				rt.AdvertiseAddrWAN = tcpAddr("1.2.3.4:8300")
				rt.TaggedAddresses = map[string]string{
					"lan": "10.0.0.1",
					"wan": "1.2.3.4",
				}
			},
		},
		{
			desc:  "-advertise and -advertise-wan",
			flags: []string{`-advertise`, `1.2.3.4`, `-advertise-wan`, `5.6.7.8`},
			patch: func(rt *RuntimeConfig) {
				rt.AdvertiseAddrLAN = tcpAddr("1.2.3.4:8300")
				rt.AdvertiseAddrWAN = tcpAddr("5.6.7.8:8300")
				rt.TaggedAddresses = map[string]string{
					"lan": "1.2.3.4",
					"wan": "5.6.7.8",
				}
			},
		},
		{
			desc:  "-bind",
			flags: []string{`-bind`, `1.2.3.4`},
			patch: func(rt *RuntimeConfig) {
				rt.BindAddr = ipAddr("1.2.3.4")
				rt.SerfBindAddrLAN = tcpAddr("1.2.3.4:8301")
				rt.SerfBindAddrWAN = tcpAddr("1.2.3.4:8302")
				rt.AdvertiseAddrLAN = tcpAddr("1.2.3.4:8300")
				rt.AdvertiseAddrWAN = tcpAddr("1.2.3.4:8300")
				rt.RPCAdvertiseAddr = tcpAddr("1.2.3.4:8300")
				rt.RPCBindAddr = tcpAddr("1.2.3.4:8300")
				rt.SerfAdvertiseAddrLAN = tcpAddr("1.2.3.4:8301")
				rt.SerfAdvertiseAddrWAN = tcpAddr("1.2.3.4:8302")
				rt.TaggedAddresses = map[string]string{
					"lan": "1.2.3.4",
					"wan": "1.2.3.4",
				}
			},
		},
		{
			desc:  "-bootstrap",
			flags: []string{`-bootstrap`, `-server`},
			patch: func(rt *RuntimeConfig) {
				rt.Bootstrap = true
				rt.ServerMode = true
				rt.LeaveOnTerm = false
				rt.SkipLeaveOnInt = true
			},
			warns: []string{"bootstrap = true: do not enable unless necessary"},
		},
		{
			desc:  "-bootstrap-expect",
			flags: []string{`-bootstrap-expect`, `3`, `-server`},
			patch: func(rt *RuntimeConfig) {
				rt.BootstrapExpect = 3
				rt.ServerMode = true
				rt.LeaveOnTerm = false
				rt.SkipLeaveOnInt = true
			},
			warns: []string{"bootstrap_expect > 0: expecting 3 servers"},
		},
		{
			desc:  "-client",
			flags: []string{`-client`, `1.2.3.4`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("1.2.3.4")}
				rt.DNSAddrs = []net.Addr{tcpAddr("1.2.3.4:8600"), udpAddr("1.2.3.4:8600")}
				rt.HTTPAddrs = []net.Addr{tcpAddr("1.2.3.4:8500")}
			},
		},
		{
			desc:  "-data-dir",
			flags: []string{`-data-dir`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.DataDir = "a"
			},
		},
		{
			desc:  "-data-dir empty",
			flags: []string{`-data-dir`, ``},
			patch: func(rt *RuntimeConfig) {
				rt.DataDir = ""
			},
			err: "data_dir: cannot be empty",
		},
		{
			desc:  "-data-dir given non-directory",
			flags: []string{`-data-dir`, `runtime_test.go`},
			patch: func(rt *RuntimeConfig) {
				rt.DataDir = "runtime_test.go"
			},
			warns: []string{`data_dir: not a directory: runtime_test.go`},
		},
		{
			desc:  "-datacenter",
			flags: []string{`-datacenter`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.Datacenter = "a"
			},
		},
		{
			desc:  "-dev",
			flags: []string{`-dev`},
			patch: func(rt *RuntimeConfig) {
				rt.AdvertiseAddrLAN = tcpAddr("127.0.0.1:8300")
				rt.AdvertiseAddrWAN = tcpAddr("127.0.0.1:8300")
				rt.BindAddr = ipAddr("127.0.0.1")
				rt.DevMode = true
				rt.DisableAnonymousSignature = true
				rt.DisableKeyringFile = true
				rt.EnableDebug = true
				rt.EnableUI = true
				rt.LeaveOnTerm = false
				rt.LogLevel = "DEBUG"
				rt.RPCAdvertiseAddr = tcpAddr("127.0.0.1:8300")
				rt.RPCBindAddr = tcpAddr("127.0.0.1:8300")
				rt.SerfAdvertiseAddrLAN = tcpAddr("127.0.0.1:8301")
				rt.SerfAdvertiseAddrWAN = tcpAddr("127.0.0.1:8302")
				rt.SerfBindAddrLAN = tcpAddr("127.0.0.1:8301")
				rt.SerfBindAddrWAN = tcpAddr("127.0.0.1:8302")
				rt.ServerMode = true
				rt.SkipLeaveOnInt = true
				rt.TaggedAddresses = map[string]string{"lan": "127.0.0.1", "wan": "127.0.0.1"}
				rt.ConsulConfig = devConsulConfig(consul.DefaultConfig())
			},
		},
		{
			desc:  "-disable-host-node-id",
			flags: []string{`-disable-host-node-id`},
			patch: func(rt *RuntimeConfig) {
				rt.DisableHostNodeID = true
			},
		},
		{
			desc:  "-disable-keyring-file",
			flags: []string{`-disable-keyring-file`},
			patch: func(rt *RuntimeConfig) {
				rt.DisableKeyringFile = true
			},
		},
		{
			desc:  "-dns-port",
			flags: []string{`-dns-port`, `123`},
			patch: func(rt *RuntimeConfig) {
				rt.DNSPort = 123
				rt.DNSAddrs = []net.Addr{tcpAddr("127.0.0.1:123"), udpAddr("127.0.0.1:123")}
			},
		},
		{
			desc:  "-domain",
			flags: []string{`-domain`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.DNSDomain = "a"
			},
		},
		{
			desc:  "-enable-script-checks",
			flags: []string{`-enable-script-checks`},
			patch: func(rt *RuntimeConfig) {
				rt.EnableScriptChecks = true
			},
		},
		{ // todo(fs): shouldn't this be '-encrypt-key'?
			desc:  "-encrypt",
			flags: []string{`-encrypt`, `i0P+gFTkLPg0h53eNYjydg==`},
			patch: func(rt *RuntimeConfig) {
				rt.EncryptKey = "i0P+gFTkLPg0h53eNYjydg=="
			},
		},
		{
			desc:  "-http-port",
			flags: []string{`-http-port`, `123`},
			patch: func(rt *RuntimeConfig) {
				rt.HTTPPort = 123
				rt.HTTPAddrs = []net.Addr{tcpAddr("127.0.0.1:123")}
			},
		},
		{
			desc:  "-join",
			flags: []string{`-join`, `a`, `-join`, `b`},
			patch: func(rt *RuntimeConfig) {
				rt.StartJoinAddrsLAN = []string{"a", "b"}
			},
		},
		{
			desc:  "-join-wan",
			flags: []string{`-join-wan`, `a`, `-join-wan`, `b`},
			patch: func(rt *RuntimeConfig) {
				rt.StartJoinAddrsWAN = []string{"a", "b"}
			},
		},
		{
			desc:  "-log-level",
			flags: []string{`-log-level`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.LogLevel = "a"
			},
		},
		{ // todo(fs): shouldn't this be '-node-name'?
			desc:  "-node",
			flags: []string{`-node`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.NodeName = "a"
			},
		},
		{
			desc:  "-node-id",
			flags: []string{`-node-id`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.NodeID = "a"
			},
		},
		{
			desc:  "-node-meta",
			flags: []string{`-node-meta`, `a:b`, `-node-meta`, `c:d`},
			patch: func(rt *RuntimeConfig) {
				rt.NodeMeta = map[string]string{"a": "b", "c": "d"}
			},
		},
		{
			desc:  "-non-voting-server",
			flags: []string{`-non-voting-server`},
			patch: func(rt *RuntimeConfig) {
				rt.NonVotingServer = true
			},
		},
		{
			desc:  "-pid-file",
			flags: []string{`-pid-file`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.PidFile = "a"
			},
		},
		{
			desc:  "-protocol",
			flags: []string{`-protocol`, `1`},
			patch: func(rt *RuntimeConfig) {
				rt.RPCProtocol = 1
			},
		},
		{
			desc:  "-raft-protocol",
			flags: []string{`-raft-protocol`, `1`},
			patch: func(rt *RuntimeConfig) {
				rt.RaftProtocol = 1
			},
		},
		{
			desc:  "-recursor",
			flags: []string{`-recursor`, `a`, `-recursor`, `b`},
			patch: func(rt *RuntimeConfig) {
				rt.DNSRecursors = []string{"a", "b"}
			},
		},
		{
			desc:  "-rejoin",
			flags: []string{`-rejoin`},
			patch: func(rt *RuntimeConfig) {
				rt.RejoinAfterLeave = true
			},
		},
		{
			desc:  "-retry-interval",
			flags: []string{`-retry-interval`, `5s`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinIntervalLAN = 5 * time.Second
			},
		},
		{
			desc:  "-retry-interval-wan",
			flags: []string{`-retry-interval-wan`, `5s`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinIntervalWAN = 5 * time.Second
			},
		},
		{
			desc:  "-retry-join",
			flags: []string{`-retry-join`, `a`, `-retry-join`, `b`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"a", "b"}
			},
		},
		{
			desc:  "-retry-join-wan",
			flags: []string{`-retry-join-wan`, `a`, `-retry-join-wan`, `b`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinWAN = []string{"a", "b"}
			},
		},
		{
			desc:  "-retry-max",
			flags: []string{`-retry-max`, `1`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinMaxAttemptsLAN = 1
			},
		},
		{
			desc:  "-retry-max-wan",
			flags: []string{`-retry-max-wan`, `1`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinMaxAttemptsWAN = 1
			},
		},
		{
			desc:  "-serf-lan-bind",
			flags: []string{`-serf-lan-bind`, `1.2.3.4`},
			patch: func(rt *RuntimeConfig) {
				rt.SerfBindAddrLAN = tcpAddr("1.2.3.4:8301")
			},
		},
		{
			desc:  "-serf-wan-bind",
			flags: []string{`-serf-wan-bind`, `1.2.3.4`},
			patch: func(rt *RuntimeConfig) {
				rt.SerfBindAddrWAN = tcpAddr("1.2.3.4:8302")
			},
		},
		{
			desc:  "-server",
			flags: []string{`-server`},
			patch: func(rt *RuntimeConfig) {
				rt.ServerMode = true
				rt.LeaveOnTerm = false
				rt.SkipLeaveOnInt = true
			},
		},
		{
			desc:  "-syslog",
			flags: []string{`-syslog`},
			patch: func(rt *RuntimeConfig) {
				rt.EnableSyslog = true
			},
		},
		{
			desc:  "-ui",
			flags: []string{`-ui`},
			patch: func(rt *RuntimeConfig) {
				rt.EnableUI = true
			},
		},
		{
			desc:  "-ui-dir",
			flags: []string{`-ui-dir`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.UIDir = "a"
			},
		},

		// ------------------------------------------------------------
		// deprecated flags
		//

		{
			desc:  "-atlas",
			flags: []string{`-atlas`, `a`},
			warns: []string{`==> DEPRECATION: "-atlas" is deprecated. Please remove it from your configuration`},
		},
		{
			desc:  "-atlas-endpoint",
			flags: []string{`-atlas-endpoint`, `a`},
			warns: []string{`==> DEPRECATION: "-atlas-endpoint" is deprecated. Please remove it from your configuration`},
		},
		{
			desc:  "-atlas-join",
			flags: []string{`-atlas-join`},
			warns: []string{`==> DEPRECATION: "-atlas-join" is deprecated. Please remove it from your configuration`},
		},
		{
			desc:  "-atlas-token",
			flags: []string{`-atlas-token`, `a`},
			warns: []string{`==> DEPRECATION: "-atlas-token" is deprecated. Please remove it from your configuration`},
		},
		{
			desc:  "-dc",
			flags: []string{`-dc`, `a`},
			warns: []string{`==> DEPRECATION: "-dc" is deprecated. Use "-datacenter" instead`},
			patch: func(rt *RuntimeConfig) {
				rt.Datacenter = "a"
			},
		},
		{
			desc:  "-retry-join-azure-tag-name",
			flags: []string{`-retry-join-azure-tag-name`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=azure tag_name=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_azure" is deprecated. Please add "provider=azure tag_name=a" to "retry_join".`},
		},
		{
			desc:  "-retry-join-azure-tag-value",
			flags: []string{`-retry-join-azure-tag-value`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=azure tag_value=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_azure" is deprecated. Please add "provider=azure tag_value=a" to "retry_join".`},
		},
		{
			desc:  "-retry-join-ec2-region",
			flags: []string{`-retry-join-ec2-region`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=aws region=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_ec2" is deprecated. Please add "provider=aws region=a" to "retry_join".`},
		},
		{
			desc:  "-retry-join-ec2-tag-key",
			flags: []string{`-retry-join-ec2-tag-key`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=aws tag_key=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_ec2" is deprecated. Please add "provider=aws tag_key=a" to "retry_join".`},
		},
		{
			desc:  "-retry-join-ec2-tag-value",
			flags: []string{`-retry-join-ec2-tag-value`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=aws tag_value=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_ec2" is deprecated. Please add "provider=aws tag_value=a" to "retry_join".`},
		},
		{
			desc:  "-retry-join-gce-credentials-file",
			flags: []string{`-retry-join-gce-credentials-file`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=gce credentials_file=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_gce" is deprecated. Please add "provider=gce credentials_file=hidden" to "retry_join".`},
		},
		{
			desc:  "-retry-join-gce-project-name",
			flags: []string{`-retry-join-gce-project-name`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=gce project_name=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_gce" is deprecated. Please add "provider=gce project_name=a" to "retry_join".`},
		},
		{
			desc:  "-retry-join-gce-tag-value",
			flags: []string{`-retry-join-gce-tag-value`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=gce tag_value=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_gce" is deprecated. Please add "provider=gce tag_value=a" to "retry_join".`},
		},
		{
			desc:  "-retry-join-gce-zone-pattern",
			flags: []string{`-retry-join-gce-zone-pattern`, `a`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=gce zone_pattern=a"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_gce" is deprecated. Please add "provider=gce zone_pattern=a" to "retry_join".`},
		},

		// ------------------------------------------------------------
		// deprecated fields
		//

		{
			desc:  "addresses.rpc",
			json:  []string{`{"addresses":{ "rpc": "a" }}`},
			hcl:   []string{`addresses = { rpc = "a" }`},
			warns: []string{`==> DEPRECATION: "addresses.rpc" is deprecated and is no longer used. Please remove it from your configuration.`},
		},
		{
			desc:  "ports.rpc",
			json:  []string{`{"ports":{ "rpc": 123 }}`},
			hcl:   []string{`ports = { rpc = 123 }`},
			warns: []string{`==> DEPRECATION: "ports.rpc" is deprecated and is no longer used. Please remove it from your configuration.`},
		},
		{
			desc: "check.service_id alias",
			json: []string{`{"check":{ "service_id":"d", "serviceid":"dd" }}`},
			hcl:  []string{`check = { service_id="d" serviceid="dd" }`},
			patch: func(rt *RuntimeConfig) {
				rt.Checks = []*structs.CheckDefinition{{ServiceID: "dd"}}
			},
			warns: []string{`==> DEPRECATION: "serviceid" is deprecated in check definitions. Please use "service_id" instead.`},
		},
		{
			desc: "check.docker_container_id alias",
			json: []string{`{"check":{ "docker_container_id":"k", "dockercontainerid":"kk" }}`},
			hcl:  []string{`check = { docker_container_id="k" dockercontainerid="kk" }`},
			patch: func(rt *RuntimeConfig) {
				rt.Checks = []*structs.CheckDefinition{{DockerContainerID: "kk"}}
			},
			warns: []string{`==> DEPRECATION: "dockercontainerid" is deprecated in check definitions. Please use "docker_container_id" instead.`},
		},
		{
			desc: "check.tls_skip_verify alias",
			json: []string{`{"check":{ "tls_skip_verify":true, "tlsskipverify":false }}`},
			hcl:  []string{`check = { tls_skip_verify=true tlsskipverify=false }`},
			patch: func(rt *RuntimeConfig) {
				rt.Checks = []*structs.CheckDefinition{{TLSSkipVerify: false}}
			},
			warns: []string{`==> DEPRECATION: "tlsskipverify" is deprecated in check definitions. Please use "tls_skip_verify" instead.`},
		},
		{
			desc: "check.deregister_critical_service_after alias",
			json: []string{`{"check":{ "deregister_critical_service_after":"5s", "deregistercriticalserviceafter": "10s" }}`},
			hcl:  []string{`check = { deregister_critical_service_after="5s" deregistercriticalserviceafter="10s"}`},
			patch: func(rt *RuntimeConfig) {
				rt.Checks = []*structs.CheckDefinition{{DeregisterCriticalServiceAfter: 10 * time.Second}}
			},
			warns: []string{`==> DEPRECATION: "deregistercriticalserviceafter" is deprecated in check definitions. Please use "deregister_critical_service_after" instead.`},
		},
		{
			desc: "http_api_response_headers",
			json: []string{`{"http_api_response_headers":{"a":"b","c":"d"}}`},
			hcl:  []string{`http_api_response_headers = {"a"="b" "c"="d"}`},
			patch: func(rt *RuntimeConfig) {
				rt.HTTPResponseHeaders = map[string]string{"a": "b", "c": "d"}
			},
			warns: []string{`==> DEPRECATION: "http_api_response_headers" is deprecated. Please use "http_config.response_headers" instead.`},
		},
		{
			desc: "retry_join_azure",
			json: []string{`{
					"retry_join_azure":{
						"tag_name": "a",
						"tag_value": "b",
						"subscription_id": "c",
						"tenant_id": "d",
						"client_id": "e",
						"secret_access_key": "f"
					}
				}`},
			hcl: []string{`
					retry_join_azure = {
						tag_name = "a"
						tag_value = "b"
						subscription_id = "c"
						tenant_id = "d"
						client_id = "e"
						secret_access_key = "f"
					}
				`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=azure client_id=e secret_access_key=f subscription_id=c tag_name=a tag_value=b tenant_id=d"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_azure" is deprecated. Please add "provider=azure client_id=hidden secret_access_key=hidden subscription_id=hidden tag_name=a tag_value=b tenant_id=hidden" to "retry_join".`},
		},
		{
			desc: "retry_join_ec2",
			json: []string{`{
					"retry_join_ec2":{
						"tag_key": "a",
						"tag_value": "b",
						"region": "c",
						"access_key_id": "d",
						"secret_access_key": "e"
					}
				}`},
			hcl: []string{`
					retry_join_ec2 = {
						tag_key = "a"
						tag_value = "b"
						region = "c"
						access_key_id = "d"
						secret_access_key = "e"
					}
				`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=aws access_key_id=d region=c secret_access_key=e tag_key=a tag_value=b"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_ec2" is deprecated. Please add "provider=aws access_key_id=hidden region=c secret_access_key=hidden tag_key=a tag_value=b" to "retry_join".`},
		},
		{
			desc: "retry_join_gce",
			json: []string{`{
					"retry_join_gce":{
						"project_name": "a",
						"zone_pattern": "b",
						"tag_value": "c",
						"credentials_file": "d"
					}
				}`},
			hcl: []string{`
					retry_join_gce = {
						project_name = "a"
						zone_pattern = "b"
						tag_value = "c"
						credentials_file = "d"
					}
				`},
			patch: func(rt *RuntimeConfig) {
				rt.RetryJoinLAN = []string{"provider=gce credentials_file=d project_name=a tag_value=c zone_pattern=b"}
			},
			warns: []string{`==> DEPRECATION: "retry_join_gce" is deprecated. Please add "provider=gce credentials_file=hidden project_name=a tag_value=c zone_pattern=b" to "retry_join".`},
		},

		{
			desc:  "telemetry.dogstatsd_addr alias",
			json:  []string{`{"dogstatsd_addr":"a", "telemetry":{"dogstatsd_addr": "b"}}`},
			hcl:   []string{`dogstatsd_addr = "a" telemetry = { dogstatsd_addr = "b"}`},
			warns: []string{`==> DEPRECATION: "dogstatsd_addr" is deprecated. Please use "telemetry.dogstatsd_addr" instead.`},
			patch: func(rt *RuntimeConfig) {
				rt.TelemetryDogstatsdAddr = "a"
			},
		},
		{
			desc:  "telemetry.dogstatsd_tags alias",
			json:  []string{`{"dogstatsd_tags":["a", "b"], "telemetry": { "dogstatsd_tags": ["c", "d"]}}`},
			hcl:   []string{`dogstatsd_tags = ["a", "b"] telemetry = { dogstatsd_tags = ["c", "d"] }`},
			warns: []string{`==> DEPRECATION: "dogstatsd_tags" is deprecated. Please use "telemetry.dogstatsd_tags" instead.`},
			patch: func(rt *RuntimeConfig) {
				rt.TelemetryDogstatsdTags = []string{"a", "b", "c", "d"}
			},
		},
		{
			desc:  "telemetry.statsd_addr alias",
			json:  []string{`{"statsd_addr":"a", "telemetry":{"statsd_addr": "b"}}`},
			hcl:   []string{`statsd_addr = "a" telemetry = { statsd_addr = "b" }`},
			warns: []string{`==> DEPRECATION: "statsd_addr" is deprecated. Please use "telemetry.statsd_addr" instead.`},
			patch: func(rt *RuntimeConfig) {
				rt.TelemetryStatsdAddr = "a"
			},
		},
		{
			desc:  "telemetry.statsite_addr alias",
			json:  []string{`{"statsite_addr":"a", "telemetry":{ "statsite_addr": "b" }}`},
			hcl:   []string{`statsite_addr = "a" telemetry = { statsite_addr = "b"}`},
			warns: []string{`==> DEPRECATION: "statsite_addr" is deprecated. Please use "telemetry.statsite_addr" instead.`},
			patch: func(rt *RuntimeConfig) {
				rt.TelemetryStatsiteAddr = "a"
			},
		},
		{
			desc:  "telemetry.statsite_prefix alias",
			json:  []string{`{"statsite_prefix":"a", "telemetry":{ "statsite_prefix": "b" }}`},
			hcl:   []string{`statsite_prefix = "a" telemetry = { statsite_prefix = "b" }`},
			warns: []string{`==> DEPRECATION: "statsite_prefix" is deprecated. Please use "telemetry.statsite_prefix" instead.`},
			patch: func(rt *RuntimeConfig) {
				rt.TelemetryStatsitePrefix = "a"
			},
		},

		// ------------------------------------------------------------
		// ports and addresses
		//

		{
			desc: "client addr and ports == 0",
			json: []string{`{
					"client_addr":"0.0.0.0",
					"ports":{}
				}`},
			hcl: []string{`
					client_addr = "0.0.0.0"
					ports {}
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("0.0.0.0")}
				rt.DNSAddrs = []net.Addr{tcpAddr("0.0.0.0:8600"), udpAddr("0.0.0.0:8600")}
				rt.HTTPAddrs = []net.Addr{tcpAddr("0.0.0.0:8500")}
			},
		},
		{
			desc: "client addr and ports < 0",
			json: []string{`{
					"client_addr":"0.0.0.0",
					"ports": { "dns":-1, "http":-2, "https":-3 }
				}`},
			hcl: []string{`
					client_addr = "0.0.0.0"
					ports { dns = -1 http = -2 https = -3 }
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("0.0.0.0")}
				rt.DNSPort = -1
				rt.DNSAddrs = nil
				rt.HTTPPort = -1
				rt.HTTPAddrs = nil
			},
		},
		{
			desc: "client addr and ports < 0",
			json: []string{`{
					"client_addr":"0.0.0.0",
					"ports": { "dns":-1, "http":-2, "https":-3 }
				}`},
			hcl: []string{`
					client_addr = "0.0.0.0"
					ports { dns = -1 http = -2 https = -3 }
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("0.0.0.0")}
				rt.DNSPort = -1
				rt.DNSAddrs = nil
				rt.HTTPPort = -1
				rt.HTTPAddrs = nil
			},
		},
		{
			desc: "client addr and ports > 0",
			json: []string{`{
					"client_addr":"0.0.0.0",
					"ports":{ "dns": 1, "http": 2, "https": 3 }
				}`},
			hcl: []string{`
					client_addr = "0.0.0.0"
					ports { dns = 1 http = 2 https = 3 }
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("0.0.0.0")}
				rt.DNSPort = 1
				rt.DNSAddrs = []net.Addr{tcpAddr("0.0.0.0:1"), udpAddr("0.0.0.0:1")}
				rt.HTTPPort = 2
				rt.HTTPAddrs = []net.Addr{tcpAddr("0.0.0.0:2")}
				rt.HTTPSPort = 3
				rt.HTTPSAddrs = []net.Addr{tcpAddr("0.0.0.0:3")}
			},
		},

		{
			desc: "client addr, addresses and ports == 0",
			json: []string{`{
					"client_addr":"0.0.0.0",
					"addresses": { "dns": "1.1.1.1", "http": "2.2.2.2", "https": "3.3.3.3" },
					"ports":{}
				}`},
			hcl: []string{`
					client_addr = "0.0.0.0"
					addresses = { dns = "1.1.1.1" http = "2.2.2.2" https = "3.3.3.3" }
					ports {}
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("0.0.0.0")}
				rt.DNSAddrs = []net.Addr{tcpAddr("1.1.1.1:8600"), udpAddr("1.1.1.1:8600")}
				rt.HTTPAddrs = []net.Addr{tcpAddr("2.2.2.2:8500")}
			},
		},
		{
			desc: "client addr, addresses and ports < 0",
			json: []string{`{
					"client_addr":"0.0.0.0",
					"addresses": { "dns": "1.1.1.1", "http": "2.2.2.2", "https": "3.3.3.3" },
					"ports": { "dns":-1, "http":-2, "https":-3 }
				}`},
			hcl: []string{`
					client_addr = "0.0.0.0"
					addresses = { dns = "1.1.1.1" http = "2.2.2.2" https = "3.3.3.3" }
					ports { dns = -1 http = -2 https = -3 }
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("0.0.0.0")}
				rt.DNSPort = -1
				rt.DNSAddrs = nil
				rt.HTTPPort = -1
				rt.HTTPAddrs = nil
			},
		},
		{
			desc: "client addr, addresses and ports",
			json: []string{`{
					"client_addr": "0.0.0.0",
					"addresses": { "dns": "1.1.1.1", "http": "2.2.2.2", "https": "3.3.3.3" },
					"ports":{ "dns":1, "http":2, "https":3 }
				}`},
			hcl: []string{`
					client_addr = "0.0.0.0"
					addresses = { dns = "1.1.1.1" http = "2.2.2.2" https = "3.3.3.3" }
					ports { dns = 1 http = 2 https = 3 }
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("0.0.0.0")}
				rt.DNSPort = 1
				rt.DNSAddrs = []net.Addr{tcpAddr("1.1.1.1:1"), udpAddr("1.1.1.1:1")}
				rt.HTTPPort = 2
				rt.HTTPAddrs = []net.Addr{tcpAddr("2.2.2.2:2")}
				rt.HTTPSPort = 3
				rt.HTTPSAddrs = []net.Addr{tcpAddr("3.3.3.3:3")}
			},
		},
		{
			desc: "client template and ports",
			json: []string{`{
					"client_addr": "{{ printf \"1.2.3.4 2001:db8::1\" }}",
					"ports":{ "dns":1, "http":2, "https":3 }
				}`},
			hcl: []string{`
					client_addr = "{{ printf \"1.2.3.4 2001:db8::1\" }}"
					ports { dns = 1 http = 2 https = 3 }
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("1.2.3.4"), ipAddr("2001:db8::1")}
				rt.DNSPort = 1
				rt.DNSAddrs = []net.Addr{tcpAddr("1.2.3.4:1"), tcpAddr("[2001:db8::1]:1"), udpAddr("1.2.3.4:1"), udpAddr("[2001:db8::1]:1")}
				rt.HTTPPort = 2
				rt.HTTPAddrs = []net.Addr{tcpAddr("1.2.3.4:2"), tcpAddr("[2001:db8::1]:2")}
				rt.HTTPSPort = 3
				rt.HTTPSAddrs = []net.Addr{tcpAddr("1.2.3.4:3"), tcpAddr("[2001:db8::1]:3")}
			},
		},
		{
			desc: "client, address template and ports",
			json: []string{`{
					"client_addr": "{{ printf \"1.2.3.4 2001:db8::1\" }}",
					"addresses": {
						"dns": "{{ printf \"1.1.1.1 unix://dns 2001:db8::10 \" }}",
						"http": "{{ printf \"2.2.2.2 unix://http 2001:db8::20 \" }}",
						"https": "{{ printf \"3.3.3.3 unix://https 2001:db8::30 \" }}"
					},
					"ports":{ "dns":1, "http":2, "https":3 }
				}`},
			hcl: []string{`
					client_addr = "{{ printf \"1.2.3.4 2001:db8::1\" }}"
					addresses = {
						dns = "{{ printf \"1.1.1.1 unix://dns 2001:db8::10 \" }}"
						http = "{{ printf \"2.2.2.2 unix://http 2001:db8::20 \" }}"
						https = "{{ printf \"3.3.3.3 unix://https 2001:db8::30 \" }}"
					}
					ports { dns = 1 http = 2 https = 3 }
				`},
			patch: func(rt *RuntimeConfig) {
				rt.ClientAddrs = []*net.IPAddr{ipAddr("1.2.3.4"), ipAddr("2001:db8::1")}
				rt.DNSPort = 1
				rt.DNSAddrs = []net.Addr{tcpAddr("1.1.1.1:1"), unixAddr("unix://dns"), tcpAddr("[2001:db8::10]:1"), udpAddr("1.1.1.1:1"), udpAddr("[2001:db8::10]:1")}
				rt.HTTPPort = 2
				rt.HTTPAddrs = []net.Addr{tcpAddr("2.2.2.2:2"), unixAddr("unix://http"), tcpAddr("[2001:db8::20]:2")}
				rt.HTTPSPort = 3
				rt.HTTPSAddrs = []net.Addr{tcpAddr("3.3.3.3:3"), unixAddr("unix://https"), tcpAddr("[2001:db8::30]:3")}
			},
		},
		{
			desc: "advertise address lan template",
			json: []string{`{ "advertise_addr": "{{ printf \"1.2.3.4\" }}" }`},
			hcl:  []string{`advertise_addr = "{{ printf \"1.2.3.4\" }}"`},
			patch: func(rt *RuntimeConfig) {
				rt.AdvertiseAddrLAN = tcpAddr("1.2.3.4:8300")
				rt.AdvertiseAddrWAN = tcpAddr("1.2.3.4:8300")
				rt.TaggedAddresses = map[string]string{
					"lan": "1.2.3.4",
					"wan": "1.2.3.4",
				}
			},
		},
		{
			desc: "advertise address wan template",
			json: []string{`{ "advertise_addr_wan": "{{ printf \"1.2.3.4\" }}" }`},
			hcl:  []string{`advertise_addr_wan = "{{ printf \"1.2.3.4\" }}"`},
			patch: func(rt *RuntimeConfig) {
				rt.AdvertiseAddrWAN = tcpAddr("1.2.3.4:8300")
				rt.TaggedAddresses = map[string]string{
					"lan": "10.0.0.1",
					"wan": "1.2.3.4",
				}
			},
		},
		{
			desc: "serf advertise address lan template",
			json: []string{`{ "advertise_addrs": { "serf_lan": "{{ printf \"1.2.3.4\" }}" } }`},
			hcl:  []string{`advertise_addrs = { serf_lan = "{{ printf \"1.2.3.4\" }}" }`},
			patch: func(rt *RuntimeConfig) {
				rt.SerfAdvertiseAddrLAN = tcpAddr("1.2.3.4:8301")
			},
		},
		{
			desc: "serf advertise address wan template",
			json: []string{`{ "advertise_addrs": { "serf_wan": "{{ printf \"1.2.3.4\" }}" } }`},
			hcl:  []string{`advertise_addrs = { serf_wan = "{{ printf \"1.2.3.4\" }}" }`},
			patch: func(rt *RuntimeConfig) {
				rt.SerfAdvertiseAddrWAN = tcpAddr("1.2.3.4:8302")
			},
		},
		{
			desc: "serf bind address lan template",
			json: []string{`{ "serf_lan": "{{ printf \"1.2.3.4\" }}" }`},
			hcl:  []string{`serf_lan = "{{ printf \"1.2.3.4\" }}"`},
			patch: func(rt *RuntimeConfig) {
				rt.SerfBindAddrLAN = tcpAddr("1.2.3.4:8301")
			},
		},
		{
			desc: "serf bind address wan template",
			json: []string{`{ "serf_wan": "{{ printf \"1.2.3.4\" }}" }`},
			hcl:  []string{`serf_wan = "{{ printf \"1.2.3.4\" }}"`},
			patch: func(rt *RuntimeConfig) {
				rt.SerfBindAddrWAN = tcpAddr("1.2.3.4:8302")
			},
		},

		// ------------------------------------------------------------
		// precedence rules
		//

		{
			desc: "precedence: merge order",
			json: []string{
				`{
						"bootstrap": true,
						"bootstrap_expect": 1,
						"datacenter": "a",
						"start_join": ["a", "b"],
						"node_meta": {"a":"b"}
					}`,
				`{
						"bootstrap": false,
						"bootstrap_expect": 0,
						"datacenter":"b",
						"start_join": ["c", "d"],
						"node_meta": {"c":"d"}
					}`,
			},
			hcl: []string{
				`
					bootstrap = true
					bootstrap_expect = 1
					datacenter = "a"
					start_join = ["a", "b"]
					node_meta = { "a" = "b" }
					`,
				`
					bootstrap = false
					bootstrap_expect = 0
					datacenter = "b"
					start_join = ["c", "d"]
					node_meta = { "c" = "d" }
					`,
			},
			patch: func(rt *RuntimeConfig) {
				rt.Bootstrap = false
				rt.BootstrapExpect = 0
				rt.Datacenter = "b"
				rt.StartJoinAddrsLAN = []string{"a", "b", "c", "d"}
				rt.NodeMeta = map[string]string{"c": "d"}
			},
		},
		{
			desc: "precedence: flag before file",
			json: []string{
				`{
						"advertise_addr": "a",
						"advertise_addr_wan": "a",
						"bootstrap":true,
						"bootstrap_expect": 3,
						"datacenter":"a",
						"node_meta": {"a":"b"},
						"recursors":["a", "b"],
						"serf_lan": "a",
						"serf_wan": "a",
						"start_join":["a", "b"]
					}`,
			},
			hcl: []string{
				`
					advertise_addr = "a"
					advertise_addr_wan = "a"
					bootstrap = true
					bootstrap_expect = 3
					datacenter = "a"
					node_meta = { "a" = "b" }
					recursors = ["a", "b"]
					serf_lan = "a"
					serf_wan = "a"
					start_join = ["a", "b"]
					`,
			},
			flags: []string{
				`-advertise`, `1.1.1.1`,
				`-advertise-wan`, `2.2.2.2`,
				`-bootstrap=false`,
				`-bootstrap-expect=0`,
				`-datacenter=b`,
				`-join`, `c`, `-join`, `d`,
				`-node-meta`, `c:d`,
				`-recursor`, `c`, `-recursor`, `d`,
				`-serf-lan-bind`, `3.3.3.3`,
				`-serf-wan-bind`, `4.4.4.4`,
			},
			patch: func(rt *RuntimeConfig) {
				rt.AdvertiseAddrLAN = tcpAddr("1.1.1.1:8300")
				rt.AdvertiseAddrWAN = tcpAddr("2.2.2.2:8300")
				rt.Datacenter = "b"
				rt.DNSRecursors = []string{"c", "d", "a", "b"}
				rt.NodeMeta = map[string]string{"c": "d"}
				rt.SerfBindAddrLAN = tcpAddr("3.3.3.3:8301")
				rt.SerfBindAddrWAN = tcpAddr("4.4.4.4:8302")
				rt.StartJoinAddrsLAN = []string{"c", "d", "a", "b"}
				rt.TaggedAddresses = map[string]string{
					"lan": "1.1.1.1",
					"wan": "2.2.2.2",
				}
			},
		},

		// ------------------------------------------------------------
		// validations
		//

		{
			desc: "datacenter is lower-cased",
			json: []string{`{ "datacenter": "A" }`},
			hcl:  []string{`datacenter = "A"`},
			patch: func(rt *RuntimeConfig) {
				rt.Datacenter = "a"
			},
		},
		{
			desc: "acl_datacenter is lower-cased",
			json: []string{`{ "acl_datacenter": "A" }`},
			hcl:  []string{`acl_datacenter = "A"`},
			patch: func(rt *RuntimeConfig) {
				rt.ACLDatacenter = "a"
			},
		},
		{
			desc: "acl_replication_token enables acl replication",
			json: []string{`{ "acl_replication_token": "a" }`},
			hcl:  []string{`acl_replication_token = "a"`},
			patch: func(rt *RuntimeConfig) {
				rt.ACLReplicationToken = "a"
				rt.EnableACLReplication = true
			},
		},
		{
			desc:     "ae_interval invalid == 0",
			jsontail: []string{`{ "ae_interval": "0s" }`},
			hcltail:  []string{`ae_interval = "0s"`},
			err:      `ae_interval: must be positive: 0s`,
		},
		{
			desc:     "ae_interval invalid < 0",
			jsontail: []string{`{ "ae_interval": "-1s" }`},
			hcltail:  []string{`ae_interval = "-1s"`},
			err:      `ae_interval: must be positive: -1s`,
		},
		{
			desc: "datacenter invalid",
			json: []string{`{ "datacenter": "%" }`},
			hcl:  []string{`datacenter = "%"`},
			err:  `datacenter: invalid value "%". Please use only [a-z0-9-_]`,
		},
		{
			desc:  "acl_datacenter invalid",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "acl_datacenter": "%" }`},
			hcl:   []string{`acl_datacenter = "%"`},
			err:   `acl_datacenter: invalid value "%". Please use only [a-z0-9-_]`,
		},
		{
			desc:  "autopilot.max_trailing_logs invalid",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "autopilot": { "max_trailing_logs": -1 } }`},
			hcl:   []string{`autopilot = { max_trailing_logs = -1 }`},
			err:   "autopilot.max_trailing_logs: cannot be negative: -1",
		},
		{
			desc:  "bind does not allow socket",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "bind_addr": "unix:///foo" }`},
			hcl:   []string{`bind_addr = "unix:///foo"`},
			err:   "bind_addr: cannot use a unix socket: /foo",
		},
		{
			desc:  "bootstrap without server",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "bootstrap": true }`},
			hcl:   []string{`bootstrap = true`},
			err:   "'bootstrap = true' requires 'server = true'",
		},
		{
			desc:  "bootstrap-expect without server",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "bootstrap_expect": 3 }`},
			hcl:   []string{`bootstrap_expect = 3`},
			err:   "'bootstrap_expect > 0' requires 'server = true'",
		},
		{
			desc:  "bootstrap-expect invalid",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "bootstrap_expect": -1 }`},
			hcl:   []string{`bootstrap_expect = -1`},
			err:   "bootstrap_expect: cannot be negative",
		},
		{
			desc:  "bootstrap-expect and dev mode",
			flags: []string{`-datacenter=a`, `-dev`},
			json:  []string{`{ "bootstrap_expect": 3, "server": true }`},
			hcl:   []string{`bootstrap_expect = 3 server = true`},
			err:   "'bootstrap_expect > 0' not allowed in dev mode",
		},
		{
			desc:  "bootstrap-expect and boostrap",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "bootstrap": true, "bootstrap_expect": 3, "server": true }`},
			hcl:   []string{`bootstrap = true bootstrap_expect = 3 server = true`},
			err:   "'bootstrap_expect > 0' and 'bootstrap = true' are mutually exclusive",
		},
		{
			desc:  "client does not allow socket",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "client_addr": "unix:///foo" }`},
			hcl:   []string{`client_addr = "unix:///foo"`},
			err:   "client_addr: cannot use a unix socket: /foo",
		},
		{
			desc:  "enable_ui and ui_dir",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "enable_ui": true, "ui_dir": "a" }`},
			hcl:   []string{`enable_ui = true ui_dir = "a"`},
			err: "Both the ui and ui-dir flags were specified, please provide only one.\n" +
				"If trying to use your own web UI resources, use the ui-dir flag.\n" +
				"If using Consul version 0.7.0 or later, the web UI is included in the binary so use ui to enable it",
		},

		// test ANY address failures
		// to avoid combinatory explosion for tests use 0.0.0.0, :: or [::] but not all of them
		{
			desc:  "advertise_addr any",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "advertise_addr": "0.0.0.0" }`},
			hcl:   []string{`advertise_addr = "0.0.0.0"`},
			err:   "advertise_addr: cannot be 0.0.0.0, :: or [::]",
		},
		{
			desc:  "advertise_addr_wan any",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "advertise_addr_wan": "::" }`},
			hcl:   []string{`advertise_addr_wan = "::"`},
			err:   "advertise_addr_wan: cannot be 0.0.0.0, :: or [::]",
		},
		{
			desc:  "advertise_addrs.rpc any",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "advertise_addrs":{ "rpc": "[::]" } }`},
			hcl:   []string{`advertise_addrs = { rpc = "[::]" }`},
			err:   "advertise_addrs.rpc: cannot be 0.0.0.0, :: or [::]",
		},
		{
			desc:  "advertise_addrs.serf_lan any",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "advertise_addrs":{ "serf_lan": "[::]" } }`},
			hcl:   []string{`advertise_addrs = { serf_lan = "[::]" }`},
			err:   "advertise_addrs.serf_lan: cannot be 0.0.0.0, :: or [::]",
		},
		{
			desc:  "advertise_addrs.serf_wan any",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "advertise_addrs":{ "serf_wan": "0.0.0.0" } }`},
			hcl:   []string{`advertise_addrs = { serf_wan = "0.0.0.0" }`},
			err:   "advertise_addrs.serf_wan: cannot be 0.0.0.0, :: or [::]",
		},
		{
			desc:  "segments.advertise any",
			flags: []string{`-datacenter=a`, `-server=true`},
			json:  []string{`{ "segments":[{ "name":"x", "advertise": "::", "port": 123 }] }`},
			hcl:   []string{`segments = [{ name = "x" advertise = "::" port = 123 }]`},
			err:   `segments[x].advertise: cannot be 0.0.0.0, :: or [::]`,
		},
		{
			desc:  "segments.advertise socket",
			flags: []string{`-datacenter=a`, `-server=true`},
			json:  []string{`{ "segments":[{ "name":"x", "advertise": "unix:///foo" }] }`},
			hcl:   []string{`segments = [{ name = "x" advertise = "unix:///foo" }]`},
			err:   `segments[x].advertise: cannot use a unix socket: /foo`,
		},
		{
			desc:  "dns_config.udp_answer_limit invalid",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "dns_config": { "udp_answer_limit": 0 } }`},
			hcl:   []string{`dns_config = { udp_answer_limit = 0 }`},
			err:   "dns_config.udp_answer_limit: must be positive: 0",
		},
		{
			desc:  "dns_config.udp_answer_limit invalid",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "dns_config": { "udp_answer_limit": 0 } }`},
			hcl:   []string{`dns_config = { udp_answer_limit = 0 }`},
			err:   "dns_config.udp_answer_limit: must be positive: 0",
		},
		{
			desc:  "performance.raft_multiplier < 0",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "performance": { "raft_multiplier": -1 } }`},
			hcl:   []string{`performance = { raft_multiplier = -1 }`},
			err:   `performance.raft_multiplier: value -1 not between 1 and 10`,
		},
		{
			desc:  "performance.raft_multiplier == 0",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "performance": { "raft_multiplier": 0 } }`},
			hcl:   []string{`performance = { raft_multiplier = 0 }`},
			err:   `performance.raft_multiplier: value 0 not between 1 and 10`,
		},
		{
			desc:  "performance.raft_multiplier > 10",
			flags: []string{`-datacenter=a`},
			json:  []string{`{ "performance": { "raft_multiplier": 20 } }`},
			hcl:   []string{`performance = { raft_multiplier = 20 }`},
			err:   `performance.raft_multiplier: value 20 not between 1 and 10`,
		},
		{
			desc:     "node_name invalid",
			flags:    []string{`-datacenter=a`},
			hostname: func() (string, error) { return "", nil },
			err:      "node_name: cannot be empty",
		},
		{
			desc:  "node_meta key too long",
			flags: []string{`-datacenter=a`},
			json: []string{
				`{ "dns_config": { "udp_answer_limit": 1 } }`,
				`{ "node_meta": { "` + randomString(130) + `": "a" } }`,
			},
			hcl: []string{
				`dns_config = { udp_answer_limit = 1 }`,
				`node_meta = { "` + randomString(130) + `" = "a" }`,
			},
			err: "Key is too long (limit: 128 characters)",
		},
		{
			desc:  "node_meta value too long",
			flags: []string{`-datacenter=a`},
			json: []string{
				`{ "dns_config": { "udp_answer_limit": 1 } }`,
				`{ "node_meta": { "a": "` + randomString(520) + `" } }`,
			},
			hcl: []string{
				`dns_config = { udp_answer_limit = 1 }`,
				`node_meta = { "a" = "` + randomString(520) + `" }`,
			},
			err: "Value is too long (limit: 512 characters)",
		},
		{
			desc:  "node_meta too many keys",
			flags: []string{`-datacenter=a`},
			json: []string{
				`{ "dns_config": { "udp_answer_limit": 1 } }`,
				`{ "node_meta": {` + metaPairs(70, "json") + `} }`,
			},
			hcl: []string{
				`dns_config = { udp_answer_limit = 1 }`,
				`node_meta = {` + metaPairs(70, "hcl") + ` }`,
			},
			err: "Node metadata cannot contain more than 64 key/value pairs",
		},
		{
			desc:  "unique listeners dns vs http",
			flags: []string{`-datacenter=a`},
			json: []string{`{
					"client_addr": "1.2.3.4",
					"ports": { "dns": 1000, "http": 1000 },
					"dns_config": { "udp_answer_limit": 1 }
				}`},
			hcl: []string{`
					client_addr = "1.2.3.4"
					ports = { dns = 1000 http = 1000 }
					dns_config = { udp_answer_limit = 1 }
				`},
			err: "HTTP address 1.2.3.4:1000 already configured for DNS",
		},
		{
			desc:  "unique listeners dns vs https",
			flags: []string{`-datacenter=a`},
			json: []string{`{
					"client_addr": "1.2.3.4",
					"ports": { "dns": 1000, "https": 1000 },
					"dns_config": { "udp_answer_limit": 1 }
				}`},
			hcl: []string{`
					client_addr = "1.2.3.4"
					ports = { dns = 1000 https = 1000 }
					dns_config = { udp_answer_limit = 1 }
				`},
			err: "HTTPS address 1.2.3.4:1000 already configured for DNS",
		},
		{
			desc:  "unique listeners http vs https",
			flags: []string{`-datacenter=a`},
			json: []string{`{
					"client_addr": "1.2.3.4",
					"ports": { "http": 1000, "https": 1000 },
					"dns_config": { "udp_answer_limit": 1 }
				}`},
			hcl: []string{`
					client_addr = "1.2.3.4"
					ports = { http = 1000 https = 1000 }
					dns_config = { udp_answer_limit = 1 }
				`},
			err: "HTTPS address 1.2.3.4:1000 already configured for HTTP",
		},
	}

	for _, tt := range tests {
		for pass, format := range []string{"json", "hcl"} {
			// when we test only flags then there are no JSON or HCL
			// sources and we need to make only one pass over the
			// tests.
			flagsOnly := len(tt.json) == 0 && len(tt.hcl) == 0
			if flagsOnly && pass > 0 {
				continue
			}

			// json and hcl sources need to be in sync
			// to make sure we're generating the same config
			if len(tt.json) != len(tt.hcl) {
				t.Fatal("JSON and HCL test case out of sync")
			}

			// select the source
			srcs, tails := tt.json, tt.jsontail
			if format == "hcl" {
				srcs, tails = tt.hcl, tt.hcltail
			}

			// build the description
			var desc []string
			if !flagsOnly {
				desc = append(desc, format)
			}
			if tt.desc != "" {
				desc = append(desc, tt.desc)
			}

			t.Run(strings.Join(desc, ":"), func(t *testing.T) {
				// add flags
				flags, err := ParseFlags(tt.flags)
				if err != nil {
					t.Fatalf("ParseFlags failed: %s", err)
				}

				// mock the hostname function unless a mock is provided
				hostnameFn := tt.hostname
				if hostnameFn == nil {
					hostnameFn = func() (string, error) { return "nodex", nil }
				}

				// create the builder with the flags
				b := &Builder{
					Flags: flags,
					Head: []Source{
						DefaultSource,
						Source{Name: "data-dir", Format: "hcl", Data: `data_dir = "h8smrzKF"`},
					},
					Tail: []Source{
						NonUserSource,
						VersionSource("abcdef", "0.9.0", "test"),
					},

					Hostname:    hostnameFn,
					PrivateIPv4: func() (*net.IPAddr, error) { return ipAddr("10.0.0.1"), nil },
					PublicIPv6:  func() (*net.IPAddr, error) { return ipAddr("2001:db8:1000"), nil },
				}

				// read the source fragements
				for i, data := range srcs {
					b.Sources = append(b.Sources, Source{
						Name:   fmt.Sprintf("%s-%d", format, i),
						Format: format,
						Data:   data,
					})
				}
				for i, data := range tails {
					b.Tail = append(b.Tail, Source{
						Name:   fmt.Sprintf("%s-%d", format, i),
						Format: format,
						Data:   data,
					})
				}

				// build/merge the config fragments
				rt, err := b.BuildAndValidate()
				if err == nil && tt.err != "" {
					t.Fatalf("got no error want %q", tt.err)
				}
				if err != nil && tt.err != "" && !strings.Contains(err.Error(), tt.err) {

					t.Fatalf("error %q does not contain %q", err.Error(), tt.err)
				}

				// check the warnings
				if !verify.Values(t, "warnings", b.Warnings, tt.warns) {
					t.FailNow()
				}

				// stop if we expected an error
				if tt.err != "" {
					return
				}

				// build the runtime config we expect
				// from the same default config and patch it
				x := &Builder{
					Head:        b.Head,
					Tail:        b.Tail,
					Hostname:    b.Hostname,
					PrivateIPv4: b.PrivateIPv4,
					PublicIPv6:  b.PublicIPv6,
				}
				wantRT, err := x.Build()
				if err != nil {
					t.Fatalf("build default failed: %s", err)
				}
				if tt.patch != nil {
					tt.patch(&wantRT)
				}
				if err := x.Validate(wantRT); err != nil {
					t.Fatalf("validate default failed: %s", err)
				}
				if got, want := rt, wantRT; !verify.Values(t, "", got, want) {
					t.FailNow()
				}
			})
		}
	}
}

// TestFullConfig tests the conversion from a fully populated JSON or
// HCL config file to a RuntimeConfig structure. All fields must be set
// to a unique non-zero value.
//
// To aid populating the fields the following bash functions can be used
// to generate random strings and ints:
//
//   random-int() { echo $RANDOM }
//   random-string() { base64 /dev/urandom | tr -d '/+' | fold -w ${1:-32} | head -n 1 }
//
// To generate a random string of length 8 run the following command in
// a terminal:
//
//   random-string 8
//
func TestFullConfig(t *testing.T) {
	flagSrc := []string{`-dev`}
	src := map[string]string{
		"json": `{
			"acl_agent_master_token": "furuQD0b",
			"acl_agent_token": "cOshLOQ2",
			"acl_datacenter": "m3urck3z",
			"acl_default_policy": "ArK3WIfE",
			"acl_down_policy": "vZXMfMP0",
			"acl_enforce_version_8": true,
			"acl_master_token": "C1Q1oIwh",
			"acl_replication_token": "LMmgy5dO",
			"acl_token": "O1El0wan",
			"acl_ttl": "18060s",
			"addresses": {
				"dns": "93.95.95.81",
				"http": "83.39.91.39",
				"https": "95.17.17.19",
				"rpc": "ZIkSmEPN"
			},
			"advertise_addr": "17.99.29.16",
			"advertise_addr_wan": "78.63.37.19",
			"advertise_addrs": {
				"rpc": "28.27.94.38",
				"serf_lan": "49.38.36.95",
				"serf_wan": "63.38.52.13"
			},
			"autopilot": {
				"cleanup_dead_servers": true,
				"disable_upgrade_migration": true,
				"last_contact_threshold": "12705s",
				"max_trailing_logs": 17849,
				"redundancy_zone_tag": "3IsufDJf",
				"server_stabilization_time": "23057s",
				"upgrade_version_tag": "W9pDwFAL"
			},
			"bind_addr": "16.99.34.17",
			"bootstrap": true,
			"bootstrap_expect": 53,
			"ca_file": "erA7T0PM",
			"ca_path": "mQEN1Mfp",
			"cert_file": "7s4QAzDk",
			"check": {
				"id": "fZaCAXww",
				"name": "OOM2eo0f",
				"notes": "zXzXI9Gt",
				"service_id": "L8G0QNmR",
				"token": "oo4BCTgJ",
				"status": "qLykAl5u",
				"script": "dhGfIF8n",
				"http": "29B93haH",
				"header": {
					"hBq0zn1q": [ "2a9o9ZKP", "vKwA5lR6" ],
					"f3r6xFtM": [ "RyuIdDWv", "QbxEcIUM" ]
				},
				"method": "Dou0nGT5",
				"tcp": "JY6fTTcw",
				"interval": "18714s",
				"docker_container_id": "qF66POS9",
				"shell": "sOnDy228",
				"tls_skip_verify": true,
				"timeout": "5954s",
				"ttl": "30044s",
				"deregister_critical_service_after": "13209s"
			},
			"checks": [
				{
					"id": "uAjE6m9Z",
					"name": "QsZRGpYr",
					"notes": "VJ7Sk4BY",
					"service_id": "lSulPcyz",
					"token": "toO59sh8",
					"status": "9RlWsXMV",
					"script": "8qbd8tWw",
					"http": "dohLcyQ2",
					"header": {
						"ZBfTin3L": [ "1sDbEqYG", "lJGASsWK" ],
						"Ui0nU99X": [ "LMccm3Qe", "k5H5RggQ" ]
					},
					"method": "aldrIQ4l",
					"tcp": "RJQND605",
					"interval": "22164s",
					"docker_container_id": "ipgdFtjd",
					"shell": "qAeOYy0M",
					"tls_skip_verify": true,
					"timeout": "1813s",
					"ttl": "21743s",
					"deregister_critical_service_after": "14232s"
				},
				{
					"id": "Cqq95BhP",
					"name": "3qXpkS0i",
					"notes": "sb5qLTex",
					"service_id": "CmUUcRna",
					"token": "a3nQzHuy",
					"status": "irj26nf3",
					"script": "FJsI1oXt",
					"http": "yzhgsQ7Y",
					"header": {
						"zcqwA8dO": [ "qb1zx0DL", "sXCxPFsD" ],
						"qxvdnSE9": [ "6wBPUYdF", "YYh8wtSZ" ]
					},
					"method": "gLrztrNw",
					"tcp": "4jG5casb",
					"interval": "28767s",
					"docker_container_id": "THW6u7rL",
					"shell": "C1Zt3Zwh",
					"tls_skip_verify": true,
					"timeout": "18506s",
					"ttl": "31006s",
					"deregister_critical_service_after": "2366s"
				}
			],
			"check_update_interval": "16507s",
			"client_addr": "93.83.18.19",
			"data_dir": "oTOOIoV9",
			"datacenter": "rzo029wg",
			"disable_anonymous_signature": true,
			"disable_coordinates": true,
			"disable_host_node_id": true,
			"disable_keyring_file": true,
			"disable_remote_exec": true,
			"disable_update_check": true,
			"domain": "7W1xXSqd",
			"dns_config": {
				"allow_stale": true,
				"disable_compression": true,
				"enable_truncate": true,
				"max_stale": "29685s",
				"node_ttl": "7084s",
				"only_passing": true,
				"recursor_timeout": "4427s",
				"service_ttl": {
					"*": "32030s"
				},
				"udp_answer_limit": 29909
			},
			"enable_acl_replication": true,
			"enable_debug": true,
			"enable_script_checks": true,
			"enable_syslog": true,
			"enable_ui": true,
			"encrypt": "A4wELWqH",
			"encrypt_verify_incoming": true,
			"encrypt_verify_outgoing": true,
			"http_config": {
				"block_endpoints": [ "RBvAFcGD", "fWOWFznh" ],
				"response_headers": {
					"M6TKa9NP": "xjuxjOzQ",
					"JRCrHZed": "rl0mTx81"
				}
			},
			"key_file": "IEkkwgIA",
			"leave_on_terminate": true,
			"limits": {
				"rpc_rate": 12029.43,
				"rpc_max_burst": 44848
			},
			"log_level": "k1zo9Spt",
			"node_id": "AsUIlw99",
			"node_meta": {
				"5mgGQMBk": "mJLtVMSG",
				"A7ynFMJB": "0Nx6RGab"
			},
			"node_name": "otlLxGaI",
			"non_voting_server": true,
			"performance": {
				"raft_multiplier": 5
			},
			"pid_file": "43xN80Km",
			"ports": {
				"dns": 7001,
				"http": 7999,
				"https": 15127,
				"rpc": 10664,
				"server": 3757
			},
			"protocol": 30793,
			"raft_protocol": 19016,
			"reconnect_timeout": "23739s",
			"reconnect_timeout_wan": "26694s",
			"recursor": "EZX7MOYF",
			"recursors": [ "FtFhoUHl", "UYkwck1k" ],
			"rejoin_after_leave": true,
			"retry_interval": "8067s",
			"retry_interval_wan": "28866s",
			"retry_join": [ "pbsSFY7U", "l0qLtWij" ],
			"retry_join_wan": [ "PFsR02Ye", "rJdQIhER" ],
			"retry_max": 913,
			"retry_max_wan": 23160,
			"segment": "BC2NhTDi",
			"segments": [
				{
					"name": "PExYMe2E",
					"bind": "36.73.36.19",
					"port": 38295,
					"rpc_listener": true,
					"advertise": "63.39.19.18"
				},
				{
					"name": "UzCvJgup",
					"bind": "37.58.38.19",
					"port": 39292,
					"rpc_listener": true,
					"advertise": "83.58.26.27"
				}
			],
			"serf_lan": "99.43.63.15",
			"serf_wan": "67.88.33.19",
			"server": true,
			"server_name": "Oerr9n1G",
			"service": {
				"id": "dLOXpSCI",
				"name": "o1ynPkp0",
				"tags": ["nkwshvM5", "NTDWn3ek"],
				"address": "cOlSOhbp",
				"token": "msy7iWER",
				"port": 24237,
				"enable_tag_override": true,
				"check": {
					"check_id": "RMi85Dv8",
					"name": "iehanzuq",
					"status": "rCvn53TH",
					"notes": "fti5lfF3",
					"script": "rtj34nfd",
					"http": "dl3Fgme3",
					"header": {
						"rjm4DEd3": ["2m3m2Fls"],
						"l4HwQ112": ["fk56MNlo", "dhLK56aZ"]
					},
					"method": "9afLm3Mj",
					"tcp": "fjiLFqVd",
					"interval": "23926s",
					"docker_container_id": "dO5TtRHk",
					"shell": "e6q2ttES",
					"tls_skip_verify": true,
					"timeout": "38483s",
					"ttl": "10943s",
					"deregister_critical_service_after": "68787s"
				},
				"checks": [
					{
						"id": "Zv99e9Ka",
						"name": "sgV4F7Pk",
						"notes": "yP5nKbW0",
						"status": "7oLMEyfu",
						"script": "NlUQ3nTE",
						"http": "KyDjGY9H",
						"header": {
							"gv5qefTz": [ "5Olo2pMG", "PvvKWQU5" ],
							"SHOVq1Vv": [ "jntFhyym", "GYJh32pp" ]
						},
						"method": "T66MFBfR",
						"tcp": "bNnNfx2A",
						"interval": "22224s",
						"docker_container_id": "ipgdFtjd",
						"shell": "omVZq7Sz",
						"tls_skip_verify": true,
						"timeout": "18913s",
						"ttl": "44743s",
						"deregister_critical_service_after": "8482s"
					},
					{
						"id": "G79O6Mpr",
						"name": "IEqrzrsd",
						"notes": "SVqApqeM",
						"status": "XXkVoZXt",
						"script": "IXLZTM6E",
						"http": "kyICZsn8",
						"header": {
							"4ebP5vL4": [ "G20SrL5Q", "DwPKlMbo" ],
							"p2UI34Qz": [ "UsG1D0Qh", "NHhRiB6s" ]
						},
						"method": "ciYHWors",
						"tcp": "FfvCwlqH",
						"interval": "12356s",
						"docker_container_id": "HBndBU6R",
						"shell": "hVI33JjA",
						"tls_skip_verify": true,
						"timeout": "38282s",
						"ttl": "1181s",
						"deregister_critical_service_after": "4992s"
					}
				]
			},
			"services": [
				{
					"id": "wI1dzxS4",
					"name": "7IszXMQ1",
					"tags": ["0Zwg8l6v", "zebELdN5"],
					"address": "9RhqPSPB",
					"token": "myjKJkWH",
					"port": 72219,
					"enable_tag_override": true,
					"check": {
						"check_id": "qmfeO5if",
						"name": "atDGP7n5",
						"status": "pDQKEhWL",
						"notes": "Yt8EDLev",
						"script": "MDu7wjlD",
						"http": "qzHYvmJO",
						"header": {
							"UkpmZ3a3": ["2dfzXuxZ"],
							"cVFpko4u": ["gGqdEB6k", "9LsRo22u"]
						},
						"method": "X5DrovFc",
						"tcp": "ICbxkpSF",
						"interval": "24392s",
						"docker_container_id": "ZKXr68Yb",
						"shell": "CEfzx0Fo",
						"tls_skip_verify": true,
						"timeout": "38333s",
						"ttl": "57201s",
						"deregister_critical_service_after": "44214s"
					}
				},
				{
					"id": "MRHVMZuD",
					"name": "6L6BVfgH",
					"tags": ["7Ale4y6o", "PMBW08hy"],
					"address": "R6H6g8h0",
					"token": "ZgY8gjMI",
					"port": 38292,
					"enable_tag_override": true,
					"checks": [
						{
							"id": "GTti9hCo",
							"name": "9OOS93ne",
							"notes": "CQy86DH0",
							"status": "P0SWDvrk",
							"script": "6BhLJ7R9",
							"http": "u97ByEiW",
							"header": {
								"MUlReo8L": [ "AUZG7wHG", "gsN0Dc2N" ],
								"1UJXjVrT": [ "OJgxzTfk", "xZZrFsq7" ]
							},
							"method": "5wkAxCUE",
							"tcp": "MN3oA9D2",
							"interval": "32718s",
							"docker_container_id": "cU15LMet",
							"shell": "nEz9qz2l",
							"tls_skip_verify": true,
							"timeout": "34738s",
							"ttl": "22773s",
							"deregister_critical_service_after": "84282s"
						},
						{
							"id": "UHsDeLxG",
							"name": "PQSaPWlT",
							"notes": "jKChDOdl",
							"status": "5qFz6OZn",
							"script": "PbdxFZ3K",
							"http": "1LBDJhw4",
							"header": {
								"cXPmnv1M": [ "imDqfaBx", "NFxZ1bQe" ],
								"vr7wY7CS": [ "EtCoNPPL", "9vAarJ5s" ]
							},
							"method": "wzByP903",
							"tcp": "2exjZIGE",
							"interval": "5656s",
							"docker_container_id": "5tDBWpfA",
							"shell": "rlTpLM8s",
							"tls_skip_verify": true,
							"timeout": "4868s",
							"ttl": "11222s",
							"deregister_critical_service_after": "68482s"
						}
					]
				}
			],
			"session_ttl_min": "26627s",
			"skip_leave_on_interrupt": true,
			"start_join": [ "LR3hGDoG", "MwVpZ4Up" ],
			"start_join_wan": [ "EbFSc3nA", "kwXTh623" ],
			"syslog_facility": "hHv79Uia",
			"tagged_addresses": {
				"7MYgHrYH": "dALJAhLD",
				"h6DdBy6K": "ebrr9zZ8"
			},
			"telemetry": {
				"circonus_api_app": "p4QOTe9j",
				"circonus_api_token": "E3j35V23",
				"circonus_api_url": "mEMjHpGg",
				"circonus_broker_id": "BHlxUhed",
				"circonus_broker_select_tag": "13xy1gHm",
				"circonus_check_display_name": "DRSlQR6n",
				"circonus_check_force_metric_activation": "Ua5FGVYf",
				"circonus_check_id": "kGorutad",
				"circonus_check_instance_id": "rwoOL6R4",
				"circonus_check_search_tag": "ovT4hT4f",
				"circonus_check_tags": "prvO4uBl",
				"circonus_submission_interval": "DolzaflP",
				"circonus_submission_url": "gTcbS93G",
				"disable_hostname": true,
				"dogstatsd_addr": "0wSndumK",
				"dogstatsd_tags": [ "3N81zSUB","Xtj8AnXZ" ],
				"filter_default": true,
				"prefix_filter": [ "+oJotS8XJ","-cazlEhGn" ],
				"statsd_address": "drce87cy",
				"statsite_address": "HpFwKB8R",
				"statsite_prefix": "ftO6DySn"
			},
			"tls_cipher_suites": "TLS_RSA_WITH_RC4_128_SHA,TLS_RSA_WITH_3DES_EDE_CBC_SHA",
			"tls_min_version": "pAOWafkR",
			"tls_prefer_server_cipher_suites": true,
			"translate_wan_addrs": true,
			"ui_dir": "11IFzAUn",
			"unix_sockets": {
				"group": "8pFodrV8",
				"mode": "E8sAwOv4",
				"user": "E0nB1DwA"
			},
			"verify_incoming": true,
			"verify_incoming_https": true,
			"verify_incoming_rpc": true,
			"verify_outgoing": true,
			"verify_server_hostname": true,
			"watches": [
				{
					"type": "key",
					"datacenter": "GyE6jpeW",
					"key": "j9lF1Tve",
					"handler": "90N7S4LN"
				}
			]
		}`,
		"hcl": `
			acl_agent_master_token = "furuQD0b"
			acl_agent_token = "cOshLOQ2"
			acl_datacenter = "m3urck3z"
			acl_default_policy = "ArK3WIfE"
			acl_down_policy = "vZXMfMP0"
			acl_enforce_version_8 = true
			acl_master_token = "C1Q1oIwh"
			acl_replication_token = "LMmgy5dO"
			acl_token = "O1El0wan"
			acl_ttl = "18060s"
			addresses = {
				dns = "93.95.95.81"
				http = "83.39.91.39"
				https = "95.17.17.19"
				rpc = "ZIkSmEPN"
			}
			advertise_addr = "17.99.29.16"
			advertise_addr_wan = "78.63.37.19"
			advertise_addrs = {
				rpc = "28.27.94.38"
				serf_lan = "49.38.36.95"
				serf_wan = "63.38.52.13"
			}
			autopilot = {
				cleanup_dead_servers = true
				disable_upgrade_migration = true
				last_contact_threshold = "12705s"
				max_trailing_logs = 17849
				redundancy_zone_tag = "3IsufDJf"
				server_stabilization_time = "23057s"
				upgrade_version_tag = "W9pDwFAL"
			}
			bind_addr = "16.99.34.17"
			bootstrap = true
			bootstrap_expect = 53
			ca_file = "erA7T0PM"
			ca_path = "mQEN1Mfp"
			cert_file = "7s4QAzDk"
			check = {
				id = "fZaCAXww"
				name = "OOM2eo0f"
				notes = "zXzXI9Gt"
				service_id = "L8G0QNmR"
				token = "oo4BCTgJ"
				status = "qLykAl5u"
				script = "dhGfIF8n"
				http = "29B93haH"
				header = {
					hBq0zn1q = [ "2a9o9ZKP", "vKwA5lR6" ]
					f3r6xFtM = [ "RyuIdDWv", "QbxEcIUM" ]
				}
				method = "Dou0nGT5"
				tcp = "JY6fTTcw"
				interval = "18714s"
				docker_container_id = "qF66POS9"
				shell = "sOnDy228"
				tls_skip_verify = true
				timeout = "5954s"
				ttl = "30044s"
				deregister_critical_service_after = "13209s"
			},
			checks = [
				{
					id = "uAjE6m9Z"
					name = "QsZRGpYr"
					notes = "VJ7Sk4BY"
					service_id = "lSulPcyz"
					token = "toO59sh8"
					status = "9RlWsXMV"
					script = "8qbd8tWw"
					http = "dohLcyQ2"
					header = {
						"ZBfTin3L" = [ "1sDbEqYG", "lJGASsWK" ]
						"Ui0nU99X" = [ "LMccm3Qe", "k5H5RggQ" ]
					}
					method = "aldrIQ4l"
					tcp = "RJQND605"
					interval = "22164s"
					docker_container_id = "ipgdFtjd"
					shell = "qAeOYy0M"
					tls_skip_verify = true
					timeout = "1813s"
					ttl = "21743s"
					deregister_critical_service_after = "14232s"
				},
				{
					id = "Cqq95BhP"
					name = "3qXpkS0i"
					notes = "sb5qLTex"
					service_id = "CmUUcRna"
					token = "a3nQzHuy"
					status = "irj26nf3"
					script = "FJsI1oXt"
					http = "yzhgsQ7Y"
					header = {
						"zcqwA8dO" = [ "qb1zx0DL", "sXCxPFsD" ]
						"qxvdnSE9" = [ "6wBPUYdF", "YYh8wtSZ" ]
					}
					method = "gLrztrNw"
					tcp = "4jG5casb"
					interval = "28767s"
					docker_container_id = "THW6u7rL"
					shell = "C1Zt3Zwh"
					tls_skip_verify = true
					timeout = "18506s"
					ttl = "31006s"
					deregister_critical_service_after = "2366s"
				}
			]
			check_update_interval = "16507s"
			client_addr = "93.83.18.19"
			data_dir = "oTOOIoV9"
			datacenter = "rzo029wg"
			disable_anonymous_signature = true
			disable_coordinates = true
			disable_host_node_id = true
			disable_keyring_file = true
			disable_remote_exec = true
			disable_update_check = true
			domain = "7W1xXSqd"
			dns_config {
				allow_stale = true
				disable_compression = true
				enable_truncate = true
				max_stale = "29685s"
				node_ttl = "7084s"
				only_passing = true
				recursor_timeout = "4427s"
				service_ttl = {
					"*" = "32030s"
				}
				udp_answer_limit = 29909
			}
			enable_acl_replication = true
			enable_debug = true
			enable_script_checks = true
			enable_syslog = true
			enable_ui = true
			encrypt = "A4wELWqH"
			encrypt_verify_incoming = true
			encrypt_verify_outgoing = true
			http_config {
				block_endpoints = [ "RBvAFcGD", "fWOWFznh" ]
				response_headers = {
					"M6TKa9NP" = "xjuxjOzQ"
					"JRCrHZed" = "rl0mTx81"
				}
			}
			key_file = "IEkkwgIA"
			leave_on_terminate = true
			limits {
				rpc_rate = 12029.43
				rpc_max_burst = 44848
			}
			log_level = "k1zo9Spt"
			node_id = "AsUIlw99"
			node_meta {
				"5mgGQMBk" = "mJLtVMSG"
				"A7ynFMJB" = "0Nx6RGab"
			}
			node_name = "otlLxGaI"
			non_voting_server = true
			performance {
				raft_multiplier = 5
			}
			pid_file = "43xN80Km"
			ports {
				dns = 7001,
				http = 7999,
				https = 15127
				rpc = 10664
				server = 3757
			}
			protocol = 30793
			raft_protocol = 19016
			reconnect_timeout = "23739s"
			reconnect_timeout_wan = "26694s"
			recursor = "EZX7MOYF"
			recursors = [ "FtFhoUHl", "UYkwck1k" ]
			rejoin_after_leave = true
			retry_interval = "8067s"
			retry_interval_wan = "28866s"
			retry_join = [ "pbsSFY7U", "l0qLtWij" ]
			retry_join_wan = [ "PFsR02Ye", "rJdQIhER" ]
			retry_max = 913
			retry_max_wan = 23160
			segment = "BC2NhTDi"
			segments = [
				{
					name = "PExYMe2E"
					bind = "36.73.36.19"
					port = 38295
					rpc_listener = true
					advertise = "63.39.19.18"
				},
				{
					name = "UzCvJgup"
					bind = "37.58.38.19"
					port = 39292
					rpc_listener = true
					advertise = "83.58.26.27"
				}
			]
			serf_lan = "99.43.63.15"
			serf_wan = "67.88.33.19"
			server = true
			server_name = "Oerr9n1G"
			service = {
				id = "dLOXpSCI"
				name = "o1ynPkp0"
				tags = ["nkwshvM5", "NTDWn3ek"]
				address = "cOlSOhbp"
				token = "msy7iWER"
				port = 24237
				enable_tag_override = true
				check = {
					check_id = "RMi85Dv8"
					name = "iehanzuq"
					status = "rCvn53TH"
					notes = "fti5lfF3"
					script = "rtj34nfd"
					http = "dl3Fgme3"
					header = {
						rjm4DEd3 = [ "2m3m2Fls" ]
						l4HwQ112 = [ "fk56MNlo", "dhLK56aZ" ]
					}
					method = "9afLm3Mj"
					tcp = "fjiLFqVd"
					interval = "23926s"
					docker_container_id = "dO5TtRHk"
					shell = "e6q2ttES"
					tls_skip_verify = true
					timeout = "38483s"
					ttl = "10943s"
					deregister_critical_service_after = "68787s"
				}
				checks = [
					{
						id = "Zv99e9Ka"
						name = "sgV4F7Pk"
						notes = "yP5nKbW0"
						status = "7oLMEyfu"
						script = "NlUQ3nTE"
						http = "KyDjGY9H"
						header = {
							"gv5qefTz" = [ "5Olo2pMG", "PvvKWQU5" ]
							"SHOVq1Vv" = [ "jntFhyym", "GYJh32pp" ]
						}
						method = "T66MFBfR"
						tcp = "bNnNfx2A"
						interval = "22224s"
						docker_container_id = "ipgdFtjd"
						shell = "omVZq7Sz"
						tls_skip_verify = true
						timeout = "18913s"
						ttl = "44743s"
						deregister_critical_service_after = "8482s"
					},
					{
						id = "G79O6Mpr"
						name = "IEqrzrsd"
						notes = "SVqApqeM"
						status = "XXkVoZXt"
						script = "IXLZTM6E"
						http = "kyICZsn8"
						header = {
							"4ebP5vL4" = [ "G20SrL5Q", "DwPKlMbo" ]
							"p2UI34Qz" = [ "UsG1D0Qh", "NHhRiB6s" ]
						}
						method = "ciYHWors"
						tcp = "FfvCwlqH"
						interval = "12356s"
						docker_container_id = "HBndBU6R"
						shell = "hVI33JjA"
						tls_skip_verify = true
						timeout = "38282s"
						ttl = "1181s"
						deregister_critical_service_after = "4992s"
					}
				]
			}
			services = [
				{
					id = "wI1dzxS4"
					name = "7IszXMQ1"
					tags = ["0Zwg8l6v", "zebELdN5"]
					address = "9RhqPSPB"
					token = "myjKJkWH"
					port = 72219
					enable_tag_override = true
					check = {
						check_id = "qmfeO5if"
						name = "atDGP7n5"
						status = "pDQKEhWL"
						notes = "Yt8EDLev"
						script = "MDu7wjlD"
						http = "qzHYvmJO"
						header = {
							UkpmZ3a3 = [ "2dfzXuxZ" ]
							cVFpko4u = [ "gGqdEB6k", "9LsRo22u" ]
						}
						method = "X5DrovFc"
						tcp = "ICbxkpSF"
						interval = "24392s"
						docker_container_id = "ZKXr68Yb"
						shell = "CEfzx0Fo"
						tls_skip_verify = true
						timeout = "38333s"
						ttl = "57201s"
						deregister_critical_service_after = "44214s"
					}
				},
				{
					id = "MRHVMZuD"
					name = "6L6BVfgH"
					tags = ["7Ale4y6o", "PMBW08hy"]
					address = "R6H6g8h0"
					token = "ZgY8gjMI"
					port = 38292
					enable_tag_override = true
					checks = [
						{
							id = "GTti9hCo"
							name = "9OOS93ne"
							notes = "CQy86DH0"
							status = "P0SWDvrk"
							script = "6BhLJ7R9"
							http = "u97ByEiW"
							header = {
								"MUlReo8L" = [ "AUZG7wHG", "gsN0Dc2N" ]
								"1UJXjVrT" = [ "OJgxzTfk", "xZZrFsq7" ]
							}
							method = "5wkAxCUE"
							tcp = "MN3oA9D2"
							interval = "32718s"
							docker_container_id = "cU15LMet"
							shell = "nEz9qz2l"
							tls_skip_verify = true
							timeout = "34738s"
							ttl = "22773s"
							deregister_critical_service_after = "84282s"
						},
						{
							id = "UHsDeLxG"
							name = "PQSaPWlT"
							notes = "jKChDOdl"
							status = "5qFz6OZn"
							script = "PbdxFZ3K"
							http = "1LBDJhw4"
							header = {
								"cXPmnv1M" = [ "imDqfaBx", "NFxZ1bQe" ],
								"vr7wY7CS" = [ "EtCoNPPL", "9vAarJ5s" ]
							}
							method = "wzByP903"
							tcp = "2exjZIGE"
							interval = "5656s"
							docker_container_id = "5tDBWpfA"
							shell = "rlTpLM8s"
							tls_skip_verify = true
							timeout = "4868s"
							ttl = "11222s"
							deregister_critical_service_after = "68482s"
						}
					]
				}
			]
			session_ttl_min = "26627s"
			skip_leave_on_interrupt = true
			start_join = [ "LR3hGDoG", "MwVpZ4Up" ]
			start_join_wan = [ "EbFSc3nA", "kwXTh623" ]
			syslog_facility = "hHv79Uia"
			tagged_addresses = {
				"7MYgHrYH" = "dALJAhLD"
				"h6DdBy6K" = "ebrr9zZ8"
			}
			telemetry {
				circonus_api_app = "p4QOTe9j"
				circonus_api_token = "E3j35V23"
				circonus_api_url = "mEMjHpGg"
				circonus_broker_id = "BHlxUhed"
				circonus_broker_select_tag = "13xy1gHm"
				circonus_check_display_name = "DRSlQR6n"
				circonus_check_force_metric_activation = "Ua5FGVYf"
				circonus_check_id = "kGorutad"
				circonus_check_instance_id = "rwoOL6R4"
				circonus_check_search_tag = "ovT4hT4f"
				circonus_check_tags = "prvO4uBl"
				circonus_submission_interval = "DolzaflP"
				circonus_submission_url = "gTcbS93G"
				disable_hostname = true
				dogstatsd_addr = "0wSndumK"
				dogstatsd_tags = [ "3N81zSUB","Xtj8AnXZ" ]
				filter_default = true
				prefix_filter = [ "+oJotS8XJ","-cazlEhGn" ]
				statsd_address = "drce87cy"
				statsite_address = "HpFwKB8R"
				statsite_prefix = "ftO6DySn"
			}
			tls_cipher_suites = "TLS_RSA_WITH_RC4_128_SHA,TLS_RSA_WITH_3DES_EDE_CBC_SHA"
			tls_min_version = "pAOWafkR"
			tls_prefer_server_cipher_suites = true
			translate_wan_addrs = true
			ui_dir = "11IFzAUn"
			unix_sockets = {
				group = "8pFodrV8"
				mode = "E8sAwOv4"
				user = "E0nB1DwA"
			}
			verify_incoming = true
			verify_incoming_https = true
			verify_incoming_rpc = true
			verify_outgoing = true
			verify_server_hostname = true
			watches = [{
				type = "key"
				datacenter = "GyE6jpeW"
				key = "j9lF1Tve"
				handler = "90N7S4LN"
			}]
		`}

	nonUserSource := map[string]Source{
		"json": Source{
			Name:   "non-user.json",
			Format: "json",
			Data: `{
				"acl_disabled_ttl": "957s",
				"check_deregister_interval_min": "27870s",
				"check_reap_interval": "10662s",
				"ae_interval": "10003s",
				"sync_coordinate_rate_target": 137.81,
				"sync_coordinate_interval_min": "27983s"
			}`,
		},
		"hcl": Source{
			Name:   "non-user.hcl",
			Format: "hcl",
			Data: `
				acl_disabled_ttl = "957s"
				check_deregister_interval_min = "27870s"
				check_reap_interval = "10662s"
				ae_interval = "10003s"
				sync_coordinate_rate_target = 137.81
				sync_coordinate_interval_min = "27983s"
		`,
		},
	}

	want := RuntimeConfig{
		// non-user configurable values
		ACLDisabledTTL:             957 * time.Second,
		CheckDeregisterIntervalMin: 27870 * time.Second,
		CheckReapInterval:          10662 * time.Second,
		AEInterval:                 10003 * time.Second,
		SyncCoordinateRateTarget:   137.81,
		SyncCoordinateIntervalMin:  27983 * time.Second,

		Revision:          "JNtPSav3",
		Version:           "R909Hblt",
		VersionPrerelease: "ZT1JOQLn",

		// user configurable values

		ACLAgentMasterToken:              "furuQD0b",
		ACLAgentToken:                    "cOshLOQ2",
		ACLDatacenter:                    "m3urck3z",
		ACLDefaultPolicy:                 "ArK3WIfE",
		ACLDownPolicy:                    "vZXMfMP0",
		ACLEnforceVersion8:               true,
		ACLMasterToken:                   "C1Q1oIwh",
		ACLReplicationToken:              "LMmgy5dO",
		ACLTTL:                           18060 * time.Second,
		ACLToken:                         "O1El0wan",
		AdvertiseAddrLAN:                 tcpAddr("17.99.29.16:3757"),
		AdvertiseAddrWAN:                 tcpAddr("78.63.37.19:3757"),
		AutopilotCleanupDeadServers:      true,
		AutopilotDisableUpgradeMigration: true,
		AutopilotLastContactThreshold:    12705 * time.Second,
		AutopilotMaxTrailingLogs:         17849,
		AutopilotRedundancyZoneTag:       "3IsufDJf",
		AutopilotServerStabilizationTime: 23057 * time.Second,
		AutopilotUpgradeVersionTag:       "W9pDwFAL",
		BindAddr:                         ipAddr("16.99.34.17"),
		Bootstrap:                        true,
		BootstrapExpect:                  53,
		CAFile:                           "erA7T0PM",
		CAPath:                           "mQEN1Mfp",
		CertFile:                         "7s4QAzDk",
		Checks: []*structs.CheckDefinition{
			&structs.CheckDefinition{
				ID:        "fZaCAXww",
				Name:      "OOM2eo0f",
				Notes:     "zXzXI9Gt",
				ServiceID: "L8G0QNmR",
				Token:     "oo4BCTgJ",
				Status:    "qLykAl5u",
				Script:    "dhGfIF8n",
				HTTP:      "29B93haH",
				Header: map[string][]string{
					"hBq0zn1q": {"2a9o9ZKP", "vKwA5lR6"},
					"f3r6xFtM": {"RyuIdDWv", "QbxEcIUM"},
				},
				Method:            "Dou0nGT5",
				TCP:               "JY6fTTcw",
				Interval:          18714 * time.Second,
				DockerContainerID: "qF66POS9",
				Shell:             "sOnDy228",
				TLSSkipVerify:     true,
				Timeout:           5954 * time.Second,
				TTL:               30044 * time.Second,
				DeregisterCriticalServiceAfter: 13209 * time.Second,
			},
			&structs.CheckDefinition{
				ID:        "uAjE6m9Z",
				Name:      "QsZRGpYr",
				Notes:     "VJ7Sk4BY",
				ServiceID: "lSulPcyz",
				Token:     "toO59sh8",
				Status:    "9RlWsXMV",
				Script:    "8qbd8tWw",
				HTTP:      "dohLcyQ2",
				Header: map[string][]string{
					"ZBfTin3L": []string{"1sDbEqYG", "lJGASsWK"},
					"Ui0nU99X": []string{"LMccm3Qe", "k5H5RggQ"},
				},
				Method:            "aldrIQ4l",
				TCP:               "RJQND605",
				Interval:          22164 * time.Second,
				DockerContainerID: "ipgdFtjd",
				Shell:             "qAeOYy0M",
				TLSSkipVerify:     true,
				Timeout:           1813 * time.Second,
				TTL:               21743 * time.Second,
				DeregisterCriticalServiceAfter: 14232 * time.Second,
			},
			&structs.CheckDefinition{
				ID:        "Cqq95BhP",
				Name:      "3qXpkS0i",
				Notes:     "sb5qLTex",
				ServiceID: "CmUUcRna",
				Token:     "a3nQzHuy",
				Status:    "irj26nf3",
				Script:    "FJsI1oXt",
				HTTP:      "yzhgsQ7Y",
				Header: map[string][]string{
					"zcqwA8dO": []string{"qb1zx0DL", "sXCxPFsD"},
					"qxvdnSE9": []string{"6wBPUYdF", "YYh8wtSZ"},
				},
				Method:            "gLrztrNw",
				TCP:               "4jG5casb",
				Interval:          28767 * time.Second,
				DockerContainerID: "THW6u7rL",
				Shell:             "C1Zt3Zwh",
				TLSSkipVerify:     true,
				Timeout:           18506 * time.Second,
				TTL:               31006 * time.Second,
				DeregisterCriticalServiceAfter: 2366 * time.Second,
			},
		},
		CheckUpdateInterval:       16507 * time.Second,
		ClientAddrs:               []*net.IPAddr{ipAddr("93.83.18.19")},
		ConsulConfig:              devConsulConfig(consul.DefaultConfig()),
		DNSAddrs:                  []net.Addr{tcpAddr("93.95.95.81:7001"), udpAddr("93.95.95.81:7001")},
		DNSAllowStale:             true,
		DNSDisableCompression:     true,
		DNSDomain:                 "7W1xXSqd",
		DNSEnableTruncate:         true,
		DNSMaxStale:               29685 * time.Second,
		DNSNodeTTL:                7084 * time.Second,
		DNSOnlyPassing:            true,
		DNSPort:                   7001,
		DNSRecursorTimeout:        4427 * time.Second,
		DNSRecursors:              []string{"EZX7MOYF", "FtFhoUHl", "UYkwck1k"},
		DNSServiceTTL:             map[string]time.Duration{"*": 32030 * time.Second},
		DNSUDPAnswerLimit:         29909,
		DataDir:                   "oTOOIoV9",
		Datacenter:                "rzo029wg",
		DevMode:                   true,
		DisableAnonymousSignature: true,
		DisableCoordinates:        true,
		DisableHostNodeID:         true,
		DisableKeyringFile:        true,
		DisableRemoteExec:         true,
		DisableUpdateCheck:        true,
		EnableACLReplication:      true,
		EnableDebug:               true,
		EnableScriptChecks:        true,
		EnableSyslog:              true,
		EnableUI:                  true,
		EncryptKey:                "A4wELWqH",
		EncryptVerifyIncoming:     true,
		EncryptVerifyOutgoing:     true,
		HTTPAddrs:                 []net.Addr{tcpAddr("83.39.91.39:7999")},
		HTTPBlockEndpoints:        []string{"RBvAFcGD", "fWOWFznh"},
		HTTPPort:                  7999,
		HTTPResponseHeaders:       map[string]string{"M6TKa9NP": "xjuxjOzQ", "JRCrHZed": "rl0mTx81"},
		HTTPSAddrs:                []net.Addr{tcpAddr("95.17.17.19:15127")},
		HTTPSPort:                 15127,
		KeyFile:                   "IEkkwgIA",
		LeaveOnTerm:               true,
		LogLevel:                  "k1zo9Spt",
		NodeID:                    types.NodeID("AsUIlw99"),
		NodeMeta:                  map[string]string{"5mgGQMBk": "mJLtVMSG", "A7ynFMJB": "0Nx6RGab"},
		NodeName:                  "otlLxGaI",
		NonVotingServer:           true,
		PerformanceRaftMultiplier: 5,
		PidFile:                   "43xN80Km",
		RPCAdvertiseAddr:          tcpAddr("28.27.94.38:3757"),
		RPCBindAddr:               tcpAddr("16.99.34.17:3757"),
		RPCProtocol:               30793,
		RPCRateLimit:              12029.43,
		RPCMaxBurst:               44848,
		RaftProtocol:              19016,
		ReconnectTimeoutLAN:       23739 * time.Second,
		ReconnectTimeoutWAN:       26694 * time.Second,
		RejoinAfterLeave:          true,
		RetryJoinIntervalLAN:      8067 * time.Second,
		RetryJoinIntervalWAN:      28866 * time.Second,
		RetryJoinLAN:              []string{"pbsSFY7U", "l0qLtWij"},
		RetryJoinMaxAttemptsLAN:   913,
		RetryJoinMaxAttemptsWAN:   23160,
		RetryJoinWAN:              []string{"PFsR02Ye", "rJdQIhER"},
		SegmentName:               "BC2NhTDi",
		Segments: []structs.NetworkSegment{
			{
				Name:        "PExYMe2E",
				Bind:        tcpAddr("36.73.36.19:38295"),
				Advertise:   tcpAddr("63.39.19.18:38295"),
				RPCListener: true,
			},
			{
				Name:        "UzCvJgup",
				Bind:        tcpAddr("37.58.38.19:39292"),
				Advertise:   tcpAddr("83.58.26.27:39292"),
				RPCListener: true,
			},
		},
		SerfPortLAN: 8301,
		SerfPortWAN: 8302,
		ServerMode:  true,
		ServerName:  "Oerr9n1G",
		ServerPort:  3757,
		Services: []*structs.ServiceDefinition{
			{
				ID:                "wI1dzxS4",
				Name:              "7IszXMQ1",
				Tags:              []string{"0Zwg8l6v", "zebELdN5"},
				Address:           "9RhqPSPB",
				Token:             "myjKJkWH",
				Port:              72219,
				EnableTagOverride: true,
				Checks: []*structs.CheckType{
					&structs.CheckType{
						CheckID: "qmfeO5if",
						Name:    "atDGP7n5",
						Status:  "pDQKEhWL",
						Notes:   "Yt8EDLev",
						Script:  "MDu7wjlD",
						HTTP:    "qzHYvmJO",
						Header: map[string][]string{
							"UkpmZ3a3": {"2dfzXuxZ"},
							"cVFpko4u": {"gGqdEB6k", "9LsRo22u"},
						},
						Method:            "X5DrovFc",
						TCP:               "ICbxkpSF",
						Interval:          24392 * time.Second,
						DockerContainerID: "ZKXr68Yb",
						Shell:             "CEfzx0Fo",
						TLSSkipVerify:     true,
						Timeout:           38333 * time.Second,
						TTL:               57201 * time.Second,
						DeregisterCriticalServiceAfter: 44214 * time.Second,
					},
				},
			},
			{
				ID:                "MRHVMZuD",
				Name:              "6L6BVfgH",
				Tags:              []string{"7Ale4y6o", "PMBW08hy"},
				Address:           "R6H6g8h0",
				Token:             "ZgY8gjMI",
				Port:              38292,
				EnableTagOverride: true,
				Checks: structs.CheckTypes{
					&structs.CheckType{
						CheckID: "GTti9hCo",
						Name:    "9OOS93ne",
						Notes:   "CQy86DH0",
						Status:  "P0SWDvrk",
						Script:  "6BhLJ7R9",
						HTTP:    "u97ByEiW",
						Header: map[string][]string{
							"MUlReo8L": {"AUZG7wHG", "gsN0Dc2N"},
							"1UJXjVrT": {"OJgxzTfk", "xZZrFsq7"},
						},
						Method:            "5wkAxCUE",
						TCP:               "MN3oA9D2",
						Interval:          32718 * time.Second,
						DockerContainerID: "cU15LMet",
						Shell:             "nEz9qz2l",
						TLSSkipVerify:     true,
						Timeout:           34738 * time.Second,
						TTL:               22773 * time.Second,
						DeregisterCriticalServiceAfter: 84282 * time.Second,
					},
					&structs.CheckType{
						CheckID: "UHsDeLxG",
						Name:    "PQSaPWlT",
						Notes:   "jKChDOdl",
						Status:  "5qFz6OZn",
						Script:  "PbdxFZ3K",
						HTTP:    "1LBDJhw4",
						Header: map[string][]string{
							"cXPmnv1M": {"imDqfaBx", "NFxZ1bQe"},
							"vr7wY7CS": {"EtCoNPPL", "9vAarJ5s"},
						},
						Method:            "wzByP903",
						TCP:               "2exjZIGE",
						Interval:          5656 * time.Second,
						DockerContainerID: "5tDBWpfA",
						Shell:             "rlTpLM8s",
						TLSSkipVerify:     true,
						Timeout:           4868 * time.Second,
						TTL:               11222 * time.Second,
						DeregisterCriticalServiceAfter: 68482 * time.Second,
					},
				},
			},
			{
				ID:                "dLOXpSCI",
				Name:              "o1ynPkp0",
				Tags:              []string{"nkwshvM5", "NTDWn3ek"},
				Address:           "cOlSOhbp",
				Token:             "msy7iWER",
				Port:              24237,
				EnableTagOverride: true,
				Checks: structs.CheckTypes{
					&structs.CheckType{
						CheckID: "Zv99e9Ka",
						Name:    "sgV4F7Pk",
						Notes:   "yP5nKbW0",
						Status:  "7oLMEyfu",
						Script:  "NlUQ3nTE",
						HTTP:    "KyDjGY9H",
						Header: map[string][]string{
							"gv5qefTz": {"5Olo2pMG", "PvvKWQU5"},
							"SHOVq1Vv": {"jntFhyym", "GYJh32pp"},
						},
						Method:            "T66MFBfR",
						TCP:               "bNnNfx2A",
						Interval:          22224 * time.Second,
						DockerContainerID: "ipgdFtjd",
						Shell:             "omVZq7Sz",
						TLSSkipVerify:     true,
						Timeout:           18913 * time.Second,
						TTL:               44743 * time.Second,
						DeregisterCriticalServiceAfter: 8482 * time.Second,
					},
					&structs.CheckType{
						CheckID: "G79O6Mpr",
						Name:    "IEqrzrsd",
						Notes:   "SVqApqeM",
						Status:  "XXkVoZXt",
						Script:  "IXLZTM6E",
						HTTP:    "kyICZsn8",
						Header: map[string][]string{
							"4ebP5vL4": {"G20SrL5Q", "DwPKlMbo"},
							"p2UI34Qz": {"UsG1D0Qh", "NHhRiB6s"},
						},
						Method:            "ciYHWors",
						TCP:               "FfvCwlqH",
						Interval:          12356 * time.Second,
						DockerContainerID: "HBndBU6R",
						Shell:             "hVI33JjA",
						TLSSkipVerify:     true,
						Timeout:           38282 * time.Second,
						TTL:               1181 * time.Second,
						DeregisterCriticalServiceAfter: 4992 * time.Second,
					},
					&structs.CheckType{
						CheckID: "RMi85Dv8",
						Name:    "iehanzuq",
						Status:  "rCvn53TH",
						Notes:   "fti5lfF3",
						Script:  "rtj34nfd",
						HTTP:    "dl3Fgme3",
						Header: map[string][]string{
							"rjm4DEd3": {"2m3m2Fls"},
							"l4HwQ112": {"fk56MNlo", "dhLK56aZ"},
						},
						Method:            "9afLm3Mj",
						TCP:               "fjiLFqVd",
						Interval:          23926 * time.Second,
						DockerContainerID: "dO5TtRHk",
						Shell:             "e6q2ttES",
						TLSSkipVerify:     true,
						Timeout:           38483 * time.Second,
						TTL:               10943 * time.Second,
						DeregisterCriticalServiceAfter: 68787 * time.Second,
					},
				},
			},
		},
		SerfAdvertiseAddrLAN:                        tcpAddr("49.38.36.95:8301"),
		SerfAdvertiseAddrWAN:                        tcpAddr("63.38.52.13:8302"),
		SerfBindAddrLAN:                             tcpAddr("99.43.63.15:8301"),
		SerfBindAddrWAN:                             tcpAddr("67.88.33.19:8302"),
		SessionTTLMin:                               26627 * time.Second,
		SkipLeaveOnInt:                              true,
		StartJoinAddrsLAN:                           []string{"LR3hGDoG", "MwVpZ4Up"},
		StartJoinAddrsWAN:                           []string{"EbFSc3nA", "kwXTh623"},
		SyslogFacility:                              "hHv79Uia",
		TelemetryCirconusAPIApp:                     "p4QOTe9j",
		TelemetryCirconusAPIToken:                   "E3j35V23",
		TelemetryCirconusAPIURL:                     "mEMjHpGg",
		TelemetryCirconusBrokerID:                   "BHlxUhed",
		TelemetryCirconusBrokerSelectTag:            "13xy1gHm",
		TelemetryCirconusCheckDisplayName:           "DRSlQR6n",
		TelemetryCirconusCheckForceMetricActivation: "Ua5FGVYf",
		TelemetryCirconusCheckID:                    "kGorutad",
		TelemetryCirconusCheckInstanceID:            "rwoOL6R4",
		TelemetryCirconusCheckSearchTag:             "ovT4hT4f",
		TelemetryCirconusCheckTags:                  "prvO4uBl",
		TelemetryCirconusSubmissionInterval:         "DolzaflP",
		TelemetryCirconusSubmissionURL:              "gTcbS93G",
		TelemetryDisableHostname:                    true,
		TelemetryDogstatsdAddr:                      "0wSndumK",
		TelemetryDogstatsdTags:                      []string{"3N81zSUB", "Xtj8AnXZ"},
		TelemetryFilterDefault:                      true,
		TelemetryAllowedPrefixes:                    []string{"oJotS8XJ"},
		TelemetryBlockedPrefixes:                    []string{"cazlEhGn"},
		TelemetryStatsdAddr:                         "drce87cy",
		TelemetryStatsiteAddr:                       "HpFwKB8R",
		TelemetryStatsitePrefix:                     "ftO6DySn",
		TLSCipherSuites:                             []uint16{tls.TLS_RSA_WITH_RC4_128_SHA, tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA},
		TLSMinVersion:                               "pAOWafkR",
		TLSPreferServerCipherSuites:                 true,
		TaggedAddresses: map[string]string{
			"7MYgHrYH": "dALJAhLD",
			"h6DdBy6K": "ebrr9zZ8",
			"lan":      "17.99.29.16",
			"wan":      "78.63.37.19",
		},
		TranslateWANAddrs:    true,
		UIDir:                "11IFzAUn",
		UnixSocketUser:       "E0nB1DwA",
		UnixSocketGroup:      "8pFodrV8",
		UnixSocketMode:       "E8sAwOv4",
		VerifyIncoming:       true,
		VerifyIncomingHTTPS:  true,
		VerifyIncomingRPC:    true,
		VerifyOutgoing:       true,
		VerifyServerHostname: true,
		Watches: []map[string]interface{}{
			map[string]interface{}{
				"type":       "key",
				"datacenter": "GyE6jpeW",
				"key":        "j9lF1Tve",
				"handler":    "90N7S4LN",
			},
		},
	}

	warns := []string{
		`==> DEPRECATION: "addresses.rpc" is deprecated and is no longer used. Please remove it from your configuration.`,
		`==> DEPRECATION: "ports.rpc" is deprecated and is no longer used. Please remove it from your configuration.`,
		`bootstrap_expect > 0: expecting 53 servers`,
	}

	// ensure that all fields are set to unique non-zero values
	// todo(fs): This currently fails since ServiceDefinition.Check is not used
	// todo(fs): not sure on how to work around this. Possible options are:
	// todo(fs):  * move first check into the Check field
	// todo(fs):  * ignore the Check field
	// todo(fs): both feel like a hack
	if err := nonZero("RuntimeConfig", nil, want); err != nil {
		t.Log(err)
	}

	for format, data := range src {
		t.Run(format, func(t *testing.T) {
			// parse the flags since this is the only way we can set the
			// DevMode flag
			var flags Flags
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			AddFlags(fs, &flags)
			if err := fs.Parse(flagSrc); err != nil {
				t.Fatalf("ParseFlags: %s", err)
			}

			// ensure that all fields are set to unique non-zero values
			// if err := nonZero("Config", nil, c); err != nil {
			// 	t.Fatal(err)
			// }

			b := &Builder{
				Flags:   flags,
				Head:    []Source{DefaultSource},
				Sources: []Source{{Name: "full", Format: format, Data: data}},
				Tail: []Source{
					nonUserSource[format],
					VersionSource("JNtPSav3", "R909Hblt", "ZT1JOQLn"),
				},
			}

			// construct the runtime config
			rt, err := b.Build()
			if err != nil {
				t.Fatalf("Build: %s", err)
			}

			// verify that all fields are set
			if !verify.Values(t, "", rt, want) {
				t.FailNow()
			}

			// at this point we have confirmed that the parsing worked
			// for all fields but the validation will fail since certain
			// combinations are not allowed. Since it is not possible to have
			// all fields with non-zero values and to have a valid configuration
			// we are patching a handful of safe fields to make validation pass.
			rt.Bootstrap = false
			rt.DevMode = false
			rt.EnableUI = false
			rt.SegmentName = ""

			// validate the runtime config
			if err := b.Validate(rt); err != nil {
				t.Fatalf("Validate: %s", err)
			}

			// check the warnings
			if got, want := b.Warnings, warns; !verify.Values(t, "warnings", got, want) {
				t.FailNow()
			}
		})
	}
}

// nonZero verifies recursively that all fields are set to unique,
// non-zero and non-nil values.
//
// struct: check all fields recursively
// slice: check len > 0 and all values recursively
// ptr: check not nil
// bool: check not zero (cannot check uniqueness)
// string, int, uint: check not zero and unique
// other: error
func nonZero(name string, uniq map[interface{}]string, v interface{}) error {
	if v == nil {
		return fmt.Errorf("%q is nil", name)
	}

	if uniq == nil {
		uniq = map[interface{}]string{}
	}

	isUnique := func(v interface{}) error {
		if other := uniq[v]; other != "" {
			return fmt.Errorf("%q and %q both use vaule %q", name, other, v)
		}
		uniq[v] = name
		return nil
	}

	val, typ := reflect.ValueOf(v), reflect.TypeOf(v)
	// fmt.Printf("%s: %T\n", name, v)
	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			fieldname := fmt.Sprintf("%s.%s", name, f.Name)
			err := nonZero(fieldname, uniq, val.Field(i).Interface())
			if err != nil {
				return err
			}
		}

	case reflect.Slice:
		if val.Len() == 0 {
			return fmt.Errorf("%q is empty slice", name)
		}
		for i := 0; i < val.Len(); i++ {
			elemname := fmt.Sprintf("%s[%d]", name, i)
			err := nonZero(elemname, uniq, val.Index(i).Interface())
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		if val.Len() == 0 {
			return fmt.Errorf("%q is empty map", name)
		}
		for _, key := range val.MapKeys() {
			keyname := fmt.Sprintf("%s[%s]", name, key.String())
			if err := nonZero(keyname, uniq, key.Interface()); err != nil {
				if strings.Contains(err.Error(), "is zero value") {
					return fmt.Errorf("%q has zero value map key", name)
				}
				return err
			}
			if err := nonZero(keyname, uniq, val.MapIndex(key).Interface()); err != nil {
				return err
			}
		}

	case reflect.Bool:
		if val.Bool() != true {
			return fmt.Errorf("%q is zero value", name)
		}
		// do not test bool for uniqueness since there are only two values

	case reflect.String:
		if val.Len() == 0 {
			return fmt.Errorf("%q is zero value", name)
		}
		return isUnique(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() == 0 {
			return fmt.Errorf("%q is zero value", name)
		}
		return isUnique(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Uint() == 0 {
			return fmt.Errorf("%q is zero value", name)
		}
		return isUnique(v)

	case reflect.Float32, reflect.Float64:
		if val.Float() == 0 {
			return fmt.Errorf("%q is zero value", name)
		}
		return isUnique(v)

	case reflect.Ptr:
		if val.IsNil() {
			return fmt.Errorf("%q is nil", name)
		}
		return nonZero("*"+name, uniq, val.Elem().Interface())

	default:
		return fmt.Errorf("%T is not supported", v)
	}
	return nil
}

func TestNonZero(t *testing.T) {
	var empty string

	tests := []struct {
		desc string
		v    interface{}
		err  error
	}{
		{"nil", nil, errors.New(`"x" is nil`)},
		{"zero bool", false, errors.New(`"x" is zero value`)},
		{"zero string", "", errors.New(`"x" is zero value`)},
		{"zero int", int(0), errors.New(`"x" is zero value`)},
		{"zero int8", int8(0), errors.New(`"x" is zero value`)},
		{"zero int16", int16(0), errors.New(`"x" is zero value`)},
		{"zero int32", int32(0), errors.New(`"x" is zero value`)},
		{"zero int64", int64(0), errors.New(`"x" is zero value`)},
		{"zero uint", uint(0), errors.New(`"x" is zero value`)},
		{"zero uint8", uint8(0), errors.New(`"x" is zero value`)},
		{"zero uint16", uint16(0), errors.New(`"x" is zero value`)},
		{"zero uint32", uint32(0), errors.New(`"x" is zero value`)},
		{"zero uint64", uint64(0), errors.New(`"x" is zero value`)},
		{"zero float32", float32(0), errors.New(`"x" is zero value`)},
		{"zero float64", float64(0), errors.New(`"x" is zero value`)},
		{"ptr to zero value", &empty, errors.New(`"*x" is zero value`)},
		{"empty slice", []string{}, errors.New(`"x" is empty slice`)},
		{"slice with zero value", []string{""}, errors.New(`"x[0]" is zero value`)},
		{"empty map", map[string]string{}, errors.New(`"x" is empty map`)},
		{"map with zero value key", map[string]string{"": "y"}, errors.New(`"x" has zero value map key`)},
		{"map with zero value elem", map[string]string{"y": ""}, errors.New(`"x[y]" is zero value`)},
		{"struct with nil field", struct{ Y *int }{}, errors.New(`"x.Y" is nil`)},
		{"struct with zero value field", struct{ Y string }{}, errors.New(`"x.Y" is zero value`)},
		{"struct with empty array", struct{ Y []string }{}, errors.New(`"x.Y" is empty slice`)},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got, want := nonZero("x", nil, tt.v), tt.err; !reflect.DeepEqual(got, want) {
				t.Fatalf("got error %v want %v", got, want)
			}
		})
	}
}

func TestConfigDecodeBytes(t *testing.T) {
	t.Parallel()
	// Test with some input
	src := []byte("abc")
	key := base64.StdEncoding.EncodeToString(src)

	result, err := decodeBytes(key)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !bytes.Equal(src, result) {
		t.Fatalf("bad: %#v", result)
	}

	// Test with no input
	result, err = decodeBytes("")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(result) > 0 {
		t.Fatalf("bad: %#v", result)
	}
}
