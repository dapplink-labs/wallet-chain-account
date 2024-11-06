package aptos

import (
	"encoding/hex"
	"fmt"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"strings"
)

type Environment string

const (
	Devnet  Environment = "devnet"
	Testnet Environment = "testnet"
	Mainnet Environment = "mainnet"
	Local   Environment = "local"
)

func NewAptosClient(networkConfig string) (*aptos.Client, error) {
	if networkConfig == "" {
		return nil, fmt.Errorf("network configuration is empty")
	}

	aptosEnv, ok := ConvertEnvironment(networkConfig)
	if !ok {
		return nil, fmt.Errorf("unsupported network environment: %s", networkConfig)
	}

	client, err := NewAptosClientEnv(aptosEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aptos client: %w", err)
	}

	return client, nil
}

func NewAptosClientEnv(env Environment) (*aptos.Client, error) {
	var config aptos.NetworkConfig

	switch env {
	case Devnet:
		config = aptos.DevnetConfig
	case Testnet:
		config = aptos.TestnetConfig
	case Mainnet:
		config = aptos.MainnetConfig
	case Local:
		// Assuming local node runs on default port
		config = aptos.LocalnetConfig
	default:
		return nil, fmt.Errorf("unsupported environment: %s", env)
	}

	return aptos.NewClient(config)
}

func ConvertEnvironment(network string) (Environment, bool) {
	validEnvs := map[string]Environment{
		"devnet":  Devnet,
		"testnet": Testnet,
		"mainnet": Mainnet,
		"local":   Local,
	}

	env, exists := validEnvs[strings.ToLower(network)]
	return env, exists
}

func AddressToAccountAddress(address string) (aptos.AccountAddress, error) {
	address = strings.TrimPrefix(address, "0x")

	if len(address) != 64 {
		return aptos.AccountAddress{}, fmt.Errorf("invalid address length: expected 64 hex chars, got %d", len(address))
	}

	bytes, err := hex.DecodeString(address)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("failed to decode hex address: %w", err)
	}

	var accountAddress aptos.AccountAddress
	copy(accountAddress[:], bytes)

	return accountAddress, nil
}

func PrivateKeyToPrivateKey(privateKey string) (*crypto.Ed25519PrivateKey, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")

	privKey := &crypto.Ed25519PrivateKey{}

	err := privKey.FromHex(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key from hex: %w", err)
	}

	return privKey, nil
}
