package cosmos

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
)

const (
	defaultConfigPath = "../../config.yml"
)

func getChainAdaptor() (chain.IChainAdaptor, error) {
	var f = flag.String("c", defaultConfigPath, "config path")
	flag.Parse()
	conf, _ := config.New(*f)
	return NewChainAdaptor(conf)
}

// success
func TestCosmos_ConvertAddress(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.ConvertAddressRequest{
		PublicKey: "032d535553c70dfbb9c13f32cb6d1002a4b421beff39009670e29a7e51fb88ec3f",
	}
	response, err := chainAdaptor.ConvertAddress(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_ValidAddress(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.ValidAddressRequest{
		Address: "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q",
	}
	response, err := chainAdaptor.ValidAddress(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_GetAccount(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.AccountRequest{
		Address: "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q",
	}
	response, err := chainAdaptor.GetAccount(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_GetBlockByNumber(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.BlockNumberRequest{
		Height: int64(22879895),
	}
	response, err := chainAdaptor.GetBlockByNumber(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetBlockByHash(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	req := &account.BlockHashRequest{
		Hash: "35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90",
	}
	response, err := chainAdaptor.GetBlockByHash(req)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetBlockHeaderByHash(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	req := &account.BlockHeaderHashRequest{
		Hash: "35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90",
	}
	response, err := chainAdaptor.GetBlockHeaderByHash(req)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_GetTxByHash(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.TxHashRequest{
		Hash: "85C84677F466D71C0BB6E744439C3040ABB35B8F2B838CC7B73CD1BFF33D0B88", //"35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90",
	}
	response, err := chainAdaptor.GetTxByHash(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetTxByAddress(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	request := &account.TxAddressRequest{
		Chain:    "Cosmos",
		Address:  "cosmos1nvcgd368m4pm5mm3ppzawhsq6grra4ejnppplx",
		Pagesize: 1,
		Page:     1,
	}

	response, err := chainAdaptor.GetTxByAddress(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetBlockByRange(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	request := &account.BlockByRangeRequest{
		Start: "22879895",
		End:   "22879896",
	}
	response, err := chainAdaptor.GetBlockByRange(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_CreateUnSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:         "cosmoshub-4",
		FromAddress:     "",
		ToAddress:       "",
		ContractAddress: "",
		Amount:          10000,
		GasLimit:        137674,
		FeeAmount:       10000,
		Sequence:        5,
		AccountNumber:   3014650,
		Decimal:         6,
		Memo:            "10086",
		PubKey:          "",
	}

	txBytes, err := json.Marshal(txStruct)
	assert.NoError(t, err)

	base64Tx := base64.StdEncoding.EncodeToString(txBytes)

	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	request := &account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  NetWork,
		Base64Tx: base64Tx,
	}

	response, err := chainAdaptor.CreateUnSignTransaction(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_BuildSignedTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:         "cosmoshub-4",
		FromAddress:     "",
		ToAddress:       "",
		ContractAddress: "",
		Amount:          10000,
		GasLimit:        137674,
		FeeAmount:       10000,
		Sequence:        5,
		AccountNumber:   3014650,
		Decimal:         6,
		Memo:            "10086",
		PubKey:          "",
	}

	txBytes, err := json.Marshal(txStruct)
	assert.NoError(t, err)

	base64Tx := base64.StdEncoding.EncodeToString(txBytes)

	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	signStr := "f95a35e4b530b25c9f81b0b572393f903b9f277fe54cf42144e861a4598e92e6647f4320c7db748979f3f796f701ce3300636af9ec7a1eab91787f7b5c88fd4a00"
	request := &account.SignedTransactionRequest{
		Chain:     ChainName,
		Network:   NetWork,
		Signature: signStr,
		Base64Tx:  base64Tx,
		PublicKey: "",
	}
	response, err := chainAdaptor.BuildSignedTransaction(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_SendTx(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	rawTx := "0a98010a8e010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126e0a2d636f736d6f73317167617338787070746e7030396c796c33326b66703630686c64676573366775753238716d6b122d636f736d6f73313974687873756e6c396c7a7977676c736e64746835613237387774617661777a7a70763434711a0e0a057561746f6d120531303030301205313030383612680a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21032d535553c70dfbb9c13f32cb6d1002a4b421beff39009670e29a7e51fb88ec3f12040a020801180512140a0e0a057561746f6d1205313030303010cab3081a40f95a35e4b530b25c9f81b0b572393f903b9f277fe54cf42144e861a4598e92e6647f4320c7db748979f3f796f701ce3300636af9ec7a1eab91787f7b5c88fd4a"
	req := &account.SendTxRequest{
		RawTx: rawTx,
	}
	response, err := chainAdaptor.SendTx(req)
	assert.NoError(t, err)
	fmt.Printf("response TxHash=%s \n", response.TxHash)
	txbytes, err := hex.DecodeString(rawTx)
	assert.NoError(t, err)
	fmt.Printf("base64 RawTx=%s\n", base64.StdEncoding.EncodeToString(txbytes))
}

func TestCosmos_validata(t *testing.T) {
	signStr := "4d8a522ce2527dff9b8c019b1809e108797fcff6a370b0b2dc2577f2a6c3458c2dbc68894046c453a554b011adeb828f199e92b66e1c2514b61bd9143040205100"
	signBytes, _ := hex.DecodeString(signStr)
	fmt.Printf("sign-1:%x\n", signBytes)
	fmt.Printf("sign-2:%x\n", signBytes[1:])
	fmt.Printf("sign-2:%s\n", signStr[:len(signStr)-2])

	rawTx := "0A98010A8E010A1C2F636F736D6F732E62616E6B2E763162657461312E4D736753656E64126E0A2D636F736D6F73317167617338787070746E7030396C796C33326B66703630686C64676573366775753238716D6B122D636F736D6F73313974687873756E6C396C7A7977676C736E64746835613237387774617661777A7A70763434711A0E0A057561746F6D120531303030301205313030383712670A500A460A1F2F636F736D6F732E63727970746F2E736563703235366B312E5075624B657912230A21032D535553C70DFBB9C13F32CB6D1002A4B421BEFF39009670E29A7E51FB88EC3F12040A020801180412130A0D0A057561746F6D12043130303010CAB3081A404D8A522CE2527DFF9B8C019B1809E108797FCFF6A370B0B2DC2577F2A6C3458C2DBC68894046C453A554B011ADEB828F199E92B66E1C2514B61BD91430402051"
	fmt.Printf("rawTx[1:]=%s\n", rawTx[1:])
}
