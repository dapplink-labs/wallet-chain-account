# Test for grpc api

## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
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

## 2.convert address

- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "network": "mainnet",
  "publicKey": "02e993166ac8fb56c438a2a0e1266f33b54dfe7b79f738d9945dbbbebf6e367c55"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
- reponse

```
{
  "code": "SUCCESS",
  "msg": "convert address success",
  "address": "0x2ec57B631580dF40d1E9e027360357eb61C7B25A"
}
```

## 3.valid address

- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "network": "mainnet",
  "address": "0x2ec57B631580dF40d1E9e027360357eb61C7B25A"
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
  "chain": "Ethereum",
  "network": "mainnet",
  "height": "0"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```
- reponse
```
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0xc8d18e5d4774c8015afa3421c896fa3c62caed9e95a996888951050d52d907f1",
    "parent_hash": "0x65c9b549a39fe009e30da2e4b8254f110807716ee13aace39d311db047cb99f6",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
    "root": "0xe8767654cb614145c913554ee37a6339345a8aed99828073494e51c78ce4dfa7",
    "tx_hash": "0x65b08324c07c9eddf9a9985dc3f39992d221ed5a8d285c06beee3b146dc5b11f",
    "receipt_hash": "0x8c07c03f43a36f20408ab979b113f4a57656b9b3e6c44f678b215fa9c8506f96",
    "parent_beacon_root": "0x1509286f7f1e87edc0a203e7404c0dc7baaa636b1393803de57ee43288649ed1",
    "difficulty": "0",
    "number": "21105467",
    "gas_limit": "30000000",
    "gas_used": "11948296",
    "time": "1730617067",
    "extra": "beaverbuild.org",
    "mix_digest": "0xbe1d6a7b7b2db6c20c76476b838c66f37230f7fcf7a8f0c2c56564d620d68d49",
    "nonce": "0",
    "base_fee": "3103331426",
    "withdrawals_hash": "0x8bcc937c30343dc790269e00894c3d8ef49de5df4fdb64a167b12ad8e5632495",
    "blob_gas_used": "655360",
    "excess_blob_gas": "67239936"
  }
}

{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0xb89c8992dea2f3f9dc48fd876b6f7dd00a746d485bf2c02330bd2bda50f98558",
    "parent_hash": "0x6e93a5274fa198be164f46760c95066518d10d47e73f3d6861d93b02ba556b29",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
    "root": "0x0a22f243ddb4b9e5251cd1ad966e9170864d7efb7462107683df06e1523ffe79",
    "tx_hash": "0xc999f11db91fc3398fb4f5c2d9cf706d75b4d40ecdc2a84851e5a7768ac7b739",
    "receipt_hash": "0x9f8c36d852fae8cffccfec05c5e84939c6a88e0cf63e7dc3963e5657d9ddc09e",
    "parent_beacon_root": "0x5f702b3489d2683f9d325e85fe2919d7c4d08198e6fa372d6b7f3d42d8834fa4",
    "difficulty": "0",
    "number": "21105407",
    "gas_limit": "30000000",
    "gas_used": "15379455",
    "time": "1730616347",
    "extra": "beaverbuild.org",
    "mix_digest": "0xe53ea84af490d6d4a7169dfa9dbe00ecec1a7b6afb5c53a1035d60bf4d799fd4",
    "nonce": "0",
    "base_fee": "3384783595",
    "withdrawals_hash": "0xbea8080e76aad28e92218ba212b9c8ccb21d26d0ea3a357e6d02866e7a9de06e",
    "blob_gas_used": "393216",
    "excess_blob_gas": "67371008"
  }
}
```

## block header by hash

- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "network": "mainnet",
  "hash": "0xb89c8992dea2f3f9dc48fd876b6f7dd00a746d485bf2c02330bd2bda50f98558"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByHash
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0xb89c8992dea2f3f9dc48fd876b6f7dd00a746d485bf2c02330bd2bda50f98558",
    "parent_hash": "0x6e93a5274fa198be164f46760c95066518d10d47e73f3d6861d93b02ba556b29",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
    "root": "0x0a22f243ddb4b9e5251cd1ad966e9170864d7efb7462107683df06e1523ffe79",
    "tx_hash": "0xc999f11db91fc3398fb4f5c2d9cf706d75b4d40ecdc2a84851e5a7768ac7b739",
    "receipt_hash": "0x9f8c36d852fae8cffccfec05c5e84939c6a88e0cf63e7dc3963e5657d9ddc09e",
    "parent_beacon_root": "0x5f702b3489d2683f9d325e85fe2919d7c4d08198e6fa372d6b7f3d42d8834fa4",
    "difficulty": "0",
    "number": "21105407",
    "gas_limit": "30000000",
    "gas_used": "15379455",
    "time": "1730616347",
    "extra": "beaverbuild.org",
    "mix_digest": "0xe53ea84af490d6d4a7169dfa9dbe00ecec1a7b6afb5c53a1035d60bf4d799fd4",
    "nonce": "0",
    "base_fee": "3384783595",
    "withdrawals_hash": "0xbea8080e76aad28e92218ba212b9c8ccb21d26d0ea3a357e6d02866e7a9de06e",
    "blob_gas_used": "393216",
    "excess_blob_gas": "67371008"
  }
}
```

