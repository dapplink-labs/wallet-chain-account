# Test for grpc api

## getSupportChains
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm"
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

## validAddress
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "address": "GAZEFFEFXCG2IOM7QBICUKCO3MOL3NVT3GORVBNS7TTETRQCQDYXPOQC"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.validAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "ValidAddress Success",
  "valid": true
}
```

## getAccount
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "address": "GDYDI34YZQCP7WU726B626KTJIE6COPSXMWL2VLR7KNRHMNB6HLNFUJB"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "code": "SUCCESS",
  "msg": "GetAccountInfo Success",
  "network": "mainnet",
  "account_number": "GDYDI34YZQCP7WU726B626KTJIE6COPSXMWL2VLR7KNRHMNB6HLNFUJB",
  "sequence": "239763383908302850",
  "balance": "15.8899801"
}
```

## getFee
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "code": "SUCCESS",
  "msg": "GetFee Success",
  "slow_fee": "100",
  "normal_fee": "100",
  "fast_fee": "200"
}
```

## getBlockByNumber
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "height": "56013360"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "GetBlockByNumber Success",
  "height": "56013360",
  "hash": "0de2a98b6f1f0381c1ad3b5e3d20fc19dd31e45cc36de9889defe570f55715eb",
  "base_fee": "not support",
  "transactions": [
    {
      "from": "GCVTUT4QMMCZFKOQTXFRHWFGBT6B4V3GUECDE6OVNC6HFKSPOBHHTBHI",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0de2a98b6f1f0381c1ad3b5e3d20fc19dd31e45cc36de9889defe570f55715eb",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GAGIOOF5O3MICGPDOM23FIMP6RYBAB4VHD27OALY4SDPMH6S6LTQL3HN",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "dc978859dc7c1579ef0a959e63f332b01171c155b019d1aad9cd9ad6d89cd1c2",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GBJN2P72VGBC2VXATZ63YOF3ZVCJW7VFG2TXKSGMMKONLLHX4IH4G3AF",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "69a8e067caf7d30d045318b3246aa7af96732a234264150799fc29c231a23476",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GC3JWRE2XYGKHIDIUJDWHJJRBZ5HX4BAIBNYYTOS3DV6K7FSENZEBRZT",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "072f6ccf882fa02b5638db1db2cdb338c3b6cd6840939a77210b06a7a19ed749",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GB2LUKAJ2MV3DS7Y5Q652J547T27EA2J7N2VL566CRYULXC24A6HIHQB",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0fa666bab32622fd2476d47c48fe7f19b886c76dd5e51f341ce446d577b1ceb7",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GAABED3SGW3GU2LPUILFPLAZJTGAO56BWPN2JH6X6KZZV4DW7VYQYAMR",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "5465a2583291d5e51f2d7a923c5ede581b51ee18b6032334507c65b12645d547",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GBDB7AWXKFGXJL347Y2BBL5FZXVVZ2AN5TP6WA473R6MO4K42PL3IBHZ",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "c06d0e2b73c7cc46c8ccddf1f8ff3f7f79ca52ec748692ea0856f28ef5e7335f",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GBITCV7JEVGP73NSBZ4X2ZCW7FNEFIBN5IO7MY2F7LIA5SAFI67FEBMX",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "be9d028ee151d9a0179ab524d364286832883d4874ae6ce209189c2687a7c23d",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GCZU4U4Q6Y74CKEW3CIR4AY5OUIKLXEXWAMMFJJI75VZ5G2O4LOMSSXW",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "c8da7f0f5d98d0c2a6cf0933b020751ef7303d67ce2d5da1e9cbc1ee048ec4e1",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    },
    {
      "from": "GD2QHQLCXFDYR7IQRP76YPRNQG5W676QMWC4UUJ5QLRHPD636N2IKH7S",
      "to": "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
      "token_address": "",
      "contract_wallet": "",
      "hash": "c7cb04d22aa5717eeb259d9a127372739000b821071f3b8405b1ab30e19af586",
      "height": "56013360",
      "amount": "Not Support In This Function(Please Use GetTransactionByHash For Detail)"
    }
  ]
}
```

## getBlockHeaderByNumber
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "height": "56013360"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "GetBlockHeaderByNumber Success",
  "block_header": {
    "hash": "36715cbd595fc4bd1c797c0d2ed92b62ec5add3202f6e42f3b4fa7cd2fdd9f00",
    "parent_hash": "c31da262a9e205a05eed437666a45ba0c0e32bc9bee44767bf5464a062dca46d",
    "uncle_hash": "",
    "coin_base": "105443902087.3472865",
    "root": "",
    "tx_hash": "",
    "receipt_hash": "",
    "parent_beacon_root": "",
    "difficulty": "",
    "number": "56013360",
    "gas_limit": "0",
    "gas_used": "0",
    "time": "0",
    "extra": "",
    "mix_digest": "",
    "nonce": "",
    "base_fee": "100",
    "withdrawals_hash": "",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}
