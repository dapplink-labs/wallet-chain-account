# Test for grpc api

## 1.support chain
- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getSupportChains
```
- response
```
{
  "msg": "Support this chain",
  "support": true
}
```

## 2.convert address

- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "publicKey": "0xe9ad4b2f85daedb54f9ba61e09d12e3fb92c28913598c350583406ad8651ad8f"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
- reponse

```
{
  "msg": "convert address success",
  "address": "0x3a8eef8a52bc873f5416e835e7ec7da6dd978e5f6a8a12d278df0c42ef01d131"
}

```


## 3.valid address

- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "address": "0x3a8eef8a52bc873f5416e835e7ec7da6dd978e5f6a8a12d278df0c42ef01d131"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.validAddress
```
- response
```
{
  "msg": "ValidAddress success",
  "valid": true
}

```


## latest block header by number

- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "height": "0"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber

./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "height": "248978488"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```
- reponse
```
{
  "msg": "GetBlockHeaderByNumber success",
  "blockHeader": {
    "number": "248978488",
    "time": "1730949333"
  }
}

{
  "msg": "GetBlockHeaderByNumber success",
  "blockHeader": {
    "hash": "0x049336ba3f1ebcdaebd3ca5e0154be616bfae95065a308b1c205daa283ef187e",
    "number": "248978488",
    "time": "1730949333870127"
  }
}

```

## block header by hash

- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "hash": "0x049336ba3f1ebcdaebd3ca5e0154be616bfae95065a308b1c205daa283ef187e"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByHash
```

- response
```
  Code: Internal
  Message: Panic err: implement me

```

## block by number
- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "height": "248978488"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
- response
```
{
  "msg": "GetBlockByNumber success",
  "height": "248978488",
  "hash": "0x049336ba3f1ebcdaebd3ca5e0154be616bfae95065a308b1c205daa283ef187e"
}

```

## get account

- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "address": "0x8d2d7bcde13b2513617df3f98cdd5d0e4b9f714c6308b9204fe18ad900d92609"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "msg": "GetAccount success",
  "network": "mainnet",
  "sequence": "24",
  "balance": "68374979"
}

```

## get tx by address
- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "address": "0xb5e1cc180e603037887c9e9eb4a8a06774ebcddafac37ceea9e33f3b6552bb25"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```

- response
```
{
  "msg": "GetTxByAddress success",
  "tx": [
    {
      "hash": "0x4320eb74382aad6c14f85765e6540705453091eb85175786ad5c49233dafaf92",
      "froms": [
        {
          "address": "0xb5e1cc180e603037887c9e9eb4a8a06774ebcddafac37ceea9e33f3b6552bb25"
        }
      ],
      "fee": "900",
      "status": "Success",
      "height": "1893091016",
      "datetime": "1730949337030508",
      "data": "{\"accumulator_root_hash\":\"0x97013687270eb0e7fdb73d9ab90fc726b5802753c8e0c4c68b7e8ba4b2fe52f4\",\"changes\":[{\"address\":\"0xb5e1cc180e603037887c9e9eb4a8a06774ebcddafac37ceea9e33f3b6552bb25\",\"state_key_hash\":\"0xae5e7cf8c139d4e5c086b1ed6d6464fb53c691f299418cddb64cd45666590067\",\"data\":{\"type\":\"0x1::coin::CoinStore\\u003c0x1::aptos_coin::AptosCoin\\u003e\",\"data\":{\"events\":{\"counter\":\"\",\"guid\":{\"id\":{\"addr\":\"\",\"creation_num\":\"\"}}},\"genesis\":\"\",\"reward_epochs\":null,\"rewards\":null}},\"type\":\"write_resource\"},{\"address\":\"0xb5e1cc180e603037887c9e9eb4a8a06774ebcddafac37ceea9e33f3b6552bb25\",\"state_key_hash\":\"0x91e5df0df2ebd5f57eebc6b724ef19e6feb5d295888195f8fcbd1c70d627dcd7\",\"data\":{\"type\":\"0x1::account::Account\",\"data\":{\"events\":{\"counter\":\"\",\"guid\":{\"id\":{\"addr\":\"\",\"creation_num\":\"\"}}},\"genesis\":\"\",\"reward_epochs\":null,\"rewards\":null}},\"type\":\"write_resource\"},{\"address\":\"0xc0125b866e871048903ef146ab24a577291ab0f37446d3f49100de0ca738a10f\",\"state_key_hash\":\"0xfca859893a759feb7783d0db7305ee3cf001b24714e75c62fad224f7e33e8d6c\",\"data\":{\"type\":\"0x1::coin::CoinStore\\u003c0x1::aptos_coin::AptosCoin\\u003e\",\"data\":{\"events\":{\"counter\":\"\",\"guid\":{\"id\":{\"addr\":\"\",\"creation_num\":\"\"}}},\"genesis\":\"\",\"reward_epochs\":null,\"rewards\":null}},\"type\":\"write_resource\"},{\"address\":\"\",\"state_key_hash\":\"0x6e4b28d40f98a106a65163530924c0dcb40c1349d3aa915d108b4d6cfc1ddb19\",\"data\":{\"type\":\"\",\"data\":{\"events\":{\"counter\":\"\",\"guid\":{\"id\":{\"addr\":\"\",\"creation_num\":\"\"}}},\"genesis\":\"\",\"reward_epochs\":null,\"rewards\":null}},\"type\":\"write_table_item\"}],\"events\":[{\"guid\":{\"creation_number\":\"3\",\"account_address\":\"0xb5e1cc180e603037887c9e9eb4a8a06774ebcddafac37ceea9e33f3b6552bb25\"},\"sequence_number\":\"394648\",\"type\":\"0x1::coin::WithdrawEvent\",\"data\":{\"execution_gas_units\":\"0\",\"io_gas_units\":\"0\",\"storage_fee_octas\":\"0\",\"storage_fee_refund_octas\":\"0\",\"total_charge_gas_units\":\"0\"}},{\"guid\":{\"creation_number\":\"2\",\"account_address\":\"0xc0125b866e871048903ef146ab24a577291ab0f37446d3f49100de0ca738a10f\"},\"sequence_number\":\"23\",\"type\":\"0x1::coin::DepositEvent\",\"data\":{\"execution_gas_units\":\"0\",\"io_gas_units\":\"0\",\"storage_fee_octas\":\"0\",\"storage_fee_refund_octas\":\"0\",\"total_charge_gas_units\":\"0\"}},{\"guid\":{\"creation_number\":\"0\",\"account_address\":\"0x0\"},\"sequence_number\":\"0\",\"type\":\"0x1::transaction_fee::FeeStatement\",\"data\":{\"execution_gas_units\":\"4\",\"io_gas_units\":\"5\",\"storage_fee_octas\":\"0\",\"storage_fee_refund_octas\":\"0\",\"total_charge_gas_units\":\"9\"}}],\"payload\":{\"function\":\"0x1::aptos_account::transfer\",\"type_arguments\":[],\"arguments\":[\"0xc0125b866e871048903ef146ab24a577291ab0f37446d3f49100de0ca738a10f\",\"60000\"],\"type\":\"entry_function_payload\"},\"signature\":{\"public_key\":\"0xf5ba4ddf7f492a8bf946fbfeaa670863bf8bd8d48550284409aa41206d5803e7\",\"signature\":\"0x6d23bc749a6eefa83c2fb0150bb903710d14e6a2a4b1d7c6a7531b8ccd3e2ead6740b585dbad2dfd05157060f7a49911d6e61f543fd51861aa0d32cf0f19b904\",\"type\":\"ed25519_signature\"},\"success\":true,\"vm_status\":\"Executed successfully\"}"
    }
  ]
}
```

