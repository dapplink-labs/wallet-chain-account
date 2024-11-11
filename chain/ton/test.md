## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Ton",
  "network": ""
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

## 2.get fee
- request
```
grpcurl -plaintext -d '{
  "chain": "Ton",
  "coin": "Ton",
  "network": "mainnet",
  "address": "UQAmR3oackYYWKyLCXMPktBt6i1YdgqjvEB1h-z_4fxIdMxh",
  "rawTx": "te6cckEBAgEAsQAB34gBZxgy0BTErESwihTr71KC90UqEBEH/jYfv8jC0d5UK9IAwvpgFQ1g48Pa/ufkFQhxd1Em3Huto1ZjG1f/4DmajMCAvHCRDoTZG+muS6a2eyakWR8G3+SkHKRVgPQ7FEegYU1NGLs5YKa4E5ZuKBwBAHhCAEgEWnejouwnkZpdjD+P5IWAj3d3f75SIUyVS3V5htS9IdzWUAAAAAAAAAAAAAAAAAAAAAAAAG1lbW8ILCE7"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "slow_fee": "",
  "normal_fee": "1030392",
  "fast_fee": ""
}
```

## 3.get account
- request
```
grpcurl -plaintext -d '{
  "chain": "Ton",
  "coin": "Ton",
  "network": "mainnet",
  "address": "UQBX63RAdgShn34EAFMV73Cut7Z15lUZd1hnVva68SEl7pGn"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get account info success",
  "network": "",
  "account_number": "",
  "sequence": "50648836000005",
  "balance": "45005.451813553"
}
```


## 4.get tx by address
- request
```
grpcurl -plaintext -d '{
  "chain": "Ton",
  "coin": "Ton",
  "network": "mainnet",
  "address": "UQBX63RAdgShn34EAFMV73Cut7Z15lUZd1hnVva68SEl7pGn",
  "contractAddress": "0x00",
  "page": 1,
  "pagesize": 10
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get transactions fail",
  "tx": [
    {
      "hash": "dWndYLJgtTYcHEX+E23/n3qvT905Amx6Hw9ZgFjOBvk=",
      "index": 0,
      "froms": [
        {
          "address": "EQBSkosi3wGpHY8pVRQoCceFKiQdNWUrTRpEbPe1HstEAaWE"
        }
      ],
      "tos": [
        {
          "address": "UQBX63RAdgShn34EAFMV73Cut7Z15lUZd1hnVva68SEl7pGn"
        }
      ],
      "fee": "2",
      "status": "Success",
      "values": [
        {
          "value": "1"
        }
      ],
      "type": 0,
      "height": "46948661",
      "contract_address": "",
      "datetime": "1730946841",
      "data": ""
    }
  ]
}
```


## 5.get tx by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Ton",
  "coin": "Ton",
  "network": "mainet",
  "hash": "c5hQxGQPrAMzj38xr/lyx4n09eb1V/l4eeYxHjqjSjI="
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get transaction by hash success",
  "tx": {
    "hash": "c5hQxGQPrAMzj38xr/lyx4n09eb1V/l4eeYxHjqjSjI=",
    "index": 0,
    "froms": [
      {
        "address": ""
      },
      {
        "address": "UQDymFdnA3EOhpTTzkTL4WCfOSpmczY9P3SLqzMDeSqftSfF"
      }
    ],
    "tos": [
      {
        "address": "UQDymFdnA3EOhpTTzkTL4WCfOSpmczY9P3SLqzMDeSqftSfF"
      },
      {
        "address": "EQCHPJRdUbvIgdFKII1WhQdV3yH6qHbTtuu22O9x1dvFGvqY"
      }
    ],
    "fee": "3997465",
    "status": "Success",
    "values": [
      {
        "value": "-298710400"
      }
    ],
    "type": 0,
    "height": "46910081",
    "contract_address": "",
    "datetime": "1730946554",
    "data": ""
  }
}
```


## 5.send raw tx 
- request
```

```

- response
```

```
