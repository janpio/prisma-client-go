package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prisma/prisma-client-go/logger"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
	"strconv"
	"time"
)

var prismaQueryEngineMap = map[int64]*QueryEngine{}

type EngineFactory interface {
	GetPrismaQueryEngine() (*QueryEngine, error)
	ReloadPrismaQueryEngine() error
	QuerySchema(param GQLRequest, result interface{}) error
}

// QueryEngineFactory 创建工厂结构体并实现工厂接口
type QueryEngineFactory struct {
	Key      int64  `json:"key"`
	DBSchema string `json:"DBSchema"`
}

func NewQueryEngineFactory(key int64, dbSchema string) QueryEngineFactory {
	return QueryEngineFactory{
		Key:      key,
		DBSchema: dbSchema,
	}
}

func (q *QueryEngineFactory) GetPrismaQueryEngine() (*QueryEngine, error) {
	// 如果不存在
	if _, ok := prismaQueryEngineMap[q.Key]; !ok {
		// 创建
		content, err := Pull(q.DBSchema)
		if err != nil {
			logger.Debug.Printf("connect fail err : ", err)
			return nil, err
		}
		queryEngine := NewQueryEngine(content, false)
		if err := queryEngine.ConnectSDK(); err != nil {
			logger.Debug.Printf("connect fail err : ", err)
			return nil, err
		}
		prismaQueryEngineMap[q.Key] = queryEngine
	}
	return prismaQueryEngineMap[q.Key], nil
}

func (q *QueryEngineFactory) ReloadPrismaQueryEngine() error {
	// 先销毁旧的引擎
	if _, ok := prismaQueryEngineMap[q.Key]; ok {
		prismaQueryEngineMap[q.Key].Disconnect()
	}

	// 创建新引擎
	content, err := Pull(q.DBSchema)
	if err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return err
	}
	queryEngine := NewQueryEngine(content, false)
	if err := queryEngine.ConnectSDK(); err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return err
	}
	// 存入
	prismaQueryEngineMap[q.Key] = queryEngine
	return nil
}

func DisConnectEngine() {
	for key, engine := range prismaQueryEngineMap {
		engine.Disconnect()
		delete(prismaQueryEngineMap, key)
	}
}

func (q *QueryEngineFactory) QuerySchema(param GQLRequest, result interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Debug.Println(err)
		}
	}()
	ctx := context.TODO()
	queryEngine, err := q.GetPrismaQueryEngine()
	if err != nil {
		return err
	}
	//err = queryEngine.DoQuery(ctx, param, result)
	err = queryEngine.DoManyQuery(ctx, param, result)
	if err != nil {
		return err
	}
	return nil
}

type GQLResult struct {
	Data       json.RawMessage        `json:"data"`
	Errors     []GQLError             `json:"errors"`
	Extensions map[string]interface{} `json:"extensions"`
}

type ErrResponse struct {
	Errors []SQLErrResult `json:"errors"`
}

type SQLErrResult struct {
	Message   string        `json:"message"`
	Path      []string      `json:"path"`
	Locations []interface{} `json:"locations"`
}

func NewSqlErrResult(errStr string) SQLErrResult {
	return SQLErrResult{
		Message:   errStr,
		Path:      make([]string, 0),
		Locations: make([]interface{}, 0),
	}
}
func (e *QueryEngine) DoManyQuery(ctx context.Context, payload GQLRequest, v interface{}) error {
	queryObj, _ := parser.ParseQuery(&ast.Source{Input: payload.Query})
	resultErr := ErrResponse{
		Errors: make([]SQLErrResult, 0),
	}
	if len(queryObj.Operations) != 1 {
		return fmt.Errorf("一次只能查询一个operation")
	}
	ope := queryObj.Operations[0]

	// 如果只有一个查询，则直接查询返回
	if len(ope.SelectionSet) == 1 {
		onePayLoad := GQLRequest{
			Query:     payload.Query,
			Variables: map[string]interface{}{},
		}
		return e.DoQuery(ctx, onePayLoad, v)
	}

	// 多个查询
	requests := make([]GQLRequest, len(ope.SelectionSet))

	selectionset := ope.SelectionSet
	for i, selection := range selectionset {
		ope.SelectionSet = ast.SelectionSet{selection}
		requests[i] = GQLRequest{
			Query:     FormatOperateionDocument(ope),
			Variables: map[string]interface{}{},
		}
	}
	type GQLBatchResult struct {
		Errors []GQLError  `json:"errors"`
		Result []GQLResult `json:"batchResult"`
	}
	var result GQLBatchResult
	payloads := GQLBatchRequest{
		Batch:       requests,
		Transaction: true,
	}

	if err := e.BatchReq(ctx, payloads, &result); err != nil {
		// 如果出现错误，则将错误返回给前端
		resultErr.Errors = append(resultErr.Errors, NewSqlErrResult(err.Error()))
		errBytes, _ := json.Marshal(resultErr)
		if err := json.Unmarshal(errBytes, v); err != nil {
			return fmt.Errorf("json unmarshal: %w", err)
		}
		return nil
	}
	if len(result.Errors) > 0 {
		// 如果出现错误，则将错误返回给前端
		resultErr.Errors = append(resultErr.Errors, NewSqlErrResult(result.Errors[0].RawMessage()))
		errBytes, _ := json.Marshal(resultErr)
		if err := json.Unmarshal(errBytes, v); err != nil {
			return fmt.Errorf("json unmarshal: %w", err)
		}
		return nil
	}
	// 合并JSON字符串
	var tmpRes string
	for idx, inner := range result.Result {
		if len(inner.Errors) > 0 {
			// 如果出现错误，则直接将错误返回给前端
			resultErr.Errors = append(resultErr.Errors, NewSqlErrResult(result.Errors[0].RawMessage()))
			errBytes, _ := json.Marshal(resultErr)
			if err := json.Unmarshal(errBytes, v); err != nil {
				return fmt.Errorf("json unmarshal: %w", err)
			}
			return nil
		}

		str := string(inner.Data)
		// 最后一条
		if idx == len(result.Result)-1 {
			tmpRes = tmpRes + str[1:] // 删除开头的{
		} else {
			// 非最后一条
			tmpRes = tmpRes + str[:len(str)-1] + "," // 删除结尾的}
		}
	}
	resultStruct := struct {
		Data interface{} `json:"data"`
	}{}
	if err := json.Unmarshal([]byte(tmpRes), &resultStruct.Data); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	resultBytes, _ := json.Marshal(resultStruct)
	if err := json.Unmarshal(resultBytes, v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}
	return nil
}

func (e *QueryEngine) DoQuery(ctx context.Context, payload GQLRequest, v interface{}) error {
	startReq := time.Now()

	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	logger.Debug.Printf("[timing] query engine request took %s", time.Since(startReq))

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}
	return nil
}

func InterfaceToString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}
}

// 位置挪走
func FormatOperateionDocument(operate *ast.OperationDefinition) string {

	query := &ast.QueryDocument{
		Operations: ast.OperationList{operate},
	}
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)

	bufstr := buf.String()

	return bufstr
}

// Do sends the http Request to the query engine and unmarshals the response
func (e *QueryEngine) BatchReq(ctx context.Context, payload interface{}, v interface{}) error {
	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}
