# Test for ZkSync grpc api

> **Important Note about Proto Method Names**:
> Please note that there is a case sensitivity issue with gRPC method names in the proto files.
> The method names in the proto files should follow consistent casing conventions:
> - Current implementation uses camelCase (e.g., `buildUnSignTransaction`)
> - Some methods might appear with PascalCase in generated code (e.g., `BuildUnSignTransaction`)
> - To avoid issues, always use the exact method name as defined in the proto file when making gRPC calls
> - If you encounter "unknown method" errors, verify the method name casing in your request matches the proto definition

## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
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
  "chain": "Zksync",
  "network": "mainnet",
  "publicKey": "048846b3ce4376e8d58c83c1c6420a784caa675d7f26c496f499585d09891af8fc9167a4b658b57b28211783cdee651caa8b5341b753fa39c995317670123f12d8"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "convert address success",
  "address": "0x82565b64e8063674CAea7003979280f4dbC3aAE7"
}
```

## 3.valid address
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "address": "0x8916B42a4DB16CA71080dBB0f3650162Ad1E7e3e"
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

## 4.get block by number
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "height": "57458640"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
- response
```
{
  "msg": "get latest block header success",
  "height": "57458640",
  "hash": "0x604e1f80266e3306e8b9b093d2ae2c5dba9e45318eca968a8bac8d46a9f53283",
  "baseFee": "0x2b275d0",
  "transactions": [
    {
      "from": "0x2921b419f43898b25d32b74a1f49657f0989a56d",
      "to": "0xebd1e414ebb98522cfd932104ba41fac10a4ef35",
      "hash": "0x9daf9fa2e49c51ec569a62fd7aece8621e7c11bd76b4904a81929fff4963b2ea",
      "height": "57458640",
      "amount": "0x3e871b540c000"
    },
    {
      "from": "0x03ac0b1b952c643d66a4dc1fbc75118109cc074c",
      "to": "0xae45cbe2d1e90358cbd216bc16f2c9267a4ea80a",
      "hash": "0x913cf8cd5d0bdbb4bc5b87279680c49080885d61a658438e58af9b3217ca4722",
      "height": "57458640",
      "amount": "0x0"
    }
  ]
}
```

## 5.get block by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "hash": "0x604e1f80266e3306e8b9b093d2ae2c5dba9e45318eca968a8bac8d46a9f53283"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByHash
```
- response
```
{
  "msg": "get block by hash success",
  "height": "57458640",
  "hash": "0x604e1f80266e3306e8b9b093d2ae2c5dba9e45318eca968a8bac8d46a9f53283",
  "baseFee": "0x2b275d0",
  "transactions": [
    {
      "from": "0x2921b419f43898b25d32b74a1f49657f0989a56d",
      "to": "0xebd1e414ebb98522cfd932104ba41fac10a4ef35",
      "hash": "0x9daf9fa2e49c51ec569a62fd7aece8621e7c11bd76b4904a81929fff4963b2ea",
      "height": "57458640",
      "amount": "0x3e871b540c000"
    },
    {
      "from": "0x03ac0b1b952c643d66a4dc1fbc75118109cc074c",
      "to": "0xae45cbe2d1e90358cbd216bc16f2c9267a4ea80a",
      "hash": "0x913cf8cd5d0bdbb4bc5b87279680c49080885d61a658438e58af9b3217ca4722",
      "height": "57458640",
      "amount": "0x0"
    }
  ]
}
```

## 6.get block header by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "hash": "0x604e1f80266e3306e8b9b093d2ae2c5dba9e45318eca968a8bac8d46a9f53283"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get block header by hash success",
  "block_header": {
    "hash": "0xebe10f45a8b632c235ee1190e1a72390b0a71a11c003374044ea42dacf4fb879",
    "parent_hash": "0xbd28c4f62cd8846a59f3ce4c1ccdd4babf315db42c1859854736f23ba4ef9477",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x0000000000000000000000000000000000000000",
    "root": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "tx_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "receipt_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "difficulty": "0",
    "number": "57458640",
    "gas_limit": "1125899906842624",
    "gas_used": "203147",
    "time": "1741663389",
    "mix_digest": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0",
    "base_fee": "45250000"
  }
}
```

## 7.get block header by number
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "height": "57458640"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get block header by number success",
  "block_header": {
    "hash": "0xebe10f45a8b632c235ee1190e1a72390b0a71a11c003374044ea42dacf4fb879",
    "parent_hash": "0xbd28c4f62cd8846a59f3ce4c1ccdd4babf315db42c1859854736f23ba4ef9477",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x0000000000000000000000000000000000000000",
    "root": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "tx_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "receipt_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "difficulty": "0",
    "number": "57458640",
    "gas_limit": "1125899906842624",
    "gas_used": "203147",
    "time": "1741663389",
    "mix_digest": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0",
    "base_fee": "45250000"
  }
}
```

## 8.get account
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "address": "0x000002c34bAE6DD7BeC72AcbA6aAAC1e01A359De",
  "contractAddress": "0x00"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "msg": "get account response success",
  "accountNumber": "0",
  "sequence": "10918",
  "balance": "152704368597482025"
}
```

## 9.get fee
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```
- response
```
{
  "msg": "get gas price success",
  "slowFee": "45250000|0",
  "normalFee": "45250000|0|*2",
  "fastFee": "45250000|0|*3"
}
```

