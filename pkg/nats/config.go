package nats

type ServerConfig struct {
	Listen        string          `json:"listen,omitempty"`
	LeafNodes     LeafNodesConfig `json:"leafnodes,omitempty"`
	Jetstream     JetstreamConfig `json:"jetstream"`
	PidFile       string          `json:"pid_file,omitempty"`
	Http          int             `json:"http,omitempty"`
	Operator      string          `json:"operator,omitempty"`
	SystemAccount string          `json:"system_account,omitempty"`
	// Resolver        ResolverConfig    `json:"resolver,omitempty"`
	ResolverPreload map[string]string `json:"resolver_preload,omitempty"`
}

type ResolverConfig struct {
	Type    string `json:"type,omitempty"`
	Dir     string `json:"dir,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

type LeafNodesConfig struct {
	Remotes []LeafNodeRemoteConfig `json:"remotes,omitempty"`
}

type LeafNodeRemoteConfig struct {
	Url         string   `json:"url,omitempty"`
	Credentials string   `json:"credentials,omitempty"`
	Account     string   `json:"account,omitempty"`
	DenyImports []string `json:"deny_imports,omitempty"`
	DenyExports []string `json:"deny_exports,omitempty"`
}

type JetstreamConfig struct {
	MaxMemory string `json:"max_mem,omitempty"`
	MaxFile   string `json:"max_file,omitempty"`
	StoreDir  string `json:"store_dir,omitempty"`
	Domain    string `json:"domain,omitempty"`
}
