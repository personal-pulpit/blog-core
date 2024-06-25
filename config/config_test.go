package config

import (
	
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Get Config", func() {
	Describe("GetConfigInstance", func() {
		Context("when configInstance is nil", func() {
			It("should return a valid Config instance", func() {
				config := GetConfigInstance()
				Expect(config).NotTo(BeNil())
			})
		})
		Context("when checking the fileds of config", func() {
			It("should match a specific string value", func() {
				config := GetConfigInstance()
				Expect(config.Mysql.Username).To(Equal("user"))
				Expect(config.Mysql.Password).To(Equal("password"))
				Expect(config.Mysql.DBName).To(Equal("blog"))
				Expect(config.Mysql.Host).To(Equal("127.0.0.1"))
				Expect(config.Mysql.Port).To(Equal(3306))
				Expect(config.Redis.DB).To(Equal(0))
				Expect(config.Redis.Port).To(Equal(6379))
				Expect(config.Redis.Host).To(Equal("127.0.0.1"))
				Expect(config.Server.Port).To(Equal(8000))
			})
		})
	})
	Describe("getEnv", func() {
		Context("when ENV is set to Development", func() {
			It("should return Development environment", func() {
				os.Setenv("ENV", "Development")
				env := getEnv()
				Expect(env).To(Equal(Development))
			})
		})

		Context("when ENV is set to Production", func() {
			It("should return Production environment", func() {
				os.Setenv("ENV", "Production")
				env := getEnv()
				Expect(env).To(Equal(Production))
			})
		})
		Context("when ENV is set to Test", func() {
			It("should return Test environment", func() {
				os.Setenv("ENV", "Test")
				env := getEnv()
				Expect(env).To(Equal(Test))
			})
		})
		Context("when ENV is set to an invalid value", func() {
			It("should panic with an error message", func() {
				os.Setenv("ENV", "Invalid")
				Expect(func() { getEnv() }).To(Panic())
			})
		})
	})
})
