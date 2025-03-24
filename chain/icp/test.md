## 1.Get Fee

- request
```
grpcurl -plaintext -d '{
  "chain": "ICP"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "slow_fee": "10000",
  "normal_fee": "10000",
  "fast_fee": "10000"
}

```


## 2.get support chain

- request
```
grpcurl -plaintext -d '{
  "chain": "ICP"
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


## 3.get tx list by address

- request
```
grpcurl -plaintext -d '{
  "chain": "ICP",
  "address": "e9ed7def415a8c323953f578c9ce0decf0031cb5b2a1f88cf6f1d89af80ed43a"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get txByAddress success",
  "tx": [
    {
      "hash": "a0c26b74eea80757f1e59c728abc63526efbf75a037760b4822760b04097be23",
      "index": 20911963,
      "froms": [
        {
          "address": "220c3a33f90601896e26f76fa619fe288742df1fa75426edfaf759d39f2455a5"
        }
      ],
      "tos": [
        {
          "address": "e9ed7def415a8c323953f578c9ce0decf0031cb5b2a1f88cf6f1d89af80ed43a"
        }
      ],
      "fee": "10000",
      "status": "Success",
      "values": [
        {
          "value": "76893000"
        }
      ],
      "type": 0,
      "height": "20911963",
      "contract_address": "",
      "datetime": "1740816936277092864",
      "data": "1740816935266"
    }
  ]
}

```

## 4.get tx by hash

- request
```
grpcurl -plaintext -d '{
  "hash": "e72ed0ae3b787b1aa1e45d5087c31d7cb13632b3f03b6a5dd871f9bfce2ef689",
  "chain": "ICP"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
- response
```
grpcurl -plaintext -d '{
  "chain": "ICP",
  "hash": "e72ed0ae3b787b1aa1e45d5087c31d7cb13632b3f03b6a5dd871f9bfce2ef689"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash

```

## 5.get account info 

- request
```
grpcurl -plaintext -d '{
  "chain": "ICP",
  "address": "3caf8e326ebe1916c47be71560ffd417a8b2d2692932334e0606296a7c01ab9e",
  "coin": "ICP"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get account success",
  "network": "ICP",
  "account_number": "",
  "sequence": "",
  "balance": "10000"
}
```

