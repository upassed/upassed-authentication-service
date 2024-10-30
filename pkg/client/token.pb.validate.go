// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: token.proto

package client

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on TokenGenerateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TokenGenerateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TokenGenerateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TokenGenerateRequestMultiError, or nil if none found.
func (m *TokenGenerateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *TokenGenerateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !_TokenGenerateRequest_Username_Pattern.MatchString(m.GetUsername()) {
		err := TokenGenerateRequestValidationError{
			field:  "Username",
			reason: "value does not match regex pattern \"^[a-zA-Z][a-zA-Z0-9_]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Password

	if len(errors) > 0 {
		return TokenGenerateRequestMultiError(errors)
	}

	return nil
}

// TokenGenerateRequestMultiError is an error wrapping multiple validation
// errors returned by TokenGenerateRequest.ValidateAll() if the designated
// constraints aren't met.
type TokenGenerateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenGenerateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenGenerateRequestMultiError) AllErrors() []error { return m }

// TokenGenerateRequestValidationError is the validation error returned by
// TokenGenerateRequest.Validate if the designated constraints aren't met.
type TokenGenerateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenGenerateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenGenerateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenGenerateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenGenerateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenGenerateRequestValidationError) ErrorName() string {
	return "TokenGenerateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e TokenGenerateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTokenGenerateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenGenerateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenGenerateRequestValidationError{}

var _TokenGenerateRequest_Username_Pattern = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]+$")

// Validate checks the field values on TokenGenerateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TokenGenerateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TokenGenerateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TokenGenerateResponseMultiError, or nil if none found.
func (m *TokenGenerateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *TokenGenerateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	// no validation rules for RefreshToken

	if len(errors) > 0 {
		return TokenGenerateResponseMultiError(errors)
	}

	return nil
}

// TokenGenerateResponseMultiError is an error wrapping multiple validation
// errors returned by TokenGenerateResponse.ValidateAll() if the designated
// constraints aren't met.
type TokenGenerateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenGenerateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenGenerateResponseMultiError) AllErrors() []error { return m }

// TokenGenerateResponseValidationError is the validation error returned by
// TokenGenerateResponse.Validate if the designated constraints aren't met.
type TokenGenerateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenGenerateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenGenerateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenGenerateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenGenerateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenGenerateResponseValidationError) ErrorName() string {
	return "TokenGenerateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e TokenGenerateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTokenGenerateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenGenerateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenGenerateResponseValidationError{}