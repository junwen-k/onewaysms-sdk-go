// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package owsms

import (
	"github.com/pkg/errors"
)

// LanguageType SMS language type type.
type LanguageType string

const (
	// LanguageTypeNormal normal text message. (160 characters as 1 MT)
	LanguageTypeNormal LanguageType = "1"

	// LanguageTypeUnicode unicode text message. (70 characters as 1 MT)
	LanguageTypeUnicode LanguageType = "2"
)

// MTTransactionStatus mobile terminating transaction status type.
type MTTransactionStatus string

const (
	// MTTransactionStatusSuccess successfully sent message.
	MTTransactionStatusSuccess MTTransactionStatus = "success"

	// MTTransactionStatusTelcoDelivered message has been delivered to Telco.
	MTTransactionStatusTelcoDelivered MTTransactionStatus = "telco_delivered"
)

// SendSMSInput send SMS input structure.
type SendSMSInput struct {
	LanguageType LanguageType // Language Type of the SMS. Refer to LanguageType for details.
	Message      string       // Content of the SMS.
	MobileNo     []string     // Phone number of recipient. Phone number must include country code. For example: 6581234567.
}

// Validate validates send SMS input's values.
func (i *SendSMSInput) Validate() error {
	if i.Message == "" {
		return errors.New("SendSMSInput: Error: Message is required")
	}
	if i.MobileNo == nil || len(i.MobileNo) <= 0 {
		return errors.New("SendSMSInput: Error: MobileNo is required")
	}
	if i.LanguageType != LanguageTypeNormal && i.LanguageType != LanguageTypeUnicode {
		return errors.New("SendSMSInput: Error: LanguageType is invalid")
	}
	return nil
}

// SendSMSOutput send SMS output structure.
type SendSMSOutput struct {
	MTIDs []int // Mobile terminating ID(s) from the send SMS result.
}

// CheckTransactionStatusInput check transaction input structure.
type CheckTransactionStatusInput struct {
	MTID int // Mobile terminating ID returned from the send SMS result.
}

// Validate validates check transaction status input's values.
func (i *CheckTransactionStatusInput) Validate() error {
	if i.MTID == 0 {
		return errors.New("CheckTransactionStatusInput: Error: MTID is required")
	}
	return nil
}

// CheckTransactionStatusOutput check transaction output structure.
type CheckTransactionStatusOutput struct {
	Status MTTransactionStatus //Status of the mobile terminating transaction. Refer to MTTransactionStatus for more details.
}

// CheckCreditBalanceOutput check credit balance output structure.
type CheckCreditBalanceOutput struct {
	CreditBalance float32 // Remaining credit balance for the account of this client's config.
}
