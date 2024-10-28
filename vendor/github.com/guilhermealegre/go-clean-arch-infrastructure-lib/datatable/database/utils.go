package database

import (
	"context"
	"net/url"
	"strconv"
)

func getIntParameter(key string, values url.Values, ctx context.Context) (intValue int) {
	if ctx != nil {
		if value := ctx.Value(key); value != nil {
			if v, ok := value.(int); ok {
				intValue = v
			}
		}
	}

	if values != nil && intValue == 0 {
		intValue, _ = strconv.Atoi(values.Get(key))
	}

	return intValue
}

func getStringParameter(key string, values url.Values, ctx context.Context) (strValue string) {
	if ctx != nil {
		if value := ctx.Value(key); value != nil {
			if v, ok := value.(string); ok {
				strValue = v
			}

			if v, ok := value.(OrderDirection); ok {
				strValue = string(v)
			}
		}
	}

	if values != nil && len(strValue) == 0 {
		strValue = values.Get(key)
	}

	return strValue
}
