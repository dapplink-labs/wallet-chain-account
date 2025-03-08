package ton

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
)

func getUserFriendly(addressBook map[string]struct {
	UserFriendly string `json:"user_friendly"`
}, key string) string {
	if entry, ok := addressBook[key]; ok {
		return entry.UserFriendly
	} else {
		return ""
	}
}

func stringToInt(amount string) *big.Int {
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}

func ParseTxMessage(ret *Tx, tx *Transactions) (*account.TxMessage, error) {
	var fromAddr string
	var toAddr string
	totalAmount := big.NewInt(0)

	fromAddr = tx.InMsg.Source
	toAddr = getUserFriendly(ret.AddressBook, tx.InMsg.Destination)

	if len(tx.InMsg.Value) > 0 {
		totalAmount = new(big.Int).Add(totalAmount, stringToInt(tx.InMsg.Value))
	}

	for _, out := range tx.OutMsgs {
		if len(out.Source) > 0 {
			fromAddr = getUserFriendly(ret.AddressBook, out.Source)
		}
		if len(out.Destination) > 0 {
			toAddr = getUserFriendly(ret.AddressBook, out.Destination)
		}
		log.Info(totalAmount.String(), "value", out.Value)
		if len(out.Value) > 0 {
			totalAmount = new(big.Int).Sub(totalAmount, stringToInt(out.Value))
		}
	}

	txMsg := &account.TxMessage{
		Hash:     tx.Hash,
		From:     fromAddr,
		To:       toAddr,
		Fee:      tx.TotalFees,
		Value:    totalAmount.String(),
		Status:   account.TxStatus_Success,
		Datetime: strconv.Itoa(tx.Now),
		Height:   strconv.Itoa(tx.BlockRef.Seqno),
	}

	return txMsg, nil
}
