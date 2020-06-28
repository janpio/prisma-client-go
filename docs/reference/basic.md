# Basic API

Find, update and delete records.

## Reading data

### Find many records

```go
users, err := client.User.FindMany(
    db.User.Name.Equals("hi"),
).Exec(ctx)
```

If no records are found, this returns an empty array without returning an error (like usual SQL queries).

### Find one record

```go
user, err := client.User.FindOne(
    db.User.ID.Equals("123"),
).Exec(ctx)

if err == db.ErrNotFound {
    log.Printf("no record with id 123")
}
```

This returns an error of type `ErrNotFound` (exported in the `db` package) if there was no such record.

### Query API

Depending on the data types of your fields, you will automatically be able to query for respective operations. For example, for integer or float fields you might want to query for a field which is less than or greater than some number.

```go
user, err := client.User.FindOne(
    // query for names containing the string "Jo"
    db.User.Name.Contains("Jo"),
).Exec(ctx)
```

Other possible queries are:

```go
// query for people who are named "John"
db.User.Name.Contains("John"),
// query for names containing the string "oh"
db.User.Name.Contains("oh"),
// query for names starting with "Jo"
db.User.Name.HasPrefix("Jo"),
// query for names ending with "Jo"
db.User.Name.HasSuffix("hn"),
// query for all users which are younger than or exactly 18
db.User.Age.LTE(18),
// query for all users which are younger than 18
db.User.Age.LT(18),
// query for all users which are older than or exactly 18
db.User.Age.GT(18),
// query for all users which are older than 18
db.User.Age.GTE(18),
// query for all users which were created in the last 6 hours
db.User.CreatedAt.After(time.Now().Add(-6 * time.Hour)),
// query for all users which were created until yesterday
db.User.CreatedAt.Before(time.Now().Truncate(24 * time.Hour)),
```

All of these queries are fully type-safe and independent of the underlying database.

### Querying for relations

In a query, you can query for relations by using "Some" or "Every". You can also query for deeply nested relations.

```go
// get a user which has at least one post with a title "My Title" and that post's comments are all "What up?"
actual, err := client.User.FindMany(
    User.Email.Equals("john@example.com"),
    User.Posts.Some(
        Post.Title.Equals("My Title"),
        Post.Comments.Every(
            Comment.Content.Contains("What up?"),
        ),
    ),
).Exec(ctx)
```

## Writing data

### Create a record

```go
created, err := client.User.CreateOne(
    // required fields
    User.Email.Set("email"),
    User.Username.Set("username"),

    // optional fields
    User.ID.Set("id"),
    User.Name.Set("name"),
    User.Stuff.Set("stuff"),
).Exec(ctx)
```

### Create a record with a relation

Use the method `Link` to connect new objects with existing ones. For example, the following query creates a new post and sets the author of the post to a user with a given ID.

```go
created, err := client.Post.CreateOne(
    Post.Title.Set(title),
    Post.Author.Link(
        User.ID.Equals(userID),
    ),
    Post.ID.Set("post"),
).Exec(ctx)
```

### Update a record

To update a record, just query for a field using FindOne or FindMany, and then just chain it by invoking `.Update()`.

```go
updated, err := client.User.FindOne(
    User.Email.Equals("john@example.com"),
).Update(
    User.Username.Set("new-username"),
    User.Name.Set("New Name"),
).Exec(ctx)
```

### Delete a record

To delete a record, just query for a field using FindOne or FindMany, and then just chain it by invoking `.Delete()`.

```go
updated, err := client.User.FindOne(
    User.Email.Equals("john@example.com"),
).Delete().Exec(ctx)
```
