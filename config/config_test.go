package config

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	var config *Config

	BeforeEach(func() {
		config = GetConfigInstance()
	})

	Describe("GetConfigInstance", func() {
		Context("when configInstance is nil", func() {
			It("should return a valid Config instance", func() {
				Expect(config).NotTo(BeNil())
			})
		})
	})

	Describe("Environment", func() {
		Describe("GetEnv", func() {
			It("should return the correct environment", func() {
				Expect(GetEnv()).To(Equal(Development))
			})
		})

		Context("when checking the fields of config-development", func() {
			It("should match specific values", func() {
				if GetEnv() == Development {
					Expect(config.Postgres.Username).To(Equal("user"))
					Expect(config.Postgres.Password).To(Equal("password"))
					Expect(config.Postgres.DBName).To(Equal("blog"))
					Expect(config.Postgres.Host).To(Equal("127.0.0.1"))
					Expect(config.Postgres.Port).To(Equal(5432))
					Expect(config.Redis.DB).To(Equal(0))
					Expect(config.Redis.Port).To(Equal(6379))
					Expect(config.Redis.Host).To(Equal("127.0.0.1"))
				}
			})
		})
	})
})
