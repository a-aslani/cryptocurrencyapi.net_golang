# cryptocurrencyapi.net with Go
Implementation of cryptocurrency payment through cryptocurrencyapi.net service and payment confirmation through IPN

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