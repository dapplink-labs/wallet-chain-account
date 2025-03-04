# Test for grpc api

## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Near"
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

## 2.latest block header by number if do not send height value. 

- request
```
grpcurl -plaintext -d '{
  "chain": "Near"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```
- reponse
```
{
  "code": "SUCCESS",
  "msg": "get block header by number success",
  "block_header": {
    "hash": "",
    "parent_hash": "AP9K4rqU2KZM4QvEexvqicMXsw7T6XZv4zytGgyhxYiu",
    "uncle_hash": "",
    "coin_base": "",
    "root": "2GhP2qigBwahK2KE21SNJCyhfJgUXfhdZ2wTiiCvVxBQ",
    "tx_hash": "FUMwq2A14uEteiNvL2NScin9zY16sbeohj1FRSVA8rAv",
    "receipt_hash": "76GxnNHvWzjL64zYYG4EL4VrkNqkknctNUvRsAMqm4Bp",
    "parent_beacon_root": "",
    "difficulty": "",
    "number": "140764801",
    "gas_limit": "0",
    "gas_used": "0",
    "time": "1740718297891677053",
    "extra": "",
    "mix_digest": "",
    "nonce": "",
    "base_fee": "",
    "withdrawals_hash": "ETwvgeWMiETBQFqqfsLWjbUdnPD78x1MQq3hpqY6xmbZ",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}
```

## 3.block header by hash

- request
```
grpcurl -plaintext -d '{
  "chain": "Near",
  "hash": "AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByHash
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get block header by hash success",
  "block_header": {
    "hash": "",
    "parent_hash": "8dFBjzfoBCfaZaPPueR6CQ1svEFeNqiUf8opdba93HUg",
    "uncle_hash": "",
    "coin_base": "",
    "root": "9HhqhMb6yEg8hipANmoUJKgAA6joTFtao4KM8zEKdD2M",
    "tx_hash": "6qUV4vaQehJNxwrod42xYh2n2QSnyrAiRYJ6HsuupRWr",
    "receipt_hash": "6gTs3rjA7uLUF3CunjRSPSSnBN41qxANpcWtpuAqU5h4",
    "parent_beacon_root": "",
    "difficulty": "",
    "number": "140083597",
    "gas_limit": "0",
    "gas_used": "0",
    "time": "1739945502159047772",
    "extra": "",
    "mix_digest": "",
    "nonce": "",
    "base_fee": "",
    "withdrawals_hash": "DZ7kyYeWRusxCWxRZXhyQXjKV21WeKK11yMZBropzYEC",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}
```

