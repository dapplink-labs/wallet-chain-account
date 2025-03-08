## 1.create unSign transaction

- request
```
grpcurl -plaintext -d '{
  "base64Tx": "eyJjaGFpbl9pZCI6IjEwIiwibm9uY2UiOjUsImZyb21fYWRkcmVzcyI6IjB4NDc0MGQ3ZUUxYkQ0NTc2YUQ5NjJmMjgwNmIxMTI5OThDYzNCNzJGYyIsInRvX2FkZHJlc3MiOiIweDgyMThhMEY0N0Y0YzBkRTBjMTc1NGY1MDg3NDcwN2NkNmU3YjJlNWUiLCJnYXNfbGltaXQiOjIxMDAwLCJtYXhfZmVlX3Blcl9nYXMiOiIyNjAwMDAwMDAwMCIsIm1heF9wcmlvcml0eV9mZWVfcGVyX2dhcyI6IjIwNTIwMDAwMDAwIiwiYW1vdW50IjoiOTAwMDAwMDAwMDAwMDAwMDAwMCIsImNvbnRyYWN0X2FkZHJlc3MiOiIweGIxMmMxM2U2NmFkZTFmNzJmNzE4MzRmMmZjNTA4MmRiOGMwOTEzNTgifQ==",
  "chain": "Optimism"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.BuildUnSignTransaction
```
- response
```
{
  "code": "SUCCESS",
  "msg": "create unsign transaction success",
  "un_sign_tx": "0xc1c234195ac9871215cd960190893c4b361699b207d4d546ad8d9de175633e08"
}

```


## 2.build signed transaction

- request
```
grpcurl -plaintext -d '{
  "chain": "Optimism",
  "base64Tx": "eyJjaGFpbl9pZCI6IjEwIiwibm9uY2UiOjUsImZyb21fYWRkcmVzcyI6IjB4NDc0MGQ3ZUUxYkQ0NTc2YUQ5NjJmMjgwNmIxMTI5OThDYzNCNzJGYyIsInRvX2FkZHJlc3MiOiIweDgyMThhMEY0N0Y0YzBkRTBjMTc1NGY1MDg3NDcwN2NkNmU3YjJlNWUiLCJnYXNfbGltaXQiOjIxMDAwLCJtYXhfZmVlX3Blcl9nYXMiOiIyNjAwMDAwMDAwMCIsIm1heF9wcmlvcml0eV9mZWVfcGVyX2dhcyI6IjIwNTIwMDAwMDAwIiwiYW1vdW50IjoiOTAwMDAwMDAwMDAwMDAwMDAwMCIsImNvbnRyYWN0X2FkZHJlc3MiOiIweGIxMmMxM2U2NmFkZTFmNzJmNzE4MzRmMmZjNTA4MmRiOGMwOTEzNTgiLCJzaWduYXR1cmUiOiI0YjJhN2ZjMjM5MzFkYTJjYWFlODNiYTgyYjViOTViMjY0NGVjMzE3MDc4ZjEwNWQxMDk4YjdjZWRlNThiYzVkMjU2NjBlZjBhYWVmMWYxMjEwMjA1MjkwMDA3Y2I4MWFiYmVhMDFiZWI5YWFkZmIyZmRkNDE1MGVkMWFmZjAzOTAwIn0="
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.buildSignedTransaction
```
- response
```
{
  "code": "SUCCESS",
  "msg": "0xa1214d8d892966da2282167e1693db707274811c0d947304af6f31293b0817d9",
  "signed_tx": "0x02f8b10a058504c7165a0085060db8840082520894b12c13e66ade1f72f71834f2fc5082db8c09135880b844a9059cbb0000000000000000000000008218a0f47f4c0de0c1754f50874707cd6e7b2e5e0000000000000000000000000000000000000000000000007ce66c50e2840000c080a04b2a7fc23931da2caae83ba82b5b95b2644ec317078f105d1098b7cede58bc5da025660ef0aaef1f1210205290007cb81abbea01beb9aadfb2fdd4150ed1aff039"
}

```

## 2.send tx

