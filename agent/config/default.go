package config

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/agent/consul"
	"github.com/hashicorp/consul/version"
)

// DefaultSource is the default agent configuration.
var DefaultSource = Source{
	Name:   "default",
	Format: "hcl",
	Data: `
		acl_default_policy = "allow"
		acl_down_policy = "extend-cache"
		acl_enforce_version8 = true
		acl_ttl = "30s"
		bind_addr = "0.0.0.0"
		bootstrap = false
		bootstrap_expect = 0
		check_update_interval = "5m"
		client_addr = "127.0.0.1"
		datacenter = "dc1"
		disable_coordinates = false
		disable_host_node_id = true
		disable_remote_exec = true
		domain = "consul."
		encrypt_verify_incoming = true
		encrypt_verify_outgoing = true
		log_level = "INFO"
		protocol =  2
		retry_interval = "30s"
		retry_interval_wan = "30s"
		server = false
		syslog_facility = "LOCAL0"
		tls_min_version = "tls10"

		dns_config = {
			allow_stale = true
			udp_answer_limit = 3
			max_stale = "87600h"
			recursor_timeout = "2s"
		}
		limits = {
			rpc_rate = -1
			rpc_max_burst = 1000
		}
		ports = {
			dns = 8600
			http = 8500
			https = -1
			serf_lan = 8301
			serf_wan = 8302
			server = 8300
		}
		telemetry = {
			statsite_prefix = "consul"
			filter_default = true
		}
	`,
}

// DevSource is the additional default configuration for dev mode.
// This should be loaded after the default configuration.
var DevSource = Source{
	Name:   "dev",
	Format: "hcl",
	Data: `
		bind_addr = "127.0.0.1"
		disable_anonymous_signature = true
		disable_keyring_file = true
		enable_debug = true
		enable_ui = true
		log_level = "DEBUG"
		server = true
	`,
}

// NonUserSource contains the values the user cannot configure.
// This needs to be merged last.
var NonUserSource = Source{
	Name:   "non-user",
	Format: "hcl",
	Data: `
		acl_disabled_ttl = "120s"
		check_deregister_interval_min = "1m"
		check_reap_interval = "30s"
		ae_interval = "1m"
		sync_coordinate_rate_target = 64
		sync_coordinate_interval_min = "15s"
	`,
}

// VersionSource creates a config source for the version parameters.
func VersionSource(rev, ver, verPre string) Source {
	return Source{
		Name:   "version",
		Format: "hcl",
		Data:   fmt.Sprintf(`revision = %q version = %q version_prerelease = %q`, rev, ver, verPre),
	}
}

// DefaultVersionSource returns the version config source for the embedded
// version numbers.
func DefaultVersionSource() Source {
	return VersionSource(version.GitCommit, version.Version, version.VersionPrerelease)
}

func DefaultRuntimeConfig() *RuntimeConfig {
	b := &Builder{
		Head: []Source{DefaultSource},
		Tail: []Source{NonUserSource, DefaultVersionSource()},
	}
	rt, _ := b.BuildAndValidate()
	return &rt
}

func devConsulConfig(conf *consul.Config) *consul.Config {
	conf.SerfLANConfig.MemberlistConfig.ProbeTimeout = 100 * time.Millisecond
	conf.SerfLANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
	conf.SerfLANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond

	conf.SerfWANConfig.MemberlistConfig.SuspicionMult = 3
	conf.SerfWANConfig.MemberlistConfig.ProbeTimeout = 100 * time.Millisecond
	conf.SerfWANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
	conf.SerfWANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond

	conf.RaftConfig.LeaderLeaseTimeout = 20 * time.Millisecond
	conf.RaftConfig.HeartbeatTimeout = 40 * time.Millisecond
	conf.RaftConfig.ElectionTimeout = 40 * time.Millisecond

	conf.CoordinateUpdatePeriod = 100 * time.Millisecond

	return conf
}
