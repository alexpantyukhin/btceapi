## Golang API of WEXnz exchange

It is the simple library wrapper for WEXnz. 

## Donate
If you like the library please donate some coins on follow addresses:

    ETH: 0xd6ed497f6a034cd28762df9df3cd2c5b5d69ce6b
    LTC: LfLDPm4rAfE5rrbjXn5NFjsx4HdXq8rV3K

## Examples

```go
    key := "Key"
    secret := "Secret"

    btceAPI := btceapi.BtceAPI{Key: key, Secret: secret}
    res, err := btceAPI.GetInfo()

    if err == nil {
        fmt.Println(res.Rights.Info)
    }
```
