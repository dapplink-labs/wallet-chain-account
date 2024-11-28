package cosmos

import (
	"context"
	"cosmossdk.io/math"
	"encoding/hex"
	"fmt"
	bftcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/log"
)

func NewTxBuilder(txStruct *TxStructure) (client.TxBuilder, client.TxConfig, error) {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	txBuilder := txConfig.NewTxBuilder()

	// coin amount
	//amount := new(big.Int).SetInt64(txStruct.Amount)
	//amount = amount.Mul(amount, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(txStruct.Decimal)), nil))
	//// fee
	//fee := new(big.Int).SetInt64(txStruct.FeeAmount)
	//fee = fee.Mul(fee, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(txStruct.Decimal)), nil))
	// set fieldx
	txBuilder.SetGasLimit(txStruct.GasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(txStruct.FeeAmount))))
	txBuilder.SetMemo(txStruct.Memo)
	fromAddr, err := sdk.AccAddressFromBech32(txStruct.FromAddress)
	if err != nil {
		log.Error("from address AccAddressFromBech32 fail", "err", err)
		return nil, nil, err
	}
	toAddr, err := sdk.AccAddressFromBech32(txStruct.ToAddress)
	if err != nil {
		log.Error("to address AccAddressFromBech32 fail", "err", err)
		return nil, nil, err
	}
	amountCoin := sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(txStruct.Amount)))
	msg := banktypes.NewMsgSend(fromAddr, toAddr, amountCoin)

	err = txBuilder.SetMsgs(msg)
	if err != nil {
		log.Error("set tx message fail", "err", err)
		return nil, nil, err
	}

	return txBuilder, txConfig, nil
}

// gasLimit=200000, feeAmount=1000, from=cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q
// to=cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2
func BuildUnSignTransaction(txStruct *TxStructure) ([]byte, error) {
	txBuilder, txConfig, err := NewTxBuilder(txStruct)
	if err != nil {
		log.Error("new tx builder fail", "err", err)
		return nil, err
	}
	pubKeyBytes, err := hex.DecodeString(txStruct.PubKey)
	if err != nil {
		log.Error("decode pubKey fail", "err", err)
		return nil, err
	}
	pubKey := &secp256k1.PubKey{Key: pubKeyBytes}
	signerData := authsigning.SignerData{
		ChainID:       txStruct.ChainId,
		AccountNumber: txStruct.AccountNumber,
		Sequence:      txStruct.Sequence,
		PubKey:        pubKey,
		Address:       txStruct.FromAddress,
	}

	signMode := signing.SignMode_SIGN_MODE_DIRECT
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: txStruct.Sequence,
	}
	sigs := []signing.SignatureV2{sig}
	if err := txBuilder.SetSignatures(sigs...); err != nil {
		return nil, err
	}

	//if err := clienttx.checkMultipleSigners(txBuilder.GetTx()); err != nil {
	//	return nil, err
	//}

	bytesToSign, err := authsigning.GetSignBytesAdapter(context.Background(), txConfig.SignModeHandler(), signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return bftcrypto.Sha256(bytesToSign), nil
}

func BuildSignTransaction(txStruct *TxStructure, signature []byte) ([]byte, error) {
	txBuilder, txConfig, err := NewTxBuilder(txStruct)
	if err != nil {
		log.Error("new tx builder fail", "err", err)
		return nil, err
	}
	signMode := signing.SignMode_SIGN_MODE_DIRECT
	pubKeyBytes, err := hex.DecodeString(txStruct.PubKey)
	if err != nil {
		log.Error("decode pubKey fail", "err", err)
		return nil, err
	}
	pubKey := &secp256k1.PubKey{Key: pubKeyBytes}
	sig := signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode:  signMode,
			Signature: signature,
		},
		Sequence: txStruct.Sequence,
	}
	err = txBuilder.SetSignatures(sig)
	// encode
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		log.Error("tx encoder fail", "err", err)
		return nil, err
	}

	fmt.Printf("signature txBytes : %X \n", txBytes)

	return txBytes, nil
}
