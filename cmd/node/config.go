package main

type Node struct {
	Port int `koanf:"port"`
}

type Backend struct {
	Port  int               `koanf:"port"`
	Users map[string]string `koanf:"users"`
}

type NATS struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type NodeConfig struct {
	Node    `koanf:"node"`
	Backend `koanf:"backend"`
	NATS    `koanf:"nats"`
}

// FIXME: move elsewhere
var DefaultConfig = NodeConfig{
	Node: Node{
		Port: 8080,
	},
	Backend: Backend{
		Port: 8080,
	},
	NATS: NATS{
		Port: 0,
	},
}
