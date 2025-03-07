## 1.get free

- request
```
grpcurl -plaintext -d '{
  "chain": "Btt"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get gas price success",
  "slow_fee": "9000000000000000|9000000000000000",
  "normal_fee": "9000000000000000|9000000000000000|*2",
  "fast_fee": "9000000000000000|9000000000000000|*3"
}

```


## 2.get support chain

- request
```
grpcurl -plaintext -d '{
  "chain": "Btt"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getSupportChains
```
- response
```
{
  "code": "SUCCESS",
  "msg": "Support Chain",
  "support": true
}

```


## 3.get tx list by address

- request
```
grpcurl -plaintext -d '{
  "chain": "Btt",
  "address": "0xDDA22000e1bCC0c70C8b1947CE7074df1DC5B80B",
  "contractAddress": "0x1Bdd8878252DaddD3Af2ba30628813271294eDc0"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get tx list success",
  "tx": []
}

```

## 4.get tx by hash

- request
```
grpcurl -plaintext -d '{
  "chain": "Btt",
  "hash": "0xfe66799cd6de5b8a6a9657bf91cb64101d8c0f511b52ab644b43bb92688d2a26"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get tx by hash success",
  "tx": {
    "hash": "0xfe66799cd6de5b8a6a9657bf91cb64101d8c0f511b52ab644b43bb92688d2a26",
    "index": 6,
    "froms": [
      {
        "address": ""
      }
    ],
    "tos": [
      {
        "address": "0x5ac40eb0dcd19a64fef09fa8a30e6ae2dd3f3afb"
      }
    ],
    "fee": "27000001",
    "status": "Success",
    "values": [
      {
        "value": "104000000"
      }
    ],
    "type": 0,
    "height": "75827131",
    "contract_address": "0x201EBa5CC46D216Ce6DC03F6a759e8E766e956aE",
    "datetime": "",
    "data": "A9059CBB0000000000000000000000005AC40EB0DCD19A64FEF09FA8A30E6AE2DD3F3AFB000000000000000000000000000000000000000000000000000000000632EA00"
  }
}

```

## 5.get account info 

- request
```
grpcurl -plaintext -d '{
  "chain": "Btt",
  "address": "0xDDA22000e1bCC0c70C8b1947CE7074df1DC5B80B",
  "contractAddress": "0x1Bdd8878252DaddD3Af2ba30628813271294eDc0"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "code": "SUCCESS",
  "msg": "",
  "network": "",
  "account_number": "",
  "sequence": "2",
  "balance": "0"
}

```

