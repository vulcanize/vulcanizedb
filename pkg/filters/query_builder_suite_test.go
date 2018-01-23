package filters_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueryBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "QueryBuilder Suite")
}
