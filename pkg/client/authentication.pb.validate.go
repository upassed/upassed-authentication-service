// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: authentication.proto

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

// Validate checks the field values on CredentialsCreateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CredentialsCreateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CredentialsCreateRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CredentialsCreateRequestMultiError, or nil if none found.
func (m *CredentialsCreateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CredentialsCreateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 5 || l > 20 {
		err := CredentialsCreateRequestValidationError{
			field:  "Username",
			reason: "value length must be between 5 and 20 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 5 || l > 50 {
		err := CredentialsCreateRequestValidationError{
			field:  "Password",
			reason: "value length must be between 5 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CredentialsCreateRequestMultiError(errors)
	}

	return nil
}

// CredentialsCreateRequestMultiError is an error wrapping multiple validation
// errors returned by CredentialsCreateRequest.ValidateAll() if the designated
// constraints aren't met.
type CredentialsCreateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CredentialsCreateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CredentialsCreateRequestMultiError) AllErrors() []error { return m }

// CredentialsCreateRequestValidationError is the validation error returned by
// CredentialsCreateRequest.Validate if the designated constraints aren't met.
type CredentialsCreateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CredentialsCreateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CredentialsCreateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CredentialsCreateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CredentialsCreateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CredentialsCreateRequestValidationError) ErrorName() string {
	return "CredentialsCreateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CredentialsCreateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCredentialsCreateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CredentialsCreateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CredentialsCreateRequestValidationError{}

// Validate checks the field values on CredentialsCreateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CredentialsCreateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CredentialsCreateResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CredentialsCreateResponseMultiError, or nil if none found.
func (m *CredentialsCreateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CredentialsCreateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CredentialsId

	if len(errors) > 0 {
		return CredentialsCreateResponseMultiError(errors)
	}

	return nil
}

// CredentialsCreateResponseMultiError is an error wrapping multiple validation
// errors returned by CredentialsCreateResponse.ValidateAll() if the
// designated constraints aren't met.
type CredentialsCreateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CredentialsCreateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CredentialsCreateResponseMultiError) AllErrors() []error { return m }

// CredentialsCreateResponseValidationError is the validation error returned by
// CredentialsCreateResponse.Validate if the designated constraints aren't met.
type CredentialsCreateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CredentialsCreateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CredentialsCreateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CredentialsCreateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CredentialsCreateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CredentialsCreateResponseValidationError) ErrorName() string {
	return "CredentialsCreateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CredentialsCreateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCredentialsCreateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CredentialsCreateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CredentialsCreateResponseValidationError{}
