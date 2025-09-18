// Package dto
package dto

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Success[T any](data T) *Response[T] {
	return &Response[T]{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func SuccessWithString[T any](msg string, data T) *Response[T] {
	return &Response[T]{
		Code: 200,
		Msg:  msg,
		Data: data,
	}
}

func Error(code int, msg string) *ErrorResponse {
	return &ErrorResponse{
		Code: code,
		Msg:  msg,
	}
}
