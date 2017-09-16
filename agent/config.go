package agent

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// SegmentLimit is the maximum number of network segments that may be declared.
	SegmentLimit = 64

	// SegmentNameLimit is the maximum segment name length.
	SegmentNameLimit = 64
)

// done(fs): // Ports is used to simplify the configuration by
// done(fs): // providing default ports, and allowing the addresses
// done(fs): // to only be specified once
// done(fs): type PortConfig struct {
// done(fs): 	DNS     int // DNS Query interface
// done(fs): 	HTTP    int // HTTP API
// done(fs): 	HTTPS   int // HTTPS API
// done(fs): 	SerfLan int `mapstructure:"serf_lan"` // LAN gossip (Client + Server)
// done(fs): 	SerfWan int `mapstructure:"serf_wan"` // WAN gossip (Server only)
// done(fs): 	Server  int // Server internal RPC
// done(fs):
// done(fs): 	// RPC is deprecated and is no longer used. It will be removed in a future
// done(fs): 	// version.
// done(fs): 	RPC int // CLI RPC
// done(fs): }
// done(fs):
// done(fs): // AddressConfig is used to provide address overrides
// done(fs): // for specific services. By default, either ClientAddress
// done(fs): // or ServerAddress is used.
// done(fs): type AddressConfig struct {
// done(fs): 	DNS   string // DNS Query interface
// done(fs): 	HTTP  string // HTTP API
// done(fs): 	HTTPS string // HTTPS API
// done(fs):
// done(fs): 	// RPC is deprecated and is no longer used. It will be removed in a future
// done(fs): 	// version.
// done(fs): 	RPC string // CLI RPC
// done(fs): }
// done(fs):
// done(fs): type AdvertiseAddrsConfig struct {
// done(fs): 	SerfLan    *net.TCPAddr `mapstructure:"-"`
// done(fs): 	SerfLanRaw string       `mapstructure:"serf_lan"`
// done(fs): 	SerfWan    *net.TCPAddr `mapstructure:"-"`
// done(fs): 	SerfWanRaw string       `mapstructure:"serf_wan"`
// done(fs): 	RPC        *net.TCPAddr `mapstructure:"-"`
// done(fs): 	RPCRaw     string       `mapstructure:"rpc"`
// done(fs): }
// done(fs):
// done(fs): // DNSConfig is used to fine tune the DNS sub-system.
// done(fs): // It can be used to control cache values, and stale
// done(fs): // reads
// done(fs): type DNSConfig struct {
// done(fs): 	// NodeTTL provides the TTL value for a node query
// done(fs): 	NodeTTL    time.Duration `mapstructure:"-"`
// done(fs): 	NodeTTLRaw string        `mapstructure:"node_ttl" json:"-"`
// done(fs):
// done(fs): 	// ServiceTTL provides the TTL value for a service
// done(fs): 	// query for given service. The "*" wildcard can be used
// done(fs): 	// to set a default for all services.
// done(fs): 	ServiceTTL    map[string]time.Duration `mapstructure:"-"`
// done(fs): 	ServiceTTLRaw map[string]string        `mapstructure:"service_ttl" json:"-"`
// done(fs):
// done(fs): 	// AllowStale is used to enable lookups with stale
// done(fs): 	// data. This gives horizontal read scalability since
// done(fs): 	// any Consul server can service the query instead of
// done(fs): 	// only the leader.
// done(fs): 	AllowStale *bool `mapstructure:"allow_stale"`
// done(fs):
// done(fs): 	// EnableTruncate is used to enable setting the truncate
// done(fs): 	// flag for UDP DNS queries.  This allows unmodified
// done(fs): 	// clients to re-query the consul server using TCP
// done(fs): 	// when the total number of records exceeds the number
// done(fs): 	// returned by default for UDP.
// done(fs): 	EnableTruncate bool `mapstructure:"enable_truncate"`
// done(fs):
// done(fs): 	// UDPAnswerLimit is used to limit the maximum number of DNS Resource
// done(fs): 	// Records returned in the ANSWER section of a DNS response. This is
// done(fs): 	// not normally useful and will be limited based on the querying
// done(fs): 	// protocol, however systems that implemented §6 Rule 9 in RFC3484
// done(fs): 	// may want to set this to `1` in order to subvert §6 Rule 9 and
// done(fs): 	// re-obtain the effect of randomized resource records (i.e. each
// done(fs): 	// answer contains only one IP, but the IP changes every request).
// done(fs): 	// RFC3484 sorts answers in a deterministic order, which defeats the
// done(fs): 	// purpose of randomized DNS responses.  This RFC has been obsoleted
// done(fs): 	// by RFC6724 and restores the desired behavior of randomized
// done(fs): 	// responses, however a large number of Linux hosts using glibc(3)
// done(fs): 	// implemented §6 Rule 9 and may need this option (e.g. CentOS 5-6,
// done(fs): 	// Debian Squeeze, etc).
// done(fs): 	UDPAnswerLimit int `mapstructure:"udp_answer_limit"`
// done(fs):
// done(fs): 	// MaxStale is used to bound how stale of a result is
// done(fs): 	// accepted for a DNS lookup. This can be used with
// done(fs): 	// AllowStale to limit how old of a value is served up.
// done(fs): 	// If the stale result exceeds this, another non-stale
// done(fs): 	// stale read is performed.
// done(fs): 	MaxStale    time.Duration `mapstructure:"-"`
// done(fs): 	MaxStaleRaw string        `mapstructure:"max_stale" json:"-"`
// done(fs):
// done(fs): 	// OnlyPassing is used to determine whether to filter nodes
// done(fs): 	// whose health checks are in any non-passing state. By
// done(fs): 	// default, only nodes in a critical state are excluded.
// done(fs): 	OnlyPassing bool `mapstructure:"only_passing"`
// done(fs):
// done(fs): 	// DisableCompression is used to control whether DNS responses are
// done(fs): 	// compressed. In Consul 0.7 this was turned on by default and this
// done(fs): 	// config was added as an opt-out.
// done(fs): 	DisableCompression bool `mapstructure:"disable_compression"`
// done(fs):
// done(fs): 	// RecursorTimeout specifies the timeout in seconds
// done(fs): 	// for Consul's internal dns client used for recursion.
// done(fs): 	// This value is used for the connection, read and write timeout.
// done(fs): 	// Default: 2s
// done(fs): 	RecursorTimeout    time.Duration `mapstructure:"-"`
// done(fs): 	RecursorTimeoutRaw string        `mapstructure:"recursor_timeout" json:"-"`
// done(fs): }
// done(fs):
// done(fs): // HTTPConfig is used to fine tune the Http sub-system.
// done(fs): type HTTPConfig struct {
// done(fs): 	// BlockEndpoints is a list of endpoint prefixes to block in the
// done(fs): 	// HTTP API. Any requests to these will get a 403 response.
// done(fs): 	BlockEndpoints []string `mapstructure:"block_endpoints"`
// done(fs):
// done(fs): 	// ResponseHeaders are used to add HTTP header response fields to the HTTP API responses.
// done(fs): 	ResponseHeaders map[string]string `mapstructure:"response_headers"`
// done(fs): }
// done(fs):
// done(fs): // RetryJoinEC2 is used to configure discovery of instances via Amazon's EC2 api
// done(fs): type RetryJoinEC2 struct {
// done(fs): 	// The AWS region to look for instances in
// done(fs): 	Region string `mapstructure:"region"`
// done(fs):
// done(fs): 	// The tag key and value to use when filtering instances
// done(fs): 	TagKey   string `mapstructure:"tag_key"`
// done(fs): 	TagValue string `mapstructure:"tag_value"`
// done(fs):
// done(fs): 	// The AWS credentials to use for making requests to EC2
// done(fs): 	AccessKeyID     string `mapstructure:"access_key_id" json:"-"`
// done(fs): 	SecretAccessKey string `mapstructure:"secret_access_key" json:"-"`
// done(fs): }
// done(fs):
// done(fs): // RetryJoinGCE is used to configure discovery of instances via Google Compute
// done(fs): // Engine's API.
// done(fs): type RetryJoinGCE struct {
// done(fs): 	// The name of the project the instances reside in.
// done(fs): 	ProjectName string `mapstructure:"project_name"`
// done(fs):
// done(fs): 	// A regular expression (RE2) pattern for the zones you want to discover the instances in.
// done(fs): 	// Example: us-west1-.*, or us-(?west|east).*.
// done(fs): 	ZonePattern string `mapstructure:"zone_pattern"`
// done(fs):
// done(fs): 	// The tag value to search for when filtering instances.
// done(fs): 	TagValue string `mapstructure:"tag_value"`
// done(fs):
// done(fs): 	// A path to a JSON file with the service account credentials necessary to
// done(fs): 	// connect to GCE. If this is not defined, the following chain is respected:
// done(fs): 	// 1. A JSON file whose path is specified by the
// done(fs): 	//		GOOGLE_APPLICATION_CREDENTIALS environment variable.
// done(fs): 	// 2. A JSON file in a location known to the gcloud command-line tool.
// done(fs): 	//    On Windows, this is %APPDATA%/gcloud/application_default_credentials.json.
// done(fs): 	//  	On other systems, $HOME/.config/gcloud/application_default_credentials.json.
// done(fs): 	// 3. On Google Compute Engine, it fetches credentials from the metadata
// done(fs): 	//    server.  (In this final case any provided scopes are ignored.)
// done(fs): 	CredentialsFile string `mapstructure:"credentials_file"`
// done(fs): }
// done(fs):
// done(fs): // RetryJoinAzure is used to configure discovery of instances via AzureRM API
// done(fs): type RetryJoinAzure struct {
// done(fs): 	// The tag name and value to use when filtering instances
// done(fs): 	TagName  string `mapstructure:"tag_name"`
// done(fs): 	TagValue string `mapstructure:"tag_value"`
// done(fs):
// done(fs): 	// The Azure credentials to use for making requests to AzureRM
// done(fs): 	SubscriptionID  string `mapstructure:"subscription_id" json:"-"`
// done(fs): 	TenantID        string `mapstructure:"tenant_id" json:"-"`
// done(fs): 	ClientID        string `mapstructure:"client_id" json:"-"`
// done(fs): 	SecretAccessKey string `mapstructure:"secret_access_key" json:"-"`
// done(fs): }
// done(fs):
// done(fs): // Limits is used to configure limits enforced by the agent.
// done(fs): type Limits struct {
// done(fs): 	// RPCRate and RPCMaxBurst control how frequently RPC calls are allowed
// done(fs): 	// to happen. In any large enough time interval, rate limiter limits the
// done(fs): 	// rate to RPCRate tokens per second, with a maximum burst size of
// done(fs): 	// RPCMaxBurst events. As a special case, if RPCRate == Inf (the infinite
// done(fs): 	// rate), RPCMaxBurst is ignored.
// done(fs): 	//
// done(fs): 	// See https://en.wikipedia.org/wiki/Token_bucket for more about token
// done(fs): 	// buckets.
// done(fs): 	RPCRate     rate.Limit `mapstructure:"rpc_rate"`
// done(fs): 	RPCMaxBurst int        `mapstructure:"rpc_max_burst"`
// done(fs): }
// done(fs):
// done(fs): // Performance is used to tune the performance of Consul's subsystems.
// done(fs): type Performance struct {
// done(fs): 	// RaftMultiplier is an integer multiplier used to scale Raft timing
// done(fs): 	// parameters: HeartbeatTimeout, ElectionTimeout, and LeaderLeaseTimeout.
// done(fs): 	RaftMultiplier uint `mapstructure:"raft_multiplier"`
// done(fs): }
// done(fs):
// done(fs): // Telemetry is the telemetry configuration for the server
// done(fs): type Telemetry struct {
// done(fs): 	// StatsiteAddr is the address of a statsite instance. If provided,
// done(fs): 	// metrics will be streamed to that instance.
// done(fs): 	StatsiteAddr string `mapstructure:"statsite_address"`
// done(fs):
// done(fs): 	// StatsdAddr is the address of a statsd instance. If provided,
// done(fs): 	// metrics will be sent to that instance.
// done(fs): 	StatsdAddr string `mapstructure:"statsd_address"`
// done(fs):
// done(fs): 	// StatsitePrefix is the prefix used to write stats values to. By
// done(fs): 	// default this is set to 'consul'.
// done(fs): 	StatsitePrefix string `mapstructure:"statsite_prefix"`
// done(fs):
// done(fs): 	// DisableHostname will disable hostname prefixing for all metrics
// done(fs): 	DisableHostname bool `mapstructure:"disable_hostname"`
// done(fs):
// done(fs): 	// PrefixFilter is a list of filter rules to apply for allowing/blocking metrics
// done(fs): 	// by prefix.
// done(fs): 	PrefixFilter    []string `mapstructure:"prefix_filter"`
// done(fs): 	AllowedPrefixes []string `mapstructure:"-" json:"-"`
// done(fs): 	BlockedPrefixes []string `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// FilterDefault is the default for whether to allow a metric that's not
// done(fs): 	// covered by the filter.
// done(fs): 	FilterDefault *bool `mapstructure:"filter_default"`
// done(fs):
// done(fs): 	// DogStatsdAddr is the address of a dogstatsd instance. If provided,
// done(fs): 	// metrics will be sent to that instance
// done(fs): 	DogStatsdAddr string `mapstructure:"dogstatsd_addr"`
// done(fs):
// done(fs): 	// DogStatsdTags are the global tags that should be sent with each packet to dogstatsd
// done(fs): 	// It is a list of strings, where each string looks like "my_tag_name:my_tag_value"
// done(fs): 	DogStatsdTags []string `mapstructure:"dogstatsd_tags"`
// done(fs):
// done(fs): 	// Circonus: see https://github.com/circonus-labs/circonus-gometrics
// done(fs): 	// for more details on the various configuration options.
// done(fs): 	// Valid configuration combinations:
// done(fs): 	//    - CirconusAPIToken
// done(fs): 	//      metric management enabled (search for existing check or create a new one)
// done(fs): 	//    - CirconusSubmissionUrl
// done(fs): 	//      metric management disabled (use check with specified submission_url,
// done(fs): 	//      broker must be using a public SSL certificate)
// done(fs): 	//    - CirconusAPIToken + CirconusCheckSubmissionURL
// done(fs): 	//      metric management enabled (use check with specified submission_url)
// done(fs): 	//    - CirconusAPIToken + CirconusCheckID
// done(fs): 	//      metric management enabled (use check with specified id)
// done(fs):
// done(fs): 	// CirconusAPIToken is a valid API Token used to create/manage check. If provided,
// done(fs): 	// metric management is enabled.
// done(fs): 	// Default: none
// done(fs): 	CirconusAPIToken string `mapstructure:"circonus_api_token" json:"-"`
// done(fs): 	// CirconusAPIApp is an app name associated with API token.
// done(fs): 	// Default: "consul"
// done(fs): 	CirconusAPIApp string `mapstructure:"circonus_api_app"`
// done(fs): 	// CirconusAPIURL is the base URL to use for contacting the Circonus API.
// done(fs): 	// Default: "https://api.circonus.com/v2"
// done(fs): 	CirconusAPIURL string `mapstructure:"circonus_api_url"`
// done(fs): 	// CirconusSubmissionInterval is the interval at which metrics are submitted to Circonus.
// done(fs): 	// Default: 10s
// done(fs): 	CirconusSubmissionInterval string `mapstructure:"circonus_submission_interval"`
// done(fs): 	// CirconusCheckSubmissionURL is the check.config.submission_url field from a
// done(fs): 	// previously created HTTPTRAP check.
// done(fs): 	// Default: none
// done(fs): 	CirconusCheckSubmissionURL string `mapstructure:"circonus_submission_url"`
// done(fs): 	// CirconusCheckID is the check id (not check bundle id) from a previously created
// done(fs): 	// HTTPTRAP check. The numeric portion of the check._cid field.
// done(fs): 	// Default: none
// done(fs): 	CirconusCheckID string `mapstructure:"circonus_check_id"`
// done(fs): 	// CirconusCheckForceMetricActivation will force enabling metrics, as they are encountered,
// done(fs): 	// if the metric already exists and is NOT active. If check management is enabled, the default
// done(fs): 	// behavior is to add new metrics as they are encoutered. If the metric already exists in the
// done(fs): 	// check, it will *NOT* be activated. This setting overrides that behavior.
// done(fs): 	// Default: "false"
// done(fs): 	CirconusCheckForceMetricActivation string `mapstructure:"circonus_check_force_metric_activation"`
// done(fs): 	// CirconusCheckInstanceID serves to uniquely identify the metrics coming from this "instance".
// done(fs): 	// It can be used to maintain metric continuity with transient or ephemeral instances as
// done(fs): 	// they move around within an infrastructure.
// done(fs): 	// Default: hostname:app
// done(fs): 	CirconusCheckInstanceID string `mapstructure:"circonus_check_instance_id"`
// done(fs): 	// CirconusCheckSearchTag is a special tag which, when coupled with the instance id, helps to
// done(fs): 	// narrow down the search results when neither a Submission URL or Check ID is provided.
// done(fs): 	// Default: service:app (e.g. service:consul)
// done(fs): 	CirconusCheckSearchTag string `mapstructure:"circonus_check_search_tag"`
// done(fs): 	// CirconusCheckTags is a comma separated list of tags to apply to the check. Note that
// done(fs): 	// the value of CirconusCheckSearchTag will always be added to the check.
// done(fs): 	// Default: none
// done(fs): 	CirconusCheckTags string `mapstructure:"circonus_check_tags"`
// done(fs): 	// CirconusCheckDisplayName is the name for the check which will be displayed in the Circonus UI.
// done(fs): 	// Default: value of CirconusCheckInstanceID
// done(fs): 	CirconusCheckDisplayName string `mapstructure:"circonus_check_display_name"`
// done(fs): 	// CirconusBrokerID is an explicit broker to use when creating a new check. The numeric portion
// done(fs): 	// of broker._cid. If metric management is enabled and neither a Submission URL nor Check ID
// done(fs): 	// is provided, an attempt will be made to search for an existing check using Instance ID and
// done(fs): 	// Search Tag. If one is not found, a new HTTPTRAP check will be created.
// done(fs): 	// Default: use Select Tag if provided, otherwise, a random Enterprise Broker associated
// done(fs): 	// with the specified API token or the default Circonus Broker.
// done(fs): 	// Default: none
// done(fs): 	CirconusBrokerID string `mapstructure:"circonus_broker_id"`
// done(fs): 	// CirconusBrokerSelectTag is a special tag which will be used to select a broker when
// done(fs): 	// a Broker ID is not provided. The best use of this is to as a hint for which broker
// done(fs): 	// should be used based on *where* this particular instance is running.
// done(fs): 	// (e.g. a specific geo location or datacenter, dc:sfo)
// done(fs): 	// Default: none
// done(fs): 	CirconusBrokerSelectTag string `mapstructure:"circonus_broker_select_tag"`
// done(fs): }
// done(fs):
// done(fs): // Autopilot is used to configure helpful features for operating Consul servers.
// done(fs): type Autopilot struct {
// done(fs): 	// CleanupDeadServers enables the automatic cleanup of dead servers when new ones
// done(fs): 	// are added to the peer list. Defaults to true.
// done(fs): 	CleanupDeadServers *bool `mapstructure:"cleanup_dead_servers"`
// done(fs):
// done(fs): 	// LastContactThreshold is the limit on the amount of time a server can go
// done(fs): 	// without leader contact before being considered unhealthy.
// done(fs): 	LastContactThreshold    *time.Duration `mapstructure:"-" json:"-"`
// done(fs): 	LastContactThresholdRaw string         `mapstructure:"last_contact_threshold"`
// done(fs):
// done(fs): 	// MaxTrailingLogs is the amount of entries in the Raft Log that a server can
// done(fs): 	// be behind before being considered unhealthy.
// done(fs): 	MaxTrailingLogs *uint64 `mapstructure:"max_trailing_logs"`
// done(fs):
// done(fs): 	// ServerStabilizationTime is the minimum amount of time a server must be
// done(fs): 	// in a stable, healthy state before it can be added to the cluster. Only
// done(fs): 	// applicable with Raft protocol version 3 or higher.
// done(fs): 	ServerStabilizationTime    *time.Duration `mapstructure:"-" json:"-"`
// done(fs): 	ServerStabilizationTimeRaw string         `mapstructure:"server_stabilization_time"`
// done(fs):
// done(fs): 	// (Enterprise-only) RedundancyZoneTag is the Meta tag to use for separating servers
// done(fs): 	// into zones for redundancy. If left blank, this feature will be disabled.
// done(fs): 	RedundancyZoneTag string `mapstructure:"redundancy_zone_tag"`
// done(fs):
// done(fs): 	// (Enterprise-only) DisableUpgradeMigration will disable Autopilot's upgrade migration
// done(fs): 	// strategy of waiting until enough newer-versioned servers have been added to the
// done(fs): 	// cluster before promoting them to voters.
// done(fs): 	DisableUpgradeMigration *bool `mapstructure:"disable_upgrade_migration"`
// done(fs):
// done(fs): 	// (Enterprise-only) UpgradeVersionTag is the node tag to use for version info when
// done(fs): 	// performing upgrade migrations. If left blank, the Consul version will be used.
// done(fs): 	UpgradeVersionTag string `mapstructure:"upgrade_version_tag"`
// done(fs): }
// done(fs):
// done(fs): // (Enterprise-only) NetworkSegment is the configuration for a network segment, which is an
// done(fs): // isolated serf group on the LAN.
// done(fs): type NetworkSegment struct {
// done(fs): 	// Name is the name of the segment.
// done(fs): 	Name string `mapstructure:"name"`
// done(fs):
// done(fs): 	// Bind is the bind address for this segment.
// done(fs): 	Bind string `mapstructure:"bind"`
// done(fs):
// done(fs): 	// Port is the port for this segment.
// done(fs): 	Port int `mapstructure:"port"`
// done(fs):
// done(fs): 	// RPCListener is whether to bind a separate RPC listener on the bind address
// done(fs): 	// for this segment.
// done(fs): 	RPCListener bool `mapstructure:"rpc_listener"`
// done(fs):
// done(fs): 	// Advertise is the advertise address of this segment.
// done(fs): 	Advertise string `mapstructure:"advertise"`
// done(fs): }
// done(fs):
// done(fs): // Config is the configuration that can be set for an Agent.
// done(fs): // Some of this is configurable as CLI flags, but most must
// done(fs): // be set using a configuration file.
// done(fs): type Config struct {
// done(fs): 	// DevMode enables a fast-path mode of operation to bring up an in-memory
// done(fs): 	// server with minimal configuration. Useful for developing Consul.
// done(fs): 	DevMode bool `mapstructure:"-"`
// done(fs):
// done(fs): 	// Limits is used to configure limits enforced by the agent.
// done(fs): 	Limits Limits `mapstructure:"limits"`
// done(fs):
// done(fs): 	// Performance is used to tune the performance of Consul's subsystems.
// done(fs): 	Performance Performance `mapstructure:"performance"`
// done(fs):
// done(fs): 	// Bootstrap is used to bring up the first Consul server, and
// done(fs): 	// permits that node to elect itself leader
// done(fs): 	Bootstrap bool `mapstructure:"bootstrap"`
// done(fs):
// done(fs): 	// BootstrapExpect tries to automatically bootstrap the Consul cluster,
// done(fs): 	// by withholding peers until enough servers join.
// done(fs): 	BootstrapExpect int `mapstructure:"bootstrap_expect"`
// done(fs):
// done(fs): 	// Server controls if this agent acts like a Consul server,
// done(fs): 	// or merely as a client. Servers have more state, take part
// done(fs): 	// in leader election, etc.
// done(fs): 	Server bool `mapstructure:"server"`
// done(fs):
// done(fs): 	// (Enterprise-only) NonVotingServer is whether this server will act as a non-voting member
// done(fs): 	// of the cluster to help provide read scalability.
// done(fs): 	NonVotingServer bool `mapstructure:"non_voting_server"`
// done(fs):
// done(fs): 	// Datacenter is the datacenter this node is in. Defaults to dc1
// done(fs): 	Datacenter string `mapstructure:"datacenter"`
// done(fs):
// done(fs): 	// DataDir is the directory to store our state in
// done(fs): 	DataDir string `mapstructure:"data_dir"`
// done(fs):
// done(fs): 	// DNSRecursors can be set to allow the DNS servers to recursively
// done(fs): 	// resolve non-consul domains. It is deprecated, and merges into the
// done(fs): 	// recursors array.
// done(fs): 	DNSRecursor string `mapstructure:"recursor"`
// done(fs):
// done(fs): 	// DNSRecursors can be set to allow the DNS servers to recursively
// done(fs): 	// resolve non-consul domains
// done(fs): 	DNSRecursors []string `mapstructure:"recursors"`
// done(fs):
// done(fs): 	// DNS configuration
// done(fs): 	DNSConfig DNSConfig `mapstructure:"dns_config"`
// done(fs):
// done(fs): 	// Domain is the DNS domain for the records. Defaults to "consul."
// done(fs): 	Domain string `mapstructure:"domain"`
// done(fs):
// done(fs): 	// HTTP configuration
// done(fs): 	HTTPConfig HTTPConfig `mapstructure:"http_config"`
// done(fs):
// done(fs): 	// Encryption key to use for the Serf communication
// done(fs): 	EncryptKey string `mapstructure:"encrypt" json:"-"`
// done(fs):
// done(fs): 	// Disables writing the keyring to a file.
// done(fs): 	DisableKeyringFile bool `mapstructure:"disable_keyring_file"`
// done(fs):
// done(fs): 	// EncryptVerifyIncoming and EncryptVerifyOutgoing are used to enforce
// done(fs): 	// incoming/outgoing gossip encryption and can be used to upshift to
// done(fs): 	// encrypted gossip on a running cluster.
// done(fs): 	EncryptVerifyIncoming *bool `mapstructure:"encrypt_verify_incoming"`
// done(fs): 	EncryptVerifyOutgoing *bool `mapstructure:"encrypt_verify_outgoing"`
// done(fs):
// done(fs): 	// LogLevel is the level of the logs to putout
// done(fs): 	LogLevel string `mapstructure:"log_level"`
// done(fs):
// done(fs): 	// Node ID is a unique ID for this node across space and time. Defaults
// done(fs): 	// to a randomly-generated ID that persists in the data-dir.
// done(fs): 	NodeID types.NodeID `mapstructure:"node_id"`
// done(fs):
// done(fs): 	// DisableHostNodeID will prevent Consul from using information from the
// done(fs): 	// host to generate a node ID, and will cause Consul to generate a
// done(fs): 	// random ID instead.
// done(fs): 	DisableHostNodeID *bool `mapstructure:"disable_host_node_id"`
// done(fs):
// done(fs): 	// Node name is the name we use to advertise. Defaults to hostname.
// done(fs): 	NodeName string `mapstructure:"node_name"`
// done(fs):
// done(fs): 	// ClientAddr is used to control the address we bind to for
// done(fs): 	// client services (DNS, HTTP, HTTPS, RPC)
// done(fs): 	ClientAddr string `mapstructure:"client_addr"`
// done(fs):
// done(fs): 	// BindAddr is used to control the address we bind to.
// done(fs): 	// If not specified, the first private IP we find is used.
// done(fs): 	// This controls the address we use for cluster facing
// done(fs): 	// services (Gossip, Server RPC)
// done(fs): 	BindAddr string `mapstructure:"bind_addr"`
// done(fs):
// done(fs): 	// SerfWanBindAddr is used to control the address we bind to.
// done(fs): 	// If not specified, the first private IP we find is used.
// done(fs): 	// This controls the address we use for cluster facing
// done(fs): 	// services (Gossip) Serf
// done(fs): 	SerfWanBindAddr string `mapstructure:"serf_wan_bind"`
// done(fs):
// done(fs): 	// SerfLanBindAddr is used to control the address we bind to.
// done(fs): 	// If not specified, the first private IP we find is used.
// done(fs): 	// This controls the address we use for cluster facing
// done(fs): 	// services (Gossip) Serf
// done(fs): 	SerfLanBindAddr string `mapstructure:"serf_lan_bind"`
// done(fs):
// done(fs): 	// AdvertiseAddr is the address we use for advertising our Serf,
// done(fs): 	// and Consul RPC IP. If not specified, bind address is used.
// done(fs): 	AdvertiseAddr string `mapstructure:"advertise_addr"`
// done(fs):
// done(fs): 	// AdvertiseAddrs configuration
// done(fs): 	AdvertiseAddrs AdvertiseAddrsConfig `mapstructure:"advertise_addrs"`
// done(fs):
// done(fs): 	// AdvertiseAddrWan is the address we use for advertising our
// done(fs): 	// Serf WAN IP. If not specified, the general advertise address is used.
// done(fs): 	AdvertiseAddrWan string `mapstructure:"advertise_addr_wan"`
// done(fs):
// done(fs): 	// TranslateWanAddrs controls whether or not Consul should prefer
// done(fs): 	// the "wan" tagged address when doing lookups in remote datacenters.
// done(fs): 	// See TaggedAddresses below for more details.
// done(fs): 	TranslateWanAddrs bool `mapstructure:"translate_wan_addrs"`
// done(fs):
// done(fs): 	// Port configurations
// done(fs): 	Ports PortConfig
// done(fs):
// done(fs): 	// Address configurations
// done(fs): 	Addresses AddressConfig
// done(fs):
// done(fs): 	// (Enterprise-only) NetworkSegment is the network segment for this client to join.
// done(fs): 	Segment string `mapstructure:"segment"`
// done(fs):
// done(fs): 	// (Enterprise-only) Segments is the list of network segments for this server to
// done(fs): 	// initialize.
// done(fs): 	Segments []NetworkSegment `mapstructure:"segments"`
// done(fs):
// done(fs): 	// Tagged addresses. These are used to publish a set of addresses for
// done(fs): 	// for a node, which can be used by the remote agent. We currently
// done(fs): 	// populate only the "wan" tag based on the SerfWan advertise address,
// done(fs): 	// but this structure is here for possible future features with other
// done(fs): 	// user-defined tags. The "wan" tag will be used by remote agents if
// done(fs): 	// they are configured with TranslateWanAddrs set to true.
// done(fs): 	TaggedAddresses map[string]string
// done(fs):
// done(fs): 	// Node metadata key/value pairs. These are excluded from JSON output
// done(fs): 	// because they can be reloaded and might be stale when shown from the
// done(fs): 	// config instead of the local state.
// done(fs): 	Meta map[string]string `mapstructure:"node_meta" json:"-"`
// done(fs):
// done(fs): 	// LeaveOnTerm controls if Serf does a graceful leave when receiving
// done(fs): 	// the TERM signal. Defaults true on clients, false on servers. This can
// done(fs): 	// be changed on reload.
// done(fs): 	LeaveOnTerm *bool `mapstructure:"leave_on_terminate"`
// done(fs):
// done(fs): 	// SkipLeaveOnInt controls if Serf skips a graceful leave when
// done(fs): 	// receiving the INT signal. Defaults false on clients, true on
// done(fs): 	// servers. This can be changed on reload.
// done(fs): 	SkipLeaveOnInt *bool `mapstructure:"skip_leave_on_interrupt"`
// done(fs):
// done(fs): 	// Autopilot is used to configure helpful features for operating Consul servers.
// done(fs): 	Autopilot Autopilot `mapstructure:"autopilot"`
// done(fs):
// done(fs): 	Telemetry Telemetry `mapstructure:"telemetry"`
// done(fs):
// done(fs): 	// Protocol is the Consul protocol version to use.
// done(fs): 	Protocol int `mapstructure:"protocol"`
// done(fs):
// done(fs): 	// RaftProtocol sets the Raft protocol version to use on this server.
// done(fs): 	RaftProtocol int `mapstructure:"raft_protocol"`
// done(fs):
// done(fs): 	// EnableDebug is used to enable various debugging features
// done(fs): 	EnableDebug bool `mapstructure:"enable_debug"`
// done(fs):
// done(fs): 	// VerifyIncoming is used to verify the authenticity of incoming connections.
// done(fs): 	// This means that TCP requests are forbidden, only allowing for TLS. TLS connections
// done(fs): 	// must match a provided certificate authority. This can be used to force client auth.
// done(fs): 	VerifyIncoming bool `mapstructure:"verify_incoming"`
// done(fs):
// done(fs): 	// VerifyIncomingRPC is used to verify the authenticity of incoming RPC connections.
// done(fs): 	// This means that TCP requests are forbidden, only allowing for TLS. TLS connections
// done(fs): 	// must match a provided certificate authority. This can be used to force client auth.
// done(fs): 	VerifyIncomingRPC bool `mapstructure:"verify_incoming_rpc"`
// done(fs):
// done(fs): 	// VerifyIncomingHTTPS is used to verify the authenticity of incoming HTTPS connections.
// done(fs): 	// This means that TCP requests are forbidden, only allowing for TLS. TLS connections
// done(fs): 	// must match a provided certificate authority. This can be used to force client auth.
// done(fs): 	VerifyIncomingHTTPS bool `mapstructure:"verify_incoming_https"`
// done(fs):
// done(fs): 	// VerifyOutgoing is used to verify the authenticity of outgoing connections.
// done(fs): 	// This means that TLS requests are used. TLS connections must match a provided
// done(fs): 	// certificate authority. This is used to verify authenticity of server nodes.
// done(fs): 	VerifyOutgoing bool `mapstructure:"verify_outgoing"`
// done(fs):
// done(fs): 	// VerifyServerHostname is used to enable hostname verification of servers. This
// done(fs): 	// ensures that the certificate presented is valid for server.<datacenter>.<domain>.
// done(fs): 	// This prevents a compromised client from being restarted as a server, and then
// done(fs): 	// intercepting request traffic as well as being added as a raft peer. This should be
// done(fs): 	// enabled by default with VerifyOutgoing, but for legacy reasons we cannot break
// done(fs): 	// existing clients.
// done(fs): 	VerifyServerHostname bool `mapstructure:"verify_server_hostname"`
// done(fs):
// done(fs): 	// CAFile is a path to a certificate authority file. This is used with VerifyIncoming
// done(fs): 	// or VerifyOutgoing to verify the TLS connection.
// done(fs): 	CAFile string `mapstructure:"ca_file"`
// done(fs):
// done(fs): 	// CAPath is a path to a directory of certificate authority files. This is used with
// done(fs): 	// VerifyIncoming or VerifyOutgoing to verify the TLS connection.
// done(fs): 	CAPath string `mapstructure:"ca_path"`
// done(fs):
// done(fs): 	// CertFile is used to provide a TLS certificate that is used for serving TLS connections.
// done(fs): 	// Must be provided to serve TLS connections.
// done(fs): 	CertFile string `mapstructure:"cert_file"`
// done(fs):
// done(fs): 	// KeyFile is used to provide a TLS key that is used for serving TLS connections.
// done(fs): 	// Must be provided to serve TLS connections.
// done(fs): 	KeyFile string `mapstructure:"key_file"`
// done(fs):
// done(fs): 	// ServerName is used with the TLS certificates to ensure the name we
// done(fs): 	// provide matches the certificate
// done(fs): 	ServerName string `mapstructure:"server_name"`
// done(fs):
// done(fs): 	// TLSMinVersion is used to set the minimum TLS version used for TLS connections.
// done(fs): 	TLSMinVersion string `mapstructure:"tls_min_version"`
// done(fs):
// done(fs): 	// TLSCipherSuites is used to specify the list of supported ciphersuites.
// done(fs): 	TLSCipherSuites    []uint16 `mapstructure:"-" json:"-"`
// done(fs): 	TLSCipherSuitesRaw string   `mapstructure:"tls_cipher_suites"`
// done(fs):
// done(fs): 	// TLSPreferServerCipherSuites specifies whether to prefer the server's ciphersuite
// done(fs): 	// over the client ciphersuites.
// done(fs): 	TLSPreferServerCipherSuites bool `mapstructure:"tls_prefer_server_cipher_suites"`
// done(fs):
// done(fs): 	// StartJoin is a list of addresses to attempt to join when the
// done(fs): 	// agent starts. If Serf is unable to communicate with any of these
// done(fs): 	// addresses, then the agent will error and exit.
// done(fs): 	StartJoin []string `mapstructure:"start_join"`
// done(fs):
// done(fs): 	// StartJoinWan is a list of addresses to attempt to join -wan when the
// done(fs): 	// agent starts. If Serf is unable to communicate with any of these
// done(fs): 	// addresses, then the agent will error and exit.
// done(fs): 	StartJoinWan []string `mapstructure:"start_join_wan"`
// done(fs):
// done(fs): 	// RetryJoin is a list of addresses to join with retry enabled.
// done(fs): 	RetryJoin []string `mapstructure:"retry_join" json:"-"`
// done(fs):
// done(fs): 	// RetryMaxAttempts specifies the maximum number of times to retry joining a
// done(fs): 	// host on startup. This is useful for cases where we know the node will be
// done(fs): 	// online eventually.
// done(fs): 	RetryMaxAttempts int `mapstructure:"retry_max"`
// done(fs):
// done(fs): 	// RetryInterval specifies the amount of time to wait in between join
// done(fs): 	// attempts on agent start. The minimum allowed value is 1 second and
// done(fs): 	// the default is 30s.
// done(fs): 	RetryInterval    time.Duration `mapstructure:"-" json:"-"`
// done(fs): 	RetryIntervalRaw string        `mapstructure:"retry_interval"`
// done(fs):
// done(fs): 	// RetryJoinWan is a list of addresses to join -wan with retry enabled.
// done(fs): 	RetryJoinWan []string `mapstructure:"retry_join_wan"`
// done(fs):
// done(fs): 	// RetryMaxAttemptsWan specifies the maximum number of times to retry joining a
// done(fs): 	// -wan host on startup. This is useful for cases where we know the node will be
// done(fs): 	// online eventually.
// done(fs): 	RetryMaxAttemptsWan int `mapstructure:"retry_max_wan"`
// done(fs):
// done(fs): 	// RetryIntervalWan specifies the amount of time to wait in between join
// done(fs): 	// -wan attempts on agent start. The minimum allowed value is 1 second and
// done(fs): 	// the default is 30s.
// done(fs): 	RetryIntervalWan    time.Duration `mapstructure:"-" json:"-"`
// done(fs): 	RetryIntervalWanRaw string        `mapstructure:"retry_interval_wan"`
// done(fs):
// done(fs): 	// ReconnectTimeout* specify the amount of time to wait to reconnect with
// done(fs): 	// another agent before deciding it's permanently gone. This can be used to
// done(fs): 	// control the time it takes to reap failed nodes from the cluster.
// done(fs): 	ReconnectTimeoutLan    time.Duration `mapstructure:"-"`
// done(fs): 	ReconnectTimeoutLanRaw string        `mapstructure:"reconnect_timeout"`
// done(fs): 	ReconnectTimeoutWan    time.Duration `mapstructure:"-"`
// done(fs): 	ReconnectTimeoutWanRaw string        `mapstructure:"reconnect_timeout_wan"`
// done(fs):
// done(fs): 	// EnableUI enables the statically-compiled assets for the Consul web UI and
// done(fs): 	// serves them at the default /ui/ endpoint automatically.
// done(fs): 	EnableUI bool `mapstructure:"ui"`
// done(fs):
// done(fs): 	// UIDir is the directory containing the Web UI resources.
// done(fs): 	// If provided, the UI endpoints will be enabled.
// done(fs): 	UIDir string `mapstructure:"ui_dir"`
// done(fs):
// done(fs): 	// PidFile is the file to store our PID in
// done(fs): 	PidFile string `mapstructure:"pid_file"`
// done(fs):
// done(fs): 	// EnableSyslog is used to also tee all the logs over to syslog. Only supported
// done(fs): 	// on linux and OSX. Other platforms will generate an error.
// done(fs): 	EnableSyslog bool `mapstructure:"enable_syslog"`
// done(fs):
// done(fs): 	// SyslogFacility is used to control where the syslog messages go
// done(fs): 	// By default, goes to LOCAL0
// done(fs): 	SyslogFacility string `mapstructure:"syslog_facility"`
// done(fs):
// done(fs): 	// RejoinAfterLeave controls our interaction with the cluster after leave.
// done(fs): 	// When set to false (default), a leave causes Consul to not rejoin
// done(fs): 	// the cluster until an explicit join is received. If this is set to
// done(fs): 	// true, we ignore the leave, and rejoin the cluster on start.
// done(fs): 	RejoinAfterLeave bool `mapstructure:"rejoin_after_leave"`
// done(fs):
// done(fs): 	// EnableScriptChecks controls whether health checks which execute
// done(fs): 	// scripts are enabled. This includes regular script checks and Docker
// done(fs): 	// checks.
// done(fs): 	EnableScriptChecks bool `mapstructure:"enable_script_checks"`
// done(fs):
// done(fs): 	// CheckUpdateInterval controls the interval on which the output of a health check
// done(fs): 	// is updated if there is no change to the state. For example, a check in a steady
// done(fs): 	// state may run every 5 second generating a unique output (timestamp, etc), forcing
// done(fs): 	// constant writes. This allows Consul to defer the write for some period of time,
// done(fs): 	// reducing the write pressure when the state is steady.
// done(fs): 	CheckUpdateInterval    time.Duration `mapstructure:"-"`
// done(fs): 	CheckUpdateIntervalRaw string        `mapstructure:"check_update_interval" json:"-"`
// done(fs):
// done(fs): 	// CheckReapInterval controls the interval on which we will look for
// done(fs): 	// failed checks and reap their associated services, if so configured.
// done(fs): 	CheckReapInterval time.Duration `mapstructure:"-"`
// done(fs):
// done(fs): 	// CheckDeregisterIntervalMin is the smallest allowed interval to set
// done(fs): 	// a check's DeregisterCriticalServiceAfter value to.
// done(fs): 	CheckDeregisterIntervalMin time.Duration `mapstructure:"-"`
// done(fs):
// done(fs): 	// ACLToken is the default token used to make requests if a per-request
// done(fs): 	// token is not provided. If not configured the 'anonymous' token is used.
// done(fs): 	ACLToken string `mapstructure:"acl_token" json:"-"`
// done(fs):
// done(fs): 	// ACLAgentMasterToken is a special token that has full read and write
// done(fs): 	// privileges for this agent, and can be used to call agent endpoints
// done(fs): 	// when no servers are available.
// done(fs): 	ACLAgentMasterToken string `mapstructure:"acl_agent_master_token" json:"-"`
// done(fs):
// done(fs): 	// ACLAgentToken is the default token used to make requests for the agent
// done(fs): 	// itself, such as for registering itself with the catalog. If not
// done(fs): 	// configured, the 'acl_token' will be used.
// done(fs): 	ACLAgentToken string `mapstructure:"acl_agent_token" json:"-"`
// done(fs):
// done(fs): 	// ACLMasterToken is used to bootstrap the ACL system. It should be specified
// done(fs): 	// on the servers in the ACLDatacenter. When the leader comes online, it ensures
// done(fs): 	// that the Master token is available. This provides the initial token.
// done(fs): 	ACLMasterToken string `mapstructure:"acl_master_token" json:"-"`
// done(fs):
// done(fs): 	// ACLDatacenter is the central datacenter that holds authoritative
// done(fs): 	// ACL records. This must be the same for the entire cluster.
// done(fs): 	// If this is not set, ACLs are not enabled. Off by default.
// done(fs): 	ACLDatacenter string `mapstructure:"acl_datacenter"`
// done(fs):
// done(fs): 	// ACLTTL is used to control the time-to-live of cached ACLs . This has
// done(fs): 	// a major impact on performance. By default, it is set to 30 seconds.
// done(fs): 	ACLTTL    time.Duration `mapstructure:"-"`
// done(fs): 	ACLTTLRaw string        `mapstructure:"acl_ttl"`
// done(fs):
// done(fs): 	// ACLDefaultPolicy is used to control the ACL interaction when
// done(fs): 	// there is no defined policy. This can be "allow" which means
// done(fs): 	// ACLs are used to black-list, or "deny" which means ACLs are
// done(fs): 	// white-lists.
// done(fs): 	ACLDefaultPolicy string `mapstructure:"acl_default_policy"`
// done(fs):
// done(fs): 	// ACLDisabledTTL is used by clients to determine how long they will
// done(fs): 	// wait to check again with the servers if they discover ACLs are not
// done(fs): 	// enabled.
// done(fs): 	ACLDisabledTTL time.Duration `mapstructure:"-"`
// done(fs):
// done(fs): 	// ACLDownPolicy is used to control the ACL interaction when we cannot
// done(fs): 	// reach the ACLDatacenter and the token is not in the cache.
// done(fs): 	// There are two modes:
// done(fs): 	//   * allow - Allow all requests
// done(fs): 	//   * deny - Deny all requests
// done(fs): 	//   * extend-cache - Ignore the cache expiration, and allow cached
// done(fs): 	//                    ACL's to be used to service requests. This
// done(fs): 	//                    is the default. If the ACL is not in the cache,
// done(fs): 	//                    this acts like deny.
// done(fs): 	ACLDownPolicy string `mapstructure:"acl_down_policy"`
// done(fs):
// done(fs): 	// EnableACLReplication is used to turn on ACL replication when using
// done(fs): 	// /v1/agent/token/acl_replication_token to introduce the token, instead
// done(fs): 	// of setting acl_replication_token in the config. Setting the token via
// done(fs): 	// config will also set this to true for backward compatibility.
// done(fs): 	EnableACLReplication bool `mapstructure:"enable_acl_replication"`
// done(fs):
// done(fs): 	// ACLReplicationToken is used to fetch ACLs from the ACLDatacenter in
// done(fs): 	// order to replicate them locally. Setting this to a non-empty value
// done(fs): 	// also enables replication. Replication is only available in datacenters
// done(fs): 	// other than the ACLDatacenter.
// done(fs): 	ACLReplicationToken string `mapstructure:"acl_replication_token" json:"-"`
// done(fs):
// done(fs): 	// ACLEnforceVersion8 is used to gate a set of ACL policy features that
// done(fs): 	// are opt-in prior to Consul 0.8 and opt-out in Consul 0.8 and later.
// done(fs): 	ACLEnforceVersion8 *bool `mapstructure:"acl_enforce_version_8"`
// done(fs):
// done(fs): 	// Watches are used to monitor various endpoints and to invoke a
// done(fs): 	// handler to act appropriately. These are managed entirely in the
// done(fs): 	// agent layer using the standard APIs.
// done(fs): 	Watches []map[string]interface{} `mapstructure:"watches"`
// done(fs):
// done(fs): 	// DisableRemoteExec is used to turn off the remote execution
// done(fs): 	// feature. This is for security to prevent unknown scripts from running.
// done(fs): 	DisableRemoteExec *bool `mapstructure:"disable_remote_exec"`
// done(fs):
// done(fs): 	// DisableUpdateCheck is used to turn off the automatic update and
// done(fs): 	// security bulletin checking.
// done(fs): 	DisableUpdateCheck bool `mapstructure:"disable_update_check"`
// done(fs):
// done(fs): 	// DisableAnonymousSignature is used to turn off the anonymous signature
// done(fs): 	// send with the update check. This is used to deduplicate messages.
// done(fs): 	DisableAnonymousSignature bool `mapstructure:"disable_anonymous_signature"`
// done(fs):
// done(fs): 	// AEInterval controls the anti-entropy interval. This is how often
// done(fs): 	// the agent attempts to reconcile its local state with the server's
// done(fs): 	// representation of our state. Defaults to every 60s.
// done(fs): 	AEInterval time.Duration `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// DisableCoordinates controls features related to network coordinates.
// done(fs): 	DisableCoordinates bool `mapstructure:"disable_coordinates"`
// done(fs):
// done(fs): 	// SyncCoordinateRateTarget controls the rate for sending network
// done(fs): 	// coordinates to the server, in updates per second. This is the max rate
// done(fs): 	// that the server supports, so we scale our interval based on the size
// done(fs): 	// of the cluster to try to achieve this in aggregate at the server.
// done(fs): 	SyncCoordinateRateTarget float64 `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// SyncCoordinateIntervalMin sets the minimum interval that coordinates
// done(fs): 	// will be sent to the server. We scale the interval based on the cluster
// done(fs): 	// size, but below a certain interval it doesn't make sense send them any
// done(fs): 	// faster.
// done(fs): 	SyncCoordinateIntervalMin time.Duration `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// Checks holds the provided check definitions
// done(fs): 	Checks []*structs.CheckDefinition `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// Services holds the provided service definitions
// done(fs): 	Services []*structs.ServiceDefinition `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// ConsulConfig can either be provided or a default one created
// done(fs): 	ConsulConfig *consul.Config `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// Revision is the GitCommit this maps to
// done(fs): 	Revision string `mapstructure:"-"`
// done(fs):
// done(fs): 	// Version is the release version number
// done(fs): 	Version string `mapstructure:"-"`
// done(fs):
// done(fs): 	// VersionPrerelease is a label for pre-release builds
// done(fs): 	VersionPrerelease string `mapstructure:"-"`
// done(fs):
// done(fs): 	// WatchPlans contains the compiled watches
// done(fs): 	WatchPlans []*watch.Plan `mapstructure:"-" json:"-"`
// done(fs):
// done(fs): 	// UnixSockets is a map of socket configuration data
// done(fs): 	UnixSockets UnixSocketConfig `mapstructure:"unix_sockets"`
// done(fs):
// done(fs): 	// Minimum Session TTL
// done(fs): 	SessionTTLMin    time.Duration `mapstructure:"-"`
// done(fs): 	SessionTTLMinRaw string        `mapstructure:"session_ttl_min"`
// done(fs):
// done(fs): 	// deprecated fields
// done(fs): 	// keep them exported since otherwise the error messages don't show up
// done(fs): 	DeprecatedAtlasInfrastructure    string            `mapstructure:"atlas_infrastructure" json:"-"`
// done(fs): 	DeprecatedAtlasToken             string            `mapstructure:"atlas_token" json:"-"`
// done(fs): 	DeprecatedAtlasACLToken          string            `mapstructure:"atlas_acl_token" json:"-"`
// done(fs): 	DeprecatedAtlasJoin              bool              `mapstructure:"atlas_join" json:"-"`
// done(fs): 	DeprecatedAtlasEndpoint          string            `mapstructure:"atlas_endpoint" json:"-"`
// done(fs): 	DeprecatedHTTPAPIResponseHeaders map[string]string `mapstructure:"http_api_response_headers"`
// done(fs): 	DeprecatedRetryJoinEC2           RetryJoinEC2      `mapstructure:"retry_join_ec2"`
// done(fs): 	DeprecatedRetryJoinGCE           RetryJoinGCE      `mapstructure:"retry_join_gce"`
// done(fs): 	DeprecatedRetryJoinAzure         RetryJoinAzure    `mapstructure:"retry_join_azure"`
// done(fs): }
// done(fs):
// done(fs): // IncomingHTTPSConfig returns the TLS configuration for HTTPS
// done(fs): // connections to consul.
// done(fs): func (c *Config) IncomingHTTPSConfig() (*tls.Config, error) {
// done(fs): 	tc := &tlsutil.Config{
// done(fs): 		VerifyIncoming:           c.VerifyIncoming || c.VerifyIncomingHTTPS,
// done(fs): 		VerifyOutgoing:           c.VerifyOutgoing,
// done(fs): 		CAFile:                   c.CAFile,
// done(fs): 		CAPath:                   c.CAPath,
// done(fs): 		CertFile:                 c.CertFile,
// done(fs): 		KeyFile:                  c.KeyFile,
// done(fs): 		NodeName:                 c.NodeName,
// done(fs): 		ServerName:               c.ServerName,
// done(fs): 		TLSMinVersion:            c.TLSMinVersion,
// done(fs): 		CipherSuites:             c.TLSCipherSuites,
// done(fs): 		PreferServerCipherSuites: c.TLSPreferServerCipherSuites,
// done(fs): 	}
// done(fs): 	return tc.IncomingTLSConfig()
// done(fs): }
// done(fs):
// done(fs): type ProtoAddr struct {
// done(fs): 	Proto, Net, Addr string
// done(fs): }
// done(fs):
// done(fs): func (p ProtoAddr) String() string {
// done(fs): 	return p.Proto + "://" + p.Addr
// done(fs): }
// done(fs):
// done(fs): func (c *Config) DNSAddrs() ([]ProtoAddr, error) {
// done(fs): 	if c.Ports.DNS <= 0 {
// done(fs): 		return nil, nil
// done(fs): 	}
// done(fs): 	a, err := c.ClientListener(c.Addresses.DNS, c.Ports.DNS)
// done(fs): 	if err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs): 	addrs := []ProtoAddr{
// done(fs): 		{"dns", "tcp", a.String()},
// done(fs): 		{"dns", "udp", a.String()},
// done(fs): 	}
// done(fs): 	return addrs, nil
// done(fs): }
// done(fs):
// done(fs): // HTTPAddrs returns the bind addresses for the HTTP server and
// done(fs): // the application protocol which should be served, e.g. 'http'
// done(fs): // or 'https'.
// done(fs): func (c *Config) HTTPAddrs() ([]ProtoAddr, error) {
// done(fs): 	var addrs []ProtoAddr
// done(fs): 	if c.Ports.HTTP > 0 {
// done(fs): 		a, err := c.ClientListener(c.Addresses.HTTP, c.Ports.HTTP)
// done(fs): 		if err != nil {
// done(fs): 			return nil, err
// done(fs): 		}
// done(fs): 		addrs = append(addrs, ProtoAddr{"http", a.Network(), a.String()})
// done(fs): 	}
// done(fs): 	if c.Ports.HTTPS > 0 && c.CertFile != "" && c.KeyFile != "" {
// done(fs): 		a, err := c.ClientListener(c.Addresses.HTTPS, c.Ports.HTTPS)
// done(fs): 		if err != nil {
// done(fs): 			return nil, err
// done(fs): 		}
// done(fs): 		addrs = append(addrs, ProtoAddr{"https", a.Network(), a.String()})
// done(fs): 	}
// done(fs): 	return addrs, nil
// done(fs): }
// done(fs):
// done(fs): // Bool is used to initialize bool pointers in struct literals.
// done(fs): func Bool(b bool) *bool {
// done(fs): 	return &b
// done(fs): }
// done(fs):
// done(fs): // Uint64 is used to initialize uint64 pointers in struct literals.
// done(fs): func Uint64(i uint64) *uint64 {
// done(fs): 	return &i
// done(fs): }
// done(fs):
// done(fs): // Duration is used to initialize time.Duration pointers in struct literals.
// done(fs): func Duration(d time.Duration) *time.Duration {
// done(fs): 	return &d
// done(fs): }
// done(fs):
// done(fs): // UnixSocketPermissions contains information about a unix socket, and
// done(fs): // implements the FilePermissions interface.
// done(fs): type UnixSocketPermissions struct {
// done(fs): 	Usr   string `mapstructure:"user"`
// done(fs): 	Grp   string `mapstructure:"group"`
// done(fs): 	Perms string `mapstructure:"mode"`
// done(fs): }
// done(fs):
// done(fs): func (u UnixSocketPermissions) User() string {
// done(fs): 	return u.Usr
// done(fs): }
// done(fs):
// done(fs): func (u UnixSocketPermissions) Group() string {
// done(fs): 	return u.Grp
// done(fs): }
// done(fs):
// done(fs): func (u UnixSocketPermissions) Mode() string {
// done(fs): 	return u.Perms
// done(fs): }
// done(fs):
// done(fs): func (s *Telemetry) GoString() string {
// done(fs): 	return fmt.Sprintf("*%#v", *s)
// done(fs): }
// done(fs):
// done(fs): // UnixSocketConfig stores information about various unix sockets which
// done(fs): // Consul creates and uses for communication.
// done(fs): type UnixSocketConfig struct {
// done(fs): 	UnixSocketPermissions `mapstructure:",squash"`
// done(fs): }
// done(fs):
// socketPath tests if a given address describes a domain socket,
// and returns the relevant path part of the string if it is.
func socketPath(addr string) string {
	if !strings.HasPrefix(addr, "unix://") {
		return ""
	}
	return strings.TrimPrefix(addr, "unix://")
}

