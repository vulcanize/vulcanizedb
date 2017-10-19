package main_test

import (
	vulcanizedb "github.com/8thlight/vulcanizedb/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vulcanizedb", func() {

	It("is an example test", func() {
		Expect(vulcanizedb.Message()).Should(Equal("Hello world"))
	})

})