## block by number and hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "height": "21105407"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "block by number success",
  "height": "0",
  "hash": "0xb89c8992dea2f3f9dc48fd876b6f7dd00a746d485bf2c02330bd2bda50f98558",
  "base_fee": "0xc9bfb2eb",
  "transactions": [
    {
      "from": "0x7d14b142cad1379e85682f4b2006cdfed38988d3",
      "to": "0xe592427a0aece92de3edee1f18e0157c05861564",
      "hash": "0xb741eac5553a3be16c5bc14e57b1dd1116c91678bfff508d35d22fd1a5b805eb",
      "amount": "0x0"
    },
    {
      "from": "0xe75ed6f453c602bd696ce27af11565edc9b46b0d",
      "to": "0x00000000009e50a7ddb7a7b0e2ee6604fd120e49",
      "hash": "0xbbc7c1df29106bfde4122bf8c9df1c878830ef2be8a6b7c8d025ebb4b7cb8a4d",
      "amount": "0xf558527"
    }
  ]
}
```

## get account 

- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "network": "mainnet",
  "address": "0x922dB1A931327CA2680343eD2d5E4541669701e9"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get account response success",
  "network": "",
  "account_number": "0",
  "sequence": "0x0",
  "balance": "0"
}
```

## get tx by address
- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "coin": "ETH",
  "network": "mainnet",
  "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb",
  "contractAddress": "0x00",
  "page": 1,
  "pagesize": 10
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get tx list success",
  "tx": [
    {
      "hash": "0x129f03e0840dc2e686c36f1c240ad175184b97150a1c6ddd0ec854923d66f0dc",
      "index": 0,
      "froms": [
        {
          "address": "0x9b0c45d46d386cedd98873168c36efd0dcba8d46"
        }
      ],
      "tos": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "fee": "0x129f03e0840dc2e686c36f1c240ad175184b97150a1c6ddd0ec854923d66f0dc",
      "status": "Success",
      "values": [
        {
          "value": "1200000000000000"
        }
      ],
      "type": 1,
      "height": "19675654",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0xeba43ea75ceb5ee36bda02f9959143b1bd2818778dca15dd11ab046068e6c254",
      "index": 0,
      "froms": [
        {
          "address": "0x9b0c45d46d386cedd98873168c36efd0dcba8d46"
        }
      ],
      "tos": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "fee": "0xeba43ea75ceb5ee36bda02f9959143b1bd2818778dca15dd11ab046068e6c254",
      "status": "Success",
      "values": [
        {
          "value": "368800000000000000"
        }
      ],
      "type": 1,
      "height": "19675742",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0xc9770f085ed67627f60feed89798db4ec810c60e5beb1b4d88560ed860d36af9",
      "index": 0,
      "froms": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "tos": [
        {
          "address": "0x1ecde88a9300451290a2d3d82fd95e615bedfc79"
        }
      ],
      "fee": "0xc9770f085ed67627f60feed89798db4ec810c60e5beb1b4d88560ed860d36af9",
      "status": "Success",
      "values": [
        {
          "value": "5000000000000000"
        }
      ],
      "type": 1,
      "height": "19732309",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0x8aea877e08616abf156d887b30165066fd44a70b0a28c85054812d4109db652e",
      "index": 0,
      "froms": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "tos": [
        {
          "address": "0x1ecde88a9300451290a2d3d82fd95e615bedfc79"
        }
      ],
      "fee": "0x8aea877e08616abf156d887b30165066fd44a70b0a28c85054812d4109db652e",
      "status": "Success",
      "values": [
        {
          "value": "5000000000000000"
        }
      ],
      "type": 1,
      "height": "19733505",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0x394e28037622810807fb9b0f774d3b222dbadf76f7d27016c77107fbc2d847ca",
      "index": 0,
      "froms": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "tos": [
        {
          "address": "0x9b0c45d46d386cedd98873168c36efd0dcba8d46"
        }
      ],
      "fee": "0x394e28037622810807fb9b0f774d3b222dbadf76f7d27016c77107fbc2d847ca",
      "status": "Success",
      "values": [
        {
          "value": "100000000000000000"
        }
      ],
      "type": 1,
      "height": "19811110",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0x4c37e5521f86a82130d83cb51b7839f85deb7be6720d57c660f16512429d7d02",
      "index": 0,
      "froms": [
        {
          "address": "0x9b0c45d46d386cedd98873168c36efd0dcba8d46"
        }
      ],
      "tos": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "fee": "0x4c37e5521f86a82130d83cb51b7839f85deb7be6720d57c660f16512429d7d02",
      "status": "Success",
      "values": [
        {
          "value": "500000000000000000000"
        }
      ],
      "type": 1,
      "height": "19888591",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0x63133204a433066b8b804b6efbaaabcad60e6c3ada96be2cf2f6ff57177f5f83",
      "index": 0,
      "froms": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "tos": [
        {
          "address": "0x8c5359cf717252540a7727cb1352d3ef239bb77d"
        }
      ],
      "fee": "0x63133204a433066b8b804b6efbaaabcad60e6c3ada96be2cf2f6ff57177f5f83",
      "status": "Success",
      "values": [
        {
          "value": "57132370000000000"
        }
      ],
      "type": 1,
      "height": "19911188",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0x218c392a661cabe96d6d5c06f63671c9cf539d8174149c93f808288856fae6e7",
      "index": 0,
      "froms": [
        {
          "address": "0x8c538da2299bab21448a03002195b5b57193b77d"
        }
      ],
      "tos": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "fee": "0x218c392a661cabe96d6d5c06f63671c9cf539d8174149c93f808288856fae6e7",
      "status": "Success",
      "values": [
        {
          "value": "1000000000"
        }
      ],
      "type": 1,
      "height": "19911214",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0x994c7f858f4b82fdc8e00b46f3871ca07d3b4ff398c9455926f6d55a6479a8d7",
      "index": 0,
      "froms": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "tos": [
        {
          "address": "0x7d007a60ebcfbb04fed24fdb16f083acd08f975d"
        }
      ],
      "fee": "0x994c7f858f4b82fdc8e00b46f3871ca07d3b4ff398c9455926f6d55a6479a8d7",
      "status": "Success",
      "values": [
        {
          "value": "6730090000000000"
        }
      ],
      "type": 1,
      "height": "19911562",
      "contract_address": "",
      "datetime": ""
    },
    {
      "hash": "0x3bfd4317c257a334baafe8a15168388cbe2acd46332b3e2047b3ffbdb7e3f859",
      "index": 0,
      "froms": [
        {
          "address": "0x7d0857e35f7fd99e27ff6e1ae1e9584a5412e75d"
        }
      ],
      "tos": [
        {
          "address": "0x2b3fed49557bd88f78b898684f82fbb355305dbb"
        }
      ],
      "fee": "0x3bfd4317c257a334baafe8a15168388cbe2acd46332b3e2047b3ffbdb7e3f859",
      "status": "Success",
      "values": [
        {
          "value": "100000000"
        }
      ],
      "type": 1,
      "height": "19911572",
      "contract_address": "",
      "datetime": ""
    }
  ]
}
```

## tx by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "coin": "ETH",
  "network": "mainnet",
  "hash": "0xd4b95878225617ebfd1d0971519609a9641e93851cb15ce7d7c2140e027bc5a9"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
-response
```
{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "0xd4b95878225617ebfd1d0971519609a9641e93851cb15ce7d7c2140e027bc5a9",
    "index": 2,
    "froms": [
      {
        "address": ""
      }
    ],
    "tos": [
      {
        "address": "0x738e79fBC9010521763944ddF13aAD7f61502221"
      }
    ],
    "fee": "23812160571",
    "status": "Success",
    "values": [
      {
        "value": "0"
      }
    ],
    "type": 0,
    "height": "21105725",
    "contract_address": "0x738e79fBC9010521763944ddF13aAD7f61502221",
    "datetime": "",
    "data": "859B7571A9DE9685D373763B424D42A9B01BB87999E110D90000000000000000000000002503B5400000000000000031B3E65CC123540000D0BE1FDED5D964619B92B3672C08C43305529BE000DAC17F958D2EE523A2206206994597C13D831EC7000BB8712BD4BEB54C6B958267D9DB0259ABDBB0BFF606"
  }
}
```

## create unsign tx

- request
```
grpcurl -plaintext -d '{
  "chain": "Ethereum",
  "network": "mainnet",
  "base64Tx": "ewogICAibm9uY2UiOjYsCiAgICJmcm9tX2FkZHJlc3MiOiIweGU5MDBBMjVhODI1ZjQ0YzM4OWRhODRCMDU4RkMxQjBkMjBjMTg1QWYiLAogICAidG9fYWRkcmVzcyI6IjB4NzJmRmFBMjg5OTkzYmNhRGEyRTAxNjEyOTk1RTVjNzVkRDgxY2RCQyIsCiAgICJnYXMiOjkxMDAwLAogICAidmFsdWUiOiIxMDAwMDAwMDAwMDAwMDAwMDAwMCIsCiAgICJnYXNfcHJpY2UiOiIxOTUwMDAwMDAwMDAiLAogICAiZ2FzX3RpcF9jYXAiOiIzMjc5OTMxNTAzMjgiLAogICAiZ2FzX2ZlZV9jYXAiOiIzMjc5OTMxNTAzMiIsCiAgICJjaGFpbklkIjoiMSIsCiAgICJ0b2tlbkFkZHJlc3MiOiIweDAwIgp9"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.createUnSignTransaction
```

- response
```
{
  "code": "SUCCESS",
  "msg": "create un sign tx success",
  "un_sign_tx": "0x47b17a22a2f89c31cb68505a461a6c47115869c348c5ea6507ba1a36af762018"
}
```
