package introspection

//type IntrospectRequest struct {
//	Id      int                     `json:"id"`
//	Jsonrpc string                  `json:"jsonrpc"`
//	Method  string                  `json:"method"`
//	Params  IntrospectRequestParams `json:"params"`
//}

type IntrospectRequest struct {
	Id      int                      `json:"id"`
	Jsonrpc string                   `json:"jsonrpc"`
	Method  string                   `json:"method"`
	Params  []map[string]interface{} `json:"params"`
}

type IntrospectRequestParams struct {
	CompositeTypeDepth int64  `json:"compositeTypeDepth"`
	Schema             string `json:"schema"`
}

type IntrospectResponse struct {
	Jsonrpc string                    `json:"jsonrpc"`
	Result  *IntrospectResponseResult `json:"result,omitempty"`
	Error   *IntrospectResponseError  `json:"error,omitempty"`
}

type IntrospectResponseResult struct {
	ExecutedSteps int         `json:"executedSteps"`
	DataModel     string      `json:"dataModel"`
	Marnings      interface{} `json:"marnings"`
	Version       string      `json:"version"`
}

type IntrospectResponseError struct {
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Data    IntrospectResponseErrorData `json:"data"`
}

type IntrospectResponseErrorData struct {
	IsPanic bool                            `json:"is_panic"`
	Message string                          `json:"message"`
	Meta    IntrospectResponseErrorDataMeta `json:"meta"`
}

type IntrospectResponseErrorDataMeta struct {
	FullError string `json:"full_error"`
}
