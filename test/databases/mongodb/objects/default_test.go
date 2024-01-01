package raw

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestObjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "types",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "id1",
					Email:    "email1",
					Username: "username1",
					Info: InfoType{
						Age:    5,
						AgeOpt: 3,
					},
					InfoOpt: &InfoType{
						Age:    5,
						AgeOpt: 3,
					},
					List: []InfoType{
						{
							Age:    5,
							AgeOpt: 3,
						},
					},
				},
			}

			user, err := client.User.CreateOne(
				User.Email.Set("id1"),
				User.Username.Set("id1"),
				User.Info.Set(InfoType{
					Age:    5,
					AgeOpt: 3,
				}),
				User.InfoOpt.Set(InfoType{
					Age:    5,
					AgeOpt: 3,
				}),
				User.List.Set([]InfoType{{
					Age:    5,
					AgeOpt: 3,
				}}),
				User.ID.Set("id1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, user)

			user, err = client.User.FindUnique(
				User.ID.Equals("id1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, user)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			mockDB := test.Start(t, test.MongoDB, client.Engine, tt.before)
			defer test.End(t, test.MongoDB, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