- request
```
grpcurl -plaintext -d '{
  "chain": "Optimism",
  "rawTx": "0x02f8b10a058504c7165a0085060db8840082520894b12c13e66ade1f72f71834f2fc5082db8c09135880b844a9059cbb0000000000000000000000008218a0f47f4c0de0c1754f50874707cd6e7b2e5e0000000000000000000000000000000000000000000000007ce66c50e2840000c080a04b2a7fc23931da2caae83ba82b5b95b2644ec317078f105d1098b7cede58bc5da025660ef0aaef1f1210205290007cb81abbea01beb9aadfb2fdd4150ed1aff039"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.SendTx
```
- response
```
{
  "code": "SUCCESS",
  "msg": "",
  "txHash": "0xbd059a5975ffdd5547b56811453afdbace10c36861a71a4256dc3a004d25977d"
}

```
## 3.convertAddress
```
grpcurl -plaintext -d '{
"chain": "Optimism",
"network": "",
"publicKey": "048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
```
{
"code": "SUCCESS",
"msg": "convert address success",
"address": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
}
```
## 4.validAddress
```
grpcurl -plaintext -d '{
"chain": "Optimism",
"address": "0x4740d7eE1bD4576aD962f2806b112998Cc3B72Fc"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.validAddress
```
```
{
"code": "SUCCESS",
"msg": "valid address",
"valid": true
}
```
## 5.getBlockByNumber
```
grpcurl -plaintext -d '{
"height": "132043390",
"chain": "Optimism"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByNumber
```
```
{
"code": "SUCCESS",
"msg": "GetBlockByNumber success",
"height": "132043390",
"hash": "0xd5fadb5b5116145952f6d2c2c689bc980f8db8e4c46005c4b5186dc468b4a3ab",
"base_fee": "0x13c9",
"transactions": [
{
"from": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001",
"to": "0x4200000000000000000000000000000000000015",
"token_address": "",
"contract_wallet": "",
"hash": "0x37d6fa72c81c1b4ebdd7d1ef9c0c808d6da1d54ddc219dfa73dadd9e18a485a6",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0x293e9accd6e0c9ce99c4f8e3d8efc0236a29a596",
"to": "0x887290c34856cd3ba1b84da78cccf43812f66324",
"token_address": "",
"contract_wallet": "",
"hash": "0xa42712654d1389fe37c83c5be524657dcebec15b6de1611607bafc1fbc7b8303",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xb8ff877ed78ba520ece21b1de7843a8a57ca47cb",
"to": "0xfe6507f094155cabb4784403cd784c2df04122dd",
"token_address": "",
"contract_wallet": "",
"hash": "0x1577c8a383d26880b2a59401452f8a2512d40e9d578a03e4b3daa6353b777e37",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xb8ff877ed78ba520ece21b1de7843a8a57ca47cb",
"to": "0xa7b5189bca84cd304d8553977c7c614329750d99",
"token_address": "",
"contract_wallet": "",
"hash": "0x0e1042a425e17e4e205b706928cc3a8c121812b651c7fac28c06673741efdf26",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0x86826caba2493439cebc9028bd89a6b7e3c0d809",
"to": "0x0daf895a78eb49151a5f1003818939770d3ca7dd",
"token_address": "",
"contract_wallet": "",
"hash": "0xd8e81b66ca00ad99ea4c1f82a0037552e75391b586c9c99750076d6f42b14b7f",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xb00a05613c746afe4c27eff93accfdf628b2ed97",
"to": "0x4200000000000000000000000000000000000006",
"token_address": "",
"contract_wallet": "",
"hash": "0x34b175d9b9b050c2a55d07cb64b0a5670be340ea39c8c1e98bc1426c16b5eab3",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0x368b1c6f9779c2895f448460a01a4c8ab7a26992",
"to": "0xdc6ff44d5d932cbd77b52e5612ba0529dc6226f1",
"token_address": "",
"contract_wallet": "",
"hash": "0xb411a435eed7b82c51853012817afd9907438c19a11d213f5b50c2bb82067b08",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xacd03d601e5bb1b275bb94076ff46ed9d753435a",
"to": "0x94b008aa00579c1307b0ef2c499ad98a8ce58e58",
"token_address": "",
"contract_wallet": "",
"hash": "0x2b2bf5871f182b9d97e8a4608a193451557b3c50867f3451e7d58c3affcd3c1e",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xacd03d601e5bb1b275bb94076ff46ed9d753435a",
"to": "0x0b2c639c533813f4aa9d7837caf62653d097ff85",
"token_address": "",
"contract_wallet": "",
"hash": "0x52e26ddbd262ffa70e3491bbad35bb57e02159a3623deb85fed7a1eec6e294ba",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xcac1a7ce337957bd737aba084fe0441392e97b21",
"to": "0x6b3872e786db187c311479ed3d1513a244e31b68",
"token_address": "",
"contract_wallet": "",
"hash": "0xc16b1d6b1698b21f1c2acab8cfe71c715ff8355f6cfe18694afdda56e175095c",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xee1049c9961459a09bcec9adb623fa48d6a056ec",
"to": "0xc5bf05cd32a14bffb705fb37a9d218895187376c",
"token_address": "",
"contract_wallet": "",
"hash": "0xa6abac09dec1ac50bd732b3317e38885e3be7e4438aac80c4c8080a71a99f41e",
"height": "132043390",
"amount": "0x3b9aca00"
},
{
"from": "0xcac1a7ce337957bd737aba084fe0441392e97b21",
"to": "0x6b3872e786db187c311479ed3d1513a244e31b68",
"token_address": "",
"contract_wallet": "",
"hash": "0x1f4a1157418577cac384a2d3b52dbd72a83d19501bc7c84f8bee7eaa522042dd",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xcac1a7ce337957bd737aba084fe0441392e97b21",
"to": "0x6b3872e786db187c311479ed3d1513a244e31b68",
"token_address": "",
"contract_wallet": "",
"hash": "0x91efa06390439ef44c8d29bc76b9448fdc77a7efde5d787e446708c68838009d",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0x6ae78e7b25862f7bd7bebd2273f5a82ba2209a67",
"to": "0x9d1b033ac8bff2b07fb7d13385b8c270db25f96f",
"token_address": "",
"contract_wallet": "",
"hash": "0xa73d7ff87db53ee3d0ae9f40ae79132e8f095238ad937f9746b6ad9898d824c7",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0x5f970cccd87685b259af4df123454de5aeecf14d",
"to": "0x887290c34856cd3ba1b84da78cccf43812f66324",
"token_address": "",
"contract_wallet": "",
"hash": "0x42713985344cc5c61f0318dda55f5bd5b2821ee911ee2e6ef2953ea9e215ffd3",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xcac1a7ce337957bd737aba084fe0441392e97b21",
"to": "0x6b3872e786db187c311479ed3d1513a244e31b68",
"token_address": "",
"contract_wallet": "",
"hash": "0xd8760c0dc2acfcb58f92f6db24b0b566f83159f489c5f07a6599ebd62509d3b1",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xcac1a7ce337957bd737aba084fe0441392e97b21",
"to": "0x6b3872e786db187c311479ed3d1513a244e31b68",
"token_address": "",
"contract_wallet": "",
"hash": "0x49163b916ec4f9e89b0eb81726bc2824b11280e71b6276758531b66df21f2286",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0x368664b287fe160fbd15c7429b273f9d2cd910d5",
"to": "0x80942a0066f72efff5900cf80c235dd32549b75d",
"token_address": "",
"contract_wallet": "",
"hash": "0xc44a62827a04aa3667a23d2f63c60eeb9faac0dd41795a413bae05085b58ffbc",
"height": "132043390",
"amount": "0x0"
},
{
"from": "0xf81f02a4e9e8eafa941e6467962d4badd3eb62f0",
"to": "0x7bbb754d03b6e4439764a56475ed8adc71c59ba8",
"token_address": "",
"contract_wallet": "",
"hash": "0xa1c1f3a734335e2ed6b2e18aa60eb3243b41adb559ee0e9d7f55e72e3998288e",
"height": "132043390",
"amount": "0x0"
}
]
}
```
## getBlockByHash
```
grpcurl -plaintext -d '{
"hash": "0x6f7e0f205950dce59896980b53c58d152baafcb085afe6ad0cd9fcd57efcd6b5",
"chain": "Optimism"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockByHash
```
```
{
"code": "SUCCESS",
"msg": "GetBlockByNumber success",
"height": "132043396",
"hash": "0x6f7e0f205950dce59896980b53c58d152baafcb085afe6ad0cd9fcd57efcd6b5",
"base_fee": "0x13bc",
"transactions": [
{
"from": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001",
"to": "0x4200000000000000000000000000000000000015",
"token_address": "",
"contract_wallet": "",
"hash": "0xb8105f4713acf821d553168a9dc4ff3be5195714eca8670792f65cfbad779889",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0x293e9accd6e0c9ce99c4f8e3d8efc0236a29a596",
"to": "0x887290c34856cd3ba1b84da78cccf43812f66324",
"token_address": "",
"contract_wallet": "",
"hash": "0x6c7ab41a2f1777eb9e3d03d2ad527ecbc22b300d9e3472c101d088b708cac7c7",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0xe37f7c80ced04c4f243c0fd04a5510d663cb88b5",
"to": "0xf1fcb4cbd57b67d683972a59b6a7b1e2e8bf27e6",
"token_address": "",
"contract_wallet": "",
"hash": "0x769106eb282e4d86a9f8dc1161348150c56797fbf3210e48cb71810b06f914c2",
"height": "132043396",
"amount": "0x17cf6739963b"
},
{
"from": "0x07ae8551be970cb1cca11dd7a11f47ae82e70e67",
"to": "0x6f26bf09b1c792e3228e5467807a900a503c0281",
"token_address": "",
"contract_wallet": "",
"hash": "0x58950cd084fefae3fa40b5288752d91e31a515713526652d068379d470997ee9",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0x4b3a308d74ae98274a590db6b25ed5d5030e1cd1",
"to": "0x9560e827af36c94d2ac33a39bce1fe78631088db",
"token_address": "",
"contract_wallet": "",
"hash": "0xd2def29da89f0b0b4e553b90619af979a7a9e5ebcba30b6407e4249e91f07bd3",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0x83681c14770b44361e21faf91d8325423365ea5c",
"to": "0xbe931744bf0b4c4c580d375de04ee7fc2f52c568",
"token_address": "",
"contract_wallet": "",
"hash": "0x7292e805fda2ba7c932f6bef3dfc4eedd8da3d34bc362b49fc6af325f85fb4ec",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0x368b1c6f9779c2895f448460a01a4c8ab7a26992",
"to": "0xdc6ff44d5d932cbd77b52e5612ba0529dc6226f1",
"token_address": "",
"contract_wallet": "",
"hash": "0xf170ce534f6cb63ad1ceeeb6370c281411f96843452b2fe258113d7bf56d8fa9",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0xf7bd3cb6730f38ab4395338f1bdb154ad4189961",
"to": "0x0daf895a78eb49151a5f1003818939770d3ca7dd",
"token_address": "",
"contract_wallet": "",
"hash": "0xd7571790af040cb9907a80ed774442ddd8c6ccca879046ecd2d9976576fdce15",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0xacd03d601e5bb1b275bb94076ff46ed9d753435a",
"to": "0x66f3fe2083513a4a1a1e872c7f58765f4e59cb4c",
"token_address": "",
"contract_wallet": "",
"hash": "0x359f14d0eb81a6e70ec9f36ce99584a8db0e8f40d186bbd46fa73371fbeac5fa",
"height": "132043396",
"amount": "0xa16f707298000"
},
{
"from": "0x32372e809d99bb6e3de073b279c6199670f463cb",
"to": "0x2dd224a09e73ef83129f47ded4731ee5a45fddf9",
"token_address": "",
"contract_wallet": "",
"hash": "0x5841a6567e9c1e4330913094f3c85ef53938fbf91ed5913cce9fa1dcc775a22b",
"height": "132043396",
"amount": "0x1a38fafb165000"
},
{
"from": "0x8f06ac73715a230df41b0c7dd10b04502c16c05d",
"to": "0xd7ba4057f43a7c4d4a34634b2a3151a60bf78f0d",
"token_address": "",
"contract_wallet": "",
"hash": "0xcb9cdad0a911bbf9cc90671868c3c2c23d906ce22e7f2662fc86c5a48eb7ccc4",
"height": "132043396",
"amount": "0x167ad3f1f0c7"
},
{
"from": "0xee1049c9961459a09bcec9adb623fa48d6a056ec",
"to": "0xc5bf05cd32a14bffb705fb37a9d218895187376c",
"token_address": "",
"contract_wallet": "",
"hash": "0xece9bfab10157f873890c1467bfb6c9f3507d43f41551ce0ce721de1ad7924d6",
"height": "132043396",
"amount": "0x3b9aca00"
},
{
"from": "0x1820532135321f7d59e77dc2f90dd54dfdbb0853",
"to": "0x8369ed7ed41c52362d4c9ce5a7565c0e00e36c5d",
"token_address": "",
"contract_wallet": "",
"hash": "0x6fe5bb2c13553cde402537c9df6c2f59e3e315e443cc9b382ac4c7c0da43cebc",
"height": "132043396",
"amount": "0x1bc83fa224c000"
},
{
"from": "0x6ae78e7b25862f7bd7bebd2273f5a82ba2209a67",
"to": "0x9d1b033ac8bff2b07fb7d13385b8c270db25f96f",
"token_address": "",
"contract_wallet": "",
"hash": "0x6987da03bc0005334128ee43f8d5c24713c851efaebbf79c923bf9863cd67678",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0x9eb645904290b2e40ef833e0f289b873a47c658d",
"to": "0xdca2e9ae8423d7b0f94d7f9fc09e698a45f3c851",
"token_address": "",
"contract_wallet": "",
"hash": "0xa32d010679935b95bad7cc0ce039c756cbcf0e649024f53b7e1d73ddc8371f1b",
"height": "132043396",
"amount": "0x6c05e0abed000"
},
{
"from": "0x5f970cccd87685b259af4df123454de5aeecf14d",
"to": "0x887290c34856cd3ba1b84da78cccf43812f66324",
"token_address": "",
"contract_wallet": "",
"hash": "0x3dc92900349dc07a3a83bbfa96597106c25a60ea0e6cedd634193b8ea818f096",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0xf70da97812cb96acdf810712aa562db8dfa3dbef",
"to": "0xcf2b5a6ac85658b8557ff96c10a96f2d01a1ef2f",
"token_address": "",
"contract_wallet": "",
"hash": "0x08c704b56ea4f4057785ab715ccd11ee8c109e3fefc9cde2026330b17084ed9b",
"height": "132043396",
"amount": "0x8ec49f12e3f1"
},
{
"from": "0x368664b287fe160fbd15c7429b273f9d2cd910d5",
"to": "0x80942a0066f72efff5900cf80c235dd32549b75d",
"token_address": "",
"contract_wallet": "",
"hash": "0x49ba8dce24be19db41dabd3f7dc16b4a740593a50ad28c6af31601b9fee00ecd",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0xf05d44ffba5909fc228be94e282263a97c0c6f12",
"to": "0x7bbb754d03b6e4439764a56475ed8adc71c59ba8",
"token_address": "",
"contract_wallet": "",
"hash": "0x96343a25570aff0af4567a8b621487d42870104c5327a1e61fa565473184e9fa",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0x0a7b03f575196c2946d5980b59a5ea02128434f4",
"to": "0x2c84370daddbcd67d729689671a9fe63df39cf13",
"token_address": "",
"contract_wallet": "",
"hash": "0x4bfe6597dde6a7d0d5772c0cb99e6d0f9018b5e6b519fc307683cc4459589a00",
"height": "132043396",
"amount": "0x0"
},
{
"from": "0xad4796d21431a3d4272a681ce67b0236f1a5783a",
"to": "0x6ff6a738233a3b728d4a67ea92b0e2416e316149",
"token_address": "",
"contract_wallet": "",
"hash": "0x337cdc2876dfd3aca1bd81c46ce068e1d80d65d792bf6ce22eba66c26e06bedd",
"height": "132043396",
"amount": "0x0"
}
]
}
```

## getBlockHeaderByHash
```
grpcurl -plaintext -d '{
"chain": "Optimism",
"hash": "0x6f7e0f205950dce59896980b53c58d152baafcb085afe6ad0cd9fcd57efcd6b5"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByHash
```

```
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0x6f7e0f205950dce59896980b53c58d152baafcb085afe6ad0cd9fcd57efcd6b5",
    "parent_hash": "0x4b5428d67ff7aec4646af1bf88f430fc10900d4e1a6fb8bd9d926b26728ddd79",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x4200000000000000000000000000000000000011",
    "root": "0x3badd5e799404edf7f08ce1d7a84c3e1b1931a82ffdd26b40d106f25c2cb1c12",
    "tx_hash": "0xd1ab479ba9102deaa98bf9dc12236e69398844740a1351de6d22a8a78ddb26ab",
    "receipt_hash": "0xb13b700167943bcc19ca2d2277b82f780eb97c94cc43f1a1034f2f962b5f1d72",
    "parent_beacon_root": "0x7f9a8c985dc4d791e2c6b8b546b90c9b09f2950d18bbbdb927ce945e21611a5a",
    "difficulty": "0",
    "number": "132043396",
    "gas_limit": "60000000",
    "gas_used": "8323941",
    "time": "1739685569",
    "extra": "AAAAAPoAAAAG",
    "mix_digest": "0xeaf900d47a0d4f32b7f8714b4e85e013531f81bb336e19046e053f2bee257060",
    "nonce": "0",
    "base_fee": "5052",
    "withdrawals_hash": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}
