package repositories_test

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/repositories/testing"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("In memory repository", func() {

	testing.AssertRepositoryBehavior(func(core.Node) repositories.Repository {
		return repositories.NewInMemory()
	})

})