```

## getTxByHash
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "hash": "47476b985c63ee571505048c179a79226e0968ca35dca0f0c9a58968bddafc6b"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "GetTransactionByHash Success",
  "tx": {
    "hash": "47476b985c63ee571505048c179a79226e0968ca35dca0f0c9a58968bddafc6b",
    "index": 0,
    "froms": [
      {
        "address": "GDYDI34YZQCP7WU726B626KTJIE6COPSXMWL2VLR7KNRHMNB6HLNFUJB"
      }
    ],
    "tos": [
      {
        "address": "GD6ARGMYT65UUC7FQDBK77GXMEONL44BL7E5G4WL2NDWMJ7NSWBUBYQQ"
      }
    ],
    "fee": "100",
    "status": "Success",
    "values": [
      {
        "value": "5.0050000"
      }
    ],
    "type": 0,
    "height": "56013159",
    "contract_address": "Sorry, is currently not supported...",
    "datetime": "",
    "data": ""
  }
}
```

## createUnSignTransaction
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "base64Tx": "eyJhZGRyRnJvbSI6IkdEWURJMzRZWlFDUDdXVTcyNkI2MjZLVEpJRTZDT1BTWE1XTDJWTFI3S05SSE1OQjZITE5GVUpCIiwiYWRkclRvIjoiR0Q2QVJHTVlUNjVVVUM3RlFEQks3N0dYTUVPTkw0NEJMN0U1RzRXTDJORFdNSjdOU1dCVUJZUVEiLCJzZXF1ZW5jZUZyb20iOjIzOTc2MzM4MzkwODMwMjg1MywiYW1vdW50IjoiMC4xMjMifQ=="
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.createUnSignTransaction
```
- response
```
{
  "code": "SUCCESS",
  "msg": "CreateUnsignTransaction Success",
  "un_sign_tx": "AAAAAgAAAADwNG+YzAT/2p/Xg+15U0oJ4Tnyuyy9VXH6mxOxofHW0gAAAGQDU8+HAAAABgAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPA0b5jMBP/an9eD7XlTSgnhOfK7LL1VcfqbE7Gh8dbSAAAAAQAAAAD8CJmYn7tKC+WAwq/812Ec1fOBX8nTcsvTR2Yn7ZWDQAAAAAAAAAAAABLEsAAAAAAAAAAA"
}
```

## buildSignedTransaction
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "base64Tx": "AAAAAgAAAADwNG+YzAT/2p/Xg+15U0oJ4Tnyuyy9VXH6mxOxofHW0gAAAGQDU8+HAAAABgAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPA0b5jMBP/an9eD7XlTSgnhOfK7LL1VcfqbE7Gh8dbSAAAAAQAAAAD8CJmYn7tKC+WAwq/812Ec1fOBX8nTcsvTR2Yn7ZWDQAAAAAAAAAAAABLEsAAAAAAAAAAA",
  "signature": "rf+1IJ/1jNpl+2B7bveES3Cd1bCIPu55u3DN4qHTcawhJ/rNLdZkosgf09NhBYHWG2SfrJF3k4Knly57EKavAg==",
  "publicKey": "GDYDI34YZQCP7WU726B626KTJIE6COPSXMWL2VLR7KNRHMNB6HLNFUJB"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.buildSignedTransaction
```
- response
```
{
  "code": "SUCCESS",
  "msg": "SignedTransaction Success",
  "signed_tx": "AAAAAgAAAADwNG+YzAT/2p/Xg+15U0oJ4Tnyuyy9VXH6mxOxofHW0gAAAGQDU8+HAAAABgAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPA0b5jMBP/an9eD7XlTSgnhOfK7LL1VcfqbE7Gh8dbSAAAAAQAAAAD8CJmYn7tKC+WAwq/812Ec1fOBX8nTcsvTR2Yn7ZWDQAAAAAAAAAAAABLEsAAAAAAAAAABofHW0gAAAECt/7Ugn/WM2mX7YHtu94RLcJ3VsIg+7nm7cM3iodNxrCEn+s0t1mSiyB/T02EFgdYbZJ+skXeTgqeXLnsQpq8C"
}
```

## SendTx
- request
```
grpcurl -plaintext -d '{
  "chain": "Xlm",
  "rawTx": "AAAAAgAAAADwNG+YzAT/2p/Xg+15U0oJ4Tnyuyy9VXH6mxOxofHW0gAAAGQDU8+HAAAABgAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPA0b5jMBP/an9eD7XlTSgnhOfK7LL1VcfqbE7Gh8dbSAAAAAQAAAAD8CJmYn7tKC+WAwq/812Ec1fOBX8nTcsvTR2Yn7ZWDQAAAAAAAAAAAABLEsAAAAAAAAAABofHW0gAAAECt/7Ugn/WM2mX7YHtu94RLcJ3VsIg+7nm7cM3iodNxrCEn+s0t1mSiyB/T02EFgdYbZJ+skXeTgqeXLnsQpq8C"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.SendTx
```
- response
```
{
  "code": "SUCCESS",
  "msg": "SendTx PENDING",
  "tx_hash": "182fdbeccdf7d16be2616992b20a73c198499bc7f9c0917b65a34a76a15755e6"
}
```
