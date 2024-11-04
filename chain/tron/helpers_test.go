package tron

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"testing"
)

func TestKeyPair(t *testing.T) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		fmt.Println(hex.EncodeToString(privKey.PubKey().SerializeCompressed()))
	}
}
