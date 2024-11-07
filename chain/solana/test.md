# Test for grpc api

## 1.support chain
- request
```
grpcurl -plaintext -d '{
  "chain": "Solana"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getSupportChains
```
- response
```
{
  "code": "SUCCESS",
  "msg": "Support solana chain",
  "support": true
}
```

## 2.convert address

- request
```
grpcurl -plaintext -d '{
  "chain": "Solana",
  "publicKey": "6488e3d824c6eb210b97f7f5c49d2e5f0ff63b289f2b2f861af31fa09421163d"
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.convertAddress
```
- reponse

```
{
  "code": "SUCCESS",
  "msg": "convert address success",
  "address": "7mSqVJpb8ziMDB7yDAEajeANyDosh1WK5ksS6mCdDHRE"
}
```

## 3.valid address

- request
```
grpcurl -plaintext -d '{
  "chain": "Solana",
  "address": "7mSqVJpb8ziMDB7yDAEajeANyDosh1WK5ksS6mCdDHRE"  
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



## 4.get account 

- request
```
grpcurl -plaintext -d '{
  "chain": "Solana",
  "address": "7mSqVJpb8ziMDB7yDAEajeANyDosh1WK5ksS6mCdDHRE",
  "contractAddress": "FaP4Ti84eCuGibYNeGMTKCZW9YyyZHgoSB6nFViGtBdy" 
}' 127.0.0.1:8189 dapplink.account.WalletAccountService.getAccount
```
>contractAddress is nonce account
- response
```
{
  "code": "SUCCESS",
  "msg": "get account response success",
  "network": "",
  "account_number": "0",
  "sequence": "DRCR8BDqFzF2LnvdgnKfZH8UPb98Ho4Ci5aFT9CDC1fV",
  "balance": "0.117326118"
}
```
