package random

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGenerateRandom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Random Suite")
}

var _ = Describe("generate random numbers", func() {
	Context("Generate a random number", func() {
		var num int

		BeforeEach(func() {
			num = generateRandomNumber(5)
		})

		It("should generate a non-zero number", func() {
			Expect(num).ToNot(Equal(0))
		})
	})
})
