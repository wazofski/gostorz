# SQL Store
SQL package helps turn any SQL database into a Store

## SQLite
```
sqlite_store = store.New(
    generated.Schema(),
    sql.Factory(sql.SqliteConnection("test.sqlite")))
```

## mySQL
```
mysql_store = store.New(
    generated.Schema(),
    sql.Factory(sql.MySqlConnection(
        "user:pass@tcp(127.0.0.1:3306)/db"))
```
