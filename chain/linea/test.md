## 1.get free

- request
```
grpcurl -plaintext -d '{
  "chain": "Linea"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get gas price success",
  "slow_fee": "103328888|104889419",
  "normal_fee": "103328888|104889419|*2",
  "fast_fee": "103328888|104889419|*3"
}

```


## 2.get support chain

- request
```
grpcurl -plaintext -d '{
  "chain": "Linea"
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
  "chain": "Linea",
  "address": "0x4640531c3A8E6C575A4cA2890f4032844123fA33",
  "contractAddress": "0x7da14988E4f390C2E34ed41DF1814467D3aDe0c3"
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
  "chain": "Linea",
  "hash": "0x2fce54eeed61bd83eacf7be8a8fadeccdca5dcfbee29831a016e96761b580f6f"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get tx by hash success",
  "tx": {
    "hash": "0x2fce54eeed61bd83eacf7be8a8fadeccdca5dcfbee29831a016e96761b580f6f",
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
    "fee": "293977842",
    "status": "Success",
    "values": [
      {
        "value": "9000000000000000000"
      }
    ],
    "type": 0,
    "height": "15979981",
    "contract_address": "0x7da14988E4f390C2E34ed41DF1814467D3aDe0c3",
    "datetime": "",
    "data": "A9059CBB0000000000000000000000008218A0F47F4C0DE0C1754F50874707CD6E7B2E5E0000000000000000000000000000000000000000000000007CE66C50E2840000"
  }
}

```

## 5.get account info 

- request
```
grpcurl -plaintext -d '{
  "chain": "Linea",
  "address": "0x4640531c3A8E6C575A4cA2890f4032844123fA33",
  "contractAddress": "0x7da14988E4f390C2E34ed41DF1814467D3aDe0c3"
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

