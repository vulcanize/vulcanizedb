package pip

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pip Suite")
}
