// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package owerr

import (
	"fmt"
	"net/http"
)

// OneWay base error.
type baseError struct {
	code       string
	message    string
	statusCode int
}

func newBaseError(code, message string, statusCode int) *baseError {
	return &baseError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

// Error returns the string representation of the error.
func (e *baseError) Error() string {
	return fmt.Sprintf("OneWaySMS: Error %d (%s): %s", e.statusCode, http.StatusText(e.statusCode), e.message)
}

// Message returns OneWay error message.
func (e *baseError) Message() string {
	return e.message
}

// Code returns OneWay error code.
func (e *baseError) Code() string {
	return e.code
}

// StatusCode returns OneWay error status code.
func (e *baseError) StatusCode() int {
	return e.statusCode
}
