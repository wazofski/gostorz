# REST Server
REST API Server that exposes Store functionality

## Usage
```
srv := rest.Server(generated.Schema(), store_to_expose,
    rest.TypeMethods(generated.WorldKind(),
        rest.ActionGet, rest.ActionCreate,
        rest.ActionDelete, rest.ActionUpdate),
    rest.TypeMethods("AnotherWorld", rest.ActionGet))
)

// use cancel function to stop server
cancel = srv.Listen(port) // does not block
```
