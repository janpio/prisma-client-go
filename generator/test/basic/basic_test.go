package test

//go:generate prisma2 generate

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cx = context.Context
type Func func(t *testing.T, client Client, ctx cx)

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			return fmt.Errorf("command %s %s failed: %w", name, args, err)
		}
		if !exit.Success() {
			return fmt.Errorf("%s %s exited with status code %d and output %s: %w", name, args, exit.ExitCode(), string(out), err)
		}
	}
	return nil
}

func TestBasic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before string
		run    Func
	}{{
		name: "Nullability",
		// language=GraphQL
		before: `
			mutation {
				createOneUser(data: {
					id: "nullability",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					stuff: null,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.User.FindOne(User.Email.Equals("john@example.com")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			name, ok := actual.Name()
			assert.Equal(t, true, ok)
			assert.Equal(t, "John", name)

			stuff, ok := actual.Stuff()
			assert.Equal(t, false, ok)
			assert.Equal(t, "", stuff)
		},
	}, {
		name: "FindOne equals",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "findOne1",
					email: "john@findOne.com",
					username: "john_doe",
				}) {
					id
				}
				b: createOneUser(data: {
					id: "findOne2",
					email: "jane@findOne.com",
					username: "jane_doe",
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.User.FindOne(User.Email.Equals("jane@findOne.com")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, "findOne2", actual.ID)
		},
	}, {
		name: "FindMany equals",
		// language=GraphQL
		before: `
				mutation {
					a: createOneUser(data: {
						id: "findMany1",
						email: "1",
						username: "john",
						name: "a",
					}) {
						id
					}
					b: createOneUser(data: {
						id: "findMany2",
						email: "2",
						username: "john",
						name: "b",
					}) {
						id
					}
				}
			`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.User.FindMany(User.Username.Equals("john")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			a := "a"
			b := "b"
			assert.Equal(t, []UserModel{{
				user{
					ID:       "findMany1",
					Email:    "1",
					Username: "john",
					Name:     &a,
				},
			}, {
				user{
					ID:       "findMany2",
					Email:    "2",
					Username: "john",
					Name:     &b,
				},
			}}, actual)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd("rm", "-rf", "dev.sqlite"); err != nil {
				log.Fatal(err)
			}
			if err := cmd("rm", "-rf", "migrations"); err != nil {
				log.Fatal(err)
			}

			if err := cmd("prisma2", "lift", "save", "--create-db", "--name", "init"); err != nil {
				t.Fatalf("could not run lift save %s", err)
			}
			if err := cmd("prisma2", "lift", "up"); err != nil {
				t.Fatalf("could not run lift up %s", err)
			}

			client := NewClient()
			if err := client.Connect(); err != nil {
				t.Fatalf("could not connect %s", err)
				return
			}

			defer func() {
				err := client.Disconnect()
				if err != nil {
					// TODO blocked by prisma-engine panicking on disconnect
					// t.Fatalf("could not disconnect %s", err)
				}
			}()

			ctx := context.Background()

			if tt.before != "" {
				response, err := client.gql.Raw(ctx, tt.before, map[string]interface{}{})
				if err != nil {
					t.Fatalf("could not send mock query %s %+v", err, response)
				}
			}

			tt.run(t, client, ctx)
		})
	}
}
