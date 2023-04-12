# Model GENerator
The model is a collection of YAML files describing `Store` Objects and Structures. 

**Modeled Objects** are used by the `Store` interface to manipulate data. Objects contain
- Metadata (**managed internally**)
    - Kind information
    - Primary key
    - Framework assigned identitier
    - Object manipulation timestamps (create, update...)
- External - any Structure, to be managed through external REST APIs (**optional**)
- Internal - to be managed by internal service code (React callbacks) (**optional**)

```
  - kind: Object
    name: World
    external: WorldExternalStruct
    internal: WorldInternalStruct
    primarykey: external.name
```

**Structures** are named collections of typed properties. Supported property types include
- Golang standard types
    - string
    - int
    - float
    - bool
- Other Structures (nesting)
- String-keyed maps
    - map[string]int (string, float...)
    - map[string]Struct
- Slices
    - []int (string, float...)
    - []Struct

```
  - kind: Struct
    name: WorldExternalStruct
    properties:
      - name: name
        type: string
      - name: description
        type: string
      - name: nested
        type: NestedWorldStruct
```


## Generated Package
Import the "generated" package to access Object interfaces and Schema.

Use <object>Factory() functions to create Model specific mutable Objects
which contain Object
- Metadata
- External / Internal
    - Property Getters/Setters