## tx by hash
- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "hash": "0x43531969ff8e93de962ea65e5609c2b05de3aa5e78933d8925613e75d3d53772"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```

- response
```
{
  "msg": "GetTxByHash success",
  "tx": {
    "hash": "0x43531969ff8e93de962ea65e5609c2b05de3aa5e78933d8925613e75d3d53772",
    "index": 388643,
    "froms": [
      {
        "address": "0xb5e1cc180e603037887c9e9eb4a8a06774ebcddafac37ceea9e33f3b6552bb25"
      }
    ],
    "tos": [
      {
        "address": "0xfe855b209ffd135f9670693d9924c66341111c2cca587986b0c48320cbf7d28e"
      }
    ],
    "fee": "198700",
    "status": "Success",
    "values": [
      {
        "value": "60000"
      }
    ],
    "type": 1,
    "datetime": "1730629953018074",
    "data": "{\"function\":\"0x1::aptos_account::transfer\",\"type_arguments\":[],\"arguments\":[\"0xfe855b209ffd135f9670693d9924c66341111c2cca587986b0c48320cbf7d28e\",\"60000\"],\"type\":\"entry_function_payload\"}"
  }
}
```

## create unsign tx

- request
```
./grpcurl.exe -plaintext -d '{
  "chain": "Aptos",
  "network": "mainnet",
  "base64Tx": "eyJGcm9tQWRkcmVzcyI6IjB4ZmY5NmFkNTE3ZGIwZjU4NzI0Y2Y1MWI3ODdiNGQ3MTM5NmY2MzRmODczMGZmMmE2ZjBlNWQxYmYzOGRjYjUzYyIsIlB1YmxpY0tleSI6IiIsIlRvQWRkcmVzcyI6IjB4Y2U2OWIwMDA1MTAyYWRjMTUwYjFiMTNiZmM0ZWE5ZjZkYzNmYjkwOWNhYTgzYmQzMzY0ZmMwZjk0ODNlN2NkOSIsIkFtb3VudCI6MTAwMDB9"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.createUnSignTransaction
```

- response
```
{
  "msg": "CreateUnSignTransaction success",
  "unSignTx": "/5atUX2w9Yckz1G3h7TXE5b2NPhzD/Km8OXRvzjctTwCAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ1hcHRvc19hY2NvdW50CHRyYW5zZmVyAAIgzmmwAFECrcFQsbE7/E6p9tw/uQnKqDvTNk/A+Ug+fNkIECcAAAAAAACghgEAAAAAAGQAAAAAAAAAWTUsZwAAAAAB"
}

```
