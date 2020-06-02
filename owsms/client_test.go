package owsms_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/junwen-k/onewaysms-sdk-go/owerr"
	"github.com/junwen-k/onewaysms-sdk-go/owsms"
	"github.com/stretchr/testify/assert"
)

func TestSendSMS(t *testing.T) {
	var (
		svc    *owsms.Client
		output *owsms.SendSMSOutput
		err    error
	)

	t.Run("With valid values", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "145712468")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "Hello World",
			MobileNo: []string{"60123456789"},
		})
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, []int{145712468}, output.MTIDs)
	})

	t.Run("With multiple mobileNo", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "145712468,145712469")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "Hello World",
			MobileNo: []string{"60123456789", "60129876543"},
		})
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, []int{145712468, 145712469}, output.MTIDs)
	})

	t.Run("With request failure", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "request failure",
			MobileNo: []string{"60123456789"},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 500 (Internal Server Error): request failure")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "request failure", owErr.Message())
		assert.Equal(t, owerr.RequestFailure, owErr.Code())
		assert.Equal(t, http.StatusInternalServerError, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With invalid user credentials", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-100")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "invalid", "invalid", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "Hello World",
			MobileNo: []string{"60123456789"},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): apiusername or apipassword is invalid")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "apiusername or apipassword is invalid", owErr.Message())
		assert.Equal(t, owerr.InvalidCredentials, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With invalid senderID", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-200")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "invalid")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "Hello World",
			MobileNo: []string{"60123456789"},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): senderid parameter is invalid")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "senderid parameter is invalid", owErr.Message())
		assert.Equal(t, owerr.InvalidSenderID, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With invalid mobileNo", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-300")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "Hello World",
			MobileNo: []string{"invalid"},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): mobileno parameter is invalid")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "mobileno parameter is invalid", owErr.Message())
		assert.Equal(t, owerr.InvalidMobileNo, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With invalid languageType", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-400")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:      "Hello World",
			MobileNo:     []string{"60123456789"},
			LanguageType: "invalid",
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): languagetype is invalid")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "languagetype is invalid", owErr.Message())
		assert.Equal(t, owerr.InvalidLanguageType, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With invalid message", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-500")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "invalid",
			MobileNo: []string{"60123456789"},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): characters in message are invalid")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "characters in message are invalid", owErr.Message())
		assert.Equal(t, owerr.InvalidMessageCharacters, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With insufficient credit balance", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-600")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "insufficient credit balance",
			MobileNo: []string{"60123456789"},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): insufficient credit balance")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "insufficient credit balance", owErr.Message())
		assert.Equal(t, owerr.InsufficientCreditBalance, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With unknown error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "random")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.SendSMS(&owsms.SendSMSInput{
			Message:  "unknown error",
			MobileNo: []string{"60123456789"},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): unknown error")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "unknown error", owErr.Message())
		assert.Equal(t, owerr.UnknownError, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})
}

func TestCheckTransactionStatus(t *testing.T) {
	var (
		svc    *owsms.Client
		output *owsms.CheckTransactionStatusOutput
		err    error
	)

	t.Run("With success status", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "0")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.CheckTransactionStatus(&owsms.CheckTransactionStatusInput{
			MTID: 145712470,
		})
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, owsms.MTTransactionStatusSuccess, output.Status)
	})

	t.Run("With telco_delivered status", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "100")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.CheckTransactionStatus(&owsms.CheckTransactionStatusInput{
			MTID: 145712471,
		})
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, owsms.MTTransactionStatusTelcoDelivered, output.Status)
	})

	t.Run("With invalid mtID", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-100")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.CheckTransactionStatus(&owsms.CheckTransactionStatusInput{
			MTID: 1,
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): mtid is invalid or not found")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "mtid is invalid or not found", owErr.Message())
		assert.Equal(t, owerr.MTInvalidNotFound, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With failed mtID", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-200")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.CheckTransactionStatus(&owsms.CheckTransactionStatusInput{
			MTID: 1,
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): message delivery failed")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "message delivery failed", owErr.Message())
		assert.Equal(t, owerr.MessageDeliveryFailure, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With unknown error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "random")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.CheckTransactionStatus(&owsms.CheckTransactionStatusInput{
			MTID: 1,
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): unknown error")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "unknown error", owErr.Message())
		assert.Equal(t, owerr.UnknownError, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})
}

func TestCheckCreditBalance(t *testing.T) {
	var (
		svc    *owsms.Client
		output *owsms.CheckCreditBalanceOutput
		err    error
	)

	t.Run("With success status", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "6500.5")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "Username", "Password", "SenderID")

		output, _, err = svc.CheckCreditBalance()
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, float32(6500.5), output.CreditBalance)
	})

	t.Run("With invalid user credentials", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "-100")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "invalid", "invalid", "SenderID")

		output, _, err = svc.CheckCreditBalance()
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): apiusername or apipassword is invalid")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "apiusername or apipassword is invalid", owErr.Message())
		assert.Equal(t, owerr.InvalidCredentials, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})

	t.Run("With unknown error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "random")
		}))
		defer ts.Close()

		svc = owsms.NewClient(ts.URL, "invalid", "invalid", "SenderID")

		output, _, err = svc.CheckCreditBalance()
		assert.Error(t, err)
		assert.EqualError(t, err, "OneWaySMS: Error 200 (OK): unknown error")
		owErr, ok := err.(owerr.Error)
		assert.True(t, ok)
		assert.Equal(t, "unknown error", owErr.Message())
		assert.Equal(t, owerr.UnknownError, owErr.Code())
		assert.Equal(t, http.StatusOK, owErr.StatusCode())
		assert.Nil(t, output)
	})
}
