# Test for grpc api

## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Oasis"
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
  "chain": "Oasis",
  "publicKey": "587c480305e3b3a1b69b95f7f1c2b2390d2e79d7465280808b0915801f269084"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
- reponse

```
{
  "code": "SUCCESS",
  "msg": "",
  "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
}
```

## 3.valid address

- request
```
grpcurl -plaintext -d '{
  "chain": "Oasis",
  "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.validAddress
```
- response
```
{
  "code": "SUCCESS",
  "msg": "",
  "valid": true
}
```


## block by number 
- request
```
grpcurl -plaintext -d '{
  "chain": "Oasis",
  "height": "23359474",
  "viewTx": true
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "",
  "height": "23359474",
  "hash": "string",
  "base_fee": "string",
  "transactions": [
    {
      "from": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr",
      "to": "oasis1qrtfe7s6cuhvyw4r4wpxs30za8dtnn0ytclausys",
      "token_address": "",
      "contract_wallet": "",
      "hash": "ad367bf760255afce848769a9cc601bea7fe16b32a808e24a247b82472123dc3",
      "height": "23359474",
      "amount": "90000000"
    }
  ]
}
```

## get account

- request
```
grpcurl -plaintext -d '{
  "chain": "Oasis",
  "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "code": "SUCCESS",
  "msg": "",
  "network": "",
  "account_number": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr",
  "sequence": "20",
  "balance": "0"
}
```

## get tx by address
- request
```
grpcurl -plaintext -d '{
  "chain": "Oasis",
  "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get tx list success",
  "tx": [
    {
      "hash": "ad367bf760255afce848769a9cc601bea7fe16b32a808e24a247b82472123dc3",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrtfe7s6cuhvyw4r4wpxs30za8dtnn0ytclausys"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "90000000"
        }
      ],
      "type": 1,
      "height": "23359474",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "a2fbbb2185aa97798b5ba8ecf409bd164fa1ec9f93e87c3a3c137ea80aed928b",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrtfe7s6cuhvyw4r4wpxs30za8dtnn0ytclausys"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "428888889"
        }
      ],
      "type": 1,
      "height": "23347735",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "44d3eae560b1f69c14aff4e93223fd4ee88ed53edbd419c5ef5011db5404ba7a",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "50000000"
        }
      ],
      "type": 1,
      "height": "23347572",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "c8e60f073efcb3603575da71d7b8ce4dea2c60ba9d82343f49cc98527c0ded9c",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "30000000"
        }
      ],
      "type": 1,
      "height": "23347557",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "fd1052c1ce82894007ac03a8c65608f279dd7a12daa2543ea0d9430df2172f45",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "20000000"
        }
      ],
      "type": 1,
      "height": "23347525",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "1fba3cf674e70aade79175b06e5f25267e94be71a112eb5c6866ff0a68580f84",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "20000000"
        }
      ],
      "type": 1,
      "height": "23347137",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "d1e75743eeffa495a6d6c073d76333c8a5fff4db466564468adb3f17d507d1d3",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "20000000"
        }
      ],
      "type": 1,
      "height": "23346889",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "d84f5252600267c716f9e3859426a70fbf28b7103b5a2ebdd496132e32a4e10e",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "20000000"
        }
      ],
      "type": 1,
      "height": "23346771",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "3646ec1a61042662be1b49980079e8372d9887ea33a13c1400e31372f1416d7b",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "10000000"
        }
      ],
      "type": 1,
      "height": "23346755",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "7c726f94ba87aa56ba995dd51a648a173b85da4042a110206af9108da3e2911d",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "1000000"
        }
      ],
      "type": 1,
      "height": "23346746",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "669e22a86d8a1ee21f23df3fd093e7c1565146bdde176aa5bfa3d68471292076",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "1000000"
        }
      ],
      "type": 1,
      "height": "23346743",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "7c47fa6de69c8f8095b108a22c778a92597eda4fb85e5132247858f50e4d22cf",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "11111111"
        }
      ],
      "type": 1,
      "height": "23346693",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "6f678d2c111432efadb626cb6c4149d6dcc37183596f0e981193a0703aab2524",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "11111"
        }
      ],
      "type": 1,
      "height": "23346690",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "b42a1b7ef8dd5537b7d20bfd4c60a26537f9194d409a2a5cbf5ac9b2c1f61f17",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "111"
        }
      ],
      "type": 1,
      "height": "23346686",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "6a8fafc29b183681738c00b4c6b809a24be2d979035bae01d0e57ed708d07d1f",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "1"
        }
      ],
      "type": 1,
      "height": "23346683",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "a7fedac2b5908b620c41329fdf3ca80d5b6d09f2c04265dfee61fe8de84a07b4",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "0"
        }
      ],
      "type": 1,
      "height": "23346653",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "87f53ce7748ee652338fbb7ee1ccf646eb77abe37d12f6832ffa3553174f7945",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "0"
        }
      ],
      "type": 1,
      "height": "23346631",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "e47f480d9b8e5f150e2fd3a6485e9506b8f3c98ba8efec1a5fed682ae4346822",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "0"
        }
      ],
      "type": 1,
      "height": "23346620",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "f16cd559696290085c4f3e279edaf7f3108bdc9a7445dedd7f75c5755c382f6a",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "tos": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "0"
        }
      ],
      "type": 1,
      "height": "23346451",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "33d68e89073ec208d351c8cee100c5d4a904a90c5c2f6138d38ebb484fa43be8",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "tos": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "700000000"
        }
      ],
      "type": 1,
      "height": "23346139",
      "contract_address": "",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "a7359e98ec2e6b0f1f19c21e25cc5b522329488d081e3b31403f8db2b0907a83",
      "index": 0,
      "froms": [
        {
          "address": "oasis1qrz33mj9092jc75kezad3h5x9sjupyer4cq07880"
        }
      ],
      "tos": [
        {
          "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
        }
      ],
      "fee": "0",
      "status": "Success",
      "values": [
        {
          "value": "10000000"
        }
      ],
      "type": 1,
      "height": "23340621",
      "contract_address": "",
      "datetime": "",
      "data": ""
    }
  ]
}
```

