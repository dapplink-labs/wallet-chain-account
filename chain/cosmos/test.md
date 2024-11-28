# Test for grpc api

## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Cosmos",
  "network": "mainnet"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getSupportChains
```
- response
```
{
  "code": "SUCCESS",
  "msg": "Support this chain",
  "support": true
}
```

## 3.valid address

- request
```
grpcurl -plaintext -d '{
  "chain": "Cosmos",
  "network": "mainnet",
  "address": "cosmos1z79jxnsw64c20upyfu8rfe89pdsel48kfmzjgu"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.validAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "valid address",
  "valid": true
}
```


## latest block header by number

- request
```
grpcurl -plaintext -d '{
  "chain": "Cosmos",
  "network": "mainnet",
  "height": "22879895"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```
- reponse
```
{
  "code": "SUCCESS",
  "msg": "get block header by number success",
  "block_header": {
    "hash": "35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90",
    "parent_hash": "34BBABD339D470E8191E9995DDBC17E801A09BC404A2A7C1F473BA72FF01372F",
    "uncle_hash": "",
    "coin_base": "",
    "root": "",
    "tx_hash": "775338910D85CE6783949F84E8038ABE90C931BE1DF26B287C368C9E235A87B6",
    "receipt_hash": "",
    "parent_beacon_root": "",
    "difficulty": "",
    "number": "22879895",
    "gas_limit": "0",
    "gas_used": "0",
    "time": "1730467750",
    "extra": "",
    "mix_digest": "",
    "nonce": "",
    "base_fee": "",
    "withdrawals_hash": "",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}
```

## get account 

- request
```
grpcurl -plaintext -d '{
  "chain": "Cosmos",
  "network": "mainnet",
  "address": "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get account success",
  "network": "mainnet",
  "account_number": "2424228",
  "sequence": "5",
  "balance": "2424228"
}
```