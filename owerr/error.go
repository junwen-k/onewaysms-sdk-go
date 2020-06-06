// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package owerr

// Error OneWay specific error.
// Switch based on code to handle specific error when using OneWayClient.
type Error interface {
	// Satisfy the generic error interface.
	error

	// Returns the error details message.
	Message() string

	// Returns the short phrase depicting the classification of the error.
	Code() string

	// Returns the status code of the HTTP response.
	StatusCode() int
}

// New initializes a new OneWayError.
func New(code, message string, statusCode int) Error {
	return newBaseError(code, message, statusCode)
}

const (
	// RequestFailure request Failure error. Error is thrown when response status is not OK.
	RequestFailure = "RequestFailure"

	// InvalidCredentials invalid Credentials error. Error is thrown when API Username or API Password is invalid.
	InvalidCredentials = "InvalidCredentials"

	// InvalidSenderID invalid Sender ID error. Error is thrown when Sender ID is not valid.
	InvalidSenderID = "InvalidSenderID"

	// InvalidMobileNo invalid Mobile Number error. Error is thrown when Mobile No is not valid.
	InvalidMobileNo = "InvalidMobileNo"

	// InvalidLanguageType invalid LanguageType error. Error is thrown when Language type is not valid.
	InvalidLanguageType = "InvalidLanguageType"

	// InvalidMessageCharacters invalid Message Characters error. Error is thrown when there are invalid characters in the request message.
	InvalidMessageCharacters = "InvalidMessageCharacters"

	// InsufficientCreditBalance insufficientCreditBalance error. Error is thrown when the user does not have sufficient credit balance to perform certain tasks.
	InsufficientCreditBalance = "InsufficientCreditBalance"

	// MTInvalidNotFound mobile terminating ID invalid or not found error. Error is thrown when Mobile Terminating ID is not a valid ID or not found.
	MTInvalidNotFound = "MTInvalidNotFound"

	// MessageDeliveryFailure message delivery failure error. Error is thrown when Message has been failed to deliver when calling check transaction status API.
	MessageDeliveryFailure = "MessageDeliveryFailure"

	// UnknownError unknown error. Unknown Response returned from OneWay API Gateway.
	UnknownError = "UnknownError"
)
