package geth_test

import (
	"path/filepath"

	"net/http"

	"fmt"

	"log"

	cfg "github.com/8thlight/vulcanizedb/pkg/config"
	"github.com/8thlight/vulcanizedb/pkg/geth"
	"github.com/ethereum/go-ethereum/accounts/abi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ABI files", func() {

	Describe("Reading ABI files", func() {

		It("loads a valid ABI file", func() {
			path := filepath.Join(cfg.ProjectRoot(), "pkg", "geth", "testing", "valid_abi.json")

			contractAbi, err := geth.ParseAbiFile(path)

			Expect(contractAbi).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("reads the contents of a valid ABI file", func() {
			path := filepath.Join(cfg.ProjectRoot(), "pkg", "geth", "testing", "valid_abi.json")

			contractAbi, err := geth.ReadAbiFile(path)

			Expect(contractAbi).To(Equal("[{\"foo\": \"bar\"}]"))
			Expect(err).To(BeNil())
		})

		It("returns an error when the file does not exist", func() {
			path := filepath.Join(cfg.ProjectRoot(), "pkg", "geth", "testing", "missing_abi.json")

			contractAbi, err := geth.ParseAbiFile(path)

			Expect(contractAbi).To(Equal(abi.ABI{}))
			Expect(err).To(Equal(geth.ErrMissingAbiFile))
		})

		It("returns an error when the file has invalid contents", func() {
			path := filepath.Join(cfg.ProjectRoot(), "pkg", "geth", "testing", "invalid_abi.json")

			contractAbi, err := geth.ParseAbiFile(path)

			Expect(contractAbi).To(Equal(abi.ABI{}))
			Expect(err).To(Equal(geth.ErrInvalidAbiFile))
		})

		Describe("Request ABI from endpoint", func() {

			var (
				server    *ghttp.Server
				client    *geth.EtherScanApi
				abiString string
			)

			BeforeEach(func() {
				server = ghttp.NewServer()
				client = geth.NewEtherScanClient(server.URL())
				path := filepath.Join(cfg.ProjectRoot(), "pkg", "geth", "testing", "sample_abi.json")
				abiString, err := geth.ReadAbiFile(path)
				_, err = geth.ParseAbi(abiString)
				if err != nil {
					log.Fatalln("Could not parse ABI")
				}
			})

			AfterEach(func() {
				server.Close()
			})

			Describe("Fetching ABI from api (etherscan)", func() {
				BeforeEach(func() {

					response := fmt.Sprintf(`{"status":"1","message":"OK","result":%q}`, abiString)
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/api", "module=contract&action=getabi&address=0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"),
							ghttp.RespondWith(http.StatusOK, response),
						),
					)
				})

				It("should make a GET request with supplied contract hash", func() {

					abi, err := client.GetAbi("0xd26114cd6EE289AccF82350c8d8487fedB8A0C07")
					Expect(server.ReceivedRequests()).Should(HaveLen(1))
					Expect(err).ShouldNot(HaveOccurred())
					Expect(abi).Should(Equal(abiString))
				})
			})
		})

		Describe("Generating etherscan endpoints based on network", func() {
			It("should return the main endpoint as the default", func() {
				url := geth.GenUrl("")
				Expect(url).To(Equal("https://api.etherscan.io"))
			})

			It("generates various test network endpoint if test network is supplied", func() {
				ropstenUrl := geth.GenUrl("ropsten")
				rinkebyUrl := geth.GenUrl("rinkeby")
				kovanUrl := geth.GenUrl("kovan")

				Expect(ropstenUrl).To(Equal("https://ropsten.etherscan.io"))
				Expect(kovanUrl).To(Equal("https://kovan.etherscan.io"))
				Expect(rinkebyUrl).To(Equal("https://rinkeby.etherscan.io"))
			})
		})
	})
})
