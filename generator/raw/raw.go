package raw

import (
	"github.com/prisma/prisma-client-go/generator/builder"
)

type Actions struct {
	Client builder.Client
}

func raw(client builder.Client, action string, query string, params ...interface{}) builder.Query {
	q := builder.NewQuery()
	q.Client = client
	q.Operation = "mutation"
	q.Method = action

	q.Inputs = append(q.Inputs, builder.Input{
		Name:  "query",
		Value: query,
	})

	// convert params to a string with an array
	var newParams = "["
	for i, param := range params {
		if i > 0 {
			newParams += ","
		}
		newParams += string(builder.Value(param))
	}
	newParams += "]"

	q.Inputs = append(q.Inputs, builder.Input{
		Name:  "parameters",
		Value: newParams,
	})

	return q
}