## tx by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Oasis",
  "hash": "ad367bf760255afce848769a9cc601bea7fe16b32a808e24a247b82472123dc3"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
-response
```
{
  "code": "SUCCESS",
  "msg": "get tx  success",
  "tx": {
    "hash": "ad367bf760255afce848769a9cc601bea7fe16b32a808e24a247b82472123dc3",
    "index": 0,
    "froms": [
      {
        "address": "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr"
      }
    ],
    "tos": [
      {
        "address": "oasis1qrtfe7s6cuhvyw4r4wpxs30za8dtnn0ytclausys"
      }
    ],
    "fee": "0",
    "status": "Success",
    "values": [
      {
        "value": "90000000"
      }
    ],
    "type": 1,
    "height": "23359474",
    "contract_address": "",
    "datetime": "",
    "data": ""
  }
}
```

## sendTx

- request
```
grpcurl -plaintext -d '{
  "chain": "Oasis",
  "rawTx": "omlzaWduYXR1cmWiaXNpZ25hdHVyZVhAQj8E0PzyVF47cqC1FPgqC3rKKAg8ALarR2XtH8Irw7BshvMyP8TdjlgNs7hI/BZ7p2WHDo+I4rhZjr5LM39TCGpwdWJsaWNfa2V5WCBYfEgDBeOzobablffxwrI5DS5510ZSgICLCRWAHyaQhHN1bnRydXN0ZWRfcmF3X3ZhbHVlWGCkY2ZlZaJjZ2FzGQYhZmFtb3VudEEAZGJvZHmiYnRvVQDFGO5FeVUsepbIutjehiwlwJMjrmZhbW91bnREAvrwgGVub25jZRBmbWV0aG9kcHN0YWtpbmcuVHJhbnNmZXI="
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.SendTx
```

- response
```
{
 "code": "SUCCESS",
  "msg": "SendTx success",
  "tx_hash": "ad367bf760255afce848769a9cc601bea7fe16b32a808e24a247b82472123dc3"
}
```
