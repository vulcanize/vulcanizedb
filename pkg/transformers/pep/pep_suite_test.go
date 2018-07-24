package pep_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPep(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pep Suite")
}