```
## getBlockHeaderByNumber
```
grpcurl -plaintext -d '{
"chain": "Optimism",
"height": "132043396"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByNumber
```


```
{
  "code": "SUCCESS",
  "msg": "get latest block header success",
  "block_header": {
    "hash": "0x6f7e0f205950dce59896980b53c58d152baafcb085afe6ad0cd9fcd57efcd6b5",
    "parent_hash": "0x4b5428d67ff7aec4646af1bf88f430fc10900d4e1a6fb8bd9d926b26728ddd79",
    "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "coin_base": "0x4200000000000000000000000000000000000011",
    "root": "0x3badd5e799404edf7f08ce1d7a84c3e1b1931a82ffdd26b40d106f25c2cb1c12",
    "tx_hash": "0xd1ab479ba9102deaa98bf9dc12236e69398844740a1351de6d22a8a78ddb26ab",
    "receipt_hash": "0xb13b700167943bcc19ca2d2277b82f780eb97c94cc43f1a1034f2f962b5f1d72",
    "parent_beacon_root": "0x7f9a8c985dc4d791e2c6b8b546b90c9b09f2950d18bbbdb927ce945e21611a5a",
    "difficulty": "0",
    "number": "132043396",
    "gas_limit": "60000000",
    "gas_used": "8323941",
    "time": "1739685569",
    "extra": "00000000fa00000006",
    "mix_digest": "0xeaf900d47a0d4f32b7f8714b4e85e013531f81bb336e19046e053f2bee257060",
    "nonce": "0",
    "base_fee": "5052",
    "withdrawals_hash": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "blob_gas_used": "0",
    "excess_blob_gas": "0"
  }
}
```

## getAccount
```
grpcurl -plaintext -d '{
"chain": "Optimism",
"address": "0x93Bf4C1383A86cC49Ce536AD5643207F9F0eB733",
"contractAddress": "0x00"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```

```
{
  "code": "SUCCESS",
  "msg": "",
  "network": "",
  "account_number": "",
  "sequence": "231",
  "balance": "2100308755671059"
}

