# Object Browser Server
Exposes a web tool for browsing objects inside a Store

## Usage
```
srv := browesr.Server(generated.Schema(), store_to_expose)

// use cancel function to stop server
cancel = srv.Listen(port) // does not block
```
