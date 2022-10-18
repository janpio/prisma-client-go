package migrate

type MigrationRequest struct {
	Id      int                    `json:"id"`
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  MigrationRequestParams `json:"params"`
}

type MigrationRequestParams struct {
	Force  bool   `json:"force"`
	Schema string `json:"schema"`
}

type MigrationResponse struct {
	Jsonrpc string                   `json:"jsonrpc"`
	Result  *MigrationResponseResult `json:"result,omitempty"`
	Error   *MigrationResponseError  `json:"error,omitempty"`
}

type MigrationResponseResult struct {
	ExecutedSteps int `json:"executedSteps"`
}

type MigrationResponseError struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    MigrationResponseErrorData `json:"data"`
}

type MigrationResponseErrorData struct {
	IsPanic bool                           `json:"is_panic"`
	Message string                         `json:"message"`
	Meta    MigrationResponseErrorDataMeta `json:"meta"`
}

type MigrationResponseErrorDataMeta struct {
	FullError string `json:"full_error"`
}