```

## getFee
```
grpcurl -plaintext -d '{
"chain": "Optimism"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getFee
```

```
{
  "code": "SUCCESS",
  "msg": "get gas price success",
  "slow_fee": "2664668|1000000",
  "normal_fee": "2664668|1000000|*2",
  "fast_fee": "2664668|1000000|*3"
}
```

## getTxByAddress
```
grpcurl -plaintext -d '{
"chain": "Optimism",
"address": "0x4740d7eE1bD4576aD962f2806b112998Cc3B72Fc",
"contractAddress": "0xb12c13e66AdE1F72f71834f2FC5082Db8C091358"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByAddress
```

```
{
  "code": "SUCCESS",
  "msg": "get tx list success",
  "tx": [
    {
      "hash": "0xa31ba448123cdf0a254cad8c5335bdf501023d21124a8d522cb7c91e75452aac",
      "index": 0,
      "froms": [
        {
          "address": "0xa86ca428512d0a18828898d2e656e9eb1b6ba6e7"
        }
      ],
      "tos": [
        {
          "address": "0x4740d7ee1bd4576ad962f2806b112998cc3b72fc"
        }
      ],
      "fee": "0xa31ba448123cdf0a254cad8c5335bdf501023d21124a8d522cb7c91e75452aac",
      "status": "Success",
      "values": [
        {
          "value": "173637638419527000000"
        }
      ],
      "type": 1,
      "height": "128544335",
      "contract_address": "0xef4461891dfb3ac8572ccf7c794664a8dd927945",
      "datetime": "",
      "data": ""
    },
    {
      "hash": "0xcec2178abe621b35506d4028bcea828656377bef402f4f121ddf2ae9ee1e5e3f",
      "index": 0,
      "froms": [
        {
          "address": "0x9cd330770dd6d3b85b999dd53169e75d113419e3"
        }
      ],
      "tos": [
        {
          "address": "0x4740d7ee1bd4576ad962f2806b112998cc3b72fc"
        }
      ],
      "fee": "0xcec2178abe621b35506d4028bcea828656377bef402f4f121ddf2ae9ee1e5e3f",
      "status": "Success",
      "values": [
        {
          "value": "87542533161815452889"
        }
      ],
      "type": 1,
      "height": "132043292",
      "contract_address": "0xb12c13e66ade1f72f71834f2fc5082db8c091358",
      "datetime": "",
      "data": ""
    }
  ]
}
```

## getTxByHash
```
grpcurl -plaintext -d '{
"hash": "0xa31ba448123cdf0a254cad8c5335bdf501023d21124a8d522cb7c91e75452aac",
"chain": "Optimism"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getTxByHash
```

```
{
  "code": "SUCCESS",
  "msg": "get tx by hash success",
  "tx": {
    "hash": "0xa31ba448123cdf0a254cad8c5335bdf501023d21124a8d522cb7c91e75452aac",
    "index": 40,
    "froms": [
      {
        "address": ""
      }
    ],
    "tos": [
      {
        "address": ""
      }
    ],
    "fee": "101564",
    "status": "Success",
    "values": [
      {
        "value": "<nil>"
      }
    ],
    "type": 0,
    "height": "128544335",
    "contract_address": "",
    "datetime": "",
    "data": "B69D1A08000000000000000000000000000000000000000000000000000000000000FDB800000000000000000000000000000000000000000000000969B4D02BC1A05FC000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000012ED3D6A52AD3F99F9D1D17D1E32BE1012449C08E8E0A0CCB4D4AA31469F0DF9634BB11F10DB39092857CE77753F8A4A1470FF2081E9867DF342A7AD47B4A0ECC8A7C01828AA6595E0F99BF3F596FA32A332A0A9542CC3897778FD52B2FA786C131BFB097724E4CDC676E32AB4AE8AEBBA6265F003F61C21ED56C5E7D7363A99098204DC0BAE2325A0DBB616816D1D42FAB3BDE9459A5B391F2FD2C8233579A82F51B48F7800909A6E072A42EB5F96090E70F2F05656D835813D2628E3137B652F3F58C0D74E18B3DB059FE504758903A649C451A52999A3E8B3C0D3B5708EFD5FF78598601469BA7B004AC508ED7F1C2C6347A1E10E46F5AFFC749C3C1B7059765188A581217B0EED315E26CBAFE5C0862DA0DFB23B39D75DAA376E96DAD728B95BF1EA02F1EDC5B4B58EC8F3493E2D60DC54288B6A15F82B167F05429469444144886CADCE7E14FE3BEE1756DF5CFEA66142C2BDAEEA40DC0BB05251530E4141D9E069D810C7339CEA624DB970C4A1DD0E1264FCE33825BB8B815E5BE30CDEDAB89593AA8C2E291B8A6002739DF7567F27F56C4464DB1D2B21BDFC6C8ED5844228944C0AA488CEECBEE28C26ED95FCFC5952895DAB6EBB5D7D896C81ECEB013106E70B0B695B149D4FBC05A026A73D8B5062AD54BA70BB44839B832526213B965511E8797521231C93FC209C76DF102073011E04A52EF7C3CCCA31891F035AFCBC4621D4B51A86903B08962F7492220BF477B99BBC6ED644649F39F77436AEB7F6D5FDB64CC77443F66CFB98233E53D39AB30406F91B945D5570153BABF71B52"
  }
}
```

## getBlockHeaderByRange
```
grpcurl -plaintext -d '{
"start": "128544330",
"end": "128544331",
"chain": "Optimism"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getBlockHeaderByRange

