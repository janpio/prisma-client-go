package raw

import (
	"context"
	"github.com/steebchen/prisma-client-go/runtime/builder"
	"github.com/steebchen/prisma-client-go/runtime/transaction"
)

func (r Raw) QueryString(query string) QueryStringExec {
	queryStr := builder.NewQueryString()
	queryStr.Engine = r.Engine
	queryStr.Query = query

	return QueryStringExec{
		query: queryStr,
	}
}

type QueryStringExec struct {
	query builder.QueryString
}

func (r QueryStringExec) ExtractQuery() builder.QueryString {
	return r.query
}

type TxQueryStringResult struct {
	query  builder.QueryString
	result *transaction.Result
}

func NewTxQueryStringResult() TxQueryStringResult {
	return TxQueryStringResult{
		result: &transaction.Result{},
	}
}

func (r TxQueryStringResult) ExtractQuery() builder.QueryString {
	return r.query
}

func (r TxQueryStringResult) IsTx() {}

func (r TxQueryStringResult) Into(v interface{}) error {
	return r.result.Get(r.query.TxResult, &v)
}

func (r QueryStringExec) Tx() TxQueryStringResult {
	v := NewTxQueryStringResult()
	v.query = r.ExtractQuery()
	v.query.TxResult = make(chan []byte, 1)
	return v
}

func (r QueryStringExec) Exec(ctx context.Context, into interface{}) error {
	return r.query.Exec(ctx, into)
}
