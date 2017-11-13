package repositories_test

import (
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	"github.com/8thlight/vulcanizedb/pkg/repositories/testing"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("In memory repository", func() {

	testing.AssertRepositoryBehavior(func() repositories.Repository {
		return repositories.NewInMemory()
	})

})
