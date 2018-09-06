package stability_fee_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStabilityFee(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StabilityFee Suite")
}