## 10.get tx by address
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "address": "0x000002c34bAE6DD7BeC72AcbA6aAAC1e01A359De"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```
- response
```
{
  "msg": "get tx list success",
  "tx": [
    {
      "hash": "0x49a01b6857333f486819ca23b3401fd9abf188c665a3cd1a7a2d8eb29973ea9e",
      "from": "0x000002c34bae6dd7bec72acba6aaac1e01a359de",
      "to": "0x0000000000000000000000000000000000008001",
      "fee": "0x49a01b6857333f486819ca23b3401fd9abf188c665a3cd1a7a2d8eb29973ea9e",
      "status": "Success",
      "value": "202964600000000",
      "type": 1,
      "height": "25341436",
      "contractAddress": "0x000000000000000000000000000000000000800a"
    }
    ...
  ]
}
```

## 11.get tx by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "hash": "0xa5d66082c85a722424675105002724f2e8c442281daf1b82ca22136f1a242342"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
- response
```
{
  "msg": "get transaction success",
  "tx": {
    "hash": "0xa5d66082c85a722424675105002724f2e8c442281daf1b82ca22136f1a242342",
    "from": "0x506289e3db9263818bcc999f72293f1497fcde62",
    "to": "0x621425a1ef6abe91058e9712575dcc4258f8d091",
    "fee": "0x2b275d0",
    "status": "Success",
    "value": "0x0",
    "height": "0x254f584",
    "contractAddress": "0x0000000000000000000000000000000000000000"
  }
}
```

## 12.build unsigned transaction
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "base64Tx": "eyJjaGFpbl9pZCI6IjMyNCIsIm5vbmNlIjoxLCJtYXhfcHJpb3JpdHlfZmVlX3Blcl9nYXMiOiIxMDAwMDAwMDAwIiwibWF4X2ZlZV9wZXJfZ2FzIjoiMjAwMDAwMDAwMDAiLCJnYXNfbGltaXQiOjIxMDAwLCJmcm9tX2FkZHJlc3MiOiIweDgyNTY1YjY0ZTgwNjM2NzRDQWVhNzAwMzk3OTI4MGY0ZGJDM2FBRTciLCJ0b19hZGRyZXNzIjoiMHg4OTE2QjQyYTREQjE2Q0E3MTA4MGRCQjBmMzY1MDE2MkFkMUU3ZTNlIiwiYW1vdW50IjoiMTAwMDAwMDAwMDAwMDAwMDAwMCIsImNvbnRyYWN0X2FkZHJlc3MiOiIweDAwIn0="
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.buildUnSignTransaction
```
- response
```
{
  "code": "SUCCESS",
  "msg": "create unsign transaction success",
  "un_sign_tx": "0x6e959617f7fdfff5379834171f28680021219479bb189a51c312a7f584224269"
}
```

## 13.build signed transaction
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "base64Tx": "eyJjaGFpbl9pZCI6IjMyNCIsIm5vbmNlIjoxLCJtYXhfcHJpb3JpdHlfZmVlX3Blcl9nYXMiOiIxMDAwMDAwMDAwIiwibWF4X2ZlZV9wZXJfZ2FzIjoiMjAwMDAwMDAwMDAiLCJnYXNfbGltaXQiOjIxMDAwLCJmcm9tX2FkZHJlc3MiOiIweDgyNTY1YjY0ZTgwNjM2NzRDQWVhNzAwMzk3OTI4MGY0ZGJDM2FBRTciLCJ0b19hZGRyZXNzIjoiMHg4OTE2QjQyYTREQjE2Q0E3MTA4MGRCQjBmMzY1MDE2MkFkMUU3ZTNlIiwiYW1vdW50IjoiMTAwMDAwMDAwMDAwMDAwMDAwMCIsImNvbnRyYWN0X2FkZHJlc3MiOiIweDAwIn0=",
  "signature": "52cf3aa0a66dfe64b6ec18f0bef0e0c90371fc5117c808a024d2c56db5e690f91af6509a8a619438b5babe2c352d2fc20fbf62bffffe72538e0eaa466ad327d601"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.buildSignedTransaction
```
- response
```
{
  "msg": "0x0b713e4c8028f64236ff84eb317abf2c030aad5566850963510ecc61cd5431f0",
  "signedTx": "0x02f87582014401843b9aca008504a817c800825208948916b42a4db16ca71080dbb0f3650162ad1e7e3e880de0b6b3a764000080c001a052cf3aa0a66dfe64b6ec18f0bef0e0c90371fc5117c808a024d2c56db5e690f9a01af6509a8a619438b5babe2c352d2fc20fbf62bffffe72538e0eaa466ad327d6"
}
```

## 14.verify signed transaction
- request
```
grpcurl -plaintext -d '{
  "chain": "Zksync",
  "network": "mainnet",
  "publicKey": "048846b3ce4376e8d58c83c1c6420a784caa675d7f26c496f499585d09891af8fc9167a4b658b57b28211783cdee651caa8b5341b753fa39c995317670123f12d8",
  "signature": "0x6e959617f7fdfff5379834171f28680021219479bb189a51c312a7f584224269:52cf3aa0a66dfe64b6ec18f0bef0e0c90371fc5117c808a024d2c56db5e690f91af6509a8a619438b5babe2c352d2fc20fbf62bffffe72538e0eaa466ad327d601"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.verifySignedTransaction
```
- response
```
{
  "msg": "verify transaction success",
  "verify": true
}
```