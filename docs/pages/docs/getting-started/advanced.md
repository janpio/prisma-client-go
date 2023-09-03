## Advanced usage

In the [quickstart](quickstart.md), we have created a simple post model and ran a few queries. However, Prisma and the
Go client are designed to work with relations between models.

We already created a post model, such as for a blog. Let's assume we want to add comments to a post, and connect these
models in a way so we can rely on SQL's foreign keys and the Go client's ability to work with relations.

So let's introduce a new comment model:

```prisma
model Comment {
    id        String   @id @default(cuid())
    createdAt DateTime @default(now())
    content   String

    post   Post   @relation(fields: [postID], references: [id])
    postID String
}
```

We will also need to add a relation from to the post model in order to make a 1:n relation between those models:

```prisma
model Post {
    // ...

    // add this to your post model
    comments Comment[]
}
```

Your full schema should look like this:

```prisma
datasource db {
    // could be postgresql or mysql
    provider = "sqlite"
    url      = "file:dev.db"
}

generator db {
    provider = "go run github.com/steebchen/prisma-client-go"
}

model Post {
    id        String    @id @default(cuid())
    createdAt DateTime  @default(now())
    updatedAt DateTime  @updatedAt
    title     String
    published Boolean
    desc      String?

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

Whenever you make changes to your model, migrate your database and re-generate your prisma code:

```shell script
go run github.com/steebchen/prisma-client-go migrate dev --name add_comment_model
```

In order to create comments, we first need to create a post, and then reference that post when creating a comment.

```go
post, err := client.Post.CreateOne(
    db.Post.Title.Set("My new post"),
    db.Post.Published.Set(true),
    db.Post.Desc.Set("Hi there."),
    db.Post.ID.Set("123"),
).Exec(ctx)
if err != nil {
    return err
}

log.Printf("post: %+v", post)

// then create a comment
comments, err := client.Comment.CreateOne(
    db.Comment.Content.Set("my description"),
    // link the post we created before
    db.Comment.Post.Link(
        db.Post.ID.Equals(post.ID),
    ),
).Exec(ctx)
if err != nil {
    return err
}

log.Printf("post: %+v", comments)
```

Now that a post and a comment are created, you can query for them as follows:

```go
// return all published posts
posts, err := client.Post.FindMany(
    db.Post.Published.Equals(true),
).Exec(ctx)
if err != nil {
    return err
}

log.Printf("published posts: %+v", posts)

// insert a few new comments
_, err = client.Comment.CreateOne(
    db.Comment.Content.Set("first comment"),
    // link the post we created before
    db.Comment.Post.Link(
        db.Post.ID.Equals("123"),
    ),
).Exec(ctx)
if err != nil {
    return err
}
_, err = client.Comment.CreateOne(
    db.Comment.Content.Set("second comment"),
    // link the post we created before
    db.Comment.Post.Link(
        db.Post.ID.Equals("123"),
    ),
).Exec(ctx)
if err != nil {
    return err
}

// return all comments from a post with a given id
comments, err := client.Comment.FindMany(
    db.Comment.Post.Where(
        db.Post.ID.Equals("123"),
    ),
).Exec(ctx)
if err != nil {
    return err
}

log.Printf("comments of post with id 123: %+v", comments)

// return the first two comments from a post with which contains a given title, and sort by descending date
orderedComments, err := client.Comment.FindMany(
    db.Comment.Post.Where(
        db.Post.ID.Equals("123"),
    ),
).Take(2).OrderBy(
    db.Comment.CreatedAt.Order(db.SortOrderDesc),
).Exec(ctx)
if err != nil {
    return err
}

log.Printf("ordered comments: %+v", orderedComments)
```

Prisma also allows you to fetch multiple things at once. Instead of doing complicated joins, you can fetch a post and a
few of their comments in just a few lines and fully type-safe:

```go
// return a post by its id including 5 of its comments
post, err := client.Post.FindUnique(
    db.Post.ID.Equals("123"),
).With(
    // also fetch 3 this post's comments
    db.Post.Comments.Fetch().Take(3),
).Exec(ctx)

// will log post and its comments
log.Printf("post: %+v", post)
```

## API reference

To explore all query capabilities, check out the [API reference](../walkthrough).
