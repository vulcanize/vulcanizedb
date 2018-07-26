package rep_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRep(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rep Suite")
}
