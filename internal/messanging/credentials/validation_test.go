package credentials__test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"testing"
)

func TestCredentialsCreateRequestUsernameValidation_Invalid(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Username = "_invalid_"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestUsernameValidation_TooLong(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Username = gofakeit.LoremIpsumSentence(50)

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestUsernameValidation_TooShort(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Username = "1"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestUsernameValidation_Empty(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Username = ""

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestUsernameValidation_Valid(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()

	err := request.Validate()
	require.Nil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_Empty(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = ""

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_TooLong(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = gofakeit.LoremIpsumSentence(50)

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_TooShort(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = "aQ1!,"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_NoUpper(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = "heavy_metal11!"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_NoLower(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = "HEAVY_METAL11!"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_NoDigit(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = "HEAVY_metal!!"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_NoPunctual(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = "heavyMETAL11"

	err := request.Validate()
	require.NotNil(t, err)
}

func TestCredentialsCreateRequestPasswordValidation_Valid(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	request.Password = "heavy_METAL11!!"

	err := request.Validate()
	require.Nil(t, err)
}
