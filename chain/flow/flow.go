package flow

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cometbft/cometbft/crypto"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/log"
	"github.com/onflow/flow-go-sdk"
	"strconv"
)

const ChainName = "Flow"

type ChainAdaptor struct {
	flowClient *FlowClient
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	ctx := context.Background()
	result, err := c.flowClient.client.GetBlockByID(ctx, flow.HexToID(req.Hash))
	if err != nil {
		return nil, err
	}
	return &account.BlockResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "success",
		Height: int64(result.Height),
		Hash:   string(result.PayloadHash),
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	ctx := context.Background()
	result, err := c.flowClient.client.GetBlockHeaderByID(ctx, flow.HexToID(req.Hash))
	if err != nil {
		return nil, err
	}
	return &account.BlockHeaderResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		BlockHeader: &account.BlockHeader{
			Hash:       string(result.PayloadHash),
			Nonce:      result.ID.String(),
			ParentHash: result.ParentID.String(),
		},
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	//ctx := context.Background()
	//result, err := c.flowClient.client.GetBlockHeaderByID(ctx, flow.HexToID(string(req.Height)))
	//if err != nil {
	//	return nil, err
	//}
	//return &account.BlockHeaderResponse{
	//	Code: common.ReturnCode_SUCCESS,
	//	Msg:  "success",
	//	BlockHeader: &account.BlockHeader{
	//		Hash:       string(result.PayloadHash),
	//		Nonce:      result.ID.String(),
	//		ParentHash: result.ParentID.String(),
	//	},
	//}, nil
	return nil, nil
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	ctx := context.Background()
	address := flow.HexToAddress(req.Address)
	fmt.Println("输出address:", address)
	accountInfo, err := c.flowClient.client.GetAccount(ctx, address)
	if err != nil {
		return nil, err
	}
	fmt.Println(accountInfo)
	//TODO implement me
	return &account.AccountResponse{
		Balance: strconv.FormatUint(accountInfo.Balance, 10),
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "success",
	}, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	var txReq flow.Transaction
	unmarshalErr := json.Unmarshal([]byte(req.RawTx), &txReq)
	if unmarshalErr != nil {
		log.Error("flows tx unmarshal  Error: %+v\n", unmarshalErr)
		panic(unmarshalErr)
	}
	unmarshalErr = c.flowClient.httpClient.SendTransaction(context.Background(), txReq)
	if unmarshalErr != nil {
		log.Error("flows tx send  Error: %+v\n", unmarshalErr)
		panic(unmarshalErr)
	}
	return &account.SendTxResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

// 根据区块id获取交易信息
func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	result, err := c.flowClient.client.GetTransactionsByBlockID(context.Background(), flow.HexToID(req.Hash))
	if err != nil {
		return nil, err
	}
	//addresses := stream.FromSlice(result).Map(func(item *flow.Transaction) *account.Address {
	//	return &account.Address{
	//		Address: item.Payer.String(),
	//	}
	//}).ToSlice()

	fmt.Println(result)
	return &account.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		Tx: &account.TxMessage{
			Hash: req.Hash,
		},
	}, nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	massageHash, err := GetMassageHash(req.Base64Tx)
	if err != nil {
		return nil, err
	}
	return &account.UnSignTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "create unsigned transaction success",
		UnSignTx: massageHash,
	}, nil
}

// 通过base64解出原始结构体并推出原始哈希
func GetMassageHash(base64Tx string) (string, error) {
	//解出base64的原始json
	obj, err := parseToflowByTx(base64Tx)
	if err != nil {
		return "error", err
	}
	//数据整理一下
	buildData(obj)
	//根据原始信息构建原始hash
	//1.结构体转[]byte
	marshal, err := json.Marshal(obj)
	if err != nil {
		return "error", err
	}
	//2.sha256哈希
	sha256 := crypto.Sha256(marshal)
	return hex.EncodeToString(sha256), nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	obj, err := parseToflowByTx(req.Base64Tx)
	buildData(obj)
	if err != nil {
		return nil, err
	}
	var envelopeSignatures []string
	//签名转base64
	envelopeSignatures = append(envelopeSignatures, base64.StdEncoding.EncodeToString([]byte(req.Signature)))
	obj.EnvelopeSignatures = envelopeSignatures
	//转为base64
	obj.Script = base64.StdEncoding.EncodeToString([]byte(obj.Script))
	signedTx, _ := json.Marshal(obj)
	fmt.Println("obj:", obj)
	return &account.SignedTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "success",
		SignedTx: string(signedTx),
	}, nil

}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	client, err := NewFlowClient(conf)
	fmt.Println("地址 client", client)
	if err != nil {
		log.Error("Failed to create Flow client", err)
		return nil, err
	}

	return &ChainAdaptor{
		flowClient: client,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

// todo
func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	return nil, nil
}

// ValidAddress 验证地址
func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	return nil, nil
}

// ValidPublicKey 验证公钥
func (c *ChainAdaptor) ValidPublicKey(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	return nil, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	ctx := context.Background()
	data, err := c.flowClient.client.GetBlockByHeight(ctx, uint64(req.Height))

	if err != nil {
		fmt.Println("err", err)
		return &account.BlockResponse{
			Code: common.ReturnCode_SUCCESS,
			Msg:  "Do not support this rpc interface",
		}, err
	}
	fmt.Println("data数据:", data)

	return &account.BlockResponse{
		Code:         common.ReturnCode_SUCCESS,
		Height:       int64(data.Height),
		Hash:         string(data.PayloadHash),
		Transactions: make([]*account.BlockInfoTransactionList, 0),
	}, nil

}

// 构建原始哈希的数据
func buildData(message *flowMessage) {
	message.Payer = hex.EncodeToString([]byte(message.Payer))
	message.ReferenceBlockId = hex.EncodeToString([]byte(message.ReferenceBlockId))
	var author []string
	for _, authorizer := range message.Authorizers {
		author = append(author, hex.EncodeToString([]byte(authorizer)))
	}
}

func parseToflowByTx(tx string) (*flowMessage, error) {
	fmt.Println("tx:", tx)
	txReqJsonByte, err := base64.StdEncoding.DecodeString(tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, err
	}

	//json转对象
	var dynamicFeeTx flowMessage
	err = json.Unmarshal(txReqJsonByte, &dynamicFeeTx)
	if err != nil {
		log.Error("parse json fail", "err", err)
		return nil, err
	}
	return &dynamicFeeTx, nil

}

// 模拟签名机签名
func Sign(msgHash []byte, privateKey string, encrypt string) (string, error) {

	//使用secp256k1算法签名
	//todo 添加其他签名算法
	var sign []byte
	var err error
	priKey, _ := hex.DecodeString(privateKey)
	if encrypt == "secp256k1" {
		sign, err = secp256k1.Sign(msgHash, priKey)
		if err != nil {
			log.Error("sign fail", "err", err)
			return "", err
		}
	} else {
		sign, err = secp256k1.Sign(msgHash, priKey)
		if err != nil {
			log.Error("sign fail", "err", err)
			return "", err
		}
	}
	return hex.EncodeToString(sign), nil
}
