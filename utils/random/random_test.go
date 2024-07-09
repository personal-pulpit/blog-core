package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateOTP(t *testing.T) {
	otp := GenerateOTP()
	require.NotZero(t, otp)

	t.Log(otp)
}

func TestGenerateUniqueId(t *testing.T) {
	uniqueId := GenerateUniqueId()

	t.Log(uniqueId)
}

func TestGenerateId(t *testing.T) {
	id := GenerateId()

	t.Log(id)
}
