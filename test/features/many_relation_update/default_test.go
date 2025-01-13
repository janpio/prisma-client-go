package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
)

func TestFindManyRelationUpdateLink(t *testing.T) {
	test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()
		mockDBName := test.Start(t, test.SQLite, client.Engine, []string{})
		defer test.End(t, test.SQLite, client.Engine, mockDBName)

		_, err := client.Post.FindMany(
			Post.ID.Equals("non-existing"),
		).Update(
			Post.User.Link(User.ID.Equals("123")),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestFindManyRelationUpdateScalar(t *testing.T) {
	test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()
		mockDBName := test.Start(t, test.SQLite, client.Engine, []string{})
		defer test.End(t, test.SQLite, client.Engine, mockDBName)

		_, err := client.Post.FindMany(
			Post.ID.Equals("non-existing"),
		).Update(
			Post.UserID.Set("123"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}
