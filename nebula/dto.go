package nebula

type Conf struct {
	Relay         Relay               `yaml:"relay"`
	Tun           Tun                 `yaml:"tun"`
	Logging       Logging             `yaml:"logging"`
	Firewall      Firewall            `yaml:"firewall"`
	Pki           Pki                 `yaml:"pki"`
	StaticHostMap map[string][]string `yaml:"static_host_map"`
	Lighthouse    Lighthouse          `yaml:"lighthouse"`
	Listen        Listen              `yaml:"listen"`
	Punchy        Punchy              `yaml:"punchy"`
	Handshakes    Handshakes          `yaml:"handshakes"`
	StaticMap     StaticMap           `yaml:"static_map"`
	Cipher        string              `yaml:"cipher"`
	Routines      int                 `yaml:"routines"`
}

type StaticMap struct {
	Network       string `yaml:"network"`
	Cadence       string `yaml:"cadence"`
	LookupTimeout string `yaml:"lookup_timeout"`
}

type Handshakes struct {
	TriggerBuffer int    `yaml:"trigger_buffer"`
	TryInterval   string `yaml:"try_interval"`
	Retries       int    `yaml:"retries"`
}

type Relay struct {
	Relays    []string `yaml:"relays"`
	AmRelay   bool     `yaml:"am_relay"`
	UseRelays bool     `yaml:"use_relays"`
}

type Firewall struct {
	OutboundAction string     `yaml:"outbound_action"`
	InboundAction  string     `yaml:"inbound_action"`
	Conntrack      Conntrack  `yaml:"conntrack"`
	Outbound       []Outbound `yaml:"outbound"`
	Inbound        []Inbound  `yaml:"inbound"`
}

type Lighthouse struct {
	AmLighthouse   bool           `yaml:"am_lighthouse"`
	Hosts          []string       `yaml:"hosts"`
	Interval       int            `yaml:"interval"`
	LocalAllowList LocalAllowList `yaml:"local_allow_list"`
	AdvertiseAddrs []string       `yaml:"advertise_addrs"`
}

type LocalAllowList struct {
	Interfaces map[string]bool `yaml:"interfaces"`
	Subnets    map[string]bool `yaml:",inline"`
}

type Listen struct {
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
}

type Outbound struct {
	Port  string `yaml:"port"`
	Proto string `yaml:"proto"`
	Host  string `yaml:"host"`
}

type Inbound struct {
	Port   string   `yaml:"port"`
	Proto  string   `yaml:"proto"`
	Host   string   `yaml:"host"`
	Groups []string `yaml:"groups"`
}

type Pki struct {
	Ca                string `yaml:"ca"`
	Cert              string `yaml:"cert"`
	Key               string `yaml:"key"`
	DisconnectInvalid bool   `yaml:"disconnect_invalid"`
}

type Punchy struct {
	Delay        string `yaml:"delay"`
	RespondDelay string `yaml:"respond_delay"`
	Punch        bool   `yaml:"punch"`
	Respond      bool   `yaml:"respond"`
}

type Tun struct {
	UnsafeRoutes       []UnsafeRoutes `yaml:"unsafe_routes"`
	Disabled           bool           `yaml:"disabled"`
	Dev                string         `yaml:"dev"`
	DropLocalBroadcast bool           `yaml:"drop_local_broadcast"`
	DropMulticast      bool           `yaml:"drop_multicast"`
	TxQueue            int            `yaml:"tx_queue"`
	Mtu                int            `yaml:"mtu"`
	Routes             []Routes       `yaml:"routes"`
}

type Routes struct {
	Mtu   int    `yaml:"mtu"`
	Route string `yaml:"route"`
}

type UnsafeRoutes struct {
	Via     string `yaml:"via"`
	Mtu     int    `yaml:"mtu"`
	Metric  int    `yaml:"metric"`
	Install bool   `yaml:"install"`
	Route   string `yaml:"route"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Conntrack struct {
	TcpTimeout     string `yaml:"tcp_timeout"`
	UdpTimeout     string `yaml:"udp_timeout"`
	DefaultTimeout string `yaml:"default_timeout"`
}
