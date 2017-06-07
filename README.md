## Golang API of btc-e.com exchange

It is the simple library wrapper for btc-e.com. 

## Donate
If you like the library please donate some coins on follow addresses:

ETH: 0xd6ed497f6a034cd28762df9df3cd2c5b5d69ce6b
Zcash: t1S7MREH6zSGQQ9Htr2jrtyvDUKnPFLJeor 
LTC: LfLDPm4rAfE5rrbjXn5NFjsx4HdXq8rV3K

## In progress
 - Working with V3 version of btc-e.com

## Examples

```go
    // You can change the target url here. For example:
    btceapi.ApiURL = "https://btc-e.nz"	
    
    key := "Key"
    secret := "Secret"

    btceAPI := btceapi.BtceAPI{Key: key, Secret: secret}
    res, err := btceApi.GetInfo()

    if err != nil {
        fmt.Println(res.Rights.Info)
    }
```