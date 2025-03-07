package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/ethereum/go-ethereum/log"
)

type Server struct {
	Port string `yaml:"port"`
}

type Node struct {
	RpcUrl       string `yaml:"rpc_url"`
	RpcUser      string `yaml:"rpc_user"`
	RpcPass      string `yaml:"rpc_pass"`
	DataApiUrl   string `yaml:"data_api_url"`
	DataApiKey   string `yaml:"data_api_key"`
	DataApiToken string `yaml:"data_api_token"`
	TimeOut      uint64 `yaml:"time_out"`
}

type WalletNode struct {
	Eth     Node `yaml:"eth"`
	Arbi    Node `yaml:"arbi"`
	Op      Node `yaml:"op"`
	Zksync  Node `yaml:"zksync"`
	Bsc     Node `yaml:"bsc"`
	Heco    Node `yaml:"heco"`
	Avax    Node `yaml:"avax"`
	Polygon Node `yaml:"polygon"`
	Tron    Node `yaml:"tron"`
	Sol     Node `yaml:"solana"`
	Cosmos  Node `yaml:"cosmos"`
	Aptos   Node `yaml:"aptos"`
	Mantle  Node `yaml:"mantle"`
	Scroll  Node `yaml:"scroll"`
	Base    Node `yaml:"evmbase"`
	Linea   Node `yaml:"linea"`
	Sui     Node `yaml:"sui"`
	Ton     Node `yaml:"ton"`
	Xlm     Node `yaml:"xlm"`
	Icp     Node `yaml:"icp"`
	Btt     Node `yaml:"btt"`
}

type Config struct {
	Server     Server     `yaml:"server"`
	WalletNode WalletNode `yaml:"wallet_node"`
	NetWork    string     `yaml:"network"`
	Chains     []string   `yaml:"chains"`
}

func New(path string) (*Config, error) {
	var config = new(Config)
	h := log.NewTerminalHandler(os.Stdout, true)
	log.SetDefault(log.NewLogger(h))

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error("read config file error", "err", err)
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Error("unmarshal config file error", "err", err)
		return nil, err
	}
	return config, nil
}

const UnsupportedChain = "Unsupport chain"
const UnsupportedOperation = UnsupportedChain
