package random_test

import (
	"blog/utils/random"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateOTP(t *testing.T) {
	otp := random.GenerateOTP()
	require.NotZero(t, otp)

	t.Log(otp)
}
