package nats

type ServerConfig struct {
	Listen    string          `json:"listen"`
	LeafNodes LeafNodesConfig `json:"leafnodes"`
	Jetstream JetstreamConfig `json:"jetstream"`
}

type LeafNodesConfig struct {
	Remotes []LeafNodeRemoteConfig `json:"remotes"`
}
type LeafNodeRemoteConfig struct {
	Url         string `json:"url"`
	Credentials string `json:"credentials"`
}

type JetstreamConfig struct {
	MaxMemory string `json:"max_mem"`
	MaxFile   string `json:"max_file"`
	StoreDir  string `json:"store_dir"`
	Domain    string `json:"domain"`
}
