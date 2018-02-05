package graphql_server_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGraphqlServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GraphqlServer Suite")
}