## 4.block by number 
- request
```
grpcurl -plaintext -d '{
  "chain": "Near",
  "height": "140083597"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
- response
```
{
  "code": "SUCCESS",
  "msg": "block by hash success",
  "height": "140083597",
  "hash": "AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY",
  "base_fee": "100000000",
  "transactions": [
    {
      "from": "30d6205b89e02d2c34ade165ae9705f8d9e80d0e304061cb10b8e56c8e703f62",
      "to": "20b9bdf32f768ac6e6ff3c9ab512d4bd7f94dbcf4e9d15bb8cd3c3b4062d585a",
      "token_address": "30d6205b89e02d2c34ade165ae9705f8d9e80d0e304061cb10b8e56c8e703f62",
      "contract_wallet": "30d6205b89e02d2c34ade165ae9705f8d9e80d0e304061cb10b8e56c8e703f62",
      "hash": "AxLS6E5EjsA6cNBXT5S4qHppCuTnWZU4qZ5n13Qs1knj",
      "height": "140083597",
      "amount": "623820480000000000000000"
    },
    {
      "from": "mytestapi.embr.playember_reserve.near",
      "to": "8e4955f1c0363b72395b578fdbf2dcb01f48cbc7bb407f3df8bca13e57829455",
      "token_address": "mytestapi.embr.playember_reserve.near",
      "contract_wallet": "mytestapi.embr.playember_reserve.near",
      "hash": "2zTry4rLcB8omeq2jGUdjYhh3tTJgF9rpsMctVorn8na",
      "height": "140083597",
      "amount": "1830000000000000000000"
    },
    {
      "from": "mytestapi.embr.playember_reserve.near",
      "to": "c2701f7b5509f05f78861d522c1ee4add612773299df0d51ea96dbd8a1ddbb3f",
      "token_address": "mytestapi.embr.playember_reserve.near",
      "contract_wallet": "mytestapi.embr.playember_reserve.near",
      "hash": "3R7A6xr5DPpzCRbwZFrfdCPu1Dy4EJ53ip5gXUf2cZ9g",
      "height": "140083597",
      "amount": "1830000000000000000000"
    },
    {
      "from": "zzkhyobi6m48.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzkhyobi6m48.users.kaiching",
      "contract_wallet": "zzkhyobi6m48.users.kaiching",
      "hash": "H2tbppTBDQAvMq1fctLPLtp1eqvf3JKqgYmywrF28rt3",
      "height": "140083597",
      "amount": "5062619225811300000000"
    },
    {
      "from": "users.kaiching",
      "to": "efhk4z7of76b.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "7MsZHAugXm5se3396JzbvhUA15MkKiBDFVCmTUw4iiPA",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzlmesya7deo.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzlmesya7deo.users.kaiching",
      "contract_wallet": "zzlmesya7deo.users.kaiching",
      "hash": "D2ufnhhmQ9oDwHxAeZDYXwpdZvVSc77o7YnPNia1Ti6x",
      "height": "140083597",
      "amount": "8204478239888400000000"
    },
    {
      "from": "users.kaiching",
      "to": "u1iqa2n667oy.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "GHUC2Sj2fbT4PTqJ81yAVMsBvThaBFkrSAdvidGxEQon",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzl5nyi5lhnx.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzl5nyi5lhnx.users.kaiching",
      "contract_wallet": "zzl5nyi5lhnx.users.kaiching",
      "hash": "4a5HTe2db6sZv6bg3MV9muSjSSVP2r7zDZS2QFNBf4UY",
      "height": "140083597",
      "amount": "5015125215024100000000"
    },
    {
      "from": "zzl5iaj6x4pc.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzl5iaj6x4pc.users.kaiching",
      "contract_wallet": "zzl5iaj6x4pc.users.kaiching",
      "hash": "AWq78hQ8wHRg7tQzpsjtmUxWNnotFcMSwT5CRPji7TxB",
      "height": "140083597",
      "amount": "9262620087134500000000"
    },
    {
      "from": "users.kaiching",
      "to": "sdc23wvkf26g.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "6kCFEi8eLnVCMjDgBQ3dnf2jWT4XDRXqix1A1m4LQU94",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "t4eiv8f0ei02.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "FyeYNtV2912F6xohiYxVAdmddd5BzHqHKE83MXSwyT6S",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzk3zqpg3l5b.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzk3zqpg3l5b.users.kaiching",
      "contract_wallet": "zzk3zqpg3l5b.users.kaiching",
      "hash": "HLEkPfUDZxDxSbPYBesUaifSUZx1io2oqHb83t4HZbQ",
      "height": "140083597",
      "amount": "9252184734919900000000"
    },
    {
      "from": "users.kaiching",
      "to": "uq3rl1d5yi1f.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "8L8kCKm6FDvQ9L7XBjXnJUUGjPBAfnzk2qSYZkA4JBoG",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "j7nl4h6g8e74.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "HhCszk6U4K9yviVCmARwMaq6qy1vJiMkmLskDL7hFFLQ",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "dfqwmizdem8m.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "Eucs7oQ5nxE8SvtjmJ72PThARhUu59uqeNUR8GkDrvWc",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "2v5bhv0hvtuy.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "25Ub4NGLiGMpHzkqUZhAJdHHaNkSyVyvGxHikm5maz6K",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzlmrwb4fu6w.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzlmrwb4fu6w.users.kaiching",
      "contract_wallet": "zzlmrwb4fu6w.users.kaiching",
      "hash": "ECy5TM8BaxDDtv2P3AAWaAiVFCoRzYZTi2djgay5uUUM",
      "height": "140083597",
      "amount": "9149602412586100000000"
    },
    {
      "from": "users.kaiching",
      "to": "5xcrg7w45g43.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "Hp16shx4xYKyBqeAsvE3dgdYYBDcjVZW9y5h47oCGWC4",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "7w6marju0fhk.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "BdGbarYSLPYoS8AEJkUmjov45erVqD5qWrwmrkBAP918",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "43i7idpodh42.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "G4BXFYKBfiaPanVmWufKivzo9mbGSFWEiW46TFD2t64D",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "je2fsfzwwy54.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "39mWYjRcHH9nFeKaBMtXNUgwC49i941wcMua4Pmm86HP",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "zwvxo22yk6tb.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "6uW6D5rmRcK9gDyeJ8yvv4REqqncDaaRKYU7QvoCnGCo",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzk4tu5mhke7.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzk4tu5mhke7.users.kaiching",
      "contract_wallet": "zzk4tu5mhke7.users.kaiching",
      "hash": "Aake48ev877zLGauYCzdZtY2RzguETcYpEdVjPtUSG1L",
      "height": "140083597",
      "amount": "8666023862372700000000"
    },
    {
      "from": "users.kaiching",
      "to": "nsoceswcdjhl.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "7vnrUa4Gt9H1L26Mt9H65fzdoWMQxi9tz4q6KmjqrAcr",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "a5rtxmmxqyg5.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "8PCmrXwycLHxNRKr15GKCXVjQgsLPLCPAPaLgXV8e9gj",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "95y4y8gilkw7.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "DcULY3XXR6BJua4KgjRkNu6YeJ7eyJxgkHUef3QDCpje",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "of1n7qwzm4ub.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "2beQPZLGfSixbpwpk22LXDhy2HgFxFRj5kQxuCUd6piD",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "1x94uxdutqv5.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "5KCCbbwu7xa9MLSfRSARiBp1LrmywvoWPx45mAPtTLtf",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "1u2makeh0al5.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "HGwrrRgYGfuRSU7zKyKPiArxnmaMEC4c2kumRTk2nB7J",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "mw1uj3us3bh7.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "7ff3jB4jhFN3Pot4cX4uRiowtwQnqunPV3aRJhZNvwcN",
      "height": "140083597",
      "amount": "0"
    }
  ]
}
```

## 4.block by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Near",
  "hash": "AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByHash
```
- response
```
{
  "code": "SUCCESS",
  "msg": "block by hash success",
  "height": "140083597",
  "hash": "AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY",
  "base_fee": "100000000",
  "transactions": [
    {
      "from": "30d6205b89e02d2c34ade165ae9705f8d9e80d0e304061cb10b8e56c8e703f62",
      "to": "20b9bdf32f768ac6e6ff3c9ab512d4bd7f94dbcf4e9d15bb8cd3c3b4062d585a",
      "token_address": "30d6205b89e02d2c34ade165ae9705f8d9e80d0e304061cb10b8e56c8e703f62",
      "contract_wallet": "30d6205b89e02d2c34ade165ae9705f8d9e80d0e304061cb10b8e56c8e703f62",
      "hash": "AxLS6E5EjsA6cNBXT5S4qHppCuTnWZU4qZ5n13Qs1knj",
      "height": "140083597",
      "amount": "623820480000000000000000"
    },
    {
      "from": "zzkhyobi6m48.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzkhyobi6m48.users.kaiching",
      "contract_wallet": "zzkhyobi6m48.users.kaiching",
      "hash": "H2tbppTBDQAvMq1fctLPLtp1eqvf3JKqgYmywrF28rt3",
      "height": "140083597",
      "amount": "5062619225811300000000"
    },
    {
      "from": "users.kaiching",
      "to": "efhk4z7of76b.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "7MsZHAugXm5se3396JzbvhUA15MkKiBDFVCmTUw4iiPA",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzlmesya7deo.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzlmesya7deo.users.kaiching",
      "contract_wallet": "zzlmesya7deo.users.kaiching",
      "hash": "D2ufnhhmQ9oDwHxAeZDYXwpdZvVSc77o7YnPNia1Ti6x",
      "height": "140083597",
      "amount": "8204478239888400000000"
    },
    {
      "from": "users.kaiching",
      "to": "u1iqa2n667oy.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "GHUC2Sj2fbT4PTqJ81yAVMsBvThaBFkrSAdvidGxEQon",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzl5nyi5lhnx.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzl5nyi5lhnx.users.kaiching",
      "contract_wallet": "zzl5nyi5lhnx.users.kaiching",
      "hash": "4a5HTe2db6sZv6bg3MV9muSjSSVP2r7zDZS2QFNBf4UY",
      "height": "140083597",
      "amount": "5015125215024100000000"
    },
    {
      "from": "zzl5iaj6x4pc.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzl5iaj6x4pc.users.kaiching",
      "contract_wallet": "zzl5iaj6x4pc.users.kaiching",
      "hash": "AWq78hQ8wHRg7tQzpsjtmUxWNnotFcMSwT5CRPji7TxB",
      "height": "140083597",
      "amount": "9262620087134500000000"
    },
    {
      "from": "users.kaiching",
      "to": "sdc23wvkf26g.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "6kCFEi8eLnVCMjDgBQ3dnf2jWT4XDRXqix1A1m4LQU94",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "t4eiv8f0ei02.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "FyeYNtV2912F6xohiYxVAdmddd5BzHqHKE83MXSwyT6S",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzk3zqpg3l5b.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzk3zqpg3l5b.users.kaiching",
      "contract_wallet": "zzk3zqpg3l5b.users.kaiching",
      "hash": "HLEkPfUDZxDxSbPYBesUaifSUZx1io2oqHb83t4HZbQ",
      "height": "140083597",
      "amount": "9252184734919900000000"
    },
    {
      "from": "users.kaiching",
      "to": "uq3rl1d5yi1f.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "8L8kCKm6FDvQ9L7XBjXnJUUGjPBAfnzk2qSYZkA4JBoG",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "j7nl4h6g8e74.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "HhCszk6U4K9yviVCmARwMaq6qy1vJiMkmLskDL7hFFLQ",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "dfqwmizdem8m.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "Eucs7oQ5nxE8SvtjmJ72PThARhUu59uqeNUR8GkDrvWc",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "2v5bhv0hvtuy.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "25Ub4NGLiGMpHzkqUZhAJdHHaNkSyVyvGxHikm5maz6K",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzlmrwb4fu6w.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzlmrwb4fu6w.users.kaiching",
      "contract_wallet": "zzlmrwb4fu6w.users.kaiching",
      "hash": "ECy5TM8BaxDDtv2P3AAWaAiVFCoRzYZTi2djgay5uUUM",
      "height": "140083597",
      "amount": "9149602412586100000000"
    },
    {
      "from": "users.kaiching",
      "to": "5xcrg7w45g43.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "Hp16shx4xYKyBqeAsvE3dgdYYBDcjVZW9y5h47oCGWC4",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "7w6marju0fhk.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "BdGbarYSLPYoS8AEJkUmjov45erVqD5qWrwmrkBAP918",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "43i7idpodh42.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "G4BXFYKBfiaPanVmWufKivzo9mbGSFWEiW46TFD2t64D",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "je2fsfzwwy54.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "39mWYjRcHH9nFeKaBMtXNUgwC49i941wcMua4Pmm86HP",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "zwvxo22yk6tb.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "6uW6D5rmRcK9gDyeJ8yvv4REqqncDaaRKYU7QvoCnGCo",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "zzk4tu5mhke7.users.kaiching",
      "to": "reserve.kaiching",
      "token_address": "zzk4tu5mhke7.users.kaiching",
      "contract_wallet": "zzk4tu5mhke7.users.kaiching",
      "hash": "Aake48ev877zLGauYCzdZtY2RzguETcYpEdVjPtUSG1L",
      "height": "140083597",
      "amount": "8666023862372700000000"
    },
    {
      "from": "users.kaiching",
      "to": "nsoceswcdjhl.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "7vnrUa4Gt9H1L26Mt9H65fzdoWMQxi9tz4q6KmjqrAcr",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "a5rtxmmxqyg5.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "8PCmrXwycLHxNRKr15GKCXVjQgsLPLCPAPaLgXV8e9gj",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "95y4y8gilkw7.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "DcULY3XXR6BJua4KgjRkNu6YeJ7eyJxgkHUef3QDCpje",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "of1n7qwzm4ub.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "2beQPZLGfSixbpwpk22LXDhy2HgFxFRj5kQxuCUd6piD",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "1x94uxdutqv5.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "5KCCbbwu7xa9MLSfRSARiBp1LrmywvoWPx45mAPtTLtf",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "1u2makeh0al5.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "HGwrrRgYGfuRSU7zKyKPiArxnmaMEC4c2kumRTk2nB7J",
      "height": "140083597",
      "amount": "0"
    },
    {
      "from": "users.kaiching",
      "to": "mw1uj3us3bh7.users.kaiching",
      "token_address": "users.kaiching",
      "contract_wallet": "users.kaiching",
      "hash": "7ff3jB4jhFN3Pot4cX4uRiowtwQnqunPV3aRJhZNvwcN",
      "height": "140083597",
      "amount": "0"
    }
  ]
}
```

