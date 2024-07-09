package random

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

const OtpLength = 6

func generateRandomNumber(length int) int {
	min := int64(math.Pow(10, float64(length)-1))
	max := int64(math.Pow(10, float64(length))) - 1

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	number := int(randomNumber.Int64()) + int(min)

	if len(strconv.Itoa(number)) != length {
		number = GenerateOTP()
	}

	return number
}

func GenerateOTP() int {
	return generateRandomNumber(OtpLength)
}

func GenerateUniqueId() int {
	return generateRandomNumber(5)
}

func GenerateId() int {
	return generateRandomNumber(6)
}
