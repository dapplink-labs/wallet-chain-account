## 1.get free

- request
```
grpcurl -plaintext -d '{
  "chain": "Scroll"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get gas price success",
  "slow_fee": "39932311|100",
  "normal_fee": "39932311|100|*2",
  "fast_fee": "39932311|100|*3"
}

```


## 2.get support chain

- request
```
grpcurl -plaintext -d '{
  "chain": "Scroll"
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
  "chain": "Scroll",
  "address": "0x4640531c3A8E6C575A4cA2890f4032844123fA33",
  "contractAddress": "0xb0643F7b3e2E2F10FE4e38728a763eC05f4ADeC3"
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
  "chain": "Scroll",
  "hash": "0x2001ed0c6416bfb072038186bb83de4ee63569ab0d5b1487a5c4c2b4f83ac9c7"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get tx by hash success",
  "tx": {
    "hash": "0x2001ed0c6416bfb072038186bb83de4ee63569ab0d5b1487a5c4c2b4f83ac9c7",
    "index": 0,
    "froms": [
      {
        "address": ""
      }
    ],
    "tos": [
      {
        "address": "0x8218a0f47f4c0de0c1754f50874707cd6e7b2e5e"
      }
    ],
    "fee": "239726923",
    "status": "Success",
    "values": [
      {
        "value": "500000000000000000"
      }
    ],
    "type": 0,
    "height": "13508886",
    "contract_address": "0xb0643F7b3e2E2F10FE4e38728a763eC05f4ADeC3",
    "datetime": "",
    "data": "A9059CBB0000000000000000000000008218A0F47F4C0DE0C1754F50874707CD6E7B2E5E00000000000000000000000000000000000000000000000006F05B59D3B20000"
  }
}

```

## 5.get account info 

- request
```
grpcurl -plaintext -d '{
  "chain": "Scroll",
  "address": "0x4640531c3A8E6C575A4cA2890f4032844123fA33",
  "contractAddress": "0xb0643F7b3e2E2F10FE4e38728a763eC05f4ADeC3"
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

