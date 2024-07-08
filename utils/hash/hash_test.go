package hash

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHashManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hash Manager Suite")
}

var _ = Describe("Hash Manager", func() {
	var (
		plainPassword  string
		hashManager    *HashManager
		hashedPassword string
		err            error
	)

	BeforeEach(func() {
		plainPassword = "password"
		hashManager = NewHashManager(DefaultHashParams)
		hashedPassword, err = hashManager.HashPassword(plainPassword)
		Expect(err).ToNot(HaveOccurred())
		Expect(hashedPassword).ToNot(BeEmpty())
	})

	It("should hash and validate password correctly", func() {
		valid := hashManager.CheckPasswordHash(plainPassword, hashedPassword)
		Expect(valid).To(BeTrue())

		valid = hashManager.CheckPasswordHash("invalid-plain-password", hashedPassword)
		Expect(valid).To(BeFalse())
	})
})
