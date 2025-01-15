package builder

import (
	"context"
	"fmt"
	"github.com/steebchen/prisma-client-go/engine"
	"github.com/steebchen/prisma-client-go/engine/protocol"
	"github.com/steebchen/prisma-client-go/logger"
	"time"
)

type QueryString struct {
	// Engine holds the implementation of how queries are processed
	Engine engine.Engine

	// Start time of the request for tracing
	Start time.Time

	// Query string
	Query string

	TxResult chan []byte
}

func NewQueryString() QueryString {
	return QueryString{
		Start: time.Now(),
	}
}

func (q QueryString) Exec(ctx context.Context, into interface{}) error {
	payload := protocol.GQLRequest{
		Query:     q.Query,
		Variables: map[string]interface{}{},
	}
	return q.Do(ctx, payload, into)
}

func (q QueryString) Do(ctx context.Context, payload interface{}, into interface{}) error {
	if q.Engine == nil {
		return fmt.Errorf("client.Prisma.Connect() needs to be called before sending queries")
	}

	err := q.Engine.Do(ctx, payload, into)
	now := time.Now()
	totalDuration := now.Sub(q.Start)
	logger.Debug.Printf("[timing] TOTAL %q", totalDuration)
	return err
}
