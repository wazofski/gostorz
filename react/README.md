# React Store
React store allows attaching callbacks to object actions

## Usage
```
func WorldCreateCb(store.Object, store.Store) error {
    // ...
    return nil
}

func WorldDeleteCb(store.Object, store.Store) error {
    // ...
    return nil
}

store := store.New(
    generated.Schema(),
    react.ReactFactory(underlying_store,
        react.Subscribe(generated.WorldKind(), react.ActionCreate, WorldCreateCb),
        react.Subscribe(generated.WorldKind(), react.ActionDelete, WorldDeleteCb),
    ))
```
