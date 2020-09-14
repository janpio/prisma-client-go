# Raw API

You can use the raw API when there's something you can't do with the current go client features. The query will be redirected to the underlying database, so everything supported by the database should work. Please note that you need to use the syntax specific to the database you're using.

## MySQL & SQLite

### Select all

```go
var users []db.UserModel
err := client.QueryRaw(`SELECT * FROM User`).Exec(ctx, &users)
```

### Select specific

```go
var users []UserModel
err := client.QueryRaw(`SELECT * FROM User WHERE id = ? AND email = ?`, "123abc", "prisma@example.com").Exec(ctx, &users)
```

### Count

```go
count, err := client.ExecuteRaw(`SELECT COUNT(*) AS count FROM User`).Exec(ctx, &actual)
```

## Postgres

### Select all

```go
var users []UserModel
err := client.QueryRaw(`SELECT * FROM "User"`).Exec(ctx, &users)
```

### Select specific

```go
var users []UserModel
err := client.QueryRaw(`SELECT * FROM "User" WHERE id = $1 AND email = $2`, "id2", "email2").Exec(ctx, &users)
```

### Count

```go
count, err := client.ExecuteRaw(`SELECT COUNT(*) AS count FROM "User"`).Exec(ctx, &result)
```
