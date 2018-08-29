package pit_file_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPitFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PitFile Suite")
}
