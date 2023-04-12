# Cache Store
Store implementation that adds a layer of caching to any Store

## Usage
```
store := store.New(
    generated.Schema(),
    cache.Factory(existing_store))
```
