package level_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLevel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Level Suite")
}
