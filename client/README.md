# REST Client Store
Store implementation that uses a remote [storz REST Server](https://github.com/wazofski/gostorz/tree/main/rest) as storage

## Usage
```
store := store.New(
    generated.Schema(),
    client.Factory("http://server-host:port",
        client.Header("A", "B"), ...// headers
    ))
```
