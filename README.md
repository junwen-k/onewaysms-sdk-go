# OneWaySMS SDK for Go

A non-official OneWaySMS SDK written for Go. Based on the latest [v1.3](http://smsd2.onewaysms.sg/api.pdf) OneWaySMS API documentation.

[![license](https://img.shields.io/github/license/junwen-k/onewaysms-sdk-go)](https://raw.githubusercontent.com/junwen-k/onewaysms-sdk-go/master/LICENSE.txt)

## Installing

Use `go get` to retrieve the SDK to add it to your `GOPATH` workspace, or project's Go module dependencies.

    go get github.com/junwen-k/onewaysms-sdk-go

To update the SDK use `go get -u` to retrieve the latest version of the SDK.

    go get -u github.com/junwen-k/onewaysms-sdk-go

## Dependencies

The SDK includes a `vendor` folder containing the runtime dependencies of the SDK. The metadata of the SDK's dependencies can be found in the Go module file `go.mod`.

## Usage and Getting Started

### Importing the SDK

Include the following import statement to use the SDK.

```go
package main

import (
  "github.com/junwen-k/onewaysms-sdk-go/owerr"
  "github.com/junwen-k/onewaysms-sdk-go/owsms"
)
```

### Initializing a new client

Initialize a new client using the `NewClient` function. For instance:

```go
func main() {
	svc := owsms.NewClient(
		"API_BASE_URL",
		"API_USERNAME",
		"API_PASSWORD",
		"SENDER_ID",
  )
  // ...
}
```

Optionally, a new client with a custom HTTP client can be initialized using the `NewClientWithHTTP` function. For instance:

```go
func main() {
	svc := owsms.NewClientWithHTTP(
		"API_BASE_URL",
		"API_USERNAME",
		"API_PASSWORD",
		"SENDER_ID",
		&http.Client{Timeout: time.Second * 30},
  )
  // ...
}
```

### Use case examples

1. **Send SMS** - Send SMS by calling OneWaySMS API gateway, returning mobile terminating ID(s) if request is successful.

   ```go
    func main() {
      // ...
      output, _, err := svc.SendSMS(&owsms.SendSMSInput{
        Message:  "Hello, 世界",
        MobileNo: []string{"60123456789", "60129876543"},
      })
      if err != nil {
        if owErr, ok := err.(owerr.Error); ok {
          switch owErr.Code() {
          case owerr.RequestFailure:
          // Handle RequestFailure
          case owerr.InvalidCredentials:
          // Handle InvalidCredentials
          case owerr.InvalidSenderID:
          // Handle InvalidSenderID
          case owerr.InvalidMobileNo:
          // Handle InvalidMobileNo
          case owerr.InvalidLanguageType:
          // Handle InvalidLanguageType
          case owerr.InvalidMessageCharacters:
          // Handle InvalidMessageCharacters
          case owerr.InsufficientCreditBalance:
          // Handle InsufficientCreditBalance
          case owerr.UnknownError:
            // Handle UnknownError
          default:
          }
        } else {
          // Handle Generic Error
        }
      }

      // MTIDs - Mobile terminating IDs
      fmt.Println(output.MTIDs)
    }
   ```

1. **Check MT Transaction Status** - Check mobile terminating transaction status based on mobile terminating ID provided. Mobile terminating ID can be obtained by calling send SMS API.

   ```go
    func main() {
      // ...
      output, _, err := svc.CheckTransactionStatus(&owsms.CheckTransactionStatusInput{
        MTID: 145712470,
      })
      if err != nil {
        if owErr, ok := err.(owerr.Error); ok {
          switch owErr.Code() {
          case owerr.MTInvalidNotFound:
          // Handle MTInvalidNotFound
          case owerr.MessageDeliveryFailure:
          // Handle MessageDeliveryFailure
          case owerr.UnknownError:
            // Handle UnknownError
          default:
          }
        } else {
          // Handle Generic Error
        }
      }

      switch output.Status {
      case owsms.MTTransactionStatusSuccess:
        // Handle success status
      case owsms.MTTransactionStatusTelcoDelivered:
        // Handle telco delivered status
      default:
      }
    }
   ```

1. **Check Credit Balance**. Check remaining credit balance for the account in the client's config.

   ```go
    func main() {
      // ...
      output, _, err := svc.CheckCreditBalance()
      if err != nil {
        if owErr, ok := err.(owerr.Error); ok {
          switch owErr.Code() {
          case owerr.InvalidCredentials:
          // Handle InvalidCredentials
          case owerr.UnknownError:
            // Handle UnknownError
          default:
          }
        } else {
          // Handle Generic Error
        }
      }

      // CreditBalance - Remaining credit balance for this account
      fmt.Println(output.CreditBalance)
    }
   ```

## License

This SDK is distributed under the MIT License, see LICENSE.txt for more information.
