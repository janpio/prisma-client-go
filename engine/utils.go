package engine

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
)

func InlineQueryDocument(query *ast.QueryDocument, variable map[string]interface{}) (string, error) {
	// 这里要循环处理，去除变量输入
	for _, operation := range query.Operations {
		operation.VariableDefinitions = ast.VariableDefinitionList{}
	}

	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(query)

	bufstr := buf.String()

	for k, v := range variable {
		// s, _ := json.MarshalIndent(v, "", "\t") // TODO:这里去掉引号
		s, _ := json.Marshal(v) // TODO:这里去掉引号
		ss := convert(string(s))
		bufstr = strings.ReplaceAll(bufstr, "$"+k, ss)
	}
	return bufstr, nil
}

func InlineQuery(str string, variable map[string]interface{}) (string, error) {
	query, err := parser.ParseQuery(&ast.Source{Input: str})
	if err != nil {
		gqlErr := err.(*gqlerror.Error)
		return "", gqlerror.List{gqlErr}
	}

	return InlineQueryDocument(query, variable)
}

// https://www.cnblogs.com/vicF/p/9517960.html
// {"id":{"equals":"ssss"}}=>{id:{equals:"ssss"}}
// ["id","id"]=>["id","id"]
func convert(s string) string {
	reg := regexp.MustCompile("\"(\\w+)\"(\\s*:\\s*)")
	res := reg.ReplaceAllString(s, "$1$2")

	return res
}

// func convert2(s string) string {
// 	var b bytes.Buffer
// 	shouldSkip := true

// 	for i := 0; i < len(s); i++ {
// 		c := string(s[i])
// 		if c == `"` {
// 			if i > 0 && string(s[i-1]) == `:` {
// 				shouldSkip = false
// 				b.WriteString(c)
// 				continue
// 			}
// 			if shouldSkip {
// 				continue
// 			}
// 			shouldSkip = true
// 		}

// 		b.WriteString(c)
// 	}

// 	return b.String()
// }
