package core_http_request

import (
	"fmt"
	"net/http"
	core_error "practice/internal/core/errors"
	"strconv"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("no key='%s' in path values :%w", core_error.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("fail to parse pathValue: %v: %w", err, core_error.ErrInvalidArgument)
	}

	return val, nil
}
