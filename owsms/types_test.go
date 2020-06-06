// Copyright (c) KwanJunWen
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package owsms_test

import (
	"testing"

	"github.com/junwen-k/onewaysms-sdk-go/owsms"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSendSMSInputValidate(t *testing.T) {
	tests := []struct {
		desc     string
		input    *owsms.SendSMSInput
		expected error
	}{
		{
			desc: "With valid values",
			input: &owsms.SendSMSInput{
				LanguageType: owsms.LanguageTypeNormal,
				Message:      "Hello World",
				MobileNo:     []string{"60123456789"},
			},
			expected: nil,
		},
		{
			desc: "With invalid LanguageType",
			input: &owsms.SendSMSInput{
				LanguageType: "invalid",
				Message:      "Hello World",
				MobileNo:     []string{"60123456789"},
			},
			expected: errors.New("SendSMSInput: Error: LanguageType is invalid"),
		},
		{
			desc: "With missing Message",
			input: &owsms.SendSMSInput{
				LanguageType: owsms.LanguageTypeNormal,
				Message:      "",
				MobileNo:     []string{"60123456789"},
			},
			expected: errors.New("SendSMSInput: Error: Message is required"),
		},
		{
			desc: "With missing MobileNo",
			input: &owsms.SendSMSInput{
				LanguageType: owsms.LanguageTypeNormal,
				Message:      "Hello World",
				MobileNo:     nil,
			},
			expected: errors.New("SendSMSInput: Error: MobileNo is required"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			actual := test.input.Validate()
			if actual != nil {
				assert.EqualError(t, test.expected, actual.Error())
			} else {
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestCheckTransactionStatusInputValidate(t *testing.T) {
	tests := []struct {
		desc     string
		input    *owsms.CheckTransactionStatusInput
		expected error
	}{
		{
			desc: "With valid values",
			input: &owsms.CheckTransactionStatusInput{
				MTID: 145712468,
			},
			expected: nil,
		},
		{
			desc: "With missing MTID",
			input: &owsms.CheckTransactionStatusInput{
				MTID: 0,
			},
			expected: errors.New("CheckTransactionStatusInput: Error: MTID is required"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			actual := test.input.Validate()
			if actual != nil {
				assert.EqualError(t, test.expected, actual.Error())
			} else {
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
