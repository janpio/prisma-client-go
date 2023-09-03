# Raw API

You can use the raw API when there's something you can't do with the current go client features. The query will be redirected to the underlying database, so everything supported by the database should work. Please note that you need to use the syntax specific to the database you're using.

NOTE: When defining your return type structure, you have to use the native database type. For example, MySQL uses `int` for `bool`.
You can also use Prisma-specific raw data types, such as `RawInt`, `RawString`, so that it works without having to think about what is used internally. If you are querying for a specific model, you can also use `Raw<Model>Model`, e.g. `RawPostModel` instead of `PostModel`.

The examples use the following prisma schema:

```prisma
model Post {
    id        String   @id @default(cuid())
    createdAt DateTime @default(now())
    updatedAt DateTime @updatedAt
    published Boolean
    title     String
    content   String?
    views     Int      @default(0)

    comments Comment[]
}

model Comment {
    id        String   @id @default(cuid())
    createdAt DateTime @default(now())
    content   String

    post   Post   @relation(fields: [postID], references: [id])
    postID String
}
```

## MySQL & SQLite

### Query

Use `QueryRaw` to query for data and automatically unmarshal it into a slice of structs.

#### Select all

```go
var posts []db.RawPostModel
err := client.Prisma.QueryRaw("SELECT * FROM `Post`").Exec(ctx, &posts)
```

#### Select with parameters

```go
// note the usage of RawPostModel instead of PostModel
var posts []db.RawPostModel
err := client.Prisma.QueryRaw("SELECT * FROM `Post` WHERE id = ? AND title = ?", "123abc", "my post").Exec(ctx, &posts)
```

#### Custom Query

The Prisma client doesn't support aggregations out of the box. But you can do that via a custom query:

```go
var res []struct{
	PostID   db.RawString `json:"post_id"`
	Comments db.RawInt    `json:"n_comments"`
}
err := client.Prisma.QueryRaw("SELECT post_id, count(*) as n_comments FROM `Comment` GROUP BY post_id").Exec(ctx, &res)
```

Note that the query uses `db.RawString` etc in the struct definition to maintain compatibility. Note also that the results are an array of structs, not a struct.

#### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `Result{Count: int}`, which contains the affected rows.

```go
result, err := client.Prisma.ExecuteRaw("UPDATE `Post` SET title = ? WHERE id = ?", "my post", "123").Exec(ctx)
println(result.Count) // 1
```

## Postgres

### Query

Use `QueryRaw` to query for data and automatically unmarshal it into a slice of structs.s

#### Select all for a model

```go
var posts []db.RawPostModel
err := client.Prisma.QueryRaw(`SELECT * FROM "Post"`).Exec(ctx, &posts)
```

#### Select with parameters

```go
var posts []db.RawPostModel
err := client.Prisma.QueryRaw(`SELECT * FROM "Post" WHERE id = $1 AND title = $2`, "id2", "title2").Exec(ctx, &posts)
```

#### Custom Query

```go
var res []struct{
	ID        db.RawString  `json:"id"`
	Published db.RawBoolean `json:"published"`
}
err := client.Prisma.QueryRaw(`SELECT id, published FROM "Post"`).Exec(ctx, &res)
```

#### Using Prisma raw types

To ensure compatibility with database and go types, you can use raw types.

```go
// note the usage of db.RawString, db.RawInt, etc.
var res []struct{
	ID        db.RawString  `json:"post_id"`
	Comments  db.RawInt     `json:"comments"`
}
err := client.Prisma.QueryRaw(`SELECT post_id, count(*) as comments FROM "Comment" GROUP BY post_id`).Exec(ctx, &res)
```

#### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `Result{Count: int}`, which contains the affected rows.

```go
result, err := client.Prisma.ExecuteRaw(`UPDATE "Post" SET title = $1 WHERE id = $2`, "my post", "123").Exec(ctx)
println(result.Count) // 1
```
