package binance

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"math/big"
)

func SignHash(hash string, privateKey string) (signature string, err error) {
	bytes, err := hexutil.Decode(hash)
	if err != nil {
		panic(err)
	}
	prkByte, err := hex.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
	prk, err := crypto.ToECDSA(prkByte)
	if err != nil {
		panic(err)
	}
	sig, err := crypto.Sign(bytes, prk)
	if err != nil {
		panic(err)
	}
	signature = hex.EncodeToString(sig)
	return signature, nil
}

func BuildErc20Data(toAddress common.Address, amount *big.Int) []byte {
	var data []byte

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodId := hash[:4]
	dataAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	dataAmount := common.LeftPadBytes(amount.Bytes(), 32)

	data = append(data, methodId...)
	data = append(data, dataAddress...)
	data = append(data, dataAmount...)

	return data
}

func BuildErc721Data(fromAddress, toAddress common.Address, tokenId *big.Int) []byte {
	var data []byte

	transferFnSignature := []byte("safeTransferFrom(address,address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodId := hash[:4]

	dataFromAddress := common.LeftPadBytes(fromAddress.Bytes(), 32)
	dataToAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	dataTokenId := common.LeftPadBytes(tokenId.Bytes(), 32)

	data = append(data, methodId...)
	data = append(data, dataFromAddress...)
	data = append(data, dataToAddress...)
	data = append(data, dataTokenId...)

	return data
}

func CreateEip1559UnSignTx(txData *types.DynamicFeeTx, chainId *big.Int) (string, error) {
	tx := types.NewTx(txData)
	// 签名者
	signer := types.LatestSignerForChainID(chainId)
	txHash := signer.Hash(tx)
	return txHash.String(), nil
}

func CreateEip1559SignedTx(txData *types.DynamicFeeTx, signature []byte, chainId *big.Int) (types.Signer, *types.Transaction, string, string, error) {
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)
	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return nil, nil, "", "", errors.New("tx with signature fail")
	}
	signedTxData, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, nil, "", "", errors.New("encode tx to byte fail")
	}
	return signer, signedTx, "0x" + hex.EncodeToString(signedTxData)[4:], signedTx.Hash().String(), nil
}

func CreateLegacyUnSignTx(txData *types.LegacyTx, chainId *big.Int) string {
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)
	txHash := signer.Hash(tx)
	return txHash.String()
}

func CreateLegacySignedTx(txData *types.LegacyTx, signature []byte, chainId *big.Int) (string, string, error) {
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(chainId)
	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return "", "", errors.New("tx with signature fail")
	}
	signedTxData, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", "", errors.New("encode tx to byte fail")
	}
	return "0x" + hex.EncodeToString(signedTxData), signedTx.Hash().String(), nil
}
