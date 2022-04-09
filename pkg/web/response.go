package web

import (
	"strconv"
)

type Response struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error error      `json:"error,omitempty"`
}

func NewResponse(code int, data interface{}, err error) Response {
	if code < 300 {
		return Response{strconv.FormatInt(int64(code), 10), data, nil}
	}
	return Response{strconv.FormatInt(int64(code), 10), nil, err}
}
