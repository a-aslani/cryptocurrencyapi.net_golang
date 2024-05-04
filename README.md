# cryptocurrencyapi.net with Go
Implementation of cryptocurrency payment through cryptocurrencyapi.net service and transaction confirmation through IPN

## IPN data verification algorithm
To check the signature you need to:
1) remember the received signature (sign field) from the data
2) remove the sign field from the data
3) sort the data by key in ascending order
4) form a string - concatenate the data values (without keys) through ':', if the value has an 'array' type (labels, ids), then use the 'Array' value for gluing
5) add to the string ':' and the MD5 hash from YOUR_API_KEY
6) get the SHA1 hash of the string (in php this is the sha1() function)
7) compare signature from step 1 with hash from step 6

```go
keys := make([]string, 0)

m := structs.Map(inputStruct)

for k, _ := range m {
    if k != "Sign" {
        keys = append(keys, k)
    }
}

sort.Strings(keys)

values := make([]string, len(keys))
for i, k := range keys {
    values[i] = fmt.Sprintf("%v", m[k])
}

signData := strings.Join(values, ":")
hashKey := fmt.Sprintf("%x", md5.Sum([]byte("YOUR_API_KEY")))
sign := fmt.Sprintf("%s:%s", signData, hashKey)
hashSign := fmt.Sprintf("%x", sha1.Sum([]byte(sign)))

if hashSign != request.Sign {
    fmt.Println("sign wrong")
} else {
    fmt.Println("OK")
}
```

## Give method
Blockchain: All

Cost: Free if the address was issued previously or 1 unit + 1 unit if the address is virtual

Provides an address for incoming payments. Returns an address, its public key and (if enabled by settings) its private key. Typically, such addresses are temporary or transit.

If you pass in the address or label or uniqID parameter a value from an address that is in the archive, the API will restore and return this address
Parameters:

- token - additionally display the balance of this token at the address
- from - from this address take a commission for sending tokens [main address]
- to - send coins to this address [from settings or main address]
- statusURL - IPN handler URL [from settings] if you specify '-', then IPN will not be sent for this address
- label - label, transferred to IPN
- uniqID - unique request ID, which corresponds to the issued address [from label]. if you repeat uniqID, the API will return the previously issued address, but the rest of the address parameters will be updated
- waitPeriod - time before archiving in minutes (max. waiting period for funds to arrive)
- period - time before archiving in minutes (max. waiting period for the next receipt of funds)
- reusable - use in the address pool (see point 7) [0 - no]
- forceGenerate - force a new one [0 - no]
- address - create a “virtual” address (see point 9) or return from archive
- privateKey - create address from such a private key
- publicKey - create a virtual address from such a public key
- isPrimary - create a primary address [0 - no]
- qr - return a QR code image [0 - no]

```
https://new.cryptocurrencyapi.net/api/trx/.give?key=YOUR_API_KEY&label=u1&period=30
```

```json
{
  "result": {
    "address":"TKavpKP2VJbfV4AGyi3MrhT6FJAP5eJ2RR",
    "publicKey":"03b7a403f5ffce0292b8f17328b7f0dd0ca52ac5549292bc6fe1fe580ef40d183f",
    "privateKey":"***",
    "QR": "data:image/png;base64, ..."
  }
}
```