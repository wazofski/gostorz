# Store Definitions

`Store` interface defined below exposes
five main operations over Objects allowing various
options like pagination, filtering and sorting.

**The interface is fully exposed from the `store` package to allow custom store implementations if needed.**

```
  Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
  List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
  Create(context.Context, Object, ...CreateOption) (Object, error)
  Delete(context.Context, ObjectIdentity, ...DeleteOption) error
  Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
```

## Create an object
```
// Given a Store
var str store.Store = ...

// Initialize the object
world := generated.WorldFactory()
world.External().SetName("abc")

world, err = str.Create(ctx, world)
```

## Update an object
```
world.External().SetDescription("abc")

ret, err = clt.Update(ctx, 
  world.Metadata().Identity(), world)
```

## Delete an object
```
err = str.Delete(ctx, generated.WorldIdentity("abc"))

err = clt.Delete(ctx, world.Metadata().Identity())
```

## Get an object by Identity or PKey
```
world, err = str.Get(ctx, generated.WorldIdentity("abc"))

world, err = str.Get(ctx, world.Metadata().Identity())
```

## List all World objects
```
world_list, err = str.List(ctx, generated.WorldKindIdentity())
```

## List World objects with a property filter
```
world_list, err = str.List(ctx,
    generated.WorldKindIdentity(),
    options.PropFilter("external.name", "abc"))
```

## List World objects with a primary key filter
```
world_list, err = str.List(ctx,
    generated.WorldKindIdentity(),
    options.KeyFilter("a", "b", "c"))
```

## List World objects and sort by a given property
```
world_list, err = str.List(ctx,
    generated.WorldKindIdentity(),
    options.OrderBy("external.name"))

world_list, err = str.List(ctx,
    generated.WorldKindIdentity(),
    options.OrderBy("external.name"),
    options.OrderDescending())
```

## List the World objects and paginate the results
```
world_list, err = str.List(ctx,
    generated.WorldKindIdentity(),
    options.PageOffset(10),
    options.PageSize(50))
```
