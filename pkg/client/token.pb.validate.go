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

// Validate checks the field values on TokenRefreshRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TokenRefreshRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TokenRefreshRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TokenRefreshRequestMultiError, or nil if none found.
func (m *TokenRefreshRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *TokenRefreshRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for RefreshToken

	if len(errors) > 0 {
		return TokenRefreshRequestMultiError(errors)
	}

	return nil
}

// TokenRefreshRequestMultiError is an error wrapping multiple validation
// errors returned by TokenRefreshRequest.ValidateAll() if the designated
// constraints aren't met.
type TokenRefreshRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenRefreshRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenRefreshRequestMultiError) AllErrors() []error { return m }

// TokenRefreshRequestValidationError is the validation error returned by
// TokenRefreshRequest.Validate if the designated constraints aren't met.
type TokenRefreshRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenRefreshRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenRefreshRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenRefreshRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenRefreshRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenRefreshRequestValidationError) ErrorName() string {
	return "TokenRefreshRequestValidationError"
}

// Error satisfies the builtin error interface
func (e TokenRefreshRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTokenRefreshRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenRefreshRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenRefreshRequestValidationError{}

// Validate checks the field values on TokenRefreshResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TokenRefreshResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TokenRefreshResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TokenRefreshResponseMultiError, or nil if none found.
func (m *TokenRefreshResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *TokenRefreshResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for NewAccessToken

	if len(errors) > 0 {
		return TokenRefreshResponseMultiError(errors)
	}

	return nil
}

// TokenRefreshResponseMultiError is an error wrapping multiple validation
// errors returned by TokenRefreshResponse.ValidateAll() if the designated
// constraints aren't met.
type TokenRefreshResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenRefreshResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenRefreshResponseMultiError) AllErrors() []error { return m }

// TokenRefreshResponseValidationError is the validation error returned by
// TokenRefreshResponse.Validate if the designated constraints aren't met.
type TokenRefreshResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenRefreshResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenRefreshResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenRefreshResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenRefreshResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenRefreshResponseValidationError) ErrorName() string {
	return "TokenRefreshResponseValidationError"
}

// Error satisfies the builtin error interface
func (e TokenRefreshResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTokenRefreshResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenRefreshResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenRefreshResponseValidationError{}

// Validate checks the field values on TokenValidateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TokenValidateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TokenValidateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TokenValidateRequestMultiError, or nil if none found.
func (m *TokenValidateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *TokenValidateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	if len(errors) > 0 {
		return TokenValidateRequestMultiError(errors)
	}

	return nil
}

// TokenValidateRequestMultiError is an error wrapping multiple validation
// errors returned by TokenValidateRequest.ValidateAll() if the designated
// constraints aren't met.
type TokenValidateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenValidateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenValidateRequestMultiError) AllErrors() []error { return m }

// TokenValidateRequestValidationError is the validation error returned by
// TokenValidateRequest.Validate if the designated constraints aren't met.
type TokenValidateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenValidateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenValidateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenValidateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenValidateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenValidateRequestValidationError) ErrorName() string {
	return "TokenValidateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e TokenValidateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTokenValidateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenValidateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenValidateRequestValidationError{}

// Validate checks the field values on TokenValidateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TokenValidateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TokenValidateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TokenValidateResponseMultiError, or nil if none found.
func (m *TokenValidateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *TokenValidateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CredentialsId

	// no validation rules for Username

	// no validation rules for AccountType

	if len(errors) > 0 {
		return TokenValidateResponseMultiError(errors)
	}

	return nil
}

// TokenValidateResponseMultiError is an error wrapping multiple validation
// errors returned by TokenValidateResponse.ValidateAll() if the designated
// constraints aren't met.
type TokenValidateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenValidateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenValidateResponseMultiError) AllErrors() []error { return m }

// TokenValidateResponseValidationError is the validation error returned by
// TokenValidateResponse.Validate if the designated constraints aren't met.
type TokenValidateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenValidateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenValidateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenValidateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenValidateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenValidateResponseValidationError) ErrorName() string {
	return "TokenValidateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e TokenValidateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTokenValidateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenValidateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenValidateResponseValidationError{}
