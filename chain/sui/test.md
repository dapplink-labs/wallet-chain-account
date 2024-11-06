# Test for grpc api

## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Sui",
  "network": "mainet"
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

## 2.convert address

- request
```
grpcurl -plaintext -d '{
  "chain": "Sui",
  "network": "mainnet",
  "publicKey": "f1f191fd812f91d8663822071d1de5c499483cba398aed5019b76af4137f4cc5"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
- response

```
{
  "code": "SUCCESS",
  "msg": "convert address successs",
  "address": "604b6f869a8848c53bf9b3e5a6c6caf02bffbb437e9d67a95052444d221c183e"
}
```

## 3.valid address

- request
```
grpcurl -plaintext -d '{
  "chain": "Sui",
  "network": "mainnet",
  "address": "0x604b6f869a8848c53bf9b3e5a6c6caf02bffbb437e9d67a95052444d221c183e"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.validAddress
```
- response

```
{
  "code": "SUCCESS",
  "msg": "valid address success",
  "valid": true
}
```

## 4.get block by number

## 5.get block by hash

## 6.get account
- request
```
grpcurl -plaintext -d '{
  "address": "604b6f869a8848c53bf9b3e5a6c6caf02bffbb437e9d67a95052444d221c183e",
  "chain": "Sui"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get account success",
  "network": "Sui",
  "account_number": "",
  "sequence": "",
  "balance": "0"
}
```
## 7.get fee

- request
```
grpcurl -plaintext -d '{
  "chain": "Sui"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response

```
{
  "code": "SUCCESS",
  "msg": "get gas price success",
  "slow_fee": "750|",
  "normal_fee": "750|*2",
  "fast_fee": "750|*3"
}
```

## 8.send tx
- request
```
grpcurl -plaintext -d '{
  "chain": "Sui",
  "rawTx": "6OczDiFarOkCdIpTiOkQUB+PMvMw4Z50o1eBM/f1mGAWt6SQPHnekezKDECHXcJZpb/w2ZfyHPzyG3L80+5vAQ=="
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.SendTx
```
- response

```
{
  "code": "SUCCESS",
  "msg": "send tx success",
  "tx_hash": ""
}
```

## 9.get tx by address

- request
```
grpcurl -plaintext -d '{
  "chain": "Sui",
  "address": "0x95f1baf8c250c06fc2558f2ca5b35b371977f7182d381cf29b0f36f2f9da434a",
  "cursor": "YxjRfteuVNyPfJdTf3gZD6grHjUrkTgi8pQKQZqGHyz",
  "limit": 10
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get transactions success",
  "tx": [
    {
      "hash": "2xVXiAivHGtWgwqiq36fVKFXGhyZAQiX1NsTwzXXAFBZ",
      "index": 0,
      "froms": [
        {
          "address": "0xd56948cebf0a3309e13980126bcc8ef4d7733305cd7b412fa00167d57741984e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17409569",
      "contract_address": "",
      "datetime": "1699135038425",
      "data": ""
    },
    {
      "hash": "3fUWneuY1wch7dUMFbysvYcYTw74Z7vtgWUEEUUdJx4P",
      "index": 0,
      "froms": [
        {
          "address": "0xd56948cebf0a3309e13980126bcc8ef4d7733305cd7b412fa00167d57741984e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17409576",
      "contract_address": "",
      "datetime": "1699135045121",
      "data": ""
    },
    {
      "hash": "FNKFJ4miT4DW8kYtVBM1kVrufZDJ8F1MEwiVMBQcyn6u",
      "index": 0,
      "froms": [
        {
          "address": "0x1d632d46ff70491033fefc4e6398dceaa4943dcf62512b4d57378b5ab703bc5e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17409640",
      "contract_address": "",
      "datetime": "1699135108519",
      "data": ""
    },
    {
      "hash": "ATwMPRTxTyCrDNB2xn1yM2e54a9mZmBv4x2iUmy9W2SV",
      "index": 0,
      "froms": [
        {
          "address": "0x1d632d46ff70491033fefc4e6398dceaa4943dcf62512b4d57378b5ab703bc5e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17409715",
      "contract_address": "",
      "datetime": "1699135182887",
      "data": ""
    },
    {
      "hash": "CMvmAW8RAot6DfN1HpNsrG1TXypdgMj1Bxrteydnbux3",
      "index": 0,
      "froms": [
        {
          "address": "0xd56948cebf0a3309e13980126bcc8ef4d7733305cd7b412fa00167d57741984e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17409992",
      "contract_address": "",
      "datetime": "1699135464267",
      "data": ""
    },
    {
      "hash": "8aCQvwKNgqmobDHRSxY3VPhKuztESoJWnrawCwVzHUwa",
      "index": 0,
      "froms": [
        {
          "address": "0x1d632d46ff70491033fefc4e6398dceaa4943dcf62512b4d57378b5ab703bc5e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17410362",
      "contract_address": "",
      "datetime": "1699135831326",
      "data": ""
    },
    {
      "hash": "3Baif93z1ePf44dnZUVAhH6bs5ZYmx6urYYrPJWftUQi",
      "index": 0,
      "froms": [
        {
          "address": "0x1d632d46ff70491033fefc4e6398dceaa4943dcf62512b4d57378b5ab703bc5e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17410423",
      "contract_address": "",
      "datetime": "1699135891765",
      "data": ""
    },
    {
      "hash": "zJ7PTf3v8VAZMGCsS5RHJF2FBpZZPtFoJNWRxkWNiVp",
      "index": 0,
      "froms": [
        {
          "address": "0x1d632d46ff70491033fefc4e6398dceaa4943dcf62512b4d57378b5ab703bc5e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17410478",
      "contract_address": "",
      "datetime": "1699135946202",
      "data": ""
    },
    {
      "hash": "E32gjDnUM9KnRTZeCg1mZ1qTmnuWwMprYibubFyjjAhy",
      "index": 0,
      "froms": [
        {
          "address": "0x1d632d46ff70491033fefc4e6398dceaa4943dcf62512b4d57378b5ab703bc5e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17410568",
      "contract_address": "",
      "datetime": "1699136035048",
      "data": ""
    },
    {
      "hash": "BaHiiysiHGQQyC7GYnTC82odJHanvPHeLVHA2jvWjbD3",
      "index": 0,
      "froms": [
        {
          "address": "0xd56948cebf0a3309e13980126bcc8ef4d7733305cd7b412fa00167d57741984e"
        }
      ],
      "tos": [],
      "fee": "793472",
      "status": "Success",
      "values": [],
      "type": 0,
      "height": "17410705",
      "contract_address": "",
      "datetime": "1699136171547",
      "data": ""
    }
  ]
}
```
## 10.get tx by hash

- request
```
grpcurl -plaintext -d '{
  "chain": "Sui",
  "hash": "6c286gRAis7AsBnTQYz3ons2DAiAfMj1DL4x7iPnwNfW"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
- response

```
{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "6c286gRAis7AsBnTQYz3ons2DAiAfMj1DL4x7iPnwNfW",
    "index": 0,
    "froms": [
      {
        "address": "0xd0581315160cd6d5399a5a0867f4fad6bd6d449bdc3f0c4e796f02de047b2926"
      },
      {
        "address": "0xd0581315160cd6d5399a5a0867f4fad6bd6d449bdc3f0c4e796f02de047b2926"
      }
    ],
    "tos": [],
    "fee": "82277993",
    "status": "Success",
    "values": [],
    "type": 0,
    "height": "76895198",
    "contract_address": "",
    "datetime": "1730899383395",
    "data": ""
  }
}
```

## 11.get block by range

## 12.create un sign transaction