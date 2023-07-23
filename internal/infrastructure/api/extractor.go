package api

import (
	"fmt"
	"net/url"
)

var missingParam = "missing query parameter: %s"

func extractParamStrFromQuery(query url.Values, queryName string) (string, error) {
	queryValue := query.Get(queryName)
	if queryValue == "" {
		return queryValue, fmt.Errorf(missingParam, queryName)
	}

	return queryValue, nil
}
