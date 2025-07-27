package config

import (
	"os"
)

type Config struct {
	RPCURL          string
	ContractAddress string
	PrivateKey      string
	ABIPath         string
	PostgresDSN     string
}

func LoadConfig() *Config {
	return &Config{
		RPCURL:          os.Getenv("BESU_NODE_URL"),
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
		PrivateKey:      os.Getenv("SIGNER_PRIVATE_KEY"),
		ABIPath:         os.Getenv("CONTRACT_ABI_PATH"),
		PostgresDSN:     os.Getenv("POSTGRES_DSN"),
	}
}
