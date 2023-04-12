# Mongo Store
Store implementation that usees an instance of Mongo DB for persistification

## Usage
```
store := store.New(
    generated.Schema(),
    mongo.Factory("mongodb://path:27017/", "mdb"))
```
