package mysql

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

const containerName = "go-client-mysql"

func setup(t *testing.T) {
	teardown(t)

	if err := hooks.Cmd("docker", "run", "--name", containerName, "-p", "3306:3306", "-e", "MYSQL_DATABASE=testing", "-e", "MYSQL_ROOT_PASSWORD=pw", "-d", "mysql:5.6"); err != nil {
		t.Fatal(err)
	}

	time.Sleep(15 * time.Second)
}

func teardown(t *testing.T) {
	if err := hooks.Cmd("docker", "stop", containerName); err != nil {
		log.Println(err)
	}

	if err := hooks.Cmd("docker", "rm", containerName, "--force"); err != nil {
		log.Println(err)
	}
}

func TestRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "raw query",
		// language=GraphQL
		before: []string{`
			mutation {
				a: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []UserModel
			err := client.Raw(`SELECT * FROM User`).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				RawUser: RawUser{
					ID:       "id1",
					Email:    "email1",
					Username: "a",
				},
			}, {
				RawUser: RawUser{
					ID:       "id2",
					Email:    "email2",
					Username: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with parameter",
		// language=GraphQL
		before: []string{`
			mutation {
				a: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []UserModel
			err := client.Raw(`SELECT * FROM User WHERE id = ?`, "id2").Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				RawUser: RawUser{
					ID:       "id2",
					Email:    "email2",
					Username: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with multiple parameters",
		// language=GraphQL
		before: []string{`
			mutation {
				a: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []UserModel
			err := client.Raw(`SELECT * FROM User WHERE id = ? AND email = ?`, "id2", "email2").Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				RawUser: RawUser{
					ID:       "id2",
					Email:    "email2",
					Username: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query count",
		// language=GraphQL
		before: []string{`
			mutation {
				a: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []struct {
				Count int `json:"count"`
			}
			err := client.Raw(`SELECT COUNT(*) AS count FROM User`).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 2, actual[0].Count)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup(t)

			client := NewClient()

			hooks.Start(t, client.Engine, tt.before)
			defer hooks.End(t, client.Engine)

			tt.run(t, client, context.Background())

			teardown(t)
		})
	}
}