// done(fs):
// done(fs): type dirEnts []os.FileInfo
// done(fs):
// done(fs): // DefaultConfig is used to return a sane default configuration
// done(fs): func DefaultConfig() *Config {
// done(fs): 	return &Config{
// done(fs): 		Limits: Limits{
// done(fs): 			RPCRate:     rate.Inf,
// done(fs): 			RPCMaxBurst: 1000,
// done(fs): 		},
// done(fs): 		Bootstrap:       false,
// done(fs): 		BootstrapExpect: 0,
// done(fs): 		Server:          false,
// done(fs): 		Datacenter:      consul.DefaultDC,
// done(fs): 		Domain:          "consul.",
// done(fs): 		LogLevel:        "INFO",
// done(fs): 		ClientAddr:      "127.0.0.1",
// done(fs): 		BindAddr:        "0.0.0.0",
// done(fs): 		Ports: PortConfig{
// done(fs): 			DNS:     8600,
// done(fs): 			HTTP:    8500,
// done(fs): 			HTTPS:   -1,
// done(fs): 			SerfLan: consul.DefaultLANSerfPort,
// done(fs): 			SerfWan: consul.DefaultWANSerfPort,
// done(fs): 			Server:  8300,
// done(fs): 		},
// done(fs): 		DNSConfig: DNSConfig{
// done(fs): 			AllowStale:      Bool(true),
// done(fs): 			UDPAnswerLimit:  3,
// done(fs): 			MaxStale:        10 * 365 * 24 * time.Hour,
// done(fs): 			RecursorTimeout: 2 * time.Second,
// done(fs): 		},
// done(fs): 		Telemetry: Telemetry{
// done(fs): 			StatsitePrefix: "consul",
// done(fs): 			FilterDefault:  Bool(true),
// done(fs): 		},
// done(fs): 		Meta:                       make(map[string]string),
// done(fs): 		SyslogFacility:             "LOCAL0",
// done(fs): 		Protocol:                   consul.ProtocolVersion2Compatible,
// done(fs): 		CheckUpdateInterval:        5 * time.Minute,
// done(fs): 		CheckDeregisterIntervalMin: time.Minute,
// done(fs): 		CheckReapInterval:          30 * time.Second,
// done(fs): 		AEInterval:                 time.Minute,
// done(fs): 		DisableCoordinates:         false,
// done(fs):
// done(fs): 		// SyncCoordinateRateTarget is set based on the rate that we want
// done(fs): 		// the server to handle as an aggregate across the entire cluster.
// done(fs): 		// If you update this, you'll need to adjust CoordinateUpdate* in
// done(fs): 		// the server-side config accordingly.
// done(fs): 		SyncCoordinateRateTarget:  64.0, // updates / second
// done(fs): 		SyncCoordinateIntervalMin: 15 * time.Second,
// done(fs):
// done(fs): 		ACLTTL:             30 * time.Second,
// done(fs): 		ACLDownPolicy:      "extend-cache",
// done(fs): 		ACLDefaultPolicy:   "allow",
// done(fs): 		ACLDisabledTTL:     120 * time.Second,
// done(fs): 		ACLEnforceVersion8: Bool(true),
// done(fs): 		DisableRemoteExec:  Bool(true),
// done(fs): 		RetryInterval:      30 * time.Second,
// done(fs): 		RetryIntervalWan:   30 * time.Second,
// done(fs):
// done(fs): 		TLSMinVersion: "tls10",
// done(fs):
// done(fs): 		EncryptVerifyIncoming: Bool(true),
// done(fs): 		EncryptVerifyOutgoing: Bool(true),
// done(fs):
// done(fs): 		DisableHostNodeID: Bool(true),
// done(fs): 	}
// done(fs): }
// done(fs):
// todo(fs): // DevConfig is used to return a set of configuration to use for dev mode.
// todo(fs): func DevConfig() *Config {
// todo(fs): 	conf := DefaultConfig()
// todo(fs): 	conf.DevMode = true
// todo(fs): 	conf.LogLevel = "DEBUG"
// todo(fs): 	conf.Server = true
// todo(fs): 	conf.EnableDebug = true
// todo(fs): 	conf.DisableAnonymousSignature = true
// todo(fs): 	conf.EnableUI = true
// todo(fs): 	conf.BindAddr = "127.0.0.1"
// todo(fs): 	conf.DisableKeyringFile = true
// todo(fs):
// todo(fs): 	conf.ConsulConfig = consul.DefaultConfig()
// todo(fs): 	conf.ConsulConfig.SerfLANConfig.MemberlistConfig.ProbeTimeout = 100 * time.Millisecond
// todo(fs): 	conf.ConsulConfig.SerfLANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
// todo(fs): 	conf.ConsulConfig.SerfLANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond
// todo(fs):
// todo(fs): 	conf.ConsulConfig.SerfWANConfig.MemberlistConfig.SuspicionMult = 3
// todo(fs): 	conf.ConsulConfig.SerfWANConfig.MemberlistConfig.ProbeTimeout = 100 * time.Millisecond
// todo(fs): 	conf.ConsulConfig.SerfWANConfig.MemberlistConfig.ProbeInterval = 100 * time.Millisecond
// todo(fs): 	conf.ConsulConfig.SerfWANConfig.MemberlistConfig.GossipInterval = 100 * time.Millisecond
// todo(fs):
// todo(fs): 	conf.ConsulConfig.RaftConfig.LeaderLeaseTimeout = 20 * time.Millisecond
// todo(fs): 	conf.ConsulConfig.RaftConfig.HeartbeatTimeout = 40 * time.Millisecond
// todo(fs): 	conf.ConsulConfig.RaftConfig.ElectionTimeout = 40 * time.Millisecond
// todo(fs):
// todo(fs): 	conf.ConsulConfig.CoordinateUpdatePeriod = 100 * time.Millisecond
// todo(fs):
// todo(fs): 	return conf
// todo(fs): }
// todo(fs):
// done(fs): // EncryptBytes returns the encryption key configured.
// done(fs): func (c *Config) EncryptBytes() ([]byte, error) {
// done(fs): 	return base64.StdEncoding.DecodeString(c.EncryptKey)
// done(fs): }
// done(fs):
// done(fs): // ClientListener is used to format a listener for a
// done(fs): // port on a ClientAddr
// done(fs): func (c *Config) ClientListener(override string, port int) (net.Addr, error) {
// done(fs): 	addr := c.ClientAddr
// done(fs): 	if override != "" {
// done(fs): 		addr = override
// done(fs): 	}
// done(fs): 	if path := socketPath(addr); path != "" {
// done(fs): 		return &net.UnixAddr{Name: path, Net: "unix"}, nil
// done(fs): 	}
// done(fs): 	ip := net.ParseIP(addr)
// done(fs): 	if ip == nil {
// done(fs): 		return nil, fmt.Errorf("Failed to parse IP: %v", addr)
// done(fs): 	}
// done(fs): 	return &net.TCPAddr{IP: ip, Port: port}, nil
// done(fs): }
// done(fs):
// done(fs): // VerifyUniqueListeners checks to see if an address was used more than once in
// done(fs): // the config
// done(fs): func (c *Config) VerifyUniqueListeners() error {
// done(fs): 	listeners := []struct {
// done(fs): 		host  string
// done(fs): 		port  int
// done(fs): 		descr string
// done(fs): 	}{
// done(fs): 		{c.Addresses.DNS, c.Ports.DNS, "DNS"},
// done(fs): 		{c.Addresses.HTTP, c.Ports.HTTP, "HTTP"},
// done(fs): 		{c.Addresses.HTTPS, c.Ports.HTTPS, "HTTPS"},
// done(fs): 		{c.AdvertiseAddr, c.Ports.Server, "Server RPC"},
// done(fs): 		{c.AdvertiseAddr, c.Ports.SerfLan, "Serf LAN"},
// done(fs): 		{c.AdvertiseAddr, c.Ports.SerfWan, "Serf WAN"},
// done(fs): 	}
// done(fs):
// done(fs): 	type key struct {
// done(fs): 		host string
// done(fs): 		port int
// done(fs): 	}
// done(fs): 	m := make(map[key]string, len(listeners))
// done(fs):
// done(fs): 	for _, l := range listeners {
// done(fs): 		if l.host == "" {
// done(fs): 			l.host = "0.0.0.0"
// done(fs): 		} else if strings.HasPrefix(l.host, "unix") {
// done(fs): 			// Don't compare ports on unix sockets
// done(fs): 			l.port = 0
// done(fs): 		}
// done(fs): 		if l.host == "0.0.0.0" && l.port <= 0 {
// done(fs): 			continue
// done(fs): 		}
// done(fs):
// done(fs): 		k := key{l.host, l.port}
// done(fs): 		v, ok := m[k]
// done(fs): 		if ok {
// done(fs): 			return fmt.Errorf("%s address already configured for %s", l.descr, v)
// done(fs): 		}
// done(fs): 		m[k] = l.descr
// done(fs): 	}
// done(fs): 	return nil
// done(fs): }
// done(fs):
// done(fs): // DecodeConfig reads the configuration from the given reader in JSON
// done(fs): // format and decodes it into a proper Config structure.
// done(fs): func DecodeConfig(r io.Reader) (*Config, error) {
// done(fs): 	var raw interface{}
// done(fs): 	if err := json.NewDecoder(r).Decode(&raw); err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs):
// done(fs): 	// Check the result type
// done(fs): 	var result Config
// done(fs): 	if obj, ok := raw.(map[string]interface{}); ok {
// done(fs): 		// Check for a "services", "service" or "check" key, meaning
// done(fs): 		// this is actually a definition entry
// done(fs): 		if sub, ok := obj["services"]; ok {
// done(fs): 			if list, ok := sub.([]interface{}); ok {
// done(fs): 				for _, srv := range list {
// done(fs): 					service, err := DecodeServiceDefinition(srv)
// done(fs): 					if err != nil {
// done(fs): 						return nil, err
// done(fs): 					}
// done(fs): 					result.Services = append(result.Services, service)
// done(fs): 				}
// done(fs): 			}
// done(fs): 		}
// done(fs): 		if sub, ok := obj["service"]; ok {
// done(fs): 			service, err := DecodeServiceDefinition(sub)
// done(fs): 			if err != nil {
// done(fs): 				return nil, err
// done(fs): 			}
// done(fs): 			result.Services = append(result.Services, service)
// done(fs): 		}
// done(fs): 		if sub, ok := obj["checks"]; ok {
// done(fs): 			if list, ok := sub.([]interface{}); ok {
// done(fs): 				for _, chk := range list {
// done(fs): 					check, err := DecodeCheckDefinition(chk)
// done(fs): 					if err != nil {
// done(fs): 						return nil, err
// done(fs): 					}
// done(fs): 					result.Checks = append(result.Checks, check)
// done(fs): 				}
// done(fs): 			}
// done(fs): 		}
// done(fs): 		if sub, ok := obj["check"]; ok {
// done(fs): 			check, err := DecodeCheckDefinition(sub)
// done(fs): 			if err != nil {
// done(fs): 				return nil, err
// done(fs): 			}
// done(fs): 			result.Checks = append(result.Checks, check)
// done(fs): 		}
// done(fs):
// done(fs): 		// A little hacky but upgrades the old stats config directives to the new way
// done(fs): 		if sub, ok := obj["statsd_addr"]; ok && result.Telemetry.StatsdAddr == "" {
// done(fs): 			result.Telemetry.StatsdAddr = sub.(string)
// done(fs): 		}
// done(fs):
// done(fs): 		if sub, ok := obj["statsite_addr"]; ok && result.Telemetry.StatsiteAddr == "" {
// done(fs): 			result.Telemetry.StatsiteAddr = sub.(string)
// done(fs): 		}
// done(fs):
// done(fs): 		if sub, ok := obj["statsite_prefix"]; ok && result.Telemetry.StatsitePrefix == "" {
// done(fs): 			result.Telemetry.StatsitePrefix = sub.(string)
// done(fs): 		}
// done(fs):
// done(fs): 		if sub, ok := obj["dogstatsd_addr"]; ok && result.Telemetry.DogStatsdAddr == "" {
// done(fs): 			result.Telemetry.DogStatsdAddr = sub.(string)
// done(fs): 		}
// done(fs):
// done(fs): 		if sub, ok := obj["dogstatsd_tags"].([]interface{}); ok && len(result.Telemetry.DogStatsdTags) == 0 {
// done(fs): 			result.Telemetry.DogStatsdTags = make([]string, len(sub))
// done(fs): 			for i := range sub {
// done(fs): 				result.Telemetry.DogStatsdTags[i] = sub[i].(string)
// done(fs): 			}
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	// Decode
// done(fs): 	var md mapstructure.Metadata
// done(fs): 	msdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
// done(fs): 		Metadata: &md,
// done(fs): 		Result:   &result,
// done(fs): 	})
// done(fs): 	if err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs):
// done(fs): 	if err := msdec.Decode(raw); err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs):
// done(fs): 	// Check for deprecations
// done(fs): 	if result.Ports.RPC != 0 {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: ports.rpc is deprecated and is "+
// done(fs): 			"no longer used. Please remove it from your configuration.")
// done(fs): 	}
// done(fs): 	if result.Addresses.RPC != "" {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: addresses.rpc is deprecated and "+
// done(fs): 			"is no longer used. Please remove it from your configuration.")
// done(fs): 	}
// done(fs): 	if result.DeprecatedAtlasInfrastructure != "" {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: atlas_infrastructure is deprecated and "+
// done(fs): 			"is no longer used. Please remove it from your configuration.")
// done(fs): 	}
// done(fs): 	if result.DeprecatedAtlasToken != "" {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: atlas_token is deprecated and "+
// done(fs): 			"is no longer used. Please remove it from your configuration.")
// done(fs): 	}
// done(fs): 	if result.DeprecatedAtlasACLToken != "" {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: atlas_acl_token is deprecated and "+
// done(fs): 			"is no longer used. Please remove it from your configuration.")
// done(fs): 	}
// done(fs): 	if result.DeprecatedAtlasJoin != false {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: atlas_join is deprecated and "+
// done(fs): 			"is no longer used. Please remove it from your configuration.")
// done(fs): 	}
// done(fs): 	if result.DeprecatedAtlasEndpoint != "" {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: atlas_endpoint is deprecated and "+
// done(fs): 			"is no longer used. Please remove it from your configuration.")
// done(fs): 	}
// done(fs):
// done(fs): 	// Check unused fields and verify that no bad configuration options were
// done(fs): 	// passed to Consul. There are a few additional fields which don't directly
// done(fs): 	// use mapstructure decoding, so we need to account for those as well. These
// done(fs): 	// telemetry-related fields used to be available as top-level keys, so they
// done(fs): 	// are here for backward compatibility with the old format.
// done(fs): 	allowedKeys := []string{
// done(fs): 		"service", "services", "check", "checks", "statsd_addr", "statsite_addr", "statsite_prefix",
// done(fs): 		"dogstatsd_addr", "dogstatsd_tags",
// done(fs): 	}
// done(fs):
// done(fs): 	var unused []string
// done(fs): 	for _, field := range md.Unused {
// done(fs): 		if !lib.StrContains(allowedKeys, field) {
// done(fs): 			unused = append(unused, field)
// done(fs): 		}
// done(fs): 	}
// done(fs): 	if len(unused) > 0 {
// done(fs): 		return nil, fmt.Errorf("Config has invalid keys: %s", strings.Join(unused, ","))
// done(fs): 	}
// done(fs):
// done(fs): 	// Handle time conversions
// done(fs): 	if raw := result.DNSConfig.NodeTTLRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("NodeTTL invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.DNSConfig.NodeTTL = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.DNSConfig.MaxStaleRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("MaxStale invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.DNSConfig.MaxStale = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.DNSConfig.RecursorTimeoutRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("RecursorTimeout invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.DNSConfig.RecursorTimeout = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if len(result.DNSConfig.ServiceTTLRaw) != 0 {
// done(fs): 		if result.DNSConfig.ServiceTTL == nil {
// done(fs): 			result.DNSConfig.ServiceTTL = make(map[string]time.Duration)
// done(fs): 		}
// done(fs): 		for service, raw := range result.DNSConfig.ServiceTTLRaw {
// done(fs): 			dur, err := time.ParseDuration(raw)
// done(fs): 			if err != nil {
// done(fs): 				return nil, fmt.Errorf("ServiceTTL %s invalid: %v", service, err)
// done(fs): 			}
// done(fs): 			result.DNSConfig.ServiceTTL[service] = dur
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.CheckUpdateIntervalRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("CheckUpdateInterval invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.CheckUpdateInterval = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.ACLTTLRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("ACL TTL invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.ACLTTL = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.RetryIntervalRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("RetryInterval invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.RetryInterval = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.RetryIntervalWanRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("RetryIntervalWan invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.RetryIntervalWan = dur
// done(fs): 	}
// done(fs):
// done(fs): 	const reconnectTimeoutMin = 8 * time.Hour
// done(fs): 	if raw := result.ReconnectTimeoutLanRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("ReconnectTimeoutLan invalid: %v", err)
// done(fs): 		}
// done(fs): 		if dur < reconnectTimeoutMin {
// done(fs): 			return nil, fmt.Errorf("ReconnectTimeoutLan must be >= %s", reconnectTimeoutMin.String())
// done(fs): 		}
// done(fs): 		result.ReconnectTimeoutLan = dur
// done(fs): 	}
// done(fs): 	if raw := result.ReconnectTimeoutWanRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("ReconnectTimeoutWan invalid: %v", err)
// done(fs): 		}
// done(fs): 		if dur < reconnectTimeoutMin {
// done(fs): 			return nil, fmt.Errorf("ReconnectTimeoutWan must be >= %s", reconnectTimeoutMin.String())
// done(fs): 		}
// done(fs): 		result.ReconnectTimeoutWan = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.Autopilot.LastContactThresholdRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("LastContactThreshold invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.Autopilot.LastContactThreshold = &dur
// done(fs): 	}
// done(fs): 	if raw := result.Autopilot.ServerStabilizationTimeRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("ServerStabilizationTime invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.Autopilot.ServerStabilizationTime = &dur
// done(fs): 	}
// done(fs):
// done(fs): 	// Merge the single recursor
// done(fs): 	if result.DNSRecursor != "" {
// done(fs): 		result.DNSRecursors = append(result.DNSRecursors, result.DNSRecursor)
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.SessionTTLMinRaw; raw != "" {
// done(fs): 		dur, err := time.ParseDuration(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("Session TTL Min invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.SessionTTLMin = dur
// done(fs): 	}
// done(fs):
// done(fs): 	if result.AdvertiseAddrs.SerfLanRaw != "" {
// done(fs): 		ipStr, err := parseSingleIPTemplate(result.AdvertiseAddrs.SerfLanRaw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("Serf Advertise LAN address resolution failed: %v", err)
// done(fs): 		}
// done(fs): 		result.AdvertiseAddrs.SerfLanRaw = ipStr
// done(fs):
// done(fs): 		addr, err := net.ResolveTCPAddr("tcp", result.AdvertiseAddrs.SerfLanRaw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("AdvertiseAddrs.SerfLan is invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.AdvertiseAddrs.SerfLan = addr
// done(fs): 	}
// done(fs):
// done(fs): 	if result.AdvertiseAddrs.SerfWanRaw != "" {
// done(fs): 		ipStr, err := parseSingleIPTemplate(result.AdvertiseAddrs.SerfWanRaw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("Serf Advertise WAN address resolution failed: %v", err)
// done(fs): 		}
// done(fs): 		result.AdvertiseAddrs.SerfWanRaw = ipStr
// done(fs):
// done(fs): 		addr, err := net.ResolveTCPAddr("tcp", result.AdvertiseAddrs.SerfWanRaw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("AdvertiseAddrs.SerfWan is invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.AdvertiseAddrs.SerfWan = addr
// done(fs): 	}
// done(fs):
// done(fs): 	if result.AdvertiseAddrs.RPCRaw != "" {
// done(fs): 		ipStr, err := parseSingleIPTemplate(result.AdvertiseAddrs.RPCRaw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("RPC Advertise address resolution failed: %v", err)
// done(fs): 		}
// done(fs): 		result.AdvertiseAddrs.RPCRaw = ipStr
// done(fs):
// done(fs): 		addr, err := net.ResolveTCPAddr("tcp", result.AdvertiseAddrs.RPCRaw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("AdvertiseAddrs.RPC is invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.AdvertiseAddrs.RPC = addr
// done(fs): 	}
// done(fs):
// done(fs): 	// Validate segment config.
// done(fs): 	if err := ValidateSegments(&result); err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs):
// done(fs): 	// Enforce the max Raft multiplier.
// done(fs): 	if result.Performance.RaftMultiplier > consul.MaxRaftMultiplier {
// done(fs): 		return nil, fmt.Errorf("Performance.RaftMultiplier must be <= %d", consul.MaxRaftMultiplier)
// done(fs): 	}
// done(fs):
// done(fs): 	if raw := result.TLSCipherSuitesRaw; raw != "" {
// done(fs): 		ciphers, err := tlsutil.ParseCiphers(raw)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("TLSCipherSuites invalid: %v", err)
// done(fs): 		}
// done(fs): 		result.TLSCipherSuites = ciphers
// done(fs): 	}
// done(fs):
// done(fs): 	// This is for backwards compatibility.
// done(fs): 	// HTTPAPIResponseHeaders has been replaced with HTTPConfig.ResponseHeaders
// done(fs): 	if len(result.DeprecatedHTTPAPIResponseHeaders) > 0 {
// done(fs): 		fmt.Fprintln(os.Stderr, "==> DEPRECATION: http_api_response_headers is deprecated and "+
// done(fs): 			"is no longer used. Please use http_config.response_headers instead.")
// done(fs): 		if result.HTTPConfig.ResponseHeaders == nil {
// done(fs): 			result.HTTPConfig.ResponseHeaders = make(map[string]string)
// done(fs): 		}
// done(fs): 		for field, value := range result.DeprecatedHTTPAPIResponseHeaders {
// done(fs): 			result.HTTPConfig.ResponseHeaders[field] = value
// done(fs): 		}
// done(fs): 		result.DeprecatedHTTPAPIResponseHeaders = nil
// done(fs): 	}
// done(fs):
// done(fs): 	// Set the ACL replication enable if they set a token, for backwards
// done(fs): 	// compatibility.
// done(fs): 	if result.ACLReplicationToken != "" {
// done(fs): 		result.EnableACLReplication = true
// done(fs): 	}
// done(fs):
// done(fs): 	// Parse the metric filters
// done(fs): 	for _, rule := range result.Telemetry.PrefixFilter {
// done(fs): 		if rule == "" {
// done(fs): 			return nil, fmt.Errorf("Cannot have empty filter rule in prefix_filter")
// done(fs): 		}
// done(fs): 		switch rule[0] {
// done(fs): 		case '+':
// done(fs): 			result.Telemetry.AllowedPrefixes = append(result.Telemetry.AllowedPrefixes, rule[1:])
// done(fs): 		case '-':
// done(fs): 			result.Telemetry.BlockedPrefixes = append(result.Telemetry.BlockedPrefixes, rule[1:])
// done(fs): 		default:
// done(fs): 			return nil, fmt.Errorf("Filter rule must begin with either '+' or '-': %q", rule)
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	// Validate node meta fields
// done(fs): 	if err := structs.ValidateMetadata(result.Meta, false); err != nil {
// done(fs): 		return nil, fmt.Errorf("Failed to parse node metadata: %v", err)
// done(fs): 	}
// done(fs):
// done(fs): 	return &result, nil
// done(fs): }
// done(fs):
// done(fs): // DecodeServiceDefinition is used to decode a service definition
// done(fs): func DecodeServiceDefinition(raw interface{}) (*structs.ServiceDefinition, error) {
// done(fs): 	rawMap, ok := raw.(map[string]interface{})
// done(fs): 	if !ok {
// done(fs): 		goto AFTER_FIX
// done(fs): 	}
// done(fs):
// done(fs): 	// If no 'tags', handle the deprecated 'tag' value.
// done(fs): 	if _, ok := rawMap["tags"]; !ok {
// done(fs): 		if tag, ok := rawMap["tag"]; ok {
// done(fs): 			rawMap["tags"] = []interface{}{tag}
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	for k, v := range rawMap {
// done(fs): 		switch strings.ToLower(k) {
// done(fs): 		case "check":
// done(fs): 			if err := FixupCheckType(v); err != nil {
// done(fs): 				return nil, err
// done(fs): 			}
// done(fs): 		case "checks":
// done(fs): 			chkTypes, ok := v.([]interface{})
// done(fs): 			if !ok {
// done(fs): 				goto AFTER_FIX
// done(fs): 			}
// done(fs): 			for _, chkType := range chkTypes {
// done(fs): 				if err := FixupCheckType(chkType); err != nil {
// done(fs): 					return nil, err
// done(fs): 				}
// done(fs): 			}
// done(fs): 		}
// done(fs): 	}
// done(fs): AFTER_FIX:
// done(fs): 	var md mapstructure.Metadata
// done(fs): 	var result structs.ServiceDefinition
// done(fs): 	msdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
// done(fs): 		Metadata: &md,
// done(fs): 		Result:   &result,
// done(fs): 	})
// done(fs): 	if err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs): 	if err := msdec.Decode(raw); err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs): 	return &result, nil
// done(fs): }

var errInvalidHeaderFormat = errors.New("agent: invalid format of 'header' field")

func FixupCheckType(raw interface{}) error {
	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}

	parseDuration := func(v interface{}) (time.Duration, error) {
		if v == nil {
			return 0, nil
		}
		switch x := v.(type) {
		case time.Duration:
			return x, nil
		case float64:
			return time.Duration(x), nil
		case string:
			return time.ParseDuration(x)
		default:
			return 0, fmt.Errorf("invalid format")
		}
	}

	parseHeaderMap := func(v interface{}) (map[string][]string, error) {
		if v == nil {
			return nil, nil
		}
		vm, ok := v.(map[string]interface{})
		if !ok {
			return nil, errInvalidHeaderFormat
		}
		m := map[string][]string{}
		for k, vv := range vm {
			vs, ok := vv.([]interface{})
			if !ok {
				return nil, errInvalidHeaderFormat
			}
			for _, vs := range vs {
				s, ok := vs.(string)
				if !ok {
					return nil, errInvalidHeaderFormat
				}
				m[k] = append(m[k], s)
			}
		}
		return m, nil
	}

	replace := func(oldKey, newKey string, val interface{}) {
		rawMap[newKey] = val
		if oldKey != newKey {
			delete(rawMap, oldKey)
		}
	}

	for k, v := range rawMap {
		switch strings.ToLower(k) {
		case "header":
			h, err := parseHeaderMap(v)
			if err != nil {
				return fmt.Errorf("invalid %q: %s", k, err)
			}
			rawMap[k] = h

		case "ttl", "interval", "timeout":
			d, err := parseDuration(v)
			if err != nil {
				return fmt.Errorf("invalid %q: %v", k, err)
			}
			rawMap[k] = d

		case "deregister_critical_service_after", "deregistercriticalserviceafter":
			d, err := parseDuration(v)
			if err != nil {
				return fmt.Errorf("invalid %q: %v", k, err)
			}
			replace(k, "DeregisterCriticalServiceAfter", d)

		case "docker_container_id":
			replace(k, "DockerContainerID", v)

		case "service_id":
			replace(k, "ServiceID", v)

		case "tls_skip_verify":
			replace(k, "TLSSkipVerify", v)
		}
	}
	return nil
}

// done(fs): // DecodeCheckDefinition is used to decode a check definition
// done(fs): func DecodeCheckDefinition(raw interface{}) (*structs.CheckDefinition, error) {
// done(fs): 	if err := FixupCheckType(raw); err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs): 	var md mapstructure.Metadata
// done(fs): 	var result structs.CheckDefinition
// done(fs): 	msdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
// done(fs): 		Metadata: &md,
// done(fs): 		Result:   &result,
// done(fs): 	})
// done(fs): 	if err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs): 	if err := msdec.Decode(raw); err != nil {
// done(fs): 		return nil, err
// done(fs): 	}
// done(fs): 	return &result, nil
// done(fs): }
// done(fs):
// done(fs): // MergeConfig merges two configurations together to make a single new
// done(fs): // configuration.
// done(fs): func MergeConfig(a, b *Config) *Config {
// done(fs): 	var result Config = *a
// done(fs):
// done(fs): 	if b.Limits.RPCRate > 0 {
// done(fs): 		result.Limits.RPCRate = b.Limits.RPCRate
// done(fs): 	}
// done(fs): 	if b.Limits.RPCMaxBurst > 0 {
// done(fs): 		result.Limits.RPCMaxBurst = b.Limits.RPCMaxBurst
// done(fs): 	}
// done(fs):
// done(fs): 	// Propagate non-default performance settings
// done(fs): 	if b.Performance.RaftMultiplier > 0 {
// done(fs): 		result.Performance.RaftMultiplier = b.Performance.RaftMultiplier
// done(fs): 	}
// done(fs):
// done(fs): 	// Copy the strings if they're set
// done(fs): 	if b.Bootstrap {
// done(fs): 		result.Bootstrap = true
// done(fs): 	}
// done(fs): 	if b.BootstrapExpect != 0 {
// done(fs): 		result.BootstrapExpect = b.BootstrapExpect
// done(fs): 	}
// done(fs): 	if b.Datacenter != "" {
// done(fs): 		result.Datacenter = b.Datacenter
// done(fs): 	}
// done(fs): 	if b.DataDir != "" {
// done(fs): 		result.DataDir = b.DataDir
// done(fs): 	}
// done(fs):
// done(fs): 	// Copy the dns recursors
// done(fs): 	result.DNSRecursors = make([]string, 0, len(a.DNSRecursors)+len(b.DNSRecursors))
// done(fs): 	result.DNSRecursors = append(result.DNSRecursors, a.DNSRecursors...)
// done(fs): 	result.DNSRecursors = append(result.DNSRecursors, b.DNSRecursors...)
// done(fs):
// done(fs): 	if b.Domain != "" {
// done(fs): 		result.Domain = b.Domain
// done(fs): 	}
// done(fs): 	if b.EncryptKey != "" {
// done(fs): 		result.EncryptKey = b.EncryptKey
// done(fs): 	}
// done(fs): 	if b.DisableKeyringFile {
// done(fs): 		result.DisableKeyringFile = true
// done(fs): 	}
// done(fs): 	if b.EncryptVerifyIncoming != nil {
// done(fs): 		result.EncryptVerifyIncoming = b.EncryptVerifyIncoming
// done(fs): 	}
// done(fs): 	if b.EncryptVerifyOutgoing != nil {
// done(fs): 		result.EncryptVerifyOutgoing = b.EncryptVerifyOutgoing
// done(fs): 	}
// done(fs): 	if b.LogLevel != "" {
// done(fs): 		result.LogLevel = b.LogLevel
// done(fs): 	}
// done(fs): 	if b.Protocol > 0 {
// done(fs): 		result.Protocol = b.Protocol
// done(fs): 	}
// done(fs): 	if b.RaftProtocol > 0 {
// done(fs): 		result.RaftProtocol = b.RaftProtocol
// done(fs): 	}
// done(fs): 	if b.NodeID != "" {
// done(fs): 		result.NodeID = b.NodeID
// done(fs): 	}
// done(fs): 	if b.DisableHostNodeID != nil {
// done(fs): 		result.DisableHostNodeID = b.DisableHostNodeID
// done(fs): 	}
// done(fs): 	if b.NodeName != "" {
// done(fs): 		result.NodeName = b.NodeName
// done(fs): 	}
// done(fs): 	if b.ClientAddr != "" {
// done(fs): 		result.ClientAddr = b.ClientAddr
// done(fs): 	}
// done(fs): 	if b.BindAddr != "" {
// done(fs): 		result.BindAddr = b.BindAddr
// done(fs): 	}
// done(fs): 	if b.AdvertiseAddr != "" {
// done(fs): 		result.AdvertiseAddr = b.AdvertiseAddr
// done(fs): 	}
// done(fs): 	if b.AdvertiseAddrWan != "" {
// done(fs): 		result.AdvertiseAddrWan = b.AdvertiseAddrWan
// done(fs): 	}
// done(fs): 	if b.SerfWanBindAddr != "" {
// done(fs): 		result.SerfWanBindAddr = b.SerfWanBindAddr
// done(fs): 	}
// done(fs): 	if b.SerfLanBindAddr != "" {
// done(fs): 		result.SerfLanBindAddr = b.SerfLanBindAddr
// done(fs): 	}
// done(fs): 	if b.TranslateWanAddrs == true {
// done(fs): 		result.TranslateWanAddrs = true
// done(fs): 	}
// done(fs): 	if b.AdvertiseAddrs.SerfLan != nil {
// done(fs): 		result.AdvertiseAddrs.SerfLan = b.AdvertiseAddrs.SerfLan
// done(fs): 		result.AdvertiseAddrs.SerfLanRaw = b.AdvertiseAddrs.SerfLanRaw
// done(fs): 	}
// done(fs): 	if b.AdvertiseAddrs.SerfWan != nil {
// done(fs): 		result.AdvertiseAddrs.SerfWan = b.AdvertiseAddrs.SerfWan
// done(fs): 		result.AdvertiseAddrs.SerfWanRaw = b.AdvertiseAddrs.SerfWanRaw
// done(fs): 	}
// done(fs): 	if b.AdvertiseAddrs.RPC != nil {
// done(fs): 		result.AdvertiseAddrs.RPC = b.AdvertiseAddrs.RPC
// done(fs): 		result.AdvertiseAddrs.RPCRaw = b.AdvertiseAddrs.RPCRaw
// done(fs): 	}
// done(fs): 	if b.Server == true {
// done(fs): 		result.Server = b.Server
// done(fs): 	}
// done(fs): 	if b.NonVotingServer == true {
// done(fs): 		result.NonVotingServer = b.NonVotingServer
// done(fs): 	}
// done(fs): 	if b.LeaveOnTerm != nil {
// done(fs): 		result.LeaveOnTerm = b.LeaveOnTerm
// done(fs): 	}
// done(fs): 	if b.SkipLeaveOnInt != nil {
// done(fs): 		result.SkipLeaveOnInt = b.SkipLeaveOnInt
// done(fs): 	}
// done(fs): 	if b.Autopilot.CleanupDeadServers != nil {
// done(fs): 		result.Autopilot.CleanupDeadServers = b.Autopilot.CleanupDeadServers
// done(fs): 	}
// done(fs): 	if b.Autopilot.LastContactThreshold != nil {
// done(fs): 		result.Autopilot.LastContactThreshold = b.Autopilot.LastContactThreshold
// done(fs): 	}
// done(fs): 	if b.Autopilot.MaxTrailingLogs != nil {
// done(fs): 		result.Autopilot.MaxTrailingLogs = b.Autopilot.MaxTrailingLogs
// done(fs): 	}
// done(fs): 	if b.Autopilot.ServerStabilizationTime != nil {
// done(fs): 		result.Autopilot.ServerStabilizationTime = b.Autopilot.ServerStabilizationTime
// done(fs): 	}
// done(fs): 	if b.Autopilot.RedundancyZoneTag != "" {
// done(fs): 		result.Autopilot.RedundancyZoneTag = b.Autopilot.RedundancyZoneTag
// done(fs): 	}
// done(fs): 	if b.Autopilot.DisableUpgradeMigration != nil {
// done(fs): 		result.Autopilot.DisableUpgradeMigration = b.Autopilot.DisableUpgradeMigration
// done(fs): 	}
// done(fs): 	if b.Autopilot.UpgradeVersionTag != "" {
// done(fs): 		result.Autopilot.UpgradeVersionTag = b.Autopilot.UpgradeVersionTag
// done(fs): 	}
// done(fs): 	if b.Telemetry.DisableHostname == true {
// done(fs): 		result.Telemetry.DisableHostname = true
// done(fs): 	}
// done(fs): 	if len(b.Telemetry.PrefixFilter) != 0 {
// done(fs): 		result.Telemetry.PrefixFilter = append(result.Telemetry.PrefixFilter, b.Telemetry.PrefixFilter...)
// done(fs): 	}
// done(fs): 	if b.Telemetry.FilterDefault != nil {
// done(fs): 		result.Telemetry.FilterDefault = b.Telemetry.FilterDefault
// done(fs): 	}
// done(fs): 	if b.Telemetry.StatsdAddr != "" {
// done(fs): 		result.Telemetry.StatsdAddr = b.Telemetry.StatsdAddr
// done(fs): 	}
// done(fs): 	if b.Telemetry.StatsiteAddr != "" {
// done(fs): 		result.Telemetry.StatsiteAddr = b.Telemetry.StatsiteAddr
// done(fs): 	}
// done(fs): 	if b.Telemetry.StatsitePrefix != "" {
// done(fs): 		result.Telemetry.StatsitePrefix = b.Telemetry.StatsitePrefix
// done(fs): 	}
// done(fs): 	if b.Telemetry.DogStatsdAddr != "" {
// done(fs): 		result.Telemetry.DogStatsdAddr = b.Telemetry.DogStatsdAddr
// done(fs): 	}
// done(fs): 	if b.Telemetry.DogStatsdTags != nil {
// done(fs): 		result.Telemetry.DogStatsdTags = b.Telemetry.DogStatsdTags
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusAPIToken != "" {
// done(fs): 		result.Telemetry.CirconusAPIToken = b.Telemetry.CirconusAPIToken
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusAPIApp != "" {
// done(fs): 		result.Telemetry.CirconusAPIApp = b.Telemetry.CirconusAPIApp
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusAPIURL != "" {
// done(fs): 		result.Telemetry.CirconusAPIURL = b.Telemetry.CirconusAPIURL
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusCheckSubmissionURL != "" {
// done(fs): 		result.Telemetry.CirconusCheckSubmissionURL = b.Telemetry.CirconusCheckSubmissionURL
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusSubmissionInterval != "" {
// done(fs): 		result.Telemetry.CirconusSubmissionInterval = b.Telemetry.CirconusSubmissionInterval
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusCheckID != "" {
// done(fs): 		result.Telemetry.CirconusCheckID = b.Telemetry.CirconusCheckID
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusCheckForceMetricActivation != "" {
// done(fs): 		result.Telemetry.CirconusCheckForceMetricActivation = b.Telemetry.CirconusCheckForceMetricActivation
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusCheckInstanceID != "" {
// done(fs): 		result.Telemetry.CirconusCheckInstanceID = b.Telemetry.CirconusCheckInstanceID
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusCheckSearchTag != "" {
// done(fs): 		result.Telemetry.CirconusCheckSearchTag = b.Telemetry.CirconusCheckSearchTag
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusCheckDisplayName != "" {
// done(fs): 		result.Telemetry.CirconusCheckDisplayName = b.Telemetry.CirconusCheckDisplayName
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusCheckTags != "" {
// done(fs): 		result.Telemetry.CirconusCheckTags = b.Telemetry.CirconusCheckTags
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusBrokerID != "" {
// done(fs): 		result.Telemetry.CirconusBrokerID = b.Telemetry.CirconusBrokerID
// done(fs): 	}
// done(fs): 	if b.Telemetry.CirconusBrokerSelectTag != "" {
// done(fs): 		result.Telemetry.CirconusBrokerSelectTag = b.Telemetry.CirconusBrokerSelectTag
// done(fs): 	}
// done(fs): 	if b.EnableDebug {
// done(fs): 		result.EnableDebug = true
// done(fs): 	}
// done(fs): 	if b.VerifyIncoming {
// done(fs): 		result.VerifyIncoming = true
// done(fs): 	}
// done(fs): 	if b.VerifyIncomingRPC {
// done(fs): 		result.VerifyIncomingRPC = true
// done(fs): 	}
// done(fs): 	if b.VerifyIncomingHTTPS {
// done(fs): 		result.VerifyIncomingHTTPS = true
// done(fs): 	}
// done(fs): 	if b.VerifyOutgoing {
// done(fs): 		result.VerifyOutgoing = true
// done(fs): 	}
// done(fs): 	if b.VerifyServerHostname {
// done(fs): 		result.VerifyServerHostname = true
// done(fs): 	}
// done(fs): 	if b.CAFile != "" {
// done(fs): 		result.CAFile = b.CAFile
// done(fs): 	}
// done(fs): 	if b.CAPath != "" {
// done(fs): 		result.CAPath = b.CAPath
// done(fs): 	}
// done(fs): 	if b.CertFile != "" {
// done(fs): 		result.CertFile = b.CertFile
// done(fs): 	}
// done(fs): 	if b.KeyFile != "" {
// done(fs): 		result.KeyFile = b.KeyFile
// done(fs): 	}
// done(fs): 	if b.ServerName != "" {
// done(fs): 		result.ServerName = b.ServerName
// done(fs): 	}
// done(fs): 	if b.TLSMinVersion != "" {
// done(fs): 		result.TLSMinVersion = b.TLSMinVersion
// done(fs): 	}
// done(fs): 	if len(b.TLSCipherSuites) != 0 {
// done(fs): 		result.TLSCipherSuites = append(result.TLSCipherSuites, b.TLSCipherSuites...)
// done(fs): 	}
// done(fs): 	if b.TLSPreferServerCipherSuites {
// done(fs): 		result.TLSPreferServerCipherSuites = true
// done(fs): 	}
// done(fs): 	if b.Checks != nil {
// done(fs): 		result.Checks = append(result.Checks, b.Checks...)
// done(fs): 	}
// done(fs): 	if b.Services != nil {
// done(fs): 		result.Services = append(result.Services, b.Services...)
// done(fs): 	}
// done(fs): 	if b.Ports.DNS != 0 {
// done(fs): 		result.Ports.DNS = b.Ports.DNS
// done(fs): 	}
// done(fs): 	if b.Ports.HTTP != 0 {
// done(fs): 		result.Ports.HTTP = b.Ports.HTTP
// done(fs): 	}
// done(fs): 	if b.Ports.HTTPS != 0 {
// done(fs): 		result.Ports.HTTPS = b.Ports.HTTPS
// done(fs): 	}
// done(fs): 	if b.Ports.RPC != 0 {
// done(fs): 		result.Ports.RPC = b.Ports.RPC
// done(fs): 	}
// done(fs): 	if b.Ports.SerfLan != 0 {
// done(fs): 		result.Ports.SerfLan = b.Ports.SerfLan
// done(fs): 	}
// done(fs): 	if b.Ports.SerfWan != 0 {
// done(fs): 		result.Ports.SerfWan = b.Ports.SerfWan
// done(fs): 	}
// done(fs): 	if b.Ports.Server != 0 {
// done(fs): 		result.Ports.Server = b.Ports.Server
// done(fs): 	}
// done(fs): 	if b.Addresses.DNS != "" {
// done(fs): 		result.Addresses.DNS = b.Addresses.DNS
// done(fs): 	}
// done(fs): 	if b.Addresses.HTTP != "" {
// done(fs): 		result.Addresses.HTTP = b.Addresses.HTTP
// done(fs): 	}
// done(fs): 	if b.Addresses.HTTPS != "" {
// done(fs): 		result.Addresses.HTTPS = b.Addresses.HTTPS
// done(fs): 	}
// done(fs): 	if b.Addresses.RPC != "" {
// done(fs): 		result.Addresses.RPC = b.Addresses.RPC
// done(fs): 	}
// done(fs): 	if b.Segment != "" {
// done(fs): 		result.Segment = b.Segment
// done(fs): 	}
// done(fs): 	if len(b.Segments) > 0 {
// done(fs): 		result.Segments = append(result.Segments, b.Segments...)
// done(fs): 	}
// done(fs): 	if b.EnableUI {
// done(fs): 		result.EnableUI = true
// done(fs): 	}
// done(fs): 	if b.UIDir != "" {
// done(fs): 		result.UIDir = b.UIDir
// done(fs): 	}
// done(fs): 	if b.PidFile != "" {
// done(fs): 		result.PidFile = b.PidFile
// done(fs): 	}
// done(fs): 	if b.EnableSyslog {
// done(fs): 		result.EnableSyslog = true
// done(fs): 	}
// done(fs): 	if b.RejoinAfterLeave {
// done(fs): 		result.RejoinAfterLeave = true
// done(fs): 	}
// done(fs): 	if b.RetryMaxAttempts != 0 {
// done(fs): 		result.RetryMaxAttempts = b.RetryMaxAttempts
// done(fs): 	}
// done(fs): 	if b.RetryInterval != 0 {
// done(fs): 		result.RetryInterval = b.RetryInterval
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinEC2.AccessKeyID != "" {
// done(fs): 		result.DeprecatedRetryJoinEC2.AccessKeyID = b.DeprecatedRetryJoinEC2.AccessKeyID
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinEC2.SecretAccessKey != "" {
// done(fs): 		result.DeprecatedRetryJoinEC2.SecretAccessKey = b.DeprecatedRetryJoinEC2.SecretAccessKey
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinEC2.Region != "" {
// done(fs): 		result.DeprecatedRetryJoinEC2.Region = b.DeprecatedRetryJoinEC2.Region
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinEC2.TagKey != "" {
// done(fs): 		result.DeprecatedRetryJoinEC2.TagKey = b.DeprecatedRetryJoinEC2.TagKey
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinEC2.TagValue != "" {
// done(fs): 		result.DeprecatedRetryJoinEC2.TagValue = b.DeprecatedRetryJoinEC2.TagValue
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinGCE.ProjectName != "" {
// done(fs): 		result.DeprecatedRetryJoinGCE.ProjectName = b.DeprecatedRetryJoinGCE.ProjectName
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinGCE.ZonePattern != "" {
// done(fs): 		result.DeprecatedRetryJoinGCE.ZonePattern = b.DeprecatedRetryJoinGCE.ZonePattern
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinGCE.TagValue != "" {
// done(fs): 		result.DeprecatedRetryJoinGCE.TagValue = b.DeprecatedRetryJoinGCE.TagValue
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinGCE.CredentialsFile != "" {
// done(fs): 		result.DeprecatedRetryJoinGCE.CredentialsFile = b.DeprecatedRetryJoinGCE.CredentialsFile
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinAzure.TagName != "" {
// done(fs): 		result.DeprecatedRetryJoinAzure.TagName = b.DeprecatedRetryJoinAzure.TagName
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinAzure.TagValue != "" {
// done(fs): 		result.DeprecatedRetryJoinAzure.TagValue = b.DeprecatedRetryJoinAzure.TagValue
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinAzure.SubscriptionID != "" {
// done(fs): 		result.DeprecatedRetryJoinAzure.SubscriptionID = b.DeprecatedRetryJoinAzure.SubscriptionID
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinAzure.TenantID != "" {
// done(fs): 		result.DeprecatedRetryJoinAzure.TenantID = b.DeprecatedRetryJoinAzure.TenantID
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinAzure.ClientID != "" {
// done(fs): 		result.DeprecatedRetryJoinAzure.ClientID = b.DeprecatedRetryJoinAzure.ClientID
// done(fs): 	}
// done(fs): 	if b.DeprecatedRetryJoinAzure.SecretAccessKey != "" {
// done(fs): 		result.DeprecatedRetryJoinAzure.SecretAccessKey = b.DeprecatedRetryJoinAzure.SecretAccessKey
// done(fs): 	}
// done(fs): 	if b.RetryMaxAttemptsWan != 0 {
// done(fs): 		result.RetryMaxAttemptsWan = b.RetryMaxAttemptsWan
// done(fs): 	}
// done(fs): 	if b.RetryIntervalWan != 0 {
// done(fs): 		result.RetryIntervalWan = b.RetryIntervalWan
// done(fs): 	}
// done(fs): 	if b.ReconnectTimeoutLan != 0 {
// done(fs): 		result.ReconnectTimeoutLan = b.ReconnectTimeoutLan
// done(fs): 		result.ReconnectTimeoutLanRaw = b.ReconnectTimeoutLanRaw
// done(fs): 	}
// done(fs): 	if b.ReconnectTimeoutWan != 0 {
// done(fs): 		result.ReconnectTimeoutWan = b.ReconnectTimeoutWan
// done(fs): 		result.ReconnectTimeoutWanRaw = b.ReconnectTimeoutWanRaw
// done(fs): 	}
// done(fs): 	if b.DNSConfig.NodeTTL != 0 {
// done(fs): 		result.DNSConfig.NodeTTL = b.DNSConfig.NodeTTL
// done(fs): 	}
// done(fs): 	if len(b.DNSConfig.ServiceTTL) != 0 {
// done(fs): 		if result.DNSConfig.ServiceTTL == nil {
// done(fs): 			result.DNSConfig.ServiceTTL = make(map[string]time.Duration)
// done(fs): 		}
// done(fs): 		for service, dur := range b.DNSConfig.ServiceTTL {
// done(fs): 			result.DNSConfig.ServiceTTL[service] = dur
// done(fs): 		}
// done(fs): 	}
// done(fs): 	if b.DNSConfig.AllowStale != nil {
// done(fs): 		result.DNSConfig.AllowStale = b.DNSConfig.AllowStale
// done(fs): 	}
// done(fs): 	if b.DNSConfig.UDPAnswerLimit != 0 {
// done(fs): 		result.DNSConfig.UDPAnswerLimit = b.DNSConfig.UDPAnswerLimit
// done(fs): 	}
// done(fs): 	if b.DNSConfig.EnableTruncate {
// done(fs): 		result.DNSConfig.EnableTruncate = true
// done(fs): 	}
// done(fs): 	if b.DNSConfig.MaxStale != 0 {
// done(fs): 		result.DNSConfig.MaxStale = b.DNSConfig.MaxStale
// done(fs): 	}
// done(fs): 	if b.DNSConfig.OnlyPassing {
// done(fs): 		result.DNSConfig.OnlyPassing = true
// done(fs): 	}
// done(fs): 	if b.DNSConfig.DisableCompression {
// done(fs): 		result.DNSConfig.DisableCompression = true
// done(fs): 	}
// done(fs): 	if b.DNSConfig.RecursorTimeout != 0 {
// done(fs): 		result.DNSConfig.RecursorTimeout = b.DNSConfig.RecursorTimeout
// done(fs): 	}
// done(fs): 	if b.EnableScriptChecks {
// done(fs): 		result.EnableScriptChecks = true
// done(fs): 	}
// done(fs): 	if b.CheckUpdateIntervalRaw != "" || b.CheckUpdateInterval != 0 {
// done(fs): 		result.CheckUpdateInterval = b.CheckUpdateInterval
// done(fs): 	}
// done(fs): 	if b.SyslogFacility != "" {
// done(fs): 		result.SyslogFacility = b.SyslogFacility
// done(fs): 	}
// done(fs): 	if b.ACLToken != "" {
// done(fs): 		result.ACLToken = b.ACLToken
// done(fs): 	}
// done(fs): 	if b.ACLAgentMasterToken != "" {
// done(fs): 		result.ACLAgentMasterToken = b.ACLAgentMasterToken
// done(fs): 	}
// done(fs): 	if b.ACLAgentToken != "" {
// done(fs): 		result.ACLAgentToken = b.ACLAgentToken
// done(fs): 	}
// done(fs): 	if b.ACLMasterToken != "" {
// done(fs): 		result.ACLMasterToken = b.ACLMasterToken
// done(fs): 	}
// done(fs): 	if b.ACLDatacenter != "" {
// done(fs): 		result.ACLDatacenter = b.ACLDatacenter
// done(fs): 	}
// done(fs): 	if b.ACLTTLRaw != "" {
// done(fs): 		result.ACLTTL = b.ACLTTL
// done(fs): 		result.ACLTTLRaw = b.ACLTTLRaw
// done(fs): 	}
// done(fs): 	if b.ACLDownPolicy != "" {
// done(fs): 		result.ACLDownPolicy = b.ACLDownPolicy
// done(fs): 	}
// done(fs): 	if b.ACLDefaultPolicy != "" {
// done(fs): 		result.ACLDefaultPolicy = b.ACLDefaultPolicy
// done(fs): 	}
// done(fs): 	if b.EnableACLReplication {
// done(fs): 		result.EnableACLReplication = true
// done(fs): 	}
// done(fs): 	if b.ACLReplicationToken != "" {
// done(fs): 		result.ACLReplicationToken = b.ACLReplicationToken
// done(fs): 	}
// done(fs): 	if b.ACLEnforceVersion8 != nil {
// done(fs): 		result.ACLEnforceVersion8 = b.ACLEnforceVersion8
// done(fs): 	}
// done(fs): 	if len(b.Watches) != 0 {
// done(fs): 		result.Watches = append(result.Watches, b.Watches...)
// done(fs): 	}
// done(fs): 	if len(b.WatchPlans) != 0 {
// done(fs): 		result.WatchPlans = append(result.WatchPlans, b.WatchPlans...)
// done(fs): 	}
// done(fs): 	if b.DisableRemoteExec != nil {
// done(fs): 		result.DisableRemoteExec = b.DisableRemoteExec
// done(fs): 	}
// done(fs): 	if b.DisableUpdateCheck {
// done(fs): 		result.DisableUpdateCheck = true
// done(fs): 	}
// done(fs): 	if b.DisableAnonymousSignature {
// done(fs): 		result.DisableAnonymousSignature = true
// done(fs): 	}
// done(fs): 	if b.UnixSockets.Usr != "" {
// done(fs): 		result.UnixSockets.Usr = b.UnixSockets.Usr
// done(fs): 	}
// done(fs): 	if b.UnixSockets.Grp != "" {
// done(fs): 		result.UnixSockets.Grp = b.UnixSockets.Grp
// done(fs): 	}
// done(fs): 	if b.UnixSockets.Perms != "" {
// done(fs): 		result.UnixSockets.Perms = b.UnixSockets.Perms
// done(fs): 	}
// done(fs): 	if b.DisableCoordinates {
// done(fs): 		result.DisableCoordinates = true
// done(fs): 	}
// done(fs): 	if b.SessionTTLMinRaw != "" {
// done(fs): 		result.SessionTTLMin = b.SessionTTLMin
// done(fs): 		result.SessionTTLMinRaw = b.SessionTTLMinRaw
// done(fs): 	}
// done(fs):
// done(fs): 	result.HTTPConfig.BlockEndpoints = append(a.HTTPConfig.BlockEndpoints,
// done(fs): 		b.HTTPConfig.BlockEndpoints...)
// done(fs): 	if len(b.HTTPConfig.ResponseHeaders) > 0 {
// done(fs): 		if result.HTTPConfig.ResponseHeaders == nil {
// done(fs): 			result.HTTPConfig.ResponseHeaders = make(map[string]string)
// done(fs): 		}
// done(fs): 		for field, value := range b.HTTPConfig.ResponseHeaders {
// done(fs): 			result.HTTPConfig.ResponseHeaders[field] = value
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	if len(b.Meta) != 0 {
// done(fs): 		if result.Meta == nil {
// done(fs): 			result.Meta = make(map[string]string)
// done(fs): 		}
// done(fs): 		for field, value := range b.Meta {
// done(fs): 			result.Meta[field] = value
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	// Copy the start join addresses
// done(fs): 	result.StartJoin = make([]string, 0, len(a.StartJoin)+len(b.StartJoin))
// done(fs): 	result.StartJoin = append(result.StartJoin, a.StartJoin...)
// done(fs): 	result.StartJoin = append(result.StartJoin, b.StartJoin...)
// done(fs):
// done(fs): 	// Copy the start join addresses
// done(fs): 	result.StartJoinWan = make([]string, 0, len(a.StartJoinWan)+len(b.StartJoinWan))
// done(fs): 	result.StartJoinWan = append(result.StartJoinWan, a.StartJoinWan...)
// done(fs): 	result.StartJoinWan = append(result.StartJoinWan, b.StartJoinWan...)
// done(fs):
// done(fs): 	// Copy the retry join addresses
// done(fs): 	result.RetryJoin = make([]string, 0, len(a.RetryJoin)+len(b.RetryJoin))
// done(fs): 	result.RetryJoin = append(result.RetryJoin, a.RetryJoin...)
// done(fs): 	result.RetryJoin = append(result.RetryJoin, b.RetryJoin...)
// done(fs):
// done(fs): 	// Copy the retry join -wan addresses
// done(fs): 	result.RetryJoinWan = make([]string, 0, len(a.RetryJoinWan)+len(b.RetryJoinWan))
// done(fs): 	result.RetryJoinWan = append(result.RetryJoinWan, a.RetryJoinWan...)
// done(fs): 	result.RetryJoinWan = append(result.RetryJoinWan, b.RetryJoinWan...)
// done(fs):
// done(fs): 	return &result
// done(fs): }
// done(fs):
// done(fs): // ReadConfigPaths reads the paths in the given order to load configurations.
// done(fs): // The paths can be to files or directories. If the path is a directory,
// done(fs): // we read one directory deep and read any files ending in ".json" as
// done(fs): // configuration files.
// done(fs): func ReadConfigPaths(paths []string) (*Config, error) {
// done(fs): 	result := new(Config)
// done(fs): 	for _, path := range paths {
// done(fs): 		f, err := os.Open(path)
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("Error reading '%s': %s", path, err)
// done(fs): 		}
// done(fs):
// done(fs): 		fi, err := f.Stat()
// done(fs): 		if err != nil {
// done(fs): 			f.Close()
// done(fs): 			return nil, fmt.Errorf("Error reading '%s': %s", path, err)
// done(fs): 		}
// done(fs):
// done(fs): 		if !fi.IsDir() {
// done(fs): 			config, err := DecodeConfig(f)
// done(fs): 			f.Close()
// done(fs):
// done(fs): 			if err != nil {
// done(fs): 				return nil, fmt.Errorf("Error decoding '%s': %s", path, err)
// done(fs): 			}
// done(fs):
// done(fs): 			result = MergeConfig(result, config)
// done(fs): 			continue
// done(fs): 		}
// done(fs):
// done(fs): 		contents, err := f.Readdir(-1)
// done(fs): 		f.Close()
// done(fs): 		if err != nil {
// done(fs): 			return nil, fmt.Errorf("Error reading '%s': %s", path, err)
// done(fs): 		}
// done(fs):
// done(fs): 		// Sort the contents, ensures lexical order
// done(fs): 		sort.Sort(dirEnts(contents))
// done(fs):
// done(fs): 		for _, fi := range contents {
// done(fs): 			// Don't recursively read contents
// done(fs): 			if fi.IsDir() {
// done(fs): 				continue
// done(fs): 			}
// done(fs):
// done(fs): 			// If it isn't a JSON file, ignore it
// done(fs): 			if !strings.HasSuffix(fi.Name(), ".json") {
// done(fs): 				continue
// done(fs): 			}
// done(fs): 			// If the config file is empty, ignore it
// done(fs): 			if fi.Size() == 0 {
// done(fs): 				continue
// done(fs): 			}
// done(fs):
// done(fs): 			subpath := filepath.Join(path, fi.Name())
// done(fs): 			f, err := os.Open(subpath)
// done(fs): 			if err != nil {
// done(fs): 				return nil, fmt.Errorf("Error reading '%s': %s", subpath, err)
// done(fs): 			}
// done(fs):
// done(fs): 			config, err := DecodeConfig(f)
// done(fs): 			f.Close()
// done(fs):
// done(fs): 			if err != nil {
// done(fs): 				return nil, fmt.Errorf("Error decoding '%s': %s", subpath, err)
// done(fs): 			}
// done(fs):
// done(fs): 			result = MergeConfig(result, config)
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	return result, nil
// done(fs): }
// done(fs):
// done(fs): // ResolveTmplAddrs iterates over the myriad of addresses in the agent's config
// done(fs): // and performs go-sockaddr/template Parse on each known address in case the
// done(fs): // user specified a template config for any of their values.
// done(fs): func (c *Config) ResolveTmplAddrs() (err error) {
// done(fs): 	parse := func(addr *string, socketAllowed bool, name string) {
// done(fs): 		if *addr == "" || err != nil {
// done(fs): 			return
// done(fs): 		}
// done(fs): 		var ip string
// done(fs): 		ip, err = parseSingleIPTemplate(*addr)
// done(fs): 		if err != nil {
// done(fs): 			err = fmt.Errorf("Resolution of %s failed: %v", name, err)
// done(fs): 			return
// done(fs): 		}
// done(fs): 		ipAddr := net.ParseIP(ip)
// done(fs): 		if !socketAllowed && ipAddr == nil {
// done(fs): 			err = fmt.Errorf("Failed to parse %s: %v", name, ip)
// done(fs): 			return
// done(fs): 		}
// done(fs): 		if socketAllowed && socketPath(ip) == "" && ipAddr == nil {
// done(fs): 			err = fmt.Errorf("Failed to parse %s, %q is not a valid IP address or socket", name, ip)
// done(fs): 			return
// done(fs): 		}
// done(fs):
// done(fs): 		*addr = ip
// done(fs): 	}
// done(fs):
// done(fs): 	if c == nil {
// done(fs): 		return
// done(fs): 	}
// done(fs): 	parse(&c.Addresses.DNS, true, "DNS address")
// done(fs): 	parse(&c.Addresses.HTTP, true, "HTTP address")
// done(fs): 	parse(&c.Addresses.HTTPS, true, "HTTPS address")
// done(fs): 	parse(&c.AdvertiseAddr, false, "Advertise address")
// done(fs): 	parse(&c.AdvertiseAddrWan, false, "Advertise WAN address")
// done(fs): 	parse(&c.BindAddr, true, "Bind address")
// done(fs): 	parse(&c.ClientAddr, true, "Client address")
// done(fs): 	parse(&c.SerfLanBindAddr, false, "Serf LAN address")
// done(fs): 	parse(&c.SerfWanBindAddr, false, "Serf WAN address")
// done(fs): 	for i, segment := range c.Segments {
// done(fs): 		parse(&c.Segments[i].Bind, false, fmt.Sprintf("Segment %q bind address", segment.Name))
// done(fs): 		parse(&c.Segments[i].Advertise, false, fmt.Sprintf("Segment %q advertise address", segment.Name))
// done(fs): 	}
// done(fs):
// done(fs): 	return
// done(fs): }
// done(fs):
// done(fs): // SetupTaggedAndAdvertiseAddrs configures advertise addresses and sets up a map of tagged addresses
// done(fs): func (cfg *Config) SetupTaggedAndAdvertiseAddrs() error {
// done(fs): 	if cfg.AdvertiseAddr == "" {
// done(fs): 		switch {
// done(fs):
// done(fs): 		case cfg.BindAddr != "" && !ipaddr.IsAny(cfg.BindAddr):
// done(fs): 			cfg.AdvertiseAddr = cfg.BindAddr
// done(fs):
// done(fs): 		default:
// done(fs): 			ip, err := consul.GetPrivateIP()
// done(fs): 			if ipaddr.IsAnyV6(cfg.BindAddr) {
// done(fs): 				ip, err = consul.GetPublicIPv6()
// done(fs): 			}
// done(fs): 			if err != nil {
// done(fs): 				return fmt.Errorf("Failed to get advertise address: %v", err)
// done(fs): 			}
// done(fs): 			cfg.AdvertiseAddr = ip.String()
// done(fs): 		}
// done(fs): 	}
// done(fs):
// done(fs): 	// Try to get an advertise address for the wan
// done(fs): 	if cfg.AdvertiseAddrWan == "" {
// done(fs): 		cfg.AdvertiseAddrWan = cfg.AdvertiseAddr
// done(fs): 	}
// done(fs):
// done(fs): 	// Create the default set of tagged addresses.
// done(fs): 	cfg.TaggedAddresses = map[string]string{
// done(fs): 		"lan": cfg.AdvertiseAddr,
// done(fs): 		"wan": cfg.AdvertiseAddrWan,
// done(fs): 	}
// done(fs): 	return nil
// done(fs): }
// done(fs):
// done(fs): // parseSingleIPTemplate is used as a helper function to parse out a single IP
// done(fs): // address from a config parameter.
// done(fs): func parseSingleIPTemplate(ipTmpl string) (string, error) {
// done(fs): 	out, err := template.Parse(ipTmpl)
// done(fs): 	if err != nil {
// done(fs): 		return "", fmt.Errorf("Unable to parse address template %q: %v", ipTmpl, err)
// done(fs): 	}
// done(fs):
// done(fs): 	ips := strings.Split(out, " ")
// done(fs): 	switch len(ips) {
// done(fs): 	case 0:
// done(fs): 		return "", errors.New("No addresses found, please configure one.")
// done(fs): 	case 1:
// done(fs): 		return ips[0], nil
// done(fs): 	default:
// done(fs): 		return "", fmt.Errorf("Multiple addresses found (%q), please configure one.", out)
// done(fs): 	}
// done(fs): }
// done(fs):
// done(fs): // Implement the sort interface for dirEnts
// done(fs): func (d dirEnts) Len() int {
// done(fs): 	return len(d)
// done(fs): }
// done(fs):
// done(fs): func (d dirEnts) Less(i, j int) bool {
// done(fs): 	return d[i].Name() < d[j].Name()
// done(fs): }
// done(fs):
// done(fs): func (d dirEnts) Swap(i, j int) {
// done(fs): 	d[i], d[j] = d[j], d[i]
// done(fs): }
// done(fs):
// ParseMetaPair parses a key/value pair of the form key:value
func ParseMetaPair(raw string) (string, string) {
	pair := strings.SplitN(raw, ":", 2)
	if len(pair) == 2 {
		return pair[0], pair[1]
	}
	return pair[0], ""
}
