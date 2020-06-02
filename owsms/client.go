package owsms

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/junwen-k/onewaysms-sdk-go/owerr"
)

const version = "0.1.0"

// doer implements http.Client Do interface.
type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client OneWaySMS client structure.
// Based on specifications found in http://smsd2.onewaysms.sg/api.pdf.
type Client struct {
	client      doer
	baseURL     string
	apiUsername string
	apiPassword string
	senderID    string
}

// NewClient initializes a new OneWaySMS client.
func NewClient(baseURL, apiUsername, apiPassword, senderID string) *Client {
	return &Client{
		client:      http.DefaultClient,
		baseURL:     baseURL,
		apiUsername: apiUsername,
		apiPassword: apiPassword,
		senderID:    senderID,
	}
}

// NewClientWithHTTP initializes a new OneWaySMS client with custom http client.
func NewClientWithHTTP(baseURL, apiUsername, apiPassword, senderID string, client doer) *Client {
	c := NewClient(baseURL, apiUsername, apiPassword, senderID)
	c.client = client
	return c
}

func (c *Client) messageToHex(message string) string {
	buf := new(bytes.Buffer)
	for _, r := range message {
		buf.WriteString(strings.Trim(fmt.Sprintf("%U", r), "U+"))
	}
	return buf.String()
}

func (c *Client) getLanguageType(message string) LanguageType {
	m := message
	for len(m) > 0 {
		_, size := utf8.DecodeRuneInString(m)
		if size > 1 {
			return LanguageTypeUnicode
		}
		m = m[size:]
	}
	return LanguageTypeNormal
}

func (c *Client) buildRequestURL(path string, urlParams map[string]string) string {
	params := url.Values{}
	for k, v := range urlParams {
		params.Add(k, v)
	}
	return fmt.Sprintf("%s/%s?%s", c.baseURL, path, params.Encode())
}

func (c *Client) buildSendSMSRequestURL(input *SendSMSInput) string {
	if input.LanguageType == "" {
		input.LanguageType = c.getLanguageType(input.Message)
	}
	if input.LanguageType == LanguageTypeUnicode {
		input.Message = c.messageToHex(input.Message)
	}

	return c.buildRequestURL("api.aspx", map[string]string{
		"apiusername":  c.apiUsername,
		"apipassword":  c.apiPassword,
		"senderid":     c.senderID,
		"mobileno":     strings.Join(input.MobileNo, ","),
		"languagetype": string(input.LanguageType),
		"message":      input.Message,
	})
}

func (c *Client) buildCheckTransactionStatusRequestURL(input *CheckTransactionStatusInput) string {
	return c.buildRequestURL("bulktrx.aspx", map[string]string{
		"mtid": strconv.Itoa(input.MTID),
	})
}

func (c *Client) buildCheckCreditBalanceRequestURL() string {
	return c.buildRequestURL("bulkcredit.aspx", map[string]string{
		"apiusername": c.apiUsername,
		"apipassword": c.apiPassword,
	})
}

func (c *Client) getRequest(requestURL string) (*http.Response, error) {
	if c.client == nil {
		c.client = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("onewaysms-sdk-go/%s", version))

	return c.client.Do(req)
}

// SendSMS Initiate send SMS request. SMS's language type will be automatically set unless it is defined in the SMS request structure.
func (c *Client) SendSMS(input *SendSMSInput) (*SendSMSOutput, *http.Response, error) {
	requestURL := c.buildSendSMSRequestURL(input)

	resp, err := c.getRequest(requestURL)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, owerr.New(owerr.RequestFailure, "request failure", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	mtIDs := make([]int, 0)
	for _, _mtID := range strings.Split(string(b), ",") {
		mtID, err := strconv.Atoi(strings.TrimSpace(_mtID))
		if err != nil {
			return nil, resp, owerr.New(owerr.UnknownError, "unknown error", resp.StatusCode)
		}
		mtIDs = append(mtIDs, mtID)
	}

	if len(mtIDs) <= 0 {
		return nil, resp, owerr.New(owerr.UnknownError, "unknown error", resp.StatusCode)
	}

	if mtIDs[0] > 0 {
		return &SendSMSOutput{MTIDs: mtIDs}, resp, nil
	}
	switch mtIDs[0] {
	case -100:
		return nil, resp, owerr.New(owerr.InvalidCredentials, "apiusername or apipassword is invalid", resp.StatusCode)
	case -200:
		return nil, resp, owerr.New(owerr.InvalidSenderID, "senderid parameter is invalid", resp.StatusCode)
	case -300:
		return nil, resp, owerr.New(owerr.InvalidMobileNo, "mobileno parameter is invalid", resp.StatusCode)
	case -400:
		return nil, resp, owerr.New(owerr.InvalidLanguageType, "languagetype is invalid", resp.StatusCode)
	case -500:
		return nil, resp, owerr.New(owerr.InvalidMessageCharacters, "characters in message are invalid", resp.StatusCode)
	case -600:
		return nil, resp, owerr.New(owerr.InsufficientCreditBalance, "insufficient credit balance", resp.StatusCode)
	default:
		return nil, resp, owerr.New(owerr.UnknownError, "unknown error", resp.StatusCode)
	}
}

// CheckTransactionStatus check transaction status based on mobile terminating ID provided.
func (c *Client) CheckTransactionStatus(input *CheckTransactionStatusInput) (*CheckTransactionStatusOutput, *http.Response, error) {
	requestURL := c.buildCheckTransactionStatusRequestURL(input)

	resp, err := c.getRequest(requestURL)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	status, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		return nil, resp, owerr.New(owerr.UnknownError, "unknown error", resp.StatusCode)
	}

	switch status {
	case 0:
		return &CheckTransactionStatusOutput{Status: MTTransactionStatusSuccess}, resp, nil
	case 100:
		return &CheckTransactionStatusOutput{Status: MTTransactionStatusTelcoDelivered}, resp, nil
	case -100:
		return nil, resp, owerr.New(owerr.MTInvalidNotFound, "mtid is invalid or not found", resp.StatusCode)
	case -200:
		return nil, resp, owerr.New(owerr.MessageDeliveryFailure, "message delivery failed", resp.StatusCode)
	default:
		return nil, resp, owerr.New(owerr.UnknownError, "unknown error", resp.StatusCode)
	}
}

// CheckCreditBalance check remaining credit balance based on API Username and Password from client's config.
func (c *Client) CheckCreditBalance() (*CheckCreditBalanceOutput, *http.Response, error) {
	requestURL := c.buildCheckCreditBalanceRequestURL()

	resp, err := c.getRequest(requestURL)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	creditBalance, err := strconv.ParseFloat(strings.TrimSpace(string(b)), 32)
	if err != nil {
		return nil, resp, owerr.New(owerr.UnknownError, "unknown error", resp.StatusCode)
	}

	if creditBalance >= 0 {
		return &CheckCreditBalanceOutput{CreditBalance: float32(creditBalance)}, resp, nil
	}

	switch creditBalance {
	case -100:
		return nil, resp, owerr.New(owerr.InvalidCredentials, "apiusername or apipassword is invalid", resp.StatusCode)
	default:
		return nil, resp, owerr.New(owerr.UnknownError, "unknown error", resp.StatusCode)
	}
}
