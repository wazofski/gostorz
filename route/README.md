# Route Store
Route store allows mapping object types to different Stores

## Usage
```
store := store.New(
    generated.Schema(),
    route.Factory(deault_store,
        route.Mapping("type1", store1),
        route.Mapping("type2", store2)))
```
