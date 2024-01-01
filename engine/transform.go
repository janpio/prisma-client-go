package engine

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/steebchen/prisma-client-go/logger"
)

// transformResponse for raw queries
// transforms all custom prisma types into native go types, such as
// [{"prisma__type":"string","prisma__value":"asdf"},{"prisma__type":"null","prisma__value":null}]
// ->
// ["asdf", null]
func transformResponse(data []byte) ([]byte, error) {
	logger.Debug.Printf("before transform: %s", data)
	var m interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	forEachValue(&m, func(k *string, i *int, v *interface{}) (interface{}, bool) {
		if v == nil {
			return nil, false
		}
		var n = *v
		o, isObject := (*v).(map[string]interface{})
		if isObject {
			var ok bool
			n, ok = handleObject(o)
			if !ok {
				return n, false
			}
		}
		return n, true
	})

	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("transform response marshal: %w", err)
	}

	logger.Debug.Printf("transformed response: %s", out)

	return out, nil
}

func handleObject(o map[string]interface{}) (interface{}, bool) {
	if t, ok := o["prisma__type"]; ok {
		if t == "bytes" {
			// bytes from prisma are base64 encoded
			bytes, ok := o["prisma__value"].(string)
			if !ok {
				panic("expected bytes")
			}
			dst := make([]byte, base64.StdEncoding.DecodedLen(len(bytes)))
			n, err := base64.StdEncoding.Decode(dst, []byte(bytes))
			if err != nil {
				panic(err)
			}
			dst = dst[:n]
			return dst, false
		}
		if t == "array" {
			value, ok := o["prisma__value"].([]interface{})
			if !ok {
				panic("expected array")
			}
			var items []interface{}
			for _, item := range value {
				item, _ := handleObject(item.(map[string]interface{}))
				items = append(items, item)
			}
			return items, false
		}
		return o["prisma__value"], false
	}
	return o, true
}

func forEachValue(obj *interface{}, handler func(*string, *int, *interface{}) (interface{}, bool)) {
	if obj == nil {
		return
	}
	var ok bool
	var n = *obj
	// Yield all key/value pairs for objects.
	o, isObject := (*obj).(map[string]interface{})
	if isObject {
		for k := range o {
			item := o[k]
			o[k], ok = handler(&k, nil, &item)
			item = o[k]
			if ok {
				forEachValue(&item, handler)
			}
		}
		n = o
	}
	// Yield each index/value for arrays.
	a, isArray := (*obj).([]interface{})
	if isArray {
		for i := range a {
			item := a[i]
			a[i], ok = handler(nil, &i, &item)
			item = a[i]
			if ok {
				forEachValue(&item, handler)
			}
		}
		n = a
	}
	*obj = n
	// Do nothing for primitives since the handler got them.
}
