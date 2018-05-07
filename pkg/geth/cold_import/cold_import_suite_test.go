package cold_import_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestColdImport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ColdImport Suite")
}
