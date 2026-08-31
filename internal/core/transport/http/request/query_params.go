package core_http_request

import (
	"fmt"
	"net/http"
	core_error "practice/internal/core/errors"
	"strconv"
	"time"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	v, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param=%s by key=%s not a valid int: %v: %w",
			param,
			key,
			err,
			core_error.ErrInvalidArgument,
		)
	}

	return &v, nil
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf("param=%s by key=%s not a valid date: %v : %w", param, key, err, core_error.ErrInvalidArgument)
	}

	return &date, nil
}
