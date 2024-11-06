package aptos

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
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
		return nil, fmt.Errorf("failed to create Aptos aptclient: %w", err)
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

func PrivateKeyToPrivateKey(privateKeyHex string) (*crypto.Ed25519PrivateKey, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privKey := &crypto.Ed25519PrivateKey{}

	err := privKey.FromHex(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key from hex: %w", err)
	}

	return privKey, nil
}

func PrivateKeyToPubKey(privateKey *crypto.Ed25519PrivateKey) (*crypto.Ed25519PublicKey, error) {
	publicKey := privateKey.PubKey().(*crypto.Ed25519PublicKey)
	return publicKey, nil
}

func PrivateKeyHexToPubKey(privateKeyHex string) (*crypto.Ed25519PublicKey, error) {
	ed25519PrivateKey, err := PrivateKeyToPrivateKey(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key from hex: %w", err)
	}
	publicKey := ed25519PrivateKey.PubKey().(*crypto.Ed25519PublicKey)
	return publicKey, nil
}

func PubKeyHexToPubKey(publicKeyHex string) (*crypto.Ed25519PublicKey, error) {
	publicKeyHex = strings.TrimPrefix(publicKeyHex, "0x")
	pubKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("decode public key failed: %w", err)
	}
	publicKey := &crypto.Ed25519PublicKey{}
	err = publicKey.FromBytes(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("create public key failed: %w", err)
	}
	return publicKey, nil
}

func PubKeyToPubKeyHex(publicKey *crypto.Ed25519PublicKey) (string, error) {
	return publicKey.ToHex(), nil
}

func PubKeyToAddress(publicKey *crypto.Ed25519PublicKey) (string, error) {
	authKey := publicKey.AuthKey()
	address := "0x" + hex.EncodeToString(authKey[:])
	return address, nil
}

func PubKeyHexToAddress(publicKeyHex string) (string, error) {
	ed25519PublicKey, err := PubKeyHexToPubKey(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("create public key failed: %w", err)
	}
	authKey := ed25519PublicKey.AuthKey()

	address := "0x" + hex.EncodeToString(authKey[:])

	return address, nil
}

func PubKeyToAccountAddress(publicKey *crypto.Ed25519PublicKey) (*aptos.AccountAddress, error) {
	authKey := publicKey.AuthKey()
	address := aptos.AccountAddress{}
	copy(address[:], authKey[:])
	return &address, nil
}

func PubKeyHexToAccountAddress(publicKeyHex string) (*aptos.AccountAddress, error) {
	ed25519PublicKey, err := PubKeyHexToPubKey(publicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("create public key failed: %w", err)
	}

	authKey := &crypto.AuthenticationKey{}
	authKey.FromPublicKey(ed25519PublicKey)
	address := &aptos.AccountAddress{}
	address.FromAuthKey(authKey)
	return address, nil

	//authKey := ed25519PublicKey.AuthKey()
	//address := aptos.AccountAddress{}
	//copy(address[:], authKey[:])
	//return &address, nil
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