## 5.get account 

- request
```
grpcurl -plaintext -d '{
  "chain": "Near",
  "network": "mainnet",
  "address": "20b9bdf32f768ac6e6ff3c9ab512d4bd7f94dbcf4e9d15bb8cd3c3b4062d585a",
  "contractAddress": "0x00"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
- response
```
{
  "code": "SUCCESS",
  "msg": "get account success",
  "network": "Near",
  "account_number": "",
  "sequence": "",
  "balance": "612590893150272500000000"
}
```

## 6.tx by hash
- request
```
grpcurl -plaintext -d '{
  "chain": "Near",
  "coin": "users.kaiching",
  "hash": "7ff3jB4jhFN3Pot4cX4uRiowtwQnqunPV3aRJhZNvwcN"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```
-response
```
{
  "code": "SUCCESS",
  "msg": "get transaction success",
  "tx": {
    "hash": "7ff3jB4jhFN3Pot4cX4uRiowtwQnqunPV3aRJhZNvwcN",
    "index": 0,
    "froms": [
      {
        "address": "users.kaiching"
      }
    ],
    "tos": [
      {
        "address": "mw1uj3us3bh7.users.kaiching"
      }
    ],
    "fee": "103000000",
    "status": "NotFound",
    "values": [],
    "type": 0,
    "height": "",
    "contract_address": "",
    "datetime": "",
    "data": "{\"final_execution_status\":\"FINAL\",\"receipts\":[{\"predecessor_id\":\"users.kaiching\",\"priority\":0,\"receipt\":{\"Action\":{\"actions\":[\"CreateAccount\",{\"AddKey\":{\"access_key\":{\"nonce\":0,\"permission\":\"FullAccess\"},\"public_key\":\"ed25519:Ei4qADmXscKQtudU98xKEpCmkgU57gC6mtTZCEewbKgE\"}},{\"Transfer\":{\"deposit\":\"0\"}}],\"gas_price\":\"103000000\",\"input_data_ids\":[],\"is_promise_yield\":false,\"output_data_receivers\":[],\"signer_id\":\"users.kaiching\",\"signer_public_key\":\"ed25519:8o7aytbrAyiqyJsb2KCDHBQNpmxSRSsFyws447A8Lz97\"}},\"receipt_id\":\"EbYVc29Tv9ogYZKhfz3qJpCNPT4LVPyAevhc59Cmt9yN\",\"receiver_id\":\"mw1uj3us3bh7.users.kaiching\"},{\"predecessor_id\":\"system\",\"priority\":0,\"receipt\":{\"Action\":{\"actions\":[{\"Transfer\":{\"deposit\":\"12524843062500000000\"}}],\"gas_price\":\"0\",\"input_data_ids\":[],\"is_promise_yield\":false,\"output_data_receivers\":[],\"signer_id\":\"users.kaiching\",\"signer_public_key\":\"ed25519:8o7aytbrAyiqyJsb2KCDHBQNpmxSRSsFyws447A8Lz97\"}},\"receipt_id\":\"6rDDUxZia28qXVezTX9h5Dss4zPdm5nuHDRiTzTqJB3H\",\"receiver_id\":\"users.kaiching\"}],\"receipts_outcome\":[{\"block_hash\":\"A96zR92ZSfgmULZopAt7v2ESUoFnuTJbJ3Pf5QQ486L3\",\"id\":\"EbYVc29Tv9ogYZKhfz3qJpCNPT4LVPyAevhc59Cmt9yN\",\"outcome\":{\"executor_id\":\"mw1uj3us3bh7.users.kaiching\",\"gas_burnt\":4174947687500,\"logs\":[],\"metadata\":{\"gas_profile\":[],\"version\":3},\"receipt_ids\":[\"6rDDUxZia28qXVezTX9h5Dss4zPdm5nuHDRiTzTqJB3H\"],\"status\":{\"SuccessValue\":\"\"},\"tokens_burnt\":\"417494768750000000000\"},\"proof\":[{\"direction\":\"Right\",\"hash\":\"Hx5Z4oSfx4YkxQzMoou2R3jhVHVdJzEiDSU5ECKTNBzy\"},{\"direction\":\"Right\",\"hash\":\"44mLEBKWGKuDtDgpQ1Fbfnu9nuShCgjbX4xNrHvkEK7F\"},{\"direction\":\"Right\",\"hash\":\"6XCfb8cjFYcg7zuvRa7C6rze6fgF599NdGKksKEWCZgE\"},{\"direction\":\"Right\",\"hash\":\"EGMnu3vNMcEgZ45xkJ6uQ36kb1bV7hZXKCnfTXSUCrL3\"},{\"direction\":\"Left\",\"hash\":\"huY1ZJ4voA4M26KrMXQ5RGz9pnE2DHvTd3CWz6WcVMH\"},{\"direction\":\"Left\",\"hash\":\"Fmuix26YDXVAp63HKYdrAPH2zs1U5i8P6hRCa1LPhJSU\"},{\"direction\":\"Left\",\"hash\":\"GYbjs4aJbQcBSBV5sYNCQaBYGEm5fKwrqTzpyEUxMMws\"},{\"direction\":\"Right\",\"hash\":\"4NdMVGvuwhbSjqRHhLh5RQJwijsK1iWJhS3tBAffJHdo\"}]},{\"block_hash\":\"5sdarQpufHBpg7747SFbZtiaXrkJVPjvX3a7agr56g8t\",\"id\":\"6rDDUxZia28qXVezTX9h5Dss4zPdm5nuHDRiTzTqJB3H\",\"outcome\":{\"executor_id\":\"users.kaiching\",\"gas_burnt\":223182562500,\"logs\":[],\"metadata\":{\"gas_profile\":[],\"version\":3},\"receipt_ids\":[],\"status\":{\"SuccessValue\":\"\"},\"tokens_burnt\":\"0\"},\"proof\":[{\"direction\":\"Right\",\"hash\":\"5QJhTnFVMMudEzzjuVaSZnnWyrEKaRsw7nrgoAzd6yf\"},{\"direction\":\"Left\",\"hash\":\"6a4f2YR7d92PgYusmYk6mjamt3YcCskE6v28qEdfEkzH\"},{\"direction\":\"Right\",\"hash\":\"4MdUFejvFgJrFHbDAYeuyBVtr55fkoLjkTDmfAgHbrua\"},{\"direction\":\"Right\",\"hash\":\"5orFz3Rvr85bwPGLDKgZgf1Mn4CHaUU43BsxdNr85RR4\"},{\"direction\":\"Left\",\"hash\":\"DYujoUW99GxwZiZ1x1f5epGzydnK7jz9gfnJTVG5WJZH\"},{\"direction\":\"Right\",\"hash\":\"4Q4VCaTWKMwdL7QFvqnag2Dbz4urdrWVcTW2XZ67Z5jc\"},{\"direction\":\"Right\",\"hash\":\"um7DadjLYgjmTNy8fwVp2ankfb9smcdpQtfGpRSKFiU\"}]}],\"status\":{\"SuccessValue\":\"\"},\"transaction\":{\"actions\":[\"CreateAccount\",{\"AddKey\":{\"access_key\":{\"nonce\":0,\"permission\":\"FullAccess\"},\"public_key\":\"ed25519:Ei4qADmXscKQtudU98xKEpCmkgU57gC6mtTZCEewbKgE\"}},{\"Transfer\":{\"deposit\":\"0\"}}],\"hash\":\"7ff3jB4jhFN3Pot4cX4uRiowtwQnqunPV3aRJhZNvwcN\",\"nonce\":99016713200550,\"priority_fee\":0,\"public_key\":\"ed25519:8o7aytbrAyiqyJsb2KCDHBQNpmxSRSsFyws447A8Lz97\",\"receiver_id\":\"mw1uj3us3bh7.users.kaiching\",\"signature\":\"ed25519:2qwPxJCYfBBkAp3SUzZkLuGzXjjQjjfLkPfPrvieY9awQmcPXw9jBTjzaKR8F4KnPPCMFMpC2dYVmD3rkBzWvPFD\",\"signer_id\":\"users.kaiching\"},\"transaction_outcome\":{\"block_hash\":\"AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY\",\"id\":\"7ff3jB4jhFN3Pot4cX4uRiowtwQnqunPV3aRJhZNvwcN\",\"outcome\":{\"executor_id\":\"users.kaiching\",\"gas_burnt\":4174947687500,\"logs\":[],\"metadata\":{\"gas_profile\":null,\"version\":1},\"receipt_ids\":[\"EbYVc29Tv9ogYZKhfz3qJpCNPT4LVPyAevhc59Cmt9yN\"],\"status\":{\"SuccessReceiptId\":\"EbYVc29Tv9ogYZKhfz3qJpCNPT4LVPyAevhc59Cmt9yN\"},\"tokens_burnt\":\"417494768750000000000\"},\"proof\":[{\"direction\":\"Right\",\"hash\":\"7Qy4WZ7XjYzdsUobuzEL1zXJNmxuVUuNnVCvL3gV1SDM\"},{\"direction\":\"Left\",\"hash\":\"CgGAMWXmTLFGw2nezThLNq4oWtxjzzJaZQYK3sB33BTm\"},{\"direction\":\"Left\",\"hash\":\"3RQ1WdQNaYKgCSxgup3VBjK1KZmVwdY64kPNmS9moz8o\"},{\"direction\":\"Left\",\"hash\":\"DDPfvq77txUm9LLGbwdTQ9BJKodBhkB4tdGDJ8TFnnWw\"},{\"direction\":\"Left\",\"hash\":\"SstQaqw83ktMkjh95jQLjewyJGoZLA5UG8tTXDD4ZJF\"},{\"direction\":\"Right\",\"hash\":\"6SiQXBuKHHZkcJyfd8oWjusbedgshHpSbRCShYZxejoc\"},{\"direction\":\"Right\",\"hash\":\"7QgeCmTjhbXrGaX2gXPoLPqcjYQaESLpnN5jkDLFiwGw\"}]}}"
  }
}
```

## 7.create send tx

- request
```
grpcurl -plaintext -d '{
  "chain": "Near",
  "rawTx": "QAAAADIwYjliZGYzMmY3NjhhYzZlNmZmM2M5YWI1MTJkNGJkN2Y5NGRiY2Y0ZTlkMTViYjhjZDNjM2I0MDYyZDU4NWEAILm98y92isbm/zyatRLUvX+U289OnRW7jNPDtAYtWFpFvQnBZ38AABAAAAByZWNlaXZlci50ZXN0bmV0/qNZXDvLObzPkHVht/hI4QThWo5Cld9SKwgqIxSpsGkBAAAAAwAAoN7Frck1NgAAAAAAAAAAw0MuGjNuAJZ+iWiXzVVCycYeQ8Vzb4BlYQiD7P55z5UywscYC53d2wTRZdslocF0TTDvAjVfBSXaowleDQNtCQ=="
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.SendTx
```

- response
```
{
  "code": "SUCCESS",
  "msg": "send tx success",
  "tx_hash": "AhfNVjhMiRzvYSpqjeJhfzKNfKhf91ZGqh1NSqQ3v7qx"
}
```

## 8.get block by range

- request
```
grpcurl -plaintext -d '{
  "chain": "Near",
  "start": "140083597",
  "end": "140083599"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByRange
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get block range success",
  "block_header": [
    {
      "hash": "",
      "parent_hash": "AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY",
      "uncle_hash": "",
      "coin_base": "",
      "root": "2g6hW9fnPDK8oRPVBgNWZEpo2Q4NeT6TacR1n5oZ2Fv5",
      "tx_hash": "CLvmjDaqz2BgqvQLQ3k6VTua7kCVibkgWrkURx24ft1Q",
      "receipt_hash": "3q5ovFj1uiHGd2G6ydxF4knz468CZryaA5NW9sV95w3f",
      "parent_beacon_root": "",
      "difficulty": "",
      "number": "140083598",
      "gas_limit": "0",
      "gas_used": "0",
      "time": "1739945503226664662",
      "extra": "",
      "mix_digest": "",
      "nonce": "",
      "base_fee": "",
      "withdrawals_hash": "BNTg7m7hkwurwWjNFmKFHoZDiVpr4cHaJXhc8PPTY7EH",
      "blob_gas_used": "0",
      "excess_blob_gas": "0"
    },
    {
      "hash": "",
      "parent_hash": "8dFBjzfoBCfaZaPPueR6CQ1svEFeNqiUf8opdba93HUg",
      "uncle_hash": "",
      "coin_base": "",
      "root": "9HhqhMb6yEg8hipANmoUJKgAA6joTFtao4KM8zEKdD2M",
      "tx_hash": "6qUV4vaQehJNxwrod42xYh2n2QSnyrAiRYJ6HsuupRWr",
      "receipt_hash": "6gTs3rjA7uLUF3CunjRSPSSnBN41qxANpcWtpuAqU5h4",
      "parent_beacon_root": "",
      "difficulty": "",
      "number": "140083597",
      "gas_limit": "0",
      "gas_used": "0",
      "time": "1739945502159047772",
      "extra": "",
      "mix_digest": "",
      "nonce": "",
      "base_fee": "",
      "withdrawals_hash": "DZ7kyYeWRusxCWxRZXhyQXjKV21WeKK11yMZBropzYEC",
      "blob_gas_used": "0",
      "excess_blob_gas": "0"
    }
  ]
}
```

## 9.get fee

- request
```
grpcurl -plaintext -d '{
  "chain": "Near"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```

- response
```
{
  "code": "SUCCESS",
  "msg": "get fee success",
  "slow_fee": "",
  "normal_fee": "100000000",
  "fast_fee": ""
}
```
