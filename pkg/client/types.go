package client

import "crypto/ed25519"

type Member struct {
	Key       string `yaml:"key" json:"key"`
	Metadata  string `yaml:"metadata" json:"metadata"`
	Signature string `yaml:"signature" json:"signature"`
}

type Body struct {
	Data      string `yaml:"data" json:"data"`
	Timestamp int64  `yaml:"timestamp" json:"timestamp"`
	Signature string `yaml:"signature" json:"signature"`
}

type Config struct {
	MasterKey string `yaml:"masterKey"`
	ServerURL string `yaml:"serverUrl"`
	Member    Member `yaml:"member"`
}

type RequestBody struct {
	MasterKey string `json:"masterKey"`
	Member    Member `json:"member"`
	Body      Body   `json:"body"`
}

type Client struct {
	Config     *Config
	PublicKey  string
	PrivateKey ed25519.PrivateKey
}
