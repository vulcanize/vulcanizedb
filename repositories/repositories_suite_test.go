package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repositories Suite")
}
