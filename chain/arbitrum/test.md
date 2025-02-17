## 1.get block by hash

- request
```
grpcurl -plaintext -d '{
  "chain": "Arbitrum",
  "hash": "0x28e353426602ca061bb2cc4549b2693097f7b84903e9e1a234f888b3a9dc69f6"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "GetBlockByNumber success",
  "height": "307055340",
  "hash": "0x28e353426602ca061bb2cc4549b2693097f7b84903e9e1a234f888b3a9dc69f6",
  "base_fee": "0x1e0b298",
  "transactions": [
    {
      "from": "0x00000000000000000000000000000000000a4b05",
      "to": "0x00000000000000000000000000000000000a4b05",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0x62f2ca2277791d6744b1640a9f9078acae713610c7dd34c591fceee1613eb095",
      "height": "307055340",
      "amount": "0x0"
    },
    {
      "from": "0xc7fd16304e3fe336f83c8dd12197a5417199c589",
      "to": "0xa669e7a0d4b3e4fa48af2de86bd4cd7126be4e13",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0x40b2ddc8ad477cecb575b7269c3da0c1d6b02a1eb5f5cbae9037a6b27dcc7309",
      "height": "307055340",
      "amount": "0x38d7ea4c68000"
    },
    {
      "from": "0x34885554b42ccb65d37fccd2ecb4bb3f33e59391",
      "to": "0x0ff73e9ec871b12e1d03a762aafd2e57a8ba5dbc",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0x2efb5c7461c400b91b937b7e99056c1275dbaac0a6b2555d5fc9557b6f45df65",
      "height": "307055340",
      "amount": "0x26a74d6728000"
    }
  ]
}

```


## 2.get block by number

- request
```
grpcurl -plaintext -d '{
  "chain": "Arbitrum",
  "height": "307055340"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "GetBlockByNumber success",
  "height": "307055340",
  "hash": "0x28e353426602ca061bb2cc4549b2693097f7b84903e9e1a234f888b3a9dc69f6",
  "base_fee": "0x1e0b298",
  "transactions": [
    {
      "from": "0x00000000000000000000000000000000000a4b05",
      "to": "0x00000000000000000000000000000000000a4b05",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0x62f2ca2277791d6744b1640a9f9078acae713610c7dd34c591fceee1613eb095",
      "height": "307055340",
      "amount": "0x0"
    },
    {
      "from": "0xc7fd16304e3fe336f83c8dd12197a5417199c589",
      "to": "0xa669e7a0d4b3e4fa48af2de86bd4cd7126be4e13",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0x40b2ddc8ad477cecb575b7269c3da0c1d6b02a1eb5f5cbae9037a6b27dcc7309",
      "height": "307055340",
      "amount": "0x38d7ea4c68000"
    },
    {
      "from": "0x34885554b42ccb65d37fccd2ecb4bb3f33e59391",
      "to": "0x0ff73e9ec871b12e1d03a762aafd2e57a8ba5dbc",
      "token_address": "",
      "contract_wallet": "",
      "hash": "0x2efb5c7461c400b91b937b7e99056c1275dbaac0a6b2555d5fc9557b6f45df65",
      "height": "307055340",
      "amount": "0x26a74d6728000"
    }
  ]
}

```


## 3.get block header by number

- request
```
grpcurl -plaintext -d '{
  "height": "307055340",
  "chain": "Arbitrum"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0x28e353426602ca061bb2cc4549b2693097f7b84903e9e1a234f888b3a9dc69f6",
    "parent_hash": "0x4302c15541f4a02992ff25f2e144d024d7bc797329611e4b23295b687e42b0e9",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0xA4b000000000000000000073657175656e636572",
    "root": "0xef1c371236053a439267054871b1ad8fe0eff4490c201d7573ac0ecd2bfe4d67",
    "tx_hash": "0xfa7347e76dbdea0d6e0c817c924d57f3fc18926afdab8c379ce472ad2f6fa46d",
    "receipt_hash": "0xbe8dee9ecf5c6053dfa9f2aba28cf56808e6fea679381141e6acd7a0b549196c",
    "parent_beacon_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "difficulty": "1",
    "number": "307055340",
    "gas_limit": "1125899906842624",
    "gas_used": "597312",
    "time": "1739808922",
    "extra": "19b65d28ce7bfdbf35522df3edc9d10c7c34fb28f0358d6fc148ce4bbe42217b",
    "mix_digest": "0x00000000000223e300000000014daad800000000000000200000000000000000",
    "nonce": "1874704",
    "base_fee": "31503000",
    "withdrawals_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}

```

## 4.get block header by hash

- request
```
grpcurl -plaintext -d '{
  "chain": "Arbitrum",
  "hash": "0x28e353426602ca061bb2cc4549b2693097f7b84903e9e1a234f888b3a9dc69f6"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0x28e353426602ca061bb2cc4549b2693097f7b84903e9e1a234f888b3a9dc69f6",
    "parent_hash": "0x4302c15541f4a02992ff25f2e144d024d7bc797329611e4b23295b687e42b0e9",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0xA4b000000000000000000073657175656e636572",
    "root": "0xef1c371236053a439267054871b1ad8fe0eff4490c201d7573ac0ecd2bfe4d67",
    "tx_hash": "0xfa7347e76dbdea0d6e0c817c924d57f3fc18926afdab8c379ce472ad2f6fa46d",
    "receipt_hash": "0xbe8dee9ecf5c6053dfa9f2aba28cf56808e6fea679381141e6acd7a0b549196c",
    "parent_beacon_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "difficulty": "1",
    "number": "307055340",
    "gas_limit": "1125899906842624",
    "gas_used": "597312",
    "time": "1739808922",
    "extra": "19b65d28ce7bfdbf35522df3edc9d10c7c34fb28f0358d6fc148ce4bbe42217b",
    "mix_digest": "0x00000000000223e300000000014daad800000000000000200000000000000000",
    "nonce": "1874704",
    "base_fee": "31503000",
    "withdrawals_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}

```
