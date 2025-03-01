## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Tezos"
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
  "chain": "Tezos",
  "rawTx": "[\n  {\"branch\":\"BLdSgBG7kCne5cHNrmBRdHwhHhmo9vbPAJQEtnSqPij7tZGUR3N\",\"protocol\":\"PsQuebecnLByd3JwTiGadoG4nGWi3HYiLXUjkibeFV8dCFeVMUg\",\"contents\":[{\"kind\":\"transaction\",\"source\":\"tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x\",\"fee\":\"50000\",\"counter\":\"156726390\",\"gas_limit\":\"20000\",\"storage_limit\":\"0\",\"destination\":\"tz1iPCKqRmQum2mMYpxu77WVQ1w7NxG1xn2M\",\"amount\":\"10000\"}],\"signature\":\"edsigu1xkUDNEgvHsf7bxGYebsedsti3cpjTSmEDd1ukmpgTampmb9Z9FWrBVr9HUuX56cV62AaBPWds1aU1rgWRSTEUwTYa6gD\"}\n\n\n]"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "slow_fee": "",
  "normal_fee": "168821",
  "fast_fee": ""
}
```
## 3.get account
- request
```
grpcurl -plaintext -d '{
  "chain": "Tezos",
  "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get account success",
  "network": "",
  "account_number": "5585967",
  "sequence": "156726389",
  "balance": "924033"
}
```

## 4.get tx by address
- request
```
grpcurl -plaintext -d '{
  "chain": "Tezos",
  "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get transaction by address success",
  "tx": [
    {
      "hash": "oocCY2DpzYxwKM7k8mAo4jozprKFqRD8eVKP9Drqk7S3E7GYk5G",
      "index": 156726389,
      "froms": [
        {
          "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
        }
      ],
      "tos": [
        {
          "address": "tz1iPCKqRmQum2mMYpxu77WVQ1w7NxG1xn2M"
        }
      ],
      "fee": "50000",
      "status": "Success",
      "values": [
        {
          "value": "10000"
        }
      ],
      "type": 0,
      "height": "8115585",
      "contract_address": "",
      "datetime": "2025-03-01T08:42:00Z",
      "data": ""
    },
    {
      "hash": "onig6ZmLpkAKmryumCgV1jPCEqkSPWQXXWA3Nde3MFvsDDdJRom",
      "index": 156726388,
      "froms": [
        {
          "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
        }
      ],
      "tos": [
        {
          "address": "tz1iPCKqRmQum2mMYpxu77WVQ1w7NxG1xn2M"
        }
      ],
      "fee": "50000",
      "status": "Success",
      "values": [
        {
          "value": "10000"
        }
      ],
      "type": 0,
      "height": "8011005",
      "contract_address": "",
      "datetime": "2025-02-19T14:40:28Z",
      "data": ""
    },
    {
      "hash": "ooFYCt4iCMf1bdKYV9TcyQV61mWLEvAafGcpBzqc5fcAgrQFCRQ",
      "index": 156726387,
      "froms": [
        {
          "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
        }
      ],
      "tos": [
        {
          "address": "tz1iPCKqRmQum2mMYpxu77WVQ1w7NxG1xn2M"
        }
      ],
      "fee": "64537",
      "status": "Success",
      "values": [
        {
          "value": "100"
        }
      ],
      "type": 0,
      "height": "8008447",
      "contract_address": "",
      "datetime": "2025-02-19T08:58:12Z",
      "data": ""
    },
    {
      "hash": "onj2uzpHCdAU2kqHMLihZobH9CWQQwSjsnGJkCB3iDv764utPwm",
      "index": 1165234,
      "froms": [
        {
          "address": "tz1VeaJWkdr2m5YaKFgeAafenhHhDcMxWfHC"
        }
      ],
      "tos": [
        {
          "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
        }
      ],
      "fee": "68750",
      "status": "Success",
      "values": [
        {
          "value": "1109001"
        }
      ],
      "type": 0,
      "height": "8000886",
      "contract_address": "",
      "datetime": "2025-02-18T16:03:48Z",
      "data": ""
    }
  ]
}
```

## 5.send raw tx
- request
```
grpcurl -plaintext -d '{
  "chain": "Tezos",
  "rawTx": "78b56aa434a132034aa3cb7ca067c11e51829af04f32215e0e18c8f27d368edb6c00c0b65ae23baa92db49173f436e6d9b48d7a9c4c9d08603f5e8dd4aa09c0100904e0000f97cb1e260da139b3a2ae4416df6adaa93eee16a000513955b9ba5695078a84ac9d3b672d532455d549ada6997eb457ce30e14c0c01c301c17d2290eab58d0ab2cf197d5a93af5ddcc91999aa77a00f0bd0a1cb10e"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.SendTx
```
- response
```
{
"code": "SUCCESS",
"msg": "send raw transaction success",
"tx_hash": "oocCY2DpzYxwKM7k8mAo4jozprKFqRD8eVKP9Drqk7S3E7GYk5G"
}
```
## 6.Convert Address
-request
```
grpcurl -plaintext -d '{
  "chain": "Tezos",
  "publicKey": "edpku1fgxoHsZyhfyos1tDGWVrvbCTtv9Ubi78Ah11Q8KiSr745v3m"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
-response
```
{
  "code": "SUCCESS",
  "msg": "convert address successs",
  "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
}
```
## 7.ValidAddress
-request
```
grpcurl -plaintext -d '{
  "chain": "Tezos",
  "address": "tz1dCztCWFnbvX66Tnyp8R6QB8k6azrD7i3x"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.validAddress
```
-reponse
```
{
  "code": "SUCCESS",
  "msg": "convert address successs",
  "valid": true
}
```
## 8.Get Block By Number
-request
```
grpcurl -plaintext -d '{
  "chain": "Tezos",
  "height": "81000"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
-response
```
{
  "code": "SUCCESS",
  "msg": "get block number success",
  "height": "81000",
  "hash": "BKzLgjXmP18Rwpmdgk8cVkpNPvCjSRpSXNZaTp8FuUtmNjhYQDM",
  "base_fee": "",
  "transactions": []
}
```
## 9.Get Block By Hash
-request
```
grpcurl -plaintext -d '{
  "chain": "Tezos",
  "hash": "BLdSgBG7kCne5cHNrmBRdHwhHhmo9vbPAJQEtnSqPij7tZGUR3N"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByHash
```
-response
```
{
  "code": "SUCCESS",
  "msg": "get block number success",
  "height": "8115474",
  "hash": "BLdSgBG7kCne5cHNrmBRdHwhHhmo9vbPAJQEtnSqPij7tZGUR3N",
  "base_fee": "",
  "transactions": []
}
```