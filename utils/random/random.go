package random

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

const OtpLength = 6

func generateRandomNumber(length int) (int, error) {
	min := int64(math.Pow(10, float64(length)-1))
	max := int64(math.Pow(10, float64(length))) - 1

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		return 0, err
	}

	number := int(randomNumber.Int64()) + int(min)

	if len(strconv.Itoa(number)) != length {
		number, err = GenerateOTP()
		if err != nil {
			return 0, err
		}
	}
	return number, nil
}

func GenerateOTP() (int, error) {
	return generateRandomNumber(OtpLength)
}

func GenerateUniqueID() (int, error) {
	return generateRandomNumber(5)
}
