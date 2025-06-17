package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateOTP(t *testing.T) {
	otp, err := GenerateOTP()
	require.NoError(t, err)
	require.NotZero(t, otp)

	t.Log(otp)
}

func TestGenerateUniqueID(t *testing.T) {
	uniqueID, err := GenerateUniqueID()
	require.NoError(t, err)
	require.NotZero(t, uniqueID)

	t.Log(uniqueID)
}