```

```
{
  "code": "SUCCESS",
  "msg": "get block range success",
  "block_header": [
    {
      "hash": "",
      "parent_hash": "0x4b69b9814ba686f479d7440c30695b1b22f74570bb3a87ee01621c5d2b4fa52d",
      "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
      "coin_base": "0x4200000000000000000000000000000000000011",
      "root": "0xbc06ef564a11352b70dac288266c3234f2d789c0a25bcdf195b381c6e3f5706f",
      "tx_hash": "0x78ccab0e58747087946aa4c130340209dd67adf1e28c1e55596d2cd3dbdf7d8d",
      "receipt_hash": "0xe30c67b381454c83d211ec5821ccf1644c242337eb6ef3cfa46f77e164b8d942",
      "parent_beacon_root": "0x65813c625b87998306e1d3215aedf7f436c23e62a4bc23791c779c6eaf78e77b",
      "difficulty": "0",
      "number": "128544330",
      "gas_limit": "60000000",
      "gas_used": "7242312",
      "time": "1732687437",
      "extra": "",
      "mix_digest": "0x5f349ed26647062ea45dc11a11a601440455be63e5d2e6e4f32558fec30c4328",
      "nonce": "0",
      "base_fee": "390",
      "withdrawals_hash": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
      "blob_gas_used": "0",
      "excess_blob_gas": "0"
    },
    {
      "hash": "",
      "parent_hash": "0x58b503c589294f62ebb76a563d112f14e6d57fdbec9482d4622e6e00afa1ee36",
      "uncle_hash": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
      "coin_base": "0x4200000000000000000000000000000000000011",
      "root": "0x7d532c1646fd1f5177341c6c51fd454941c9c43108f4a62edaaddc9ce2fabcde",
      "tx_hash": "0xf419a15d4c250ccfd1b23b2513ed9553ccfa67d625dc19f84cb4fa18b0fed218",
      "receipt_hash": "0x56e9b6f8348212981d1c7c8f40544860b1c974bd91bacc7f549a60ace4e70900",
      "parent_beacon_root": "0x65813c625b87998306e1d3215aedf7f436c23e62a4bc23791c779c6eaf78e77b",
      "difficulty": "0",
      "number": "128544331",
      "gas_limit": "60000000",
      "gas_used": "7820656",
      "time": "1732687439",
      "extra": "",
      "mix_digest": "0x5f349ed26647062ea45dc11a11a601440455be63e5d2e6e4f32558fec30c4328",
      "nonce": "0",
      "base_fee": "390",
      "withdrawals_hash": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
      "blob_gas_used": "0",
      "excess_blob_gas": "0"
    }
  ]
}
```

