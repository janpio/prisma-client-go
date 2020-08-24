package raw

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

type RawUserModel struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Name     *string `json:"name"`
	Stuff    *string `json:"stuff"`
	Str      string  `json:"str"`
	StrOpt   *string `json:"strOpt"`
	Int      int     `json:"int"`
	IntOpt   *int    `json:"intOpt"`
	Float    string  `json:"float"`
	FloatOpt *string `json:"floatOpt"`
	Bool     bool    `json:"bool"`
	BoolOpt  *bool   `json:"boolOpt"`
}

func TestRaw(t *testing.T) {
	t.Parallel()

	strOpt := "strOpt"
	i := 5
	f := "5.5"
	b := false

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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			err := client.QueryRaw(`SELECT * FROM "User"`).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			strOpt := "strOpt"
			i := 5
			f := "5.5"
			b := false
			expected := []RawUserModel{{
				ID:       "id1",
				Email:    "email1",
				Username: "a",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     true,
				BoolOpt:  &b,
			}, {
				ID:       "id2",
				Email:    "email2",
				Username: "b",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     true,
				BoolOpt:  &b,
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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			err := client.QueryRaw(`SELECT * FROM "User" WHERE id = $1`, "id2").Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:       "id2",
				Email:    "email2",
				Username: "b",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     true,
				BoolOpt:  &b,
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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			err := client.QueryRaw(`SELECT * FROM "User" WHERE id = $1 AND email = $2`, "id2", "email2").Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:       "id2",
				Email:    "email2",
				Username: "b",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     true,
				BoolOpt:  &b,
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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
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
					str: "str",
					strOpt: "strOpt",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []struct {
				Count int `json:"count"`
			}
			err := client.QueryRaw(`SELECT COUNT(*) AS count FROM "User"`).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 2, actual[0].Count)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			mockDB := test.Start(t, test.PostgreSQL, client.Engine, tt.before)
			defer test.End(t, test.PostgreSQL, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
