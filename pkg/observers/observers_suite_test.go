package observers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestObservers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Observers Suite")
}
