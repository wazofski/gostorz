# Memory Store
Store implementation that usees an in-memory hash map for storage
Useful for temporary storage scenarios like caching and testing

## Usage
```
store := store.New(
    generated.Schema(),
    memory.Factory())
```
